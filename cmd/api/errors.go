package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error){
	log.Printf("internal error: %s path: %s", r.Method, r.URL.Path, err)

	writeJSONError(w, http.StatusInternalServerError, "the server encountered a problem")
}

func (app *application) badRequestError(w http.ResponseWriter, r *http.Request, err error){
	log.Printf("bad request error: %s path: %s", r.Method, r.URL.Path, err)

	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundError(w http.ResponseWriter, r *http.Request, err error){
	log.Printf("not found error: %s path: %s", r.Method, r.URL.Path, err)

	writeJSONError(w, http.StatusNotFound, "not found")
}

func (app *application) conflictError(w http.ResponseWriter, r *http.Request, err error){
	log.Printf("resource already exits error: %s path: %s", r.Method, r.URL.Path, err)

	writeJSONError(w, http.StatusConflict, "already exists")
}