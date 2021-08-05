package views

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/tnyie/journaler-api/models"
	"github.com/tnyie/journaler-api/util"
)

func GetEntry(w http.ResponseWriter, r *http.Request) {
	var entry models.Entry
	entry.ID = chi.URLParam(r, "id")

	err := entry.Get()
	if err != nil {
		log.Println("Error getting entry\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if entry.OwnerID != util.GetUserID(r) {
		w.WriteHeader(http.StatusNotFound)
		log.Println(fmt.Errorf("user not authorized to access entry"))
		return
	}

	encoded, err := json.Marshal(&entry)
	if err != nil {
		log.Println("Error marshalling entry\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	respondJSON(w, encoded, http.StatusOK)
}

func PostEntry(w http.ResponseWriter, r *http.Request) {
	var entry models.Entry
	bd, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Body could not be read\n", err)
		return
	}

	err = json.Unmarshal(bd, &entry)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Body could not unmarshal to entry\n", err)
		return
	}

	entry.ID = ""
	entry.OwnerID = util.GetUserID(r)

	err = entry.Create()
	if err != nil {
		log.Println("Error creating entry\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	parentJournal := &models.Journal{
		ID: entry.JournalID,
	}

	err = parentJournal.AddEntry(entry.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Couldn't add entry to journal\n", err)
		return
	}

	encoded, err := json.Marshal(entry)
	if err != nil {
		log.Println("Couldn't marshal entry to json\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	respondJSON(w, encoded, http.StatusCreated)
}
