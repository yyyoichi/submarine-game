// import type { Route } from "./+types/";

// export async function clientLoader({ params }: Route.) {
//   const { gameId, playerId } = params;
//   try {
//     const logs = await battleClient.logs(
//       new LogsRequest({
//         gameId,
//         playerId,
//       }),
//     );
//     return logs;
//   } catch (err) {
//     const connectErr = new ConnectError(err as string);
//     console.error(connectErr.message);
//   }
// }

// export async function action({ request, params }: ActionFunctionArgs) {
//   const formData = await request.formData();
//   const { gameId, playerId } = params;

//   try {
//     switch (request.method) {
//       case "POST": {
//         switch (formData.get("type")?.toString()) {
//           case "first": {
//             const at = formData.get("at")?.toString();
//             const [mine1, mine2] = formData
//               .get("mines")
//               ?.toString()
//               .split(",") || ["0", "0"];
//             const req = new DeployRequest({
//               gameId,
//               playerId,
//               at: Number(at),
//               mines: [Number(mine1), Number(mine2)],
//             });
//             await battleClient.deploy(req, { signal: request.signal });
//             break;
//           }
//           case "action": {
//             const at = formData.get("at")?.toString();
//             const strActionType = formData.get("act")?.toString();
//             const actionType = Number(strActionType);
//             const req = new ActionRequest({
//               type: actionType,
//               gameId,
//               playerId,
//               at: Number(at),
//             });
//             await battleClient.action(req, { signal: request.signal });
//             break;
//           }
//         }
//         break;
//       }
//       case "PATCH": {
//         const req = new WaitRequest({
//           gameId,
//           playerId,
//         });
//         for await (const _ of battleClient.wait(req, {
//           signal: request.signal,
//         })) {
//         }
//         break;
//       }

//       case "DELETE": {
//         break;
//       }
//     }
//   } catch (e) {
//     if (e instanceof ConnectError) {
//       console.error(e.message);
//     } else if (e instanceof Error) {
//       const ce = new ConnectError(e.message);
//       console.error(ce.message);
//     } else {
//       console.error(e);
//     }
//   }
//   return null;
// }
