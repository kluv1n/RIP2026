export const ROUTES = {
  /** Каталог типов аккумуляторов (`battery_life_types`). */
  BATTERY_TYPES: "/",
  /** Карточка одного типа. */
  BATTERY_TYPE: "/battery/:id",
} as const;

export type RouteKeyType = keyof typeof ROUTES;

export const ROUTE_LABELS: { [key in RouteKeyType]: string } = {
  BATTERY_TYPES: "Каталог типов аккумуляторов",
  BATTERY_TYPE: "Тип аккумулятора",
};
