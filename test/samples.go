package test

import (
	"github.com/chabad360/covey/models"
)

var (
	//revive:disable:exported
	//revive:disable:line-length-limit
	U1 = models.User{
		Username: "user",
		Password: "password",
	}
	U2 = models.User{
		Username: "user",
		Password: "pass",
	}
	U3 = models.User{
		Username: "user2",
		Password: "password",
	}
	J1 = models.Job{
		Name:  "update",
		ID:    "3778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e",
		Nodes: []string{"node1"},
		Tasks: map[string]models.JobTask{
			"update": {
				Plugin:  "shell",
				Details: map[string]string{"command": "sudo apt update && sudo apt upgrade -y"},
			},
		},
	}
	J2 = models.Job{
		Name:  "add",
		ID:    "3748ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e",
		Nodes: []string{"node1"},
		Tasks: map[string]models.JobTask{
			"update": {
				Plugin:  "shell",
				Details: map[string]string{"command": "sudo apt update && sudo apt upgrade -y"},
			},
		},
	}
	N1 = models.Node{
		Name:       "node",
		ID:         "3778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e",
		PrivateKey: []byte("-----BEGIN RSA PRIVATE KEY-----\nMIIBOAIBAAJAc3MzlPc5PMH9Xc82hmxOBZXV7q6XnP+rr8GKzeaUkk4Q3jSJrTt8\nELVbZH2OPV3wo0sFnNCsSD3izlgp8eidVQIDAQABAkBCZCVtrR5FSmmh4N/CPdZA\ncAIu2EhoCL96uxpPfiJCX8qcUc6zu6ZY84wy6iN8I2iiBHCWsXyU/VHdbysOYIOh\nAiEAxoMoORbc0Dy+qi9khliIG/8UFtEcKUBXlyWctT3GdLsCIQCU4in24yM1R3rC\njXemM12Ks3Mt3T4+aJ0NQc22CcAdLwIgXL4F2rYdr4PRp/zAQCu4WywOnKJRP8x5\nn3nI/ru/reUCIAOa8m8zEuAwae2aJWKV7db0/34F1IMIX305sbSNyeQrAiAkhE+Z\nLe0VcQNyzkRTu+piHtcReomihMNOAs5KII5cMw==\n-----END RSA PRIVATE KEY-----"),
		PublicKey:  []byte("ssh-rsa MFswDQYJKoZIhvcNAQEBBQADSgAwRwJAc3MzlPc5PMH9Xc82hmxOBZXV7q6XnP+r\nr8GKzeaUkk4Q3jSJrTt8ELVbZH2OPV3wo0sFnNCsSD3izlgp8eidVQIDAQAB"),
		HostKey:    []byte{0, 0, 0, 19, 101, 99, 100, 115, 97, 45, 115, 104, 97, 50, 45, 110, 105, 115, 116, 112, 50, 53, 54, 0, 0, 0, 8, 110, 105, 115, 116, 112, 50, 53, 54, 0, 0, 0, 65, 4, 64, 50, 181, 238, 210, 94, 208, 142, 196, 54, 29, 159, 126, 106, 126, 39, 247, 37, 213, 99, 188, 3, 63, 119, 127, 226, 177, 43, 221, 97, 200, 108, 22, 4, 118, 198, 208, 128, 177, 54, 30, 164, 171, 158, 137, 236, 16, 64, 81, 118, 46, 203, 10, 69, 149, 245, 58, 22, 160, 108, 149, 154, 7, 4},
		Username:   "user",
		IP:         "127.0.0.1",
	}
	N2 = models.Node{
		Name:       "n",
		ID:         "3773ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e",
		PrivateKey: []byte("-----BEGIN RSA PRIVATE KEY-----\nMIIBOAIBAAJAc3MzlPc5PMH9Xc82hmxOBZXV7q6XnP+rr8GKzeaUkk4Q3jSJrTt8\nELVbZH2OPV3wo0sFnNCsSD3izlgp8eidVQIDAQABAkBCZCVtrR5FSmmh4N/CPdZA\ncAIu2EhoCL96uxpPfiJCX8qcUc6zu6ZY84wy6iN8I2iiBHCWsXyU/VHdbysOYIOh\nAiEAxoMoORbc0Dy+qi9khliIG/8UFtEcKUBXlyWctT3GdLsCIQCU4in24yM1R3rC\njXemM12Ks3Mt3T4+aJ0NQc22CcAdLwIgXL4F2rYdr4PRp/zAQCu4WywOnKJRP8x5\nn3nI/ru/reUCIAOa8m8zEuAwae2aJWKV7db0/34F1IMIX305sbSNyeQrAiAkhE+Z\nLe0VcQNyzkRTu+piHtcReomihMNOAs5KII5cMw==\n-----END RSA PRIVATE KEY-----"),
		PublicKey:  []byte("ssh-rsa MFswDQYJKoZIhvcNAQEBBQADSgAwRwJAc3MzlPc5PMH9Xc82hmxOBZXV7q6XnP+r\nr8GKzeaUkk4Q3jSJrTt8ELVbZH2OPV3wo0sFnNCsSD3izlgp8eidVQIDAQAB"),
		HostKey:    []byte{0, 0, 0, 19, 101, 99, 100, 115, 97, 45, 115, 104, 97, 50, 45, 110, 105, 115, 116, 112, 50, 53, 54, 0, 0, 0, 8, 110, 105, 115, 116, 112, 50, 53, 54, 0, 0, 0, 65, 4, 215, 161, 109, 241, 199, 126, 150, 159, 152, 155, 70, 165, 21, 247, 205, 47, 77, 24, 72, 213, 225, 33, 237, 144, 129, 237, 217, 216, 167, 101, 85, 191, 97, 226, 3, 181, 111, 132, 236, 182, 26, 62, 191, 99, 165, 194, 110, 76, 151, 85, 235, 202, 174, 148, 41, 73, 8, 133, 54, 185, 128, 100, 242, 134},
		Username:   "user",
		IP:         "127.0.0.2",
	}
	T1 = models.Task{
		ID:       "3778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e",
		State:    models.StateQueued,
		Plugin:   "test",
		Details:  map[string]string{"test": "test"},
		ExitCode: 258,
	}
	T2 = models.Task{
		ID:       "277dffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6c",
		State:    models.StateRunning,
		Plugin:   "test",
		Node:     "test",
		Details:  map[string]string{"test": "test"},
		ExitCode: 0,
	}
	JWT1 = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjb3ZleS1hcGkiLCJzdWIiOiIzIiwiYXVkIjoiYWxsIiwiZXhwIjo0NDY5NTQyMzYzLCJpYXQiOjE1OTU1MzI3NjMsImp0aSI6InNseUM0QVd2N0NhU3RKWG9yeXF5QzFPOWZLUFJzdFZQIn0.AAU9I8yub7VmTCnT833F54W6uQbhGVFKR8DSsi9pDJI"
	JWT2 = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjb3ZleS1hcGkiLCJzdWIiOiIxIiwiYXVkIjoiYWxsIiwiZXhwIjowLCJpYXQiOjE1OTE5MTI5NzAsImp0aSI6InBIWWp4ZVVCclZmZHdVeldIUmloRkRQUkhCTXVFV21hIn0.XiNKXNDmsxXul8ceyQUgBWJBfrUmBsHWyLC34_Qy5qo"
	JWT3 = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjb3ZleS1hcGkiLCJzdWIiOiIxIiwiYXVkIjoiYWxsIiwiZXhwIjowLCJpYXQiOjE1OTE5MTI5NzAsImp0aSI6InBIWWp4ZVVCclZmZHdVeldIUmloRkRQUkhCTXVFV21hIna.XiNKXNDmsxXul8ceyQUgBWJBfrUmBsHWyLC34_Qy5qo"
	JWT4 = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjb3ZleS1hcGkiLCJzdWIiOiIxIiwiYXVkIjoiYWxsIiwiZXhwIjowLCJpYXQiOjE1OTE5MTI5NzAsImp0aSI6InBIWWp4ZVVCclZmZHdVeldIUmloRkRQUkhCTXVFV21hIna.XiNKXNDmsxXul8ceyQUgBWJBfrUmBsHWyLC34_Qy5qo"
	JWT5 = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjb3ZleS11c2VyIiwic3ViIjoiMSIsImF1ZCI6ImFsbCIsImV4cCI6NDQ3MjIxNTg5OSwiaWF0IjoxNTk4MjA2Mjk5LCJqdGkiOiJ1cElBclJJNWZYNTBXaTNzTExPZHhycUZnSFZ6WVVZaCJ9.LB1Ji0YvNwLfnAyogQ4XZ9pXILkxVAd0kUAz2FYnHWM"
)
