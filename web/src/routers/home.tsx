import { useEffect, useState } from "react";
import { Form, redirect } from "react-router-dom";
import { getGameClient, leaveEffect } from "../api/connect";
import { JoinRequest } from "../gen/api/v1/game_pb";
import { ConnectError } from "@connectrpc/connect";

function Home() {
  const [isLoading, setIsLoading] = useState(false);
  useEffect(leaveEffect, []);
  return (
    <>
      <h1>🪖潜水艦ゲーム</h1>
      <Form
        method="post"
        onSubmit={() => {
          setIsLoading(true);
        }}
      >
        <label>
          {isLoading ? (
            <div>通信相手探し中...👀</div>
          ) : (
            <input type="submit" value={"START GAME"} />
          )}
        </label>
      </Form>
      <div>
        <div>{"ゲーム説明:「潜水艦ゲーム」は、2人で対戦する戦略ボードゲームです。6x6のボード上で、潜水艦を操作して相手の潜水艦を魚雷で撃沈することを目指します。"}</div>
        <h2>{"遊び方"}</h2>
        <ol>
          <li>{"基本ルール"}</li>
          <ul>
            <li>
              {"プレイヤーはターン制で、「行動」か「魚雷攻撃」を選択します。"}
            </li>
            <li>
              {"行動: 潜水艦を上下左右に一マス移動させることができます。"}
            </li>
            <li>
              {
                "魚雷攻撃:同じ位置にとどまり、隣接する上下左右斜めのマスに魚雷を撃ちます。"
              }
            </li>
          </ul>
          <li>{"行動履歴"}</li>
          <ul>
            <li>
              {
                "プレイヤーの動きは「行動履歴」として記録され、相手に表示されます。"
              }
            </li>
            <li>{"例: 相手が左に動いた場合、「西進」と表示されます。"}</li>
            <li>
              {
                "例:上に動くと「北進」、右に動くと「東進」、下に動くと「南進」と表示されます。"
              }
            </li>
          </ul>
          <li>{"行動結果"}</li>
          <ul>
            <li>{"魚雷攻撃の結果に基づいて、次の情報が得られます。"}</li>
            <li>
              {"相手が魚雷を撃ったボード上に場所（「海域」）が分かります。"}
            </li>
            <li>
              {
                "自分の魚雷攻撃が相手の上下左右斜めの範囲に入った場合は「面舵一杯！」、それ以外は「ヨーソロー！」と表示されます。"
              }
            </li>
          </ul>
        </ol>
      </div>
    </>
  );
}

export async function action() {
  try {
    const client = getGameClient();
    const stream = client.join(new JoinRequest());
    for await (const resp of stream) {
      console.log(resp);
      if (resp.gameId === "") {
        continue;
      }
      return redirect(`/playground/${resp.gameId}/${resp.userId}`);
    }
  } catch (e) {
    if (e instanceof ConnectError) {
      console.error(e.message);
    } else if (e instanceof Error) {
      const ce = new ConnectError(e.message);
      console.error(ce.message);
    } else {
      console.error(e);
    }
  }
  return null;
}

export default Home;
