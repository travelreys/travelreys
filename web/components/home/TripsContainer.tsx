import React, { FC } from 'react';
import { format, parse } from 'date-fns';

import {
  Box,
  Grid,
  Typography
} from '@mui/material';



interface TripCardProps {
  trip: any
}


const TripCard: FC<TripCardProps> = (props: TripCardProps) => {

  const cardSx = {
    padding: "15px",
    marginTop: "30px",
    borderRadius: "6px",
    boxShadow: "0 0 16px 1px rgba(0, 0, 0, 0.1)",
  }

  const cardImgBoxSx = {
    marginTop: "-30px",
    marginBottom: "15px",
  }

  const cardImgSx = {
    width: "100%",
    borderRadius: "6px",
    objectFit: "cover" as "cover",
    boxShadow: "0 16px 38px -12px rgb(0 0 0 / 56%), 0 4px 25px 0px rgb(0 0 0 / 12%), 0 8px 10px -5px rgb(0 0 0 / 20%)",
  }

  return (
    <Box sx={cardSx} color="primary.main">
      <Box sx={cardImgBoxSx}>
        <img
          style={cardImgSx}
          src="https://source.unsplash.com/collection/582860/660x440"
          alt=""
        />
      </Box>
      <Box>
        <Typography variant='body1'><b>{props.trip.name}</b></Typography>
        <Typography variant='body1'><b>{props.trip.name}</b></Typography>
      </Box>
    </Box>
  );
}






interface TripsContainerProps {
  trips: Array<any>,
}

const TripsContainer: FC<TripsContainerProps> = (props: TripsContainerProps) => {

  // Event Handlers

  console.log(props.trips)

  // Renderers
  const renderTripsTable = () => {
    const cards = props.trips.map((trip: any) => {
      return (
        <Grid item xs={6} md={4}>
          <TripCard trip={trip} key={trip.id} />
        </Grid>
      );
    })

    return (
      <Grid container spacing={2}>
        {cards}
      </Grid>
    );
  }

  return (
    <>
      <Typography
        variant="h5"
        color="info.main"
        gutterBottom
      >
        <b>Continue Planning</b>
      </Typography>
      {renderTripsTable()}
    </>
  );
}

export default TripsContainer;

