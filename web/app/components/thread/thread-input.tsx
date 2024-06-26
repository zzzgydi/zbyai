import clsx from "clsx";
import { forwardRef, useImperativeHandle, useRef, useState } from "react";
import { useFetcher } from "@remix-run/react";
import { Forward } from "lucide-react";
import { BaseInput } from "../base/base-input";
import { Button } from "../ui/button";

interface Props {
  threadId: string;
}

export interface ThreadInputRef {
  reset: () => void;
  startLoading: () => void;
  resetLoading: () => void;
}

export const ThreadInput = forwardRef<ThreadInputRef, Props>((props, ref) => {
  const { threadId } = props;
  const fetcher = useFetcher<{ error?: string }>();

  const [value, setValue] = useState("");
  const [focus, setFocus] = useState(false);
  const [loading, setLoading] = useState(false);
  const inputRef = useRef<HTMLTextAreaElement>(null);

  const handleAppend = async () => {
    const query = value.trim();
    if (!query || fetcher.state !== "idle") return;

    setLoading(true);
    fetcher.submit(
      { id: threadId, query },
      {
        method: "POST",
        encType: "application/json",
        action: "/action/thread_append",
      }
    );
  };

  useImperativeHandle(ref, () => ({
    reset: () => {
      setValue("");
      setFocus(false);
      // inputRef.current?.focus();
    },
    startLoading: () => setLoading(true),
    resetLoading: () => setLoading(false),
  }));

  return (
    <div
      className={clsx(
        "relative pt-3 pb-8 pl-5 pr-3 border-solid border border-muted rounded-2xl",
        "backdrop-blur supports-[backdrop-filter]:bg-muted/85 dark:supports-[backdrop-filter]:bg-muted/85",
        focus ? "ring-2 ring-muted-foreground" : "thread-shadow"
      )}
    >
      <BaseInput
        className="w-full input-scrollbar"
        ref={inputRef}
        value={value}
        placeholder="Ask follow-up..."
        autoFocus={false}
        onChange={(v) => setValue(v)}
        onEnter={handleAppend}
        onFocus={setFocus}
      />

      <div className="absolute bottom-2 w-8 h-8 right-4 flex items-center justify-end">
        {loading ? (
          <div className="second-loader"></div>
        ) : (
          <Button
            variant="ghost"
            className="w-8 h-8 p-0"
            onClick={handleAppend}
            aria-label="Send"
          >
            <Forward className="text-muted-foreground" />
          </Button>
        )}
      </div>
    </div>
  );
});
