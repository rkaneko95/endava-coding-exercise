package api

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"rkaneko/endava-coding-exercise/config"
)

type Service struct {
	Config config.ServerConfig
	Log    *logrus.Logger
}

func (s *Service) Init() {
	s.routes()
	s.Log.Error(http.ListenAndServe(fmt.Sprintf("%s:%d", s.Config.Host, s.Config.Port), nil))
}

func (s *Service) HealthCheckHandle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Hello, World!"))
		w.WriteHeader(http.StatusOK)
		s.Log.Debugf("Servier is working")
	}
}
