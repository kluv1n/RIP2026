# AI_NOTES — new_front (RIP2026)

## Что сделано

- Собран SPA на **React 19 + TypeScript + Vite 7** с **React-Bootstrap** по той же схеме слоёв, что и в `RIP_Frontend2026`: `layouts/`, `pages/<Page>/<Page>.tsx` + CSS, `components/<Name>/`, `modules/` (API-типы + mock).
- Предметная область и визуал выровнены с бэкенд-шаблоном `rip2026/templates` и `rip2026/resources/css/styles.css`: тёмная тема ShardDB, карточки стратегий с коэффициентами, детальная страница (`detail-card`), страница заявки (`system-load-detail`, `load-table`, `btn-delete`).
- Три страницы для сценария лабораторной 5 (аналог «услуг» / просмотр для гостя):
  1. Каталог стратегий с **mock**-данными и **поиском по названию** (как в `templates/index.html`, поле «Поиск стратегий…»). Для всех стратегий одни и те же локальные медиа: `public/mock/strategy-cover.png`, `public/mock/composite_sharding.mp4` (без MinIO).
  2. Карточка стратегии: видео при наличии `video`, иначе изображение; **заглушка**, если `photo_url` пустой.
  3. Страница заявки (`/system_load/:id`) с таблицей позиций и удалением (mock, переход на главную).
- **Navbar** (React-Bootstrap) в шапке и самописные **breadcrumbs** без Redux/Context.
- В `vite.config.ts` настроен **proxy** `/api` → `http://localhost:8080` для демонстрации с бэкендом; в `strategiesApi.ts` добавлен пример `fetchStrategiesByTitle` (основной сценарий лабы — работа на mock без сервера).

## Почему так

- Структура папок повторяет эталонный фронтенд-проект, чтобы совпадали границы слоёв и привычный разбор на защите; имена сущностей заменены на домен RIP2026 (стратегии шардирования, system load), как в шаблонах Go/HTML.
- Стили шаблона вынесены в `src/shard-template.css` (копия `styles.css` бэкенда), поверх — минимальные правки в `index_style.css` (хлебные крошки, корзина, фильтры, Navbar).
- Константы маршрутов вынесены в `routePaths.ts` вместо `Routes.tsx`, чтобы проходить правило ESLint `react-refresh/only-export-components`.

## Риски и ограничения

- Тексты стратегий в mock совпадают с сидом **`rip2026/migrations/init-up.sql`**. На карточках каталога по шаблону **нет** абзаца описания — только заголовок и коэффициенты; полное описание — на странице стратегии.
- На странице заявки в Go-шаблоне `system_load.html` ячейки были только для чтения; в mock добавлены редактируемые поля, соответствующие API: **описание заявки** и по строкам — **data_volume**, **query_count**, **response_time**.
- Отображение медиа с MinIO: для относительных ключей из API задайте `VITE_MINIO_BASE` (например, `http://localhost:9000/test` — как дефолт в конфиге Go). Иначе `resolveMediaUrl` вернёт SVG-заглушку.
- «Добавить в заявку» в mock только имитирует задержку и шлёт событие; счётчик в `CartRow` по-прежнему берётся из статического `MOCK_CART` (как в эталонном примере с заявкой).
- Логотип из шаблона (`/static/img/logo.png`) в репозитории бэкенда в workspace не найден; в шапке используется текстовый бренд **ShardDB** в стиле шаблона.

## Как проверить

```powershell
cd c:\Users\Nikolay\GolangProjects\RIP2026\new_front
npm install
npm run dev
```

Открыть `http://localhost:3000`: каталог, фильтры, карточки, переход на стратегию, иконка заявки → `/system_load/1`.

```powershell
npm run build
npm run lint
```

С бэкендом (порт 8080): поднять `rip2026`, в другом терминале `npm run dev`, в DevTools → Network убедиться, что запросы идут на `/api/...` (proxy). Для ключей MinIO задать в `.env.local`:

```env
VITE_MINIO_BASE=http://localhost:9000/test
```
