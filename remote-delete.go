package main

import (
	"context"
	"flag"

	"github.com/google/subcommands"
)

func (p *RemoteDelete) Name() string {
	return "remote-delete"
}

func (p *RemoteDelete) Synopsis() string {
	return "Delete the archive that push would create from the remote"
}

func (p *RemoteDelete) Usage() string {
	return `remote-delete [env]
	Deletes the archive that push would create from the remote.
`
}

func (p *RemoteDelete) SetFlags(f *flag.FlagSet) {
}

func (p *RemoteDelete) Execute(ctx context.Context, f *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {
	return subcommands.ExitSuccess
}
