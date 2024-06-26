import { serialize } from "cookie";
import { ActionFunctionArgs, redirect } from "@remix-run/node";
import { COOKIE_SESSION } from "~/services/auth.server";
import { createSupabaseServerClient } from "~/services/supabase.server";

export const action = async ({ request }: ActionFunctionArgs) => {
  const { supabaseClient, headers } = createSupabaseServerClient(request);
  // check if user is logged in
  const {
    data: { session },
  } = await supabaseClient.auth.getSession();
  if (!session?.user) return;

  // sign out
  await supabaseClient.auth.signOut();

  headers.append(
    "Set-Cookie",
    serialize(COOKIE_SESSION, "", {
      path: "/",
      maxAge: 0,
      sameSite: "lax",
      httpOnly: true,
      expires: new Date(0),
    })
  );
  return redirect("/", { headers });
};
