//go:build go1.18
// +build go1.18

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/microsoft/go-mssqldb/azuread"

	log "github.com/sirupsen/logrus"
)

var server string
var user string
var password string
var port int = 1433
var database string
var tenant string

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.000",
	})
	//log.SetFormatter(&log.JSONFormatter{})
	//log.SetReportCaller(true)

	server = os.Getenv("AZR_DB_SERVER")
	user = os.Getenv("AZR_SP_USER")
	password = os.Getenv("AZR_SP_PASSWORD")
	database = os.Getenv("AZR_DB_DATABASE")
	tenant = os.Getenv("AZR_SP_TENANT")
}

func Healthz(w http.ResponseWriter, r *http.Request) {
	log.Info("API Health is OK")
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"alive": true}`)
}

func TestConnection(w http.ResponseWriter, r *http.Request) {
	log.Info("#####################################################")
	log.Info("# Starting Test Connection ")
	log.Info("#####################################################")
	connString := fmt.Sprintf("server=%s;user id=%s@%s;password=%s;port=%d;database=%s;fedauth=ActiveDirectoryServicePrincipal;", server, user, tenant, password, port, database)

	conn, err := sql.Open(azuread.DriverName, connString)
	if err != nil {
		log.Fatal("Open connection failed:", err.Error())
	}
	defer conn.Close()

	stmt, err := conn.Prepare("select @@VERSION")
	if err != nil {
		log.Fatal("Prepare failed:", err.Error())
	}
	defer stmt.Close()

	row := stmt.QueryRow()
	var result string
	err = row.Scan(&result)
	if err != nil {
		log.Fatal("Scan failed:", err.Error())
	}
	log.Info("Result: " + result)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
	log.Info("#####################################################")
}

func main() {
	log.Info("Starting GO DEMO MSSQL server")
	router := mux.NewRouter()
	router.HandleFunc("/healthz", Healthz).Methods("GET")
	router.HandleFunc("/connection", TestConnection).Methods("GET")
	http.ListenAndServe(":8000", router)
}
