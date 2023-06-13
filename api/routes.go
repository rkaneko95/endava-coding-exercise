package api

import "net/http"

func (s *Service) routes() {
	http.HandleFunc("/health-check", s.HealthCheckHandle())
}
