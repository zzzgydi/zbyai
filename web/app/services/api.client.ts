import { API_BASE_URL } from "./constants";

const baseURL = API_BASE_URL;

export interface StreamAnswer {
  id: string;
  key: string;
  status: number;
  delta?: string;
  model?: string;
  errMsg?: string;
}

export interface StreamSearch {
  status: number;
  query?: string;
  search?: ThreadSearch[];
  errMsg?: string;
}

export type StreamHandler = {
  onQuery?: (data: string) => void;
  onSetting?: (data: ThreadSetting) => void;
  onSearch?: (data: StreamSearch) => void;
  onAnswer?: (data: StreamAnswer) => void;
  onError?: (error: string) => void;
  onDone?: () => void;
};

export class ApiClient {
  static async rewriteThread(
    id: string,
    runId: number
  ): Promise<{ id: string; runId: number }> {
    const resp = await fetch(`${baseURL}/inner/rewrite_thread`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ id, runId }),
      credentials: "include",
    });
    const data = await resp.json();
    if (data?.error) throw new Error(data.error);
    return data.result;
  }

  static async deleteThread(threadId: string) {
    const resp = await fetch(`${baseURL}/inner/delete_thread`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ id: threadId }),
      credentials: "include",
    });
    if (resp.status > 204) {
      throw new Error(`${resp.status} ${resp.statusText}`);
    }
    const data = await resp.json();
    if (data?.code !== 0) throw new Error(data?.msg || "unknown error");
    return data.data as Thread[];
  }

  static async listThread() {
    const resp = await fetch(`${baseURL}/inner/list_thread`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({}),
      credentials: "include",
    });
    if (resp.status > 204) {
      throw new Error(`${resp.status} ${resp.statusText}`);
    }
    const data = await resp.json();
    if (data?.code !== 0) throw new Error(data?.msg || "unknown error");
    return data.data as Thread[];
  }

  static async streamThread(
    id: string,
    runId: number,
    handler: StreamHandler
  ): Promise<boolean> {
    try {
      return await streamThreadInner(id, runId, handler);
    } catch (error: any) {
      console.error(error);
      handler.onError?.(error.message || error.toString());
    }
    return false;
  }
}

export async function streamThread(
  id: string,
  runId: number,
  handler: StreamHandler
): Promise<boolean> {
  try {
    return await streamThreadInner(id, runId, handler);
  } catch (error: any) {
    console.error(error);
    handler.onError?.(error.message || error.toString());
  }
  return false;
}

async function streamThreadInner(
  id: string,
  runId: number,
  handler: StreamHandler
): Promise<boolean> {
  const stream = await fetch(`${baseURL}/inner/stream_thread`, {
    method: "POST",
    headers: { "Content-Type": "text/event-stream" },
    body: JSON.stringify({ id, runId }),
    credentials: "include",
  });

  const reader = stream.body?.pipeThrough(new TextDecoderStream()).getReader();
  if (!reader) throw new Error("No reader");

  let finish = false;
  while (true) {
    const { value, done } = await reader.read();
    if (done) break;
    if (!value.startsWith("data: ")) continue;
    const newValue = value.slice(6);
    for (const line of newValue.split("\n\ndata: ")) {
      try {
        const data = JSON.parse(line) as any;
        if (data.type === "query") handler.onQuery?.(data.data);
        else if (data.type === "setting") handler.onSetting?.(data.data);
        else if (data.type === "search") handler.onSearch?.(data.data);
        else if (data.type === "answer") handler.onAnswer?.(data.data);
        else if (data.type === "error") {
          finish = true;
          handler.onError?.(data.data?.errMsg || "unknown error");
        } else if (data.type === "done") {
          finish = true;
          handler.onDone?.();
        } else console.error("Unknown type", data);
      } catch (error) {
        console.error("Error parsing JSON", line, error);
      }
    }
  }
  return finish;
}
