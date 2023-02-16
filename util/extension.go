package util

import "strings"

var magicTable = map[string]string{
	"\xff\xd8\xff":      "image/jpeg",
	"\x89PNG\r\n\x1a\n": "image/png",
	"GIF87a":            "image/gif",
	"GIF89a":            "image/gif",
}

func MimeFromIncipit(incipit []byte) string {
	incipitStr := []byte(incipit)
	for magic, mime := range magicTable {
		if strings.HasPrefix(string(incipitStr), magic) {
			return mime
		}
	}

	return ""
}
