/*
Package lib provides Quay.io API client functionality.

This file covers ORGANIZATION MARKETPLACE endpoints:

Marketplace Management:
  - GET    /api/v1/organization/{orgname}/marketplace                          - GetOrganizationMarketplace()
  - POST   /api/v1/organization/{orgname}/marketplace                          - CreateOrganizationMarketplaceSubscription()
  - POST   /api/v1/organization/{orgname}/marketplace/batchremove              - BatchRemoveOrganizationMarketplaceSubscriptions()
  - DELETE /api/v1/organization/{orgname}/marketplace/{subscription_id}        - DeleteOrganizationMarketplaceSubscription()
*/
package lib

import (
	"fmt"
)

// GetOrganizationMarketplace gets marketplace information for an organization
func (c *Client) GetOrganizationMarketplace(orgname string) (*MarketplaceInfo, error) {
	req, err := newRequest("GET", c.buildURL("/organization/%s/marketplace", orgname), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get organization marketplace request: %w", err)
	}

	var marketplace MarketplaceInfo
	if err := c.get(req, &marketplace); err != nil {
		return nil, fmt.Errorf("failed to get organization marketplace: %w", err)
	}

	return &marketplace, nil
}

// CreateOrganizationMarketplaceSubscription creates a marketplace subscription
func (c *Client) CreateOrganizationMarketplaceSubscription(orgname string, subscription *MarketplaceSubscriptionRequest) error {
	req, err := newRequestWithBody("POST", c.buildURL("/organization/%s/marketplace", orgname), subscription)
	if err != nil {
		return fmt.Errorf("failed to create marketplace subscription request: %w", err)
	}

	if err := c.post(req, nil); err != nil {
		return fmt.Errorf("failed to create marketplace subscription: %w", err)
	}

	return nil
}

// BatchRemoveOrganizationMarketplaceSubscriptions removes multiple marketplace subscriptions
func (c *Client) BatchRemoveOrganizationMarketplaceSubscriptions(orgname string, subscriptionIDs []string) error {
	body := struct {
		SubscriptionIDs []string `json:"subscription_ids"`
	}{
		SubscriptionIDs: subscriptionIDs,
	}
	req, err := newRequestWithBody("POST", c.buildURL("/organization/%s/marketplace/batchremove", orgname), body)
	if err != nil {
		return fmt.Errorf("failed to create batch remove marketplace subscriptions request: %w", err)
	}

	if err := c.post(req, nil); err != nil {
		return fmt.Errorf("failed to batch remove marketplace subscriptions: %w", err)
	}

	return nil
}

// DeleteOrganizationMarketplaceSubscription removes a specific marketplace subscription
func (c *Client) DeleteOrganizationMarketplaceSubscription(orgname, subscriptionID string) error {
	req, err := newRequest("DELETE", c.buildURL("/organization/%s/marketplace/%s", orgname, subscriptionID), nil)
	if err != nil {
		return fmt.Errorf("failed to create delete marketplace subscription request: %w", err)
	}

	if err := c.delete(req); err != nil {
		return fmt.Errorf("failed to delete marketplace subscription: %w", err)
	}

	return nil
}
