package dbctx

import (
	"context"

	"gorm.io/gorm"
)

var dbKey = new(int)

type DBGetter func(ctx context.Context) *gorm.DB

var Get = func(ctx context.Context) *gorm.DB {
	value := ctx.Value(dbKey)
	return value.(*gorm.DB).WithContext(ctx)
}

func Set(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, dbKey, db)
}
