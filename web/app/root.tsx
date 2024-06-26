import { useState } from "react";
import { createBrowserClient } from "@supabase/ssr";
import {
  Outlet,
  isRouteErrorResponse,
  useLoaderData,
  useRouteError,
} from "@remix-run/react";
import {
  LinksFunction,
  LoaderFunctionArgs,
  MetaFunction,
  defer,
} from "@remix-run/node";
import { authBegin } from "@/services/auth.server";
import { createSupabaseServerClient } from "@/services/supabase.server";
import { metaDescription } from "@/services/constants";
import { themeSessionResolver } from "@/components/theme/theme.server";
import { NavBarWrapper } from "@/components/base/nav-bar";
import { NavUser } from "@/components/base/nav-user";
import { BaseLayout } from "@/components/base/base-layout";
import { PageLoading } from "@/components/base/page-loading";
import favicon from "@/assets/images/logo-new.svg?url";
import styles from "@/assets/styles/global.css?url";

export const meta: MetaFunction = () => {
  return [
    { title: "ZByAI - AI-Enhanced Search" },
    { name: "description", content: metaDescription },
    { property: "og:title", content: "ZByAI - AI-Enhanced Search" },
    { property: "og:description", content: metaDescription },
    { property: "og:image", content: "https://www.zbyai.com/og-zbyai.jpg" },
    { property: "og:type", content: "website" },
    { property: "og:url", content: "https://www.zbyai.com/" },
    { property: "og:site_name", content: "ZByAI - AI-Enhanced Search" },
    { property: "twitter:title", content: "ZByAI - AI-Enhanced Search" },
    { property: "twitter:description", content: metaDescription },
    {
      property: "twitter:image",
      content: "https://www.zbyai.com/og-zbyai.jpg",
    },
    { property: "twitter:card", content: "summary_large_image" },
    { property: "twitter:url", content: "https://www.zbyai.com/" },
  ];
};

export async function loader({ request }: LoaderFunctionArgs) {
  const { supabaseClient, headers, supabaseEnv } =
    createSupabaseServerClient(request);

  const authHeader = await authBegin(request);
  if (authHeader) {
    Object.entries(authHeader).forEach(([key, value]) => {
      headers.append(key, value);
    });
  }

  const session = supabaseClient.auth
    .getSession()
    .then(({ data }) => data.session)
    .catch(() => null);

  const theme = await themeSessionResolver(request)
    .then(({ getTheme }) => getTheme())
    .catch(() => null);

  return defer({ theme, env: supabaseEnv, session }, { headers });
}

export const links: LinksFunction = () => [
  { rel: "icon", type: "image/svg+xml", href: favicon },
  { rel: "canonical", href: "https://www.zbyai.com/" },
  { rel: "stylesheet", href: styles },
  {
    rel: "search",
    type: "application/opensearchdescription+xml",
    title: "ZByAI",
    href: "/opensearch.xml",
  },
];

export default function App() {
  const { theme, env, session } = useLoaderData<typeof loader>() || {};

  const [supabase] = useState(() =>
    createBrowserClient(env.SUPABASE_URL, env.SUPABASE_ANON_KEY)
  );

  return (
    <BaseLayout theme={theme}>
      <div className="relative flex min-h-screen min-h-dvh flex-col bg-background">
        <NavBarWrapper>
          <NavUser supabase={supabase} session={session} />
        </NavBarWrapper>

        <Outlet context={{ supabase }} />
        <PageLoading />
      </div>
    </BaseLayout>
  );
}

export function ErrorBoundary() {
  const error = useRouteError();

  return (
    <BaseLayout theme={null}>
      <div className="relative flex min-h-screen min-h-dvh flex-col bg-background">
        <NavBarWrapper />

        {isRouteErrorResponse(error) ? (
          <div className="container">
            <h1>
              {error.status} {error.statusText}
            </h1>
            <p>{error.data}</p>
          </div>
        ) : error instanceof Error ? (
          <div className="container">
            <h1>Error</h1>
            <p>{error.message}</p>
            <p>The stack trace is:</p>
            <pre>{error.stack}</pre>
          </div>
        ) : (
          <div className="container">
            <h1>Unknown Error</h1>
          </div>
        )}
      </div>
    </BaseLayout>
  );
}
