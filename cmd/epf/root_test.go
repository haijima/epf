package main

import (
	"bytes"
	"testing"

	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRootCmd(t *testing.T) {
	v := viper.New()
	fs := afero.NewMemMapFs()
	v.SetFs(fs)
	cmd := NewRootCmd(v, fs)

	assert.Equal(t, "epf", cmd.Name(), "NewRootCmd() should return command named \"epf\". but: %q", cmd.Name())
}

func TestNewRootCmd_Flag(t *testing.T) {
	v := viper.New()
	fs := afero.NewMemMapFs()
	v.SetFs(fs)
	cmd := NewRootCmd(v, fs)
	formatFlag := cmd.Flags().Lookup("format")
	dirFlag := cmd.Flags().Lookup("dir")
	patternFlag := cmd.Flags().Lookup("pattern")

	assert.True(t, cmd.HasAvailableFlags(), "epf command should have available flag")
	assert.NotNil(t, formatFlag, "epf command should have \"format\" flag")
	assert.Equal(t, "string", formatFlag.Value.Type(), "\"format\" flag is string")
	assert.NotNil(t, dirFlag, "epf command should have \"dir\" flag")
	assert.Equal(t, "string", dirFlag.Value.Type(), "\"dir\" flag is string")
	assert.NotNil(t, patternFlag, "epf command should have \"pattern\" flag")
	assert.Equal(t, "string", patternFlag.Value.Type(), "\"pattern\" flag is string")
}

func TestRunRoot(t *testing.T) {
	v := viper.New()
	fs := afero.NewOsFs()
	v.SetFs(fs)
	stdout := new(bytes.Buffer)
	cmd := NewRootCmd(v, fs)
	cmd.SetOut(stdout)
	cmd.SetArgs([]string{"--format", "csv", "--dir", "./testdata/src"})
	err := cmd.Execute()

	require.NoError(t, err)
	assert.Equal(t, "#,Method,Path,Function,Declared Package,Declared Position\n1,POST,/api/users,CreateUser,g/h/s/i/testdata,main.go:18:6\n2,GET,/api/users,GetUsers,g/h/s/i/testdata,main.go:22:6\n3,GET,/api/users/:id,GetUser,g/h/s/i/testdata,main.go:26:6\n4,GET,/api/items,GetItems,g/h/s/i/testdata,main.go:30:6\n", stdout.String())
}
