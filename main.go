package main

import (
	"net/http"
	"os"

	"github.com/anaryk/ajtacka-grilovacky-demo/internal/handlers"
	"github.com/anaryk/ajtacka-grilovacky-demo/internal/models"
	"github.com/anaryk/ajtacka-grilovacky-demo/pkg/websocket"
	"github.com/rs/zerolog"
)

func main() {
	log := zerolog.New(os.Stdout).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	db, err := models.InitDB("root:Start123@tcp(10.181.21.90:55000)/alkoholapp")
	if err != nil {
		log.Fatal().Err(err).Msg("Database connection failed")
	}
	defer db.Close()

	hub := websocket.NewHub(&log)
	go hub.Run()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	http.HandleFunc("/stats", handlers.StatsHandler(db, hub, &log))
	http.HandleFunc("/drink-page", handlers.DrinkHandler(db, hub, &log))
	http.HandleFunc("/alkoholik", handlers.AlkoholikHandler(db, &log))
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(hub, w, r, &log)
	})

	log.Info().Msg("Server is running on :8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start the server")
	}
}
