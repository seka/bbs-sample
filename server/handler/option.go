package handler

import (
	"github.com/gorilla/sessions"
	"github.com/seka/bbs-sample/database"
)

// Option ...
type Option struct {
	CookieStore sessions.Store
	DB          database.Database
}
