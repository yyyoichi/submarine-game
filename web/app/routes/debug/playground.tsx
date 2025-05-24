import { useState } from "react";
import { useSubmit } from "react-router";
import { PlayingPage } from "../playgrounds/component";
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
  const [turn, setTurn] = useState(1);
  const playingPageProps: React.ComponentProps<typeof PlayingPage> = {
    turn: turn,
    Leading: {
      inMyTrun: true,
      Alerm: {
        useAlerm: true,
        startPingMilliSec: 1000 * 30, // 30s
        finishDatetime: new Date(Date.now() + 1000 * 30), // 30s
      },
      children: "C1に敵艦魚雷着弾！ヨーソロー！C1付近の敵艦を攻撃せよ！",
    },
    Ocean: {
      sectorCount: 6,
      islands: islands,
      me: 10,
      enables: enables,
      sector: 11,
      type: "torpedo",
    },
    Controller: {
      Direction: {
        North: {},
        South: {},
        East: {},
        West: {},
      },
      AButton: {
        onClick: () => setTurn((v) => v + 1),
      },
      BButton: {},
      CommandWindow: {
        commands: ["move", "fire-torpedo", "trigger-mine"],
        use: "fire-torpedo",
      },
    },
    OverlayedLeading: {
      children: "C1に敵艦魚雷着弾！ヨーソロー！",
    },
    OverlayedOcean: {
      gapSector: 2,
      ringSector: 2,
      title: "魚雷攻撃",
      titlePosition: "left",
    },
  };
  return (
    <>
      {/* <PreparingPage {...props} /> */}
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
