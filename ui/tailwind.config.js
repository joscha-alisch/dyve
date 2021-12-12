module.exports = {
  content: ['./src/**/*.{js,jsx,ts,tsx}', './public/index.html'],
  theme: {
    extend: {
      fontSize: {
        'tiny': '.6rem',
      }
    },
  },
  plugins: [
      require("tailwindcss-nested-groups"),
      require("@tailwindcss/forms")
  ],
}
