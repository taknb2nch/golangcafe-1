package main

import (
	"fmt"
	"log"
	"net/http"
)

func init() {
    http.HandleFunc("/", errorHandler(betterHandler))
}

func errorHandler(f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        err := f(w, r)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            log.Printf("handling %q: %v", r.RequestURI, err)
        }
    }
}

func betterHandler(w http.ResponseWriter, r *http.Request) error {
    if err := doThis(); err != nil {
        return fmt.Errorf("doing this: %v", err)
    }

    if err := doThat(); err != nil {
        return fmt.Errorf("doing that: %v", err)
    }
    return nil
}

func main() {
	http.ListenAndServe(":12345", nil)
}

type HogeError struct {
	Message string
}

func (e HogeError) Error() string {
	return e.Message
}

func doThis() error {
//	return HogeError{ Message: "Fire Error doThis"}
	return nil
}

func doThat() error {
	return HogeError{ Message: "Fire Error doThat"}
//	return nil
}