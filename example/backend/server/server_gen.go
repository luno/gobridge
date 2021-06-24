package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"gobridge/example/backend"
	"gobridge/example/backend/second"
)

func New(api backend.Example, a AuthConfig) *Server {
	s := &Server{
		Auth: a,
		API: api,
	}

	s.registerHandlers()

	return s
}

type AuthConfig map[Endpoint]func(token string) (bool, error)

type Server struct {
	Auth AuthConfig
	API backend.Example
}

type Endpoint int

var (
	HasPermissionEndpoint Endpoint = 0
	WhatsTheTimeEndpoint Endpoint = 1
	AllEndpoints Endpoint = 2
)

func (ep Endpoint) Path() string {
	switch ep {
	case AllEndpoints:
		return "**"
	case HasPermissionEndpoint:
		return "/backend/haspermission"
	case WhatsTheTimeEndpoint:
		return "/backend/whatsthetime"
	default:
		return ""
	}
}

func (s *Server) registerHandlers() {
	http.HandleFunc("/backend/haspermission", s.Wrap(HasPermissionEndpoint, HandleHasPermission(s.API)))
	http.HandleFunc("/backend/whatsthetime", s.Wrap(WhatsTheTimeEndpoint, HandleWhatsTheTime(s.API)))
}

func (s *Server) Wrap(e Endpoint, fn func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Kind, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		// Check to see if the 'AllEndpoints' type was set
		authFunc, ok := s.Auth[AllEndpoints]
		if ok {
			allow, msg, reason := checkAuth(w, r, authFunc)
			if !allow {
				http.Error(w, msg, reason)
				return
			}
		} else {
			// Check to see if there is auth setup for this endpoint as there 
			// is no config for all the routes.
			authFunc, ok = s.Auth[e]
			if ok {
				allow, msg, reason := checkAuth(w, r, authFunc)
				if !allow {
					http.Error(w, msg, reason)
					return
				}
			}
		}

		fn(w, r)
	}
}

func checkAuth(w http.ResponseWriter, r *http.Request, authFunc func(token string) (bool, error)) (bool, string, int) {
	ah := r.Header.Get("Authorization")
	if ah == "" {
		return false, "unauthorised", http.StatusUnauthorized
	}

	t := strings.TrimSpace(ah)
	if t == "" {
		return false, "no authorization token present", http.StatusUnauthorized
	}

	allow, err := authFunc(t)
	if err != nil {
		http.Error(w, "unauthorised", http.StatusUnauthorized)
		return false, "no authorization token present", http.StatusUnauthorized
	}

	return allow, "", 0
}

type HasPermissionRequest struct {
	R []backend.Role
	U backend.User
}

type HasPermissionResponse struct {
	Bool bool
}

func HandleHasPermission(api backend.Example) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		var req HasPermissionRequest
		err = json.Unmarshal(b, &req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		uqid, err := api.HasPermission(r.Context(), req.R, req.U)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		var resp HasPermissionResponse
		resp.Bool, _ = uqid, err
	
		respBody, err := json.Marshal(resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(respBody)
		if err != nil {
			_, _ = w.Write([]byte(err.Error()))
			return
		}
	}
}

type WhatsTheTimeRequest struct {
	Date time.Time
	Toy second.Toy
}

type WhatsTheTimeResponse struct {
	Bool bool
}

func HandleWhatsTheTime(api backend.Example) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		var req WhatsTheTimeRequest
		err = json.Unmarshal(b, &req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		epfq, err := api.WhatsTheTime(r.Context(), req.Date, req.Toy)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		var resp WhatsTheTimeResponse
		resp.Bool, _ = epfq, err
	
		respBody, err := json.Marshal(resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(respBody)
		if err != nil {
			_, _ = w.Write([]byte(err.Error()))
			return
		}
	}
}

