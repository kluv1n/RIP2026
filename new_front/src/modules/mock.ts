import type { BatteryTypeMock } from "./batteryApi";

/** Видео: Mixkit, набор батарейных элементов (крупный план, 360p). */
export const MOCK_VIDEO = "/mock/battery-default.mp4";

/** Обложка: фото с аккумулятором/зарядкой (Unsplash → `public/mock/battery-default.jpg`). */
export const MOCK_COVER = "/mock/battery-default.jpg";

/** Типы аккумуляторов для каталога (mock `battery_life_types`). */
export const BATTERIES_MOCK: BatteryTypeMock[] = [
  {
    battery_id: 1,
    is_deleted: false,
    title: "Li-ion 18650",
    short_description: "Универсальный цилиндр 18650 для powerbank и ноутбуков.",
    description:
      "Цилиндрический литий-ионный элемент 18650. Часто используется в powerbank, ноутбуках и переносной электронике.",
    photo_url: MOCK_COVER,
    video: MOCK_VIDEO,
    capacity_mah: 3500,
    voltage_v: 3.7,
    price_rub: 890,
    listed_at: "2026-01-10T12:00:00.000Z",
    detail_current_a_str: "0,45",
    detail_runtime_hours_str: "12,40",
  },
  {
    battery_id: 2,
    is_deleted: false,
    title: "Li-Po пакет",
    short_description: "Плоский полимерный элемент для компактных устройств.",
    description: "Полимерно-литиевый пакет: компактная форма для носимых устройств и дронов.",
    photo_url: MOCK_COVER,
    video: MOCK_VIDEO,
    capacity_mah: 5000,
    voltage_v: 3.85,
    price_rub: 1240,
    listed_at: "2026-02-05T09:00:00.000Z",
    detail_current_a_str: "0,80",
    detail_runtime_hours_str: "5,20",
  },
  {
    battery_id: 3,
    is_deleted: false,
    title: "LiFePO₄ блок",
    short_description: "Безопасная химия для стационарных и транспортных систем.",
    description: "Литий-железо-фосфатный аккумуляторный блок: безопасность и долгий ресурс циклов.",
    photo_url: MOCK_COVER,
    video: MOCK_VIDEO,
    capacity_mah: 100_000,
    voltage_v: 12.8,
    price_rub: 42_500,
    listed_at: "2026-01-22T15:30:00.000Z",
    detail_current_a_str: "2,50",
    detail_runtime_hours_str: "36,00",
  },
  {
    battery_id: 4,
    is_deleted: false,
    title: "Ni-MH AA",
    short_description: "Формат AA, 1.2 В — пульты, фонари, бытовая электроника.",
    description: "Никель-металлгидридный элемент формата AA: 1.2 В, удобен для пультов и фонарей.",
    photo_url: MOCK_COVER,
    video: MOCK_VIDEO,
    capacity_mah: 2500,
    voltage_v: 1.2,
    price_rub: 320,
    listed_at: "2025-12-18T11:00:00.000Z",
    detail_current_a_str: "0,12",
    detail_runtime_hours_str: "18,50",
  },
];

export interface BatteryFilters {
  title: string;
}

export function getMockBattery(id: number): BatteryTypeMock | undefined {
  return BATTERIES_MOCK.find((b) => b.battery_id === id);
}

export function filterMockBatteries(filters: BatteryFilters): BatteryTypeMock[] {
  let list = BATTERIES_MOCK.filter((b) => !b.is_deleted);
  const t = filters.title.trim().toLowerCase();
  if (t) {
    list = list.filter(
      (b) =>
        b.title.toLowerCase().includes(t) ||
        b.short_description.toLowerCase().includes(t),
    );
  }
  return list;
}
