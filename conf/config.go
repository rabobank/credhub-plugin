package conf

import (
	"fmt"
	"strconv"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
)

const (
	ADD_SECRET = iota
	DEL_SECRET
	LIST_SECRETS
	LIST_VERSIONS
	REINSTATE_VERSION
)

var COMMANDS = []plugin.Command{
	{
		Name:     "add-credhub-secrets",
		Alias:    "acs",
		HelpText: "Add secrets to credhub service",
		UsageDetails: plugin.Usage{
			Usage: "\n" +
				"cf add-credhub-secrets <SERVICE_INSTANCE> <JSON_OBJECT>\n" +
				"cf add-credhub-secrets <SERVICE_INSTANCE> <KEY> <VALUE>\n" +
				"\n" +
				"  SERVICE_INSTANCE - Credhub service instance name the keys are being added to.\n" +
				"\n" +
				"  JSON_OBJECT      - A well formed json object map. Key values will either replace existing keys or added to the existing credentials if not present\n" +
				"                     This will only be interpreted as a json object it the KEY/VALUE parameters are not provided.\n" +
				"  KEY              - When a VALUE is provided, instead of JSON_OBJECT, the first parameter will be interpreted as the secret key.\n" +
				"                     If updating/setting encapsulated values, dots may be used to reference the inner-keys (i.e. a.b to reference {\"a\":{\"b\":\"value\"}})\n" +
				"  VALUE            - Secret value.\n",
		},
	},
	{
		Name:     "delete-credhub-secrets",
		Alias:    "dcs",
		HelpText: "Delete a key from the credhub service instance",
		UsageDetails: plugin.Usage{
			Usage: "\n" +
				"cf delete-credhub-secrets <SERVICE_INSTANCE> <KEYS>...\n" +
				"\n" +
				"  SERVICE_INSTANCE - Credhub service instance name the keys are being deleted from.\n" +
				"\n" +
				"  KEYS             - Secret keys to delete. Multiple keys can be provided separated by spaces.\n",
		},
	},
	{
		Name:     "list-credhub-secrets",
		Alias:    "lcs",
		HelpText: "List all secret keys in the credhub service instance",
		UsageDetails: plugin.Usage{
			Usage: "\n" +
				"cf list-credhub-secrets <SERVICE_INSTANCE>\n" +
				"\n" +
				"  SERVICE_INSTANCE - Credhub service instance name.\n",
		},
	},
	{
		Name:     "list-credhub-secrets-versions",
		Alias:    "lcv",
		HelpText: "List up to 20 latest versions for a credhub service instance credentials",
		UsageDetails: plugin.Usage{
			Usage: "\n" +
				"cf list-credhub-secrets-versions <SERVICE_INSTANCE>\n" +
				"\n" +
				"  SERVICE_INSTANCE - Credhub service instance name.\n",
		},
	},
	{
		Name:     "reinstate-credhub-secrets-version",
		Alias:    "rcv",
		HelpText: "Reinstate a previous version of the credhub service instance credentials",
		UsageDetails: plugin.Usage{
			Usage: "\n" +
				"cf reinstate-credhub-secrets-version <SERVICE_INSTANCE> <VERSION_ID>\n" +
				"\n" +
				"  SERVICE_INSTANCE - Credhub service instance name.\n" +
				"  VERSION_ID       - The credentials version id to reinstate. Can be obtained from the list-credhub-secrets-versions command.\n",
		},
	},
}

var (
	VERSION  = "0.0.0"
	COMMIT   = "dev"
	Metadata plugin.PluginMetadata
)

func Initialize() {
	versionParts := strings.Split(VERSION, ".")
	var major, minor, build int
	var e error

	if strings.HasPrefix(versionParts[0], "v") {
		versionParts[0] = versionParts[0][1:]
	}

	major, e = strconv.Atoi(versionParts[0])
	if e != nil {
		fmt.Printf("invalid major version : %s. Defaulting to 0\n", versionParts[0])
		major = 0
	}

	if len(versionParts) > 1 {
		minor, e = strconv.Atoi(versionParts[1])
		if e != nil {
			fmt.Printf("invalid minor version : %s. Defaulting to 0\n", versionParts[1])
			minor = 0
		}
	}

	if len(versionParts) > 2 {
		if dashPosition := strings.Index(versionParts[2], "-"); dashPosition > 0 {
			versionParts[2] = versionParts[2][:dashPosition]
		}
		build, e = strconv.Atoi(versionParts[2])
		if e != nil {
			fmt.Printf("invalid build version : %s. Defaulting to 0\n", versionParts[2])
			build = 0
		}
	}

	Metadata = plugin.PluginMetadata{
		Name: "credhub-plugin",
		Version: plugin.VersionType{
			Major: major,
			Minor: minor,
			Build: build,
		},
		MinCliVersion: plugin.VersionType{
			Major: 7,
			Minor: 1,
			Build: 0,
		},
		Commands: COMMANDS,
	}
}
