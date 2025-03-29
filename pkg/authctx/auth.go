package authctx

import (
	"context"

	"github.com/llkhacquan/cisab/pkg/models"
)

var authKey = new(int)

func Get(ctx context.Context) AuthMD {
	value := ctx.Value(authKey)
	if value == nil {
		return AuthMD{}
	}
	return value.(AuthMD)
}

func Set(ctx context.Context, value AuthMD) context.Context {
	return context.WithValue(ctx, authKey, value)
}

type AuthMD struct {
	User models.User
}
