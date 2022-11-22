package mail_test

import (
	"testing"

	"github.com/colynn/pontus/pkg/mail"

	"github.com/stretchr/testify/assert"
)

func TestSendMail(t *testing.T) {
	mailClint := mail.NewClient(
		"smtp.example.com",
		"dev-notify@example.com",
		"password@123",
		"mail-unit-test-单元测试",
		465,
		true,
	)
	flag := mailClint.SendMail([]string{"colynn.liu@example.cn"}, "unit-test", "hello world", "attachment")
	assert := assert.New(t)
	assert.Equal(flag, true)
}
