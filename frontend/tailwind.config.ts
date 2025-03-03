import type { Config } from 'tailwindcss';

export default {
  darkMode: ['class'],
  content: ['./app/**/{**,.client,.server}/**/*.{js,jsx,ts,tsx}'],
  theme: {
    extend: {
      fontFamily: {
        sans: [
          'Inter',
          'ui-sans-serif',
          'system-ui',
          'sans-serif',
          'Apple Color Emoji',
          'Segoe UI Emoji',
          'Segoe UI Symbol',
          'Noto Color Emoji',
        ],
      },
      borderRadius: {
        lg: 'var(--radius)',
        md: 'calc(var(--radius) - 2px)',
        sm: 'calc(var(--radius) - 4px)',
      },
      colors: {
        background: "hsl(var(--background), <alpha-value>)",
        foreground: "hsl(var(--foreground), <alpha-value>)",
        primary: {
          DEFAULT: "hsl(var(--primary-a0), <alpha-value>)",
          a0: "hsl(var(--primary-a0), <alpha-value>)",
          a10: "hsl(var(--primary-a10), <alpha-value>)",
          a20: "hsl(var(--primary-a20), <alpha-value>)",
          a30: "hsl(var(--primary-a30), <alpha-value>)",
          a40: "hsl(var(--primary-a40), <alpha-value>)",
          a50: "hsl(var(--primary-a50), <alpha-value>)",
          foreground: "hsl(var(--primary-foreground), <alpha-value>)"
        },
        surface: {
          a0: "hsl(var(--surface-a0), <alpha-value>)",
          a10: "hsl(var(--surface-a10), <alpha-value>)",
          a20: "hsl(var(--surface-a20), <alpha-value>)",
          a30: "hsl(var(--surface-a30), <alpha-value>)",
          a40: "hsl(var(--surface-a40), <alpha-value>)",
          a50: "hsl(var(--surface-a50), <alpha-value>)",
        },
        accent: {
          DEFAULT: "hsl(var(--surface-a10), <alpha-value>)",
          foreground: "hsl(var(--primary-a0), <alpha-value>)",
        }
      },
    },
  },
  plugins: [require('tailwindcss-animate')],
} satisfies Config;
