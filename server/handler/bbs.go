package handler

import (
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	gsess "github.com/gorilla/sessions"
	"github.com/inconshreveable/log15"
	"github.com/justinas/nosurf"

	"github.com/seka/bbs-sample/model"
)

// BBS ...
type BBS struct {
	cookieStore  gsess.Store
	messageModel *model.MessageModel
	logger       log15.Logger
}

// NewBBS ...
func NewBBS(opt Option) *BBS {
	return &BBS{
		cookieStore:  opt.CookieStore,
		messageModel: model.NewMessageModel(opt.DB),
		logger:       log15.New("module", "handler", "handler", "bbs"),
	}
}

func (b *BBS) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sess, err := b.cookieStore.Get(r, "user")
	if err != nil || sess.IsNew {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	switch r.Method {
	case "GET":
		b.show(sess, w, r)
	case "POST":
		b.post(sess, w, r)
	default:
		http.NotFound(w, r)
	}
}

func (b *BBS) show(sess *gsess.Session, w http.ResponseWriter, r *http.Request) {
	msgs, err := b.messageModel.FindAll()
	if err != nil {
		b.logger.Error("find all messages error", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl, err := template.ParseFiles(filepath.Join("server", "view", "bbs.html"))
	if err != nil {
		b.logger.Error("parse template error", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := &struct {
		Name      string
		Messages  []*model.Message
		CsrfToken string
	}{
		Name:      sess.Values["name"].(string),
		Messages:  msgs,
		CsrfToken: nosurf.Token(r),
	}
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (b *BBS) post(sess *gsess.Session, w http.ResponseWriter, r *http.Request) {
	msg := &model.Message{
		UserID:    sess.Values["id"].(int),
		Message:   r.FormValue("message"),
		CreatedAt: time.Now().String(),
	}
	if err := b.messageModel.Save(msg); err != nil {
		b.logger.Error("save message error", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/bbs", http.StatusFound)
}

var _ http.Handler = (*BBS)(nil)
