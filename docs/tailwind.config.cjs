/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './src/**/*.{js,jsx,ts,tsx,md,mdx}',
    './docs/**/*.{md,mdx}',
    './blog/**/*.{md,mdx}',
  ],
  theme: {
    extend: {
      colors: {
        brand: {
          50: '#eff8ff',
          100: '#dbeeff',
          200: '#bfdfff',
          300: '#93c8ff',
          400: '#60a7ff',
          500: '#3b86ff',
          600: '#215ee8',
          700: '#1e4bc4',
          800: '#1f419f',
          900: '#203b7d',
        },
      },
      boxShadow: {
        soft: '0 10px 30px rgba(2, 8, 23, 0.08)',
      },
      backgroundImage: {
        'hero-grid':
          'radial-gradient(circle at 20% 20%, rgba(59,130,246,0.18), transparent 30%), radial-gradient(circle at 80% 10%, rgba(45,212,191,0.15), transparent 25%)',
      },
    },
  },
  darkMode: ['class', '[data-theme="dark"]'],
  plugins: [],
};
