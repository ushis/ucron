package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"unicode"

	"github.com/robfig/cron/v3"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <filename>|-\n", os.Args[0])
		os.Exit(1)
	}
	crontab, err := openCrontab(os.Args[1])

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer crontab.Close()

	entries, err := readCrontab(crontab)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	c := cron.New()

	for _, entry := range entries {
		if _, err := c.AddJob(entry[0], newJob(entry[1])); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
	c.Start()

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	<-sigc

	fmt.Println("ucron: shutdown: waiting for jobs to complete...")
	<-c.Stop().Done()
}

func openCrontab(filename string) (io.ReadCloser, error) {
	if filename == "-" {
		return os.Stdin, nil
	}
	return os.Open(filename)
}

func readCrontab(r io.Reader) ([][2]string, error) {
	b := bufio.NewReader(r)
	entries := [][2]string{}

	for {
		line, err := b.ReadString('\n')

		if err == io.EOF {
			return entries, nil
		}
		if err != nil {
			return nil, err
		}
		if strings.HasPrefix(line, "#") {
			continue
		}
		entry, err := parseEntry(line)

		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
}

func parseEntry(entry string) ([2]string, error) {
	for count, i := 0, 0; ; {
		j := strings.IndexFunc(entry[i:], func(r rune) bool {
			return !unicode.IsSpace(r)
		})

		if j < 0 {
			return [2]string{}, fmt.Errorf("invalid crontab entry: %s", entry)
		}
		k := strings.IndexFunc(entry[i+j:], unicode.IsSpace)

		if k < 0 {
			return [2]string{}, fmt.Errorf("invalid crontab entry: %s", entry)
		}
		i += j + k

		if count == 4 {
			return [2]string{strings.TrimSpace(entry[:i]), strings.TrimSpace(entry[i:])}, nil
		}
		count++
	}
}

type job struct {
	command string
}

func newJob(command string) *job {
	return &job{command: command}
}

func (j *job) Run() {
	fmt.Printf("ucron: /bin/sh -c \"%s\"\n", j.command)
	cmd := exec.Command("/bin/sh", "-c", j.command)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true, Pgid: 0}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout

	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
