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

* [ ] Fix the plugin system (in progress, checkout [github.com/chabad360/plugins](https://github.com/chabad360/plugins)
  for more details)
* [x] Provide configuration system
* [x] Switch to `//go:embed` (this might need to wait)
* [x] CI/CD
* [x] Big Refactor
* [ ] Add and refactor tests (Aim for 80% Coverage)
* [ ] Complete the web UI
* [x] ~~Add a Makefile (?)~~
* [x] Refactor agent

---

## State of the Project

Covey is in active development, it's written in Go, and uses Postgres as the database. If you are interested in helping
with development, open a PR with your changes. At the moment, I'm currently polishing off the plugin system and fixing
up the CI/CD process.

### Installation Instructions

```console
$ git clone https://github.com/chabad360/covey
$ cd covey
$ goreleser build --snapshot --rm-dist
$ createdb covey
$ ./dist/Covey_linux_amd64/Covey # This will crash, it's meant to (will be fixed)
$ psql -U postgres covey <<EOF
INSERT INTO users(username, password_hash) VALUES(<username>, crypt(<password>, gen_salt('bf')));
EOF
$ ./dist/Covey_linux_amd64/Covey -plugins-folder=./plugins
```

Use the following command to build covey with live file system changes support:

```shell
$ goreleaser build --snapshot --rm-dist
$ go build -tags=live -trimpath github.com/chabad360/covey
```

--- 

### Build Signing

The builds are signed with the following public key:

https://keys.openpgp.org/vks/v1/by-fingerprint/63192DE48E0FBE543127B71C8CD4B302F58542C6

```
-----BEGIN PGP PUBLIC KEY BLOCK-----

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
hjCsLQ1WqqTk7lkBYQ==
=d1Jz
-----END PGP PUBLIC KEY BLOCK-----
```

### Fossa

[![FOSSA Status](https://app.fossa.com/api/projects/custom%2B10111%2Fgithub.com%2Fchabad360%2Fcovey.svg?type=large)](https://app.fossa.com/projects/custom%2B10111%2Fgithub.com%2Fchabad360%2Fcovey?ref=badge_large)