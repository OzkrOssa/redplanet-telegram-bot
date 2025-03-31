package port

import "github.com/OzkrOssa/redplanet-telegram-bot/internal/core/domain"

type MikrotikService interface {
	GetIdentity() (*string, error)
	GetResources() (*domain.Resource, error)
	GetTraffic(mikrotikInterface string) (*domain.Traffic, error)
	ChangeMangleRuleStatus(status string) error
	ChangeStatusStaticRoutes(event string) error
}
