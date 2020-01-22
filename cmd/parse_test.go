package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func testSuccessToParseWatchOption(t *testing.T) {
	options, err := parseWatchOption([]string {"-n", "4", "-d", "ls", "-al"})
	assert.Nil(t, err)
	assert.Equal(
		t,
		options,
		watchOptions{ 4, true, []string {"ls", "-al"}})
}
