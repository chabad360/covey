# Covey

## A lightweight Linux cluster orchestration server written in Go

[![Codacy Badge](https://app.codacy.com/project/badge/Grade/b6e797a0fb5a498199b2a2d3ae494c82)](https://www.codacy.com/manual/chabad360/covey?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=chabad360/covey&amp;utm_campaign=Badge_Grade) [![Codacy Badge](https://app.codacy.com/project/badge/Coverage/b6e797a0fb5a498199b2a2d3ae494c82)](https://www.codacy.com/manual/chabad360/covey?utm_source=github.com&utm_medium=referral&utm_content=chabad360/covey&utm_campaign=Badge_Coverage)

Covey is a project designed to fill a certain void, the lack of a nice lightweight cluster management system. The are tools like Rundeck (which Covey takes after), that are quite capable, but are far too heavy to be useful.

### Features

Covey has, and will gain (in the coming weeks/months) a variety of features, including:

* RESTful API
* Plug-able Modules
* Web Interface (coming soon with v0.3)
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

* [ ] Design and Implement the Web Interface
* [x] Implement Basic Authentication
* [x] Begin Adding Tests
* [x] Switch away from mux (too slow...)

#### V0.4 Monitoring

* [ ] Create Node Agent
* [ ] Integrate With [Netdata](https://github.com/netdata/netdata) for Monitoring
* [ ] Add and Refactor Tests (Aim for 80% Coverage)

#### V0.5 A Better API

* [ ] Evaluate GraphQL for the API
* [ ] Fully Implement (and test) the API
* [ ] Swagger (OpenAPI) (Might save that for later)
* [ ] Fully Document the API

#### V0.6 Alpha

* [ ] Provide Configuration Methods
* [ ] Provide Build Artifacts
* [ ] Fix Some Issues With the Plugin System
* [ ] Refactor
* [ ] Add an AUR Package

---

## State of the Project

Covey is in active development, it's written in Go, and uses Postgres as the database. If you are interested in helping with development, open a PR with your changes. At the moment, I've been adding some tests, and doing some refactoring here and there.

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
