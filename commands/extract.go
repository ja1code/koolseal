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
				Name:     "ns",
				Usage:    "New secrets namespace",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "file",
				Usage:    "file with new values",
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

		err = os.WriteFile(cCtx.String("file"), jsonBytes, 0777)
		if err != nil {
			fmt.Println(err.Error())
			return nil
		}

		return nil
	}
}
