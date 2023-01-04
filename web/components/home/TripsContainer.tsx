import React, { FC } from 'react';

import Container from '@mui/material/Container';
import Typography from '@mui/material/Typography';

// Trips

interface TripsContainerProps {
  trips: Array<any>,
}

const TripsContainer: FC<TripsContainerProps> = (props: TripsContainerProps) => {

  // Event Handlers

  // Renderers

  const renderTripsTable = () => {
    return (
      <>
        <Typography
          variant="h5"
          color="info.main"
          gutterBottom
        >
          <b>Continue Planning</b>
        </Typography>
        {/* Trips table here */}
      </>
    );
  }

  return (
    <Container disableGutters>
      {renderTripsTable()}
    </Container>
  );
}

export default TripsContainer;

