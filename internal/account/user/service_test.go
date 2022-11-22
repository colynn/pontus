package user

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseUploadFile(t *testing.T) {
	value, err := parseUploadFile("/Users/colynn/Documents/Code/pontus/backend/assets/displayMapLdapAccount.txt")
	assert.Equal(t, nil, err)
	fmt.Printf("value length : %v", value[1])
}
