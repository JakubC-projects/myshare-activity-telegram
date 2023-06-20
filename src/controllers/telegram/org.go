package telegram

import (
	"context"
	"fmt"
	"strconv"

	"github.com/JakubC-projects/myshare-activity-telegram/src/api"
	"github.com/JakubC-projects/myshare-activity-telegram/src/db"
	"github.com/JakubC-projects/myshare-activity-telegram/src/models"
	"github.com/JakubC-projects/myshare-activity-telegram/src/telegram"
	"github.com/samber/lo"
)

func startChangeOrg(ctx context.Context, u models.User) error {
	userOrgs, err := api.GetOrgs(ctx, u)
	if err != nil {
		return fmt.Errorf("cannot get orgs :%w", err)
	}
	_, err = telegram.SendChangeOrgMessage(u, userOrgs, telegram.Edit)
	return err
}

func changeOrg(ctx context.Context, u models.User, orgIdString string) error {
	orgId, err := strconv.Atoi(orgIdString)
	if err != nil {
		return fmt.Errorf("invalid org id %s: %w", orgIdString, err)
	}
	userOrgs, err := api.GetOrgs(ctx, u)
	if err != nil {
		return fmt.Errorf("cannot get orgs :%w", err)
	}

	selectedOrg, found := lo.Find(userOrgs, func(t models.Org) bool {
		return t.Id == orgId
	})
	if !found {
		return fmt.Errorf("selected invalid org: %d", orgId)
	}
	u.Org = &selectedOrg
	err = db.SaveUser(ctx, u)
	if err != nil {
		return fmt.Errorf("cannot save user: %w", err)

	}
	_, err = telegram.SendMenuMessage(u, telegram.Edit)
	if err != nil {
		return fmt.Errorf("cannot send menu message: %w", err)
	}

	return nil
}
