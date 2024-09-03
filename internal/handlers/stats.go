package handlers

import (
	"encoding/base64"
	"html/template"
	"net/http"

	"github.com/anaryk/ajtacka-grilovacky-demo/internal/models"
	"github.com/anaryk/ajtacka-grilovacky-demo/pkg/websocket"
	"github.com/anaryk/ajtacka-grilovacky-demo/web/templates"
	"github.com/rs/zerolog"
	"github.com/skip2/go-qrcode"
)

func StatsHandler(db *models.DB, hub *websocket.Hub, logger *zerolog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Debug().Msg("Fetching all alcoholics for stats view")

		alkoholici, err := db.GetAllAlkoholici()
		if err != nil {
			logger.Error().Err(err).Msg("Failed to fetch alcoholics")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		url := isHTTPS(r) + r.Host + "/alkoholik"
		qr, err := qrcode.New(url, qrcode.Medium)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to generate QR code")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		qrPng, err := qr.PNG(256)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to encode QR code as PNG")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		qrBase64 := base64.StdEncoding.EncodeToString(qrPng)
		qrDataURL := template.URL("data:image/png;base64," + qrBase64)

		tmpl, err := templates.NewTemplate("layout.html", "stats.html")
		if err != nil {
			logger.Error().Err(err).Msg("Failed to load templates")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if err := tmpl.ExecuteTemplate(w, "layout.html", map[string]interface{}{
			"Alkoholici": alkoholici,
			"QRCode":     qrDataURL,
			"Content":    "stats.html",
		}); err != nil {
			logger.Error().Err(err).Msg("Failed to execute template")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}

func isHTTPS(r *http.Request) string {
	if r.TLS != nil {
		return "https://"
	} else {
		return "http://"
	}

}
