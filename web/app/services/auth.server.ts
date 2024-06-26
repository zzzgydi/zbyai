import jwt from "jsonwebtoken";
import { parse } from "cookie";
import { createSupabaseServerClient } from "./supabase.server";
import { API_BASE_URL } from "./constants";

export const COOKIE_SESSION = "session_token";
const baseURL = API_BASE_URL;

interface AuthClaims {
  id: string;
  name: string;
  auth: 0 | 1 | 2;
  exp: number;
}

export const getAuthFromJWT = (token: string) => {
  if (!process.env.JWT_SECRET) throw new Error("JWT_SECRET is not set");
  if (!token) return null;
  try {
    const decoded = jwt.verify(token, process.env.JWT_SECRET!) as AuthClaims;
    // 没过期
    if (decoded.exp * 1000 >= Date.now()) {
      return decoded;
    }
  } catch (err: any) {
    console.error(`getAuthFromJWT error: ${err}`);
  }
  return null;
};

export const getAuthFromCookie = (headers: Headers) => {
  const cookies = parse(headers.get("Cookie") ?? "");
  return getAuthFromJWT(cookies[COOKIE_SESSION]);
};

export const authHandler = async (request: Request) => {
  // 如果有session，且没有过期，且是supabase的，就直接传
  // 是游客的话，再看看下面是不是有supabase
  const auth = getAuthFromCookie(request.headers);
  if (auth && auth.auth > 0) {
    return;
  }

  // 如果没有，就获取supabase的session，传过去
  const { supabaseClient } = createSupabaseServerClient(request);
  try {
    const {
      data: { session },
    } = await supabaseClient.auth.getSession();
    if (!session?.user) return;
    if (!session.access_token) return;

    request.headers.append("Authorization", `Bearer ${session.access_token}`);
  } catch (err: any) {
    console.error(`authHandler error: ${err}`);
  }
};

export const authBegin = async (request: Request) => {
  // 刚登录的时候，去向服务端获取session token
  const cookies = parse(request.headers.get("Cookie") ?? "");
  if (cookies[COOKIE_SESSION]) return;

  // 判断是否包含supabase的session
  // 没包含也没法获取，就算了
  if (!Object.keys(cookies).some((key) => key.includes("-auth-token."))) return;

  const { supabaseClient } = createSupabaseServerClient(request);
  try {
    const {
      data: { session },
    } = await supabaseClient.auth.getSession();
    if (!session?.user) return;
    if (!session.access_token) return;

    request.headers.append("Authorization", `Bearer ${session.access_token}`);
    const respHeaders = await authAPI(request.headers);
    return respHeaders;
  } catch (err: any) {
    console.error(`authHandler error: ${err}`);
  }
};

async function authAPI(headers: Headers): Promise<Record<string, any>> {
  const newHeaders = new Headers(headers);
  newHeaders.set("Content-Type", "application/json");
  newHeaders.delete("host"); // 这玩意儿会导致请求又转到当前服务上

  try {
    const resp = await fetch(`${baseURL}/inner/auth`, {
      method: "POST",
      headers: newHeaders,
      credentials: "include",
    });
    const traceId = resp.headers.get("X-Trace-Id") || "";
    const setCookie = resp.headers.get("Set-Cookie");
    if (!setCookie) return { "X-Trace-Id": traceId };
    return { "Set-Cookie": setCookie, "X-Trace-Id": traceId };
  } catch (e) {
    console.error(e);
    return {};
  }
}
