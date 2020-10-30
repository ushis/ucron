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

	fmt.Println("> shutdown. waiting for jobs to complete...")
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

		if err != nil && err != io.EOF {
			return nil, err
		}
		if line == "" && err == io.EOF {
			return entries, nil
		}
		line = strings.TrimSpace(line)

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := fieldsN(line, 6)

		if len(fields) != 6 {
			return nil, fmt.Errorf("invalid crontab entry: %s", line)
		}
		entries = append(entries, [2]string{
			strings.Join(fields[0:5], " "),
			fields[5],
		})
	}
}

func fieldsN(s string, n int) []string {
	fields := []string{}
	i := 0

	for ; n != 1; n-- {
		j := strings.IndexFunc(s[i:], func(r rune) bool {
			return !unicode.IsSpace(r)
		})

		if j < 0 {
			return fields
		}
		k := strings.IndexFunc(s[i+j:], unicode.IsSpace)

		if k < 0 {
			return append(fields, s[i+j:])
		}
		fields = append(fields, s[i+j:i+j+k])
		i += j + k
	}
	return append(fields, strings.TrimSpace(s[i:]))
}

type job struct {
	command string
}

func newJob(command string) *job {
	return &job{command: command}
}

func (j *job) Run() {
	fmt.Printf("> /bin/sh -c \"%s\"\n", j.command)
	cmd := exec.Command("/bin/sh", "-c", j.command)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true, Pgid: 0}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout

	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
