package host

import (
	"fmt"
	"log"
	"plugin"

	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/chabad360/covey/host/types"
	"github.com/gorilla/mux"
)

var (
	hosts []types.Host
)

func pluginNewHost(host *types.NewHostInfo) (*types.Host, error) {
	p, err := plugin.Open("./plugins/host/" + host.Plugin + ".so")
	if err != nil {
		return nil, err
	}

	n, err := p.Lookup("Plugin")
	if err != nil {
		return nil, err
	}

	var s types.HostPlugin
	s, ok := n.(types.HostPlugin)
	if !ok {
		return nil, fmt.Errorf(host.Plugin, " does not provide a HostPlugin")
	}

	h, err := s.NewHost(host)
	if err != nil {
		return nil, err
	}

	return &h, nil
}

// NewHost adds a new host using the specified plugin.
func NewHost(w http.ResponseWriter, r *http.Request) {
	var host types.NewHostInfo
	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &host)

	h, err := pluginNewHost(&host)
	if err != nil {
		log.Fatal(err)
	}

	hosts = append(hosts, *h)
	// fmt.Fprintf(w, host.server)
	j, err := json.MarshalIndent(h, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, string(j))
}

// AddHandlers adds the mux handlers for the host module.
func AddHandlers(r *mux.Router) {
	r.HandleFunc("/add", NewHost).Methods("POST")
}
