package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/vlamitin/secrets-provider/internal/persistence"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func InitAndListen(serverHost string, serverPort int) {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeGetHandler).Methods("GET")
	r.HandleFunc("/", HomePostHandler).Methods("POST")

	r.HandleFunc("/{secret}", SecretGetHandler).Methods("GET")
	r.HandleFunc("/{secret}", SecretPostHandler).Methods("POST")
	r.HandleFunc("/{secret}", SecretDeleteHandler).Methods("DELETE")

	r.NotFoundHandler = http.HandlerFunc(NotFoundHandler)

	addr := fmt.Sprintf("%s:%d", serverHost, serverPort)

	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	fmt.Printf("ListenAndServe at %s ...\n", addr)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	srv.Shutdown(ctx)

	log.Println("shutting down")
	os.Exit(0)
}

func HomeGetHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Success!")
}

func HomePostHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Success POST!")
}

func SecretGetHandler(w http.ResponseWriter, r *http.Request) {
	cryptKey := r.Header.Get("X-Crypt-Key")
	if !persistence.CheckCryptKey(cryptKey) {
		http.Error(w, "bad crypt key provided", http.StatusForbidden)
		return
	}

	params := mux.Vars(r)

	secret, notFoundErr, decryptErr := persistence.GetSecret(params["secret"])

	if notFoundErr != nil {
		http.Error(w, fmt.Sprintf("%v", notFoundErr), http.StatusNotFound)
		return
	}

	if decryptErr != nil {
		http.Error(w, fmt.Sprintf("%v", decryptErr), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, secret)
}

type SecretPostRequest struct {
	Secret string `json:"secret"`
}

func SecretPostHandler(w http.ResponseWriter, r *http.Request) {
	cryptKey := r.Header.Get("X-Crypt-Key")
	if !persistence.CheckCryptKey(cryptKey) {
		http.Error(w, "bad crypt key provided", http.StatusForbidden)
		return
	}

	params := mux.Vars(r)

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var sr SecretPostRequest
	err := dec.Decode(&sr)
	if err != nil {
		http.Error(w, fmt.Sprintf("bad json: %v", err), http.StatusBadRequest)
		return
	}

	encryptErr := persistence.SetSecret(params["secret"], sr.Secret)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", encryptErr), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Success POST!")
}

func SecretDeleteHandler(w http.ResponseWriter, r *http.Request) {
	cryptKey := r.Header.Get("X-Crypt-Key")
	if !persistence.CheckCryptKey(cryptKey) {
		http.Error(w, "bad crypt key provided", http.StatusForbidden)
		return
	}

	params := mux.Vars(r)

	persistence.RemoveSecret(params["secret"])

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Success POST!")
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "Not found! (404)")
}
