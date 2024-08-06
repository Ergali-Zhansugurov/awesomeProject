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
	CategoryResponse := NewCategoryResource(s.store, s.cache, s.broker)
	r.Mount("/", CategoryResponse.Routes())
	TitleResourse := NewTitleResource(s.store, s.cache, s.broker)
	r.Mount("/", TitleResourse.Routes())
	UserResourse := NewUserResourse(s.store, s.cache, s.broker)
	r.Mount("/", UserResourse.Routes())
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

func (srv *Server) ListenCtxForGt(s *http.Server) {
	<-srv.ctx.Done()
	if err := s.Shutdown(context.Background()); err != nil {
		log.Printf("[HTTP] Got err while shutting down %v", err)
		return
	}
}
func (srv *Server) WaitForGraceFulTarmination() {
	<-srv.idleConnCh
}
