package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"locsvc/model"
	"log"
	"net/http"
)

var dal *model.ModelDAL

func main() {

	dal = model.OpenDB("data/locsvc.db")

	r := mux.NewRouter()

	r.HandleFunc("/", helloHandler)
	r.HandleFunc("/register", registerHandler).Methods("POST")
	r.HandleFunc("/login", loginHandler).Methods("POST")
	r.HandleFunc("/marker/create", createMarkerHandler).Methods("POST")
	r.HandleFunc("/marker/search", searchLocationHandler).Methods("POST")

	http.Handle("/", r)

	log.Println("Listening...")
	http.ListenAndServe(":3000", nil)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello World!")
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	var user model.User

	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	//hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	user.Password = string(hash)

	newUser, err := dal.CreateUser(&user)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	log.Printf("%#v", newUser)
	w.WriteHeader(http.StatusOK)
	return
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var user model.User

	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	foundUser, err := dal.GetUserByName(user.Name)
	if err != nil {
		http.Error(w, "User Not Found", 400)
		return
	}

	// Comparing the password with the hash
	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password)); err != nil {
		http.Error(w, "Invalid Password", 400)
		return
	}

	//create session
	newSession, err := dal.CreateSession(foundUser.Id)
	if err != nil {
		http.Error(w, "Unable to login", 400)
		return
	}

	retvalJson, _ := json.Marshal(newSession)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(retvalJson))

	return

}

//places a marker (text owned by a user) at a specific location
func createMarkerHandler(w http.ResponseWriter, r *http.Request) {

	var createMarker model.CreateMarker

	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&createMarker)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	//verify session
	session, err := dal.GetSessionBySessionId(createMarker.SessionId)
	if err != nil {
		http.Error(w, "Invalid Session!", 400)
		return
	}

	//place marker
	newMarker, err := dal.CreateMarker(&createMarker, session.UserId)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	retvalJson, _ := json.Marshal(newMarker)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(retvalJson))

	return

}

//searches for any markers in the vicinity (default -1 mile)
func searchLocationHandler(w http.ResponseWriter, r *http.Request) {
	var searchMarkers model.SearchMarkers

	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&searchMarkers)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	//verify session
	_, err = dal.GetSessionBySessionId(searchMarkers.SessionId)
	if err != nil {
		http.Error(w, "Invalid Session!", 400)
		return
	}

	markers, err := dal.SearchMarkersByLoc(searchMarkers.Lat, searchMarkers.Lon)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	retvalJson, _ := json.Marshal(markers)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(retvalJson))

	return

}
