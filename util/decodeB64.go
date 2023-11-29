package util

import "encoding/base64"

func DecodeB64(s string) string {
	o, _ := base64.StdEncoding.DecodeString(s)
	return string(o)
}
