# ucron

**unprivileged cron** aka **ushi's cron**

A cron implementation meant to be executed by unprivileged users. This is how you do cron inside a docker container.

## Usage

```shell
$ cat path/to/crontab
# min   hour  dom   month  dow   command
0-58/2  *     *     *      *     echo even minutes
$ ucron path/to/crontab
ucron: /bin/sh -c "echo even minutes"
even minutes
```

## Features / Design

- No root required
- Reads a single crontab from file or stdin
- Logs all job output to stdout
- Waits for all running jobs to complete on shutdown
