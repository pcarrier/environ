package main

import (
	"context"
	"flag"

	"github.com/google/subcommands"
)

func (p *Pull) Name() string {
	return "pull"
}

func (p *Pull) Synopsis() string {
	return "Pull files from the remote into the local directory"
}

func (p *Pull) Usage() string {
	return `pull [env]
	Pull files from the remote into the local directory.
`
}

func (p *Pull) SetFlags(f *flag.FlagSet) {
}

func (p *Pull) Execute(ctx context.Context, f *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {
	return subcommands.ExitSuccess
}
