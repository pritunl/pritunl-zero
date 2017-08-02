package auth

import (
	"fmt"
	"github.com/pritunl/pritunl-zero/cookie"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/node"
	"github.com/pritunl/pritunl-zero/service"
	"github.com/pritunl/pritunl-zero/session"
	"net/http"
	"net/url"
)

func Get(db *database.Database, state string) (tokn *Token, err error) {
	coll := db.Tokens()
	tokn = &Token{}

	err = coll.FindOneId(state, tokn)
	if err != nil {
		return
	}

	return
}

func CookieSession(db *database.Database,
	w http.ResponseWriter, r *http.Request) (
	cook *cookie.Cookie, sess *session.Session, err error) {

	cook, err = cookie.Get(w, r)
	if err != nil {
		return
	}

	sess, err = cook.GetSession(db)
	if err != nil {
		switch err.(type) {
		case *errortypes.NotFoundError:
			sess = nil
			err = nil
			break
		}
		return
	}

	return
}

func CookieSessionProxy(db *database.Database, srvc *service.Service,
	w http.ResponseWriter, r *http.Request) (
	cook *cookie.Cookie, sess *session.Session, err error) {

	cook, err = cookie.GetProxy(srvc, w, r)
	if err != nil {
		return
	}

	sess, err = cook.GetSession(db)
	if err != nil {
		switch err.(type) {
		case *errortypes.NotFoundError:
			sess = nil
			err = nil
			break
		}
		return
	}

	return
}

func CsrfCheck(w http.ResponseWriter, r *http.Request, domain string) bool {
	scheme := ""
	port := ""
	if node.Self.Protocol == "http" {
		scheme = "http"
		if node.Self.Port != 80 {
			port += fmt.Sprintf(":%d", node.Self.Port)
		}
	} else {
		scheme = "https"
		if node.Self.Port != 443 {
			port += fmt.Sprintf(":%d", node.Self.Port)
		}
	}
	match := fmt.Sprintf("%s://%s%s", scheme, domain, port)

	origin := r.Header.Get("Origin")
	if origin != "" {
		u, err := url.Parse(origin)
		if err != nil {
			http.Error(w, "Server error", 500)
			return false
		}
		origin = fmt.Sprintf("%s://%s", u.Scheme, u.Host)
	}

	referer := r.Header.Get("Referer")
	if referer != "" {
		u, err := url.Parse(referer)
		if err != nil {
			http.Error(w, "Server error", 500)
			return false
		}
		referer = fmt.Sprintf("%s://%s", u.Scheme, u.Host)
	}

	if origin != "" && origin != match {
		http.Error(w, "CSRF origin error", 401)
		return false
	}
	if referer != "" && referer != match {
		http.Error(w, "CSRF referer error", 401)
		return false
	}

	return true
}
