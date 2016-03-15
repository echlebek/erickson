package main

import (
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/echlebek/erickson/db"
	"github.com/echlebek/erickson/mail"
	"github.com/echlebek/erickson/server"
	"github.com/gorilla/csrf"
)

const usage = `erickson code review
====================

Usage:
  $ erickson config           # Prints a configuration template to stdout
  $ erickson config file.cfg  # Runs erickson with file.cfg as configuration
`

type Configuration struct {
	Server serverCfg `toml:"server"`
	Mail   mailCfg   `toml:"mail"`
}

type serverCfg struct {
	Database   string `toml:"database"`
	SessionKey string `toml:"session_key"`
	TLSCert    string `toml:"tls_cert"`
	TLSKey     string `toml:"tls_key"`
	Port       string `toml:"port"`
	URLRoot    string `toml:"url_root"`
}

type mailCfg struct {
	Server   string `toml:"server"`
	Username string `toml:"username"`
	Password string `toml:"password"`
}

var config = Configuration{
	Server: serverCfg{
		Database: "erickson.db",

		// SessionKey should be a 32 byte random key
		SessionKey: "12345678901234567890123456789012",

		TLSCert: "example.crt",
		TLSKey:  "example.key",
		Port:    "8080",

		// URLRoot should be the root the service is hosted at.
		// Used for building URLs.
		URLRoot: "https://localhost:8081",
	},
	Mail: mailCfg{
		Server:   "mail.example.com",
		Username: "erickson",
		Password: "secret",
	},
}

func exec() {
	scfg := config.Server
	db, err := db.NewBoltDB(scfg.Database)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	smtpAuth := smtp.PlainAuth("", config.Mail.Username, config.Mail.Password, config.Mail.Server)
	mailer := mail.NewMailer(config.Mail.Server, config.Mail.Username, smtpAuth)

	handler := server.NewRootHandler(db, ".", []byte(scfg.SessionKey), mailer, scfg.URLRoot)
	if scfg.TLSCert != "" && scfg.TLSKey != "" && len(scfg.SessionKey) == 32 {
		CSRF := csrf.Protect([]byte(scfg.SessionKey))
		log.Fatal(http.ListenAndServeTLS(":"+scfg.Port, scfg.TLSCert, scfg.TLSKey, CSRF(handler)))
	} else {
		log.Fatal(http.ListenAndServe(":"+scfg.Port, handler))
	}
}

func loadConfig(path string) {
	if _, err := toml.DecodeFile(path, &config); err != nil {
		log.Fatalf("couldn't parse config file: %s", err)
	}
}

func printConfigTemplate() {
	enc := toml.NewEncoder(os.Stdout)
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
