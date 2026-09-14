import { createTheme } from "@mui/material/styles";

export const theme = createTheme({
  palette: {
    mode: "light",
    primary: {
      main: "#5b3fd1",
      dark: "#3f289e",
      light: "#8c78e8",
    },
    secondary: {
      main: "#087e74",
    },
    background: {
      default: "#f6f4ef",
      paper: "#fffdfa",
    },
    text: {
      primary: "#1f1d2b",
      secondary: "#696575",
    },
  },
  shape: {
    borderRadius: 14,
  },
  typography: {
    fontFamily:
      'Inter, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif',
    h1: {
      fontSize: "clamp(2.25rem, 5vw, 4.6rem)",
      fontWeight: 760,
      letterSpacing: "-0.055em",
      lineHeight: 0.98,
    },
    h2: {
      fontWeight: 720,
      letterSpacing: "-0.035em",
    },
    button: {
      fontWeight: 700,
      textTransform: "none",
    },
  },
  components: {
    MuiButton: {
      defaultProps: { disableElevation: true },
      styleOverrides: {
        root: { borderRadius: 999, paddingInline: 20 },
      },
    },
    MuiCard: {
      styleOverrides: {
        root: {
          border: "1px solid rgba(31, 29, 43, 0.08)",
          boxShadow: "0 18px 55px rgba(40, 32, 67, 0.07)",
        },
      },
    },
  },
});
