package http

import (
	"context"
	"encoding/json"
	"fmt"
	"hw-6/internal/models"
	"hw-6/internal/store"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type Server struct {
	ctx         context.Context
	idleConnsCh chan struct{}
	store       store.Store

	Address string
}

func NewServer(ctx context.Context, address string, store store.Store) *Server {
	return &Server{
		ctx:         ctx,
		idleConnsCh: make(chan struct{}),
		store:       store,

		Address: address,
	}
}

func (s *Server) basicHandler() chi.Router {
	r := chi.NewRouter()

	r.Post("/transactions", func(w http.ResponseWriter, r *http.Request) {
		transaction := new(models.Transaction)
		if err := json.NewDecoder(r.Body).Decode(transaction); err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		s.store.Create(r.Context(), transaction)
	})
	r.Get("/transactions", func(w http.ResponseWriter, r *http.Request) {
		transactions, err := s.store.All(r.Context())
		if err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		render.JSON(w, r, transactions)
	})
	r.Get("/transactions/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		transaction, err := s.store.ByID(r.Context(), id)
		if err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		render.JSON(w, r, transaction)
	})
	r.Put("/transactions", func(w http.ResponseWriter, r *http.Request) {
		transaction := new(models.Transaction)
		if err := json.NewDecoder(r.Body).Decode(transaction); err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		s.store.Update(r.Context(), transaction)
	})
	r.Delete("/transactions/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		s.store.Delete(r.Context(), id)
	})

	// ACCOUNTS OPERATIONS
	r.Post("/accounts", func(w http.ResponseWriter, r *http.Request) {
		body := new(models.CreateAccountBody)
		if err := json.NewDecoder(r.Body).Decode(body); err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		s.store.CreateAccount(r.Context(), body.Name)
	})
	r.Get("/accounts", func(w http.ResponseWriter, r *http.Request) {
		accounts, err := s.store.GetAccounts(r.Context())
		if err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		render.JSON(w, r, accounts)
	})

	r.Get("/accounts/{accountId}", func(w http.ResponseWriter, r *http.Request) {
		accountIdStr := chi.URLParam(r, "accountId")
		id, err := strconv.Atoi(accountIdStr)
		if err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}
		account, err := s.store.GetAccount(r.Context(), id)
		if err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		render.JSON(w, r, account)
	})

	r.Delete("/accounts/{accountId}", func(w http.ResponseWriter, r *http.Request) {
		accountIdStr := chi.URLParam(r, "accountId")
		id, err := strconv.Atoi(accountIdStr)
		if err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}
		s.store.DeleteAccount(r.Context(), id)
	})

	r.Get("/accounts/{accountId}/amount", func(w http.ResponseWriter, r *http.Request) {
		accountIdStr := chi.URLParam(r, "accountId")
		id, err := strconv.Atoi(accountIdStr)
		if err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}
		amount, err := s.store.GetCurrentAmountOfAccount(r.Context(), id)
		amountBody := &models.CurrentAmount{Amount: amount}
		if err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		render.JSON(w, r, amountBody)
	})

	// ACCOUNT TRANSACTIONS
	r.Post("/accounts/{accountId}/transactions", func(w http.ResponseWriter, r *http.Request) {
		accountIdStr := chi.URLParam(r, "accountId")
		id, err := strconv.Atoi(accountIdStr)
		if err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		transaction := new(models.Transaction)
		if err := json.NewDecoder(r.Body).Decode(transaction); err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		s.store.CreateTransaction(r.Context(), transaction, id)
	})

	r.Get("/accounts/{accountId}/transactions", func(w http.ResponseWriter, r *http.Request) {
		accountIdStr := chi.URLParam(r, "accountId")
		id, err := strconv.Atoi(accountIdStr)
		if err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}
		accounts, err := s.store.GetTransactionsByAccount(r.Context(), id)
		if err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		render.JSON(w, r, accounts)
	})

	r.Put("/transactions", func(w http.ResponseWriter, r *http.Request) {
		transaction := new(models.Transaction)
		if err := json.NewDecoder(r.Body).Decode(transaction); err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		s.store.Update(r.Context(), transaction)
	})

	r.Put("/accounts/{accountId}/transactions", func(w http.ResponseWriter, r *http.Request) {
		accountIdStr := chi.URLParam(r, "accountId")
		accountId, err := strconv.Atoi(accountIdStr)
		if err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		transaction := new(models.Transaction)
		if err := json.NewDecoder(r.Body).Decode(transaction); err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		s.store.UpdateTransaction(r.Context(), transaction, accountId)
	})

	return r
}

func (s *Server) Run() error {
	srv := &http.Server{
		Addr:         s.Address,
		Handler:      s.basicHandler(),
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 30,
	}
	go s.ListenCtxForGT(srv)

	log.Println("[HTTP] Server running on", s.Address)
	return srv.ListenAndServe()
}

func (s *Server) ListenCtxForGT(srv *http.Server) {
	<-s.ctx.Done()

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Printf("[HTTP] Got err while shutting down^ %v", err)
	}

	log.Println("[HTTP] Proccessed all idle connections")
	close(s.idleConnsCh)
}

func (s *Server) WaitForGracefulTermination() {
	<-s.idleConnsCh
}
