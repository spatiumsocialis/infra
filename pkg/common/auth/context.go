package auth

import (
	"context"
	"errors"
)

// GetToken retrieves the access token from a request context. It returns an error if the token isn't found
func GetToken(ctx context.Context) (*Token, error) {
	token := ctx.Value(contextKeyAuthToken)
	if token == nil {
		err := errors.New("Error: context doesn't contain a token")
		return nil, err
	}
	t := token.(*Token)
	return t, nil
}

// WithToken adds an access token to a request context
func WithToken(ctx context.Context, token *Token) context.Context {
	return context.WithValue(ctx, contextKeyAuthToken, token)
}
