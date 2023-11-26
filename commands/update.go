package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/ja1code/koolseal/entity"
	"github.com/ja1code/koolseal/util"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

func UpdateCommand() *cli.Command {
	command := cli.Command{
		Name:    "update",
		Aliases: []string{"u"},
		Usage:   "update a pre-existent secrets file",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "cert",
				Usage:    "The certificate used to update secrets",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "secrets",
				Usage:    "The namespace/name of the secrets to update",
				Required: true,
			},
			&cli.BoolFlag{
				Name:  "publish",
				Usage: "commit updates",
			},
			&cli.StringFlag{
				Name:  "file",
				Usage: "file with new values",
			},
			&cli.StringFlag{
				Name:  "key",
				Usage: "new secret key",
			},
			&cli.StringFlag{
				Name:  "value",
				Usage: "new secret value",
			},
		},
		Action: updateAction(),
	}

	return &command
}

func updateAction() func(cCtx *cli.Context) error {
	return func(cCtx *cli.Context) error {

		err := validateCall(cCtx)
		if err != nil {
			fmt.Println(err.Error())
			return nil
		}

		secretName := strings.Split(cCtx.String("secrets"), "/")
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

		if cCtx.String("key") != "" {
			key := cCtx.String("key")
			value := cCtx.String("value")
			secrets.Data[key] = value
		}

		if cCtx.String("file") != "" {
			valueMapRaw, err := os.ReadFile(cCtx.String("file"))
			if err != nil {
				fmt.Println(err.Error())
				return nil
			}

			var valueMap map[string]string
			err = json.Unmarshal(valueMapRaw, &valueMap)
			if err != nil {
				fmt.Println(err.Error())
				return nil
			}

			for key, value := range valueMap {
				secrets.Data[key] = value
			}
		}

		eYaml, err := yaml.Marshal(&secrets)
		if err != nil {
			fmt.Println(err.Error())
			return nil
		}

		err = os.WriteFile("temp.yaml", eYaml, 0777)
		if err != nil {
			fmt.Println(err.Error())
			return nil
		}

		encryptedSecrets, err := util.CallCmd("kubeseal", "--cert", cCtx.String("cert"), "-o", "yaml", "-f", "temp.yaml")
		if err != nil {
			fmt.Println(err.Error())
			return nil
		}

		err = os.Remove("temp.yaml")
		if err != nil {
			fmt.Println(err.Error())
			return nil
		}

		err = os.WriteFile(cCtx.Args().Get(0), []byte(encryptedSecrets), 0755)
		if err != nil {
			fmt.Println(err.Error())
			return nil
		}

		if cCtx.Bool("publish") {
			err = util.PublishChanges(cCtx.Args().First(), fmt.Sprintf("\"KoolSeal: updating %s secrets\"", strings.Join(secretName, "/")))
			if err != nil {
				fmt.Println(err.Error())
				return nil
			}
		}

		return nil
	}
}

func validateCall(cCtx *cli.Context) error {
	if cCtx.String("cert") == "" {
		return fmt.Errorf("the cert flag is required")
	}

	if cCtx.String("secrets") == "" {
		return fmt.Errorf("the secrets flag is required")

	}

	if cCtx.String("file") == "" && cCtx.String("key") == "" {
		return fmt.Errorf("a new secret key or a file with multiple needs to be provided")
	}

	if cCtx.Args().Get(0) == "" {
		return fmt.Errorf("a destination must be provided")
	}

	return nil
}
