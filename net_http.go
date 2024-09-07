package epf

import (
	"fmt"
	"go/token"
	"log/slog"
	"regexp"
	"strings"

	"github.com/haijima/analysisutil/ssautil"
	"golang.org/x/tools/go/ssa"
)

type NetHttpExtractor struct{}

func (e *NetHttpExtractor) Extract(callInfo ssautil.CallInfo, parent *ssa.Function, pos token.Pos) (*Endpoint, bool) {
	if callInfo.Match("(*net/http.ServeMux).Handle") ||
		callInfo.Match("(*net/http.ServeMux).HandleFunc") ||
		callInfo.Match("net/http.Handle") ||
		callInfo.Match("net/http.HandleFunc") {

		if arg, ok := valueToStringConst(callInfo.Arg(0)); ok {
			s := strings.Split(arg, " ")
			e := &Endpoint{
				Method:            s[0],
				Path:              s[len(s)-1],
				PathRegexpPattern: netHttRegexpPattern(s[len(s)-1]),
				DeclarePos:        ssautil.NewPos(parent, pos),
			}
			if len(s) == 1 {
				e.Method = "ANY"
			}
			switch t := callInfo.Arg(1).(type) {
			case *ssa.Function:
				e.FuncName = t.Name()
				e.DeclarePos = ssautil.NewPos(t, t.Pos(), pos)
			case *ssa.MakeInterface:
				fmt.Printf("MakeInterface: %T\n", t.X)
				switch tt := t.X.(type) {
				case *ssa.Function:
					e.FuncName = tt.Name()
					e.DeclarePos = ssautil.NewPos(tt, tt.Pos(), pos)
				case *ssa.Alloc:
					e.FuncName = tt.String()
					e.DeclarePos = ssautil.NewPos(tt.Parent(), tt.Pos(), pos)
				}
			}
			return e, true
		}
		slog.Warn("failed to parse path", "arg", callInfo.Arg(0))
	}
	return nil, false
}

var netHttpPathPatternEndsWithSlash = regexp.MustCompile("/$")
var netHttpPathPatternEndsWithDollar = regexp.MustCompile("\\{\\$\\}$")
var netHttpPathPatternEndsWithDotsParam = regexp.MustCompile("\\{([a-zA-Z0-9_-]+)\\.\\.\\.\\}$")
var netHttpPathParamPattern = regexp.MustCompile("\\{([a-zA-Z0-9_-]+)\\}")

func netHttRegexpPattern(path string) string {
	path = netHttpPathPatternEndsWithSlash.ReplaceAllString(path, "/(.*)")
	path = netHttpPathPatternEndsWithDollar.ReplaceAllString(path, "")
	path = netHttpPathPatternEndsWithDotsParam.ReplaceAllString(path, "(.+)")
	path = netHttpPathParamPattern.ReplaceAllString(path, "([^/]+)")
	return fmt.Sprintf("^%s$", path)
}
