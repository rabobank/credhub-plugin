package commands

import (
	"encoding/json"
	"fmt"
	"strings"

	"code.cloudfoundry.org/cli/cf/errors"
	"github.com/rabobank/credhub-plugin/conf"
	"github.com/rabobank/credhub-plugin/util"
)

type CredhubCommand struct {
	Command     int
	ServiceName string
	Payload     interface{}
}

var commandMapping = map[string]int{
	conf.COMMANDS[conf.ADD_SECRET].Name:   conf.ADD_SECRET,
	conf.COMMANDS[conf.DEL_SECRET].Name:   conf.DEL_SECRET,
	conf.COMMANDS[conf.LIST_SECRETS].Name: conf.LIST_SECRETS,
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

	fmt.Println(secrets)

	return secrets, nil
}

func parseDeleteParameters(args []string) (interface{}, error) {
	if len(args) < 1 {
		// should get at least one key to delete
		return nil, BadCommandSyntaxError
	}

	return args, nil
}

func ListSecrets(brokerUrl, serviceGuid, token string) error {
	keys := make([]string, 0)
	if e := util.Request(brokerUrl, "api", serviceGuid, "keys").WithAuthorization(token).GetJson(&keys); e != nil {
		return e
	}

	fmt.Println()
	for _, key := range keys {
		fmt.Println(key)
	}
	fmt.Println()

	return nil
}

func AddSecrets(brokerUrl, serviceGuid, token string, payload interface{}) error {
	return nil
}

func DeleteSecrets(brokerUrl, serviceGuid, token string, payload interface{}) error {
	return nil
}
