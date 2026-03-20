/** @type {import('tailwindcss').Config} */
export default {
  content: ["./index.html", "./src/**/*.{ts,tsx}"],
  theme: {
    extend: {
      colors: {
        ink: "#141522",
        mist: "#F5F7FF",
        primary: {
          50: "#EEF1FF",
          100: "#DCE3FF",
          200: "#BBC9FF",
          300: "#8CA2FF",
          400: "#5F7CFF",
          500: "#3E63F4",
          600: "#2F4DCC",
          700: "#243B9F"
        },
        success: "#47C97E",
        warning: "#F6A44D",
        danger: "#F06464"
      },
      boxShadow: {
        panel: "0 24px 80px rgba(33, 56, 123, 0.14)",
        card: "0 18px 40px rgba(34, 60, 128, 0.10)"
      },
      fontFamily: {
        sans: ["Manrope", "Segoe UI", "system-ui", "sans-serif"]
      },
      backgroundImage: {
        halo: "radial-gradient(circle at top left, rgba(98, 123, 255, 0.2), transparent 32%), radial-gradient(circle at bottom right, rgba(82, 194, 255, 0.18), transparent 28%)"
      }
    }
  },
  plugins: []
};
