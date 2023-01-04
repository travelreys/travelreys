import * as React from 'react';

import AppBar from '@mui/material/AppBar';
import Container from '@mui/material/Container';
import Toolbar from '@mui/material/Toolbar';
import Typography from '@mui/material/Typography';

import PublicIcon from '@mui/icons-material/Public';

// Styles

const appBarStyle = {
  background: "none",
  color: "#3A4159",
  boxShadow: "none",
  borderBottom: "1px solid #e5e5e5",
}

export default function ButtonAppBar() {
  return (
    <AppBar sx={appBarStyle} position="static">
      <Container>
        <Toolbar>
          <PublicIcon
            color="primary"
            sx={{ display: { md: 'flex' }, mr: 1 }}
          />
          <Typography
            variant="h6"
            href="/"
            component="a"
            color="primary.main"
            sx={{
              flexGrow: 1,
              textDecoration: 'none',
              fontWeight: 900,
            }}
          >
            tiinyplanet
          </Typography>
        </Toolbar>
      </Container>
    </AppBar>
  );
}