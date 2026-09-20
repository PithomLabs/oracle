package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/PithomLabs/oracle/domain-pack"
bmistv1 "github.com/PithomLabs/oracle/domain-pack/bmist/v1"
	bmistv11 "github.com/PithomLabs/oracle/domain-pack/bmist/v1.1.0"
	"github.com/PithomLabs/oracle/internal/application"
	"github.com/PithomLabs/oracle/internal/mcp"
	"github.com/PithomLabs/oracle/internal/migrations"
	"github.com/PithomLabs/oracle/internal/ui"
	"github.com/PithomLabs/oracle/seed"
	"github.com/PithomLabs/oracle/internal/solventmigrations"
	"github.com/PithomLabs/oracle/verifier"
	"github.com/PithomLabs/oracle/verifier/physics/v1"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: argus <command> [flags]")
		fmt.Fprintln(os.Stderr, "commands: serve, mcp, verify, migrate, reset")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "serve":
		cmdServe(os.Args[2:])
	case "mcp":
		cmdMCP(os.Args[2:])
	case "verify":
		cmdVerify(os.Args[2:])
	case "migrate":
		cmdMigrate(os.Args[2:])
	case "reset":
		cmdReset(os.Args[2:])
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func cmdServe(args []string) {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	dbURL := fs.String("db", "", "database URL (default: local CRDB)")
	listen := fs.String("listen", ":8080", "HTTP listen address")
	fs.Parse(args)

	resolvedURL, managedCRDB, err := bootstrap(*dbURL)
	if err != nil {
		log.Fatalf("%v", err)
	}
	if managedCRDB != nil {
		defer func() {
			log.Println("stopping managed CockroachDB...")
			managedCRDB.Process.Signal(os.Kill)
			done := make(chan struct{})
			go func() {
				managedCRDB.Wait()
				close(done)
			}()
			select {
			case <-done:
			case <-time.After(5 * time.Second):
				log.Println("warning: CRDB did not exit in time")
			}
			os.Remove(".argus-pids/crdb.pid")
		}()
	}

	db, err := sql.Open("pgx", resolvedURL)
	if err != nil {
		log.Fatalf("ARGUS: open db: %v", err)
	}
	defer db.Close()

	registry := verifier.NewArtifactRegistry()
	specs := []application.VerifierSpec{
		{VerifierID: "physics-v1", MinVersion: "0.1.0"},
	}
	app := application.NewWithRegistry(db, registry, specs)

	// Load domain packs (embedded at compile time)
	packReg := domainpack.NewRegistry()
	packReg.RegisterDecoder("bmist", func(data []byte) (domainpack.Pack, error) {
		return bmistv1.ParsePack(data)
	})
	if err := packReg.LoadEmbedded(bmistv1.LoadEmbedded(), "bmist"); err != nil {
		log.Fatalf("embedded pack load failed: %v", err)
	}
	if err := packReg.LoadEmbedded(bmistv11.LoadEmbedded(), "bmist"); err != nil {
		log.Fatalf("embedded v1.1.0 pack load failed: %v", err)
	}
	app.SetPackRegistry(packReg)

	// Seed the BM-IST-AS POC if DB is empty.
	if err := seed.SeedIfEmpty(db); err != nil {
		log.Fatalf("seed: %v", err)
	}

	operatorToken := os.Getenv("ARGUS_OPERATOR_TOKEN")
	if operatorToken == "" {
		operatorToken = "argus-local-operator"
		log.Printf("ARGUS_OPERATOR_TOKEN not set, using default local operator token")
	}
	uiServer := ui.NewServer(app, operatorToken)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"status":"ok"}`)
	})
	uiServer.RegisterRoutes(mux)

	srv := &http.Server{
		Addr:         *listen,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		log.Println("shutting down...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		srv.Shutdown(shutdownCtx)
	}()

	log.Printf("argus serve listening on %s", *listen)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("serve: %v", err)
	}
}

func cmdMCP(args []string) {
	fs := flag.NewFlagSet("mcp", flag.ExitOnError)
	dbURL := fs.String("db", "", "database URL (default: local CRDB)")
	fs.Parse(args)

	resolvedURL, managedCRDB, err := bootstrap(*dbURL)
	if err != nil {
		log.Fatalf("%v", err)
	}
	if managedCRDB != nil {
		defer func() {
			managedCRDB.Process.Signal(os.Kill)
			done := make(chan struct{})
			go func() {
				managedCRDB.Wait()
				close(done)
			}()
			select {
			case <-done:
			case <-time.After(5 * time.Second):
			}
			os.Remove(".argus-pids/crdb.pid")
		}()
	}

	db, err := sql.Open("pgx", resolvedURL)
	if err != nil {
		log.Fatalf("ARGUS: open db: %v", err)
	}
	defer db.Close()

	app := application.New(db)

	packReg := domainpack.NewRegistry()
	packReg.RegisterDecoder("bmist", func(data []byte) (domainpack.Pack, error) {
		return bmistv1.ParsePack(data)
	})
	if err := packReg.LoadEmbedded(bmistv1.LoadEmbedded(), "bmist"); err != nil {
		log.Fatalf("embedded pack load failed: %v", err)
	}
	if err := packReg.LoadEmbedded(bmistv11.LoadEmbedded(), "bmist"); err != nil {
		log.Fatalf("embedded v1.1.0 pack load failed: %v", err)
	}
	app.SetPackRegistry(packReg)

	adapter := mcp.NewAdapter(app)
	fmt.Fprintln(os.Stderr, "argus mcp: starting MCP stdio adapter")
	runMCPStdio(adapter)
}

func cmdVerify(args []string) {
	fs := flag.NewFlagSet("verify", flag.ExitOnError)
	_ = fs.String("db", "postgres://root@localhost:26257/argus?sslmode=disable", "database URL")
	inputFile := fs.String("input", "", "verification input file (JSON)")
	timeout := fs.Duration("timeout", 30*time.Second, "verification timeout")
	fs.Parse(args)

	if *inputFile == "" {
		fmt.Fprintln(os.Stderr, "usage: argus verify --input <file> [--timeout 30s]")
		os.Exit(1)
	}

	data, err := os.ReadFile(*inputFile)
	if err != nil {
		log.Fatalf("read input: %v", err)
	}

	var input physicsv1.VerifierInput
	if err := json.Unmarshal(data, &input); err != nil {
		log.Fatalf("parse input: %v", err)
	}

	registry := verifier.NewArtifactRegistry()

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	artifact, err := verifier.RunPhysicsVerifier(ctx, registry, input)
	if err != nil {
		log.Fatalf("verify: %v", err)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(artifact); err != nil {
		log.Fatalf("encode artifact: %v", err)
	}
}

func cmdMigrate(args []string) {
	// Non-destructive: applies Solvent + ARGUS migrations idempotently.
	fs := flag.NewFlagSet("migrate", flag.ExitOnError)
	dbURL := fs.String("db", "", "database URL (default: local CRDB)")
	fs.Parse(args)

	resolvedURL, managedCRDB, err := bootstrap(*dbURL)
	if err != nil {
		log.Fatalf("%v", err)
	}
	if managedCRDB != nil {
		defer func() {
			managedCRDB.Process.Signal(os.Kill)
			done := make(chan struct{})
			go func() {
				managedCRDB.Wait()
				close(done)
			}()
			select {
			case <-done:
			case <-time.After(5 * time.Second):
			}
			os.Remove(".argus-pids/crdb.pid")
		}()
	}

	_ = resolvedURL
	log.Println("argus migrate: done (applied during bootstrap)")
}

func cmdReset(args []string) {
	// Destructive: drops all tables and re-applies.
	fs := flag.NewFlagSet("reset", flag.ExitOnError)
	dbURL := fs.String("db", "", "database URL (default: local CRDB)")
	fs.Parse(args)

	resolvedURL, managedCRDB, err := bootstrap(*dbURL)
	if err != nil {
		log.Fatalf("%v", err)
	}
	if managedCRDB != nil {
		defer func() {
			managedCRDB.Process.Signal(os.Kill)
			done := make(chan struct{})
			go func() {
				managedCRDB.Wait()
				close(done)
			}()
			select {
			case <-done:
			case <-time.After(5 * time.Second):
			}
			os.Remove(".argus-pids/crdb.pid")
		}()
	}

	db, err := sql.Open("pgx", resolvedURL)
	if err != nil {
		log.Fatalf("ARGUS: open db: %v", err)
	}
	defer db.Close()

	ctx := context.Background()

	log.Println("argus reset: dropping all tables...")
	tables := []string{
		"conductor_dependency", "conductor_task",
		"submission_idempotency", "belief_retirement_proposal",
		"evidence", "belief_edge", "belief",
	}
	for _, t := range tables {
		if _, err := db.ExecContext(ctx, "DROP TABLE IF EXISTS "+t+" CASCADE"); err != nil {
			log.Fatalf("drop %s: %v", t, err)
		}
	}

	log.Println("argus reset: applying Solvent migrations...")
	if err := solventmigrations.Apply(ctx, db); err != nil {
		log.Fatalf("solvent migrations: %v", err)
	}

	log.Println("argus reset: applying ARGUS migrations...")
	if err := migrations.Apply(ctx, db); err != nil {
		log.Fatalf("argus migrations: %v", err)
	}

	log.Println("argus reset: done")
}
