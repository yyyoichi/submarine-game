import { matchingClient } from "@/lib/connectapi";
import { ConnectError } from "@connectrpc/connect";
import { redirect, useFetcher } from "react-router";
import type { Route } from "./+types/index";
import { HomePage } from "./component";

export default function Home() {
  const fetcher = useFetcher();
  const props: React.ComponentProps<typeof HomePage> = {
    StartButton: {
      disabled: fetcher.state !== "idle",
      onClick: () => {
        fetcher.submit({}, { method: "post" });
      },
    },
  };
  return <HomePage {...props} />;
}

export async function clientAction({ request }: Route.ClientActionArgs) {
  try {
    switch (request.method) {
      case "POST": {
        const stream = matchingClient.waitEnemy(
          {},
          {
            signal: request.signal,
          },
        );
        for await (const resp of stream) {
          if (resp.gameId) {
            return redirect(`/playgrounds/${resp.gameId}/${resp.playerId}`);
          }
        }
      }
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
