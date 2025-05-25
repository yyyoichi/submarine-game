import { useFetcher } from "react-router";
import { PlayingPage, PreparingPage } from "../playgrounds/component";
import {
  usePlayingPageProps,
  usePreparingPageProps,
} from "../playgrounds/props";

const islands = [
  Math.floor(Math.random() * 6 ** 2),
  Math.floor(Math.random() * 6 ** 2),
];
export default function Playground() {
  const fetcher = useFetcher();
  const enables: number[] = [];
  for (let i = 0; i < 6 ** 2; i++) {
    if (islands.includes(i)) {
      continue;
    }
    enables.push(i);
  }
  const prepareingPageProps = usePreparingPageProps({
    timeout: new Date(Date.now() + 1000 * 60), // 60s
    milliSecondPerTurn: 1000 * 30, // 30s
    boardWidth: 6,
    mineCount: 2,
    islandSectors: islands,
    enableDeploySectors: enables,
    enableMineSectors: enables,
    deployAction: ({ deployAt, mines }) => {
      fetcher.submit({ deployAt, mines }, { method: "post" });
    },
    loading: fetcher.state !== "idle",
  });
  const playingPageProps = usePlayingPageProps({
    timeout: new Date(Date.now() + 1000 * 60), // 60s,
    milliSecondPerTurn: 1000 * 30, // 30s
    boardWidth: 6,
    turn: 3,
    inMyTurn: true,
    at: 10,
    islandSectors: islands,
    enableMoveSectors: [10, 17, 13],
    enableTorpedoSectors: [10, 11, 12, 16, 17, 18],
    enableMineSectors: [11, 17],
    prevAction: {
      me: false,
      mode: "show",
      type: "fire-torpedo",
      at: 11,
      result: "hard-to-starboard",
    },
    action: ({ at, type }) => {
      fetcher.submit({ at, type }, { method: "post" });
    },
    loading: fetcher.state !== "idle",
  });
  return (
    <>
      <PreparingPage {...prepareingPageProps} />
      <PlayingPage {...playingPageProps} />
    </>
  );
}

export async function clientAction() {
  console.log("deploy action");
  return new Promise((resolve) => {
    setTimeout(() => {
      console.log("deploy action done");
      resolve(true);
    }, 1000);
  });
}
