import {
	Form,
	useLoaderData,
	type LoaderFunctionArgs,
	type ActionFunctionArgs,
} from "react-router-dom";
import { getGameClient } from "../api/connect";
import {
	ActionRequest,
	ActionType,
	CampStatus,
	HistoryRequest,
	type HistoryResponse,
} from "../gen/api/v1/game_pb";
import { ConnectError } from "@connectrpc/connect";
import { useState } from "react";

function Home() {
	const history = useLoaderData() as HistoryResponse;
	const calcCamp = (row: number, col: number) => {
		const lineSize = history.camps.length;
		return row * lineSize + col;
	};

	const [clickCamp, setClickCamp] = useState<number | null>();
	let enableStatus: CampStatus[] = [];
	for (let i = 0; i < history.camps.length; i++) {
		for (let j = 0; j < history.camps[i].camps.length; j++) {
			const c = history.camps[i].camps[j].camp;
			if (clickCamp === c) {
				enableStatus = history.camps[i].camps[j].status;
			}
		}
	}
	return (
		<>
			<h2>{history.description}</h2>
			<Form method="post">
				<table>
					<tbody>
						{history.camps.map((line, row) => (
							// biome-ignore lint/suspicious/noArrayIndexKey: <explanation>
							<tr key={row}>
								{line.camps.map((s, col) => (
									// biome-ignore lint/suspicious/noArrayIndexKey: <explanation>
									<td key={col}>
										{!s.status.includes(CampStatus.MOVE) &&
										!s.status.includes(CampStatus.BOMB) &&
										!s.status.includes(CampStatus.PLACE) ? (
											<>{calcCamp(row, col)}</>
										) : (
											<label>
												<input
													type="radio"
													name="camp"
													value={s.camp}
													checked={clickCamp === s.camp}
													onChange={() => {
														setClickCamp(s.camp);
													}}
												/>
												{s.camp}
											</label>
										)}
									</td>
								))}
							</tr>
						))}
					</tbody>
				</table>
				<div>
					{enableStatus.includes(CampStatus.PLACE) && (
						<label>
							<input type="radio" name="action" value={ActionType.PLACE} />
							{"配置"}
						</label>
					)}
					{enableStatus.includes(CampStatus.MOVE) && (
						<label>
							<input type="radio" name="action" value={ActionType.MOVE} />
							{"移動"}
						</label>
					)}
					{enableStatus.includes(CampStatus.BOMB) && (
						<label>
							<input type="radio" name="action" value={ActionType.BOMB} />
							{"魚雷"}
						</label>
					)}
				</div>
				<div>
					<label>
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
	try {
		const history = await client.history(
			new HistoryRequest({
				gameId: gameId ?? "",
				userId: userId ?? "",
			}),
		);
		console.log(history);
		return history;
	} catch (err) {
		const connectErr = new ConnectError(err as string);
		throw connectErr.message;
	}
}

export async function action({ request, params }: ActionFunctionArgs) {
	const formData = await request.formData();
	const { gameId, userId } = params;
	let actionType = ActionType.UNSPECIFIED;
	const strActionType = formData.get("action")?.toString();
	if (strActionType === "1") {
		actionType = ActionType.MOVE;
	} else if (strActionType === "2") {
		actionType = ActionType.BOMB;
	} else if (strActionType === "4") {
		actionType = ActionType.PLACE;
	}
	let camp = 99999;
	const strCamp = formData.get("camp")?.toString();
	try {
		camp = Number(strCamp);
	} catch (e) {
		console.error(e);
	}
	try {
		const clinet = getGameClient();
		const req = new ActionRequest({
			type: actionType,
			gameId,
			userId,
			camp: camp,
		});
		console.log(req);
		await clinet.action(req);
	} catch (e) {
		if (e instanceof ConnectError) {
			window.alert(e.message);
		} else if (e instanceof Error) {
			const ce = new ConnectError(e.message);
			window.alert(ce.message);
		} else {
			window.alert(e);
		}
	}
	return null;
}

export default Home;
