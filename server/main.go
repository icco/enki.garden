package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/icco/enki.garden/client"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

type StatusMsg struct {
	Msg string `json:"msg"`
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	const tpl = `
<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8">
    <title>enki.garden</title>
  </head>
  <body>
    <p>Welcome!</p>
  </body>
</html>
`

	tmpl, err := template.New("name").Parse(tpl)
	if err != nil {
		log.Printf("Template Parse: %+v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Printf("Template Exec: %+v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	// Only allow POSTs
	if r.Method != http.MethodPost {
		http.Error(w, "Not a valid method.", http.StatusBadRequest)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Request Body Read Error: %+v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, err := validateJson(data)
	if err != nil {
		http.Error(w, "Not valid data.", http.StatusBadRequest)
		return
	}

	k := datastore.NewKey(ctx, "File", "stringID", 0, nil)
	e := new(Entity)
	if err := datastore.Get(ctx, k, e); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	old := e.Value
	e.Value = r.URL.Path

	if _, err := datastore.Put(ctx, k, e); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "old=%q\nnew=%q\n", old, e.Value)

	status := StatusMsg{Msg: "Success"}

	json_str, err := json.Marshal(status)
	if err != nil {
		log.Printf("JSON Marshal Error: %+v", err)
		http.Error(w, "Could not turn struct to string.", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s", json_str) // send data to client side
}

func validateJson(msg []byte) (client.EnkiFile, error) {
	var m EventMsg
	err := json.Unmarshal(msg, &m)
	if err != nil {
		log.Printf("Unmarshaling Error: \"%+v\" for message: %s", err, msg)
	} else {
		log.Printf("Parsed Data: %+v", m)
	}

	return m, err
}

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/api", apiHandler)
	http.HandleFunc("/health", healthHandler)

	err := http.ListenAndServe(":9090", nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
