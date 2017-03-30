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

// User ...
type User struct {
	cookieStore sessions.Store
	userModel   *model.UserModel
	logger      log15.Logger
}

// NewUser ...
func NewUser(opt Option) *User {
	return &User{
		cookieStore: opt.CookieStore,
		userModel:   model.NewUserModel(opt.DB),
		logger:      log15.New("module", "handler", "handler", "user"),
	}
}

func (u *User) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		u.show(w, r)
	case "POST":
		u.save(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (u *User) show(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(filepath.Join("server", "view", "user.html"))
	if err != nil {
		u.logger.Error("Parse template error", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, nil); err != nil {
		u.logger.Error("Template execute error", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (u *User) save(w http.ResponseWriter, r *http.Request) {
	passwd := r.FormValue("password")
	if passwd != r.FormValue("confirm") {
		http.Error(w, "Difference password", http.StatusInternalServerError)
		return
	}
	modelUser := &model.User{
		Name:  r.FormValue("name"),
		Email: r.FormValue("email"),
	}
	if exists := u.userModel.Exists(modelUser); exists {
		http.Error(w, "User already exists", http.StatusInternalServerError)
		return
	}
	modelUser.Password = cryptoutil.GenerateHash(passwd)
	if err := u.userModel.Save(modelUser); err != nil {
		u.logger.Error("Save user error", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}
