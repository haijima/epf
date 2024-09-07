package epf

import (
	"github.com/haijima/analysisutil/ssautil"
	"golang.org/x/tools/go/ssa"
)

func valueToStringConst(v ssa.Value) (string, bool) {
	if ss, ok := ssautil.ValueToStrings(v); ok && len(ss) == 1 {
		return ss[0], true
	}
	return "", false
}
