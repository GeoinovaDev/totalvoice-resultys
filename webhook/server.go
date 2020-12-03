package webhook

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"

	"git.resultys.com.br/lib/lower/convert/decode"
	"git.resultys.com.br/lib/lower/exception"
	"git.resultys.com.br/lib/lower/promise"
	"git.resultys.com.br/sdk/totalvoice-golang/payload"
)

// Server struct
type Server struct {
	Port string

	hooks map[int]*promise.Promise
	mutex *sync.Mutex
}

// New ...
func New(port string) *Server {
	s := &Server{
		Port:  port,
		mutex: &sync.Mutex{},
		hooks: make(map[int]*promise.Promise),
	}

	s.Start()

	return s
}

// AddHook ...
func (s *Server) AddHook(messageID int) *promise.Promise {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	p := promise.New()

	if s.ExistHook(messageID) {
		p = s.hooks[messageID]
		s.RemoveHook(messageID)
		return p
	}

	s.hooks[messageID] = p

	return p
}

// RemoveHook ...
func (s *Server) RemoveHook(messageID int) {
	delete(s.hooks, messageID)
}

// ResolveHook ...
func (s *Server) ResolveHook(messageID int, response interface{}) {
	s.hooks[messageID].Resolve(response)
	s.RemoveHook(messageID)
}

// ExistHook ...
func (s *Server) ExistHook(messageID int) bool {
	if _, ok := s.hooks[messageID]; ok {
		return true
	}

	return false
}

// Start ...
func (s *Server) Start() {
	go http.ListenAndServe(s.Port, s)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	body := buf.String()

	// saveInFile("/tmp/response.webhook.json", body)

	go s.process(body)

	w.Write([]byte(`{"status": "ok", "code": 200}`))
}

func (s *Server) process(body string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	defer func() {
		err := recover()
		msg := ""

		switch err.(type) {
		case string:
			msg = err.(string)
		case []string:
			msg = strings.Join(err.([]string), ". ")
		case error:
			msg = fmt.Sprint(err)
		default:
			msg = "erro de runtime"
		}

		if err != nil {
			exception.Raise(msg, exception.WARNING)
			fmt.Println(err)
		}
	}()

	response := payload.CallResponse{}
	decode.JSON(body, &response)

	if s.ExistHook(response.ID) {
		s.ResolveHook(response.ID, response)
	} else {
		p := promise.New()
		s.hooks[response.ID] = p
		s.hooks[response.ID].Resolve(response)
	}
}

func saveInFile(filename string, content string) {
	f, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = f.WriteString(content)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
}
