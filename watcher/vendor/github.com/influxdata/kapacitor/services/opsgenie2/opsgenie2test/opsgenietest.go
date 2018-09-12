package opsgenie2test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
)

type Server struct {
	mu       sync.Mutex
	ts       *httptest.Server
	URL      string
	requests []Request
	closed   bool
}

func NewServer() *Server {
	s := new(Server)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		or := Request{
			URL:           r.URL.String(),
			Authorization: r.Header.Get("Authorization"),
		}
		dec := json.NewDecoder(r.Body)
		dec.Decode(&or.PostData)
		s.mu.Lock()
		s.requests = append(s.requests, or)
		s.mu.Unlock()
	}))
	s.ts = ts
	s.URL = ts.URL
	return s
}
func (s *Server) Requests() []Request {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.requests
}
func (s *Server) Close() {
	if s.closed {
		return
	}
	s.closed = true
	s.ts.Close()
}

type Request struct {
	URL           string
	PostData      PostData
	Authorization string
}

type PostData struct {
	Message     string              `json:"message"`
	Entity      string              `json:"entity"`
	Alias       string              `json:"alias"`
	Note        string              `json:"note"`
	Priority    string              `json:"priority"`
	Description string              `json:"description"`
	Details     map[string]string   `json:"details"`
	Responders  []map[string]string `json:"responders"`
}
