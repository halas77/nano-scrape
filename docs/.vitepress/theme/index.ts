import { h } from 'vue'
import DefaultTheme from 'vitepress/theme'
import type { EnhanceAppContext } from 'vitepress'
import HeroCode from './HeroCode.vue'
import BenchmarkSection from './BenchmarkSection.vue'
import './custom.css'

export default {
  extends: DefaultTheme,
  Layout() {
    return h(DefaultTheme.Layout, null, {
      'home-hero-image': () => h(HeroCode)
    })
  },
  enhanceApp({ app }: EnhanceAppContext) {
    app.component('BenchmarkSection', BenchmarkSection)
  }
}
