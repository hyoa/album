const colors = require('tailwindcss/colors')

module.exports = {
  purge: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  darkMode: false, // or 'media' or 'class'
  theme: {
    extend: {
      inset: {
        'full': '100%',
        '-20%': '-20%'
      },
      opacity: {
        '10': '0.1',
        '90': '0.9'
      },
      colors: {
        primary: '#487CA0',
        'darker-primary': '#4299e1',
        'dark-primary': '#2c5282',
        'light-primary': '#A0D2DB',
        'lighter-primary': '#BEE7E8',
        'lightest-primary': '#E7F7F7',
        transparent: 'transparent',
        current: 'currentColor',
        black: colors.black,
        white: colors.white,
        gray: colors.trueGray,
        indigo: colors.indigo,
        red: colors.rose,
        yellow: colors.amber,
      }
    },
  },
  variants: {
    extend: {},
  },
  plugins: [],
}
