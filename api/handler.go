package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Response struct {
	Token   string `json:"token,omitempty"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

type ResponseValidation struct {
	IssuedAt  *time.Time `json:"issuedAt"`
	ExpiredAt *time.Time `json:"expiredAt"`
	Active    bool       `json:"active"`
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
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Hello, World!"))
		s.Log.Debugf("servier is working")
	}
}

func (s *Service) GenerateKeysHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		id, err := s.GenerateKeys()
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			writeErrorResponse(w, "was not able to generate data", err.Error())
			s.Log.Errorf("error in GenerateKeys: %s", err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(&struct{ UUID, Message string }{
			UUID:    id,
			Message: "Keys generated",
		})
		s.Log.Debugf("Keys generated")
	}
}

func (s *Service) TokenHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		msg := "failed to generate the token"
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			errMsg := "method not allowed"
			writeErrorResponse(w, msg, errMsg)
			s.Log.Errorf(errMsg)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			errMsg := "missing Authorization header"
			writeErrorResponse(w, msg, errMsg)
			s.Log.Errorf(errMsg)
			return
		}

		token, err := s.CreateToken(authHeader)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			writeErrorResponse(w, msg, err.Error())
			s.Log.Errorf("error in CreateToken: %s", err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(&Response{
			Token:   token,
			Message: "token generated",
		})
		s.Log.Infof("token generated")
	}
}

func (s *Service) VerifyTokenHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		msg := "failed to verify the token"
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			errMsg := "Method not allowed"
			writeErrorResponse(w, msg, errMsg)
			s.Log.Errorf(errMsg)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			errMsg := "missing Authorization header"
			writeErrorResponse(w, msg, errMsg)
			s.Log.Errorf(errMsg)
			return
		}

		issue, expired, err := s.VerifyToken(authHeader)
		if err != nil {
			if err.Error() == expiredErr {
				w.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(w).Encode(&ResponseValidation{
					IssuedAt:  issue,
					ExpiredAt: expired,
					Active:    false,
				})
				s.Log.Warnf("token has expired")
				return
			}

			w.WriteHeader(http.StatusUnauthorized)
			writeErrorResponse(w, msg, err.Error())
			s.Log.Errorf("error in VerifyToken: %s", err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(&ResponseValidation{
			IssuedAt:  issue,
			ExpiredAt: expired,
			Active:    true,
		})
		s.Log.Infof("valid token")
	}
}

func (s *Service) ListSigningKeysHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		msg := "failed to get signing keys"
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			errMsg := "Method not allowed"
			writeErrorResponse(w, msg, errMsg)
			s.Log.Errorf(errMsg)
			return
		}

		keys, err := s.ListSigningKeys()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			writeErrorResponse(w, msg, err.Error())
			s.Log.Errorf("error in ListSigningKeys: %s", err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(keys)
		s.Log.Infof("there are %d signing keys", len(keys))
	}
}

func writeErrorResponse(w http.ResponseWriter, msg, errMsg string) {
	_ = json.NewEncoder(w).Encode(&Response{
		Message: msg,
		Error:   errMsg,
	})
}
