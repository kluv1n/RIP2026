import type { FormEvent } from "react";
import Form from "react-bootstrap/Form";
import "./BatteryTypesFilterBar.css";

export interface BatteryTypesFilterBarProps {
  query: string;
  onQueryChange: (query: string) => void;
  onSearch: () => void;
}

export default function BatteryTypesFilterBar({
  query,
  onQueryChange,
  onSearch,
}: BatteryTypesFilterBarProps) {
  const handleSubmit = (e: FormEvent) => {
    e.preventDefault();
    onSearch();
  };

  return (
    <div className="battery-types-filter-bar toolbar__search-form">
      <Form className="search-form battery-types-filter-bar__form" onSubmit={handleSubmit}>
        <span className="search-bar">
          <Form.Control
            type="text"
            name="query"
            className="search-input"
            placeholder="Поиск по типу аккумулятора"
            value={query}
            onChange={(e) => onQueryChange(e.target.value)}
          />
          <button type="submit" className="search-btn" aria-label="Найти">
            <svg className="search-btn__icon" viewBox="0 0 24 24" aria-hidden="true">
              <circle cx="11" cy="11" r="6" fill="none" stroke="currentColor" strokeWidth="2" />
              <path
                d="M16 16l4 4"
                fill="none"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
              />
            </svg>
          </button>
        </span>
      </Form>
    </div>
  );
}
