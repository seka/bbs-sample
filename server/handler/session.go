package handler

import (
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/gorilla/sessions"
	"github.com/inconshreveable/log15"
	"github.com/seka/bbs-sample/internal/cryptoutil"
	"github.com/seka/bbs-sample/model"
)

// Session ...
type Session struct {
	cookieStore sessions.Store
	userModel   *model.UserModel
	logger      log15.Logger
}

// NewSession ...
func NewSession(opt Option) *Session {
	return &Session{
		cookieStore: opt.CookieStore,
		userModel:   model.NewUserModel(opt.DB),
		logger:      log15.New("module", "handler", "handler", "session"),
	}
}

// ServeHTTP ...
func (s *Session) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		s.show(w, r)
	case "POST":
		if r.FormValue("_method") == "DELETE" {
			s.doSignout(w, r)
			return
		}
		s.doSingin(w, r)
	case "DELETE":
		s.doSignout(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (s *Session) show(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(filepath.Join("server", "view", "index.html"))
	if err != nil {
		s.logger.Error("Parse template error", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, nil); err != nil {
		s.logger.Error("Template execute error", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Session) doSingin(w http.ResponseWriter, r *http.Request) {
	user, err := s.userModel.Find(&model.User{
		Email:    r.FormValue("email"),
		Password: cryptoutil.GenerateHash(r.FormValue("password")),
	})
	if err != nil {
		s.logger.Error("Find user error", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if user.ID == 0 {
		s.logger.Error("NotFound user", "email", user.Email, "password", user.Password)
		http.NotFound(w, r)
		return
	}
	if err := s.saveCookie(w, r, user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/bbs", http.StatusFound)
}

func (s *Session) doSignout(w http.ResponseWriter, r *http.Request) {
	if err := s.removeCookie(w, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func (s *Session) saveCookie(w http.ResponseWriter, r *http.Request, users *model.User) error {
	sess, err := s.cookieStore.New(r, "user")
	if err != nil {
		s.logger.Error("NewCookieStore error", "err", err)
		return err
	}
	sess.Values["id"] = users.ID
	sess.Values["name"] = users.Name
	if err := sess.Save(r, w); err != nil {
		s.logger.Error("Save cookie store error", "err", err)
		return err
	}
	return nil
}

func (s *Session) removeCookie(w http.ResponseWriter, r *http.Request) error {
	sess, err := s.cookieStore.Get(r, "user")
	if err != nil {
		s.logger.Error("Get cookie error", "err", err)
		return err
	}
	sess.Options.MaxAge = -1
	if err := sess.Save(r, w); err != nil {
		s.logger.Error("Remove cookie store error", "err", err)
		return err
	}
	return nil
}

var _ http.Handler = (*Session)(nil)
