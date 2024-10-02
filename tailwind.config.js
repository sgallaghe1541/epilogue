/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./views/**/*.templ", "./internal/timeentry/*.templ"],
  theme: {
    extend: {
      gridTemplateColumns: {
        'mainGrid': 'theme(space.24) 1fr'
      },
      gridTemplateRows: {
        'mainGrid': '80px 1fr 32px'
      }
    },
  },
  plugins: [],
}

