<template>
  <section class="benchmark-container">
    <div class="benchmark-header">
      <span class="benchmark-badge">BENCHMARKS</span>
      <h2 class="benchmark-title">Engineered for Raw Speed</h2>
      <p class="benchmark-subtitle">
        Independent benchmarks on Linux/amd64 (12th Gen Intel® i5-12500) parsing
        ~11KB HTML payloads.
      </p>
    </div>

    <!-- Stat Highlights Grid -->
    <div class="stats-grid">
      <div class="stat-card highlight">
        <div class="stat-icon">⚡</div>
        <div class="stat-value">483<span class="stat-unit">µs</span></div>
        <div class="stat-label">End-to-End Lifecycle</div>
        <div class="stat-desc">40% faster response + parse than Colly</div>
      </div>

      <div class="stat-card">
        <div class="stat-icon">🏎️</div>
        <div class="stat-value">257<span class="stat-unit">µs</span></div>
        <div class="stat-label">Pure CPU Parsing</div>
        <div class="stat-desc">15% faster CSS extraction than GoQuery</div>
      </div>

      <div class="stat-card">
        <div class="stat-icon">📦</div>
        <div class="stat-value">1,464</div>
        <div class="stat-label">Memory Allocations</div>
        <div class="stat-desc">Lowest allocation count per operation</div>
      </div>
    </div>

    <!-- Visual Bar Comparison -->
    <div class="bars-card">
      <h3 class="bars-title">End-to-End Speed Comparison (Lower is Better)</h3>

      <div class="bar-group">
        <div class="bar-info">
          <span class="lib-name font-bold text-brand">
            Nano Scrape <span class="rank-badge gold">🏆 1st</span>
          </span>
          <span class="bar-metrics">483,119 ns &bull; 1,571 allocs</span>
        </div>
        <div class="bar-track">
          <div class="bar-fill nano-bar" style="width: 59%"></div>
        </div>
      </div>

      <div class="bar-group">
        <div class="bar-info">
          <span class="lib-name">
            GoQuery <span class="rank-badge silver">🥈 2nd</span>
          </span>
          <span class="bar-metrics">583,903 ns &bull; 1,604 allocs</span>
        </div>
        <div class="bar-track">
          <div class="bar-fill goquery-bar" style="width: 71%"></div>
        </div>
      </div>

      <div class="bar-group">
        <div class="bar-info">
          <span class="lib-name">
            Colly <span class="rank-badge bronze">🥉 3rd</span>
          </span>
          <span class="bar-metrics">812,200 ns &bull; 1,827 allocs</span>
        </div>
        <div class="bar-track">
          <div class="bar-fill colly-bar" style="width: 100%"></div>
        </div>
      </div>
    </div>

    <!-- Respectful Library Overview -->
    <div class="philosophy-grid">
      <div class="philo-card">
        <h4>⚡ Nano Scrape</h4>
        <p>
          Focused on maximum throughput, zero dependency overhead, and a 4-line
          data extraction API.
        </p>
      </div>
      <div class="philo-card">
        <h4>🛡️ GoQuery</h4>
        <p>
          The established Go classic. Exceptional streaming memory efficiency
          (~94KB B/op) with jQuery semantics.
        </p>
      </div>
      <div class="philo-card">
        <h4>🌐 Colly</h4>
        <p>
          The feature-packed scraper framework. Built for complex web crawlers
          requiring async collector hooks.
        </p>
      </div>
    </div>
  </section>
</template>

<script setup>
// Interactive benchmark summary component
</script>
