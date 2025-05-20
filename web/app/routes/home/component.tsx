import type React from "react";

import { Button } from "@/components/ui/button";
import { ChevronsRight, Loader2 } from "lucide-react";

type Props = {
  StartButton: StartButtonProps;
};
export const HomePage = (props: Props) => {
  return (
    <div className="flex flex-col">
      <h1 className="font-train-one pb-2 text-5xl font-extrabold">
        🪖潜水艦ゲーム
      </h1>
      <div className="flex justify-center py-6">
        <StartButton {...props.StartButton} />
      </div>
      <div className="px-1">
        <div className="border-b-2 my-5">
          <h2 className="font-train-one text-3xl pb-2 font-bold">遊び方</h2>
        </div>
        <Section>
          <p className="my-1">
            「潜水艦ゲーム」は、2人で通信対戦する戦略ボードゲームです。
          </p>
          <p className="my-1">
            6x6のボード上で、潜水艦を操作して相手の潜水艦を魚雷で撃沈することを目指します。
          </p>
        </Section>
        <Section>
          <ul className="list-disc mx-4">
            <h3 className="font-train-one text-2xl pb-2 font-bold">
              基本ルール
            </h3>
            <HowToList span="潜航">
              潜水艦を上下左右に一マス移動させることができます。
            </HowToList>
            <HowToList span="魚雷攻撃">
              同じ位置にとどまり、隣接する上下左右斜めのマスに魚雷を撃ちます。
            </HowToList>
            <HowToList span="機雷作動">
              同じ位置にとどまり、あらかじめ敷設した機雷を作動させます。
            </HowToList>
          </ul>
        </Section>
        <Section>
          <ul className="list-disc mx-4">
            <h3 className="font-train-one text-2xl pb-2 font-bold">
              勝敗の分かれ目
            </h3>
            <HowToList span="面舵一杯！">
              魚雷・機雷攻撃の上下左右斜めの範囲に相手が潜行しています。
            </HowToList>
            <HowToList span="ヨーソロー！">
              少なくとも魚雷・機雷攻撃の上下左右斜めの範囲に相手は潜行していません。
            </HowToList>
          </ul>
        </Section>
      </div>
    </div>
  );
};

type StartButtonProps = Pick<
  React.ComponentProps<typeof Button>,
  "onClick" | "disabled"
>;

const StartButton = ({ disabled, onClick }: StartButtonProps) => {
  return (
    <Button disabled={disabled} onClick={onClick}>
      {disabled ? (
        <span className="flex items-center gap-1">
          <Loader2 className="animate-spin" />
          {"対戦相手を待っています..."}
        </span>
      ) : (
        <span className="flex items-center gap-1">
          {"開始"}
          <ChevronsRight />
        </span>
      )}
    </Button>
  );
};

const Section = (props: React.ComponentProps<"section">) => (
  <section {...props} className="border-b-2 my-5 pb-5" />
);

type HowToListProps = {
  span: React.ReactNode;
  children: React.ReactNode;
};
const HowToList = ({ span, children }: HowToListProps) => {
  return (
    <li className="my-1">
      <span className="border-b-1">{span}</span>
      {": "}
      {children}
    </li>
  );
};
