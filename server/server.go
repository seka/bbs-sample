package server

import (
	"context"
	"net"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/inconshreveable/log15"

	"github.com/seka/bbs-sample/database"
	"github.com/seka/bbs-sample/server/handler"
)

// Options ...
type Options struct {
	Addr        string
	CookieStore sessions.Store
	DB          database.Database
	CSRF        func(http.Handler) http.Handler
}

// Server ...
type Server struct {
	addr        string
	cookieStore sessions.Store
	csrf        func(http.Handler) http.Handler
	db          database.Database
	server      http.Server
	logger      log15.Logger
	started     chan struct{}
	stopped     chan struct{}
}

// New ...
func New(opt Options) *Server {
	return &Server{
		addr:        opt.Addr,
		cookieStore: opt.CookieStore,
		csrf:        opt.CSRF,
		db:          opt.DB,
		server: http.Server{
			Addr: opt.Addr,
		},
		logger:  log15.New("module", "server"),
		started: make(chan struct{}),
		stopped: make(chan struct{}),
	}
}

// Run ...
func (s *Server) Run(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		l, err := net.Listen("tcp", s.addr)
		if err != nil {
			errCh <- err
		}
		close(s.started)
		s.logger.Info("Listening for client connections on", "addr", s.addr)
		s.setupHandler()
		errCh <- s.server.Serve(l)
	}()
	select {
	case err := <-errCh:
		s.logger.Info("server err", "err", err)
		s.stop(ctx)
		return err
	case <-ctx.Done():
		s.stop(ctx)
		return ctx.Err()
	}
}

// HasStarted ...
func (s *Server) HasStarted() <-chan struct{} {
	return s.started
}

// HasStopped ...
func (s *Server) HasStopped() <-chan struct{} {
	return s.stopped
}

func (s *Server) stop(ctx context.Context) {
	s.server.Shutdown(ctx)
	close(s.stopped)
}

func (s *Server) setupHandler() {
	static := http.FileServer(http.Dir("server/static/"))
	http.Handle("/stylesheets/", static)
	http.Handle("/javascripts/", static)

	opt := handler.Option{
		CookieStore: s.cookieStore,
		DB:          s.db,
	}
	http.Handle("/", handler.NewSession(opt))
	http.Handle("/user", handler.NewUser(opt))
	http.Handle("/bbs", handler.NewBBS(opt))
}
