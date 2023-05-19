/*
	auth.go
	Purpose: Methods to authenticate users.

	@author Evan Chen

	MODIFICATION HISTORY
	   Date        Ver    Name     Description
	---------- ------- ----------- -------------------------------------------
	2023/03/02  v1.0.0 Evan Chen   Initial release
*/

package auth

import (
	"context"
	"strings"

	"app/core/errors"
	"app/core/util"

	"golang.org/x/exp/slog"
)

type Group int

const (
	CUSTOM Group = iota
	USER
	ADMIN
	SYSTEM

	SKIP Group = -1
)

var auto *UserInfo

func SetAutoLogin(u *UserInfo) {
	auto = u
}

var guards = map[string]Group{}

// Skip marks a method that should not be blocked by the middleware
// regardless of the authentication.
func Skip(methods ...string) {
	for _, method := range methods {
		guards[method] = SKIP
	}
}

// Guard sets a guard Group to @methods
func Guard(level Group, methods ...string) {
	for _, method := range methods {
		guards[method] = level
	}
}

// Authenticate checks if a request has sufficient credentials.
//
// This method should be called after [SetUser]
func Authenticate(ctx context.Context, method string) error {
	lvl := guards[method]
	if lvl == SKIP {
		return nil
	}

	usr, ok := GetUser(ctx)
	if !ok {
		return errors.ErrUnauthorized
	}

	if usr.Group < lvl {
		return errors.ErrForbidden
	}

	return nil
}

// GetUserFromToken gets the user from jwt token
func GetUserFromToken(tok string) (*UserInfo, bool) {
	// auto login set
	if auto != nil {
		return auto, true
	}
	tok = strings.TrimPrefix(tok, "Bearer ")
	if tok == "" {
		return nil, false
	}
	c, err := ParseToken(tok)
	if err != nil {
		slog.Error("ParseToken failed", util.ErrAtrr(err), slog.String("mod", "auth"))
		return nil, false
	}
	return &UserInfo{
		Username: c.User,
		Group:    Group(c.Group),
		Dept:     c.Dept,
	}, true
}

type user_key int

const ctx_user_key user_key = 0

// SetUser sets the user info to the context with the provided @tok string
func SetUser(ctx context.Context, tok string) context.Context {

	// auto login set
	if auto != nil {
		return context.WithValue(ctx, ctx_user_key, auto)
	}

	if usr, ok := GetUserFromToken(tok); !ok {
		return ctx
	} else {
		return context.WithValue(ctx, ctx_user_key, usr)
	}
}

// GetUser gets the userinfo fron the context.
func GetUser(ctx context.Context) (*UserInfo, bool) {
	u, ok := ctx.Value(ctx_user_key).(*UserInfo)
	return u, ok
}
