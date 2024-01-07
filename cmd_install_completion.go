package main

import (
	"embed"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//go:embed dssh.lua
var dsshLua embed.FS

const (
	CLINK_COMPLETIONS_DIR = "CLINK_COMPLETIONS_DIR"
)

func copyLua(clinkDir string) error {

	src, _ := dsshLua.Open("dssh.lua")
	defer src.Close()
	destFile := filepath.Join(clinkDir, "dssh.lua")

	dst, err := os.Create(destFile)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)

	return err
}

func installCompletion() error {
	envs := exec.Command("start").Environ()
	dstDir, ok := getEnvValue(CLINK_COMPLETIONS_DIR, envs)
	if !ok {
		return fmt.Errorf("Clink is not installed")
	}

	if err := copyLua(dstDir); err != nil {
		return err
	}

	return nil
}

func getEnvValue(key string, data []string) (string, bool) {
	for _, d := range data {
		if !strings.HasPrefix(d, key) {
			continue
		}

		return strings.TrimPrefix(d, key+"="), true
	}

	return "", false
}
