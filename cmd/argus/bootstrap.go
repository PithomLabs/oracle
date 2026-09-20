package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/PithomLabs/oracle/internal/migrations"
	"github.com/PithomLabs/oracle/internal/solventmigrations"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	defaultDBURL  = "postgres://root@localhost:26257/argus?sslmode=disable"
	adminDBURL    = "postgres://root@localhost:26257/defaultdb?sslmode=disable"
	crdbSQLPort   = "26257"
	crdbAdminPort = "8081"
	crdbHost      = "localhost"
	waitTimeout   = 30 * time.Second
)

// resolveDBURL returns the effective database URL.
// Priority: explicit flag > ARGUS_DB_URL env > default.
func resolveDBURL(flagValue string, flagWasSet bool) string {
	if flagWasSet && flagValue != "" {
		return flagValue
	}
	if env := os.Getenv("ARGUS_DB_URL"); env != "" {
		return env
	}
	return defaultDBURL
}

// isLocalURL returns true if the URL targets the standard local CRDB.
func isLocalURL(dbURL string) bool {
	u, err := url.Parse(dbURL)
	if err != nil {
		return false
	}
	host := u.Hostname()
	return host == "localhost" || host == "127.0.0.1"
}

// findCockroachBinary locates the cockroach executable.
func findCockroachBinary() (string, error) {
	if path, err := exec.LookPath("cockroach"); err == nil {
		return path, nil
	}
	home, _ := os.UserHomeDir()
	candidates := []string{
		filepath.Join(home, ".local", "bin", "cockroach"),
		"/usr/local/bin/cockroach",
		"/opt/cockroach/cockroach",
	}
	if runtime.GOOS == "windows" {
		candidates = append(candidates,
			filepath.Join(home, "cockroach", "cockroach.exe"),
			"C:\\cockroach\\cockroach.exe",
		)
	}
	for _, c := range candidates {
		if info, err := os.Stat(c); err == nil && info.Mode().IsRegular() {
			return c, nil
		}
	}
	return "", fmt.Errorf("ARGUS: cockroach executable not found in PATH.\nInstall CockroachDB or provide --db/ARGUS_DB_URL")
}

// isCockroachRunning checks if CRDB is accepting TCP connections on the given host:port.
func isCockroachRunning(host string) bool {
	conn, err := net.DialTimeout("tcp", host+":"+crdbSQLPort, 500*time.Millisecond)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// startLocalCockroach launches a single-node CRDB in the background.
// It writes the PID to .argus-pids/crdb.pid and returns the process.
func startLocalCockroach() (*exec.Cmd, error) {
	bin, err := findCockroachBinary()
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(".argus-pids", 0755); err != nil {
		return nil, fmt.Errorf("create pid directory: %w", err)
	}

	cmd := exec.Command(bin, "start-single-node", "--insecure",
		"--listen-addr", ":"+crdbSQLPort,
		"--http-addr", ":"+crdbAdminPort,
		"--store=.cockroach-data",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start CockroachDB: %w", err)
	}

	pidFile := filepath.Join(".argus-pids", "crdb.pid")
	if err := os.WriteFile(pidFile, []byte(fmt.Sprintf("%d", cmd.Process.Pid)), 0644); err != nil {
		log.Printf("warning: could not write CRDB pid file: %v", err)
	}

	log.Printf("ARGUS: started CockroachDB (PID %d) on :%s", cmd.Process.Pid, crdbSQLPort)
	return cmd, nil
}

// waitForSQL polls until CRDB accepts SQL connections or timeout.
func waitForSQL(host string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	dsn := "postgres://root@" + host + ":" + crdbSQLPort + "/defaultdb?sslmode=disable"

	for time.Now().Before(deadline) {
		db, err := sql.Open("pgx", dsn)
		if err == nil {
			err = db.Ping()
			db.Close()
			if err == nil {
				return nil
			}
		}
		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("CockroachDB did not become ready within %v", timeout)
}

// ensureDatabase creates the target database if it does not exist.
func ensureDatabase(dbURL string) error {
	dbName := extractDBName(dbURL)
	if dbName == "" {
		dbName = "argus"
	}

	adminDB, err := sql.Open("pgx", adminDBURL)
	if err != nil {
		return fmt.Errorf("connect to admin database: %w", err)
	}
	defer adminDB.Close()

	if err := adminDB.Ping(); err != nil {
		return fmt.Errorf("ping admin database: %w", err)
	}

	_, err = adminDB.ExecContext(context.Background(),
		fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %q", dbName))
	if err != nil {
		return fmt.Errorf("create database %q: %w", dbName, err)
	}

	log.Printf("ARGUS: database %q ready", dbName)
	return nil
}

// applyMigrations runs all Solvent and ARGUS migrations (idempotent).
func applyMigrations(dbURL string) error {
	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		return fmt.Errorf("open db for migrations: %w", err)
	}
	defer db.Close()

	ctx := context.Background()

	if err := solventmigrations.Apply(ctx, db); err != nil {
		return fmt.Errorf("solvent migrations: %w", err)
	}
	if err := migrations.Apply(ctx, db); err != nil {
		return fmt.Errorf("argus migrations: %w", err)
	}

	log.Println("ARGUS: migrations applied")
	return nil
}

// bootstrap resolves the DB URL, ensures CRDB is running, creates the
// database if needed, and applies migrations. Returns the resolved URL
// and a non-nil *exec.Cmd only if ARGUS started its own CRDB process
// (caller must clean it up).
func bootstrap(dbURLFlag string) (string, *exec.Cmd, error) {
	resolvedURL := resolveDBURL(dbURLFlag, dbURLFlag != "")

	// Step 1: Check if CRDB is reachable
	crdbRunning := isCockroachRunning(crdbHost)
	var managedCRDB *exec.Cmd

	// Step 2: If not reachable and this is the local URL, start CRDB
	if !crdbRunning {
		if !isLocalURL(resolvedURL) {
			return "", nil, fmt.Errorf(
				"cannot connect to database at %s.\nStart CockroachDB or check the URL.", resolvedURL)
		}

		var err error
		managedCRDB, err = startLocalCockroach()
		if err != nil {
			return "", nil, err
		}

		if err := waitForSQL(crdbHost, waitTimeout); err != nil {
			managedCRDB.Process.Signal(os.Kill)
			managedCRDB.Wait()
			os.Remove(filepath.Join(".argus-pids", "crdb.pid"))
			return "", nil, err
		}
	} else {
		log.Println("ARGUS: existing CockroachDB detected on :"+crdbSQLPort+", reusing")
	}

	// Step 3: Always ensure database exists (whether we started CRDB or found it)
	if err := ensureDatabase(resolvedURL); err != nil {
		if managedCRDB != nil {
			managedCRDB.Process.Signal(os.Kill)
			managedCRDB.Wait()
		}
		return "", nil, err
	}

	// Step 4: Always apply migrations (whether we started CRDB or found it)
	if err := applyMigrations(resolvedURL); err != nil {
		if managedCRDB != nil {
			managedCRDB.Process.Signal(os.Kill)
			managedCRDB.Wait()
		}
		return "", nil, err
	}

	return resolvedURL, managedCRDB, nil
}

// extractDBName parses the database name from a PostgreSQL URL.
func extractDBName(dbURL string) string {
	u, err := url.Parse(dbURL)
	if err != nil {
		return ""
	}
	name := strings.TrimPrefix(u.Path, "/")
	if name == "" {
		return "defaultdb"
	}
	return name
}
