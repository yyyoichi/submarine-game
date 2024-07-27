import { useLoaderData } from "react-router-dom";
import type { History, HistoryResponse } from "../../gen/api/v1/game_pb";

export function HistoryComponent() {
  const history = useLoaderData() as HistoryResponse;

  const iamTheFirst = history.myTurn === (history.histories.length % 2 === 0);
  const histories = history.histories.map((x) => {
    if (x.userId === "") {
      x.impact = "";
    }
    return x;
  });
  const firstHistories = histories
    .filter((x) => x.turn % 2 === 0)
    .sort((a, b) => a.turn - b.turn);
  const secondHistories = histories
    .filter((x) => x.turn % 2 === 1)
    .sort((a, b) => a.turn - b.turn);

  return (
    <div style={{ display: "flex", gap: 10 }}>
      <div>
        <h3>
          {"先攻"}
          {iamTheFirst && "📍"}
        </h3>
        {firstHistories.map((x) => (
          <ActionComponent x={x} key={x.turn} />
        ))}
      </div>
      <div>
        <h3>
          {"後攻"}
          {!iamTheFirst && "📍"}
        </h3>
        {secondHistories.map((x) => (
          <ActionComponent x={x} key={x.turn} />
        ))}
      </div>
    </div>
  );
}

function ActionComponent({ x }: { x: History }) {
  return (
    <div>
      <p>{`${x.description}${x.impact && ` >> ${x.impact}`}`}</p>
    </div>
  );
}
