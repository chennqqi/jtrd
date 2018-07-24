package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/google/subcommands"
	"github.com/malice-plugins/go-plugin-utils/utils"
)

type crackCmd struct {
}

func (p *crackCmd) Name() string {
	return "crack"
}

func (p *crackCmd) Synopsis() string {
	return "crack a password file"
}

func (p *crackCmd) Usage() string {
	return "crack <targetfile>"
}

func (p *crackCmd) SetFlags(f *flag.FlagSet) {
}

func (p *crackCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	args := flag.Args()
	if len(args) != 1 {
		fmt.Println("target password file is must")
		return subcommands.ExitUsageError
	}
	file := args[0]
	fd, err := os.Open(file)
	if os.IsNotExist(err) || err != nil {
		fmt.Println("target password file is must")
		return subcommands.ExitUsageError
	}
	fd.Close()

	/*
		docker run -it -v `pwd`/yourfiletocrack:/crackme.txt adamoss/john-the-ripper /crackme.txt
	*/
	filemap := fmt.Sprintf(`%v:/crackme.txt`, file)
	params := []string{
		"run",
		"-it",
		"--privileged",
		"-v",
		filemap,
		"adamoss/john-the-ripper",
		"/crackme.txt",
	}

	ctx := context.TODO()
	r, err := utils.RunCommand(ctx, "docker", params...)
	if err != nil {
		fmt.Println(r, err)
	} else {
		fmt.Println(r)
	}
	return subcommands.ExitSuccess
}
