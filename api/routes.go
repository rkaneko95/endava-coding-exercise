package api

import "net/http"

func (s *Service) routes() {
	http.HandleFunc("/health-check", s.HealthCheckHandler())
	http.HandleFunc("/token", s.TokenHandler())
	http.HandleFunc("/verify-token", s.VerifyTokenHandler())
}
