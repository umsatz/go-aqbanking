package aqbanking

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGwenDate(t *testing.T) {
	assert := assert.New(t)

	goDate := time.Date(2009, time.November, 10, 0, 0, 0, 0, time.Local)
	gwenDate := newGwenDate(goDate)

	assert.Equal("20091110", gwenDate.String())
	assert.Equal(goDate, gwenDate.toTime())
}
