package checkactivites

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	err := handleActivitiesCheck(context.Background())
	assert.NoError(t, err)
}
