package task

import (
	"github.com/chabad360/covey/models"
)

var (
	t1 = &models.Task{
		ID:       "2778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6c",
		State:    models.StateRunning,
		Plugin:   "test",
		Node:     "test1",
		Details:  map[string]string{"test": "test"},
		ExitCode: 0,
	}
	t2 = &models.Task{
		ID:       "277dffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6c",
		State:    models.StateRunning,
		Plugin:   "test",
		Node:     "test",
		Details:  map[string]string{"test": "test"},
		ExitCode: 0,
	}
	// revive:disable:line-length-limit
	node1 = &models.Node{
		Name:       "node",
		ID:         "3778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e",
		PrivateKey: []byte("-----BEGIN RSA PRIVATE KEY-----\nMIIBOAIBAAJAc3MzlPc5PMH9Xc82hmxOBZXV7q6XnP+rr8GKzeaUkk4Q3jSJrTt8\nELVbZH2OPV3wo0sFnNCsSD3izlgp8eidVQIDAQABAkBCZCVtrR5FSmmh4N/CPdZA\ncAIu2EhoCL96uxpPfiJCX8qcUc6zu6ZY84wy6iN8I2iiBHCWsXyU/VHdbysOYIOh\nAiEAxoMoORbc0Dy+qi9khliIG/8UFtEcKUBXlyWctT3GdLsCIQCU4in24yM1R3rC\njXemM12Ks3Mt3T4+aJ0NQc22CcAdLwIgXL4F2rYdr4PRp/zAQCu4WywOnKJRP8x5\nn3nI/ru/reUCIAOa8m8zEuAwae2aJWKV7db0/34F1IMIX305sbSNyeQrAiAkhE+Z\nLe0VcQNyzkRTu+piHtcReomihMNOAs5KII5cMw==\n-----END RSA PRIVATE KEY-----"),
		PublicKey:  []byte("ssh-rsa MFswDQYJKoZIhvcNAQEBBQADSgAwRwJAc3MzlPc5PMH9Xc82hmxOBZXV7q6XnP+r\nr8GKzeaUkk4Q3jSJrTt8ELVbZH2OPV3wo0sFnNCsSD3izlgp8eidVQIDAQAB"),
		Username:   "user",
		IP:         "127.0.0.1",
	}
	// revive:enable:line-length-limit
)

//func TestGetTask(t *testing.T) {
//	type args struct {
//		identifier string
//	}
//	tests := []struct {
//		name  string
//		args  args
//		want  *models.Task
//		want1 bool
//	}{
//		{"db", args{"3778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e"}, task, true},
//		{"noDB", args{"31b079725d0a20bfe6c3b6e"}, nil, false},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, got1 := getTask(tt.args.identifier)
//			if got1 == tt.want1 && got.ID != tt.want.ID {
//				t.Errorf("getTask() got = %v, want %v, got1 = %v, want %v", got, tt.want, got1, tt.want1)
//			}
//		})
//	}
//}
