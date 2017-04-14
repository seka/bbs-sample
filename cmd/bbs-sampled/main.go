package main

import (
	"context"
	"flag"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

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
	flag.StringVar(&args.Database.Addr, "database-addr", "localhost:3306", "specify the address of a database")
	flag.StringVar(&args.Database.Name, "database-name", "bbs-sample", "specify the name of a database")
	flag.StringVar(&args.Database.User, "database-user", "bbs-sample-user", "specify the username to connect for database")
	flag.StringVar(&args.Database.Password, "database-password", "bbs-sample-password", "specify the password to connect for database")
}

func main() {
	flag.Parse()
	logutil.SeetupRootLogger(args.LogLevel)
	if err := newMain(args).Run(); err != nil {
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
	return &Main{
		appSecret: args.AppSecret,
		db:        db,
		logger:    log15.New("module", "main"),
		server: server.New(server.Options{
			Addr:        net.JoinHostPort("", args.Port),
			CookieStore: sessions.NewCookieStore([]byte(args.AppSecret)),
			DB:          db,
		}),
	}
}

// Run ...
func (m *Main) Run() error {
	signalCtx, cancelFunc := m.createSignalHandler()
	var wg sync.WaitGroup
	dbErrCh := make(chan error, 1)
	go func() {
		wg.Add(1)
		dbErrCh <- m.runDatabase(signalCtx)
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

func (m *Main) runDatabase(ctx context.Context) error {
	if err := m.db.Connect(); err != nil {
		return err
	}
	for {
		if err := m.db.Ping(); err != nil {
			return err
		}
		select {
		case <-time.After(10 * time.Second):
			// noop
		case <-ctx.Done():
			if err := m.db.Disconnect(); err != nil {
				return err
			}
			return ctx.Err()
		}
	}
}
