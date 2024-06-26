import React from "react";
import ReactMarkdown from "react-markdown";
import mermaid from "mermaid";
import RemarkMath from "remark-math";
import RemarkBreaks from "remark-breaks";
import RehypeKatex from "rehype-katex";
import RemarkGfm from "remark-gfm";
import RehypeHighlight from "rehype-highlight";
import { useRef, useState, useEffect } from "react";
import { useDebounceFn, useThrottleFn } from "ahooks";
import { Check, Copy, SquareArrowOutUpRight } from "lucide-react";
import { Button } from "../ui/button";
import { toast } from "../ui/use-toast";

const MarkdownInner = (props: { content: string }) => {
  return (
    <ReactMarkdown
      className="markdown-body"
      remarkPlugins={[RemarkMath, RemarkGfm, RemarkBreaks]}
      rehypePlugins={[
        RehypeKatex,
        [RehypeHighlight, { detect: false, ignoreMissing: true }],
      ]}
      components={{
        pre: PreCode,
        a: (aProps) => {
          return (
            <a
              href={aProps.href}
              target="_blank"
              className="inline-flex items-center gap-1"
            >
              {aProps.children}
              <SquareArrowOutUpRight className="w-3 h-3 text-[var(--color-accent-fg)]" />
            </a>
          );
        },
      }}
    >
      {props.content}
    </ReactMarkdown>
  );
};

export const MarkdownContent = React.memo(MarkdownInner);

function Mermaid(props: { code: string }) {
  const ref = useRef<HTMLDivElement>(null);
  const [hasError, setHasError] = useState(false);

  useEffect(() => {
    if (props.code && ref.current) {
      mermaid
        .run({
          nodes: [ref.current],
          suppressErrors: true,
        })
        .catch((e) => {
          setHasError(true);
          console.error("[Mermaid] ", e.message);
        });
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [props.code]);

  function viewSvgInNewWindow() {
    const svg = ref.current?.querySelector("svg");
    if (!svg) return;
    const text = new XMLSerializer().serializeToString(svg);
    const blob = new Blob([text], { type: "image/svg+xml" });
    const url = URL.createObjectURL(blob);
    const win = window.open(url);
    if (win) {
      win.onload = () => URL.revokeObjectURL(url);
    }
  }

  if (hasError) {
    return null;
  }

  return (
    <div
      className="no-dark mermaid"
      style={{
        cursor: "pointer",
        overflow: "auto",
      }}
      ref={ref}
      onClick={() => viewSvgInNewWindow()}
    >
      {props.code}
    </div>
  );
}
function PreCode(props: { children: any }) {
  const ref = useRef<HTMLPreElement>(null);
  const refText = ref.current?.innerText;
  const [mermaidCode, setMermaidCode] = useState("");

  const [copyed, setCopyed] = useState(false);

  const { run: renderMermaid } = useDebounceFn(
    () => {
      if (!ref.current) return;
      const mermaidDom = ref.current.querySelector("code.language-mermaid");
      if (mermaidDom) {
        setMermaidCode((mermaidDom as HTMLElement).innerText);
      }
    },
    { wait: 600 }
  );

  useEffect(() => {
    setTimeout(renderMermaid, 1);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [refText]);

  const handleCopyed = () => {
    setCopyed(true);
    setTimeout(() => {
      setCopyed(false);
    }, 2000);
  };

  return (
    <>
      {mermaidCode.length > 0 && (
        <Mermaid code={mermaidCode} key={mermaidCode} />
      )}
      <pre ref={ref} className="relative">
        <div className="absolute top-1 right-1">
          <Button
            size="icon"
            className="w-8 h-8"
            variant="ghost"
            aria-label="Copy"
            onClick={async () => {
              if (!ref.current) return;
              const code = ref.current.innerText;
              try {
                await navigator.clipboard.writeText(code);
                handleCopyed();
              } catch {
                toast({ title: "Failed to copy to clipboard!" });
              }
            }}
          >
            {copyed ? (
              <Check className="w-5 h-5" />
            ) : (
              <Copy className="w-5 h-5" />
            )}
          </Button>
        </div>

        {props.children}
      </pre>
    </>
  );
}
