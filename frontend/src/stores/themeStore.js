import { defineStore } from 'pinia'

export const useThemeStore = defineStore('theme', {
    state: () => ({
        isDarkMode: false
    }),

    actions: {
        toggleDarkMode() {
            this.isDarkMode = !this.isDarkMode

            // Apply theme to document
            if (this.isDarkMode) {
                document.documentElement.setAttribute('data-theme', 'dark')
            } else {
                document.documentElement.removeAttribute('data-theme')
            }

            // Save preference to localStorage
            localStorage.setItem('darkMode', this.isDarkMode)
        },

        initTheme() {
            // Check for saved preference or system preference
            const savedTheme = localStorage.getItem('darkMode')

            if (savedTheme !== null) {
                this.isDarkMode = savedTheme === 'true'
            } else {
                // Check system preference
                this.isDarkMode = window.matchMedia('(prefers-color-scheme: dark)').matches
            }

            // Apply initial theme
            if (this.isDarkMode) {
                document.documentElement.setAttribute('data-theme', 'dark')
            }
        }
    }
})

