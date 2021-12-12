module.exports = {
  purge: ['./src/**/*.{js,jsx,ts,tsx}', './public/index.html'],
  darkMode: false, // or 'media' or 'class'
  theme: {
    extend: {
      fontSize: {
        'tiny': '.6rem',
      }
    },
  },
  variants: {
    extend: {
      border: ["hover"],
      ring: ["hover"],
      display: ["group-hover"]
    },
  },
  plugins: [
      require("tailwindcss-nested-groups"),
      require("@tailwindcss/forms")
  ],
}
