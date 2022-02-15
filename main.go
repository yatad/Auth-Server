package main

import (
	"auth-service/db"
	"auth-service/logger"
	"auth-service/server"
	"auth-service/utils"
	"database/sql"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Env struct {
	Db *sql.DB
}

func main() {

	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}
	//initialize DB
	logger.NewLogger(config.LogPath)
	dbInstance, err := db.Connect(config)
	if err != nil {
		log.Fatalf("Error %s when trying to initialize the DB", err)
	}
	env := &Env{Db: dbInstance}

	//start the server
	mux := http.NewServeMux()
	mux.HandleFunc("/auth", env.CheckAuth)
	auth := server.VerifyHeader(mux)
	log.Fatal(http.ListenAndServe("localhost:8082", auth))
}

func (h *Env) CheckAuth(w http.ResponseWriter, r *http.Request) {
	var p server.Payload
	p.BasicAuth = r.Header["Authorization"]
	username, passwd := server.ParseBasic(p.BasicAuth[0])
	isVerified, err := db.VerifyPasswd(h.Db, username, passwd)
	if err != nil {
		log.Errorf("Failed to verify password %s", err)
		w.WriteHeader(401)
		return
	}
	if isVerified {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(401)
	}

}
