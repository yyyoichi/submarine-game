import { createClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import { BattleService, MatchingService } from "../gen/api/v2/game_connect";

const transport = createConnectTransport({
  baseUrl: `${window.location.origin}/rpc`,
});

export const matchingClient = createClient(MatchingService, transport);
export const battleClient = createClient(BattleService, transport);
