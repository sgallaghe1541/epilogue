/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./views/**/*.templ", "./internal/timeentry/*.templ"],
  theme: {
    extend: {
      gridTemplateColumns: {
        'main': 'theme(space.24) 1fr',
        'timeEntry': 'theme(space.60) repeat(2, minmax(theme(space.20), theme(space.20))) repeat(5, minmax(theme(space.24), 1fr)) theme(space.20)',
        'messages': 'minmax(min-content, 1fr) theme(space.8)'
      },
      gridTemplateRows: {
        'main': 'theme(space.20) calc(100vh - theme(space.8) - theme(space.20)) theme(space.8)',
        'timeEntry': 'minmax(min-content, max-content) 1fr',
        'hours': 'repeat(4, minmax(min-content, max-content))'
      }
    },
  },
  plugins: [],
}

