import { createConnectTransport } from "@connectrpc/connect-web";
import { GameService } from "../gen/api/v1/game_connect";
import { createPromiseClient } from "@connectrpc/connect";

const transport = createConnectTransport({
	baseUrl: `${window.location.origin}/rpc`,
});

export const getGameClient = () => createPromiseClient(GameService, transport);