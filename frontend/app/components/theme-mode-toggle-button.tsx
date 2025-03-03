import { Moon, Sun } from "lucide-react";
import { Theme, useTheme } from "remix-themes";

import { Button } from "./ui/button";

export function ThemeModeToggleButton() {
  const [theme, setTheme] = useTheme();

  if (theme === Theme.LIGHT) {
    return (
      <Button variant="ghost" size="icon" onClick={() => setTheme(Theme.DARK)}>
        <Sun />
      </Button>
    );
  } else {
    return (
      <Button variant="ghost" size="icon" onClick={() => setTheme(Theme.LIGHT)}>
        <Moon />
      </Button>
    );
  }
}
