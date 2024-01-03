package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ja1code/koolseal/entity"
	"github.com/ja1code/koolseal/util"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

func CreateCommand() *cli.Command {
	command := cli.Command{
		Name:    "create",
		Aliases: []string{"c"},
		Usage:   "create a new secrets file",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "cert",
				Aliases:  []string{"c"},
				Usage:    "The certificate used to create secrets",
				Required: true,
			},
			&cli.StringFlag{
				Name:        "namespace",
				Aliases:     []string{"ns"},
				Usage:       "Namespace to add the new secrets",
				DefaultText: "default",
				Required:    true,
			},
			&cli.StringFlag{
				Name:     "name",
				Aliases:  []string{"n"},
				Usage:    "Secret name",
				Required: true,
			},
			&cli.BoolFlag{
				Name:    "publish",
				Aliases: []string{"p"},
				Usage:   "Commit new secrets to github",
			},
			&cli.StringFlag{
				Name:     "file",
				Aliases:  []string{"f"},
				Usage:    "File with new values",
				Required: true,
			},
		},
		Action: createAction(),
	}

	return &command
}

func createAction() func(cCtx *cli.Context) error {
	return func(cCtx *cli.Context) error {
		secretName := cCtx.String("name")

		secretNamespace := cCtx.String("namespace")

		secrets := entity.SecretsDeclaration{
			ApiVersion: "",
			Kind:       "",
			Metadata: entity.Metadata{
				Name:      secretName,
				Namespace: secretNamespace,
			},
			Type: "Opaque",
			Data: map[string]string{},
		}

		valueMapRaw, err := os.ReadFile(cCtx.String("file"))
		if err != nil {
			fmt.Println("Error while reading file with new values", err.Error())
			return nil
		}

		var valueMap map[string]string
		err = json.Unmarshal(valueMapRaw, &valueMap)
		if err != nil {
			fmt.Println("Error parsing new values", err.Error())
			return nil
		}

		for key, value := range valueMap {
			secrets.Data[key] = value
		}

		eYaml, err := yaml.Marshal(&secrets)
		if err != nil {
			fmt.Println("Internal YAML error", err.Error())
			return nil
		}

		err = os.WriteFile("temp.yaml", eYaml, 0777)
		if err != nil {
			fmt.Println("Error when writting temporary files", err.Error())
			return nil
		}

		encryptedSecrets, err := util.CallCmd("kubeseal", "--cert", cCtx.String("cert"), "-o", "yaml", "-f", "temp.yaml")
		if err != nil {
			fmt.Println("Error when calling kubeseal", err.Error())
			return nil
		}

		err = os.Remove("temp.yaml")
		if err != nil {
			fmt.Println("Error deleting temporary files", err.Error())
			return nil
		}

		err = os.WriteFile(cCtx.Args().Get(0), []byte(encryptedSecrets), 0755)
		if err != nil {
			fmt.Println("Error writing secrets file", err.Error())
			return nil
		}

		if cCtx.Bool("publish") {
			err = util.PublishChanges(cCtx.Args().First(), fmt.Sprintf("\"KoolSeal: updating %s secrets\"", fmt.Sprintf("%s/%s", secretNamespace, secretName)))
			if err != nil {
				fmt.Println("Error when publishing to Github", err.Error())
				return nil
			}
		}

		return nil
	}
}
