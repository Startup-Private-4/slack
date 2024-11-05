package slack

import (
	"context"
	"encoding/json"
)

type CreateChannelCanvasParams struct {
	ChannelID       string           `json:"channel_id"`
	DocumentContent *DocumentContent `json:"document_content"`
}

type DocumentContent struct {
	Type     string `json:"type"`
	Markdown string `json:"markdown"`
}

func (api *Client) CreateChannelCanvasContext(ctx context.Context, params CreateChannelCanvasParams) (*string, error) {
	encoded, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	resp := &struct {
		SlackResponse
		CanvasID *string `json:"canvas_id"`
	}{}

	endpoint := api.endpoint + "conversations.canvases.create"
	err = postJSON(ctx, api.httpclient, endpoint, api.token, encoded, resp, api)
	if err != nil {
		return nil, err
	}
	return resp.CanvasID, resp.Err()
}
