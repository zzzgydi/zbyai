import { authHandler } from "./auth.server";
import { API_BASE_URL } from "./constants";

const baseURL = API_BASE_URL;

export class ApiServer {
  request: Request;
  response: Response;

  constructor(request: Request, response?: Response) {
    this.request = request;
    this.response = response ?? new Response();
  }

  responseHeaders(): Record<string, string> {
    const traceId = this.response.headers.get("X-Trace-Id") || "";
    const setCookie = this.response.headers.get("Set-Cookie");
    if (!setCookie) return { "X-Trace-Id": traceId };
    return { "Set-Cookie": setCookie, "X-Trace-Id": traceId };
  }

  async get<T>(url: string, params: any) {
    await authHandler(this.request);
    const headers = this.request.headers || new Headers();
    headers.set("Content-Type", "application/json");
    headers.delete("host"); // 这玩意儿会导致请求又转到当前服务上

    const resp = await fetch(`${baseURL}${url}`, {
      method: "GET",
      headers,
      credentials: "include",
    });

    for (const [key, value] of resp.headers) {
      this.response.headers.append(key, value);
    }

    const data = await resp.json();
    if (data?.code !== 0) throw new Error(data?.msg || "unknown error");
    return data.data as T;
  }

  async post<T>(url: string, body: any) {
    await authHandler(this.request);
    const headers = this.request.headers || new Headers();
    headers.set("Content-Type", "application/json");
    headers.delete("host"); // 这玩意儿会导致请求又转到当前服务上

    const resp = await fetch(`${baseURL}${url}`, {
      method: "POST",
      headers,
      credentials: "include",
      body: JSON.stringify(body),
    });

    for (const [key, value] of resp.headers) {
      this.response.headers.append(key, value);
    }

    if (resp.status > 204) {
      throw new Error(`${resp.status} ${resp.statusText}`);
    }

    const data = await resp.json();
    if (data?.code !== 0) throw new Error(data?.msg || "unknown error");
    return data.data as T;
  }

  async listThread() {
    return this.post<Thread[]>("/inner/list_thread", {});
  }

  async createThread(query: string) {
    return this.post<{ id: string }>("/inner/create_thread", { query });
  }

  async appendThread(id: string, query: string) {
    if (!id) throw new Error("No thread ID");
    return this.post<{ id: string; runId: number }>("/inner/append_thread", {
      id,
      query,
    });
  }

  async rewriteThread(id: string, runId: number) {
    if (!id) throw new Error("No thread ID");
    if (!runId) throw new Error("No run ID");
    return this.post<{ id: string; runId: number }>("/inner/rewrite_thread", {
      id,
      runId,
    });
  }

  async detailThread(id: string) {
    if (!id) throw new Error("No thread ID");
    return this.post<ThreadDetail>("/inner/detail_thread", { id });
  }
}
