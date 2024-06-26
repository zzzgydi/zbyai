import { ActionFunctionArgs, json, redirect } from "@remix-run/node";
import { ApiServer } from "~/services/api.server";

// Append Thread
export async function action({ request }: ActionFunctionArgs) {
  const { id, query } = await request.json();
  const api = new ApiServer(request);
  try {
    const result = await api.appendThread(id, query);
    // maybe fork thread
    if (result.id !== id) {
      return redirect(`/thread/${result.id}`, {
        headers: api.responseHeaders(),
      });
    }
    return json(result, { status: 200, headers: api.responseHeaders() });
  } catch (error: any) {
    console.error(`appendThread error: ${error}`);
    return json(
      { error: error.message || error.toString() },
      { status: 200, headers: api.responseHeaders() }
    );
  }
}
