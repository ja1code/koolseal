package entity

type EncryptedSecretsDeclaration struct {
	ApiVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Metadata   Metadata `yaml:"metadata"`
	Spec       struct {
		EncryptedData map[string]string `yaml:"encryptedData"`
		Template      Metadata          `yaml:"template"`
	} `yaml:"spec"`
}

type Metadata struct {
	CreationTimestamp interface{} `yaml:"creationTimestamp" json:"creationTimestamp"`
	Name              string      `yaml:"name" json:"name"`
	Namespace         string      `yaml:"namespace" json:"namespace"`
}

type SecretsDeclaration struct {
	ApiVersion string            `json:"apiVersion" yaml:"apiVersion"`
	Kind       string            `json:"kind" yaml:"kind"`
	Metadata   Metadata          `json:"metadata" yaml:"metadata"`
	Type       string            `json:"type" yaml:"type"`
	Data       map[string]string `json:"data" yaml:"stringData"`
}
