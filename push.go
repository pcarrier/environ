package main

import (
	"context"
	"flag"

	"github.com/google/subcommands"
)

func (p *Push) Name() string {
	return "push"
}

func (p *Push) Synopsis() string {
	return "Push an archive of the tracked files to the remote"
}

func (p *Push) Usage() string {
	return `push [env]
	Pushes an archive of the tracked files to the remote.
`
}

func (p *Push) SetFlags(f *flag.FlagSet) {
}

func (p *Push) Execute(ctx context.Context, f *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {
	return subcommands.ExitSuccess
}
