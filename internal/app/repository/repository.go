package repository

import (
	"fmt"
	"strings"
)

type Repository struct{}

func NewRepository() (*Repository, error) {
	return &Repository{}, nil
}

// BatteryType — тип аккумулятора (услуга: ёмкость и напряжение)
type BatteryType struct {
	ID               int
	Title            string
	CapacityMah      int     // ёмкость, мА·ч
	VoltageV         float64 // напряжение, В
	DemoLoadMa       int     // типовой ток нагрузки для блока «Ток / время» на карточке услуги, мА
	Photo            string
	Video            string
	ShortDescription string
	Description      string
}

// BatteryLife — расчёт времени работы (battery life)
type BatteryLife struct {
	ID                int
	Title             string
	Description       string
	Items             []BatteryLifeItem
	ItemCount         int
	TotalRuntimeHours float64 // итоговое время работы по всем аккумуляторам (ч)
}

// BatteryLifeItem — строка расчёта: аккумулятор + потребляемый ток (заявка) + м-м (количество) → время работы (ч)
type BatteryLifeItem struct {
	Battery      BatteryType
	CurrentMa    int     // потребляемый ток устройства, мА (поле "Заявка")
	Mm           string  // м-м: количество (для отображения)
	Quantity     int     // количество — используется для расчёта
	RuntimeHours float64 // время работы одного аккумулятора, ч (ёмкость / ток)
	RuntimeTotal float64 // время × количество, ч (вклад в сумму)
}

func (r *Repository) GetBatteryTypes() ([]BatteryType, error) {
	batteries := []BatteryType{
		{
			ID:               1,
			Title:            "Li-ion (литий-ионный)",
			CapacityMah:      3000,
			VoltageV:         3.7,
			DemoLoadMa:       500, // ~6 ч при 3000 мА·ч (оценка: ч = мА·ч / мА)
			Photo:            "li_ion.jpg",
			Video:            "li_ion.mp4",
			ShortDescription: "Высокая энергия, малый саморазряд.",
			Description:      "Литий-ионные аккумуляторы широко применяются в смартфонах, ноутбуках, электротранспорте. Характеризуются высокой плотностью энергии, отсутствием эффекта памяти. Номинальное напряжение одной ячейки обычно 3,6–3,7 В.",
		},
		{
			ID:               2,
			Title:            "Li-Po (литий-полимерный)",
			CapacityMah:      1500,
			VoltageV:         3.7,
			DemoLoadMa:       250,
			Photo:            "li_po.jpg",
			Video:            "li_po.mp4",
			ShortDescription: "Лёгкий, гибкий, высокий ток.",
			Description:      "Литий-полимерные аккумуляторы позволяют делать батареи тонкими и гибкими. Часто применяются в дронах, носимой электронике. По удельной энергии и напряжению близки к Li-ion.",
		},
		{
			ID:               3,
			Title:            "Ni-MH (никель-металлгидридный)",
			CapacityMah:      2500,
			VoltageV:         1.2,
			DemoLoadMa:       200,
			Photo:            "ni_mh.jpg",
			Video:            "ni_mh.mp4",
			ShortDescription: "Без кадмия, стабилен в холоде.",
			Description:      "Никель-металлгидридные аккумуляторы — перезаряжаемая альтернатива без кадмия. Номинальное напряжение элемента 1,2 В. Используются в бытовой технике, гибридном транспорте.",
		},
	}
	return batteries, nil
}

func (r *Repository) GetBattery(id int) (BatteryType, error) {
	batteries, err := r.GetBatteryTypes()
	if err != nil {
		return BatteryType{}, err
	}
	for _, b := range batteries {
		if b.ID == id {
			return b, nil
		}
	}
	return BatteryType{}, fmt.Errorf("тип аккумулятора не найден")
}

func (r *Repository) GetBatteryByTitle(title string) ([]BatteryType, error) {
	batteries, err := r.GetBatteryTypes()
	if err != nil {
		return nil, err
	}
	var result []BatteryType
	lower := strings.ToLower(title)
	for _, b := range batteries {
		if strings.Contains(strings.ToLower(b.Title), lower) {
			result = append(result, b)
		}
	}
	return result, nil
}

func RuntimeHours(capacityMah, currentMa int) float64 {
	if currentMa <= 0 {
		return 0
	}
	return float64(capacityMah) / float64(currentMa)
}

func (r *Repository) buildBatteryLife(id int, title, description string, entries []struct {
	BatteryID int
	CurrentMa int
	Quantity  int
}) (BatteryLife, error) {
	batteries, err := r.GetBatteryTypes()
	if err != nil {
		return BatteryLife{}, err
	}
	batteryMap := make(map[int]BatteryType)
	for _, b := range batteries {
		batteryMap[b.ID] = b
	}
	var items []BatteryLifeItem
	var totalRuntime float64
	for _, e := range entries {
		b, ok := batteryMap[e.BatteryID]
		if !ok {
			continue
		}
		qty := e.Quantity
		if qty <= 0 {
			qty = 1
		}
		rh := RuntimeHours(b.CapacityMah, e.CurrentMa)
		rt := rh * float64(qty)
		items = append(items, BatteryLifeItem{
			Battery:      b,
			CurrentMa:    e.CurrentMa,
			Mm:           fmt.Sprintf("%d", qty), // м-м: количество (в лабе 2+ будет редактироваться)
			Quantity:     qty,
			RuntimeHours: rh,
			RuntimeTotal: rt,
		})
		totalRuntime += rt
	}
	return BatteryLife{
		ID:                id,
		Title:             title,
		Description:       description,
		Items:             items,
		ItemCount:         len(items),
		TotalRuntimeHours: totalRuntime,
	}, nil
}

func (r *Repository) GetBatteryLives() ([]BatteryLife, error) {
	entries := []struct {
		BatteryID int
		CurrentMa int
		Quantity  int
	}{
		{1, 500, 1},
		{2, 300, 2},
		{3, 200, 3},
	}
	app, err := r.buildBatteryLife(
		1,
		"Расчёт времени работы в часах для устройства с указанным потребляемым током и выбранным типом аккумулятора",
		"Расчёт времени работы в часах для устройства с указанным потребляемым током и выбранным типом аккумулятора. Время (ч) = ёмкость (мА·ч) / ток (мА).",
		entries,
	)
	if err != nil {
		return nil, err
	}
	return []BatteryLife{app}, nil
}

func (r *Repository) GetBatteryLife(id int) (BatteryLife, error) {
	lives, err := r.GetBatteryLives()
	if err != nil {
		return BatteryLife{}, err
	}
	for _, life := range lives {
		if life.ID == id {
			return life, nil
		}
	}
	return BatteryLife{}, fmt.Errorf("battery life не найдена")
}

func (r *Repository) GetBatteryLifeForBattery(batteryID int) (*BatteryLifeItem, error) {
	lives, err := r.GetBatteryLives()
	if err != nil {
		return nil, err
	}
	for _, life := range lives {
		for i := range life.Items {
			if life.Items[i].Battery.ID == batteryID {
				return &life.Items[i], nil
			}
		}
	}
	return nil, fmt.Errorf("аккумулятор не найден в battery life")
}
