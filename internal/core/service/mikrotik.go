package service

import (
	"fmt"

	"github.com/OzkrOssa/redplanet-telegram-bot/internal/adapter/config"
	"github.com/OzkrOssa/redplanet-telegram-bot/internal/core/domain"
	"github.com/go-routeros/routeros/v3"
)

type MikrotikService struct {
	client *routeros.Client
}

func NewMikrotikService(address string, config config.RouterOsApi) (*MikrotikService, error) {
	formatAddress := fmt.Sprintf("%s:%s", address, config.Port)

	client, err := routeros.Dial(formatAddress, config.User, config.Password)
	if err != nil {
		return nil, err
	}

	return &MikrotikService{
		client: client,
	}, nil
}

func (ms *MikrotikService) GetIndentity() (*string, error) {
	reply, err := ms.client.Run("/system/identity/print")
	if err != nil {
		return nil, err
	}

	identity, ok := reply.Re[0].Map["name"]
	if !ok {
		return nil, fmt.Errorf("error to get name: key 'name' not found in map")
	}

	return &identity, nil
}

func (ms *MikrotikService) GetResources() (*domain.Resource, error) {

	reply, err := ms.client.Run("/system/resource/print")
	if err != nil {
		return nil, err
	}

	identity, err := ms.GetIndentity()
	if err != nil {
		return nil, err
	}

	cpu, ok := reply.Re[0].Map["cpu-load"]
	if !ok {
		return nil, fmt.Errorf("error to get name: key 'cpu-load' not found in map")
	}

	uptime, ok := reply.Re[0].Map["uptime"]
	if !ok {
		return nil, fmt.Errorf("error to get name: key 'uptime' not found in map")
	}

	resource := &domain.Resource{
		Source: identity,
		Cpu:    cpu,
		Uptime: uptime,
	}

	return resource, nil
}

func (ms *MikrotikService) GetTraffic(mikrotikInterface string) (*domain.Traffic, error) {

	reply, err := ms.client.Run("/interface/monitor-traffic", "=interface="+mikrotikInterface, "=once")

	if err != nil {
		return nil, err
	}

	identity, err := ms.GetIndentity()
	if err != nil {
		return nil, err
	}

	rx, ok := reply.Re[0].Map["rx-bits-per-second"]
	if !ok {
		return nil, fmt.Errorf("error to get name: key 'rx-bits-per-second' not found in map")
	}

	tx, ok := reply.Re[0].Map["tx-bits-per-second"]
	if !ok {
		return nil, fmt.Errorf("error to get name: key 'tx-bits-per-second' not found in map")
	}

	traffic := &domain.Traffic{
		Source: identity,
		Rx:     rx,
		Tx:     tx,
	}

	return traffic, nil
}
