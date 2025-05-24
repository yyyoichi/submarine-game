import { cn } from "@/lib/utils";
import type { ClassValue } from "clsx";
import { useEffect, useId, useRef, useState } from "react";

type Props = {
  inMyTrun: boolean;
  Alerm: AlermProps;
} & GuideLineProps;

export const LeadingComponent = ({
  inMyTrun,
  Alerm: alermProps,
  ...props
}: Props) => {
  // Content hight = 22(5.5rem)
  return (
    <div className="w-full">
      {/* 計算の結果、GuideLineはh-11になるが、Overlayedとの整合性のため、明示する。 */}
      <div className="h-11">
        <GuideLine {...props} />
      </div>
      <div className="relative w-full">
        {/* ターンメッセージとアラーム */}
        {/* this content hight 11 = 2.75rem = 1.25rem top padding + 1.5rem p hight */}
        <div className="w-fit px-3 pt-5">
          <p className="">{inMyTrun ? "あなたのターン" : "あいてのターン"}</p>
          <div className="absolute top-[1.45rem] left-[.5rem] opacity-75">
            <div className="relative z-20">
              <Alerm {...alermProps} />
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

type OverlayedLeadingProps = {
  fadeout?: boolean;
  absolute?: boolean;
} & GuideLineProps;

export const OverlayedLeadingComponent = ({
  fadeout,
  absolute,
  ...props
}: OverlayedLeadingProps) => {
  return (
    <div
      className={cn(
        "w-full h-[5.5rem] bg-foreground",
        fadeout ? "animate-fadeout" : "",
        absolute ? "absolute z-100 top-0 left-0 " : "",
      )}
    >
      {props.children && (
        <div className="h-11">
          <GuideLine {...props} />
        </div>
      )}
    </div>
  );
};

type GuideLineProps = Pick<React.PropsWithChildren, "children"> & {
  contentWidth?: number;
};

const GuideLine = (props: GuideLineProps) => {
  const id = useId();
  const containerRef = useRef<HTMLDivElement>(null);
  const contentRef = useRef<HTMLDivElement>(null);
  const contentWidth =
    props.contentWidth || contentRef.current?.scrollWidth || 0;
  const [animationDuration, setAnimationDuration] = useState(0); // <= 0 でアニメーションしない
  const shouldAnimate = animationDuration > 0;
  const speed = 100; // px/s

  // biome-ignore lint/correctness/useExhaustiveDependencies: <explanation>
  useEffect(() => {
    if (!containerRef.current) {
      return;
    }

    const containerWidth = containerRef.current.clientWidth;

    // コンテンツがコンテナより広い場合のみアニメーション
    if (contentWidth > containerWidth) {
      // 速度（px/秒）に基づいて所要時間を計算
      // 全体の移動距離は (contentWidth+containerWidth)
      const duration = (contentWidth + containerWidth) / speed;
      setAnimationDuration(duration);
    } else {
      setAnimationDuration(0);
    }
  }, [containerRef.current?.clientWidth, contentWidth]);

  return (
    <div
      ref={containerRef}
      className="w-full bg-foreground px-1 py-2 h-[2.75rem] relative overflow-hidden"
    >
      <div
        ref={contentRef}
        className="text-background inline-block whitespace-nowrap text-lg"
        style={
          shouldAnimate
            ? {
                animation: `marquee-${id} ${animationDuration * 1000}ms linear infinite`,
              }
            : {}
        }
        {...props}
      />
      {shouldAnimate && (
        <style>{`
        @keyframes marquee-${id} {
          0% {
            transform: translateX(${containerRef.current?.clientWidth || 0}px);
          }
          100% {
            transform: translateX(-${contentWidth}px);
          }
        }
      `}</style>
      )}
    </div>
  );
};

type AlermProps =
  | {
      useAlerm: true;
      startPingMilliSec: number;
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
      const msec = Math.floor(diff);
      if (msec < 0) {
        return;
      }
      if (msec < 2000) {
        setIntervaStage(4);
        return;
      }
      const p = msec / (props.startPingMilliSec - 2000);
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
    <>
      <span
        className={cn(
          "absolute inline-flex bg-red-500  w-4 h-4 rounded-full opacity-75 transform -translate-x-1/2 -translate-y-1/2 ",
          ...addedCn,
        )}
      />
      <span className="absolute bg-red-500 w-4 h-4 rounded-full opacity-75 transform -translate-x-1/2 -translate-y-1/2 " />
    </>
  );
};
