package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/anaryk/ajtacka-grilovacky-demo/internal/models"
	"github.com/anaryk/ajtacka-grilovacky-demo/pkg/websocket"
	"github.com/anaryk/ajtacka-grilovacky-demo/web/templates"
	"github.com/rs/zerolog"
)

type Response struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func DrinkHandler(db *models.DB, hub *websocket.Hub, logger *zerolog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			cookie, err := r.Cookie("alkoholik_id")
			if err != nil {
				logger.Error().Err(err).Msg("Session not found")
				http.Redirect(w, r, "/alkoholik", http.StatusSeeOther)
				return
			}

			alkoholikID, err := strconv.Atoi(cookie.Value)
			if err != nil {
				logger.Error().Err(err).Msg("Invalid session ID")
				http.Redirect(w, r, "/alkoholik", http.StatusSeeOther)
				return
			}

			alkoholik, err := db.GetAlkoholikByID(alkoholikID)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to load user data. It seems user not exist. Redirection to /alkoholik to pass registration process")
				clearCookie(w)
				http.Redirect(w, r, "/alkoholik", http.StatusSeeOther)
				return
			}

			data := map[string]interface{}{
				"Jmeno": alkoholik.Jmeno,
				"Fotka": alkoholik.Fotka,
			}

			tmpl, err := templates.NewTemplate("layout.html", "drink.html")
			if err != nil {
				logger.Error().Err(err).Msg("Failed to load templates")
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			if err := tmpl.ExecuteTemplate(w, "layout.html", data); err != nil {
				logger.Error().Err(err).Msg("Failed to render the drink page")
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		} else if r.Method == http.MethodPost {
			cookie, err := r.Cookie("alkoholik_id")
			if err != nil {
				logger.Error().Err(err).Msg("Session not found")
				jsonResponse(w, http.StatusUnauthorized, "error", "Session not found")
				return
			}

			alkoholikID, err := strconv.Atoi(cookie.Value)
			if err != nil {
				logger.Error().Err(err).Msg("Invalid session ID")
				jsonResponse(w, http.StatusBadRequest, "error", "Invalid session ID")
				return
			}

			lastAdded, err := db.GetLastDrinkTime(alkoholikID)
			if err != nil {
				logger.Error().Err(err).Msg("Error retrieving last drink time")
				jsonResponse(w, http.StatusInternalServerError, "error", "Error retrieving last drink time")
				return
			}

			if time.Since(lastAdded) < 5*time.Minute {
				logger.Info().Msg("Attempt to add drink too soon")
				jsonResponse(w, http.StatusTooManyRequests, "error", "You can add a drink only every 5 minutes")
				return
			}

			drinkType := r.FormValue("type")
			if err := db.AddDrink(alkoholikID, drinkType); err != nil {
				logger.Error().Err(err).Msg("Error adding drink")
				jsonResponse(w, http.StatusInternalServerError, "error", "Error adding drink")
				return
			}

			logger.Info().Str("drinkType", drinkType).Int("alkoholikID", alkoholikID).Msg("Drink added successfully")

			alkoholici, err := db.GetAllAlkoholici()
			if err != nil {
				logger.Error().Err(err).Msg("Error fetching updated alcoholics")
				jsonResponse(w, http.StatusInternalServerError, "error", "Error fetching updated alcoholics")
				return
			}

			hub.BroadcastUpdate(alkoholici)
			jsonResponse(w, http.StatusOK, "success", "Drink added successfully")
		}
	}
}

func jsonResponse(w http.ResponseWriter, statusCode int, respType string, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(Response{Type: respType, Message: message})
}
