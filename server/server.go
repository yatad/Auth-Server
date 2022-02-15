package server

import (
	"encoding/base64"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

type Server struct {
	server *http.Server
	cfg    *Config
}

type Config struct {
	Host              string
	Port              string
	TemplatesBasePath string
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	DialTimeout       time.Duration
}

type Payload struct {
	BasicAuth []string
	Token     string
	OpenID    string
}

type Verifier struct {
	handler http.Handler
}

func (v *Verifier) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Header["Authorization"]
	if ok {
		v.handler.ServeHTTP(w, r)
	} else {
		w.WriteHeader(401)
	}

}

func VerifyHeader(handle http.Handler) *Verifier {
	return &Verifier{handle}
}

func ParseBasic(s string) (string, string) {
	str := strings.Split(s, " ")
	out, err := base64.URLEncoding.DecodeString(str[1])
	if err != nil {
		logrus.Errorf("%s failed to decode auth string", err)
	}
	f := strings.Split(string(out), ":")
	username := f[0]
	passwd := f[1]
	return username, passwd
}
