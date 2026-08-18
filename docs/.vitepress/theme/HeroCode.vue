<template>
  <div class="hero-code-wrapper">
    <div class="glow-backdrop"></div>
    <div class="code-card">
      <div class="code-header">
        <div class="window-controls">
          <span class="control close"></span>
          <span class="control minimize"></span>
          <span class="control maximize"></span>
        </div>
        <div class="tab-title">
          <svg
            class="go-icon"
            viewBox="0 0 24 24"
            width="14"
            height="14"
            fill="currentColor"
          >
            <path d="M2 6h2v12H2zm4 0h2v12H6zm4 5h12v2H10z" />
          </svg>
          main.go
        </div>
        <button
          class="copy-btn"
          @click="copyCode"
          :title="copied ? 'Copied!' : 'Copy Code'"
        >
          <svg
            v-if="!copied"
            width="13"
            height="13"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
          >
            <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
            <path
              d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"
            ></path>
          </svg>
          <svg
            v-else
            width="13"
            height="13"
            viewBox="0 0 24 24"
            fill="none"
            stroke="#10b981"
            stroke-width="2.5"
            stroke-linecap="round"
            stroke-linejoin="round"
          >
            <polyline points="20 6 9 17 4 12"></polyline>
          </svg>
          <span>{{ copied ? "Copied" : "Copy" }}</span>
        </button>
      </div>

      <div class="code-body">
        <div class="code-line">
          <span class="line-num">1</span>
          <span class="code-content">
            <span class="token var">doc</span><span class="token op">,</span>
            <span class="token var">_</span> <span class="token op">:= </span>
            <span class="token pkg">nano</span><span class="token punct">.</span
            ><span class="token func">LoadDocument</span
            ><span class="token punct">(</span
            ><span class="token string">"https://example.com"</span
            ><span class="token punct">)</span>
          </span>
        </div>
        <div class="code-line">
          <span class="line-num">2</span>
          <span class="code-content">
            <span class="token var">items</span>
            <span class="token op"> :=</span> <span class="token var">doc</span
            ><span class="token punct">.</span
            ><span class="token func">SelectAll</span
            ><span class="token punct">(</span
            ><span class="token string">".product-card"</span
            ><span class="token punct">)</span>
          </span>
        </div>
        <div class="code-line">
          <span class="line-num">3</span>
          <span class="code-content">
            <span class="token var">data</span>
            <span class="token op"> := </span>
            <span class="token var">items</span
            ><span class="token punct">.</span
            ><span class="token func">Map</span
            ><span class="token punct">(</span
            ><span class="token type">map</span
            ><span class="token punct">[</span
            ><span class="token type">string</span
            ><span class="token punct">]</span
            ><span class="token type">string</span
            ><span class="token punct">{</span
            ><span class="token key">"price"</span
            ><span class="token op">:</span>
            <span class="token string">".price"</span
            ><span class="token punct">})</span>
          </span>
        </div>
        <div class="code-line">
          <span class="line-num">4</span>
          <span class="code-content">
            <span class="token pkg">nano</span><span class="token punct">.</span
            ><span class="token func">WriteMappedJSON</span
            ><span class="token punct">(</span
            ><span class="token string">"product-prices.json"</span
            ><span class="token punct">,</span>
            <span class="token var">data</span
            ><span class="token punct">)</span>
          </span>
        </div>
      </div>

      <div class="code-footer">
        <div class="status-indicator">
          <span class="pulse-dot"></span>
          <span class="status-text">Simple &amp; Expressive</span>
        </div>
        <div class="speed-badge">~500µs benchmark</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from "vue";

const copied = ref(false);
const rawCode = `doc, _ := nano.LoadDocument("https://example.com")
items  := doc.SelectAll(".product-card")
data   := items.Map(map[string]string{"price": ".price"})
nano.WriteMappedJSON("product-prices.json", data)`;

function copyCode() {
  if (navigator.clipboard) {
    navigator.clipboard.writeText(rawCode);
    copied.value = true;
    setTimeout(() => {
      copied.value = false;
    }, 2000);
  }
}
</script>
