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

func (ms *MikrotikService) ChangeMangleRuleStatus(status string) error {
	cmd := fmt.Sprintf("/ip/firewall/mangle/%s", status)
	_, err := ms.client.Run(cmd, "=numbers=0")
	if err != nil {
		return err
	}
	return nil
}

func (ms *MikrotikService) ChangeStaticRoutesStatus(event string) error {
	routesList, err := ms.client.Run("/ip/route/print")
	if err != nil {
		return err
	}
	switch event {
	case "normal":
		for _, route := range routesList.Re {
			if route.Map["disabled"] == "true" {
				_, err := ms.client.Run("/ip/route/set", "=numbers="+route.Map[".id"], "=disabled=false")
				if err != nil {
					return err
				}
			}
		}
	case "azt_down":
		aztGateways := map[string]bool{
			"177.93.60.209":   true,
			"186.179.103.241": true,
			"200.69.92.1":     true,
		}

		for _, route := range routesList.Re {
			if aztGateways[route.Map["gateway"]] {
				_, err := ms.client.Run("/ip/route/set", "=numbers="+route.Map[".id"], "=disabled=true")
				if err != nil {
					return err
				}
			}
		}
	case "ufinet_down":
		for _, route := range routesList.Re {
			if route.Map["gateway"] == "190.60.55.193" {
				_, err := ms.client.Run("/ip/route/set", "=numbers="+route.Map[".id"], "=disabled=true")
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (ms *MikrotikService) ChangeSimpleQueuesStatus(event string) error {
	queuesList, err := ms.client.Run("/queue/simple/print")
	if err != nil {
		return err
	}
	switch event {
	case "normal":
		for _, queue := range queuesList.Re {
			if queue.Map["comment"] == "ONDOWN" {
				_, err := ms.client.Run("/queue/simple/set", "=numbers="+queue.Map[".id"], "=disabled=true")
				if err != nil {
					panic(err)
				}
			}
		}
	case "azt_down":
		for _, queue := range queuesList.Re {
			if queue.Map["comment"] == "ONDOWN" {
				_, err := ms.client.Run("/queue/simple/set", "=numbers="+queue.Map[".id"], "=disabled=false")
				if err != nil {
					panic(err)
				}
			}
		}
	case "ufinet_down":
		for _, queue := range queuesList.Re {
			if queue.Map["comment"] == "ONDOWN" {
				_, err := ms.client.Run("/queue/simple/set", "=numbers="+queue.Map[".id"], "=disabled=false")
				if err != nil {
					panic(err)
				}
			}
		}
	}
	return nil
}

func (ms *MikrotikService) ChangeNametoAddressList(event string) error {
	addressList, err := ms.client.Run("/ip/firewall/address-list/print")
	if err != nil {
		return err
	}

	routersName := map[string]bool{
		"PuebloViejo": true,
		"Cabuyal":     true,
		"Clavijo":     true,
		"Blandon":     true,
		"Tabuyo":      true,
		"Calera":      true,
	}
	switch event {
	case "normal":
		for _, list := range addressList.Re {
			if routersName[list.Map["comment"]] {
				_, err := ms.client.Run("/ip/firewall/address-list/set", "=numbers="+list.Map[".id"], "=list=2")
				if err != nil {
					panic(err)
				}
			}
		}
	case "azt_down":
		for _, queue := range addressList.Re {
			switch queue.Map["comment"] {
			case "moron":
				_, err := ms.client.Run("/ip/firewall/address-list/set", "=numbers="+queue.Map[".id"], "=list=1")
				if err != nil {
					panic(err)
				}
			case "PuebloViejo", "Clavijo", "Calera", "Cabuyal":
				_, err := ms.client.Run("/ip/firewall/address-list/set", "=numbers="+queue.Map[".id"], "=list=4")
				if err != nil {
					panic(err)
				}
			case "Blandon", "Tabuyo":
				_, err := ms.client.Run("/ip/firewall/address-list/set", "=numbers="+queue.Map[".id"], "=list=5")
				if err != nil {
					panic(err)
				}

			}

		}
	}
	return nil
}
