package main

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfigCorrect(t *testing.T) {
	curDir, err := filepath.Abs(".")
	require.NoError(t, err)
	cfgDir := filepath.Join(curDir, "test_files", "correct")

	cfg, err := parseConfigFiles(cfgDir)
	require.NoError(t, err)

	expected := []Host{
		{Name: "proxy_server", Address: "192.168.1.1", UserName: "ubuntu", Password: "pass",
			fileName: "correct.hosts", lineStart: 1, lineEnd: 5},
		{Name: "ftp_server", Address: "172.16.10.10:24", UserName: "anonymous", Password: "anonymous",
			fileName: "correct.hosts", lineStart: 7, lineEnd: 10},
	}

	require.Equal(t, expected, cfg.Hosts)
}

func TestConfigIncorrect(t *testing.T) {
	curDir, err := filepath.Abs(".")
	require.NoError(t, err)
	cfgDir := filepath.Join(curDir, "test_files", "incorrect")

	_, actualErr := parseConfigFiles(cfgDir)
	require.Error(t, actualErr)
	//t.Log(actualErr)
	// expected := &configFileError{
	// 	fileName: "incorrect.hosts",
	// 	err:      &yaml.TypeError{},
	// }
	// require.Equal(t, expected, actualErr)
}
