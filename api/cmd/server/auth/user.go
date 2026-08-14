package auth

import "context"

type User struct {
	ID       string `json:"-"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

type contextKey string

const userKey contextKey = "current-user"

func WithUser(ctx context.Context, user User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func CurrentUser(ctx context.Context) (User, bool) {
	user, ok := ctx.Value(userKey).(User)
	return user, ok
}
