import { createCookieSessionStorage } from "@remix-run/node";
import { createThemeSessionResolver } from "remix-themes";

const isProduction = process.env.NODE_ENV === "production";
const secret = process.env.SESSION_SECRET ?? "s3cr3ttt";

const sessionStorage = createCookieSessionStorage({
  cookie: {
    name: "__theme",
    path: "/",
    httpOnly: true,
    sameSite: "lax",
    secrets: [secret],
    ...(isProduction ? { domain: "your-production-domain.com", secure: true } : {}), // TODO: Set domain
  },
});

export const themeSessionResolver = createThemeSessionResolver(sessionStorage);
