package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/ja1code/koolseal/entity"
	"github.com/ja1code/koolseal/util"
	"github.com/urfave/cli/v2"
)

func ExtractCommand() *cli.Command {
	command := cli.Command{
		Name:    "extract",
		Aliases: []string{"e"},
		Usage:   "extract a secrets file",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "namespace",
				Aliases:     []string{"ns"},
				Usage:       "secrets namespace",
				Required:    true,
				DefaultText: "default",
			},
			&cli.StringFlag{
				Name:     "name",
				Aliases:  []string{"n"},
				Usage:    "secrets name",
				Required: true,
			},
		},
		Action: extractAction(),
	}

	return &command
}

func extractAction() func(cCtx *cli.Context) error {
	return func(cCtx *cli.Context) error {
		secretName := strings.Split(cCtx.String("ns"), "/")
		if len(secretName) != 2 {
			fmt.Println("the secret flag should be <namespace>/<name>")
			return nil
		}

		secretsRaw, err := util.CallCmd("kubectl", "get", "secret", secretName[1], "-o", "json", "-n", secretName[0])
		if err != nil {
			fmt.Println(err.Error())
			return nil
		}

		var secrets entity.SecretsDeclaration
		err = json.Unmarshal([]byte(secretsRaw), &secrets)
		if err != nil {
			fmt.Println(err.Error())
			return nil
		}

		for key, secret := range secrets.Data {
			secrets.Data[key] = util.DecodeB64(secret)
		}

		jsonBytes, err := json.Marshal(secrets)
		if err != nil {
			fmt.Println(err.Error())
			return nil
		}

		if cCtx.Args().First() != "" {
			err = os.WriteFile(cCtx.Args().First(), jsonBytes, 0777)
			if err != nil {
				fmt.Println(err.Error())
				return nil
			}
		} else {
			fmt.Println(string(jsonBytes))
		}

		return nil
	}
}
