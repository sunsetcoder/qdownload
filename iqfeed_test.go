package main

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testRequestId                    = "999"
	testProtocolMessage              = "S,CURRENT PROTOCOL,5.1"
	testEndMessage                   = "999,!ENDMSG!,"
	testNoDataMessage                = "999,E,!NO_DATA!,,"
	testValidIqfeedEodBar            = "999,2019-02-21,24.0600,23.8038,23.8700,24.0000,29183,0,"
	testTooFewColumnsIqfeedEodBar    = "999,2019-02-21,24.0600,23.8038,23.8700,24.0000,29183"
	testValidIqfeedMinuteBar         = "999,2019-02-26 12:22:00,23.8000,23.8000,23.8000,23.8000,13578,100,0,"
	testTooFewColumnsIqfeedMinuteBar = "999,2019-02-26 12:22:00,23.8000,23.8000,23.8000,23.8000"
	testValidIqfeedTick              = "999,2019-02-25 11:30:06.691,23.8800,12,6714,23.8700,23.9700,6,O,25,3D87,"
	testTooFewColumnsIqfeedTick      = "999,2019-02-25 11:30:06.691,23.8800,12,6714,23.8700,23.9700,6,O,25"
	testIncorrectRequestIdIqfeedTick = "111,2019-02-25 11:30:06.691,23.8800,12,6714,23.8700,23.9700,6,O,25,3D87,"
)

func TestMapRow(t *testing.T) {
	t.Run("valid tick to csv", func(t *testing.T) {
		columns := strings.Split(testValidIqfeedTick, ",")

		mappedRow, err := mapRow(columns, testRequestId, mapTick, false)

		assert.Equal(t, "2019-02-25 11:30:06.691,23.8800,12,6714,23.8700,23.9700,6,O,25,3D87", mappedRow)
		assert.Nil(t, err)
	})

	t.Run("valid tick to tsv", func(t *testing.T) {
		columns := strings.Split(testValidIqfeedTick, ",")

		mappedRow, err := mapRow(columns, testRequestId, mapTick, true)

		assert.Equal(t, "2019-02-25 11:30:06.691\t23.8800\t12\t6714\t23.8700\t23.9700\t6\tO\t25\t3D87", mappedRow)
		assert.Nil(t, err)
	})

	t.Run("no columns", func(t *testing.T) {
		var columns []string

		mappedRow, err := mapRow(columns, testRequestId, mapTick, false)

		assert.Equal(t, "", mappedRow)
		assert.Errorf(t, err, "empty row")
	})

	t.Run("protocol message", func(t *testing.T) {
		columns := strings.Split(testProtocolMessage, ",")

		mappedRow, err := mapRow(columns, testRequestId, mapTick, false)

		assert.Equal(t, "", mappedRow)
		assert.Nil(t, err)
	})

	t.Run("incorrect request id", func(t *testing.T) {
		columns := strings.Split(testIncorrectRequestIdIqfeedTick, ",")

		mappedRow, err := mapRow(columns, testRequestId, mapTick, false)

		assert.Equal(t, "", mappedRow)
		assert.Errorf(t, err, "incorrect request id")
	})

	t.Run("end message", func(t *testing.T) {
		columns := strings.Split(testEndMessage, ",")

		mappedRow, err := mapRow(columns, testRequestId, mapTick, false)

		assert.Equal(t, "", mappedRow)
		assert.Equal(t, io.EOF, err)
	})

	t.Run("no data", func(t *testing.T) {
		columns := strings.Split(testNoDataMessage, ",")

		mappedRow, err := mapRow(columns, testRequestId, mapTick, false)

		assert.Equal(t, "", mappedRow)
		assert.Errorf(t, err, "iqfeed error: !NO DATA!")
	})

	t.Run("too few columns", func(t *testing.T) {
		columns := strings.Split(testTooFewColumnsIqfeedTick, ",")

		mappedRow, err := mapRow(columns, testRequestId, mapTick, false)

		assert.Equal(t, "", mappedRow)
		assert.Nil(t, err)
	})
}

func TestEodBarMapper(t *testing.T) {
	t.Run("valid eod bar", func(t *testing.T) {
		columns := strings.Split(testValidIqfeedEodBar, ",")

		mappedRow, err := mapEodBar(columns)

		assert.Equal(t, "2019-02-21,23.8700,24.0600,23.8038,24.0000,29183,0", mappedRow)
		assert.Nil(t, err)
	})

	t.Run("too few columns", func(t *testing.T) {
		columns := strings.Split(testTooFewColumnsIqfeedEodBar, ",")

		mappedRow, err := mapEodBar(columns)

		assert.Equal(t, "", mappedRow)
		assert.Errorf(t, err, "too few columns")
	})

	t.Run("no columns", func(t *testing.T) {
		var columns []string

		mappedRow, err := mapEodBar(columns)

		assert.Equal(t, "", mappedRow)
		assert.Errorf(t, err, "too few columns")
	})
}

func TestMinuteBarMapper(t *testing.T) {
	t.Run("valid minute bar", func(t *testing.T) {
		columns := strings.Split(testValidIqfeedMinuteBar, ",")

		mappedRow, err := mapMinuteBar(columns)

		assert.Equal(t, "2019-02-26 12:21:00,23.8000,23.8000,23.8000,23.8000,100", mappedRow)
		assert.Nil(t, err)
	})

	t.Run("too few columns", func(t *testing.T) {
		columns := strings.Split(testTooFewColumnsIqfeedMinuteBar, ",")

		mappedRow, err := mapMinuteBar(columns)

		assert.Equal(t, "", mappedRow)
		assert.Errorf(t, err, "too few columns")
	})

	t.Run("no columns", func(t *testing.T) {
		var columns []string

		mappedRow, err := mapMinuteBar(columns)

		assert.Equal(t, "", mappedRow)
		assert.Errorf(t, err, "too few columns")
	})
}

func TestTickMapper(t *testing.T) {
	t.Run("valid tick", func(t *testing.T) {
		columns := strings.Split(testValidIqfeedTick, ",")

		mappedRow, err := mapTick(columns)

		assert.Equal(t, "2019-02-25 11:30:06.691,23.8800,12,6714,23.8700,23.9700,6,O,25,3D87", mappedRow)
		assert.Nil(t, err)
	})

	t.Run("too few columns", func(t *testing.T) {
		columns := strings.Split(testTooFewColumnsIqfeedTick, ",")

		mappedRow, err := mapTick(columns)

		assert.Equal(t, "", mappedRow)
		assert.Errorf(t, err, "too few columns")
	})

	t.Run("no columns", func(t *testing.T) {
		var columns []string

		mappedRow, err := mapTick(columns)

		assert.Equal(t, "", mappedRow)
		assert.Errorf(t, err, "too few columns")
	})
}
