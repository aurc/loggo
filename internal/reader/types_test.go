package reader

import (
	"testing"

	uuid2 "github.com/google/uuid"

	"github.com/stretchr/testify/assert"
)

func Test_SaveReadLoad_Params(t *testing.T) {
	sp := SavedParams{
		From:     "10m",
		Filter:   `resource.labels.namespace_name="somenamespace" resource.labels.container_name="mycontainer" NOT textPayload:"ignoring custom sampler for span"`,
		Project:  "someGCPProjectID",
		Template: "",
	}

	uuid := uuid2.New().String()
	err := Save(uuid, &sp)
	assert.NoError(t, err)

	l, err := List()
	assert.NoError(t, err)
	found := false
	for _, v := range l {
		if v.Name == uuid {
			found = true
			break
		}
	}
	assert.True(t, found)
	err = Remove(uuid)
	assert.NoError(t, err)
}
