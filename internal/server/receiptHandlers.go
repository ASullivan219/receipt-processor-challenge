package server

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/asullivan219/receiptProcessor/internal/models"
	"github.com/asullivan219/receiptProcessor/internal/store"
	"github.com/google/uuid"
)

type IdResponse struct {
	Id string `json:"id"`
}

// Handles A Post request to /receipts/process
func processReceipt(w http.ResponseWriter, r *http.Request, s store.Store) {

	var receipt models.Receipt
	err := json.NewDecoder(r.Body).Decode(&receipt)
	if err != nil {
		slog.Error("Error decoding response body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	validReceipt, err := receipt.ValidateReceipt()
	if err != nil {
		slog.Error("Error invalid receipt")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := uuid.NewUUID()
	if err != nil {
		slog.Error("Error making uuid")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	points := validReceipt.ScoreReceipt()
	recDetails, _ := io.ReadAll(r.Body)

	err = s.PutReceipt(
		store.NewDbReceipt(id.String(), string(recDetails), points),
	)

	if err != nil {
		slog.Error("Error storing receipt")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := IdResponse{
		Id: id.String(),
	}

	slog.Info("processed receipt",
		"id", id.String(),
		"score", points,
		"response", resp,
	)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&resp)
	return
}

type PointResponse struct {
	Points int `json:"points"`
}

// Handles a Get request to /receipts/{id}/points
func getReceiptScore(w http.ResponseWriter, r *http.Request, s store.Store) {
	id := r.PathValue("id")
	receipt, err := s.GetReceipt(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	resp := PointResponse{
		Points: receipt.Score,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&resp)
	return
}
