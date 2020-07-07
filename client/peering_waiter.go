package client

import (
	"context"
	"time"
)

type WaitForPeeringStateRequest struct {
	OrganizationID string
	ProjectID      string
	PeeringID      string
	State          string
}

func (c *Client) PeeringWaitForState(ctx context.Context, req *WaitForPeeringStateRequest) (*Peering, error) {
	getRequest := &GetPeeringRequest{
		OrganizationID: req.OrganizationID,
		ProjectID:      req.ProjectID,
		PeeringID:      req.PeeringID,
	}

	for {
		resp, err := c.PeeringGet(ctx, getRequest)
		if err != nil {
			return nil, err
		}

		if resp.Peering.Status != req.State {
			time.Sleep(5 * time.Second)
			continue
		}

		return &resp.Peering, nil
	}
}
