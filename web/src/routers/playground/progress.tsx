import { useState, useEffect } from "react";
import { useLoaderData } from "react-router-dom";
import type { HistoryResponse } from "../../gen/api/v1/game_pb";
import { Progress } from "@chakra-ui/react";

type ProgressBarProps = {
  callback: () => void;
};
export function ProgressBar(props: ProgressBarProps) {
  const history = useLoaderData() as HistoryResponse;
  const [timeLeft, setTimeLeft] = useState(
    Number(history.timeout) - Date.now(),
  );
  useEffect(() => {
    const interval = setInterval(() => {
      if (history.winner !== "") {
        return;
      }
      const currentTime = Number(history.timeout) - Date.now();
      setTimeLeft(currentTime);
      if (currentTime > 0) return;
      props.callback();
    }, 100); // 毎秒更新

    return () => clearInterval(interval);
  }, [history.timeout, history.winner, props.callback]);

  return (
    <Progress
      hasStripe
      value={0 < timeLeft ? (timeLeft / 1000 / 30) * 100 : 0}
    />
  );
}
