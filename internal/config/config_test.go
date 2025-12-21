package config

import (
	"flag"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func resetFlag(args []string) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = args
}

func TestLoadConfig_DefaultValues(t *testing.T) {
	resetFlag([]string{"cmd"})

	cfg := LoadConfig()

	require.NotNil(t, cfg)

	assert.Equal(t, "localhost", cfg.Host)
	assert.Equal(t, 7879, cfg.Port)
	assert.Equal(t, 15*time.Minute, cfg.SessionTTL)
	assert.Equal(t, "data/", cfg.DataStore)
}

func TestLoadConfig_CustomValues(t *testing.T) {
	resetFlag([]string{
		"cmd",
		"-host=0.0.0.0",
		"-port=8181",
		"-session_ttl=30m",
		"-data_store=/var/lib/simplevault",
	})

	cfg := LoadConfig()

	require.NotNil(t, cfg)

	assert.Equal(t, "0.0.0.0", cfg.Host)
	assert.Equal(t, 8181, cfg.Port)
	assert.Equal(t, 30*time.Minute, cfg.SessionTTL)
	assert.Equal(t, "/var/lib/simplevault", cfg.DataStore)
}

func TestGetAddr(t *testing.T) {
	cfg := &Config{
		Host: "127.0.0.1",
		Port: 9999,
	}

	addr := cfg.GetAddr()

	assert.Equal(t, "127.0.0.1:9999", addr)
}
