import {
	Form,
	useLoaderData,
	type LoaderFunctionArgs,
	type ActionFunctionArgs,
	useSubmit,
} from "react-router-dom";
import { getGameClient } from "../../api/connect";
import {
	ActionRequest,
	ActionType,
	CampStatus,
	HistoryRequest,
	type HistoryResponse_Camp,
	type HistoryResponse,
	WaitRequest,
} from "../../gen/api/v1/game_pb";
import { ConnectError } from "@connectrpc/connect";
import type React from "react";
import { useEffect, useRef, useState } from "react";
import { HistoryComponent } from "./history";

function Home() {
	const history = useLoaderData() as HistoryResponse;

	const formRef = useRef<HTMLFormElement>(null);
	const submit = useSubmit();
	useEffect(() => {
		if (formRef.current && !history.myTurn) {
			submit(formRef.current, { method: "PATCH" });
		}
	}, [history.myTurn, submit]);

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
			<Form method="POST" ref={formRef}>
				<table>
					<tbody>
						{history.camps.map((line, row) => (
							// biome-ignore lint/suspicious/noArrayIndexKey: <explanation>
							<tr key={row}>
								{line.camps.map((camp, col) => (
									// biome-ignore lint/suspicious/noArrayIndexKey: <explanation>
									<td key={col}>
										<CampInput
											camp={camp}
											type="radio"
											name="camp"
											value={camp.camp}
											checked={clickCamp === camp.camp}
											onChange={() => {
												setClickCamp(camp.camp);
											}}
										/>
									</td>
								))}
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
			<HistoryComponent />
		</>
	);
}

function CampInput({
	camp,
	...props
}: { camp: HistoryResponse_Camp } & React.ComponentProps<"input">) {
	if (camp.status.length === 0) {
		return <>🌊{camp.camp}</>;
	}
	if (camp.status.includes(CampStatus.ISLAND)) {
		return <>🏝️{camp.camp}</>;
	}
	if (camp.status.includes(CampStatus.SUBMARINE)) {
		return <>📍{camp.camp}</>;
	}
	return (
		<label>
			<input {...props} />
			{camp.camp}
		</label>
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
	try {
		switch (request.method) {
			case "POST": {
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
				const clinet = getGameClient();
				const req = new ActionRequest({
					type: actionType,
					gameId,
					userId,
					camp: camp,
				});
				console.log(req);
				await clinet.action(req, { signal: request.signal });
				break;
			}

			case "PATCH": {
				const { gameId, userId } = params;
				const clinet = getGameClient();
				const req = new WaitRequest({
					gameId,
					userId,
				});
				console.log(req);
				for await (const _ of clinet.wait(req)) {
				}
				break;
			}
		}
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
