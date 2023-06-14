package api

import "net/http"

func (s *Service) routes() {
	http.HandleFunc("/health-check", s.HealthCheckHandler())
	http.HandleFunc("/generate-keys", s.GenerateKeysHandler())

	http.HandleFunc("/token", s.TokenHandler())
	http.HandleFunc("/verify-token", s.VerifyTokenHandler())
	http.HandleFunc("/signing-keys", s.ListSigningKeysHandler())
}
