package http

import (
	"awesomeProject/internal/message_broker/broker_models"
	"awesomeProject/internal/store"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	lru "github.com/hashicorp/golang-lru"
)

type Server struct {
	ctx        context.Context
	idleConnCh chan struct{}
	store      store.Store
	broker     broker_models.Broker
	cache      *lru.TwoQueueCache
	Addres     string
}

func NewServer(ctx context.Context, opts ...ServerOption) *Server {
	srv := &Server{
		ctx:        ctx,
		idleConnCh: make(chan struct{}),
	}

	for _, opt := range opts {
		opt(srv)
	}

	return srv
}

func (s *Server) basicHandler() chi.Router {
	r := chi.NewRouter()
	Baceresourse := NewBaseResource(s.store, s.cache, s.broker)
	r.Mount("/categories", Baceresourse.CategoryRoutes())
	r.Mount("/titles", Baceresourse.TitleRoutes())
	r.Mount("/users", Baceresourse.UserRoutes())
	r.Mount("/publisher", Baceresourse.PublisherRoutes())
	return r
}
func (s *Server) Run() error {
	srv := &http.Server{
		Addr:         s.Addres,
		Handler:      s.basicHandler(),
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 30,
	}

	log.Println("[HTTP] Server runing on", s.Addres)

	return srv.ListenAndServe()
}

func (s *Server) WaitForGraceFulTarmination() {
	<-s.idleConnCh
}

func (s *Server) ListenCtxForGT(srv *http.Server) {
	<-s.ctx.Done()
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Printf("[HTTP] Got err while shutting down %v", err)
		return
	}
	close(s.idleConnCh)
}
