/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{js,ts,jsx,tsx}'],
  theme: {
    extend: {
      colors: {
        // Primary brand colors from palette
        brand: {
          dark: '#4A154B',      // Deep purple - primary dark
          DEFAULT: '#64C3EB',   // Light blue - primary
          light: '#8DD8F0',     // Lighter blue
          lighter: '#B5EBF5',   // Very light blue
        },
        // Success - Green
        success: {
          DEFAULT: '#5BB381',
          light: '#7CD19B',
          dark: '#4A8F68',
        },
        // Warning - Gold/Yellow
        warning: {
          DEFAULT: '#E3B34C',
          light: '#EBC872',
          dark: '#B8903D',
        },
        // Danger - Red/Pink
        danger: {
          DEFAULT: '#CE375C',
          light: '#D8607E',
          dark: '#A32C49',
        },
        // Neutral grays for UI
        neutral: {
          50: '#FAFAFA',
          100: '#F5F5F5',
          200: '#E5E5E5',
          300: '#D4D4D4',
          400: '#A3A3A3',
          500: '#737373',
          600: '#525252',
          700: '#404040',
          800: '#262626',
          900: '#171717',
          950: '#0A0A0A',
        },
        // Semantic aliases for easier usage
        primary: {
          50: '#EFF8FC',
          100: '#D9F1F9',
          200: '#B3E3F3',
          300: '#8DD8F0',
          400: '#64C3EB',
          500: '#4AA8D0',
          600: '#3D8DB8',
          700: '#3172A0',
          800: '#285A88',
          900: '#224A75',
          950: '#1A3558',
        },
        surface: {
          DEFAULT: '#FFFFFF',
          elevated: '#FAFAFA',
          sunken: '#F5F5F5',
        },
      },
      animation: {
        'scan': 'scan 0.3s ease-in-out',
        'fade-in': 'fadeIn 0.2s ease-out',
        'slide-up': 'slideUp 0.3s ease-out',
        'slide-down': 'slideDown 0.3s ease-out',
      },
      keyframes: {
        scan: {
          '0%': { transform: 'scaleX(0)' },
          '100%': { transform: 'scaleX(1)' },
        },
        fadeIn: {
          '0%': { opacity: '0' },
          '100%': { opacity: '1' },
        },
        slideUp: {
          '0%': { transform: 'translateY(10px)', opacity: '0' },
          '100%': { transform: 'translateY(0)', opacity: '1' },
        },
        slideDown: {
          '0%': { transform: 'translateY(-10px)', opacity: '0' },
          '100%': { transform: 'translateY(0)', opacity: '1' },
        },
      },
      boxShadow: {
        'card': '0 1px 3px 0 rgb(0 0 0 / 0.08), 0 1px 2px -1px rgb(0 0 0 / 0.08)',
        'card-hover': '0 4px 6px -1px rgb(0 0 0 / 0.1), 0 2px 4px -2px rgb(0 0 0 / 0.1)',
        'elevated': '0 10px 15px -3px rgb(0 0 0 / 0.1), 0 4px 6px -4px rgb(0 0 0 / 0.1)',
      },
      borderRadius: {
        'xl': '1rem',
        '2xl': '1.5rem',
      },
    },
  },
  plugins: [],
};