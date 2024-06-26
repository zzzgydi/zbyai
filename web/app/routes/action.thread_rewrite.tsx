import { ActionFunctionArgs, json } from "@remix-run/node";
import { ApiServer } from "~/services/api.server";

// Rewrite Thread
export async function action({ request }: ActionFunctionArgs) {
  const { id, runId } = await request.json();
  const api = new ApiServer(request);
  try {
    const result = await api.rewriteThread(id, runId);
    return json(result, { status: 200, headers: api.responseHeaders() });
  } catch (error: any) {
    console.error(`rewriteThread error: ${error}`);
    return json(
      { error: error.message || error.toString() },
      { status: 200, headers: api.responseHeaders() }
    );
  }
}
