package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/chennqqi/goutils/closeevent"
	"github.com/google/subcommands"
)

type webCmd struct {
	w        Web
	port     int
	scanto   string
	callback string
	datadir  string
	indexdir string
}

func (p *webCmd) Name() string {
	return "web"
}

func (p *webCmd) Synopsis() string {
	return "web"
}

func (p *webCmd) Usage() string {
	return "web -p port"
}

func (p *webCmd) SetFlags(f *flag.FlagSet) {
	f.IntVar(&p.port, "p", 8080, "set port")
	f.StringVar(&p.scanto, "timeout", "2h", "set scan dir timeout")
	f.StringVar(&p.callback, "callback", "", "set callback addr")
	f.StringVar(&p.datadir, "data", "/dev/shm", "set data dir")
	f.StringVar(&p.indexdir, "index", "/dev/shm/.persist", "set index dir")
}

func (p *webCmd) Execute(context.Context, *flag.FlagSet, ...interface{}) subcommands.ExitStatus {
	to, err := time.ParseDuration(p.scanto)
	if err != nil {
		to, _ = time.ParseDuration("2h")
	}
	if p.callback == "" {
		p.callback = os.Getenv("JTRD_CALLBACK")
	}

	if p.callback == "" {
		p.callback = os.Getenv("HMBD_CALLBACK")
	}

	w, err := NewWeb(p.datadir, p.indexdir, p.callback, to)
	if err != nil {
		fmt.Println("new web error:", err)
		return subcommands.ExitFailure
	}

	ctx, cancel := context.WithCancel(context.Background())
	go w.Run(p.port, ctx)

	closeevent.Wait(func(s os.Signal) {
		defer cancel()
		ctx := context.Background()
		w.Shutdown(ctx)
	}, os.Interrupt)
	return subcommands.ExitSuccess
}
