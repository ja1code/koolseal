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

func CreateCommand() *cli.Command {
	command := cli.Command{
		Name:    "create",
		Aliases: []string{"c"},
		Usage:   "create a new secrets file",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "cert",
				Usage:    "The certificate used to create secrets",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "name",
				Usage:    "The name of the secrets to create",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "namespace",
				Usage:    "New secrets namespace",
				Required: true,
			},
			&cli.StringFlag{
				Name:        "destination",
				Usage:       "The destination dir for the secrets file",
				DefaultText: "./",
				Required:    true,
			},
			&cli.BoolFlag{
				Name:  "publish",
				Usage: "commit new secrets",
			},
			&cli.StringFlag{
				Name:     "file",
				Usage:    "file with new values",
				Required: true,
			},
		},
		Action: createAction(),
	}

	return &command
}

func createAction() func(cCtx *cli.Context) error {
	return func(cCtx *cli.Context) error {
		secretName := strings.Split(cCtx.String("secrets"), "/")
		if len(secretName) != 2 {
			fmt.Println("the secret flag should be <namespace>/<name>")
			return nil
		}

		secrets := entity.SecretsDeclaration{
			ApiVersion: "",
			Kind:       "",
			Metadata: entity.Metadata{
				Name:      secretName[1],
				Namespace: secretName[0],
			},
			Type: "Opaque",
			Data: map[string]string{},
		}

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
