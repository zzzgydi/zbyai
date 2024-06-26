import { createServerClient, parse, serialize } from "@supabase/ssr";

export const SUPABASE_URL = process.env.SUPABASE_URL!;
export const SUPABASE_ANON_KEY = process.env.SUPABASE_ANON_KEY!;

export const supabaseEnv = {
  SUPABASE_URL,
  SUPABASE_ANON_KEY,
};

export const createSupabaseServerClient = (request: Request) => {
  const cookies = parse(request.headers.get("Cookie") ?? "");
  const headers = new Headers();
  const supabaseClient = createServerClient(
    supabaseEnv.SUPABASE_URL!,
    supabaseEnv.SUPABASE_ANON_KEY!,
    {
      cookies: {
        get(key) {
          return cookies[key];
        },
        set(key, value, options) {
          headers.append(
            "Set-Cookie",
            serialize(key, value, { ...options, domain: ".zbyai.com" })
          );
        },
        remove(key, options) {
          headers.append(
            "Set-Cookie",
            serialize(key, "", { ...options, domain: ".zbyai.com" })
          );
        },
      },
    }
  );
  return { supabaseClient, headers, supabaseEnv };
};
