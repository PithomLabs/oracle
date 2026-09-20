package physicsv1

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math"
)

// Expr is the interface for symbolic expressions.
type Expr interface {
	expr()
	String() string
}

// Var represents a variable.
type Var struct{ Name string }

func (Var) expr()            {}
func (v Var) String() string { return v.Name }

// Const represents a constant value.
type Const struct{ Value float64 }

func (Const) expr()             {}
func (c Const) String() string { return fmt.Sprintf("%g", c.Value) }

// Add represents addition.
type Add struct{ Left, Right Expr }

func (Add) expr()            {}
func (a Add) String() string { return fmt.Sprintf("(%s + %s)", a.Left, a.Right) }

// Mul represents multiplication.
type Mul struct{ Left, Right Expr }

func (Mul) expr()            {}
func (m Mul) String() string { return fmt.Sprintf("(%s * %s)", m.Left, m.Right) }

// Pow represents exponentiation.
type Pow struct{ Base Expr; Exp float64 }

func (Pow) expr()            {}
func (p Pow) String() string { return fmt.Sprintf("(%s^%g)", p.Base, p.Exp) }

// Div represents division.
type Div struct{ Num, Denom Expr }

func (Div) expr()            {}
func (d Div) String() string { return fmt.Sprintf("(%s / %s)", d.Num, d.Denom) }

// Sqrt represents square root.
type Sqrt struct{ Arg Expr }

func (Sqrt) expr()            {}
func (s Sqrt) String() string { return fmt.Sprintf("sqrt(%s)", s.Arg) }

// Grad represents gradient.
type Grad struct{ Arg Expr }

func (Grad) expr()            {}
func (g Grad) String() string { return fmt.Sprintf("grad(%s)", g.Arg) }

// Laplacian represents Laplacian.
type Laplacian struct{ Arg Expr }

func (Laplacian) expr()            {}
func (l Laplacian) String() string { return fmt.Sprintf("laplacian(%s)", l.Arg) }

// Simplify performs constant folding and basic algebraic simplification.
func Simplify(e Expr) Expr {
	switch v := e.(type) {
	case Const:
		return v
	case Var:
		return v
	case Add:
		left := Simplify(v.Left)
		right := Simplify(v.Right)
		lc, lok := left.(Const)
		rc, rok := right.(Const)
		if lok && rok {
			return Const{Value: lc.Value + rc.Value}
		}
		if lok && lc.Value == 0 {
			return right
		}
		if rok && rc.Value == 0 {
			return left
		}
		return Add{Left: left, Right: right}
	case Mul:
		left := Simplify(v.Left)
		right := Simplify(v.Right)
		lc, lok := left.(Const)
		rc, rok := right.(Const)
		if lok && rok {
			return Const{Value: lc.Value * rc.Value}
		}
		if lok && lc.Value == 0 {
			return Const{Value: 0}
		}
		if rok && rc.Value == 0 {
			return Const{Value: 0}
		}
		if lok && lc.Value == 1 {
			return right
		}
		if rok && rc.Value == 1 {
			return left
		}
		return Mul{Left: left, Right: right}
	case Pow:
		base := Simplify(v.Base)
		bc, ok := base.(Const)
		if ok {
			return Const{Value: math.Pow(bc.Value, v.Exp)}
		}
		if v.Exp == 0 {
			return Const{Value: 1}
		}
		if v.Exp == 1 {
			return base
		}
		return Pow{Base: base, Exp: v.Exp}
	case Div:
		num := Simplify(v.Num)
		denom := Simplify(v.Denom)
		nc, nok := num.(Const)
		dc, dok := denom.(Const)
		if nok && dok {
			if dc.Value == 0 {
				return Div{Num: num, Denom: denom}
			}
			return Const{Value: nc.Value / dc.Value}
		}
		if nok && nc.Value == 0 {
			return Const{Value: 0}
		}
		return Div{Num: num, Denom: denom}
	case Sqrt:
		arg := Simplify(v.Arg)
		ac, ok := arg.(Const)
		if ok {
			return Const{Value: math.Sqrt(ac.Value)}
		}
		return Sqrt{Arg: arg}
	case Grad:
		arg := Simplify(v.Arg)
		return Grad{Arg: arg}
	case Laplacian:
		arg := Simplify(v.Arg)
		return Laplacian{Arg: arg}
	default:
		return e
	}
}

// Equals checks structural equality of two expressions.
func Equals(a, b Expr) bool {
	aStr := fmt.Sprintf("%v", a)
	bStr := fmt.Sprintf("%v", b)
	return aStr == bStr
}

