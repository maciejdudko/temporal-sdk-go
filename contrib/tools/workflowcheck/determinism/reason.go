package determinism

import (
	"encoding/gob"
	"fmt"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// PackageNonDeterminisms contains func/var non-determinisms keyed by name.
type PackageNonDeterminisms map[string]NonDeterminisms

// AFact is for implementing golang.org/x/tools/go/analysis.Fact.
func (*PackageNonDeterminisms) AFact() {}

func (n *PackageNonDeterminisms) String() string {
	if n == nil || len(*n) == 0 {
		return "0 non-deterministic vars/funcs"
	} else if len(*n) == 1 {
		return "1 non-deterministic var/func"
	}
	return strconv.Itoa(len(*n)) + " non-deterministic vars/funcs"
}

// NonDeterminisms is a set of reasons why a function/var is non-deterministic.
type NonDeterminisms []Reason

// AFact is for implementing golang.org/x/tools/go/analysis.Fact.
func (*NonDeterminisms) AFact() {}

// String returns all reasons as a comma-delimited string.
func (n *NonDeterminisms) String() string {
	if n == nil {
		return "<none>"
	}
	var str string
	for _, reason := range *n {
		if str != "" {
			str += ", "
		}
		str += reason.String()
	}
	return str
}

// AppendChildReasonLines appends to lines the set of reasons in this slice.
// This will include newlines and indention based on depth.
func (n NonDeterminisms) AppendChildReasonLines(
	subject string,
	s []string,
	depth int,
	depthRepeat string,
	includePos bool,
	pkg *types.Package,
	lookupCache *PackageLookupCache,
	seenPos map[string]bool,
) []string {
	for _, reason := range n {
		reasonStr := reason.String()
		// Relativize path if it at least starts with working dir
		pos := reason.Pos()
		filename := pos.Filename
		if wd, err := os.Getwd(); err == nil && strings.HasPrefix(filename, wd) {
			if relFilename, err := filepath.Rel(wd, filename); err == nil {
				filename = relFilename
			}
		}
		posStr := fmt.Sprintf("%v:%v:%v", filename, pos.Line, pos.Column)
		if includePos {
			reasonStr += " at " + posStr
		}
		s = append(s, fmt.Sprintf("%v is non-deterministic, reason: %v",
			strings.Repeat(depthRepeat, depth)+subject, reasonStr))
		// Recurse if func call and we haven't seen this pos str before
		seen := seenPos[posStr]
		seenPos[posStr] = true
		if funcCall, _ := reason.(*ReasonFuncCall); funcCall != nil && !seen {
			childPkg, childPkgNonDet := lookupCache.PackageNonDeterminismsFromName(pkg, funcCall.PackageName())
			if childNonDet := childPkgNonDet[funcCall.FuncName]; len(childNonDet) > 0 {
				s = childNonDet.AppendChildReasonLines(funcCall.FuncName, s, depth+1,
					depthRepeat, includePos, childPkg, lookupCache, seenPos)
			}
		}
	}
	return s
}

// Reason represents a reason for non-determinism.
type Reason interface {
	Pos() *token.Position
	// String is expected to just include the brief reason, not any child reasons.
	String() string
}

// ReasonDecl represents a function or var that was explicitly marked
// non-deterministic via config.
type ReasonDecl struct {
	SourcePos *token.Position
}

// Pos returns the source position.
func (r *ReasonDecl) Pos() *token.Position { return r.SourcePos }

// String returns the reason.
func (r *ReasonDecl) String() string {
	return "declared non-deterministic"
}

// ReasonFuncCall represents a call to a non-deterministic function.
type ReasonFuncCall struct {
	SourcePos *token.Position
	// Fully qualified name
	FuncName string
}

// Pos returns the source position.
func (r *ReasonFuncCall) Pos() *token.Position { return r.SourcePos }

// String returns the reason.
func (r *ReasonFuncCall) String() string {
	return "calls non-deterministic function " + r.FuncName
}

func (r *ReasonFuncCall) PackageName() string {
	pkgPrefixedName := r.FuncName
	// If there is a ending paren it's a method, take the receiver as the name
	if endParen := strings.Index(r.FuncName, ")"); endParen >= 0 {
		pkgPrefixedName = strings.TrimLeft(r.FuncName[:endParen], "(*")
	}
	// Take up until the last dot as the package name
	lastDot := strings.LastIndex(pkgPrefixedName, ".")
	if lastDot == -1 {
		return pkgPrefixedName
	}
	return pkgPrefixedName[:lastDot]
}

// ReasonVarAccess represents accessing a non-deterministic global variable.
type ReasonVarAccess struct {
	SourcePos *token.Position
	// Fully qualified name
	VarName string
}

// Pos returns the source position.
func (r *ReasonVarAccess) Pos() *token.Position { return r.SourcePos }

// String returns the reason.
func (r *ReasonVarAccess) String() string {
	return "accesses non-deterministic var " + r.VarName
}

// ReasonConcurrency represents a non-deterministic concurrency construct.
type ReasonConcurrency struct {
	SourcePos *token.Position
	Kind      ConcurrencyKind
}

// Pos returns the source position.
func (r *ReasonConcurrency) Pos() *token.Position { return r.SourcePos }

// String returns the reason.
func (r *ReasonConcurrency) String() string {
	switch r.Kind {
	case ConcurrencyKindGo:
		return "starts goroutine"
	case ConcurrencyKindRecv:
		return "receives from channel"
	case ConcurrencyKindSend:
		return "sends to channel"
	case ConcurrencyKindRange:
		return "iterates over channel"
	default:
		return "<unknown-kind>"
	}
}

// ConcurrencyKind is a construct that is non-deterministic for
// ReasonConcurrency.
type ConcurrencyKind int

const (
	ConcurrencyKindGo ConcurrencyKind = iota
	ConcurrencyKindRecv
	ConcurrencyKindSend
	ConcurrencyKindRange
)

// ReasonMapRange represents iterating over a map via range.
type ReasonMapRange struct {
	SourcePos *token.Position
}

// Pos returns the source position.
func (r *ReasonMapRange) Pos() *token.Position { return r.SourcePos }

// String returns the reason.
func (r *ReasonMapRange) String() string {
	return "iterates over map"
}

func init() {
	// Needed for go vet usage
	gob.Register(&ReasonDecl{})
	gob.Register(&ReasonFuncCall{})
	gob.Register(&ReasonVarAccess{})
	gob.Register(&ReasonConcurrency{})
	gob.Register(&ReasonMapRange{})
}
