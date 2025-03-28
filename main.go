package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/subcommands"
)

type Config struct {
	Imports []string `json:"imports,omitempty"`
	Hash    string   `json:"hash,omitempty"`
	Paths   []string `json:"paths,omitempty"`
	Remote  string   `json:"remote,omitempty"`
	AWS     struct {
		Profile string `json:"profile,omitempty"`
	} `json:"aws,omitempty"`
}

func loadConfig(env string) (*Config, error) {
	filename := ".environ"
	if env != "" {
		filename = fmt.Sprintf(".environ.%s", env)
	}

	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	var configPath string
	for {
		path := filepath.Join(dir, filename)
		fi, err := os.Stat(path)
		if err == nil && !fi.IsDir() {
			configPath = path
			break
		}

		if err != nil && !os.IsNotExist(err) {
			return nil, err
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// We've reached the root.
			break
		}
		dir = parent
	}

	if configPath == "" {
		return nil, fmt.Errorf("no config file found named %s in any ancestor of %s", filename, dir)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	if err := json.Unmarshal(data, config); err != nil {
		return nil, err
	}

	importedPaths := []string{}
	for _, importPath := range config.Imports {
		if contains(importedPaths, importPath) {
			continue
		}
		importedPaths = append(importedPaths, importPath)

		importedConfig, err := loadConfig(importPath)
		if err != nil {
			return nil, err
		}

		if config.Hash == "" {
			config.Hash = importedConfig.Hash
		}
		if config.Remote == "" {
			config.Remote = importedConfig.Remote
		}
		if config.AWS.Profile == "" {
			config.AWS.Profile = importedConfig.AWS.Profile
		}
		config.Paths = append(config.Paths, importedConfig.Paths...)
	}

	return config, nil
}

func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

type Pull struct{}
type Push struct{}
type Delete struct{}
type RemoteDelete struct{}

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

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(&Pull{}, "")
	subcommands.Register(&Push{}, "")
	subcommands.Register(&Delete{}, "")
	subcommands.Register(&RemoteDelete{}, "")

	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
