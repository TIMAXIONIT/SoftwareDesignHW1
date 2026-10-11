package service

import (
	"fmt"

	"vyshkat/domain"
)

type Fleet struct {
	items []domain.InventoryItem
}

func NewFleet() *Fleet {
	return &Fleet{
		items: make([]domain.InventoryItem, 0),
	}
}

func (f *Fleet) Add(item domain.InventoryItem) error {
	if item == nil {
		return fmt.Errorf("нельзя добавить пустой объект")
	}

	number := item.InventoryNumber()
	if number == "" {
		return fmt.Errorf("инвентарный номер не может быть пустым")
	}

	if _, exists := f.FindByInventoryNumber(number); exists {
		return fmt.Errorf(
			"инвентарный номер %q уже занят",
			number,
		)
	}

	f.items = append(f.items, item)
	return nil
}

func (f *Fleet) FindByInventoryNumber(
	number string,
) (domain.InventoryItem, bool) {
	for _, item := range f.items {
		if item.InventoryNumber() == number {
			return item, true
		}
	}

	return nil, false
}

func (f *Fleet) All() []domain.InventoryItem {
	result := make([]domain.InventoryItem, len(f.items))
	copy(result, f.items)
	return result
}
