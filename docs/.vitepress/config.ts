import { defineConfig } from "vitepress";
const nav = [
  {
    text: "Blog",
    link: "https://micromdm.io/blog",
  },
  {
    text: "Docs",
    activeMatch: `^/(guide|style-guide|cookbook|examples)/`,
    items: [
      { text: "Guide", link: "/user-guide/introduction" },
      { text: "Quick Start", link: "/user-guide/quickstart" },
    ],
  },
  {
    text: "Code",
    link: "https://github.com/micromdm/micromdm",
  },
  {
    text: "Releases",
    link: "https://github.com/micromdm/micromdm/releases/latest",
  },
];

export default defineConfig({
  lang: "en-US",
  title: "MicroMDM",
  description: "bootstrap your mac deployment",

  themeConfig: {
    logo: "/logo.svg",
    nav,
    sidebar: getSidebar(),
  },
});

function getSidebar() {
  return [
    {
      text: "User Guide",
      children: [
        { text: "Introduction", link: "/user-guide/introduction" },
        { text: "Quick Start", link: "/user-guide/quickstart" },
        { text: "Enrolling Devices", link: "/user-guide/enrolling-devices" },
      ],
    },
    {
      text: "Architechture",
      children: [
        {
          text: "no_device_configured_without_blueprint",
          link: "/architecture/2019-05-03_no_device_configured_without_blueprint",
        },
      ],
    },
  ];
}
