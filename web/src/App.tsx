import { useState } from "react";
import reactLogo from "./assets/react.svg";
import viteLogo from "/vite.svg";
import "./App.css";
import { createConnectTransport } from "@connectrpc/connect-web";
import { HelloService } from "./gen/api/v1/game_connect";
import { createPromiseClient } from "@connectrpc/connect";

console.log(`${window.location.origin}/rpc`);
const transport = createConnectTransport({
	baseUrl: `${window.location.origin}/rpc`,
});
const client = createPromiseClient(HelloService, transport);

function App() {
	const [inputValue, setInputValue] = useState("");

	return (
		<>
			<div>
				<a href="https://vitejs.dev" target="_blank" rel="noreferrer">
					<img src={viteLogo} className="logo" alt="Vite logo" />
				</a>
				<a href="https://react.dev" target="_blank" rel="noreferrer">
					<img src={reactLogo} className="logo react" alt="React logo" />
				</a>
			</div>
			<h1>Vite + React</h1>
			<div className="card">
				<form
					onSubmit={async (e) => {
						e.preventDefault();
						const resp = await client.say({
							name: inputValue,
						});
						window.alert(resp.hello);
					}}
				>
					<input
						value={inputValue}
						onChange={(e) => setInputValue(e.target.value)}
					/>
					<button type="submit">Send</button>
				</form>
			</div>
		</>
	);
}

export default App;
