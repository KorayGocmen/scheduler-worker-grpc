package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	contentTypeHeader     = "Content-Type"
	applicationJSONHeader = "application/json"
)

func apiStartJob(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	startJobReq := apiStartJobReq{}

	w.Header().Set(contentTypeHeader, applicationJSONHeader)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(apiError{Error: err.Error()})
		return
	}

	err = json.Unmarshal(body, &startJobReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(apiError{Error: err.Error()})
		return
	}

	jobID, err := startJobOnWorker(startJobReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(apiError{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(apiStartJobRes{JobID: jobID})
}

func apiStopJob(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	stopJobReq := apiStopJobReq{}

	w.Header().Set(contentTypeHeader, applicationJSONHeader)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(apiError{Error: err.Error()})
		return
	}

	err = json.Unmarshal(body, &stopJobReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(apiError{Error: err.Error()})
		return
	}

	if err := stopJobOnWorker(stopJobReq); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(apiError{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(apiStopJobRes{Success: true})
}

func apiQueryJob(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	queryJobReq := apiQueryJobReq{}

	w.Header().Set(contentTypeHeader, applicationJSONHeader)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(apiError{Error: err.Error()})
		return
	}

	err = json.Unmarshal(body, &queryJobReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(apiError{Error: err.Error()})
		return
	}

	jobDone, jobError, jobErrorText, err := queryJobOnWorker(queryJobReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(apiError{Error: err.Error()})
		return
	}

	queryJobRes := apiQueryJobRes{
		Done:      jobDone,
		Error:     jobError,
		ErrorText: jobErrorText,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(queryJobRes)
}

func createRouter() *httprouter.Router {
	router := httprouter.New()

	router.POST("/start", apiStartJob)
	router.POST("/stop", apiStopJob)
	router.POST("/query", apiQueryJob)

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
