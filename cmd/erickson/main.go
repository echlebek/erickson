package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/echlebek/erickson/db"
	"github.com/echlebek/erickson/server"
	"github.com/gorilla/csrf"
)

const usage = `erickson code review
====================

Usage:
  $ erickson config           # Prints a configuration template to stdout
  $ erickson config file.cfg  # Runs erickson with file.cfg as configuration
`

type serverCfg struct {
	Database   string `toml:"database"`
	SessionKey string `toml:"session_key"`
	TLSCert    string `toml:"tls_cert"`
	TLSKey     string `toml:"tls_key"`
	Port       string `toml:"port"`
}

var config = serverCfg{
	Database: "erickson.db",

	// SessionKey should be a 32 byte random key
	SessionKey: "12345678901234567890123456789012",

	TLSCert: "example.crt",
	TLSKey:  "example.key",
	Port:    "8080",
}

func exec() {
	db, err := db.NewBoltDB(config.Database)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	handler := server.NewRootHandler(db, ".", []byte(config.SessionKey))
	if config.TLSCert != "" && config.TLSKey != "" && len(config.SessionKey) == 32 {
		CSRF := csrf.Protect([]byte(config.SessionKey))
		log.Fatal(http.ListenAndServeTLS(":"+config.Port, config.TLSCert, config.TLSKey, CSRF(handler)))
	} else {
		log.Fatal(http.ListenAndServe(":"+config.Port, handler))
	}
}

func loadConfig(path string) {
	cfg := struct {
		Server serverCfg `toml:"server"`
	}{}

	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		log.Fatalf("couldn't parse config file: %s", err)
	}

	config = cfg.Server
}

func printConfigTemplate() {
	enc := toml.NewEncoder(os.Stdout)
	config := struct { // Wrap config to give it a nice heading
		Server serverCfg `toml:"server"`
	}{config}
	if err := enc.Encode(config); err != nil {
		// Shouldn't ever happen, so make some noise
		panic(err)
	}
	fmt.Fprintln(os.Stdout)
}

func main() {
	args := os.Args
	if len(args) > 1 && len(args) < 4 && args[1] == "config" {
		if len(args) > 2 {
			loadConfig(args[2])
			exec()
		} else {
			printConfigTemplate()
		}
	} else {
		fmt.Fprintln(os.Stdout, usage)
	}
}
