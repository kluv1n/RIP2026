import { BrowserRouter, Navigate, Route, Routes } from "react-router-dom";
import MainLayout from "./layouts/MainLayout";
import BatteryTypesPage from "./pages/BatteryTypesPage/BatteryTypesPage";
import BatteryTypePage from "./pages/BatteryTypePage/BatteryTypePage";
import { ROUTES } from "./routePaths";
import "bootstrap/dist/css/bootstrap.min.css";
import "./index_style.css";

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route element={<MainLayout />}>
          <Route path={ROUTES.BATTERY_TYPES} element={<BatteryTypesPage />} />
          <Route path="/catalog" element={<Navigate to="/" replace />} />
          <Route path={ROUTES.BATTERY_TYPE} element={<BatteryTypePage />} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
}

export default App;
