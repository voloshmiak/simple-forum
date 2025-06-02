package handler

import (
	"log/slog"
	"net/http"
)

func InternalServer(logger *slog.Logger, rw http.ResponseWriter, msg string, err error) {
	logger.Error(msg, "error", err)
	http.Error(rw, msg, http.StatusInternalServerError)
}
