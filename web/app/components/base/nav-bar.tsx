import LogoSvg from "@/assets/images/icon-logo.svg";
import { Link } from "@remix-run/react";
import { ThemeToggle } from "../theme/theme-toggle";

interface Props {
  children?: React.ReactNode;
}

export const NavBarWrapper = (props: Props) => {
  return (
    <header className="sticky top-0 z-50 w-full border-border/40 bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="container px-4 md:px-8 w-full flex h-14 max-w-[1128px] items-center justify-between">
        <Link to="/" prefetch="render">
          <div className="flex items-center gap-1">
            <img className="w-5 h-5" src={LogoSvg} alt="logo" />
            <h1 className="text-base font-semibold">ZByAI</h1>
          </div>
        </Link>

        <div className="flex items-center gap-4">
          <ThemeToggle />

          {props.children}
        </div>
      </div>
    </header>
  );
};
