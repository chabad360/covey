package node

import (
	"encoding/hex"
	"github.com/chabad360/covey/models"
	"github.com/chabad360/covey/storage"
	"github.com/chabad360/covey/test"
	"github.com/google/go-cmp/cmp"
	json "github.com/json-iterator/go"
	"github.com/ory/dockertest/v3"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

var (
	// revive:disable:line-length-limit
	n = &models.Node{
		Name:       "node",
		ID:         "3778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e",
		PrivateKey: []byte("-----BEGIN RSA PRIVATE KEY-----\nMIIBOAIBAAJAc3MzlPc5PMH9Xc82hmxOBZXV7q6XnP+rr8GKzeaUkk4Q3jSJrTt8\nELVbZH2OPV3wo0sFnNCsSD3izlgp8eidVQIDAQABAkBCZCVtrR5FSmmh4N/CPdZA\ncAIu2EhoCL96uxpPfiJCX8qcUc6zu6ZY84wy6iN8I2iiBHCWsXyU/VHdbysOYIOh\nAiEAxoMoORbc0Dy+qi9khliIG/8UFtEcKUBXlyWctT3GdLsCIQCU4in24yM1R3rC\njXemM12Ks3Mt3T4+aJ0NQc22CcAdLwIgXL4F2rYdr4PRp/zAQCu4WywOnKJRP8x5\nn3nI/ru/reUCIAOa8m8zEuAwae2aJWKV7db0/34F1IMIX305sbSNyeQrAiAkhE+Z\nLe0VcQNyzkRTu+piHtcReomihMNOAs5KII5cMw==\n-----END RSA PRIVATE KEY-----"),
		PublicKey:  []byte("ssh-rsa MFswDQYJKoZIhvcNAQEBBQADSgAwRwJAc3MzlPc5PMH9Xc82hmxOBZXV7q6XnP+r\nr8GKzeaUkk4Q3jSJrTt8ELVbZH2OPV3wo0sFnNCsSD3izlgp8eidVQIDAQAB"),
		Username:   "user",
		IP:         "127.0.0.1",
	}
	n2 = &models.Node{
		Name:       "n",
		ID:         "3773ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e",
		PrivateKey: []byte("-----BEGIN RSA PRIVATE KEY-----\nMIIBOAIBAAJAc3MzlPc5PMH9Xc82hmxOBZXV7q6XnP+rr8GKzeaUkk4Q3jSJrTt8\nELVbZH2OPV3wo0sFnNCsSD3izlgp8eidVQIDAQABAkBCZCVtrR5FSmmh4N/CPdZA\ncAIu2EhoCL96uxpPfiJCX8qcUc6zu6ZY84wy6iN8I2iiBHCWsXyU/VHdbysOYIOh\nAiEAxoMoORbc0Dy+qi9khliIG/8UFtEcKUBXlyWctT3GdLsCIQCU4in24yM1R3rC\njXemM12Ks3Mt3T4+aJ0NQc22CcAdLwIgXL4F2rYdr4PRp/zAQCu4WywOnKJRP8x5\nn3nI/ru/reUCIAOa8m8zEuAwae2aJWKV7db0/34F1IMIX305sbSNyeQrAiAkhE+Z\nLe0VcQNyzkRTu+piHtcReomihMNOAs5KII5cMw==\n-----END RSA PRIVATE KEY-----"),
		PublicKey:  []byte("ssh-rsa MFswDQYJKoZIhvcNAQEBBQADSgAwRwJAc3MzlPc5PMH9Xc82hmxOBZXV7q6XnP+r\nr8GKzeaUkk4Q3jSJrTt8ELVbZH2OPV3wo0sFnNCsSD3izlgp8eidVQIDAQAB"),
		Username:   "user",
		IP:         "127.0.0.2",
	}
	// revive:enable:line-length-limit
	resource *dockertest.Resource
)

func TestNodeNew(t *testing.T) {
	var tests = []struct {
		name string
		body string
		want string
	}{
		// revive:disable:line-length-limit
		{"success", `{"ip": "127.0.0.1","port": "` + resource.GetPort("22/tcp") + `","username": "root","password":"password","name": "node1"}`,
			``},
		{"sshError", `{"ip": "127.0.0.1","port": "1000","username": "root","password": "","name": "node1"}`,
			`{"error":"dial tcp 127.0.0.1:1000: connect: connection refused"}
`},
		{"errorDuplicate", `{"ip": "127.0.0.1","port": "` + resource.GetPort("22/tcp") + `","username": "root","password":"password","name": "node1"}`,
			`{"error":"duplicate node: node1"}
`},
		{"error", `{"name":}`,
			`{"error":"models.Node.Name: ReadString: expects \" or n, but found }, error found in #9 byte of ...|{\"name\":}|..., bigger context ...|{\"name\":}|..."}
`},
		// revive:enable:line-length-limit
	}

	h := test.PureBoilerplate("POST", "/api/v1/nodes", nodeNew)

	for _, tt := range tests {
		testname := tt.name
		t.Run(testname, func(t *testing.T) {
			rr, req, err := test.HTTPBoilerplate("POST", "/api/v1/nodes", strings.NewReader(tt.body))
			if err != nil {
				t.Error(err)
			}

			h.ServeHTTP(rr, req)
			//cmp.Equal(rr.Body.Bytes()[0:5], []byte(tt.want)[0:5]) &&
			if rr.Body.String() != tt.want {
				t.Errorf("nodeNew body = %v, want %v", rr.Body.String(), tt.want)
			}
		})
	}
}

func TestNodesGet(t *testing.T) {
	storage.DB.Delete(&models.Node{}, "id != ''")
	storage.AddNode(n)
	storage.AddNode(n2)
	js, _ := json.Marshal(n2)

	var tests = []struct {
		name   string
		params string
		want   string
	}{
		// revive:disable:line-length-limit
		{"success", "sortby=name", `["` + n2.ID + `","` + n.ID + `"]
`},
		{"onlyOne", "sortby=name&limit=1", `["` + n2.ID + `"]
`},
		{"offsetOne", "sortby=name&limit=1&offset=1", `["` + n.ID + `"]
`},
		{"expandOne", "sortby=name&limit=1&expand=true", `[` + string(js) + `]
`},
		// revive:enable:line-length-limit
	}

	h := test.PureBoilerplate("GET", "/api/v1/nodes", nodesGet)

	for _, tt := range tests {
		testname := tt.name
		t.Run(testname, func(t *testing.T) {
			rr, req, err := test.HTTPBoilerplate("GET", "/api/v1/nodes?"+tt.params, nil)
			if err != nil {
				t.Fatal(err)
			}

			h.ServeHTTP(rr, req)
			if !cmp.Equal(rr.Body.Bytes(), []byte(tt.want)) && rr.Body.String() != tt.want {
				t.Errorf("nodesGet body = %v, want %v", rr.Body.String(), tt.want)
			}
		})
	}
}

func TestNodeGet(t *testing.T) {
	var tests = []struct {
		name string
		id   string
		want string
	}{
		// revive:disable:line-length-limit
		{"success", "node",
			`{"name":"node","id":"3778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e","ip":"127.0.0.1","username":"user","port":"22","CreatedAt":"2020-07-28T16:42:04.301694-07:00","UpdatedAt":"2020-07-28T16:42:04.301694-07:00"}
`},
		{"fail", "3", `{"error":"404 3 not found"}
`},
		// revive:enable:line-length-limit
	}

	h := test.PureBoilerplate("GET", "/api/v1/nodes/:node", nodeGet)

	for _, tt := range tests {
		testname := tt.name
		t.Run(testname, func(t *testing.T) {
			rr, req, err := test.HTTPBoilerplate("GET", "/api/v1/nodes/"+tt.id, nil)
			if err != nil {
				t.Fatal(err)
			}

			h.ServeHTTP(rr, req)
			if !cmp.Equal(rr.Body.Bytes()[0:10], []byte(tt.want)[0:10]) && rr.Body.String() != tt.want {
				t.Errorf("nodeGet body = %v, want %v", rr.Body.String(), tt.want)
			}
		})
	}
}

//
//func TestJobUpdate(t *testing.T) {
//	var tests = []struct {
//		name string
//		id   string
//		body string
//		want string
//	}{
//revive:disable:line-length-limit
//		{"success", "update", `{"name":"update","cron":"5 * * * *","nodes": ["test"],"tasks": {"update": {"plugin": "test","details": {"command": "hello"}}}}`,
//			`{"name":"update","id":"240875a9cf2c26d484a78b3f7f5aad21dd8f6e74031a7a5669f787d33e1b4cda","cron":"5 * * * *","nodes":["test"],"tasks":{"update":{"plugin":"test","details":{"command":"hello"}}},"task_history":[]}
//`},
//		{"error", "cron", `{"name":}`,
//			`{"error":"models.Job.Name: ReadString: expects \" or n, but found }, error found in #9 byte of ...|{\"name\":}|..., bigger context ...|{\"name\":}|..."}
//`},
//		{"404", "c", "", `{"error":"404 c not found"}
//`},
//revive:enable:line-length-limit
//	}
//
//	h := test.PureBoilerplate("PUT", "/api/v1/nodes/:node", nodeUpdate)
//
//	for _, tt := range tests {
//		testname := tt.name
//		t.Run(testname, func(t *testing.T) {
//			rr, req, err := test.HTTPBoilerplate("PUT", "/api/v1/nodes/"+tt.id, strings.NewReader(tt.body))
//			if err != nil {
//				t.Fatal(err)
//			}
//
//			h.ServeHTTP(rr, req)
//			if !cmp.Equal(rr.Body.Bytes()[0:10], []byte(tt.want)[0:10]) && rr.Body.String() != tt.want {
//				t.Errorf("nodeNew body = %v, want %v", rr.Body.String(), tt.want)
//			}
//		})
//	}
//}

func TestNodeDelete(t *testing.T) {
	//revive:disable:line-length-limit
	var tests = []struct {
		name string
		id   string
		want string
	}{
		{"success", "node", `"node"
`},
		{"fail", "3", `{"error":"404 3 not found"}
`},
	}
	//revive:enable:line-length-limit

	h := test.PureBoilerplate("DELETE", "/api/v1/nodes/:node", nodeDelete)

	for _, tt := range tests {
		testname := tt.name
		t.Run(testname, func(t *testing.T) {
			rr, req, err := test.HTTPBoilerplate("DELETE", "/api/v1/nodes/"+tt.id, nil)
			if err != nil {
				t.Fatal(err)
			}

			h.ServeHTTP(rr, req)
			if rr.Body.String() != tt.want {
				t.Errorf("nodeDelete body = %v, want %v", rr.Body.String(), tt.want)
			}
		})
	}
}

func TestMain(m *testing.M) {
	pool, r, pdb, err := test.Boilerplate()
	storage.DB = pdb
	if err != nil {
		log.Fatalf("Could not setup DB connection: %s", err)
	}

	resource, err = pool.Run("chabad360/sshd", "latest", []string{"SSH_ENABLE_ROOT=true", "SSH_ENABLE_PASSWORD_AUTH=true"})
	if err != nil {
		log.Fatalf("Could not setup sshd container: %v", err)
	}

	time.Sleep(time.Second * 5)

	// revive:disable:line-length-limit
	n.HostKey, err = hex.DecodeString("0000001365636473612d736861322d6e69737470323536000000086e6973747032353600000041044032b5eed25ed08ec4361d9f7e6a7e27f725d563bc033f777fe2b12bdd61c86c160476c6d080b1361ea4ab9e89ec104051762ecb0a4595f53a16a06c959a0704")
	n2.HostKey, err = hex.DecodeString("0000001365636473612d736861322d6e69737470323536000000086e697374703235360000004104d7a16df1c77e969f989b46a515f7cd2f4d1848d5e121ed9081edd9d8a76555bf61e203b56f84ecb61a3ebf63a5c26e4c9755ebcaae942949088536b98064f286") // revive:enable:line-length-limit
	storage.AddNode(n)

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
	if err := pool.Purge(r); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}
