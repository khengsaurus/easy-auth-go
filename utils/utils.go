package utils

import (
	"context"
	"fmt"

	"github.com/khengsaurus/easy-auth-middlewares/consts"
	"github.com/khengsaurus/easy-auth-middlewares/types"
)

func GetRequestUser(ctx context.Context) (*types.User, error) {
	user, ok := ctx.Value(consts.ContextKeyUser).(*types.User)
	if !ok || user.Id == "" || user.Username == "" {
		return nil, fmt.Errorf(consts.ErrorUserInfoMissing)
	}
	return user, nil
}
