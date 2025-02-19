/** @type {import('tailwindcss').Config} */
export const content = [
  "./internal/omniview/templates/*.html",
];
export const theme = {
  container: {
    center: true,
    padding: {
      DEFAULT: "1rem",
      mobile: "2rem",
      tablet: "4rem",
      desktop: "5rem",
    },
  },
  extend: {},
};
export const colors = {
};
export const plugins = [require("@tailwindcss/forms"), require("@tailwindcss/typography")];
