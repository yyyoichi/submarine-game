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
	type History,
	type HistoryResponse,
} from "../gen/api/v1/game_pb";
import { ConnectError } from "@connectrpc/connect";
import { useState } from "react";

function Home() {
	const history = useLoaderData() as HistoryResponse;

	const iamTheFirst = history.myTurn === (history.histories.length % 2 === 1);
	const histories = history.histories.map((x) => {
		if (x.userId === "") {
			x.impact = "";
		}
		return x;
	});
	const firstHistories = histories
		.filter((x) => x.turn % 2 === 0)
		.sort((a, b) => a.turn - b.turn);
	const secondHistories = histories
		.filter((x) => x.turn % 2 === 1)
		.sort((a, b) => a.turn - b.turn);

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
								{line.camps.map((s, col) => {
									const Cell = () => {
										if (s.status.length === 0) {
											return <>🌊{s.camp}</>;
										}
										if (s.status.includes(CampStatus.ISLAND)) {
											return <>🏝️{s.camp}</>;
										}
										if (s.status.includes(CampStatus.SUBMARINE)) {
											return <>📍{s.camp}</>;
										}
										return (
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
										);
									};
									return (
										// biome-ignore lint/suspicious/noArrayIndexKey: <explanation>
										<td key={col}>
											<Cell />
										</td>
									);
								})}
							</tr>
						))}
					</tbody>
				</table>
				<div>
					{history.myTurn && enableStatus.includes(CampStatus.PLACE) && (
						<label>
							<input type="radio" name="action" value={ActionType.PLACE} />
							{"配置"}
						</label>
					)}
					{history.myTurn && enableStatus.includes(CampStatus.MOVE) && (
						<label>
							<input type="radio" name="action" value={ActionType.MOVE} />
							{"移動"}
						</label>
					)}
					{history.myTurn && enableStatus.includes(CampStatus.BOMB) && (
						<label>
							<input type="radio" name="action" value={ActionType.BOMB} />
							{"魚雷"}
						</label>
					)}
				</div>
				{history.myTurn && (
					<div>
						<label>
							<input type="submit" value={"行動"} />
						</label>
					</div>
				)}
			</Form>
			<div style={{ display: "flex", gap: 10 }}>
				<div>
					<h3>
						{"先攻"}
						{iamTheFirst && "📍"}
					</h3>
					{firstHistories.map((x) => (
						<ActionComponent x={x} key={x.turn} />
					))}
				</div>
				<div>
					<h3>
						{"後攻"}
						{!iamTheFirst && "📍"}
					</h3>
					{secondHistories.map((x) => (
						<ActionComponent x={x} key={x.turn} />
					))}
				</div>
			</div>
		</>
	);
}

function ActionComponent({ x }: { x: History }) {
	return (
		<div>
			<p>{`${x.description}${x.impact && ` >> ${x.impact}`}`}</p>
		</div>
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
