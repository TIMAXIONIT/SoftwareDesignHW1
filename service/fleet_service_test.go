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

	scooter := domain.Scooter{
		NameValue:          "Электросамокат",
		InventoryNumberVal: "S-001",
		BeginnerLevelVal:   8,
		DailyEnergyVal:     0.5,
	}

	err := service.AddItem(scooter)
	if err != nil {
		t.Fatalf("не ожидали ошибку, получили: %v", err)
	}

	if len(service.AllItems()) != 1 {
		t.Fatal("объект должен быть добавлен в парк")
	}

	if inspection.Calls != 1 {
		t.Fatalf(
			"ожидали 1 вызов техосмотра, получили %d",
			inspection.Calls,
		)
	}
}

func TestAddItemInspectionRejected(t *testing.T) {
	fleet := NewFleet()
	inspection := &FakeInspectionService{
		Err: errors.New("неисправны тормоза"),
	}
	service := NewFleetService(fleet, inspection)

	scooter := domain.Scooter{
		NameValue:          "Электросамокат",
		InventoryNumberVal: "S-002",
		BeginnerLevelVal:   8,
		DailyEnergyVal:     0.5,
	}

	err := service.AddItem(scooter)
	if err == nil {
		t.Fatal("ожидали отказ в приёмке")
	}

	if !strings.Contains(err.Error(), "неисправны тормоза") {
		t.Fatalf("причина отказа потеряна: %v", err)
	}

	if len(service.AllItems()) != 0 {
		t.Fatal("отклонённый объект не должен попасть в парк")
	}
}

func TestDuplicateInventoryNumber(t *testing.T) {
	fleet := NewFleet()
	inspection := &FakeInspectionService{}
	service := NewFleetService(fleet, inspection)

	first := domain.Bicycle{
		NameValue:          "Велосипед",
		InventoryNumberVal: "B-001",
		BeginnerLevelVal:   7,
	}

	second := domain.Scooter{
		NameValue:          "Электросамокат",
		InventoryNumberVal: "B-001",
		BeginnerLevelVal:   8,
		DailyEnergyVal:     0.5,
	}

	if err := service.AddItem(first); err != nil {
		t.Fatalf("первый объект не добавлен: %v", err)
	}

	err := service.AddItem(second)
	if err == nil {
		t.Fatal("ожидали отказ из-за дубликата")
	}

	if !strings.Contains(err.Error(), "уже занят") {
		t.Fatalf("неожиданная причина отказа: %v", err)
	}

	if len(service.AllItems()) != 1 {
		t.Fatal("при дубликате состав парка не должен меняться")
	}

	if inspection.Calls != 1 {
		t.Fatalf(
			"техосмотр должен вызываться только для первого объекта; вызовов: %d",
			inspection.Calls,
		)
	}
}

func TestFleetReports(t *testing.T) {
	fleet := NewFleet()
	service := NewFleetService(fleet, NewServiceCenter())

	items := []domain.InventoryItem{
		domain.Scooter{
			NameValue:          "Самокат A",
			InventoryNumberVal: "S-010",
			BeginnerLevelVal:   8,
			DailyEnergyVal:     0.5,
		},
		domain.Bicycle{
			NameValue:          "Велосипед B",
			InventoryNumberVal: "B-010",
			BeginnerLevelVal:   6,
		},
		domain.Bicycle{
			NameValue:          "Велосипед C",
			InventoryNumberVal: "B-011",
			BeginnerLevelVal:   4,
		},
		domain.Helmet{
			NameValue:          "Шлем",
			InventoryNumberVal: "H-010",
		},
		domain.ChargingStation{
			NameValue:          "Зарядная станция",
			InventoryNumberVal: "C-010",
			DailyEnergyVal:     2.0,
		},
	}

	for _, item := range items {
		if err := service.AddItem(item); err != nil {
			t.Fatalf("не удалось добавить %s: %v", item.Name(), err)
		}
	}

	if got, want := service.TotalDailyEnergy(), 2.5; got != want {
		t.Errorf("расход энергии: получили %v, ожидали %v", got, want)
	}

	beginnerItems := service.BeginnerFriendlyItems()
	if len(beginnerItems) != 2 {
		t.Fatalf(
			"ожидали 2 транспорта для новичков, получили %d",
			len(beginnerItems),
		)
	}

	if len(service.AllItems()) != 5 {
		t.Fatal("все пять объектов должны присутствовать в реестре")
	}
}

func TestAddItemEmptyInventoryNumber(t *testing.T) {
	fleet := NewFleet()
	inspection := &FakeInspectionService{}
	fleetService := NewFleetService(fleet, inspection)

	err := fleetService.AddItem(domain.Scooter{
		NameValue:        "Самокат",
		BeginnerLevelVal: 8,
		DailyEnergyVal:   0.5,
	})

	if err == nil {
		t.Fatal("ожидалась ошибка из-за пустого номера")
	}

	if len(fleetService.AllItems()) != 0 {
		t.Fatal("объект с пустым номером не должен добавляться")
	}

	if inspection.Calls != 0 {
		t.Fatal("техосмотр не должен запускаться для объекта с пустым номером")
	}
}

func TestBeginnerFriendlyBoundary(t *testing.T) {
	fleet := NewFleet()
	inspection := &FakeInspectionService{}
	fleetService := NewFleetService(fleet, inspection)

	items := []domain.Bicycle{
		{
			NameValue:          "Подходящий",
			InventoryNumberVal: "BI-001",
			BeginnerLevelVal:   6,
		},
		{
			NameValue:          "Не подходящий",
			InventoryNumberVal: "BI-002",
			BeginnerLevelVal:   5,
		},
	}

	for _, item := range items {
		if err := fleetService.AddItem(item); err != nil {
			t.Fatal(err)
		}
	}

	beginnerItems := fleetService.BeginnerFriendlyItems()

	if len(beginnerItems) != 1 {
		t.Fatalf("ожидался 1 объект, получено %d", len(beginnerItems))
	}

	if beginnerItems[0].InventoryNumber() != "BI-001" {
		t.Fatal("в список попал неправильный объект")
	}
}

func TestEmptyInventoryReports(t *testing.T) {
	fleet := NewFleet()
	inspection := &FakeInspectionService{}
	fleetService := NewFleetService(fleet, inspection)

	if len(fleetService.AllItems()) != 0 {
		t.Fatal("инвентарь должен быть пустым")
	}

	if got := fleetService.TotalDailyEnergy(); got != 0 {
		t.Fatalf("ожидалось энергопотребление 0, получено %f", got)
	}

	if got := len(fleetService.BeginnerFriendlyItems()); got != 0 {
		t.Fatalf("ожидался пустой список новичков, получено %d объектов", got)
	}
}
