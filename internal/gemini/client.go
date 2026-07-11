package gemini

import "context"

type Client interface {
	GenerateContent(
		ctx context.Context,
		prompt string,
	) (string, error)
}
