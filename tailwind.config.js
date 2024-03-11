/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./components/**/*.templ",
            "./pb_public/**/*.templ",
            "./app/**/*.templ",
            "./auth/**/*.templ",
            "./views/**/*.templ",
            "./lib/**/*.templ",
            "./pb_public/*.html"
  ],
  theme: {
    container: {
      center: true,
    },
    fontFamily: {
      oswald: ["Oswald"],
      poppins: ["Poppins-Thin"],
    },
    extend: {
    },
  },
  plugins: [],
}

