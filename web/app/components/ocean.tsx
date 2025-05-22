import { sector } from "@/lib/utils";
import type { ClassValue } from "clsx";
import { cn } from "src/lib/utils";

type Props = {
  displaySectorName: boolean;
  sectorCount: number; // 1辺のセクター数
  Sectors: Record<number, { embed?: boolean }>;
};

export const OceanComponent = ({
  Sectors,
  displaySectorName,
  sectorCount,
}: Props) => {
  return (
    <div className="grid grid-cols-6 py-2 px-3 gap-2 aspect-square bg-muted-foreground rounded-xs">
      {new Array(sectorCount * sectorCount).fill("").map((_, i) => {
        const sname = sector(sectorCount, i);
        const p: OceanSectorProps = {
          type: Sectors[i]?.embed ? "embed" : undefined,
          sectorName: displaySectorName ? sname : undefined,
        };
        return <OceanSector key={sname} {...p} />;
      })}
    </div>
  );
};

type OverlayProps = {
  gapSector?: number; // 唯一空けるセクター
  sectorCount: number; // 1辺のセクター数
} & (
  | {
      title?: undefined;
      titlePosition?: undefined;
    }
  | {
      title: string;
      titlePosition: "right" | "left";
    }
);

export const OverlayedOceanComponent = ({
  sectorCount,
  gapSector,
  ...props
}: OverlayProps) => {
  const titleCn: ClassValue[] = [];
  if (props.title) {
    const topOrBottom =
      (gapSector || 0) < sectorCount ** 2 / 2 ? "bottom" : "top";
    if (props.titlePosition === "right") {
      titleCn.push("justify-end");
    }
    if (topOrBottom === "bottom") {
      titleCn.push("items-end");
    }
  }
  return (
    <div className="grid grid-cols-6 py-2 px-3 gap-2 aspect-square bg-foreground absolute w-full top-0 left-0 z-100 animate-fadeout">
      {new Array(sectorCount * sectorCount).fill("").map((_, i) => {
        const sname = sector(sectorCount, i);
        const p: OceanSectorProps = {
          type: gapSector !== i ? "overlay" : undefined,
        };
        return <OceanSector key={sname} {...p} />;
      })}
      {props.title && (
        <div
          className={cn(
            "absolute w-full h-full top-0 left-0 flex text-background text-4xl",
            ...titleCn,
          )}
        >
          {props.title}
        </div>
      )}
    </div>
  );
};

type OceanSectorProps = {
  type?: "embed" | "overlay";
  sectorName?: string;
};

const OceanSector = (props: OceanSectorProps) => {
  const addedClass: ClassValue[] = [];
  switch (props.type) {
    case "embed": {
      //  枠を残して背景を薄めで潰す
      addedClass.push("border-backgraund", "bg-muted-foreground");
      break;
    }
    case "overlay":
      // 枠を残さず背景を濃い色で潰す
      addedClass.push("border-foreground", "bg-foreground");
      break;
    default:
      // 枠を背景をと同じ色
      addedClass.push("border-backgraund", "bg-background");
      break;
  }
  return (
    <div
      className={cn(
        "w-full h-full relative border-1 rounted-xs",
        ...addedClass,
      )}
    >
      {/* iconとセクター位置補助をレイヤーする */}
      <div className={"absolete w-full h-full text-muted-foreground z-10"}>
        <div className=" flex justify-center items-center">{""}</div>
      </div>
      {props.sectorName && (
        <div
          className={
            "absolute top-0 w-full h-full flex justify-center items-center text-muted-foreground z-0"
          }
        >
          {props.sectorName}
        </div>
      )}
    </div>
  );
};
