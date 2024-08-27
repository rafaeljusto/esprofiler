# 📊 Elasticsearch Profiler

[![Go Reference](https://pkg.go.dev/badge/github.com/rafaeljusto/esprofiler.svg)](https://pkg.go.dev/github.com/rafaeljusto/esprofiler)
![Test](https://github.com/rafaeljusto/esprofiler/actions/workflows/test.yml/badge.svg)

This is a tool to visualize the Elasticsearch profiler output described
[here](https://www.elastic.co/guide/en/elasticsearch/reference/current/search-profile.html).

[![Elasticsearch Profiler Video](https://img.youtube.com/vi/0qegA-EDIl4/0.jpg)](https://www.youtube.com/watch?v=0qegA-EDIl4)

## ⚙️ Features

* Use [Flame Graphs](https://www.brendangregg.com/flamegraphs.html) to visualize the profiler output.
* Split the results per indexes and shards.

## ⚡️ Quick start

The easiest way to start using this tool is to use the
[Docker](https://www.docker.com/) image. Just run the following command:

```bash
docker run -p 8080:80 rafaeljusto/esprofiler:latest
```

Open your browser at [http://localhost:8080](http://localhost:8080).