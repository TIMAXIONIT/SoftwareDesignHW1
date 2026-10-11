package service

import (
	"errors"
	"strings"
	"testing"

	"vyshkat/domain"
)

type FakeInspectionService struct {
	Err   error
	Calls int
}

func (f *FakeInspectionService) Inspect(
	item domain.InventoryItem,
) error {
	f.Calls++
	return f.Err
}

func TestAddItemAccepted(t *testing.T) {
	fleet := NewFleet()
	inspection := &FakeInspectionService{}
	service := NewFleetService(fleet, inspection)

	item := domain.Scooter{
		NameValue:          "Самокат",
		InventoryNumberVal: "SC-001",
		BeginnerLevelVal:   8,
		DailyEnergyVal:     0.5,
	}

	if err := service.AddItem(item); err != nil {
		t.Fatalf("не ожидалась ошибка: %v", err)
	}

	if len(service.AllItems()) != 1 {
		t.Fatal("объект должен быть добавлен")
	}

	if inspection.Calls != 1 {
		t.Fatalf("ожидался 1 вызов техосмотра, получено %d",
			inspection.Calls)
	}
}

func TestAddItemInspectionRejected(t *testing.T) {
	fleet := NewFleet()
	inspection := &FakeInspectionService{
		Err: errors.New("неисправны тормоза"),
	}
	service := NewFleetService(fleet, inspection)

	err := service.AddItem(domain.Scooter{
		NameValue:          "Самокат",
		InventoryNumberVal: "SC-001",
		BeginnerLevelVal:   8,
		DailyEnergyVal:     0.5,
	})

	if err == nil {
		t.Fatal("ожидался отказ техосмотра")
	}

	if !strings.Contains(err.Error(), "неисправны тормоза") {
		t.Fatalf("причина отказа потеряна: %v", err)
	}

	if len(service.AllItems()) != 0 {
		t.Fatal("отклонённый объект не должен попасть в инвентарь")
	}
}

func TestDuplicateInventoryNumber(t *testing.T) {
	fleet := NewFleet()
	inspection := &FakeInspectionService{}
	service := NewFleetService(fleet, inspection)

	first := domain.Bicycle{
		NameValue:          "Велосипед 1",
		InventoryNumberVal: "BI-001",
		BeginnerLevelVal:   7,
	}
	second := domain.Scooter{
		NameValue:          "Самокат",
		InventoryNumberVal: "BI-001",
		BeginnerLevelVal:   8,
		DailyEnergyVal:     0.5,
	}

	if err := service.AddItem(first); err != nil {
		t.Fatal(err)
	}

	err := service.AddItem(second)
	if err == nil {
		t.Fatal("ожидалась ошибка дубликата")
	}

	if !strings.Contains(err.Error(), "BI-001") {
		t.Fatalf("в ошибке должен быть указан номер: %v", err)
	}

	if len(service.AllItems()) != 1 {
		t.Fatal("дубликат не должен добавляться")
	}

	if inspection.Calls != 1 {
		t.Fatalf("техосмотр должен вызываться один раз, получено %d",
			inspection.Calls)
	}
}

func TestAddItemEmptyInventoryNumber(t *testing.T) {
	fleet := NewFleet()
	inspection := &FakeInspectionService{}
	service := NewFleetService(fleet, inspection)

	err := service.AddItem(domain.Scooter{
		NameValue:        "Самокат",
		BeginnerLevelVal: 8,
	})

	if err == nil {
		t.Fatal("ожидалась ошибка пустого номера")
	}

	if len(service.AllItems()) != 0 {
		t.Fatal("объект с пустым номером не должен добавляться")
	}

	if inspection.Calls != 0 {
		t.Fatal("техосмотр не должен запускаться")
	}
}

func TestReports(t *testing.T) {
	fleet := NewFleet()
	inspection := &FakeInspectionService{}
	service := NewFleetService(fleet, inspection)

	items := []domain.InventoryItem{
		domain.Scooter{
			NameValue:          "Самокат",
			InventoryNumberVal: "SC-001",
			BeginnerLevelVal:   8,
			DailyEnergyVal:     0.5,
		},
		domain.Bicycle{
			NameValue:          "Велосипед для новичка",
			InventoryNumberVal: "BI-001",
			BeginnerLevelVal:   6,
		},
		domain.Bicycle{
			NameValue:          "Сложный велосипед",
			InventoryNumberVal: "BI-002",
			BeginnerLevelVal:   5,
		},
		domain.Helmet{
			NameValue:          "Шлем",
			InventoryNumberVal: "HE-001",
		},
		domain.ChargingStation{
			NameValue:          "Зарядная станция",
			InventoryNumberVal: "CH-001",
			DailyEnergyVal:     2.0,
		},
	}

	for _, item := range items {
		if err := service.AddItem(item); err != nil {
			t.Fatal(err)
		}
	}

	if got := len(service.AllItems()); got != 5 {
		t.Fatalf("ожидалось 5 объектов, получено %d", got)
	}

	if got := service.TotalDailyEnergy(); got != 2.5 {
		t.Fatalf("ожидалось 2.5 кВт·ч, получено %f", got)
	}

	beginners := service.BeginnerFriendlyItems()
	if len(beginners) != 2 {
		t.Fatalf("ожидалось 2 объекта для новичков, получено %d",
			len(beginners))
	}
}

func TestEmptyInventoryReports(t *testing.T) {
	fleet := NewFleet()
	inspection := &FakeInspectionService{}
	service := NewFleetService(fleet, inspection)

	if got := len(service.AllItems()); got != 0 {
		t.Fatalf("ожидался пустой инвентарь, получено %d объектов", got)
	}

	if got := service.TotalDailyEnergy(); got != 0 {
		t.Fatalf("ожидалось энергопотребление 0, получено %f", got)
	}

	if got := len(service.BeginnerFriendlyItems()); got != 0 {
		t.Fatalf("ожидался пустой список, получено %d объектов", got)
	}
}

func TestDemoServiceCenter(t *testing.T) {
	center := NewDemoServiceCenter()

	rejected := domain.Helmet{
		NameValue:          "Бракованный шлем",
		InventoryNumberVal: "FAIL-001",
	}

	if err := center.Inspect(rejected); err == nil {
		t.Fatal("объект с номером FAIL-001 должен быть отклонён")
	}

	accepted := domain.Helmet{
		NameValue:          "Обычный шлем",
		InventoryNumberVal: "HE-001",
	}

	if err := center.Inspect(accepted); err != nil {
		t.Fatalf("объект должен пройти техосмотр: %v", err)
	}
}

func TestInvalidBeginnerLevel(t *testing.T) {
	fleet := NewFleet()
	inspection := &FakeInspectionService{}
	service := NewFleetService(fleet, inspection)

	err := service.AddItem(domain.Bicycle{
		NameValue:          "Велосипед",
		InventoryNumberVal: "BI-001",
		BeginnerLevelVal:   11,
	})

	if err == nil {
		t.Fatal("ожидалась ошибка недопустимого уровня")
	}

	if len(service.AllItems()) != 0 {
		t.Fatal("объект с недопустимым уровнем не должен добавляться")
	}

	if inspection.Calls != 0 {
		t.Fatal("техосмотр не должен запускаться до проверки данных")
	}
}
