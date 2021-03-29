package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"gobridge/example"
)

func New(api example.Example, a AuthConfig) *Server {
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
	API example.Example
}

type Endpoint int

var (
	HasPermissionEndpoint Endpoint = 0
	AllEndpoints Endpoint = 1
)

func (ep Endpoint) Path() string {
	switch ep {
	case AllEndpoints:
		return "**"
	case HasPermissionEndpoint:
		return "/example/haspermission"
	default:
		return ""
	}
}

func (s *Server) registerHandlers() {
	http.HandleFunc("/example/haspermission", s.Wrap(HasPermissionEndpoint, HandleHasPermission(s.API)))
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
	R []example.Role
}

type HasPermissionResponse struct {
	Bool bool
}

func HandleHasPermission(api example.Example) func(http.ResponseWriter, *http.Request) {
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

		uqid, err := api.HasPermission(r.Context(), req.R)
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

