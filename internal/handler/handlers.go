// Package classification of Product API
//
// Documentation for Product API
//
//  Schemes: http
//  BasePath: /
//  Version: 1.0.0
//
//  Consumes:
//  - application/json
//
//  Produces:
//  - application/json
// swagger:meta

package handler

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/kartochnik010/outstaff-task/docs"
	"github.com/kartochnik010/outstaff-task/internal/config"
	"github.com/kartochnik010/outstaff-task/internal/domain/models"
	"github.com/kartochnik010/outstaff-task/internal/pkg/js"
	"github.com/kartochnik010/outstaff-task/internal/pkg/logger"
	"github.com/kartochnik010/outstaff-task/internal/service"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	client *http.Client
	s      *service.Service
	l      *logrus.Logger
	cfg    *config.Config
}

func NewHandler(client *http.Client, s *service.Service, l *logrus.Logger, cfg *config.Config) *Handler {
	return &Handler{
		client: client,
		s:      s,
		l:      l,
		cfg:    cfg,
	}
}

// @Summary Fetch and store rates
// @Description Fetch and store rates by date
// @Tags rates
// @Accept json
// @Produce json
// @Param date path string true "date. Example: '01-01-2022'"
// @Success 200 {object} js.JSON{success=bool}
// @Failure 500 {object} js.JSON{error=string}
// @Router /currency/save/{date} [get]
func (h *Handler) StoreMusic(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLoggerFromCtx(r.Context()).WithField("op", "handler.StoreMusic")
	log.Debug()

	music := models.Music{}
	if err := js.ReadJSON(w, r, &music); err != nil {
		js.WriteJSON(w, 400, js.JSON{"success": false, "err": err.Error()}, nil)
		log.WithError(err).Error("failed to read json")
		return
	}

	// validate

	// service method
	id, err := h.s.Music.StoreMusic(r.Context(), music)
	if err != nil {
		js.WriteJSON(w, 500, js.JSON{"success": false, "err": err.Error()}, nil)
		log.WithError(err).Error("error saving music")
		return
	}

	js.WriteJSON(w, 200, js.JSON{"success": true, "id": id}, nil)
}

// @Summary Get rates
// @Description Get rates by date
// @Tags rates
// @Accept json
// @Produce json
// @Param date path string true "date. Example: '01-01-2022'"
// @Param code path string false "code. Example: 'USD'"
// @Success 200 {object} js.JSON{music=[]models.Music}
// @Failure 500 {object} js.JSON{error=string}
// @Router /currency/{date}/{code} [get]
func (h *Handler) GetMusic(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLoggerFromCtx(r.Context()).WithField("op", "Handler.GetMusic")

	meta, err := getSearchMeta(r)
	if err != nil {
		log.WithError(err).Error("failed to parse metadata")
		js.WriteJSON(w, 400, js.JSON{"error": err.Error()}, nil)
		return
	}
	music, err := h.s.Music.GetMusic(r.Context(), meta)
	if err != nil {
		log.WithError(err).Error("failed to get music")
		js.WriteJSON(w, 500, js.JSON{"error": err.Error()}, nil)
		return
	}

	js.WriteJSON(w, 200, js.JSON{"music": music, "meta": meta}, nil)
}

func (h *Handler) DeleteMusic(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLoggerFromCtx(r.Context()).WithField("op", "Handler.DeleteMusic")

	vars := mux.Vars(r)
	if vars["id"] == "" {
		log.Error("parameter id is empty")
		js.WriteJSON(w, 400, js.JSON{"success": false, "error": "id is empty"}, nil)
		return
	}

	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		log.WithError(err).Error("failed to parse id")
		js.WriteJSON(w, 400, js.JSON{"success": false, "error": "failed to parse id"}, nil)
		return
	}

	err = h.s.Music.DeleteMusic(r.Context(), id)
	if err != nil {
		log.WithError(err).Error("failed to delete music")
		js.WriteJSON(w, 500, js.JSON{"success": false, "error": err.Error()}, nil)
		return
	}

	js.WriteJSON(w, 200, js.JSON{"success": true}, nil)
}

func (h *Handler) UpdateMusic(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLoggerFromCtx(r.Context()).WithField("op", "Handler.UpdateMusic")

	vars := mux.Vars(r)
	if vars["id"] == "" {
		log.Error("parameter id is empty")
		js.WriteJSON(w, 400, js.JSON{"success": false, "error": "id is empty"}, nil)
		return
	}

	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		log.WithError(err).Error("failed to parse id")
		js.WriteJSON(w, 400, js.JSON{"success": false, "error": "failed to parse id"}, nil)
		return
	}

	music := models.Music{}
	if err := js.ReadJSON(w, r, &music); err != nil {
		js.WriteJSON(w, 400, js.JSON{"success": false, "err": err.Error()}, nil)
		log.WithError(err).Error("failed to read json")
		return
	}

	music.ID = id

	err = h.s.Music.UpdateMusicByID(r.Context(), music)
	if err != nil {
		log.WithError(err).Error("failed to update music")
		js.WriteJSON(w, 500, js.JSON{"success": false, "error": err.Error()}, nil)
		return
	}

	js.WriteJSON(w, 200, js.JSON{"success": true}, nil)
}
