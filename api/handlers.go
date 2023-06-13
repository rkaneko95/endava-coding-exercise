package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Payload struct {
	Token string `json:"token,omitempty"`
	Error string `json:"errors,omitempty"`
}

func (s *Service) Init() {
	s.routes()
	s.Log.Error(
		http.ListenAndServe(
			fmt.Sprintf("%s:%d", s.Config.Host, s.Config.Port),
			nil),
	)
}

func (s *Service) HealthCheckHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Hello, World!"))
		w.WriteHeader(http.StatusOK)
		s.Log.Debugf("Servier is working")
	}
}

func (s *Service) TokenHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			errMsg := "Method not allowed"
			writeErrorResponse(w, errMsg)
			s.Log.Errorf(errMsg)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			errMsg := "Must contain Authorization"
			writeErrorResponse(w, errMsg)
			s.Log.Errorf(errMsg)
			return
		}

		token, err := s.CreateToken(authHeader)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			writeErrorResponse(w, err.Error())
			s.Log.Errorf("Error in CreateToken: %s", err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(&Payload{
			Token: token,
		})
		s.Log.Infof("Token generated")
	}
}

func writeErrorResponse(w http.ResponseWriter, errMsg string) {
	_ = json.NewEncoder(w).Encode(&Payload{
		Error: errMsg,
	})
}
