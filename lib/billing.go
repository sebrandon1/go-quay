/*
Package lib provides Quay.io API client functionality.

This file covers BILLING AND SUBSCRIPTION endpoints:

Organization Billing:
  - GET /api/v1/organization/{orgname}/plan           - GetOrganizationBilling()
  - GET /api/v1/organization/{orgname}/plan           - GetOrganizationSubscription()
  - GET /api/v1/organization/{orgname}/invoices       - GetOrganizationInvoices()

User Billing:
  - GET /api/v1/user/plan                            - GetUserBilling()
  - GET /api/v1/user/plan                            - GetUserSubscription()

Plans:
  - GET /api/v1/plans                                - GetAvailablePlans()

Note: User invoice endpoint is not available in Quay API.
*/
package lib

import (
	"fmt"
)

// GetOrganizationBilling returns billing information for an organization
func (c *Client) GetOrganizationBilling(orgname string) (*BillingInfo, error) {
	// Get new request
	req, err := newRequest("GET", c.buildURL("/organization/%s/plan", orgname), nil)
	if err != nil {
		return nil, err
	}

	var billing BillingInfo
	if err := c.get(req, &billing); err != nil {
		return nil, err
	}

	return &billing, nil
}

// GetUserBilling returns billing information for the current user
func (c *Client) GetUserBilling() (*BillingInfo, error) {
	// Get new request
	req, err := newRequest("GET", c.BaseURL+"/user/plan", nil)
	if err != nil {
		return nil, err
	}

	var billing BillingInfo
	if err := c.get(req, &billing); err != nil {
		return nil, err
	}

	return &billing, nil
}

// GetOrganizationSubscription returns subscription details for an organization
func (c *Client) GetOrganizationSubscription(orgname string) (*Subscription, error) {
	// Get new request
	req, err := newRequest("GET", c.buildURL("/organization/%s/plan", orgname), nil)
	if err != nil {
		return nil, err
	}

	var subscription Subscription
	if err := c.get(req, &subscription); err != nil {
		return nil, err
	}

	return &subscription, nil
}

// GetUserSubscription returns subscription details for the current user
func (c *Client) GetUserSubscription() (*Subscription, error) {
	// Get new request
	req, err := newRequest("GET", c.BaseURL+"/user/plan", nil)
	if err != nil {
		return nil, err
	}

	var subscription Subscription
	if err := c.get(req, &subscription); err != nil {
		return nil, err
	}

	return &subscription, nil
}

// GetOrganizationInvoices returns invoices for an organization
func (c *Client) GetOrganizationInvoices(orgname string) ([]Invoice, error) {
	// Get new request
	req, err := newRequest("GET", c.buildURL("/organization/%s/invoices", orgname), nil)
	if err != nil {
		return nil, err
	}

	var invoicesResp struct {
		Invoices []Invoice `json:"invoices"`
	}
	if err := c.get(req, &invoicesResp); err != nil {
		return nil, err
	}

	return invoicesResp.Invoices, nil
}

// GetUserInvoices - NOT AVAILABLE in Quay API (returns 404)
func (c *Client) GetUserInvoices() ([]Invoice, error) {
	return nil, fmt.Errorf("user invoices endpoint not available in Quay API")
}

// GetAvailablePlans returns available subscription plans
func (c *Client) GetAvailablePlans() ([]Subscription, error) {
	// Get new request
	req, err := newRequest("GET", c.BaseURL+"/plans", nil)
	if err != nil {
		return nil, err
	}

	var plansResp struct {
		Plans []Subscription `json:"plans"`
	}
	if err := c.get(req, &plansResp); err != nil {
		return nil, err
	}

	return plansResp.Plans, nil
}
