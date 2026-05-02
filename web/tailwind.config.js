import forms from '@tailwindcss/forms';

/** @type {import('tailwindcss').Config} */
export default {
  content: ['./src/**/*.{html,js,svelte,ts}'],
  theme: {
    extend: {
      colors: {
        neon: {
          cyan: '#3cf2ff',
          pink: '#ff4fd8',
          lime: '#8dff65',
          amber: '#ffb84d',
          bg: '#070b14',
          panel: '#10192b'
        }
      },
      boxShadow: {
        glow: '0 0 35px rgba(60, 242, 255, 0.18)'
      }
    }
  },
  plugins: [forms]
};
