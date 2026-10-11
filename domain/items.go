package domain

type Scooter struct {
	NameValue          string
	InventoryNumberVal string
	BeginnerLevelVal   int
	DailyEnergyVal     float64
}

func (s Scooter) Name() string {
	return s.NameValue
}

func (s Scooter) InventoryNumber() string {
	return s.InventoryNumberVal
}

func (s Scooter) BeginnerLevel() int {
	return s.BeginnerLevelVal
}

func (s Scooter) DailyEnergyConsumption() float64 {
	return s.DailyEnergyVal
}

type Bicycle struct {
	NameValue          string
	InventoryNumberVal string
	BeginnerLevelVal   int
}

func (b Bicycle) Name() string {
	return b.NameValue
}

func (b Bicycle) InventoryNumber() string {
	return b.InventoryNumberVal
}

func (b Bicycle) BeginnerLevel() int {
	return b.BeginnerLevelVal
}

type Helmet struct {
	NameValue          string
	InventoryNumberVal string
}

func (h Helmet) Name() string {
	return h.NameValue
}

func (h Helmet) InventoryNumber() string {
	return h.InventoryNumberVal
}

type ChargingStation struct {
	NameValue          string
	InventoryNumberVal string
	DailyEnergyVal     float64
}

func (c ChargingStation) Name() string {
	return c.NameValue
}

func (c ChargingStation) InventoryNumber() string {
	return c.InventoryNumberVal
}

func (c ChargingStation) DailyEnergyConsumption() float64 {
	return c.DailyEnergyVal
}
