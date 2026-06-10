import type {ReactNode} from 'react';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import Heading from '@theme/Heading';

function HomepageHeader() {
  const {siteConfig} = useDocusaurusContext();
  return (
    <header className="relative overflow-hidden bg-slate-50 dark:bg-slate-950">
      <div className="container py-20 md:py-24">
        <div className="mx-auto max-w-5xl text-center">
          <span className="inline-flex items-center rounded-full border border-brand-200 bg-brand-50 px-4 py-1 text-xs font-semibold uppercase tracking-[0.2em] text-brand-700 dark:border-slate-700 dark:bg-slate-800 dark:text-slate-200">
            Go HTML Scraping Toolkit
          </span>
        </div>
        <Heading as="h1" className="hero__title mt-6 text-center text-4xl font-black tracking-tight md:text-6xl">
          {siteConfig.title}
        </Heading>
        <p className="hero__subtitle mx-auto mt-4 max-w-2xl text-center text-lg text-slate-600 dark:text-slate-300">
          {siteConfig.tagline}
        </p>
        <p className="mx-auto mt-3 max-w-3xl text-lg leading-8 text-slate-700 dark:text-slate-300">
          Parse HTML from raw strings or URLs, traverse with CSS selectors or
          tag/attribute filters, and ship clean structured data to your
          pipeline in minutes.
        </p>
        <div className="mt-7 flex flex-wrap items-center justify-center gap-3">
          <Link
            className="inline-flex items-center rounded-xl bg-brand-600 px-6 py-3 text-sm font-semibold text-white shadow-soft transition hover:bg-brand-700"
            to="/docs/intro">
            Get Started
          </Link>
          <Link
            className="inline-flex items-center rounded-xl border border-slate-300 bg-white px-6 py-3 text-sm font-semibold text-slate-800 transition hover:border-brand-500 hover:text-brand-700 dark:border-slate-700 dark:bg-slate-900 dark:text-slate-200"
            to="/docs/api-overview">
            Explore API
          </Link>
        </div>

        <div className="mx-auto mt-12 grid max-w-5xl gap-4 md:grid-cols-3">
          <article className="rounded-2xl border border-slate-200 bg-white p-4 shadow-soft dark:border-slate-700 dark:bg-slate-900">
            <p className="text-3xl font-black text-brand-600 dark:text-brand-300">3</p>
            <p className="mt-1 text-sm text-slate-600 dark:text-slate-300">Selection styles: CSS, tags, and attributes</p>
          </article>
          <article className="rounded-2xl border border-slate-200 bg-white p-4 shadow-soft dark:border-slate-700 dark:bg-slate-900">
            <p className="text-3xl font-black text-brand-600 dark:text-brand-300">JSON · CSV · MD</p>
            <p className="mt-1 text-sm text-slate-600 dark:text-slate-300">Built-in exports for reporting and data pipelines</p>
          </article>
          <article className="rounded-2xl border border-slate-200 bg-white p-4 shadow-soft dark:border-slate-700 dark:bg-slate-900">
            <p className="text-3xl font-black text-brand-600 dark:text-brand-300">HTTP Ready</p>
            <p className="mt-1 text-sm text-slate-600 dark:text-slate-300">Headers, cookies, and proxy rotation included</p>
          </article>
        </div>
      </div>
    </header>
  );
}

export default function Home(): ReactNode {
  const {siteConfig} = useDocusaurusContext();
  return (
    <Layout
      title={`${siteConfig.title} Documentation`}
      description="Official documentation for nano scrape, an ergonomic HTML scraping library for Go.">
      <HomepageHeader />
      <main>
        <section className="py-10 md:py-16">
          <div className="container">
            <div className="mx-auto mb-10 max-w-2xl text-center">
              <Heading as="h2" className="text-3xl font-black tracking-tight md:text-4xl">
                Designed for real scraping workflows
              </Heading>
              <p className="mt-3 text-slate-600 dark:text-slate-300">
                Move from raw HTML to clean datasets with an API that stays readable under production pressure.
              </p>
            </div>
            <div className="grid gap-5 md:grid-cols-3">
              <article className="rounded-2xl border border-slate-200 bg-white p-6 shadow-soft transition hover:-translate-y-1 hover:shadow-xl dark:border-slate-700 dark:bg-slate-900">
                <Heading as="h3">Flexible Selection</Heading>
                <p>
                  Mix high-level CSS selectors with strict tag/attribute filters when pages are messy or inconsistent.
                </p>
              </article>
              <article className="rounded-2xl border border-slate-200 bg-white p-6 shadow-soft transition hover:-translate-y-1 hover:shadow-xl dark:border-slate-700 dark:bg-slate-900">
                <Heading as="h3">Structured Export</Heading>
                <p>
                  Convert extraction results into JSON, CSV, or Markdown without building custom serializers.
                </p>
              </article>
              <article className="rounded-2xl border border-slate-200 bg-white p-6 shadow-soft transition hover:-translate-y-1 hover:shadow-xl dark:border-slate-700 dark:bg-slate-900">
                <Heading as="h3">HTTP + Proxy Support</Heading>
                <p>
                  Handle realistic request flows with custom headers, cookie-aware sessions, and proxy rotation.
                </p>
              </article>
            </div>
          </div>
        </section>
      </main>
    </Layout>
  );
}
