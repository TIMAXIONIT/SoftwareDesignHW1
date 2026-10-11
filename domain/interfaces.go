package domain

// InventoryItem — любой объект, который учитывается
// в реестре имущества.
type InventoryItem interface {
	Name() string
	InventoryNumber() string
}

// EnergyConsumer — объект, который потребляет энергию.
type EnergyConsumer interface {
	DailyEnergyConsumption() float64
}

// Transport — транспорт, который можно оценить
// по простоте использования для новичка.
type Transport interface {
	InventoryItem
	BeginnerLevel() int
}

type InspectionService interface {
	Inspect(item InventoryItem) error
}

type Inventory interface {
	Add(item InventoryItem) error
	FindByInventoryNumber(number string) (InventoryItem, bool)
	All() []InventoryItem
}
