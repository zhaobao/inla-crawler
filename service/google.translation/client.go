package google_translation

import (
	"cloud.google.com/go/translate"
	"context"
	"fmt"
	"golang.org/x/text/language"
	"google.golang.org/api/option"
)

type client struct {
	*translate.Client
	ctx context.Context
}

func New(credential string) (*client, error) {
	ctx := context.Background()
	c, err := translate.NewClient(ctx, option.WithCredentialsFile(credential))
	if err != nil {
		return nil, err
	}
	return &client{Client: c, ctx: ctx}, nil
}

func (c *client) TranslateText(target language.Tag, text []string) ([]translate.Translation, error) {
	defer func() {
		_ = c.Close()
	}()

	resp, err := c.Translate(c.ctx, text, target, nil)
	if err != nil {
		return nil, fmt.Errorf("translate: %v", err)
	}
	if len(resp) == 0 {
		return nil, fmt.Errorf("translate returned empty response to text: %s", text)
	}
	return resp, nil
}
