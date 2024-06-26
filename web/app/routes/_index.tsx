import clsx from "clsx";
import Logo from "~/assets/images/icon-logo.svg?react";
import { Link, useFetcher } from "@remix-run/react";
import { useRef, useState } from "react";
import { Forward } from "lucide-react";
import { BaseInput } from "~/components/base/base-input";
import { Button } from "~/components/ui/button";

export default function Index() {
  const fetcher = useFetcher<{ error: string }>();

  const [value, setValue] = useState("");
  const [focus, setFocus] = useState(true);
  const inputRef = useRef<HTMLTextAreaElement>(null);

  const handleCreate = async () => {
    const query = value.trim();
    if (!query || fetcher.state !== "idle") return;
    fetcher.submit(
      { query },
      {
        method: "POST",
        encType: "application/json",
        action: "/action/thread_create",
      }
    );
  };

  return (
    <div className="flex-auto container px-4 md:px-8 w-full mx-auto max-w-[980px] flex flex-col">
      <div className="pt-[14vh] lg:pt-[15vh] pb-4 text-center mb-5 lg:mb-8 overflow-hidden transition-all">
        <div className="flex items-center justify-center gap-0.5">
          <Logo className="w-[50px] h-[50px] lg:w-[64px] lg:h-[64px]" />
          <div className="zbyai-text text-[52px] lg:text-[64px]">
            <span className="hidden">Z</span>
            <span>ByAI</span>
          </div>
        </div>
        <p className="text-lg lg:text-2xl text-secondary-foreground mt-1">
          Search. Elevate. Discover with AI.
        </p>
      </div>

      <div className="pb-40">
        <div
          className={clsx(
            "relative pt-3 pb-8 pl-5 pr-3 border-solid border border-muted rounded-2xl",
            "bg-muted",
            focus && "ring-2 ring-muted-foreground dark:ring-muted-foreground"
          )}
        >
          <BaseInput
            className="w-full input-scrollbar"
            ref={inputRef}
            value={value}
            autoFocus
            placeholder="Search Anything..."
            onChange={(v) => setValue(v)}
            onEnter={handleCreate}
            onFocus={setFocus}
          />

          <div className="absolute bottom-2 w-8 h-8 right-4 flex items-center justify-end">
            {fetcher.state !== "idle" ? (
              <div className="second-loader"></div>
            ) : (
              <Button
                variant="ghost"
                className="w-8 h-8 p-0"
                onClick={handleCreate}
                aria-label="Send"
              >
                <Forward className="text-muted-foreground" />
              </Button>
            )}
          </div>
        </div>

        {fetcher.data?.error && (
          <div className="flex items-center justify-center mt-2 px-5">
            <div className="text-destructive truncate">
              {fetcher.data?.error}
            </div>
          </div>
        )}
      </div>

      <footer className="flex-none p-3 mt-auto">
        <div className="text-balance text-center text-sm leading-loose text-muted-foreground">
          <div className="flex items-center gap-4 justify-center">
            {/* <p className="">&copy; 2024 ZByAI. All rights reserved.</p> */}
            <Link to="/" className="hover:underline">
              Terms of Service
            </Link>
            <Link to="/" className="hover:underline">
              Privacy Policy
            </Link>
          </div>
        </div>
      </footer>
    </div>
  );
}
