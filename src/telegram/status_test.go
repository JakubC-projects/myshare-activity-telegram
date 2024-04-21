package telegram

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetStatusForToday(t *testing.T) {
	d, _ := time.Parse(time.RFC3339, "2024-04-21T12:00:00Z")
	expectedP := 47.7
	p := getPeacefulRoadPercentForDate(d)
	assert.InDelta(t, expectedP, p, 0.1)
}
