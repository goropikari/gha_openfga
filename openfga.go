package openfga

import (
	"context"
	"encoding/json"
	"os"

	fga "github.com/openfga/go-sdk"
	"github.com/openfga/go-sdk/client"
)

type Client struct {
	fgaClient *client.OpenFgaClient
}

func NewClient() (*Client, error) {
	url := os.Getenv("FGA_API_URL")
	if url == "" {
		url = "http://localhost:8080" // default to local server if not set
	}

	fgaClient, err := client.NewSdkClient(&client.ClientConfiguration{
		ApiUrl: url,
	})
	if err != nil {
		return nil, err
	}

	return &Client{
		fgaClient: fgaClient,
	}, nil
}

func (c *Client) CreateStore(ctx context.Context) error {
	resp, err := c.fgaClient.CreateStore(ctx).Body(client.ClientCreateStoreRequest{Name: "FGA Demo"}).Execute()
	if err != nil {
		return err
	}

	c.fgaClient.SetStoreId(resp.GetId())
	return nil
}

func (c *Client) WriteAuthorizationModel(ctx context.Context) error {
	// https://openfga.dev/docs/getting-started/configure-model
	writeAuthorizationModelRequestString := "{\"schema_version\":\"1.1\",\"type_definitions\":[{\"type\":\"user\"},{\"type\":\"document\",\"relations\":{\"reader\":{\"this\":{}},\"writer\":{\"this\":{}},\"owner\":{\"this\":{}}},\"metadata\":{\"relations\":{\"reader\":{\"directly_related_user_types\":[{\"type\":\"user\"}]},\"writer\":{\"directly_related_user_types\":[{\"type\":\"user\"}]},\"owner\":{\"directly_related_user_types\":[{\"type\":\"user\"}]}}}}]}"
	var body fga.WriteAuthorizationModelRequest
	if err := json.Unmarshal([]byte(writeAuthorizationModelRequestString), &body); err != nil {
		return err
	}

	if _, err := c.fgaClient.WriteAuthorizationModel(ctx).
		Body(body).
		Execute(); err != nil {
		return err
	}

	return nil
}

func (c *Client) WriteTuple(ctx context.Context, tuple fga.TupleKey) error {
	body := client.ClientWriteRequest{
		Writes: []client.ClientTupleKey{
			{
				User:     tuple.GetUser(),
				Relation: tuple.GetRelation(),
				Object:   tuple.GetObject(),
			},
		},
	}
	if _, err := c.fgaClient.Write(ctx).
		Body(body).
		Execute(); err != nil {
		return err
	}
	return nil
}

func (c *Client) Check(ctx context.Context, tuple fga.TupleKey) (bool, error) {
	body := client.ClientCheckRequest{
		User:     tuple.GetUser(),
		Relation: tuple.GetRelation(),
		Object:   tuple.GetObject(),
	}
	resp, err := c.fgaClient.Check(ctx).
		Body(body).
		Execute()
	if err != nil {
		return false, err
	}

	return resp.GetAllowed(), nil
}
