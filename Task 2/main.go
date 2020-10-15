package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

var currid int = 1

// struct
type Article struct {
	Title     string `json:"title"`
	SubTitle  string `json:"subtitle"`
	ID        string `json:"id"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}

type articleHandlers struct {
	sync.Mutex
	store map[string]Article
}

func (h *articleHandlers) articles(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.get(w, r)
		return
	case "POST":
		h.post(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
		return
	}
}

func (h *articleHandlers) get(w http.ResponseWriter, r *http.Request) {
	articles := make([]Article, len(h.store))

	h.Lock()
	i := 0
	for _, article := range h.store {
		articles[i] = article
		i++
	}
	h.Unlock()

	jsonBytes, err := json.Marshal(articles)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *articleHandlers) getArticle(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.String(), "/")
	if len(parts) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	h.Lock()
	article, ok := h.store[parts[2]]
	h.Unlock()
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	jsonBytes, err := json.Marshal(article)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *articleHandlers) post(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	var article Article
	err = json.Unmarshal(bodyBytes, &article)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	article.ID = strconv.Itoa(currid)
	currid++
	article.Timestamp = fmt.Sprintf("%d", time.Now().UnixNano())
	h.Lock()
	h.store[article.ID] = article
	defer h.Unlock()
}

func newArticleHandlers() *articleHandlers {
	return &articleHandlers{
		store: map[string]Article{
			"1": Article{
				Title:     "Daddy",
				SubTitle:  "Big Daddy",
				ID:        "1",
				Timestamp: "Today",
				Content:   "sfjkhasdjkhfgdsakjfvkjsdbkvh",
			},
			"2": Article{
				Title:     "Daddy",
				SubTitle:  "Big Daddy",
				ID:        "2",
				Timestamp: "Today",
				Content:   "sfjkhasdjkhfgdsakjfvkjsdbkvh",
			},
			"3": Article{
				Title:     "Daddy",
				SubTitle:  "Big Daddy",
				ID:        "3",
				Timestamp: "Today",
				Content:   "sfjkhasdjkhfgdsakjfvkjsdbkvh",
			},
			"4": Article{
				Title:     "Daddy",
				SubTitle:  "Big Daddy",
				ID:        "4",
				Timestamp: "Today",
				Content:   "sfjkhasdjkhfgdsakjfvkjsdbkvh",
			},
			"5": Article{
				Title:     "234234",
				SubTitle:  "234addy",
				ID:        "5",
				Timestamp: "Today23423",
				Content:   "sdfas",
			},
		},
	}
}

func main() {
	articleHandlers := newArticleHandlers()
	http.HandleFunc("/articles", articleHandlers.articles)
	http.HandleFunc("/articles/", articleHandlers.getArticle)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
	println("Works")
}
