import dayjs from "dayjs";
import relativeTime from "dayjs/plugin/relativeTime";
import { LoaderFunctionArgs, MetaFunction, json } from "@remix-run/node";
import { Link, useLoaderData } from "@remix-run/react";
import { ApiServer } from "~/services/api.server";
import { metaDescription } from "~/services/constants";

dayjs.extend(relativeTime);

export const loader = async ({ request }: LoaderFunctionArgs) => {
  const api = new ApiServer(request);

  try {
    const threads = await api.listThread();
    return json({ threads, error: null }, { headers: api.responseHeaders() });
  } catch (error: any) {
    console.error(`listThread error: ${error}`);
    return json(
      {
        threads: [] as Thread[],
        error: error.message || error.toString(),
      },
      { headers: api.responseHeaders() }
    );
  }
};

export const meta: MetaFunction<typeof loader> = ({ data }) => {
  return [
    { title: "ZByAI - Thread History" },
    { name: "description", content: metaDescription },
  ];
};

export default function ThreadHistoryPage() {
  const { threads, error } = useLoaderData<typeof loader>() || {};

  if (error) {
    return (
      <div className="relative container mx-auto max-w-[800px] pb-10">
        <p className="text-muted-foreground py-10 text-center">
          Sign in to view thread history...
        </p>
      </div>
    );
  }

  return (
    <div className="relative container mx-auto max-w-[980px] pb-10">
      <h2 className="text-xl font-semibold mb-6">Thread History</h2>
      <ul className="space-y-4">
        {threads?.map((thread) => (
          <li key={thread.id}>
            <Link
              to={`/thread/${thread.id}`}
              className="block rounded-lg border p-3 text-sm transition-all hover:bg-accent"
            >
              <div className="flex items-center justify-between gap-4">
                <h3 className="text-base font-semibold truncate">
                  {thread.title || "New Thread"}
                </h3>
                <p className="flex-none font-normal text-muted-foreground">
                  {dayjs(thread.created_at).fromNow()}
                </p>
              </div>
            </Link>
          </li>
        ))}
      </ul>
    </div>
  );
}
