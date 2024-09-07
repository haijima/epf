package epf

import (
	"go/token"

	"github.com/haijima/analysisutil/ssautil"
	"golang.org/x/tools/go/ssa"
)

type Endpoint struct {
	Method            string
	Path              string
	PathRegexpPattern string
	FuncName          string
	Comment           string
	DeclarePos        *ssautil.Posx
	//PathParams    []string
	//FuncPos    *ssautil.Posx
}

func (e *Endpoint) String() string {
	return e.Method + " " + e.Path
}

type Extractor interface {
	Extract(callInfo ssautil.CallInfo, parent *ssa.Function, pos token.Pos) (*Endpoint, bool)
}

func FindEndpoints(dir, pattern string, ext Extractor) ([]*Endpoint, error) {
	ssaProgs, err := ssautil.LoadBuildSSAs(dir, pattern)
	if err != nil {
		return nil, err
	}

	result := make([]*Endpoint, 0)
	for _, ssaProg := range ssaProgs {
		for _, fn := range ssaProg.SrcFuncs {
			for _, b := range fn.Blocks {
				for _, instr := range b.Instrs {
					if call, ok := instr.(*ssa.Call); ok {
						if p, ok := ext.Extract(ssautil.GetCallInfo(&call.Call), call.Parent(), call.Pos()); ok {
							result = append(result, p)
						}
					}
				}
			}
		}
	}

	return result, nil
}
