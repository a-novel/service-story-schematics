import { defineConfig } from "vitepress";
import { viteStaticCopy } from "vite-plugin-static-copy";

// https://vitepress.dev/reference/site-config
export default defineConfig({
  title: "Service JSON Keys",
  titleTemplate: "A-Novel",
  themeConfig: {
    // https://vitepress.dev/reference/default-theme-config
    nav: [
      { text: "Reference", link: "/service/containerized" },
      { text: "API Reference", link: "/api.html", target: "_blank" },
    ],

    sidebar: [
      {
        text: "Run As a Service",
        items: [{ text: "Containerized", link: "/service/containerized" }],
      },
      {
        text: "Package",
        items: [{ text: "React client", link: "/package/react" }],
      },
    ],

    socialLinks: [{ icon: "github", link: "https://github.com/a-novel/service-story-schematics" }],
  },

  head: [["link", { rel: "icon", href: "./icon.png" }]],

  base: "/service-story-schematics/",

  vite: {
    plugins: [
      viteStaticCopy({
        targets: [
          {
            src: "api.yaml",
            dest: "./",
          },
        ],
      }),
    ],
  },
});
