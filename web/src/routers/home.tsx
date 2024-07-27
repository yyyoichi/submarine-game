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
      <Form
        method="post"
        onSubmit={() => {
          setIsLoading(true);
        }}
      >
        <label>
          {isLoading ? (
            <div>....👀</div>
          ) : (
            <input type="submit" value={"START GAME"} />
          )}
        </label>
      </Form>
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