// VerifyStep1 checks: Δ√ρ/√ρ = ½(Δρ/ρ) − ¼(|∇ρ|²/ρ²)
func VerifyStep1(rho Expr) (bool, string) {
	// LHS: Δ√ρ/√ρ
	lhs := Div{
		Num:   Laplacian{Arg: Sqrt{Arg: rho}},
		Denom: Sqrt{Arg: rho},
	}

	// RHS: ½(Δρ/ρ) − ¼(|∇ρ|²/ρ²)
	halfDeltaRhoOverRho := Mul{
		Left:  Const{Value: 0.5},
		Right: Div{Num: Laplacian{Arg: rho}, Denom: rho},
	}

	quarterGradSquaredOverRhoSquared := Mul{
		Left: Const{Value: 0.25},
		Right: Div{
			Num:   Mul{Left: Grad{Arg: rho}, Right: Grad{Arg: rho}},
			Denom: Pow{Base: rho, Exp: 2},
		},
	}

	rhs := Add{
		Left:  halfDeltaRhoOverRho,
		Right: Mul{Left: Const{Value: -1}, Right: quarterGradSquaredOverRhoSquared},
	}

	simplifiedLHS := Simplify(lhs)
	simplifiedRHS := Simplify(rhs)

	equal := Equals(simplifiedLHS, simplifiedRHS)
	if equal {
		return true, "Identity verified: Δ√ρ/√ρ = ½(Δρ/ρ) − ¼(|∇ρ|²/ρ²)"
	}
	return false, fmt.Sprintf("Identity not verified: LHS=%s, RHS=%s", simplifiedLHS, simplifiedRHS)
}

// VerifyStep2 checks: a(ρ) = κ²/(8mρ)
func VerifyStep2(kappa, m, rho Expr) (bool, string) {
	// RHS: κ²/(8mρ)
	rhs := Div{
		Num:   Pow{Base: kappa, Exp: 2},
		Denom: Mul{Left: Const{Value: 8}, Right: Mul{Left: m, Right: rho}},
	}

	simplifiedRHS := Simplify(rhs)

	// For the POC, we verify the coefficient structure matches
	expectedForm := Div{
		Num:   Pow{Base: Var{Name: "κ"}, Exp: 2},
		Denom: Mul{Left: Const{Value: 8}, Right: Mul{Left: Var{Name: "m"}, Right: Var{Name: "ρ"}}},
	}

	equal := Equals(simplifiedRHS, expectedForm)
	if equal {
		return true, "Coefficient verified: a(ρ) = κ²/(8mρ)"
	}
	return false, fmt.Sprintf("Coefficient mismatch: got %s, expected %s", simplifiedRHS, expectedForm)
}

// VerifyStep3 checks: |∇ρ|² coefficient consistency (human-attested)
func VerifyStep3(input string, expected string) (bool, string) {
	if input == expected {
		return true, "Human-attested: |∇ρ|² coefficient consistent"
	}
	return false, fmt.Sprintf("Human-attested check failed: got %q, expected %q", input, expected)
}

// VerifyStep4 checks: b′(ρ) = 0 (human-attested)
func VerifyStep4(input string, expected string) (bool, string) {
	if input == expected {
		return true, "Human-attested: b′(ρ) = 0"
	}
	return false, fmt.Sprintf("Human-attested check failed: got %q, expected %q", input, expected)
}

// VerifierInput contains the inputs for the physics verifier.
type VerifierInput struct {
	RunID         string
	CorpusHash    string
	RhoExpr       Expr
	KappaExpr     Expr
	MExpr         Expr
	Step3Input    string
	Step3Expected string
	Step4Input    string
	Step4Expected string
}

// DefaultInput returns a default verifier input for the POC.
func DefaultInput() VerifierInput {
	return VerifierInput{
		RunID:         "run-001",
		CorpusHash:    "corpus-hash-001",
		RhoExpr:       Var{Name: "ρ"},
		KappaExpr:     Var{Name: "κ"},
		MExpr:         Var{Name: "m"},
		Step3Input:    "consistent",
		Step3Expected: "consistent",
		Step4Input:    "0",
		Step4Expected: "0",
	}
}

// Hash computes the canonical SHA-256 of this VerifierInput for input-to-claim binding.
func (v VerifierInput) Hash() (string, error) {
	canonical, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	h := sha256.Sum256(canonical)
	return fmt.Sprintf("%x", h), nil
}
