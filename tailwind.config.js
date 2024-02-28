const { addDynamicIconSelectors } = require('@iconify/tailwind')

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./**/*.templ'],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        bgcolor: '#080f25',
        colorcard: '#101935',
        indigo2: '#6C72FF',
        indigo3: '#aeb9e1',
        blue2: '#212c4d',
        gray2: '#343b4f',
        gray3: '#37446b',
        purple: '#7478d8',
        gray4: '#e8e8e8',
        green: {
          150: '#0e3b3f',
          250: '0c5647',
          350: '#14ca74',
        },
        red: {
          150: '#40263e',
          250: '#663046',
          350: '#ff5a65',
        },
      },
      width: {
        300: '300px',
      },
      height: {
        150: '150px',
      },
      fontFamily: {
        inter: ['Inter', 'sans-serif'],
      },
      boxShadow: {
        cardshadow: '0 2px 7px rgba(20, 20, 43, .06)',
      },
      borderRadius: {
        xii: '12px',
        greenbadge: '2.322px',
        dropdown: '4.644px',
        point: '50%',
      },
      borderWidth: {
        0.6: '0.6px',
        1.16: '1.16px',
        2: '2px',
      },
      borderColor: {
        gray2: '#343b4f',
        green1: 'rgba(5, 193, 104, 0.20)',
        red1: 'rgba(255, 90, 101, 0.20)',
        dropdown: '#37446B',
      },
      backgroundColor: {
        colorcard: '#101935',
        dropdown: '#212C4D',
      },
    },
  },
  plugins: [
    require('@tailwindcss/forms'),
    // Iconify plugin
    addDynamicIconSelectors(),
  ],
  corePlugins: {
    preflight: true,
  },
}
