package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/netip"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/sys/windows"
	"golang.org/x/term"

	gotitle "github.com/lxi1400/GoTitle"
)

// all available algorithms in ssh package
var preferredKexAlgos = []string{
	"diffie-hellman-group1-sha1",
	"diffie-hellman-group14-sha1",
	"diffie-hellman-group14-sha256",
	"diffie-hellman-group16-sha512",
	"ecdh-sha2-nistp256",
	"ecdh-sha2-nistp384",
	"ecdh-sha2-nistp521",
	"curve25519-sha256@libssh.org",
	"curve25519-sha256",
	"diffie-hellman-group-exchange-sha1",
	"diffie-hellman-group-exchange-sha256",
}

func connect(hostname, addrStr, username, pass string) error {
	_, err := netip.ParseAddrPort(addrStr)

	if err != nil {
		_, err := netip.ParseAddr(addrStr)
		if err != nil {
			return err
		}
		addrStr += ":22"
	}

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
		Timeout:         3 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	config.KeyExchanges = preferredKexAlgos

	client, err := ssh.Dial("tcp", addrStr, config)
	if err != nil {
		return fmt.Errorf("Unable to connect with (%s): %w", addrStr, err)
	}
	defer client.Close()

	s, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("Open session: %w", err)
	}

	state, err := makeRaw(os.Stdin)
	if err != nil {
		return fmt.Errorf("Can't make raw: %w", err)
	}
	defer term.Restore(int(os.Stdin.Fd()), state)

	w, h := 100, 100
	w, h, err = term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		slog.Warn("Get terminal size: %w", err)
	}

	modes := ssh.TerminalModes{
		ssh.ECHO: 1, // Disable echoing
	}
	err = s.RequestPty("vt100", h, w, modes)
	if err != nil {
		return fmt.Errorf("Request PTY: %w", err)
	}

	s.Stdout = os.Stdout
	s.Stdin = os.Stdin
	s.Stderr = os.Stderr

	// replace backspace by Ctrl+H
	if pipeR, pipeW, err := os.Pipe(); err == nil {
		s.Stdin = pipeR

		go func() {
			var err error
			var keyBackspace = byte(127)
			var keyCtrlH = byte(8)
			data := make([]byte, 1, 1)
			for err == nil {
				_, err = os.Stdin.Read(data)
				if data[0] == keyBackspace {
					data[0] = keyCtrlH
				}
				_, err = pipeW.Write(data)
			}
		}()
	} else {
		slog.Warn(fmt.Sprintf("Backspace will NOT be replaced by CTRL+H: %s", err))
	}

	if err = s.Shell(); err != nil {
		return fmt.Errorf("Open shell: %w", err)
	}

	// ignore errors
	gotitle.SetTitle(hostname)

	defer func() {
		slog.Info(fmt.Sprintf("Connection to %s closed.\n", addrStr))
	}()

	var exitMissingError *ssh.ExitMissingError
	if err = s.Wait(); err != nil && !errors.As(err, &exitMissingError) {
		return err
	}

	return nil
}

func makeRaw(file *os.File) (*term.State, error) {
	// need to transfer Ctrl+C and etc to target system
	state, err := term.MakeRaw(int(file.Fd()))
	if err != nil {
		return nil, err
	}

	// need to transfer navigation key: "home", "end" and etc
	var st uint32
	if err := windows.GetConsoleMode(windows.Handle(int(file.Fd())), &st); err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	st = st | windows.ENABLE_VIRTUAL_TERMINAL_INPUT
	if err := windows.SetConsoleMode(windows.Handle(int(file.Fd())), st); err != nil {
		term.Restore(int(file.Fd()), state)
		return nil, err
	}

	return state, nil
}
