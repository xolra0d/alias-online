import type { PaletteMode } from "@mui/material";
import { alpha } from "@mui/material/styles";

export const getDesignTokens = (mode: PaletteMode) => {
  const isLight = mode === "light";

  const primaryMain = isLight ? "#B9771E" : "#F2C15B";
  const primaryLight = isLight ? "#E0A94B" : "#F7D789";
  const primaryDark = isLight ? "#8B5912" : "#D39B2B";

  const secondaryMain = isLight ? "#2C8B92" : "#57D5DA";
  const secondaryLight = isLight ? "#4EB8BF" : "#7CE5E8";
  const secondaryDark = isLight ? "#1D6469" : "#2AA8B0";

  const backgroundDefault = isLight ? "#F3E8D6" : "#0E1013";
  const backgroundPaper = isLight ? "#FFF8EF" : "#171A1F";
  const textPrimary = isLight ? "#201812" : "#F3E8D7";
  const textSecondary = isLight
    ? "rgba(32, 24, 18, 0.66)"
    : "rgba(243, 232, 215, 0.72)";
  const divider = isLight ? "rgba(63, 44, 25, 0.12)" : "rgba(242, 193, 91, 0.16)";
  const successMain = isLight ? "#2F7A4D" : "#67C58B";
  const warningMain = isLight ? "#C58A2D" : "#E1AA4A";
  const errorMain = isLight ? "#B45444" : "#E07A68";

  const headingFont = '"Josefin Sans", Georgia, serif';
  const bodyFont = '"Lora", Georgia, serif';

  return {
    palette: {
      mode,
      primary: {
        main: primaryMain,
        light: primaryLight,
        dark: primaryDark,
        contrastText: isLight ? "#1E160F" : "#0E1013",
      },
      secondary: {
        main: secondaryMain,
        light: secondaryLight,
        dark: secondaryDark,
        contrastText: isLight ? "#FFF8EF" : "#0E1013",
      },
      success: {
        main: successMain,
        contrastText: isLight ? "#FFF8EF" : "#0E1013",
      },
      warning: {
        main: warningMain,
        contrastText: isLight ? "#1E160F" : "#0E1013",
      },
      error: {
        main: errorMain,
        contrastText: isLight ? "#FFF8EF" : "#0E1013",
      },
      info: {
        main: secondaryMain,
        contrastText: isLight ? "#FFF8EF" : "#0E1013",
      },
      divider,
      background: {
        default: backgroundDefault,
        paper: backgroundPaper,
      },
      text: {
        primary: textPrimary,
        secondary: textSecondary,
      },
      action: {
        hover: alpha(primaryMain, isLight ? 0.08 : 0.14),
        selected: alpha(primaryMain, isLight ? 0.14 : 0.2),
        disabled: alpha(textPrimary, 0.38),
        disabledBackground: alpha(textSecondary, 0.1),
      },
    },
    shape: {
      borderRadius: 18,
    },
    typography: {
      fontFamily: bodyFont,
      h1: {
        fontFamily: headingFont,
        fontWeight: 700,
        letterSpacing: "-0.04em",
        lineHeight: 1.02,
      },
      h2: {
        fontFamily: headingFont,
        fontWeight: 700,
        letterSpacing: "-0.03em",
        lineHeight: 1.05,
      },
      h3: {
        fontFamily: headingFont,
        fontWeight: 700,
        letterSpacing: "-0.02em",
        lineHeight: 1.1,
      },
      h4: {
        fontFamily: headingFont,
        fontWeight: 700,
        letterSpacing: "-0.01em",
        lineHeight: 1.1,
      },
      h5: {
        fontFamily: headingFont,
        fontWeight: 700,
        letterSpacing: "-0.01em",
        lineHeight: 1.15,
      },
      h6: {
        fontFamily: headingFont,
        fontWeight: 700,
        letterSpacing: "0.01em",
        lineHeight: 1.2,
      },
      subtitle1: {
        fontFamily: headingFont,
        fontWeight: 600,
        letterSpacing: "0.06em",
      },
      subtitle2: {
        fontFamily: headingFont,
        fontWeight: 600,
        letterSpacing: "0.08em",
      },
      overline: {
        fontFamily: headingFont,
        fontWeight: 700,
        letterSpacing: "0.28em",
        textTransform: "uppercase",
      },
      button: {
        fontFamily: headingFont,
        fontWeight: 700,
        letterSpacing: "0.08em",
        textTransform: "none" as const,
      },
      body1: {
        lineHeight: 1.75,
      },
      body2: {
        lineHeight: 1.65,
      },
    },
    components: {
      MuiCssBaseline: {
        styleOverrides: {
          "html, body, #root": {
            minHeight: "100%",
          },
          body: {
            backgroundColor: backgroundDefault,
            backgroundImage: [
              `radial-gradient(circle at 12% 8%, ${alpha(primaryMain, isLight ? 0.18 : 0.12)} 0%, transparent 34%)`,
              `radial-gradient(circle at 88% 12%, ${alpha(secondaryMain, isLight ? 0.15 : 0.12)} 0%, transparent 30%)`,
              `radial-gradient(circle at 50% 104%, ${alpha(primaryMain, isLight ? 0.12 : 0.08)} 0%, transparent 38%)`,
              `repeating-linear-gradient(135deg, ${alpha(textPrimary, isLight ? 0.025 : 0.04)} 0 1px, transparent 1px 14px)`,
            ].join(", "),
            backgroundAttachment: "fixed",
            color: textPrimary,
            fontFamily: bodyFont,
            textRendering: "optimizeLegibility",
          },
          "a, a:visited": {
            color: "inherit",
          },
          "*::selection": {
            backgroundColor: alpha(primaryMain, isLight ? 0.22 : 0.28),
            color: textPrimary,
          },
        },
      },
      MuiAppBar: {
        styleOverrides: {
          root: {
            boxShadow: "none",
            backgroundImage: "none",
            backgroundColor: alpha(backgroundPaper, isLight ? 0.88 : 0.82),
            backdropFilter: "blur(20px)",
            borderBottom: `1px solid ${divider}`,
          },
        },
      },
      MuiToolbar: {
        styleOverrides: {
          root: {
            minHeight: 80,
          },
        },
      },
      MuiPaper: {
        styleOverrides: {
          root: {
            backgroundImage: "none",
            borderRadius: 18,
            border: `1px solid ${divider}`,
            boxShadow: "none",
            backgroundColor: backgroundPaper,
          },
        },
      },
      MuiButton: {
        styleOverrides: {
          root: {
            borderRadius: 12,
            boxShadow: "none",
            padding: "0.85rem 1.2rem",
          },
          containedPrimary: {
            backgroundColor: primaryMain,
            color: isLight ? "#1E160F" : "#0E1013",
            "&:hover": {
              backgroundColor: primaryDark,
              boxShadow: "none",
            },
          },
          containedSecondary: {
            backgroundColor: secondaryMain,
            color: isLight ? "#FFF8EF" : "#0E1013",
            "&:hover": {
              backgroundColor: secondaryDark,
            },
          },
          outlined: {
            borderColor: divider,
            backgroundColor: alpha(backgroundPaper, isLight ? 0.76 : 0.5),
          },
        },
      },
      MuiToggleButton: {
        styleOverrides: {
          root: {
            borderRadius: 12,
            borderColor: divider,
            textTransform: "none",
            fontFamily: headingFont,
            fontWeight: 700,
            letterSpacing: "0.08em",
            color: textSecondary,
            backgroundColor: alpha(backgroundPaper, isLight ? 0.7 : 0.42),
            "&.Mui-selected": {
              color: textPrimary,
              backgroundColor: alpha(primaryMain, isLight ? 0.14 : 0.22),
              borderColor: primaryMain,
              "&:hover": {
                backgroundColor: alpha(primaryMain, isLight ? 0.18 : 0.26),
              },
            },
          },
        },
      },
      MuiOutlinedInput: {
        styleOverrides: {
          root: {
            borderRadius: 14,
            backgroundColor: alpha(backgroundPaper, isLight ? 0.82 : 0.9),
            "& .MuiOutlinedInput-notchedOutline": {
              borderColor: divider,
            },
            "&:hover .MuiOutlinedInput-notchedOutline": {
              borderColor: alpha(primaryMain, 0.55),
            },
            "&.Mui-focused .MuiOutlinedInput-notchedOutline": {
              borderColor: primaryMain,
              borderWidth: 1,
            },
          },
          input: {
            fontFamily: bodyFont,
          },
        },
      },
      MuiInputLabel: {
        styleOverrides: {
          root: {
            fontFamily: headingFont,
            fontWeight: 600,
            letterSpacing: "0.08em",
          },
        },
      },
      MuiFormLabel: {
        styleOverrides: {
          root: {
            fontFamily: headingFont,
            fontWeight: 600,
            letterSpacing: "0.08em",
          },
        },
      },
      MuiFormControlLabel: {
        styleOverrides: {
          label: {
            fontFamily: headingFont,
            letterSpacing: "0.04em",
          },
        },
      },
      MuiChip: {
        styleOverrides: {
          root: {
            borderRadius: 10,
            fontFamily: headingFont,
            fontWeight: 700,
            letterSpacing: "0.08em",
          },
        },
      },
      MuiAlert: {
        styleOverrides: {
          root: {
            borderRadius: 16,
            border: `1px solid ${divider}`,
            backgroundImage: "none",
            alignItems: "center",
          },
        },
      },
      MuiDivider: {
        styleOverrides: {
          root: {
            borderColor: divider,
          },
        },
      },
      MuiListItem: {
        styleOverrides: {
          root: {
            borderRadius: 14,
          },
        },
      },
      MuiListItemText: {
        styleOverrides: {
          primary: {
            fontFamily: headingFont,
            fontWeight: 700,
          },
        },
      },
    },
  };
};
