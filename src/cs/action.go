package cs

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/codegangsta/cli"
	"github.com/deis/riak/src/replace"
)

const (
	confFilePath = "/etc/riak-cs/riak-cs.conf"
)

func Action(ctx *cli.Context) {
	conf, err := getConfig()
	if err != nil {
		log.Printf("Error: getting config (%s)", err)
		os.Exit(1)
	}
	stanchionURL := fmt.Sprintf("%s:%d", conf.StanchionHost, conf.StanchionPort)

	adminKey, err := ioutil.ReadFile(conf.AdminKeyLocation)
	if err != nil {
		log.Printf("Error: reading admin key from %s (%s)", conf.AdminKeyLocation, err)
		os.Exit(1)
	}

	adminSecret, err := ioutil.ReadFile(conf.AdminSecretLocation)
	if err != nil {
		log.Printf("Error: reading admin secret from %s (%s)", conf.AdminSecretLocation, err)
		os.Exit(1)
	}

	confFile, err := ioutil.ReadFile(confFilePath)
	if err != nil {
		log.Printf("Error: reading Riak CS config file from %s (%s)", confFilePath, err)
		os.Exit(1)
	}

	replacements := []replace.Replacement{
		replace.FmtReplacement("listener = 127.0.0.1:8080", "listener = %s:%d", conf.ListenHost, conf.ListenPort),
		replace.FmtReplacement("stanchion_host = 127.0.0.1:8085", "stanchion_host = %s", stanchionURL),
		replace.FmtReplacement("stanchion_ssl = on", "stanchion_ssl = off"),
		replace.FmtReplacement("admin.key = admin-key", "admin.key = %s", adminKey),
		replace.FmtReplacement("admin.secret = admin-secret", "admin.secret = %s", adminSecret),
	}
	newConfFile := replace.String(string(confFile), replacements...)
	if err := ioutil.WriteFile(confFilePath, []byte(newConfFile), os.ModePerm); err != nil {
		log.Printf("Error: writing new config file to %s (%s)", confFilePath, err)
		os.Exit(1)
	}

	log.Printf("Starting Riak CS...")
	startCmd := exec.Command("riak-cs", "start")
	startCmd.Stdout = os.Stdout
	startCmd.Stderr = os.Stderr
	if err := startCmd.Run(); err != nil {
		log.Printf("Error: starting Riak CS (%s)", err)
		os.Exit(1)
	}
	log.Printf("started")
	select {}

}