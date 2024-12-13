package handler

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func Routes(h *Handler) *mux.Router {
	r := mux.NewRouter()
	r.Use(h.AssignLoggerMiddleware, h.rateLimit, h.WriteToConsole)

	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://localhost:%v/swagger/doc.json", h.cfg.Port)),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)

	r.HandleFunc("/music", h.GetMusic).Methods("GET")
	r.HandleFunc("/music/save", h.StoreMusic).Methods("POST")
	r.HandleFunc("/music/delete/{id}", h.DeleteMusic).Methods("DELETE")
	r.HandleFunc("/music/edit/{id}", h.UpdateMusic).Methods("PUT")

	return r
}
