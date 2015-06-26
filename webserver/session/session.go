package session

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"strings"
)
import "github.com/toomore/hashvalues"

func makeSession(value string, resp *http.Request) *http.Cookie {
	return &http.Cookie{
		Name:     "session",
		Value:    value,
		Path:     "/",
		Domain:   strings.Split(resp.Host, ":")[0],
		HttpOnly: true,
		//Expires: time.Now().Add(time.Duration(expires) * time.Second),
	}
}

type Session struct {
	Hashvalues *hashvalues.HashValues
	w          http.ResponseWriter
	resp       *http.Request
}

func New(hashkey []byte, w http.ResponseWriter, resp *http.Request) *Session {
	s := &Session{
		Hashvalues: hashvalues.New(hashkey, sha256.New),
		w:          w,
		resp:       resp,
	}
	s.parse()
	return s
}

func (s *Session) Set(key, value string) {
	s.Hashvalues.Set(key, value)
}

func (s *Session) Save() {
	code, msg := s.Hashvalues.Encode()
	http.SetCookie(s.w, makeSession(fmt.Sprintf("%s|%s", code, msg), s.resp))
}

func (s *Session) parse() {
	if rawcookie, err := s.resp.Cookie("session"); err == nil {
		cookies := strings.Split(rawcookie.String()[8:], "|")
		if err := s.Hashvalues.Decode([]byte(cookies[0]), []byte(cookies[1])); err != nil {
			s.Set("", "")
			s.Save()
		}
	}
}
