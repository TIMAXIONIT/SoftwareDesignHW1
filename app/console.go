package app

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"vyshkat/domain"

	"go.uber.org/fx"
)

type FleetOperations interface {
	AddItem(item domain.InventoryItem) error
	AllItems() []domain.InventoryItem
	TotalDailyEnergy() float64
	BeginnerFriendlyItems() []domain.Transport
}

type Console struct {
	fleet      FleetOperations
	scanner    *bufio.Scanner
	shutdowner fx.Shutdowner
}

func NewConsole(
	fleet FleetOperations,
	shutdowner fx.Shutdowner,
) *Console {
	return &Console{
		fleet:      fleet,
		scanner:    bufio.NewScanner(os.Stdin),
		shutdowner: shutdowner,
	}
}

func (c *Console) readLine(prompt string) (string, error) {
	fmt.Print(prompt)

	if !c.scanner.Scan() {
		if err := c.scanner.Err(); err != nil {
			return "", err
		}
		return "", fmt.Errorf("ввод завершён")
	}

	return strings.TrimSpace(c.scanner.Text()), nil
}

func (c *Console) readInt(prompt string) (int, error) {
	for {
		value, err := c.readLine(prompt)
		if err != nil {
			return 0, err
		}

		number, err := strconv.Atoi(value)
		if err == nil {
			return number, nil
		}

		fmt.Println("Введите целое число.")
	}
}

func (c *Console) readFloat(prompt string) (float64, error) {
	for {
		value, err := c.readLine(prompt)
		if err != nil {
			return 0, err
		}

		number, err := strconv.ParseFloat(
			strings.ReplaceAll(value, ",", "."),
			64,
		)
		if err == nil {
			return number, nil
		}

		fmt.Println("Введите число, например 1.5.")
	}
}

func (c *Console) Run() {
	for {
		fmt.Println("\n=== ВышКат ===")
		fmt.Println("1. Добавить самокат")
		fmt.Println("2. Добавить велосипед")
		fmt.Println("3. Добавить шлем")
		fmt.Println("4. Добавить зарядную станцию")
		fmt.Println("5. Показать весь инвентарь")
		fmt.Println("6. Показать общее количество объектов")
		fmt.Println("7. Показать суточное энергопотребление")
		fmt.Println("8. Показать транспорт для новичков")
		fmt.Println("0. Выход")

		choice, err := c.readInt("Выберите действие: ")
		if err != nil {
			fmt.Println("\nЗавершение работы.")
			return
		}

		switch choice {
		case 1:
			c.addScooter()
		case 2:
			c.addBicycle()
		case 3:
			c.addHelmet()
		case 4:
			c.addChargingStation()
		case 5:
			c.showInventory()
		case 6:
			fmt.Printf("Всего объектов: %d\n", len(c.fleet.AllItems()))
		case 7:
			fmt.Printf("Суточное энергопотребление: %.2f кВт·ч\n",
				c.fleet.TotalDailyEnergy())
		case 8:
			c.showBeginnerFriendly()
		case 0:
			fmt.Println("Завершение работы.")
			_ = c.shutdowner.Shutdown()
			return
		default:
			fmt.Println("Неизвестная команда.")
		}
	}
}

func (c *Console) addScooter() {
	name, number, ok := c.readCommonFields()
	if !ok {
		return
	}

	level, err := c.readInt("Простота для новичка (1–10): ")
	if err != nil {
		fmt.Println(err)
		return
	}

	energy, err := c.readFloat("Суточное энергопотребление (кВт·ч): ")
	if err != nil {
		fmt.Println(err)
		return
	}

	item := domain.Scooter{
		NameValue:          name,
		InventoryNumberVal: number,
		BeginnerLevelVal:   level,
		DailyEnergyVal:     energy,
	}

	c.addItem(item)
}

func (c *Console) addBicycle() {
	name, number, ok := c.readCommonFields()
	if !ok {
		return
	}

	level, err := c.readInt("Простота для новичка (1–10): ")
	if err != nil {
		fmt.Println(err)
		return
	}

	item := domain.Bicycle{
		NameValue:          name,
		InventoryNumberVal: number,
		BeginnerLevelVal:   level,
	}

	c.addItem(item)
}

func (c *Console) addHelmet() {
	name, number, ok := c.readCommonFields()
	if !ok {
		return
	}

	item := domain.Helmet{
		NameValue:          name,
		InventoryNumberVal: number,
	}

	c.addItem(item)
}

func (c *Console) addChargingStation() {
	name, number, ok := c.readCommonFields()
	if !ok {
		return
	}

	energy, err := c.readFloat("Суточное энергопотребление (кВт·ч): ")
	if err != nil {
		fmt.Println(err)
		return
	}

	item := domain.ChargingStation{
		NameValue:          name,
		InventoryNumberVal: number,
		DailyEnergyVal:     energy,
	}

	c.addItem(item)
}

func (c *Console) readCommonFields() (string, string, bool) {
	name, err := c.readLine("Название: ")
	if err != nil {
		fmt.Println(err)
		return "", "", false
	}

	number, err := c.readLine("Инвентарный номер: ")
	if err != nil {
		fmt.Println(err)
		return "", "", false
	}

	return name, number, true
}

func (c *Console) addItem(item domain.InventoryItem) {
	if err := c.fleet.AddItem(item); err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	fmt.Println("Объект успешно прошёл техосмотр и добавлен в парк.")
}

func (c *Console) showInventory() {
	items := c.fleet.AllItems()

	if len(items) == 0 {
		fmt.Println("Инвентарь пуст.")
		return
	}

	fmt.Println("\n=== Инвентарь ===")
	for _, item := range items {
		fmt.Printf("%s | № %s", item.Name(), item.InventoryNumber())

		if transport, ok := item.(domain.Transport); ok {
			fmt.Printf(" | Простота для новичка: %d/10",
				transport.BeginnerLevel())
		}

		if consumer, ok := item.(domain.EnergyConsumer); ok {
			fmt.Printf(" | Энергопотребление: %.2f кВт·ч/сутки",
				consumer.DailyEnergyConsumption())
		}

		fmt.Println()
	}
}

func (c *Console) showBeginnerFriendly() {
	items := c.fleet.BeginnerFriendlyItems()

	if len(items) == 0 {
		fmt.Println("Подходящего для новичков транспорта пока нет.")
		return
	}

	fmt.Println("\n=== Транспорт для новичков ===")
	for _, item := range items {
		fmt.Printf("%s | № %s | Простота: %d/10\n",
			item.Name(),
			item.InventoryNumber(),
			item.BeginnerLevel())
	}
}
