import { useEffect, useState } from "react";
import CatalogChrome from "../../components/CatalogChrome/CatalogChrome";
import BatteryTypesFilterBar from "../../components/BatteryTypesFilterBar/BatteryTypesFilterBar";
import BatteryTypesList from "../../components/BatteryTypesList/BatteryTypesList";
import type { BatteryTypeMock } from "../../modules/batteryApi";
import { BATTERIES_MOCK, filterMockBatteries, type BatteryFilters } from "../../modules/mock";

const initialFilters = (): BatteryFilters => ({ title: "" });

export default function BatteryTypesPage() {
  const [batteries, setBatteries] = useState<BatteryTypeMock[]>(BATTERIES_MOCK);
  const [filters, setFilters] = useState<BatteryFilters>(initialFilters);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    setBatteries(filterMockBatteries(initialFilters()));
  }, []);

  const applyFilters = () => {
    setLoading(true);
    try {
      setBatteries(filterMockBatteries(filters));
    } finally {
      setLoading(false);
    }
  };

  const toolbarForm = (
    <BatteryTypesFilterBar
      query={filters.title}
      onQueryChange={(q) => setFilters((f) => ({ ...f, title: q }))}
      onSearch={applyFilters}
    />
  );

  return (
    <>
      <CatalogChrome toolbarForm={toolbarForm} />
      <div className="space">
        <h2 className="section-title">Типы аккумуляторов</h2>
        {loading ? (
          <p style={{ color: "var(--neter-text-muted)" }}>Загрузка…</p>
        ) : batteries.length > 0 ? (
          <BatteryTypesList batteries={batteries} />
        ) : (
          <p style={{ color: "var(--neter-text-muted)" }}>По заданным фильтрам типы не найдены.</p>
        )}
      </div>
    </>
  );
}
