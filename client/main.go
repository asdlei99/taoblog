package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func parseFlags() {
	flag.Parse()
}

type JSONRet struct {
	Code int         `json:"code"`
	Msgs string      `json:"msgs"`
	Data interface{} `json:"data"`
}

type xInitConfig struct {
	api    string // the api root, like https://taoblog/apiv2
	verify bool   // verify host key
	key    string // the key
}

var initConfig xInitConfig

func readInitConfig() {
	var err error

	usr, err := user.Current()
	path := filepath.Join(usr.HomeDir, "/.taoblog.cfg")
	fp, err := os.Open(path)
	if err != nil {
		panic("cannot read init config: " + path)
	}

	defer fp.Close()

	buf := bufio.NewScanner(fp)
	for buf.Scan() {
		line := strings.TrimSpace(buf.Text())
		if line == "" {
			continue
		}
		toks := strings.SplitN(line, ":", 2)
		if len(toks) < 2 {
			log.Printf("invalid config: %s\n", line)
			continue
		}

		switch toks[0] {
		case "api":
			initConfig.api = toks[1]
		case "verify":
			initConfig.verify = toks[1] == "1"
		case "key":
			initConfig.key = toks[1]
		default:
			log.Printf("unrecognized config: %s\n", line)
		}
	}
}

func main() {
	parseFlags()
	readInitConfig()

	if len(os.Args) >= 2 {
		command := os.Args[1]
		if command == "post" {
			evalPost(os.Args[2:])
		} else if command == "backup" {
			evalBackup(os.Args[2:])
		}
	}
}
