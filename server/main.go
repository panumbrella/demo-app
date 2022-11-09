package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

func NewHTTPServer(addr string) *http.Server {
	r := mux.NewRouter()
	r.HandleFunc("/", handleFunc).Methods("POST")
	return &http.Server{
		Addr:    addr,
		Handler: r,
	}
}

type FuncRequest struct {
	Input int `json:"input"`
}

type FuncResponse struct {
	Output float64 `json:"output"`
}

func handleFunc(w http.ResponseWriter, r *http.Request) {
	var req FuncRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	output, err := processFunc(req.Input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res := FuncResponse{Output: output}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func processFunc(input int) (float64, error) {
	var wg sync.WaitGroup
	now := time.Now()
	for i := 0; i < input; i++ {
		wg.Add(1)
		go partialProcessFunc(&wg)
	}
	wg.Wait()
	timetaken := time.Since(now).Seconds()
	return timetaken, nil
}

func partialProcessFunc(wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(1 * time.Second)
}

func main() {
	srv := NewHTTPServer(":8080")
	log.Fatal(srv.ListenAndServe())
}
