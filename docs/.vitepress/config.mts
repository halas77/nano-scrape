import { defineConfig } from "vitepress";

// https://vitepress.dev/reference/site-config
export default defineConfig({
  title: "Nano Scrape",
  description:
    "A high-performance, lightweight web scraper and HTML parser for Go.",
  themeConfig: {
    // https://vitepress.dev/reference/default-theme-config
    nav: [
      { text: "Home", link: "/" },
      { text: "Guide", link: "/guide/getting-started" },
      { text: "API Reference", link: "/api/reference" },
    ],

    sidebar: [
      {
        text: "Guides",
        items: [
          { text: "Getting Started", link: "/guide/getting-started" },
          { text: "Request Client", link: "/guide/request-client" },
          {
            text: "Navigation & Selection",
            link: "/guide/navigation-selection",
          },
          { text: "Mapping & Exporting", link: "/guide/mapping-exporting" },
        ],
      },
      {
        text: "Reference",
        items: [{ text: "API Reference", link: "/api/reference" }],
      },
    ],

    socialLinks: [
      { icon: "github", link: "https://github.com/halas77/nano-scrape" },
    ],

    footer: {
      message: "Released under the MIT License.",
      copyright: "Copyright © 2026-present",
    },
  },
});
