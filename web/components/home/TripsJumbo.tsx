import React, { useState, FC } from 'react';

import Button from '@mui/material/Button';
import Container from '@mui/material/Container';
import Typography from '@mui/material/Typography';

interface TripsJumboProps {
  onCreateTripBtnClick: any,
}

const TripsJumbo: FC<TripsJumboProps> = (props: TripsJumboProps) => {
  return (
    <Container disableGutters>
      <Typography
          variant="h4"
          color="info.main"
          gutterBottom
        >
          <b>Plan your next adventure!</b>
        </Typography>
        <Button
          disableElevation
          variant="contained"
          onClick={props.onCreateTripBtnClick}
          sx={{
            textTransform: 'none',
            padding: "0.25em 1.5em",
            borderRadius: "2em",
            fontWeight: "900"
          }}
        >
          + Create new trip
        </Button>
    </Container>
  );
}

export default TripsJumbo;
