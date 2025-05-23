import type React from "react";
import { useEffect, useState } from "react";
import { PreparingPage } from "../playgrounds/component";

export default function Playground() {
  const [viewGuide, setViewGuide] = useState(false);
  const [message, setMessage] = useState("敵魚雷A1に着弾！接近しています。");
  useEffect(() => {
    const interval = setTimeout(() => {
      setMessage("魚雷B1に着弾！付近から反響音！");
    }, 7000);
    return () => {
      clearInterval(interval);
    };
  }, []);
  const props: React.ComponentProps<typeof PreparingPage> = {
    Leading: {
      inMyTrun: true,
      Alerm: {
        useAlerm: true,
        startPingSec: 30,
        finishDatetime: new Date(Date.now() + 1000 * 30),
      },
      GuideLine: {
        children: message,
      },
    },
    OverlayedLeading: {
      fadeout: true,
      GuideLine: {
        children: "魚雷A1に着弾！接近しています。",
      },
    },
    Ocean: {
      displaySectorName: viewGuide,
      sectorCount: 6,
      Sectors: {
        1: { embed: true },
        16: { embed: true },
      },
    },
    OverlayedOcean: {
      fadeout: true,
      gapSector: 17,
      ringSector: 17,
      sectorCount: 6,
      title: "魚雷発射",
      titlePosition: "right",
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
