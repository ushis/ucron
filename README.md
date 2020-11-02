# ucron

**unprivileged cron** aka **ushi's cron**

A cron implementation meant to be executed by unprivileged users. The project aims to provide a hassle free solution for running cron inside a docker container.

## Install

Download a build from the [Releases](https://github.com/ushis/ucron/releases) page.

```shell
$ curl -LO https://github.com/ushis/ucron/releases/download/<version>/ucron
$ chmod +x ucron
```

Download the build from the latest release.

```shell
$ curl https://api.github.com/repos/ushis/ucron/releases/latest |
  jq '.assets | map(select(.name == "ucron")) | first | .browser_download_url' |
  xargs curl -LO
$ chmod +x ucron
```

Or install it via `go get`.

```shell
$ go get github.com/ushis/ucron
```

## Usage

```shell
$ ucron path/to/crontab
```

## Features / Design

- No root required
- Reads a single crontab from file or stdin
- Logs all job output to stdout
- Waits for all running jobs to complete on shutdown

## Q&A

**Is ucron capable of running multiple crontabs?** Of course!

```shell
$ cat crontab1 crontab2 crontab3 | ucron -
```
