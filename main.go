package main

import (
	"fmt"
	"os"

	"code.cloudfoundry.org/cli/cf/errors"
	"code.cloudfoundry.org/cli/plugin"
	plugin_models "code.cloudfoundry.org/cli/plugin/models"
	"github.com/rabobank/credhub-plugin/commands"
	"github.com/rabobank/credhub-plugin/conf"
	"github.com/rabobank/credhub-plugin/util"
)

var (
	OperationInProgressError = errors.New("operation in progress")
	UnsupportedServiceError  = errors.New("service is not a credhub service instance")
)

type CredhubPlugin struct{}

func (c *CredhubPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	if args[0] == "CLI-MESSAGE-UNINSTALL" {
		return
	}

	command, e := commands.ParseCommand(args)
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}

	service, e := validateServiceInstance(cliConnection, command.ServiceName)
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}

	token, e := cliConnection.AccessToken()
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}

	brokerUrl, e := getBrokerUrl(cliConnection, token, service)
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}

	switch command.Command {
	case conf.LIST_SECRETS:
		e = commands.ListSecrets(brokerUrl, service.Guid, token)
	case conf.ADD_SECRET:
		e = commands.AddSecrets(brokerUrl, service.Guid, token, command.Payload)
	case conf.DEL_SECRET:
		e = commands.DeleteSecrets(brokerUrl, service.Guid, token, command.Payload)
	case conf.LIST_VERSIONS:
		e = commands.ListVersions(brokerUrl, service.Guid, token)
	case conf.REINSTATE_VERSION:
		e = commands.ReinstateVersion(brokerUrl, service.Guid, token, command.Payload)
	}

	if e != nil {
		fmt.Printf("Failed!\n%v\n", e)
		os.Exit(1)
	}
}

func validateServiceInstance(cliConnection plugin.CliConnection, ServiceName string) (*plugin_models.GetService_Model, error) {
	serviceInstance, e := cliConnection.GetService(ServiceName)
	if e != nil {
		return nil, e
	}

	if serviceInstance.ServiceOffering.Name != "credhub" {
		return nil, UnsupportedServiceError
	}

	if serviceInstance.LastOperation.State == "in progress" {
		return nil, OperationInProgressError
	}

	return &serviceInstance, nil
}

func getBrokerUrl(cliConnection plugin.CliConnection, token string, service *plugin_models.GetService_Model) (string, error) {
	cfUrl, e := cliConnection.ApiEndpoint()
	if e != nil {
		return "", e
	}
	response := make(map[string]interface{})
	if e = util.Request(cfUrl, "v3", "service_instances", service.Guid).WithAuthorization(token).GetJson(&response); e != nil {
		return "", e
	}
	if e = util.Request(response["links"].(map[string]interface{})["service_plan"].(map[string]interface{})["href"].(string)).WithAuthorization(token).GetJson(&response); e != nil {
		return "", e
	}
	if e = util.Request(response["links"].(map[string]interface{})["service_offering"].(map[string]interface{})["href"].(string)).WithAuthorization(token).GetJson(&response); e != nil {
		return "", e
	}
	if e = util.Request(response["links"].(map[string]interface{})["service_broker"].(map[string]interface{})["href"].(string)).WithAuthorization(token).GetJson(&response); e != nil {
		return "", e
	}

	return response["url"].(string), nil
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
