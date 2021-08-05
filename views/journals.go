package views

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/tnyie/journaler-api/middleware"
	"github.com/tnyie/journaler-api/models"
)

func GetOwnJournals(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.AuthCtx{}).(string)

	journals, err := models.GetOwnJournals(userID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println("Couldn't get user journals\n", err)
		return
	}

	encoded, err := json.Marshal(journals)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Couldn't marshal journals to json\n", err)
		return
	}

	respondJSON(w, encoded, http.StatusOK)
}

func GetJournalInfo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	journal := &models.Journal{
		ID: id,
	}

	journal.Get()

	children := []*models.Journal{}
	entries := []*models.Entry{}

	for _, childID := range journal.Children {
		childJournal := &models.Journal{
			ID: childID,
		}

		childJournal.Get()
		children = append(children, childJournal)
	}

	for _, entryID := range journal.Entries {
		entry := &models.Entry{
			ID: entryID,
		}

		entry.Get()
		entries = append(entries, entry)
	}

	data := make(map[string]interface{})

	data["journal"] = journal
	data["children"] = children
	data["entries"] = entries

	encoded, err := json.Marshal(&data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Couldn't marshal json\n", err)
		return
	}

	respondJSON(w, encoded, http.StatusOK)
}

func CreateJournal(w http.ResponseWriter, r *http.Request) {
	var journal models.Journal
	bd, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Body could not be read\n", err)
		return
	}

	err = json.Unmarshal(bd, &journal)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Body could not unmarshal to journal\n", err)
		return
	}

	journal.ID = ""
	journal.OwnerID = r.Context().Value(middleware.AuthCtx{}).(string)

	journal.Create()

	parentJournal := &models.Journal{
		ID: journal.ParentID,
	}

	err = parentJournal.AddChild(journal.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Couldn't add child to parent journal\n", err)
		return
	}

	encoded, err := json.Marshal(journal)
	if err != nil {
		log.Println("Couldn't marshal journal to json\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	respondJSON(w, encoded, http.StatusCreated)
}
