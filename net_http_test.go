package epf_test

import (
	"fmt"
	"testing"

	"github.com/haijima/epf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenMatchingGroupFromNetHttp(t *testing.T) {
	ext := &epf.NetHttpExtractor{}
	endpoints, err := FindEndpoints("testdata/src/net_http", "./...", ext)

	require.NoError(t, err)
	require.Equal(t, 9, len(endpoints))
	assert.Equal(t, "GET /foo/{id}", endpoints[0].String())
	assert.Equal(t, "POST /foo/{id}", endpoints[1].String())
	assert.Equal(t, "GET /foo/{foo_id}/bar/{bar_id}", endpoints[2].String())
	assert.Equal(t, "GET /foo/{id...}", endpoints[3].String())
	assert.Equal(t, "GET /foo/{$}", endpoints[4].String())
	assert.Equal(t, "ANY /foo/{id}", endpoints[5].String())
	assert.Equal(t, "GET /foo/", endpoints[6].String())
	assert.Equal(t, "GET /foo", endpoints[7].String())
	assert.Equal(t, "GET /foo/bar", endpoints[8].String())

	for i, mg := range endpoints {
		fmt.Printf("%d: %s %s %s\n", i, mg, mg.FuncName, mg.DeclarePos.PositionString())
	}
}
