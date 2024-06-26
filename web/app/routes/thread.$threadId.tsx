import { useEffect, useRef, useState } from "react";
import { LoaderFunctionArgs, MetaFunction, json } from "@remix-run/node";
import { Link, useLoaderData, useRevalidator } from "@remix-run/react";
import { ThreadInput, ThreadInputRef } from "~/components/thread/thread-input";
import { metaDescription } from "~/services/constants";
import { ThreadRunItem } from "~/components/thread/thread-run";
import { Separator } from "~/components/ui/separator";
import { ApiServer } from "~/services/api.server";
import { Button } from "~/components/ui/button";
import katexStyle from "katex/dist/katex.min.css?url";
import markdownStyle from "~/assets/styles/markdown.scss?url";
import highlightStyle from "~/assets/styles/highlight.scss?url";

export const loader = async ({ params, request }: LoaderFunctionArgs) => {
  const { threadId } = params;
  const url = new URL(request.url);
  let query = url.searchParams.get("query") || "";

  const api = new ApiServer(request);

  try {
    const detail = await api.detailThread(threadId!);
    if (!query) query = detail.title;
    return json(
      { threadId, query, detail, error: null },
      { headers: api.responseHeaders() }
    );
  } catch (error: any) {
    console.error(`detailThread error: ${error}`);
    return json(
      {
        threadId,
        query,
        detail: null,
        error: error.message || error.toString(),
      },
      { headers: api.responseHeaders() }
    );
  }
};

export const links = () => [
  { rel: "stylesheet", href: katexStyle },
  { rel: "stylesheet", href: markdownStyle },
  { rel: "stylesheet", href: highlightStyle },
];

export const meta: MetaFunction<typeof loader> = ({ data }) => {
  const shakeQuery = (query: string) => {
    const maxSize = 250;
    if (query.length > maxSize) return query.slice(0, maxSize) + "...";
    return query;
  };

  const title = data?.query
    ? `ZByAI - ${shakeQuery(data.query)}`
    : "ZByAI - AI-Enhanced Search";
  const url = data?.threadId
    ? `https://www.zbyai.com/thread/${data?.threadId}`
    : "https://www.zbyai.com/";
  return [
    { title },
    { name: "description", content: metaDescription },
    { property: "og:title", content: title },
    { property: "og:description", content: metaDescription },
    { property: "og:image", content: "https://www.zbyai.com/og-zbyai.jpg" },
    { property: "og:type", content: "website" },
    { property: "og:url", content: url },
    { property: "og:site_name", content: "ZByAI - AI-Enhanced Search" },
    { property: "twitter:title", content: title },
    { property: "twitter:description", content: metaDescription },
    {
      property: "twitter:image",
      content: "https://www.zbyai.com/og-zbyai.jpg",
    },
    { property: "twitter:card", content: "summary_large_image" },
    { property: "twitter:url", content: url },
  ];
};

export default function ThreadPage() {
  const { threadId, detail, error } = useLoaderData<typeof loader>() || {};
  const [loading, setLoading] = useState(false);

  const revalidator = useRevalidator();
  const inputRef = useRef<ThreadInputRef>(null);

  useEffect(() => {
    inputRef.current?.reset();
  }, [detail]);

  if (error != null) {
    return (
      <div className="container px-4 md:px-8 w-full mx-auto max-w-[980px] py-20 text-center">
        <h3 className="text-lg font-semibold text-destructive">
          Error: {error}
        </h3>

        <div className="mt-6 flex items-center gap-3 justify-center">
          <Button onClick={() => revalidator.revalidate()} variant="outline">
            Retry
          </Button>
          <Link to="/">
            <Button variant="secondary">Back</Button>
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className="relative container px-4 md:px-8 mx-auto w-full max-w-[980px] pb-[135px]">
      <div className="space-y-6">
        {detail?.history?.map((run, idx) => (
          <div key={run.id}>
            <ThreadRunItem
              loading={loading}
              threadId={threadId!}
              item={run}
              onRefresh={() => {
                setLoading(false);
                inputRef.current?.resetLoading();
                revalidator.revalidate();
              }}
              onStartLoading={() => {
                setLoading(true);
                inputRef.current?.startLoading();
              }}
              onDone={() => {
                setLoading(false);
                inputRef.current?.resetLoading();
              }}
            />
            {idx < detail?.history.length - 1 && <Separator className="mt-8" />}
          </div>
        ))}
      </div>

      <div className="fixed left-0 right-0 bottom-0 thread-mask" />
      <div className="fixed left-0 right-0 bottom-10">
        <div className="container px-4 md:px-8 w-full max-w-[980px]">
          <div className="px-1 md:px-3">
            <ThreadInput ref={inputRef} threadId={threadId!} />
          </div>
        </div>
      </div>
    </div>
  );
}
