/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./views/**/*.templ", "./internal/timeentry/*.templ"],
  theme: {
    extend: {
      gridTemplateColumns: {
        'main': 'theme(space.24) 1fr'
      },
      gridTemplateRows: {
        'main': 'theme(space.20) calc(100vh - theme(space.8) - theme(space.20)) theme(space.8)',
        'timeEntry': 'minmax(theme(space.16), max-content) minmax(theme(space.8), max-content) 1fr',
        'hours': 'minmax(min-content, max-content) minmax(min-content, max-content)'
      }
    },
  },
  plugins: [],
}

