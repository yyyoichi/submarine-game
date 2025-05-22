import { cn } from "@/lib/utils";
import type { ClassValue } from "clsx";
import { useEffect, useState } from "react";

type Props = {
  inMyTrun: boolean;
  Alerm: AlermProps;
};

export const LeadingComponent = (props: Props) => {
  const children = props.inMyTrun ? "あなたのターン" : "あいてのターン";
  return (
    <div className="w-full">
      <div className="w-full bg-foreground px-1 py-2 mb-5">
        <p className="text-background">text</p>
      </div>
      <div className="relative w-fit px-3">
        <p className="">{children}</p>
        <Alerm {...props.Alerm} />
      </div>
    </div>
  );
};

type AlermProps =
  | {
      useAlerm: true;
      startPingSec: number;
      finishDatetime: Date;
    }
  | {
      useAlerm: false;
    };

const Alerm = (props: AlermProps) => {
  const [intervaStage, setIntervaStage] = useState<0 | 1 | 2 | 3 | 4>(0);

  const addedCn: ClassValue[] = [];
  switch (intervaStage) {
    case 1:
      addedCn.push("animate-[ping_1s_ease-in-out_infinite]");
      break;
    case 2:
      addedCn.push("animate-[ping_500ms_ease-in-out_infinite]");
      break;
    case 3:
      addedCn.push("animate-[ping_300ms_ease-in-out_infinite]");
      break;
    case 4:
      addedCn.push("animate-[ping_300ms_ease-in-out_infinite]");
      addedCn.push("w-8 h-8");
      break;
  }

  useEffect(() => {
    if (!props.useAlerm) {
      return;
    }
    const interval = setInterval(() => {
      const now = new Date();
      const diff = props.finishDatetime.getTime() - now.getTime();
      const sec = Math.floor(diff / 1000);
      if (sec < 0) {
        return;
      }
      if (sec < 2) {
        setIntervaStage(4);
        return;
      }
      const p = sec / (props.startPingSec - 2);
      if (p < 0.2) {
        setIntervaStage(3);
        return;
      }
      if (p < 0.4) {
        setIntervaStage(2);
        return;
      }
      if (p < 1) {
        setIntervaStage(1);
        return;
      }
    }, 100);

    return () => clearInterval(interval);
  });

  if (!props.useAlerm) {
    return <></>;
  }
  return (
    <span className="absolute top-[.2rem] left-[.5rem] opacity-75">
      <span className="relative z-20">
        <span
          className={cn(
            "absolute inline-flex bg-red-500  w-4 h-4 rounded-full opacity-75 transform -translate-x-1/2 -translate-y-1/2 ",
            ...addedCn,
          )}
        />
        <span className="absolute bg-red-500 w-4 h-4 rounded-full opacity-75 transform -translate-x-1/2 -translate-y-1/2 " />
      </span>
    </span>
  );
};
