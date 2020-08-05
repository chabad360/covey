package storage

import (
	"github.com/chabad360/covey/models"
	"github.com/google/go-cmp/cmp"
	"testing"
)

// revive:disable:line-length-limit
var n = &models.Node{
	Name:       "node",
	ID:         "3778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e",
	PrivateKey: []byte("-----BEGIN RSA PRIVATE KEY-----\nMIIBOAIBAAJAc3MzlPc5PMH9Xc82hmxOBZXV7q6XnP+rr8GKzeaUkk4Q3jSJrTt8\nELVbZH2OPV3wo0sFnNCsSD3izlgp8eidVQIDAQABAkBCZCVtrR5FSmmh4N/CPdZA\ncAIu2EhoCL96uxpPfiJCX8qcUc6zu6ZY84wy6iN8I2iiBHCWsXyU/VHdbysOYIOh\nAiEAxoMoORbc0Dy+qi9khliIG/8UFtEcKUBXlyWctT3GdLsCIQCU4in24yM1R3rC\njXemM12Ks3Mt3T4+aJ0NQc22CcAdLwIgXL4F2rYdr4PRp/zAQCu4WywOnKJRP8x5\nn3nI/ru/reUCIAOa8m8zEuAwae2aJWKV7db0/34F1IMIX305sbSNyeQrAiAkhE+Z\nLe0VcQNyzkRTu+piHtcReomihMNOAs5KII5cMw==\n-----END RSA PRIVATE KEY-----"),
	PublicKey:  []byte("ssh-rsa MFswDQYJKoZIhvcNAQEBBQADSgAwRwJAc3MzlPc5PMH9Xc82hmxOBZXV7q6XnP+r\nr8GKzeaUkk4Q3jSJrTt8ELVbZH2OPV3wo0sFnNCsSD3izlgp8eidVQIDAQAB"),
	Username:   "user",
	IP:         "127.0.0.1",
}

// revive:enable:line-length-limit

func TestAddNode(t *testing.T) {
	//revive:disable:line-length-limit
	var tests = []struct {
		id   string
		want *models.Node
	}{
		{"3778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e", n},
		{"3", &models.Node{}},
	}
	//revive:enable:line-length-limit

	testError := AddNode(n)

	for _, tt := range tests {
		testname := tt.id
		t.Run(testname, func(t *testing.T) {
			var got models.Node
			if DB.Where("id = ?", tt.id).First(&got); cmp.Equal(got, tt.want) {
				t.Errorf("addNode() = %v, want %v, error: %v", got, tt.want, testError)
			}
		})
	}
}

func TestGetNodeIDOrName(t *testing.T) {
	var tests = []struct {
		name  string
		id    string
		field string
		want  string
		want2 bool
	}{
		{"ok_ID", n.Name, "id", n.ID, true},
		{"notok_ID", "n", "id", "", false},
		{"ok_Name", n.ID, "name", n.Name, true},
		{"notok_Name", "n", "name", "", false},
	}

	for _, tt := range tests {
		testname := tt.name
		t.Run(testname, func(t *testing.T) {
			got, got2 := GetNodeIDorName(tt.id, tt.field)
			if got2 != tt.want2 {
				t.Errorf("GetNodeIDorName() = %v, want %v", got2, tt.want2)
			}
			if got != tt.want {
				t.Errorf("GetNodeIDorName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetNode(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name  string
		args  args
		want  *models.Node
		want1 bool
	}{
		{"success", args{n.ID}, n, true},
		{"fail", args{"3"}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got1 := GetNode(tt.args.id)
			if got1 != tt.want1 {
				t.Errorf("GetNode() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
