import {themes as prismThemes} from 'prism-react-renderer';
import type {Config} from '@docusaurus/types';
import type * as Preset from '@docusaurus/preset-classic';

// This runs in Node.js - Don't use client-side code here (browser APIs, JSX...)

const config: Config = {
  title: 'nano scrape',
  tagline: 'Fast, ergonomic HTML scraping for Go',
  favicon: 'img/nano-scrape.jpg',

  // Future flags, see https://docusaurus.io/docs/api/docusaurus-config#future
  future: {
    v4: true, // Improve compatibility with the upcoming Docusaurus v4
  },

  // Set the production url of your site here
  url: 'https://halas77.github.io',
  // Set the /<baseUrl>/ pathname under which your site is served
  // For GitHub pages deployment, it is often '/<projectName>/'
  baseUrl: '/',

  // GitHub pages deployment config.
  // If you aren't using GitHub pages, you don't need these.
  organizationName: 'halas77', // Usually your GitHub org/user name.
  projectName: 'goscrape', // Usually your repo name.

  onBrokenLinks: 'throw',

  // Even if you don't use internationalization, you can use this field to set
  // useful metadata like html lang. For example, if your site is Chinese, you
  // may want to replace "en" with "zh-Hans".
  i18n: {
    defaultLocale: 'en',
    locales: ['en'],
  },

  presets: [
    [
      'classic',
      {
        docs: {
          sidebarPath: './sidebars.ts',
          editUrl: 'https://github.com/halas77/goscrape/tree/main/docs/',
        },
        blog: false,
        theme: {
          customCss: './src/css/tailwind.css',
        },
      } satisfies Preset.Options,
    ],
  ],

  themeConfig: {
    // Replace with your project's social card
    image: 'img/docusaurus-social-card.jpg',
    colorMode: {
      defaultMode: 'light',
      disableSwitch: false,
      respectPrefersColorScheme: false,
    },
    announcementBar: {
      id: 'github-star',
      content:
        '⭐ If nano scrape helps you, consider starring it on <a target="_blank" rel="noopener noreferrer" href="https://github.com/halas77/goscrape">GitHub</a>.',
      backgroundColor: '#0f172a',
      textColor: '#e2e8f0',
      isCloseable: true,
    },
    navbar: {
      title: 'nano scrape',
      logo: {
        alt: 'nano scrape logo',
        src: 'img/nano-scrape.jpg',
        width: 32,
        height: 32,
      },
      items: [
        {
          to: '/docs/intro',
          label: 'Docs',
          position: 'left',
        },
        {
          to: '/docs/api-overview',
          label: 'API',
          position: 'left',
        },
        {
          to: '/docs/http-requests',
          label: 'HTTP',
          position: 'left',
        },
        {
          type: 'dropdown',
          label: 'Guides',
          position: 'left',
          items: [
            {
              label: 'Getting Started',
              to: '/docs/getting-started',
            },
            {
              label: 'Exporting Data',
              to: '/docs/exporting-data',
            },
            {
              label: 'Deployment',
              to: '/docs/tutorial-basics/deploy-your-site',
            },
          ],
        },
        {
          to: '/docs/getting-started',
          label: 'Quickstart',
          position: 'right',
        },
        {
          href: 'https://github.com/halas77/goscrape',
          label: 'GitHub',
          position: 'right',
        },
      ],
    },
    footer: {
      style: 'dark',
      links: [
        {
          title: 'Docs',
          items: [
            {
              label: 'Introduction',
              to: '/docs/intro',
            },
            {
              label: 'Getting Started',
              to: '/docs/getting-started',
            },
            {
              label: 'API Overview',
              to: '/docs/api-overview',
            },
          ],
        },
        {
          title: 'Guides',
          items: [
            {
              label: 'HTTP Requests',
              to: '/docs/http-requests',
            },
            {
              label: 'Exporting Data',
              to: '/docs/exporting-data',
            },
            {
              label: 'Deploy your site',
              to: '/docs/tutorial-basics/deploy-your-site',
            },
          ],
        },
        {
          title: 'Project',
          items: [
            {
              label: 'GitHub Repository',
              href: 'https://github.com/halas77/goscrape',
            },
            {
              label: 'Issues',
              href: 'https://github.com/halas77/goscrape/issues',
            },
            {
              label: 'Discussions',
              href: 'https://github.com/halas77/goscrape/discussions',
            },
          ],
        },
      ],
      copyright: `Copyright © ${new Date().getFullYear()} nano scrape. Crafted for practical Go scraping workflows.`,
    },
    prism: {
      theme: prismThemes.github,
      darkTheme: prismThemes.dracula,
      additionalLanguages: ['go', 'bash'],
    },
  } satisfies Preset.ThemeConfig,
};

export default config;
