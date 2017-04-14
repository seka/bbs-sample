package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/inconshreveable/log15"

	"github.com/seka/bbs-sample/database"
)

var (
	args   Arguments
	logger = log15.New("module", "main")
)

// Arguments ...
type Arguments struct {
	RetryCount    int
	RetryInterval int
	Database      database.Options
}

func init() {
	flag.IntVar(&args.RetryCount, "retry-count", 10, "specify the retry count of connect database")
	flag.IntVar(&args.RetryInterval, "retry-interval", 1, "specify the retry interval second of connect database")
	flag.StringVar(&args.Database.Addr, "database-addr", "localhost:3306", "specify the address of a database")
	flag.StringVar(&args.Database.Name, "database-name", "bbs-sample", "specify the name of a database")
	flag.StringVar(&args.Database.User, "database-user", "bbs-sample-user", "specify the username to connect for database")
	flag.StringVar(&args.Database.Password, "database-password", "bbs-sample-password", "specify the password to connect for database")
}

func main() {
	flag.Parse()
	db := database.NewMySQL(args.Database)
	sigs := createSignalHandler()
	for {
		if args.RetryCount == 0 {
			log15.Error("Used up retry count of connect database")
			return
		}
		err := db.Connect()
		if err == nil {
			log15.Info("Establish database connect")
			return
		}
		log15.Info("Waiting for database ...", "err", err)
		select {
		case <-time.After(time.Duration(args.RetryInterval) * time.Second):
			args.RetryCount--
		case <-sigs:
			log15.Info("Cancel database connect")
			return
		}
	}
}

func createSignalHandler() chan os.Signal {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	return sigs
}
