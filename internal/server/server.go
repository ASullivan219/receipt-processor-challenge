package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/asullivan219/receiptProcessor/internal/store"
)

type Server struct {
	mux *http.ServeMux
	Srv *http.Server
}

func NewServer(port int, dbFile string) Server {
	s := store.NewStore(dbFile)
	serverString := fmt.Sprintf("0.0.0.0:%d", port)
	mux := http.NewServeMux()

	mux.HandleFunc(
		"GET /receipts/{id}/points",
		func(w http.ResponseWriter, r *http.Request) {
			getReceiptScore(w, r, s)
		})

	mux.HandleFunc(
		"POST /receipts/process",
		func(w http.ResponseWriter, r *http.Request) {
			processReceipt(w, r, s)
		})

	srv := &http.Server{
		Handler:      mux,
		Addr:         serverString,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	return Server{
		mux: mux,
		Srv: srv,
	}
}
