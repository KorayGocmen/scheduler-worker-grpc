package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type rStartJobReq struct {
	Command  string `json:"command"`
	Path     string `json:"path"`
	WorkerID string `json:"worker_id"`
}

func rPong(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "pong\n")
}

func rStartJob(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	startJobReq := rStartJobReq{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &startJobReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if success, err := startJobOnWorker(startJobReq); !success {
		http.Error(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "ok\n")
}

func createRouter() *httprouter.Router {
	router := httprouter.New()
	router.GET("/ping", rPong)
	router.POST("/start", rStartJob)
	return router
}

func api() {
	APIAddr := ":3000"
	srv := &http.Server{
		Addr:    APIAddr,
		Handler: createRouter(),
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
