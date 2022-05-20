package auth

import "context"

func ExtractToken(ctx context.Context) string {
	s, ok := ctx.Value("authorization_header").(string)
	if !ok {
		return ""
	}

	return s
}
