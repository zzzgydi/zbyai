import { redirect, type LoaderFunctionArgs } from "@remix-run/node";
import { ApiServer } from "~/services/api.server";

export const loader = async ({ request }: LoaderFunctionArgs) => {
  const requestUrl = new URL(request.url);
  const query = requestUrl.searchParams.get("q");

  if (!query?.trim()) return redirect("/");

  const api = new ApiServer(request);
  const q = encodeURIComponent(query);

  try {
    const { id } = await api.createThread(query);

    return redirect(`/thread/${id}?query=${q}`, {
      headers: api.responseHeaders(),
    });
  } catch (error: any) {
    console.error(`createThread error: ${error}`);
    return redirect("/?q=" + q);
  }
};
