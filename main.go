package main

import (
	"flag"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"

	"common"
	"serve"
)

func main() {
	cfgPathPtr := flag.String("config", "/etc/localmod.yaml", "path to config")
	flag.Parse()
	cmd := flag.Arg(0)

	if cmd != "serve" {
		log.Fatalln("Bad command")
	}

	data, err := ioutil.ReadFile(*cfgPathPtr)
	common.HandleError(err)
	cfg := common.Config{}
	err = yaml.Unmarshal(data, &cfg)
	common.HandleError(err)

	serve.Serve(cfg)
}
