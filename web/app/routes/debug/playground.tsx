import type React from "react";
import { PreparingPage } from "../playgrounds/component";

export default function Playground() {
  const props: React.ComponentProps<typeof PreparingPage> = {
    Ocean: {
      displaySectorName: true,
      sectorCount: 6,
      Sectors: {
        1: { color: "embed" },
        16: { color: "embed" },
      },
    },
  };
  return <PreparingPage {...props} />;
}
