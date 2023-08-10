package main

import (
	"fmt"
	"os"

	"code.cloudfoundry.org/cli/plugin"
	"github.com/rabobank/credhub-plugin/commands"
	"github.com/rabobank/credhub-plugin/conf"
)

type CredhubPlugin struct{}

func (c *CredhubPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	if args[0] == "CLI-MESSAGE-UNINSTALL" {
		return
	}

	_, e := commands.ParseCommand(args)
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}

}

func (c *CredhubPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "credhub-plugin",
		Version: plugin.VersionType{
			Major: 1,
			Minor: 0,
			Build: 0,
		},
		MinCliVersion: plugin.VersionType{
			Major: 7,
			Minor: 1,
			Build: 0,
		},
		Commands: conf.COMMANDS,
	}
}

func main() {
	if len(os.Args) == 1 {
		_, _ = fmt.Fprintf(os.Stderr, "This executable is a cf plugin.\n"+
			"Run `cf install-plugin %s` to install it",
			os.Args[0])
		os.Exit(1)
	}

	plugin.Start(new(CredhubPlugin))
}
