package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/smrrazavian/url-shortener/internal/models"
	"github.com/smrrazavian/url-shortener/pkg/idgen"
)

var (
	urlStore = make(map[string]models.URL, 0)
	mutex    = sync.Mutex{}
)

// TODO: write docs
type SaveRequest struct {
	URL models.CustomURL `json:"URL"`
	TTL *int             `json:"TTL,omitempty"`
}

// TODO: Complete docs
func SaveURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SaveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON Body", http.StatusBadRequest)
		return
	}

	var ttlSeconds int
	if req.TTL == nil {
		ttlSeconds = 5 * 60
	} else {
		ttlSeconds = *req.TTL
	}
	expiresAt := time.Now().Add(time.Duration(ttlSeconds) * time.Second)

	shortId, err := idgen.GenerateID()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()
	urlStore[shortId] = models.URL{URL: req.URL, ExpiresAt: expiresAt}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"id": shortId})
}

// TODO: write docs
func GetURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	id := r.URL.Path[1:]
	urlData, ok := urlStore[id]
	if !ok {
		http.Error(w, "ID not found", http.StatusNotFound)
		return
	}

	if urlData.ExpiresAt.Before(time.Now()) {
		http.Error(w, "URL is Gone", http.StatusGone)
		return
	}

	if urlData.URL.IsNil() {
		http.Error(w, "Invalid Stored URL", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, urlData.URL.String(), http.StatusMovedPermanently)
}

func StoreToFile(filename string) error {
	mutex.Lock()
	defer mutex.Unlock()

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(urlStore)
}

func LoadFromFile(filename string) error {
	mutex.Lock()
	defer mutex.Unlock()

	// https://stackoverflow.com/questions/1821811/how-to-read-write-from-to-a-file-using-go
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()

	return json.NewDecoder(file).Decode(&urlStore)
}
