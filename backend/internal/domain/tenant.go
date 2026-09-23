package domain

import "time"

type Tenant struct {
	ID                    int64      `json:"id"`
	Name                  string     `json:"name"`
	Slug                  string     `json:"slug"`
	Domain                string     `json:"domain"`
	LogoURL               string     `json:"logo_url"`
	SubscriptionPlan      string     `json:"subscription_plan"`
	SubscriptionExpiresAt *time.Time `json:"subscription_expires_at,omitempty"`
	MaxUsers              int        `json:"max_users"`
	MaxProducts           int        `json:"max_products"`
	MaxBranches           int        `json:"max_branches"`
	Settings              string     `json:"settings"`
	Status                string     `json:"status"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
}
