package web

import (
	"context"
	"net/http"

	"github.com/marvinmarpol/golang-boilerplate/internal/protection/entity"
	"github.com/marvinmarpol/golang-boilerplate/internal/protection/service"
	"github.com/sirupsen/logrus"

	"github.com/go-chi/render"
	"github.com/go-pg/pg/v10"
)

type Server struct {
	service service.Services
}

func NewServer(service service.Services) *Server {
	return &Server{service}
}

func (s *Server) Deidentify(w http.ResponseWriter, r *http.Request) {
	var payload interface{}

	if err := render.DecodeJSON(r.Body, &payload); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, err.Error())
		return
	}

	result, err := s.service.Deidentify(r.Context(), payload)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, err.Error())
		return
	}

	render.JSON(w, r, result)
}

func (s *Server) Reidentify(w http.ResponseWriter, r *http.Request) {
	var payload interface{}

	if err := render.DecodeJSON(r.Body, &payload); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, err.Error())
		return
	}

	result, err := s.service.Reidentify(r.Context(), payload)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, err.Error())
		return
	}

	render.JSON(w, r, result)
}

func (s *Server) GetCypher(w http.ResponseWriter, r *http.Request) {
	var payload entity.GetCypherPayload

	if err := render.DecodeJSON(r.Body, &payload); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, err.Error())
		return
	}

	result, err := s.service.GetCypher(r.Context(), payload)
	if err != nil {
		if err == pg.ErrNoRows {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, err.Error())
			return
		}
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, err.Error())
		return
	}

	render.JSON(w, r, map[string]interface{}{"value": result})
}

func (s *Server) RotateKeys(w http.ResponseWriter, r *http.Request) {
	var payload entity.RotatePayload
	if err := render.DecodeJSON(r.Body, &payload); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, err.Error())
	}

	go func() {
		_, err := s.service.RotateKeys(context.Background(), payload)
		if err != nil {
			logrus.WithContext(context.Background()).WithField("err", err).Error("Failed to rotate keys")
			return
		}
	}()

	render.JSON(w, r, map[string]interface{}{"result": http.StatusText(http.StatusOK)})
}
