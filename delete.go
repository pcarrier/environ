package main

import (
	"context"
	"flag"

	"github.com/google/subcommands"
)

func (p *Delete) Name() string {
	return "delete"
}

func (p *Delete) Synopsis() string {
	return "Delete tracked paths from the local directory"
}

func (p *Delete) Usage() string {
	return `delete [env]
	Deletes tracked paths from the local directory.
`
}

func (p *Delete) SetFlags(f *flag.FlagSet) {
}

func (p *Delete) Execute(ctx context.Context, f *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {
	return subcommands.ExitSuccess
}
