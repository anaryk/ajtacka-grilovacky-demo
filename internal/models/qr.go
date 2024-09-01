package models

import (
	"encoding/base64"

	"github.com/skip2/go-qrcode"
)

func GenerateQRCode(url string) string {
	var png []byte
	png, err := qrcode.Encode(url, qrcode.Medium, 256)
	if err != nil {
		return ""
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(png)
}
