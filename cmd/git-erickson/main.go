package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/echlebek/erickson/client"
)

var (
	errNoConfig   = errors.New("no config")
	errBadCommand = errors.New("bad command")
)

type cfg struct {
	url      string
	username string
}

func post(config cfg, c *client.Client, args []string) error {
	diffArgs := append([]string{"diff"}, args...)
	cmd := exec.Command("git", diffArgs...)
	output, err := cmd.Output()
	if err != nil {
		return err
	}
	diffText := string(output)
	logArgs := append([]string{"log"}, args...)
	cmd = exec.Command("git", logArgs...)
	output, err = cmd.Output()
	if err != nil {
		return err
	}
	commits := string(output)
	cmd = exec.Command("git", "config", "remote.origin.url")
	output, err = cmd.Output()
	if err != nil {
		return err
	}
	repo := string(output)
	location, err := c.PostReview(diffText, config.username, commits, repo)
	if err != nil {
		return err
	}
	fmt.Println(location)
	return nil
}

var modes = map[string]func(cfg, *client.Client, []string) error{
	"post": post,
}

func printUsage() {
	fmt.Println("TODO")
}

func main() {
	if len(os.Args) == 1 {
		printUsage()
		return
	}
	mode := os.Args[1]
	config, err := readConfig()
	if err == errNoConfig {
		if config, err = askConfig(); err != nil {
			log.Fatal(err)
		}
	} else if err != nil {
		log.Fatal(err)
	}
	client, err := client.New(config.url)
	if err != nil {
		log.Fatal(err)
	}
	if err := authenticate(client, config.username); err != nil {
		log.Fatal(err)
	}
	f, ok := modes[mode]
	if !ok {
		log.Fatal(errBadCommand)
	}
	if err := f(config, client, os.Args[2:]); err != nil {
		log.Fatal(err)
	}
}

func authenticate(client *client.Client, username string) error {
	password, err := askCredentials(username)
	if err != nil {
		return err
	}
	response, err := client.Authenticate(username, password)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		// We should have been redirected but we weren't
		return errors.New(response.Status)
	}
	return nil
}

func readConfig() (cfg, error) {
	var config cfg
	cmd := exec.Command("git", "config", "erickson.url")
	output, err := cmd.Output()
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return cfg{}, errNoConfig
		} else {
			return cfg{}, err
		}
	}
	config.url = string(bytes.Trim(output, "\n"))
	cmd = exec.Command("git", "config", "erickson.username")
	output, err = cmd.Output()
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			// the config option doesn't exist, that's ok
			return config, nil
		} else {
			return cfg{}, err
		}
	}
	config.username = string(bytes.Trim(output, "\n"))
	return config, nil
}

func askConfig() (cfg, error) {
	stdinState, err := terminal.MakeRaw(syscall.Stdin)
	if err != nil {
		return cfg{}, err
	}
	defer terminal.Restore(syscall.Stdin, stdinState)
	stdoutState, err := terminal.MakeRaw(syscall.Stdout)
	if err != nil {
		return cfg{}, err
	}
	defer terminal.Restore(syscall.Stdout, stdoutState)
	t := struct {
		io.Reader
		io.Writer
	}{os.Stdin, os.Stdout}

	term := terminal.NewTerminal(t, "")

	msg := "Configure git-erickson for first time use.\nErickson server URL: "
	if _, err := term.Write([]byte(msg)); err != nil {
		return cfg{}, err
	}

	url, err := term.ReadLine()
	if err != nil {
		return cfg{}, err
	}

	cmd := exec.Command("git", "config", "--global", "erickson.url", string(url))
	if err := cmd.Run(); err != nil {
		return cfg{}, err
	}

	if _, err := term.Write([]byte("Erickson username: ")); err != nil {
		return cfg{}, err
	}

	username, err := term.ReadLine()
	if err != nil {
		return cfg{}, err
	}

	cmd = exec.Command("git", "config", "--global", "erickson.username", string(username))
	if err := cmd.Run(); err != nil {
		return cfg{}, err
	}

	return cfg{url: url, username: username}, nil
}

func askCredentials(username string) (password string, err error) {
	stdinState, err := terminal.MakeRaw(syscall.Stdin)
	if err != nil {
		return
	}
	defer terminal.Restore(syscall.Stdin, stdinState)
	stdoutState, err := terminal.MakeRaw(syscall.Stdout)
	if err != nil {
		return
	}
	defer terminal.Restore(syscall.Stdout, stdoutState)
	t := struct {
		io.Reader
		io.Writer
	}{os.Stdin, os.Stdout}
	term := terminal.NewTerminal(t, "")
	msg := fmt.Sprintf("Password for %s: ", username)
	password, err = term.ReadPassword(msg)
	return
}
