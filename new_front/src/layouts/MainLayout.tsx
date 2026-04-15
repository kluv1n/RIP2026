import { Outlet } from "react-router-dom";

/** Как в Gin-шаблонах: без общего хедера — каждая страница сама подключает свой кусок вёрстки. */
export default function MainLayout() {
  return <Outlet />;
}
