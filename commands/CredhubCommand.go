package commands

import (
	"encoding/json"
	"fmt"
	"strings"

	"code.cloudfoundry.org/cli/cf/errors"
	"github.com/rabobank/credhub-plugin/conf"
	"github.com/rabobank/credhub-plugin/util"
	"github.com/rabobank/credhub-service-broker/model"
)

type DeleteResponse struct {
	IgnoredKeys []string
}

type CredhubCommand struct {
	Command     int
	ServiceName string
	Payload     interface{}
}

var commandMapping = map[string]int{
	conf.COMMANDS[conf.ADD_SECRET].Name:        conf.ADD_SECRET,
	conf.COMMANDS[conf.DEL_SECRET].Name:        conf.DEL_SECRET,
	conf.COMMANDS[conf.LIST_SECRETS].Name:      conf.LIST_SECRETS,
	conf.COMMANDS[conf.LIST_VERSIONS].Name:     conf.LIST_VERSIONS,
	conf.COMMANDS[conf.REINSTATE_VERSION].Name: conf.REINSTATE_VERSION,
}

var (
	MissingServiceInstanceError = errors.New("missing service instance name")
	UnknownCommandError         = errors.New("unknown command")
	BadCommandSyntaxError       = errors.New("bad command syntax")
	BadJsonObjectError          = errors.New("invalid JSON object")
)

func ParseCommand(args []string) (*CredhubCommand, error) {

	if len(args) < 2 {
		return nil, MissingServiceInstanceError
	}

	command := &CredhubCommand{
		ServiceName: args[1],
	}

	var found bool
	command.Command, found = commandMapping[args[0]]
	if !found {
		return nil, UnknownCommandError
	}

	var parseError error
	switch command.Command {
	case conf.ADD_SECRET:
		command.Payload, parseError = parseAddParameters(args[2:])
	case conf.DEL_SECRET:
		command.Payload, parseError = parseDeleteParameters(args[2:])
	case conf.REINSTATE_VERSION:
		command.Payload, parseError = parseReinstateParameters(args[2:])
	default:
		if len(args) > 2 {
			parseError = BadCommandSyntaxError
		}
	}

	return command, parseError
}

func parseAddParameters(args []string) (interface{}, error) {
	numArgs := len(args)
	if numArgs > 2 || numArgs < 1 {
		// add should either get a json objects or a key/value pair
		return nil, BadCommandSyntaxError
	}

	secrets := make(map[string]interface{})
	if numArgs == 1 {
		if e := json.Unmarshal([]byte(args[0]), &secrets); e != nil {
			fmt.Println(e)
			return nil, BadJsonObjectError
		}
	} else { // it's either 1 or two arguments
		pointer := secrets
		keys := strings.Split(args[0], ".")
		numberOfKeys := len(keys)
		for i, key := range keys {
			if i+1 == numberOfKeys {
				pointer[key] = args[1]
			} else {
				pointer[key] = make(map[string]interface{})
				pointer = pointer[key].(map[string]interface{})
			}
		}
	}

	return secrets, nil
}

func parseDeleteParameters(args []string) (interface{}, error) {
	if len(args) < 1 {
		// should get at least one key to delete
		return nil, BadCommandSyntaxError
	}

	return args, nil
}

func parseReinstateParameters(args []string) (interface{}, error) {
	if len(args) != 1 {
		// should get a single id to reinstate
		return nil, BadCommandSyntaxError
	}

	return args[0], nil
}

func printKeys(keys []string) {
	fmt.Println()
	for _, key := range keys {
		fmt.Println("  ", key)
	}
	fmt.Println()
}

func ListSecrets(brokerUrl, serviceGuid, token string, ignoreSsl bool) error {
	keys := make([]string, 0)
	if e := util.Request(brokerUrl, "api", serviceGuid, "keys").IgnoringSsl(ignoreSsl).WithAuthorization(token).GetJson(&keys); e != nil {
		return e
	}
	printKeys(keys)

	return nil
}

func AddSecrets(brokerUrl, serviceGuid, token string, payload interface{}, ignoreSsl bool) error {
	content, e := json.Marshal(payload)
	if e != nil {
		return e
	}
	if _, e = util.Request(brokerUrl, "api", serviceGuid, "keys").IgnoringSsl(ignoreSsl).WithAuthorization(token).Sending("application/json").WithContent(content).Put(); e != nil {
		return e
	}
	return nil
}

func DeleteSecrets(brokerUrl, serviceGuid, token string, payload interface{}, ignoreSsl bool) error {
	content, e := json.Marshal(payload)
	if e != nil {
		return e
	}
	response := model.DeleteResponse{}
	if e = util.Request(brokerUrl, "api", serviceGuid, "keys").IgnoringSsl(ignoreSsl).WithAuthorization(token).Sending("application/json").WithContent(content).DeleteReturningJson(&response); e != nil {
		return e
	}

	if response.IgnoredKeys != nil {
		fmt.Println("Ignored keys (not existing):")
		printKeys(response.IgnoredKeys)
	}
	return nil
}

func ListVersions(brokerUrl, serviceGuid, token string, ignoreSsl bool) error {
	versions := make([]model.SecretsVersionKeys, 0)
	if e := util.Request(brokerUrl, "api", serviceGuid, "versions").IgnoringSsl(ignoreSsl).WithAuthorization(token).GetJson(&versions); e != nil {
		return e
	}

	fmt.Println()
	for _, version := range versions {
		fmt.Println("ID:", version.ID)
		fmt.Println("Created:", version.VersionCreatedAt)
		printKeys(version.Keys)
	}
	fmt.Println()

	return nil
}

func ReinstateVersion(brokerUrl, serviceGuid, token string, payload interface{}, ignoreSsl bool) error {
	if _, e := util.Request(brokerUrl, "api", serviceGuid, "version", payload.(string)).IgnoringSsl(ignoreSsl).WithAuthorization(token).Put(); e != nil {
		return e
	}

	return nil
}
