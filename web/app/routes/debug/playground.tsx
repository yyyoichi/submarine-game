import type React from "react";
import { useState } from "react";
import { PreparingPage } from "../playgrounds/component";

export default function Playground() {
  const [viewGuide, setViewGuide] = useState(false);
  const time30secAgo = new Date(Date.now() - 1000 * 30);
  const props: React.ComponentProps<typeof PreparingPage> = {
    Leading: {
      inMyTrun: true,
      Alerm: {
        useAlerm: true,
        startPingSec: 30,
        finishDatetime: new Date(Date.now() + 1000 * 30),
      },
    },
    Ocean: {
      displaySectorName: viewGuide,
      sectorCount: 6,
      Sectors: {
        1: { color: "embed" },
        16: { color: "embed" },
      },
    },
    Controller: {
      Direction: {
        North: {},
        West: {},
        South: {},
        East: {},
      },
      Compass: {
        onClick: () => setViewGuide((v) => !v),
      },
      visibleDirection: viewGuide,
    },
  };
  return <PreparingPage {...props} />;
}
