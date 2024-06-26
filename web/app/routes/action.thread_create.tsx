import type { ActionFunctionArgs } from "@remix-run/node";
import { json, redirect } from "@remix-run/react";
import { ApiServer } from "~/services/api.server";

export async function action({ request }: ActionFunctionArgs) {
  const { query } = await request.json();
  const api = new ApiServer(request);
  try {
    const { id } = await api.createThread(query);
    const queryEnc = encodeURIComponent(query);

    return redirect(`/thread/${id}?query=${queryEnc}`, {
      headers: api.responseHeaders(),
    });
  } catch (error: any) {
    console.error(`createThread error: ${error}`);
    return json(
      { error: error.message || error.toString() },
      { status: 200, headers: api.responseHeaders() }
    );
  }
}
