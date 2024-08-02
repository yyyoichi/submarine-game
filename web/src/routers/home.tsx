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
        <h2>{"遊び方"}</h2>
        <div>
          <p>
            {
              "ゲーム説明:「潜水艦ゲーム」は、2人で対戦する戦略ボードゲームです。"
            }
          </p>
          <p>
            {
              "6x6のボード上で、潜水艦を操作して相手の潜水艦を魚雷で撃沈することを目指します。"
            }
          </p>
        </div>
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
          <li>{"魚雷攻撃"}</li>
          <ul>
            <li>{"魚雷攻撃の結果に、基づいて次の情報が得られます。"}</li>
            <li>
              {
                "面舵一杯！: 自分の魚雷攻撃の上下左右斜めの範囲に相手が潜行しています。"
              }
            </li>
            <li>
              {
                "ヨーソロー！: 少なくとも自分の魚雷攻撃の上下左右斜めの範囲に相手は潜行していません。"
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
