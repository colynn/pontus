package cmdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetNeedDeleteTags(t *testing.T) {
	originTags := []AssetTag{
		{
			TagKey:   "key01",
			TagValue: "value01",
		},
		{
			TagKey:   "key02",
			TagValue: "value02",
		},
	}

	currentTags := []AssetTag{
		{
			TagKey:   "key01",
			TagValue: "value01",
		},
		{
			TagKey:   "key03",
			TagValue: "value03",
		},
	}
	assert := assert.New(t)
	tags := getNeedDeleteTags(originTags, currentTags, "")
	assert.Equal(len(tags), 1)
	assert.Equal(tags[0].TagValue, "value02")
}
