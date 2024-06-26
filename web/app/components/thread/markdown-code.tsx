import type { FC, HTMLAttributes, ReactElement } from "react";
import { Children, useId } from "react";
import type { Language } from "prism-react-renderer";
import { Highlight, themes } from "prism-react-renderer";

export function getLanguageFromClassName(className: string) {
  const match = className.match(/language-(\w+)/);
  return match ? match[1] : "";
}

function isLanguageSupported(lang: string): lang is Language {
  return (
    lang === "markup" ||
    lang === "bash" ||
    lang === "clike" ||
    lang === "c" ||
    lang === "cpp" ||
    lang === "css" ||
    lang === "javascript" ||
    lang === "jsx" ||
    lang === "coffeescript" ||
    lang === "actionscript" ||
    lang === "css-extr" ||
    lang === "diff" ||
    lang === "git" ||
    lang === "go" ||
    lang === "graphql" ||
    lang === "handlebars" ||
    lang === "json" ||
    lang === "less" ||
    lang === "makefile" ||
    lang === "markdown" ||
    lang === "objectivec" ||
    lang === "ocaml" ||
    lang === "python" ||
    lang === "reason" ||
    lang === "sass" ||
    lang === "scss" ||
    lang === "sql" ||
    lang === "stylus" ||
    lang === "tsx" ||
    lang === "typescript" ||
    lang === "wasm" ||
    lang === "yaml"
  );
}

interface Props {
  code: string;
  language: string;
  className?: string;
}

export const CodeBlock = (props: Props) => {
  const lang = isLanguageSupported(props.language) ? props.language : "bash";

  return (
    <Highlight theme={themes.shadesOfPurple} code={props.code} language={lang}>
      {({ className, tokens, getLineProps, getTokenProps }) => (
        <pre className={`overflow-scroll relative ${className}`} style={{}}>
          <code className={className} style={{}}>
            {tokens.map((line, i) => (
              <div key={i} {...getLineProps({ line, key: i })} style={{}}>
                {line.map((token, key) => (
                  <span
                    key={key}
                    {...getTokenProps({ token, key })}
                    style={{}}
                  />
                ))}
              </div>
            ))}
          </code>
        </pre>
      )}
    </Highlight>
  );
};
