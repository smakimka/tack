package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/smakimka/tack/internal/model"
	"github.com/smakimka/tack/internal/monitor/storage"
)

type UpdateHandler struct {
	s storage.Storage
}

func NewUpdateHandler(s storage.Storage) *UpdateHandler {
	return &UpdateHandler{s}
}

func (h *UpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	log.Info().Msg("processing update request")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Err(err).Msg("error reading body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var data model.SendData
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Err(err).Msg("error unmarshaling json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.s.Update(data)
	log.Info().Msg("update succesful")
	w.WriteHeader(http.StatusOK)
}
