import { type RouteConfig, index, prefix } from "@react-router/dev/routes";

export default [
  index("routes/home/index.tsx"),
  ...prefix("debug", [index("routes/debug/index.tsx")]),
] satisfies RouteConfig;
