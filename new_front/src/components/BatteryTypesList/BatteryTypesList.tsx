import type { BatteryTypeMock } from "../../modules/batteryApi";
import BatteryTypeCard from "../BatteryTypeCard/BatteryTypeCard";

export default function BatteryTypesList({ batteries }: { batteries: BatteryTypeMock[] }) {
  return (
    <div className="container">
      {batteries.map((b) => (
        <BatteryTypeCard key={b.battery_id} battery={b} />
      ))}
    </div>
  );
}
