import { createTheme } from "@mui/material/styles";

export const theme = createTheme({
  components: {
    MuiButton: {
      styleOverrides: {
        root: {
          textTransform: "none",
          fontWeight: 600,
          borderRadius: 12,
        },
      },
      variants: [
        {
          props: { variant: "contained", color: "primary" },
          style: {
            backgroundColor: "#2563eb",
            boxShadow: "0 4px 12px rgba(37, 99, 235, 0.2)",
            "&:hover": {
              backgroundColor: "#1d4ed8",
              boxShadow: "0 4px 12px rgba(37, 99, 235, 0.2)",
            },
          },
        },
      ],
    },
    MuiTabs: {
      styleOverrides: {
        root: {
          borderBottom: "1px solid #e2e8f0",
        },
        indicator: {
          backgroundColor: "#2563eb",
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
          "&.Mui-selected": {
            color: "#2563eb",
          },
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
  },
});
