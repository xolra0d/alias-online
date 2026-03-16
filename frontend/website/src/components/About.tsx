import { Typography, Paper, Box, Divider } from "@mui/material";

export default function About() {
  return (
    <Box sx={{ maxWidth: 800, mx: "auto" }}>
      <Paper elevation={0} sx={{ p: { xs: 3, sm: 4 }, borderRadius: 4, border: "1px solid", borderColor: "divider" }}>
        <Typography variant="h3" gutterBottom sx={{ fontWeight: 800 }}>
          About Alias Online
        </Typography>
        <Typography variant="h6" color="text.secondary" paragraph>
          The classic word explanation game, now online and easier to play with friends.
        </Typography>
        
        <Divider sx={{ my: 3 }} />

        <Typography variant="h5" gutterBottom sx={{ fontWeight: 600 }}>
          How to Play
        </Typography>
        <Typography variant="body1" paragraph>
          1. Create a room with your preferred language and settings.
        </Typography>
        <Typography variant="body1" paragraph>
          2. Share the link with your friends.
        </Typography>
        <Typography variant="body1" paragraph>
          3. One player explains words while others try to guess them.
        </Typography>
        <Typography variant="body1" paragraph>
          4. Every correct guess gives points to both the guesser and the explainer!
        </Typography>
      </Paper>
    </Box>
  );
}
