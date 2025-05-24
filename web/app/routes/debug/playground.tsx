import { useSubmit } from "react-router";
import { PreparingPage } from "../playgrounds/component";
import { usePreparingPageProps } from "../playgrounds/props";

const islands = [
  Math.floor(Math.random() * 6 ** 2),
  Math.floor(Math.random() * 6 ** 2),
];
export default function Playground() {
  const submit = useSubmit();
  const enables: number[] = [];
  for (let i = 0; i < 6 ** 2; i++) {
    if (islands.includes(i)) {
      continue;
    }
    enables.push(i);
  }
  const props = usePreparingPageProps({
    timeout: new Date(Date.now() + 1000 * 60), // 60s
    milliSecondPerTurn: 1000 * 30, // 30s
    boardWidth: 6,
    mineCount: 2,
    islandSectors: islands,
    enableDeploySectors: enables,
    enableMineSectors: enables,
    deployAction: ({ deployAt, mines }) => {
      submit({ deployAt, mines }, { method: "post" });
    },
  });
  return <PreparingPage {...props} />;
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
