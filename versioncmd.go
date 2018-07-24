package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/google/subcommands"
	"io/ioutil"
)

type versionCmd struct {
}

func (p *versionCmd) Name() string {
	return "version"
}

func (p *versionCmd) Synopsis() string {
	return `version `
}

func (p *versionCmd) Usage() string {
	return `version`
}

func (p *versionCmd) SetFlags(*flag.FlagSet) {
}

func (p *versionCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	txt, _ := ioutil.ReadFile("/opt/jtrd/VERSION")
	fmt.Println(string(txt))
	return subcommands.ExitSuccess
}
