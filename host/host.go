package host

import (
	"fmt"
	"io"
	"log"
	"os"
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

func loadPlugin(pluginName string) (types.HostPlugin, error) {
	p, err := plugin.Open("./plugins/host/" + pluginName + ".so")
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
		return nil, fmt.Errorf(pluginName, " does not provide a HostPlugin")
	}

	return s, nil
}

// NewHost adds a new host using the specified plugin.
func NewHost(w http.ResponseWriter, r *http.Request) {
	var host types.HostInfo
	reqBody, _ := ioutil.ReadAll(r.Body)
	if err := json.Unmarshal(reqBody, &host); err != nil {
		log.Fatal(err)
	}

	for _, h := range hosts {
		if h.Name == host.Name {
			errorWriter(w, fmt.Errorf("Duplicate host: %v", host.Name))
		}
	}

	p, err := loadPlugin(host.Plugin)
	if err != nil {
		log.Fatal(err)
	}

	h, err := p.NewHost(reqBody)
	if err != nil {
		log.Fatal(err)
	}

	hosts = append(hosts, h)
	j, err := json.MarshalIndent(hosts, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.Create("./config/hosts.json")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if err = f.Chmod(0600); err != nil {
		log.Fatal(err)
	}
	if _, err = f.Write(j); err != nil {
		log.Fatal(err)
	}

	j, err = json.MarshalIndent(h, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, string(j))
}

// RegisterHandlers adds the mux handlers for the host module.
func RegisterHandlers(r *mux.Router) {
	r.HandleFunc("/add", NewHost).Methods("POST")
}

// LoadConfig loads up the stored hosts
func LoadConfig() {
	log.Println("Loading Host Config")
	f, err := os.Open("./config/hosts.json")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var h []json.RawMessage
	if err = json.NewDecoder(f).Decode(&h); err != nil {
		log.Fatal(err)
	}

	var plugins = make(map[string]types.HostPlugin)
	p, err := loadPlugin("ssh")
	if err != nil {
		log.Fatal(err)
	}
	plugins["ssh"] = p

	for _, host := range h {
		var z types.HostInfo
		j, err := host.MarshalJSON()
		if err != nil {
			log.Fatal(err)
		}
		if err := json.Unmarshal(j, &z); err != nil {
			log.Fatal(err)
		}

		t, err := plugins[z.Plugin].LoadHost(j)
		if err != nil {
			log.Fatal(err)
		}
		hosts = append(hosts, t)

		r, err := t.Run([]string{"echo", "Hello World"})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(r)

	}
}

func errorWriter(w io.Writer, err error) {
	fmt.Fprintf(w, "{'error':'%s'}", err)
}
