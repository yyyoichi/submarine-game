import {
	type ActionFunctionArgs,
	Form,
	useLoaderData,
	type LoaderFunctionArgs,
} from "react-router-dom";
import { getGameClient } from "../api/connect";
import {
	ActionRequest,
	ActionType,
	HistoryRequest,
	type HistoryResponse,
} from "../gen/api/v1/game_pb";

function Home() {
	const history = useLoaderData() as HistoryResponse;
	console.log(history);
	return (
		<>
			<h2>{history.description}</h2>
			<Form method="post">
				<table>
					<tbody>
						{new Array(6).fill("").map((_, i) => {
							return (
								// biome-ignore lint/suspicious/noArrayIndexKey: <explanation>
								<tr key={i}>
									{new Array(6).fill("").map((_, j) => {
										const n = i * 6 + j;
										const isIsland = history.island.includes(BigInt(n));
										return (
											<td key={n}>
												{isIsland ? (
													<>■</>
												) : (
													<input type="radio" name="camp" value={i} />
												)}
											</td>
										);
									})}
								</tr>
							);
						})}
					</tbody>
				</table>
				<div>
					<label>
						<input type="radio" name="action" value={ActionType.MOVE} />
						移動
					</label>
					<label>
						<input type="radio" name="action" value={ActionType.BOMB} />
						魚雷
					</label>
				</div>
				<div>
					<label>
						<input type="hidden" name="type" value="action" />
						<input type="submit" value={"行動"} />
					</label>
				</div>
			</Form>
		</>
	);
}

export async function loader({ params }: LoaderFunctionArgs) {
	const { gameId, userId } = params;
	const client = getGameClient();
	const history = await client.history(
		new HistoryRequest({
			gameId: gameId ?? "",
			userId: userId ?? "",
		}),
	);
	return history;
}

export const action = async ({ request, params }: ActionFunctionArgs) => {
	const formData = await request.formData();
	const { gameId, userId } = params;
	const strActionType = formData.get("action")?.toString();
	const strCamp = formData.get("camp")?.toString();
	try {
		const clinet = getGameClient();
		const req = new ActionRequest({
			type:
				strActionType === String(ActionType.BOMB)
					? ActionType.BOMB
					: ActionType.MOVE,
			camp: Number(strCamp || -1),
			gameId,
			userId,
		});
		console.log(req);
		await clinet.action(req);
	} catch (e) {
		throw new Error(e as string);
	}
};

export default Home;
