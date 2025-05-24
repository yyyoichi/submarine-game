import { useEffect, useId, useState } from "react";

type Props = {
  trigger: string;
  delay?: number; // delay for fadein
  duration?: number; // duration for showing children component
  fadein?: number; // duration for fadein
  fadeout?: number; // duration for fadeout
} & Pick<React.ComponentProps<"div">, "children" | "className">;

export const FadeInOutTrigger = (props: Props) => {
  const id = useId();
  const [viewable, setViewable] = useState(false);

  // default .. 遅延0で、100msでfadein、1000ms表示、300msでfadeout
  const delay = props.delay || 0;
  const fadein = props.fadein || 100;
  const duration = props.duration || 1000;
  const fadeout = props.fadeout || 300;
  // アニメーション時間
  const animationDuration = fadein + duration + fadeout;
  // 表示開始
  const percentageFadein = Math.floor((fadein / animationDuration) * 100);
  // 表示終了
  const percentageFadeout =
    100 - Math.floor((fadeout / animationDuration) * 100);

  // biome-ignore lint/correctness/useExhaustiveDependencies: <explanation>
  useEffect(() => {
    setViewable(true);
    const outTimeout = setTimeout(
      () => {
        setViewable(false); // fadeoutしたら要素削除
      },
      delay + animationDuration + 100,
    );
    return () => {
      clearTimeout(outTimeout);
    };
  }, [props.trigger]);

  if (!viewable) {
    return <></>;
  }
  return (
    <>
      <div
        className={props.className}
        style={{
          animation: `fadeinout-${id} ${Math.floor(animationDuration)}ms linear ${Math.floor(delay)}ms forwards`,
        }}
      >
        {props.children}
      </div>
      <style>
        {`
        @keyframes fadeinout-${id} {
          0%,100% {
            opacity: 0;
          }
          ${percentageFadein}%,${percentageFadeout}% {
            opacity: 1;
          }
        }
      `}
      </style>
    </>
  );
};
