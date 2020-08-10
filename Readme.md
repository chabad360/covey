# Covey

## A lightweight Linux cluster orchestration server written in Go

[![Codacy Badge](https://app.codacy.com/project/badge/Grade/b6e797a0fb5a498199b2a2d3ae494c82)](https://www.codacy.com/manual/chabad360/covey?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=chabad360/covey&amp;utm_campaign=Badge_Grade)
[![Codacy Badge](https://app.codacy.com/project/badge/Coverage/b6e797a0fb5a498199b2a2d3ae494c82)](https://www.codacy.com/manual/chabad360/covey?utm_source=github.com&utm_medium=referral&utm_content=chabad360/covey&utm_campaign=Badge_Coverage)
[![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/4095/badge)](https://bestpractices.coreinfrastructure.org/projects/4095)
[![Chat on discord](https://img.shields.io/discord/727820939013783582?logo=discord&logoColor=white)](https://discord.gg/kWXPrWg)

Covey is a project designed to fill a certain void, the lack of a nice lightweight cluster management system.
There are tools like Rundeck (which Covey takes after), that are quite capable, but are far too heavy to be useful.

### Features

Covey has, and will gain (in the coming weeks/months) a variety of features, including:

* RESTful API (with swagger docs!)
* Plug-able Modules (ish, it's very broken at the moment)
* Web Interface (basic)
* Node Monitoring (using Netdata)
* Automated Setup (Almost there)
* Crash-only design (pretty much)

---

### Current Roadmap

#### V0.1 MVP

* [x] Create MVP

#### V0.2 The Refactor

* [x] Major Refactor

#### V0.3 Web Interface

* [x] Design and implement the basic web interface
* [x] Implement basic authentication
* [x] Begin adding tests
* [x] Switch away from mux (too slow...)

#### V0.4 Monitoring

* [x] Fix some issues with the plugin system
* [x] Rework the task module
* [x] Create Node Agent
* [x] Persistent queue
* [x] Add relevant UI elements
* [x] Automatically install agent
* [x] Add SystemD service file to the agent
* [x] Integrate with [Netdata](https://github.com/netdata/netdata) for monitoring
* [ ] ~~Add tests~~

#### V0.5 A Better API

* [x] ~~Evaluate designing a very basic framework (for keeping things cleaner)~~
* [x] ~~Evaluate GraphQL for the API~~
* [x] Redesign DB using Gorm
* [x] Fully implement (and test) the API
* [x] Swagger (OpenAPI)
* [x] Fully document the API

#### V0.6 Alpha

* [ ] Fix the plugin system
* [x] Provide configuration system
* [x] Deal with packed files (`.gitignore` and then include it on build?)
* [ ] CI/CD
* [x] Big Refactor
* [ ] Add and refactor tests (Aim for 80% Coverage)
* [ ] Complete the web UI
* [ ] Add a Makefile (?)
* [x] Refactor agent

---

## State of the Project

Covey is in active development, it's written in Go, and uses Postgres as the database.
If you are interested in helping with development, open a PR with your changes.
At the moment, I've been finishing off most of my tests, and fixing bugs.

### Installation Instructions

```bash
git clone https://github.com/chabad360/covey
cd covey

go mod download
go get github.com/omeid/go-resources/cmd/resources

go build -ldflags="-s -w" -trimpath -v -o assets/agent github.com/chabad360/covey/agent
upx assets/agent/agent

resources -declare -package=asset -output=asset/asset.go -tag="!live" -trim assets/ ./assets/*

go build -trimpath -buildmode=plugin -o plugins/task/shell.so github.com/chabad360/covey/plugins/task/shell
go build -trimpath github.com/chabad360/covey

createdb covey

./covey # This will crash, it's meant to (will be fixed)
psql -U postgres covey <<EOF
INSERT INTO users(username, password_hash) VALUES(<username>, crypt(<password>, gen_salt('bf')));
EOF

./covey -plugins-folder=./plugins
```

Use the following command to build covey with live file system changes support:

```bash
go build -tags live -trimpath github.com/chabad360/covey
```

Use the following for a fancy release build:

```bash
go build -trimpath -ldflags="-s -w" github.com/chabad360/covey && upx covey
```

--- 

### Fossa

[![FOSSA Status](https://app.fossa.com/api/projects/custom%2B10111%2Fgithub.com%2Fchabad360%2Fcovey.svg?type=large)](https://app.fossa.com/projects/custom%2B10111%2Fgithub.com%2Fchabad360%2Fcovey?ref=badge_large)