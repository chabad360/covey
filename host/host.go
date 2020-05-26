package host

import (
	"log"
	"plugin"

	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	hosts []Host
)

// NewHost adds a new host using the specified plugin.
func NewHost(w http.ResponseWriter, r *http.Request) {
	var host NewHostInfo
	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &host)

	p, err := plugin.Open("./plugins/host/" + host.Plugin + ".so")
	if err != nil {
		log.Fatal(err)
	}

	n, err := p.Lookup("Plugin")
	if err != nil {
		log.Fatal(err)
	}

	var s HostPlugin
	s, ok := n.(HostPlugin)
	if !ok {
		log.Fatal(host.Plugin, "does not provide a HostPlugin")
	}

	h, err := s.NewHost(&host)
	if err != nil {
		log.Fatal(err)
	}

	hosts = append(hosts, h)
	// fmt.Fprintf(w, host.server)
	json.NewEncoder(w).Encode(host)
}

// AddHandlers adds the mux handlers for the host module.
func AddHandlers(r *mux.Router) {
	r.HandleFunc("/add", NewHost).Methods("POST")
}
