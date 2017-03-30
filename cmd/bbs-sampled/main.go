package main

import (
	"context"
	"flag"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gorilla/sessions"
	"github.com/inconshreveable/log15"

	"github.com/seka/bbs-sample/database"
	"github.com/seka/bbs-sample/internal/logutil"
	"github.com/seka/bbs-sample/server"
)

var (
	args Arguments
)

func init() {
	flag.StringVar(&args.LogLevel, "log-level", "info", "spcify the application log-level")
	flag.StringVar(&args.Port, "port", "8080", "specify the application listening port")
	flag.StringVar(&args.AppSecret, "app-secret", "", "specify the authentication key provided should be 32 bytes long")
	flag.StringVar(&args.Database.Host, "database-host", "localhost", "specify the host of a database")
	flag.StringVar(&args.Database.Port, "database-port", "3306", "specify the port of a database")
	flag.StringVar(&args.Database.Name, "database-name", "bbs-sample", "specify the name of a database")
	flag.StringVar(&args.Database.User, "database-user", "bbs-sample-user", "specify the username to connect for database")
	flag.StringVar(&args.Database.Password, "database-password", "bbs-sample-password", "specify the password to connect for database")
}

func main() {
	flag.Parse()
	logutil.SeetupRootLogger(args.LogLevel)
	m := newMain(args)
	if err := m.Run(); err != nil {
		os.Exit(1)
	}
}

// Arguments ...
type Arguments struct {
	Port      string
	LogLevel  string
	AppSecret string
	Database  database.Options
}

// Main ...
type Main struct {
	appSecret string
	db        database.Database
	logger    log15.Logger
	server    *server.Server
}

func newMain(args Arguments) *Main {
	db := database.NewMySQL(args.Database)
	server := server.New(server.Options{
		Addr:        net.JoinHostPort("", args.Port),
		CookieStore: sessions.NewCookieStore([]byte(args.AppSecret)),
		DB:          db,
	})
	return &Main{
		appSecret: args.AppSecret,
		db:        db,
		logger:    log15.New("module", "main"),
		server:    server,
	}
}

// Run ...
func (m *Main) Run() error {
	signalCtx, cancelFunc := m.createSignalHandler()
	var wg sync.WaitGroup
	dbErrCh := make(chan error, 1)
	go func() {
		wg.Add(1)
		dbErrCh <- m.db.Connect(signalCtx)
		wg.Done()
	}()
	serverErrCh := make(chan error, 1)
	go func() {
		wg.Add(1)
		serverErrCh <- m.server.Run(signalCtx)
		wg.Done()
	}()
	select {
	case err := <-serverErrCh:
		m.logger.Error("Server error", "err", err)
		cancelFunc()
		wg.Wait()
		return err
	case err := <-dbErrCh:
		m.logger.Error("Database error", "err", err)
		cancelFunc()
		wg.Wait()
		return err
	case <-signalCtx.Done():
		wg.Wait()
		return signalCtx.Err()
	}
}

func (m *Main) createSignalHandler() (context.Context, context.CancelFunc) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		s := <-sigs
		m.logger.Info("Got signal", "sigs", s)
		cancel()
	}()
	return ctx, cancel
}
