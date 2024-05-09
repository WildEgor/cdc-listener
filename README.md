<h1 align="center"> CDC Listener </h1> <br>
<hr>
<div>
    <img src="assets/logo.png" width="200" height="200" style="display: block;margin-left: auto;margin-right: auto;">
</div>
<p align="center">
  A service helps for tracking MongoDB data changes asynchronously
</p>
<hr>

[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/WildEgor/cdc-listener)
[![Go Report Card](https://goreportcard.com/badge/github.com/WildEgor/cdc-listener)](https://goreportcard.com/report/github.com/WildEgor/cdc-listener)
![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/WildEgor/cdc-listener)
[![Publish Docker image](https://github.com/WildEgor/cdc-listener/actions/workflows/publish.yml/badge.svg)](https://github.com/WildEgor/cdc-listener/actions/workflows/publish.yml)

## Table of Contents
- [Introduction](#introduction)
- [Features](#features)
- [Requirements](#requirements)
- [Quick Start](#quick-start)
- [Contributing](#contributing)

## Introduction

A service helps for tracking MongoDB data changes asynchronously.
Change Data Capture (CDC) is a design pattern that allows you to track and capture change streams from a source data 
system so that downstream systems can efficiently process these changes.
It's useful for **"Event-Driven architecture"** and implementation patterns like **"Outbox"** - just insert or updates data in database collections and subscribe to topic(s).
Service publish event to topic (see below):
```json
{
  "id": "ObjectID string",
  "db": "db name",
  "collection": "collection name",
  "action": "insert/update/delete operation types",
  "data": {},
  "data_old": {},
  "event_time": "ISO time"
}
```

Messages are published to the broker at least once!

### Filter configuration example

```yaml
filter:
- db: "test"
  collections:
    test:
      - insert
      - update
```
This filter means that we only process events occurring with the `test` database and `test` collection,
and in particular `insert` and `update` data.

### Topic mapping
By default, using topic provided to `publisher.topic` option,
but if you want to send all update in one topic you should be configured the topic map:
```yaml
topicsMap:
  test-test: "notifier" # [database name]-[collection name]: [topic name]
```

#### Available metrics

| name                        | description                          | fields             | status  |
|-----------------------------|--------------------------------------|--------------------|---------|
| published_events_total      | the total number of published events | `subject`, `table` | * WIP * |
| filter_skipped_events_total | the total number of skipped events   | `table`            | * WIP * |

## Features

- [x] Handle database(s) collection(s) changes;
- [x] Publish using RabbitMQ `type=rabbitmq`;
- [**TODO**] Support other publishers;
- [**TODO**] Old changes included with new document in event;
- [**TODO**] Resume token logic;
- [**TODO**] Metrics;
- [**TODO**] Health checks.

## Requirements

- [Git](http://git-scm.com/)
- [Go >= 1.22](https://go.dev/dl/)
- [Docker](https://www.docker.com/products/docker-desktop/)
- [MongoDB](https://www.mongodb.com/)

## Quick start

1. MongoDB must be replicated;
2. Create `config.yml`
```yaml
# Application level configs
app:
  name: "cdc-listener"
  port: 8080
  mode: "develop"
logger:
  level: info
  format: json
# MongoDB listener config
listener:
  filter:
    - db: "test"
      collections:
        test:
          - insert
          - update
  topicsMap:
    test-test: "notifier"
# MongoDB config
database:
  uri: "mongodb://127.0.0.1:27017,127.0.0.1:27018,127.0.0.1:27019/?replicaSet=rs0"
  debug: false
# Publisher config
publisher:
  type: rabbitmq
  uri: "amqp://guest:guest@localhost:5672/"
  topic: "notifier" # used as default topic
  topicPrefix: ""
# Prometheus config
monitoring:
  promAddr: ":2112"
```

3. Run using `air` or `docker`.
```shell
docker-compose up -d build
```

## Contributing

Please, use ```git cz``` for commit messages!

```shell
git clone https://github.com/WildEgor/cdc-listener
cd cdc-listener
git checkout -b feature-or-fix-branch
git add .
git cz
git push --set-upstream-to origin/feature-or-fix-branch
```