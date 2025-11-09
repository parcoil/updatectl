import { defineConfig } from "vitepress";

const currentTimestamp = Date.now();
const currentDate = new Date(currentTimestamp);
const currentYear = currentDate.getFullYear();

export default defineConfig({
  title: "Updatectl",
  description: "A CLI tool for automating project updates",
  base: "/",
  head: [
    ["link", { rel: "icon", href: "/icon.ico" }],
    [
      "link",
      {
        rel: "stylesheet",
        href: "https://fonts.googleapis.com/css2?family=Material+Symbols+Rounded:opsz,wght,FILL,GRAD@20..48,100..700,0..1,-50..200",
      },
    ],
  ],

  themeConfig: {
    logo: "/logo.svg",
    search: {
      provider: "local",
    },
    footer: {
      copyright: `Copyright © ${currentYear} Parcoil`,
    },
    nav: [
      { text: "Home", link: "/" },
      { text: "Getting Started", link: "/quickstart" },
    ],
    sidebar: [
      {
        text: "Guide",
        link: "/guide",
      },
      {
        text: "Getting Started",
        link: "/quickstart",
      },
      {
        text: "User Guide",
        items: [
          { text: "Configuration", link: "/configuration" },
          { text: "Project Types", link: "/project-types" },
          { text: "Troubleshooting", link: "/troubleshooting" },
          { text: "Custom Commands", link: "/custom-commands" },
          { text: "Monitoring", link: "/monitoring" },
        ],
      },
      {
        text: "API Reference",
        items: [
          { text: "CLI Commands", link: "/cli" },
          { text: "Configuration Schema", link: "/schema" },
        ],
      },
      {
        text: "Contributing",
        link: "/contributing",
      },
    ],
  },
});
