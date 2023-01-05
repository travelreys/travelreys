import { Source_Sans_Pro, Roboto } from '@next/font/google'
import { createTheme } from '@mui/material/styles';
import { fontSize } from '@mui/system';

const ssp = Source_Sans_Pro({
  weight: ['300', '400', '600', '700', '900'],
  style: ['normal', 'italic'],
  subsets: ['latin'],
  fallback: ['Helvetica', 'Arial', 'sans-serif'],
})

// Create a theme instance.
const theme = createTheme({
  palette: {
    primary: {
      main: "#A1A5D8",
      contrastText: "#f6f6f7"
    },
    secondary: {
      main: "#D28088",
    },
    info: {
      main: "#505b77",
    },
    success: {
      main: "#66ac79",
    },
    warning: {
      main: "#e39c41"
    },
    error: {
      main:"#f44336",
    },
    background: {
      paper: "#f6f6f7"
    },
    text: {
      primary: "#505b77",
    }
  },
  typography: {
    fontFamily: ssp.style.fontFamily,
    fontSize: 18,
  },
});

export default theme;