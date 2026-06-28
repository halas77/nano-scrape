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
        text: "Guides & Features",
        items: [
          { text: "Getting Started", link: "/guide/getting-started" },
          { text: "Document Initialization", link: "/init-doc/init-document" },
          { text: "Request Client", link: "/requests/request" },
          { text: "Proxy Rotation", link: "/requests/proxies" },
          { text: "CSS Selection", link: "/traversals/css-traversal" },
          { text: "Attribute Selection", link: "/traversals/attribute-traversal" },
          { text: "Data Exporting & Mapping", link: "/export/exporting" },
          { text: "Formatting & Print Helpers", link: "/helper/helpers" },
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
