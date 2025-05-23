import type React from "react";
import { useEffect, useState } from "react";
import { PreparingPage } from "../playgrounds/component";

export default function Playground() {
  const [viewGuide, setViewGuide] = useState(false);
  const [step, setStep] = useState<1 | 2 | 3>(1);
  useEffect(() => {
    const interval = setInterval(() => {
      setStep((v) => {
        if (v === 3) {
          return 1;
        }
        return (v + 1) as 2 | 3;
      });
    }, 10000);
    return () => {
      clearInterval(interval);
    };
  }, []);
  const props: React.ComponentProps<typeof PreparingPage> = {
    // @ts-ignore
    preparingStep: ["me", "mines", "done"][step - 1],
    Leading: {
      Alerm: {
        useAlerm: true,
        startPingSec: 30,
        finishDatetime: new Date(Date.now() + 1000 * 30),
      },
    },
    Ocean: {
      sectorCount: 6,
      islands: [10, 28],
    },
    MeOcean: {
      sectors: [12],
    },
    MinesOcean: {
      sectors: [18, 32],
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
