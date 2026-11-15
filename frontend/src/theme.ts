import { createTheme } from "@mui/material/styles";

export const theme = createTheme({
  palette: {
    primary: {
      main: "#2563eb",
      dark: "#1d4ed8",
    },
  },
  components: {
    MuiButton: {
      styleOverrides: {
        root: {
          textTransform: "none",
          fontWeight: 600,
          borderRadius: 12,
        },
        sizeSmall: ({ theme }) => ({
          minWidth: 0,
          fontSize: "0.875rem",
          padding: "4px 8px",
          [theme.breakpoints.down("sm")]: {
            fontSize: "0.75rem",
            padding: "2px 6px",
          },
        }),
      },
      variants: [
        {
          props: { variant: "contained", color: "primary" },
          style: {
            boxShadow: "0 4px 12px rgba(37, 99, 235, 0.2)",
            "&:hover": {
              boxShadow: "0 4px 12px rgba(37, 99, 235, 0.2)",
            },
          },
        },
      ],
    },
    MuiIconButton: {
      styleOverrides: {
        sizeSmall: ({ theme }) => ({
          padding: 4,
          [theme.breakpoints.down("sm")]: {
            padding: 2,
          },
        }),
      },
    },
    MuiTabs: {
      styleOverrides: {
        root: {
          borderBottom: "1px solid #e2e8f0",
        },
        indicator: {
          height: 3,
          borderRadius: "3px 3px 0 0",
        },
      },
    },
    MuiTab: {
      styleOverrides: {
        root: {
          textTransform: "none",
          fontWeight: 600,
          color: "#64748b",
          fontSize: "0.95rem",
        },
      },
    },
    MuiDialog: {
      styleOverrides: {
        paper: {
          borderRadius: 24,
          padding: 16,
          background: "rgba(255, 255, 255, 0.85)",
          backdropFilter: "blur(20px)",
          border: "1px solid rgba(255, 255, 255, 0.4)",
          boxShadow: "0 24px 64px rgba(0,0,0,0.15)",
        },
      },
    },
    MuiOutlinedInput: {
      styleOverrides: {
        root: {
          borderRadius: 12,
        },
      },
    },
    MuiMenu: {
      styleOverrides: {
        paper: {
          borderRadius: 16,
          marginTop: 4,
        },
      },
    },
    MuiMenuItem: {
      styleOverrides: {
        root: {
          borderRadius: 8,
        },
      },
    },
  },
});
