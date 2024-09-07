package epf_test

import (
	"testing"

	"github.com/haijima/epf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenMatchingGroupFromEchoV4(t *testing.T) {
	ext := &epf.EchoExtractor{}
	endpoints, err := FindEndpoints("testdata/src/echo_simple", "./...", ext)

	require.NoError(t, err)
	require.Equal(t, 4, len(endpoints))
	assert.Equal(t, "POST /api/users", endpoints[0].String())
	assert.Equal(t, "GET /api/users", endpoints[1].String())
	assert.Equal(t, "GET /api/users/:id", endpoints[2].String())
	assert.Equal(t, "GET /api/items", endpoints[3].String())
}

func TestGenMatchingGroupFromEchoV4_complex(t *testing.T) {
	ext := &epf.EchoExtractor{}
	endpoints, err := FindEndpoints("testdata/src/echo_complex", "./...", ext)

	require.NoError(t, err)
	require.Equal(t, 9, len(endpoints))
	assert.Equal(t, "POST /api/groups", endpoints[0].String())
	assert.Equal(t, "GET /api/groups", endpoints[1].String())
	assert.Equal(t, "GET /api/groups/:group_id/users", endpoints[2].String())
	assert.Equal(t, "GET /index", endpoints[3].String())
	assert.Equal(t, "GET /api/groups/:group_id/users/:user_id/tasks", endpoints[4].String())
	assert.Equal(t, "GET /view/screen1", endpoints[5].String())
	assert.Equal(t, "POST /auth/login", endpoints[6].String())
	assert.Equal(t, "POST /auth/logout", endpoints[7].String())
	assert.Equal(t, "POST /auth/admin/login", endpoints[8].String())
}
