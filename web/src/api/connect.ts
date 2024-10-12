import { createClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import { MatchingService } from "../gen/api/v2/game_connect";

const transport = createConnectTransport({
  baseUrl: `${window.location.origin}/rpc`,
});

export const matchingClient = createClient(MatchingService, transport);
