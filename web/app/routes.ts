import {
  type RouteConfig,
  index,
  prefix,
  route,
} from "@react-router/dev/routes";

export default [
  index("routes/home/index.tsx"),
  route("playgrounds/:gameId/:playerId", "routes/playgrounds/index.tsx"),
  ...prefix("debug", [
    index("routes/debug/index.tsx"),
    route("playground", "routes/debug/playground.tsx"),
  ]),
] satisfies RouteConfig;
