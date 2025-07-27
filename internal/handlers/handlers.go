package handlers

import (
	"crypto/rand"
	"encoding/json"
	"math/big"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jackkayser2005/learn_go/internal/store"
	"github.com/jackkayser2005/learn_go/internal/util"
)

type Handlers struct {
	Store *store.MemStore
}

func (h Handlers) Shorten(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var body store.UserLink
	if err := dec.Decode(&body); err != nil {
		util.WriteError(w, http.StatusBadRequest, "bad json")
		return
	}
	if strings.TrimSpace(body.Text) == "" {
		util.WriteError(w, http.StatusBadRequest, "missing url")
		return
	}
	u, err := url.ParseRequestURI(body.Text)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
		util.WriteError(w, http.StatusBadRequest, "invalid url")
		return
	}
	if len(body.Text) > 2048 {
		util.WriteError(w, http.StatusBadRequest, "url too long")
		return
	}

	code, err := generateCode(6)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, "code gen fail")
		return
	}

	h.Store.Save(code, body)
	util.WriteJSON(w, http.StatusCreated, map[string]string{"code": code})
}

func (h Handlers) Redirect(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	link, ok := h.Store.Find(code)
	if !ok {
		util.WriteError(w, http.StatusNotFound, "not found")
		return
	}
	http.Redirect(w, r, link.Text, http.StatusFound)
}

// simple crypto-safe base62 gen
func generateCode(n int) (string, error) {
	const alphabet = "abcdefghijklmnopqrstuvwxyz1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		if err != nil {
			return "", err
		}
		b[i] = alphabet[idx.Int64()]
	}
	return string(b), nil
}
