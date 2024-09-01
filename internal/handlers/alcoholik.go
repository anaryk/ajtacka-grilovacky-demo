package handlers

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/anaryk/ajtacka-grilovacky-demo/internal/models"
	"github.com/anaryk/ajtacka-grilovacky-demo/web/templates"
	"github.com/rs/zerolog"
)

// AlkoholikHandler handles the registration of a new user
func AlkoholikHandler(db *models.DB, logger *zerolog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Debug().Msg("Handling /alkoholik request")

		if cookie, err := r.Cookie("alkoholik_id"); err == nil {
			alkoholikID, err := strconv.Atoi(cookie.Value)
			if err != nil {
				logger.Error().Err(err).Msg("Invalid session ID in cookie")
				clearCookie(w)
				http.Redirect(w, r, "/alkoholik", http.StatusSeeOther)
				return
			}

			if _, err := db.GetAlkoholikByID(alkoholikID); err == nil {
				logger.Debug().Msg("User already registered, redirecting to drink page")
				http.Redirect(w, r, "/drink-page", http.StatusSeeOther)
				return
			} else {
				logger.Warn().Int("id", alkoholikID).Msg("No user found with this ID, clearing cookie")
				clearCookie(w)
			}
		}

		if r.Method == "POST" {
			r.ParseMultipartForm(10 << 20)
			jmeno := r.FormValue("jmeno")
			file, _, err := r.FormFile("fotka")
			if err != nil {
				logger.Error().Err(err).Msg("Invalid image upload")
				http.Error(w, "Invalid image", 400)
				return
			}
			defer file.Close()
			fileBytes, err := ioutil.ReadAll(file)
			if err != nil {
				logger.Error().Err(err).Msg("Error reading file")
				http.Error(w, "Error reading file", 500)
				return
			}
			fotka := base64.StdEncoding.EncodeToString(fileBytes)
			id, err := db.CreateAlkoholik(jmeno, fotka)
			if err != nil {
				logger.Error().Err(err).Msg("Unable to register alcoholik")
				http.Error(w, "Unable to register", 500)
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:     "alkoholik_id",
				Value:    strconv.Itoa(id),
				Expires:  time.Now().Add(24 * time.Hour),
				Path:     "/",
				HttpOnly: true,
			})

			logger.Debug().Int("id", id).Msg("Alkoholik registered and cookie set")
			http.Redirect(w, r, "/drink-page", http.StatusSeeOther)

		} else {
			tmpl, err := templates.NewTemplate("layout.html", "alcoholik.html")
			if err != nil {
				logger.Error().Err(err).Msg("Failed to load templates")
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			if err := tmpl.ExecuteTemplate(w, "layout.html", nil); err != nil {
				logger.Error().Err(err).Msg("Template execution error")
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}

func clearCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "alkoholik_id",
		Value:    "",
		Expires:  time.Unix(0, 0),
		Path:     "/",
		HttpOnly: true,
	})
}
