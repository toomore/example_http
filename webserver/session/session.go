package session

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"strings"
)
import "github.com/toomore/hashvalues"

const sessionName = "session"
const sessionSplitSign = "|"

func makeSession(value string, resp *http.Request) *http.Cookie {
	return &http.Cookie{
		Name:     sessionName,
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
	http.SetCookie(s.w, makeSession(joinCodeMsg(code, msg), s.resp))
}

func (s *Session) parse() {
	if rawcookie, err := s.resp.Cookie(sessionName); err == nil {
		code, msg := splitCodeMsg(rawcookie.String())
		if err := s.Hashvalues.Decode([]byte(code), []byte(msg)); err != nil {
			s.Set("", "")
			s.Save()
		}
	}
}

func splitCodeMsg(rawcookie string) (code, msg string) {
	cookies := strings.Split(rawcookie[len(sessionName)+1:], sessionSplitSign)
	return cookies[0], cookies[1]
}

func joinCodeMsg(code, msg []byte) string {
	return fmt.Sprintf("%s%s%s", code, sessionSplitSign, msg)
}
