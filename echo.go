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

type EchoExtractor struct{}

func (e *EchoExtractor) Extract(callInfo ssautil.CallInfo, parent *ssa.Function, pos token.Pos) (*Endpoint, bool) {
	if c, ok := callInfo.(*ssautil.StaticMethodCall); ok {
		if c.Method().Pkg().Path() != "github.com/labstack/echo/v4" {
			return nil, false
		}

		e := Endpoint{}
		m := c.Method().Name()
		// Method
		switch m {
		case "GET", "POST", "PUT", "PATCH", "DELETE":
			e.Method = m
		case "Add":
			if method, ok := valueToStringConst(c.Arg(0)); ok {
				e.Method = method
			}
		case "Static", "File":
			e.Method = "-"
		case "Any":
			e.Method = "*"
		case "Match":
			e.Method = "-"
		default:
			return nil, false
		}
		// Path
		var pathArgIdx int
		if m == "Add" || m == "Match" {
			pathArgIdx = 1
		}
		if path, ok := path(c, c.Arg(pathArgIdx)); ok {
			e.Path = path
			e.PathRegexpPattern = echoRegexpPattern(e.Path)
		} else {
			e.Path = "-"
		}
		// Function name
		if m == "Static" || m == "File" {
			e.FuncName = "-"
			e.DeclarePos = ssautil.NewPos(c.StaticCallee(), c.Pos(), pos)
		} else {
			fnArgIdx := 1
			if m == "Add" || m == "Match" {
				fnArgIdx = 2
			}

			e.FuncName = "-"
			e.DeclarePos = ssautil.NewPos(c.Arg(fnArgIdx).Parent(), c.Arg(fnArgIdx).Pos(), pos)
			switch v := c.Arg(fnArgIdx).(type) {
			case *ssa.ChangeType:
				switch t := v.X.(type) {
				case *ssa.Function:
					e.FuncName = t.Name()
					e.DeclarePos = ssautil.NewPos(t, t.Pos(), pos)
				case *ssa.MakeClosure:
					switch fn := t.Fn.(type) {
					case *ssa.Function:
						e.FuncName = fn.Name()
						e.DeclarePos = ssautil.NewPos(t.Parent(), fn.Pos(), pos)
					}
				}
			default:
				e.DeclarePos = ssautil.NewPos(v.Parent(), v.Pos(), pos)
				fmt.Printf("3  %T", v)
			}
		}
		return &e, true
	}
	return nil, false
}

func path(c *ssautil.StaticMethodCall, arg ssa.Value) (string, bool) {
	if arg, ok := valueToStringConst(arg); ok {
		switch c.Recv().Type().String() {
		case "*github.com/labstack/echo/v4.Echo":
			return arg, true
		case "*github.com/labstack/echo/v4.Group":
			paths, ok := groupPrefixes(c)
			return strings.Join(append(paths, arg), ""), ok
		}
	} else {
		slog.Warn("failed to parse path", "arg", arg)
	}
	return "", false
}

func groupPrefixes(staticMethodCall *ssautil.StaticMethodCall) ([]string, bool) {
	if call, ok := ssautil.ValueToCallCommon(staticMethodCall.Recv()); ok {
		if c, ok := ssautil.GetCallInfo(call).(*ssautil.StaticMethodCall); ok {
			if s, ok := valueToStringConst(c.Arg(0)); ok {
				switch c.Name() {
				case "(*github.com/labstack/echo/v4.Echo).Group":
					return []string{s}, true
				case "(*github.com/labstack/echo/v4.Group).Group":
					prefixes, ok := groupPrefixes(c)
					return append(prefixes, s), ok
				}
			} else {
				slog.Warn("failed to parse path", "arg", c.Arg(0))
			}
		}
	}
	return nil, false
}

var echoPathParamPattern = regexp.MustCompile(":([^/]+)")

func echoRegexpPattern(path string) string {
	return fmt.Sprintf("^%s$", echoPathParamPattern.ReplaceAllString(path, "([^/]+)"))
}
