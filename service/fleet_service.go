package service

import (
	"fmt"

	"vyshkat/domain"
)

type FleetService struct {
	inventory  domain.Inventory
	inspection domain.InspectionService
}

func NewFleetService(
	inventory domain.Inventory,
	inspection domain.InspectionService,
) *FleetService {
	return &FleetService{
		inventory:  inventory,
		inspection: inspection,
	}
}

func (s *FleetService) AddItem(
	item domain.InventoryItem,
) error {
	if item == nil {
		return fmt.Errorf("нельзя добавить пустой объект")
	}

	number := item.InventoryNumber()
	if number == "" {
		return fmt.Errorf("инвентарный номер не может быть пустым")
	}

	if transport, ok := item.(domain.Transport); ok {
		level := transport.BeginnerLevel()
		if level < 1 || level > 10 {
			return fmt.Errorf(
				"уровень простоты должен быть от 1 до 10, получено %d",
				level,
			)
		}
	}

	if consumer, ok := item.(domain.EnergyConsumer); ok {
		if consumer.DailyEnergyConsumption() < 0 {
			return fmt.Errorf("энергопотребление не может быть отрицательным")
		}
	}

	if _, exists := s.inventory.FindByInventoryNumber(number); exists {
		return fmt.Errorf(
			"инвентарный номер %q уже занят",
			number,
		)
	}

	if err := s.inspection.Inspect(item); err != nil {
		return fmt.Errorf("техосмотр не пройден: %w", err)
	}

	if err := s.inventory.Add(item); err != nil {
		return fmt.Errorf("не удалось добавить объект: %w", err)
	}

	return nil
}

func (s *FleetService) AllItems() []domain.InventoryItem {
	return s.inventory.All()
}

func (s *FleetService) TotalDailyEnergy() float64 {
	var total float64

	for _, item := range s.inventory.All() {
		consumer, ok := item.(domain.EnergyConsumer)
		if !ok {
			continue
		}

		total += consumer.DailyEnergyConsumption()
	}

	return total
}

func (s *FleetService) BeginnerFriendlyItems() []domain.Transport {
	result := make([]domain.Transport, 0)

	for _, item := range s.inventory.All() {
		transport, ok := item.(domain.Transport)
		if !ok {
			continue
		}

		if transport.BeginnerLevel() >= 6 {
			result = append(result, transport)
		}
	}

	return result
}
