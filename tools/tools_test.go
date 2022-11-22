package tools_test

import (
	"fmt"
	"testing"

	"github.com/colynn/pontus/tools"

	"github.com/stretchr/testify/assert"
)

func TestParseStrToDateTime(t *testing.T) {
	value, err := tools.ParseStrToDateTime("2019-01-02 15:22:05")
	assert.Equal(t, nil, err)
	zoneName, _ := value.Zone()
	assert.Equal(t, 2019, value.Year())
	assert.Equal(t, 2, value.Day())
	assert.Equal(t, "CST", zoneName)
	fmt.Printf("value: %v", value)
}

func TestParseStrToDate(t *testing.T) {
	value, err := tools.ParseStrToDate("2020-10-22")
	assert.Equal(t, nil, err)
	zoneName, _ := value.Zone()
	assert.Equal(t, 2020, value.Year())
	assert.Equal(t, 22, value.Day())
	assert.Equal(t, "CST", zoneName)
}

func TestGetCurrentTime(t *testing.T) {
	value := tools.GetCurrentTime("")
	zoneName, _ := value.Zone()
	assert.Equal(t, "CST", zoneName)
}
