import { sector } from "@/lib/utils";
import { cn } from "src/lib/utils";

type Props = {
  displaySectorName: boolean;
  sectorCount: number; // 1辺のセクター数
  Sectors: Record<number, Omit<OceanSectorProps, "sectorName">>;
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
          ...Sectors[i],
          sectorName: displaySectorName ? sname : undefined,
        };
        return <OceanSector key={sname} {...p} />;
      })}
    </div>
  );
};

type OceanSectorProps = {
  color?: "embed";
  sectorName?: string;
};

const OceanSector = ({ color, ...props }: OceanSectorProps) => {
  const colorClass =
    color === "embed" ? "bg-muted-foreground" : "bg-background";
  return (
    <div
      className={cn(
        "w-full h-full relative border-1 border-backgraund rounted-xs",
        colorClass,
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
