package service

import (
	"fmt"
	"strings"

	"vyshkat/domain"
)

type ServiceCenter struct{}

func NewServiceCenter() *ServiceCenter {
	return &ServiceCenter{}
}

func (s *ServiceCenter) Inspect(item domain.InventoryItem) error {
	return nil
}

type DemoServiceCenter struct{}

func NewDemoServiceCenter() *DemoServiceCenter {
	return &DemoServiceCenter{}
}

func (s *DemoServiceCenter) Inspect(item domain.InventoryItem) error {
	if strings.HasPrefix(item.InventoryNumber(), "FAIL-") {
		return fmt.Errorf(
			"демонстрационный отказ: объект %q не прошёл техосмотр",
			item.InventoryNumber(),
		)
	}

	return nil
}
