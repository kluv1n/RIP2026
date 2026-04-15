import { Link } from "react-router-dom";
import { useEffect, useState } from "react";
import { fallbackImageUrl, resolveMediaUrl, type BatteryTypeMock } from "../../modules/batteryApi";

function photoSrc(photo_url: string, imageError: boolean): string {
  if (imageError || !photo_url?.trim()) return fallbackImageUrl();
  return resolveMediaUrl(photo_url);
}

/** Карточка типа аккумулятора в каталоге (гость: без «Добавить в заявку»). */
export default function BatteryTypeCard({ battery }: { battery: BatteryTypeMock }) {
  const [imageError, setImageError] = useState(false);
  const [imageUrl, setImageUrl] = useState(photoSrc(battery.photo_url, false));

  useEffect(() => {
    setImageError(false);
    setImageUrl(photoSrc(battery.photo_url, false));
  }, [battery.photo_url]);

  const handleImageError = () => {
    setImageError(true);
    setImageUrl(fallbackImageUrl());
  };

  return (
    <div className="card">
      <Link to={`/battery/${battery.battery_id}`} className="card__link">
        <div className="card__media">
          <img
            className="card__photo"
            src={imageError ? fallbackImageUrl() : imageUrl}
            alt={battery.title}
            width={400}
            height={300}
            decoding="async"
            onError={handleImageError}
          />
        </div>
        <h1>{battery.title}</h1>
        <p className="card__employees">
          Ёмкость и напряжение: {battery.capacity_mah} мА·ч, {battery.voltage_v} В
        </p>
        <p className="card__description">{battery.short_description}</p>
      </Link>
    </div>
  );
}
