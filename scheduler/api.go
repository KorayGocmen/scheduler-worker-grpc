package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func rPong(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "pong\n")
}

func rStartJob(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	startJobReq := apiStartJobReq{}

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

	ok, startErr, jobID := startJobOnWorker(startJobReq)
	if !ok {
		http.Error(w, startErr, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(apiStartJobRes{JobID: jobID})
}

func rStopJob(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	stopJobReq := apiStopJobReq{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &stopJobReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if ok, err := stopJobOnWorker(stopJobReq); !ok {
		http.Error(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(apiStopJobRes{Success: true})
}

func rQueryJob(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	queryJobReq := apiQueryJobReq{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &queryJobReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ok, queryError, done := queryJobOnWorker(queryJobReq)
	if !ok {
		http.Error(w, queryError, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(apiQueryJobRes{Done: done})
}

func createRouter() *httprouter.Router {
	router := httprouter.New()

	router.GET("/ping", rPong)
	router.POST("/start", rStartJob)
	router.POST("/stop", rStopJob)
	router.POST("/query", rQueryJob)

	return router
}

func api() {
	srv := &http.Server{
		Addr:    config.HTTPServer.Addr,
		Handler: createRouter(),
	}

	log.Println("HTTP Server listening on", config.HTTPServer.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
