/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./views/**/*.templ", "./internal/timeentry/*.templ"],
  theme: {
    extend: {
      gridTemplateColumns: {
        'main': 'theme(space.24) 1fr'
      },
      gridTemplateRows: {
        'main': '80px 1fr 32px',
        'timeEntry': 'theme(space.16) theme(space.8) 1fr'
      }
    },
  },
  plugins: [],
}

