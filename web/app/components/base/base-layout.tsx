import clsx from "clsx";
import {
  Theme,
  ThemeProvider,
  PreventFlashOnWrongTheme,
  useTheme,
} from "remix-themes";
import { Links, Meta, Scripts, ScrollRestoration } from "@remix-run/react";
import { Toaster } from "../ui/toaster";

interface Props {
  theme: Theme | null;
  children: React.ReactNode;
}

export const BaseLayout = (props: Props) => {
  return (
    <ThemeProvider specifiedTheme={props.theme} themeAction="/action/theme">
      <LayoutInner {...props} />
    </ThemeProvider>
  );
};

const LayoutInner = (props: Props) => {
  const { children } = props;

  const [theme] = useTheme();

  return (
    <html lang="en" className={clsx(theme)}>
      <head>
        <meta charSet="utf-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <Meta />
        <PreventFlashOnWrongTheme ssrTheme={Boolean(theme)} />
        <Links />
      </head>
      <body>
        {children}
        <Toaster />
        <ScrollRestoration />
        <Scripts />
      </body>
    </html>
  );
};
