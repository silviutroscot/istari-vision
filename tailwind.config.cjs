module.exports = {
    mode: 'jit',
    purge: {
        enabled: true,
        mode: 'all',
        preserveHtmlElements: false,
        options: {
            keyframes: true,
        },
        content: [
            './frontend/**/*.html',
            './frontend/**/*.js',
        ],
    },
    theme: {},
    darkMode: 'media', // or 'media' or 'class'
    plugins: [
        require('@tailwindcss/forms'),
    ],
}