# Covey

## A lightweight Linux cluster orchestration server written in Go

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
* [ ] Implement Basic Authentication
* [ ] Add Tests

#### V0.4 Monitoring

* [ ] Create Node Agent
* [ ] Implement Plug-able Monitoring Interface

---

## State of the Project

Covey is in active development, it's written in Go, and uses Postgres as the database. If you are interested in helping with development, open a PR with your changes. At the moment, I'm beginning to work on refacoring the code, and will be accepting PRs for that.

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
