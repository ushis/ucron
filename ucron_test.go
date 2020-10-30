package main

import (
	"bytes"
	"fmt"
	"testing"
)

func assertEquality(actual, expected interface{}) error {
	if actual == expected {
		return nil
	}
	return fmt.Errorf("expected %q to equal %q", actual, expected)
}

const crontab = `
* * * * * echo test
`

func TestReadCrontab(t *testing.T) {
	entries, err := readCrontab(bytes.NewBufferString(crontab))

	if err != nil {
		t.Fatal(err)
	}
	if err := assertEquality(len(entries), 1); err != nil {
		t.Fatal(err)
	}
	if err := assertEquality(entries[0][0], "* * * * *"); err != nil {
		t.Fatal(err)
	}
	if err := assertEquality(entries[0][1], "echo test"); err != nil {
		t.Fatal(err)
	}
}

const crontabWithComment = `
# some comment
* * * * * echo test
`

func TestReadCrontabWithComment(t *testing.T) {
	entries, err := readCrontab(bytes.NewBufferString(crontabWithComment))

	if err != nil {
		t.Fatal(err)
	}
	if err := assertEquality(len(entries), 1); err != nil {
		t.Fatal(err)
	}
	if err := assertEquality(entries[0][0], "* * * * *"); err != nil {
		t.Fatal(err)
	}
	if err := assertEquality(entries[0][1], "echo test"); err != nil {
		t.Fatal(err)
	}
}

const crontabWithBlankLine = `
* * * * * echo hello

5 4 3 2 1 echo world
`

func TestReadCrontabWithBlankLine(t *testing.T) {
	entries, err := readCrontab(bytes.NewBufferString(crontabWithBlankLine))

	if err != nil {
		t.Fatal(err)
	}
	if err := assertEquality(len(entries), 2); err != nil {
		t.Fatal(err)
	}
	if err := assertEquality(entries[0][0], "* * * * *"); err != nil {
		t.Fatal(err)
	}
	if err := assertEquality(entries[0][1], "echo hello"); err != nil {
		t.Fatal(err)
	}
	if err := assertEquality(entries[1][0], "5 4 3 2 1"); err != nil {
		t.Fatal(err)
	}
	if err := assertEquality(entries[1][1], "echo world"); err != nil {
		t.Fatal(err)
	}
}

const crontabWithoutNewlineTerminator = `
* * * * * echo test`

func TestReadCrontabWithoutNewlineTerminator(t *testing.T) {
	entries, err := readCrontab(bytes.NewBufferString(crontabWithoutNewlineTerminator))

	if err != nil {
		t.Fatal(err)
	}
	if err := assertEquality(len(entries), 1); err != nil {
		t.Fatal(err)
	}
	if err := assertEquality(entries[0][0], "* * * * *"); err != nil {
		t.Fatal(err)
	}
	if err := assertEquality(entries[0][1], "echo test"); err != nil {
		t.Fatal(err)
	}
}

const crontabWithRedundantSpaces = `
   *     * 	*   *  *       echo test
`

func TestReadCrontabWithRedundantSpaces(t *testing.T) {
	entries, err := readCrontab(bytes.NewBufferString(crontabWithRedundantSpaces))

	if err != nil {
		t.Fatal(err)
	}
	if err := assertEquality(len(entries), 1); err != nil {
		t.Fatal(err)
	}
	if err := assertEquality(entries[0][0], "* * * * *"); err != nil {
		t.Fatal(err)
	}
	if err := assertEquality(entries[0][1], "echo test"); err != nil {
		t.Fatal(err)
	}
}

const crontabWithMissingFields = `
   *     * 	*   *
`

func TestReadCrontabWithMissingFields(t *testing.T) {
	if _, err := readCrontab(bytes.NewBufferString(crontabWithMissingFields)); err == nil {
		t.Fatal("expected to receive an error")
	}
}
