-- Добавить один элемент во второй (неполный) ряд — Аккумулятор LiFePO4.
-- Выполнить, если в списке карточек только 3 штуки и нет второго ряда.
INSERT INTO battery_types (id, title, capacity_mah, voltage_v, photo, video, short_description, description, is_deleted) VALUES
  (4, 'Аккумулятор LiFePO4', 3200, 3.2, 'Li_fe_po.jpg', 'Li_fe_po.mp4', 'Долгий срок службы, безопасность, стабильное напряжение 3,2 В', 'Литий-железо-фосфатные (LiFePO4) аккумуляторы отличаются высокой термостабильностью, большим числом циклов заряда-разряда и отсутствием риска возгорания. Номинальное напряжение ячейки 3,2 В. Широко применяются в накопителях энергии, электровелосипедах и силовом оборудовании.', false)
ON CONFLICT (id) DO NOTHING;

SELECT setval(pg_get_serial_sequence('battery_types', 'id'), (SELECT COALESCE(MAX(id), 1) FROM battery_types));
