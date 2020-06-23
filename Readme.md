# Covey

## A lightweight Linux cluster orchestration server written in Go

[![Codacy Badge](https://app.codacy.com/project/badge/Grade/b6e797a0fb5a498199b2a2d3ae494c82)](https://www.codacy.com/manual/chabad360/covey?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=chabad360/covey&amp;utm_campaign=Badge_Grade) [![Codacy Badge](https://app.codacy.com/project/badge/Coverage/b6e797a0fb5a498199b2a2d3ae494c82)](https://www.codacy.com/manual/chabad360/covey?utm_source=github.com&utm_medium=referral&utm_content=chabad360/covey&utm_campaign=Badge_Coverage)

Covey is a project designed to fill a certain void, the lack of a nice lightweight cluster management system. The are tools like Rundeck (which Covey takes after), that are quite capable, but are far too heavy to be useful.

### Features

Covey has, and will gain (in the coming weeks/months) a variety of features, including:

* RESTful API
* Plug-able Modules
* Web Interface (basic and doesn't work with nodes yet)
* Node Monitoring (planned for v0.4)
* Automated Setup
* Crash-only design

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

* [ ] Fix some issues With the plugin system
* [ ] Rework the task module
* [ ] Create Node Agent
* [ ] Integrate with [Netdata](https://github.com/netdata/netdata) for monitoring
* [ ] Add tests
* [ ] Add relevant UI elements

#### V0.5 A Better API

* [ ] Evaluate designing a very basic framework (for keeping things cleaner)
* [ ] Evaluate GraphQL for the API
* [ ] Fully implement (and test) the API
* [ ] Swagger (OpenAPI) (might save that for later)
* [ ] Fully document the API

#### V0.6 Alpha

* [ ] Fix the plugin system
* [ ] Provide configuration system
* [ ] CI
* [ ] Refactor
* [ ] Add an AUR package
* [ ] Add and refactor tests (Aim for 80% Coverage)
* [ ] Complete the web UI

---

## State of the Project

Covey is in active development, it's written in Go, and uses Postgres as the database. If you are interested in helping with development, open a PR with your changes. At the moment, I've begun working on the web interface.

### Installation Instructions

```bash
git clone https://github.com/chabad360/covey
cd covey

go build -trimpath -buildmode=plugin -o plugins/task/shell.so github.com/chabad360/covey/plugins/task/shell
go build -trimpath -buildmode=plugin -o plugins/node/ssh.so github.com/chabad360/covey/plugins/node/ssh
go build -trimpath github.com/chabad360/covey

createdb covey
psql covey < structure.sql

./covey
```

Use the following command to build covey with live file system changes support:

```bash
go build -tags live -trimpath github.com/chabad360/covey
```

Use the following for a fancy release build:

```bash
go build -trimpath -ldflags="-s -w" github.com/chabad360/covey && upx covey
```
