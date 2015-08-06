package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
)

type Color int

const (
	None = Color(iota)
	Red
	Blue
)

func handleSignal() {
	sc := make(chan os.Signal, 10)
	signal.Notify(sc, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
	go func() {
		<-sc
		os.Exit(0)
	}()
}

func ready() error {
	vendor, err := filepath.Abs(vendorFolder)
	if err != nil {
		return err
	}

	binPath := strings.Join(
		[]string{filepath.Join(vendor, "bin"), os.Getenv("PATH")},
		string(filepath.ListSeparator),
	)

	if *verbose {
		fmt.Printf("setenv PATH=%s\n", binPath)
	}
	err = os.Setenv("PATH", binPath)
	if err != nil {
		return err
	}

	gopath := strings.Join(
		[]string{vendor, os.Getenv("GOPATH")},
		string(filepath.ListSeparator),
	)
	if *verbose {
		fmt.Printf("setenv GOPATH=%s\n", gopath)
	}
	err = os.Setenv("GOPATH", gopath)
	if err != nil {
		return err
	}

	return nil
}

var stdout = os.Stdout
var stderr = os.Stderr
var stdin = os.Stdin

func run(args []string, c Color) error {
	if err := ready(); err != nil {
		return err
	}
	if len(args) == 0 {
		usage()
	}
	if *verbose {
		fmt.Printf("%q\n", args)
	}
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.Stdin = stdin
	err := cmd.Run()
	return err
}
