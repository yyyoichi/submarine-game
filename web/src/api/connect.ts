import { createConnectTransport } from "@connectrpc/connect-web";
import { GameService } from "../gen/api/v1/game_connect";
import { createPromiseClient } from "@connectrpc/connect";
import { LeaveRequest } from "../gen/api/v1/game_pb";

const transport = createConnectTransport({
	baseUrl: `${window.location.origin}/rpc`,
});

export const getGameClient = () => createPromiseClient(GameService, transport);

export const leaveEffect = () => {
	const leaveFromGame = async () => {
		const client = getGameClient();
		client.leave(new LeaveRequest());
	};
	window.addEventListener("beforeunload", leaveFromGame);
	return () => {
		window.removeEventListener("beforeunload", leaveFromGame);
	};
}