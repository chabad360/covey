name: Go

on:
  push:
    branches: [ master ]
  pull_request:
    branches: [ master ]

jobs:

  build:
    name: Build
    runs-on: ubuntu-latest
    steps:

    - name: Set up Go 1.x
      uses: actions/setup-go@v2
      with:
        go-version: ^1.15

    - name: Check out code into the Go module directory
      uses: actions/checkout@v2

    - name: Get dependencies
      run: |
        go get -v -t -d ./...
        if [ -f Gopkg.toml ]; then
            curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
            dep ensure
        fi
        go get github.com/t-yuki/gocover-cobertura
        go mod download
        go get github.com/omeid/go-resources/cmd/resources
        
    - name: Build agent
      run: go build -ldflags="-s -w" -trimpath -v -o assets/agent github.com/chabad360/covey/agent
    
    - name: UPX agent
      uses: crazy-max/ghaction-upx@v1
      with:
        version: latest
        file: assets/agent/agent

    - name: Pack assets
      run: resources -declare -package=asset -output=asset/asset.go -tag="!live" -trim assets/ ./assets/*
    
    - name: Test
      if: success() || failure()
      run: go test -trimpath -v -coverprofile=coverage.txt -covermode count `go list ./... | grep -v test`
      
    - name: Upload coverage
      if: success() || failure()
      env:
        CODACY_PROJECT_TOKEN: ${{ secrets.CODACY_PROJECT_TOKEN }}
      run: |
        gocover-cobertura < coverage.txt > coverage.xml
        bash <(curl -Ls https://coverage.codacy.com/get.sh) \
          report --language Go --force-language -r coverage.xml
