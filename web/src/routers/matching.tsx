import { useState } from "react";
import { getGameClient } from "../api/connect";
import { JoinRequest } from "../gen/api/v1/game_pb";
import { Link } from "react-router-dom";

function Home() {
	const [isLoading, setIsLoading] = useState(false);
	const [url, setUrl] = useState("");

	const join = async () => {
		setIsLoading(true);
		try {
			const client = getGameClient();
			const stream = client.join(new JoinRequest());
			for await (const resp of stream) {
				console.log(resp);
				if (resp.gameId === "") {
					continue;
				}
				setUrl(`/playground/${resp.gameId}/${resp.userId}`);
			}
		} catch (e) {
			console.log(e);
		}
		setIsLoading(false);
	};
	return (
		<>
			<div>
				{!isLoading && url === "" ? (
					<h2 onKeyDown={join} onClick={join}>
						START JOIN
					</h2>
				) : isLoading ? (
					<div>wait...</div>
				) : (
					<Link to={url}>START</Link>
				)}
			</div>
			{/* <div>{url === "" ? <div>wait..</div> : <Link to={url}>START</Link>}</div> */}
		</>
	);
}

export default Home;
