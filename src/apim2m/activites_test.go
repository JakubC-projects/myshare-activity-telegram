package apim2m

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestGetActivities(t *testing.T) {

	res, err := GetActivities(context.Background(), ActivitiesQueryParams{
		Filter: &ActivitiesFilter{
			Created: &Filter[time.Time]{
				Gt: lo.ToPtr(time.Now().AddDate(0, 0, -1)),
			},
			Start: &Filter[time.Time]{
				Gt: lo.ToPtr(time.Now().AddDate(0, 0, -1)),
			},
		},
	})
	assert.NoError(t, err)
	for _, a := range res {
		fmt.Printf("%d. %s - %s, (%s - %s)\n", a.Id, a.Name,
			a.Created.Format(time.RFC3339),
			a.Start.Format(time.RFC3339),
			a.Finish.Format(time.RFC3339),
		)
	}
}
