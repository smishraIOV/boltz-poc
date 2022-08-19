package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"context"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/patogallaiov/boltz-poc/connectors"
	"github.com/patogallaiov/boltz-poc/storage"
	log "github.com/sirupsen/logrus"
)

const (
	svcStatusOk          = "ok"
	svcStatusDegraded    = "degraded"
	svcStatusUnreachable = "unreachable"
)

type Server struct {
	srv   http.Server
	boltz connectors.BoltzConnector
	rsk   connectors.RSKConnector
	db    storage.DBConnector
	now   func() time.Time
}

func New(boltz connectors.BoltzConnector, rsk connectors.RSKConnector, db storage.DBConnector) Server {
	return newServer(boltz, rsk, db, time.Now)
}

func newServer(boltz connectors.BoltzConnector, rsk connectors.RSKConnector, db storage.DBConnector, now func() time.Time) Server {
	return Server{
		boltz: boltz,
		rsk:   rsk,
		db:    db,
		now:   now,
	}
}

func (s *Server) Start(port uint) error {
	r := mux.NewRouter()
	r.Path("/health").Methods(http.MethodGet).HandlerFunc(s.checkHealthHandler)
	r.Path("/swapinfo").Methods(http.MethodGet).HandlerFunc(s.getSwapInfo)
	w := log.StandardLogger().WriterLevel(log.DebugLevel)
	h := handlers.LoggingHandler(w, r)
	defer func(w *io.PipeWriter) {
		_ = w.Close()
	}(w)

	s.srv = http.Server{
		Addr:    ":" + fmt.Sprint(port),
		Handler: h,
	}
	log.Info("server started at localhost:", s.srv.Addr)

	err := s.srv.ListenAndServe()
	if err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) Shutdown() {
	log.Info("stopping server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.srv.Shutdown(ctx); err != nil {
		log.Fatal("server shutdown failed: ", err)
	}
	log.Info("server stopped")
}

func (s *Server) checkHealthHandler(w http.ResponseWriter, _ *http.Request) {
	type services struct {
		Db  string `json:"db"`
		Rsk string `json:"rsk"`
	}
	type healthRes struct {
		Status   string   `json:"status"`
		Services services `json:"services"`
	}

	dbSvcStatus := svcStatusOk
	rskSvcStatus := svcStatusOk

	if err := s.db.CheckConnection(); err != nil {
		log.Error("error checking db connection status: ", err.Error())
		dbSvcStatus = svcStatusUnreachable
	}

	if err := s.rsk.CheckConnection(); err != nil {
		log.Error("error checking rsk connection status: ", err.Error())
		rskSvcStatus = svcStatusUnreachable
	}

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	response := healthRes{
		Services: services{
			Db:  dbSvcStatus,
			Rsk: rskSvcStatus,
		},
	}
	err := enc.Encode(response)
	if err != nil {
		log.Error("error encoding response: ", err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func (s *Server) getSwapInfo(w http.ResponseWriter, _ *http.Request) {
	info, err := s.boltz.GetReverseSwapInfo()
	if err != nil {
		log.Error("error boltz response: ", err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	err = enc.Encode(info)
	if err != nil {
		log.Error("error encoding response: ", err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
