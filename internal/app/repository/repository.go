package repository

import (
	"fmt"
	"strings"
)

type Repository struct{}

func NewRepository() (*Repository, error) {
	return &Repository{}, nil
}

// BatteryType — тип аккумулятора (услуга)
type BatteryType struct {
	ID               int
	Title            string
	CapacityMah      int     // ёмкость, мА·ч
	VoltageV         float64 // напряжение, В
	Photo            string
	Video            string
	ShortDescription string
	Description      string
}

// Application — заявка на расчёт времени работы
type Application struct {
	ID                int
	Title             string
	Description       string
	Items             []ApplicationItem
	ItemCount         int
	TotalRuntimeHours float64 // сумма времени работы по всем аккумуляторам (ч)
}

// ApplicationItem — строка заявки: аккумулятор + потребляемый ток + м-м (количество) → время работы (ч)
type ApplicationItem struct {
	Battery      BatteryType
	CurrentMa    int     // потребляемый ток устройства, мА
	Mm           string  // м-м: количество/порядок/комментарий (в последующих лабах меняет пользователь)
	Quantity     int     // количество — используется для расчёта (в м-м отображаем как вариант "количества")
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
			Photo:            "li_ion.jpg",
			Video:            "li_ion.mp4",
			ShortDescription: "Высокая удельная энергия, малый саморазряд",
			Description:      "Литий-ионные аккумуляторы широко применяются в смартфонах, ноутбуках, электротранспорте. Характеризуются высокой плотностью энергии, отсутствием эффекта памяти. Номинальное напряжение одной ячейки обычно 3,6–3,7 В.",
		},
		{
			ID:               2,
			Title:            "Li-Po (литий-полимерный)",
			CapacityMah:      1500,
			VoltageV:         3.7,
			Photo:            "li_po.jpg",
			Video:            "li_po.mp4",
			ShortDescription: "Гибкая форма, малый вес, высокая токоотдача",
			Description:      "Литий-полимерные аккумуляторы позволяют делать батареи тонкими и гибкими. Часто применяются в дронах, носимой электронике. По удельной энергии и напряжению близки к Li-ion.",
		},
		{
			ID:               3,
			Title:            "Ni-MH (никель-металлгидридный)",
			CapacityMah:      2500,
			VoltageV:         1.2,
			Photo:            "ni_mh.jpg",
			Video:            "ni_mh.mp4",
			ShortDescription: "Экологичность, перезаряжаемость, стабильность при низких температурах",
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

func (r *Repository) buildApplication(id int, title, description string, entries []struct {
	BatteryID int
	CurrentMa int
	Quantity  int
}) (Application, error) {
	batteries, err := r.GetBatteryTypes()
	if err != nil {
		return Application{}, err
	}
	batteryMap := make(map[int]BatteryType)
	for _, b := range batteries {
		batteryMap[b.ID] = b
	}
	var items []ApplicationItem
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
		items = append(items, ApplicationItem{
			Battery:      b,
			CurrentMa:    e.CurrentMa,
			Mm:           fmt.Sprintf("%d", qty), // м-м: количество (в лабе 2+ будет редактироваться)
			Quantity:     qty,
			RuntimeHours: rh,
			RuntimeTotal: rt,
		})
		totalRuntime += rt
	}
	return Application{
		ID:                id,
		Title:             title,
		Description:       description,
		Items:             items,
		ItemCount:         len(items),
		TotalRuntimeHours: totalRuntime,
	}, nil
}

func (r *Repository) GetApplications() ([]Application, error) {
	entries := []struct {
		BatteryID int
		CurrentMa int
		Quantity  int
	}{
		{1, 500, 1},
		{2, 300, 2},
		{3, 200, 3},
	}
	app, err := r.buildApplication(
		1,
		"Расчёт времени работы в часах для устройства с указанным потребляемым током и выбранным типом аккумулятора",
		"Расчёт времени работы в часах для устройства с указанным потребляемым током и выбранным типом аккумулятора. Время (ч) = ёмкость (мА·ч) / ток (мА).",
		entries,
	)
	if err != nil {
		return nil, err
	}
	return []Application{app}, nil
}

func (r *Repository) GetApplication(id int) (Application, error) {
	apps, err := r.GetApplications()
	if err != nil {
		return Application{}, err
	}
	for _, app := range apps {
		if app.ID == id {
			return app, nil
		}
	}
	return Application{}, fmt.Errorf("заявка не найдена")
}

func (r *Repository) GetApplicationForBattery(batteryID int) (*ApplicationItem, error) {
	apps, err := r.GetApplications()
	if err != nil {
		return nil, err
	}
	for _, app := range apps {
		for i := range app.Items {
			if app.Items[i].Battery.ID == batteryID {
				return &app.Items[i], nil
			}
		}
	}
	return nil, fmt.Errorf("аккумулятор не найден в заявке")
}
