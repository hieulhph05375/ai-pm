/** @type {import('tailwindcss').Config} */
export default {
    content: ['./src/**/*.{html,js,svelte,ts}'],
    theme: {
        extend: {
            colors: {
                "primary": "#1337ec",
                "background-light": "#f6f6f8",
                "background-dark": "#101322",
            },
            fontFamily: {
                "display": ["Inter", "sans-serif"],
                "outfit": ["Outfit", "sans-serif"]
            },
            borderRadius: {
                "DEFAULT": "0.5rem",
                "lg": "1rem",
                "xl": "1.5rem",
                "full": "9999px"
            },
        },
    },
    plugins: [],
}
