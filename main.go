package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/hpcloud/tail"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

type ModelTransferJson struct {
	SenderID int `json:"sender_id,omitempty"`
	UID      int `json:"resep_id,omitempty"`
	Value    int `json:"val,omitempty"`
}

var (
	model ModelTransferJson
	body  []byte
	conn  *sql.DB
	err   error
)

const (
	logPath  = "development.log"
	httpPort = 4000
)

func main() {
	defer conn.Close()
	conn, err = initDb()
	if err != nil {
		log.Fatal(err)
	}

	openLogFile(logPath)

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/payment", payHandler)
	http.HandleFunc("/transfer", transferHandler)

	fmt.Printf("listening on %v\n", httpPort)
	fmt.Printf("Logging to %v\n", logPath)

	checkLog(conn)

	err := http.ListenAndServe(fmt.Sprintf(":%d", httpPort), logRequest(http.DefaultServeMux))
	if err != nil {
		log.Fatal(err)
	}
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ = ioutil.ReadAll(r.Body)
		if len(body) != 0 {
			request := strings.ReplaceAll(string(body), "\n", "")
			log.Println(r.URL.Path, request)
		}
		handler.ServeHTTP(w, r)
	})
}

func openLogFile(logfile string) {
	if logfile != "" {
		lf, err := os.OpenFile(logfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0640)

		if err != nil {
			log.Fatal("OpenLogfile: os.OpenFile:", err)
		}

		log.SetOutput(lf)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Hello World</h1>")
}

func payHandler(w http.ResponseWriter, r *http.Request) {
	if len(body) != 0 {
		if err := json.Unmarshal([]byte(string(body)), &model); err != nil {
			log.Println(err)
		}
	}
	insert(conn, model)
}
func transferHandler(w http.ResponseWriter, r *http.Request) {
	if len(body) != 0 {
		if err := json.Unmarshal([]byte(string(body)), &model); err != nil {
			log.Println(err)
		}
	}
	transfer(conn, model)
}

func checkLog(conn *sql.DB) {
	t, _ := tail.TailFile(logPath, tail.Config{Follow: false})
	var (
		reTime *regexp.Regexp = regexp.MustCompile(`^(.*?) main.go`)
		reVal  *regexp.Regexp = regexp.MustCompile(`/payment (.*)$`)
		reTrns *regexp.Regexp = regexp.MustCompile(`/transfer (.*)$`)
		res    []string
		pay    []string
		trns   []string
	)
	lastTransaction := selectLastTransaction(conn)
	currTime, _ := time.Parse(time.RFC3339Nano, lastTransaction)

	for line := range t.Lines {
		if res = reTime.FindStringSubmatch(line.Text); len(res) > 0 {
			secondTime, _ := time.Parse("2006/01/02 15:04:05", res[1])
			if secondTime.After(currTime) {
				if pay = reVal.FindStringSubmatch(line.Text); len(pay) > 0 {
					if err := json.Unmarshal([]byte(pay[1]), &model); err != nil {
						log.Println(err)
					}
					insert(conn, model)
				} else if trns = reTrns.FindStringSubmatch(line.Text); len(trns) > 0 {
					if err := json.Unmarshal([]byte(trns[1]), &model); err != nil {
						log.Println(err)
					}
					transfer(conn, model)
				}
			}
		}
	}
}
