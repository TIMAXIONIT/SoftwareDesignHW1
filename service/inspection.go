package service

import "vyshkat/domain"

type ServiceCenter struct{}

func NewServiceCenter() *ServiceCenter {
	return &ServiceCenter{}
}

func (s *ServiceCenter) Inspect(
	item domain.InventoryItem,
) error {
	return nil
}
