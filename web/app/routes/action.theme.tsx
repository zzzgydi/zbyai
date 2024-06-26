import { createThemeAction } from "remix-themes";
import { themeSessionResolver } from "@/components/theme/theme.server";

export const action = createThemeAction(themeSessionResolver);
