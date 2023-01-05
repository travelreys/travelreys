import React, { FC, useState } from 'react';
import Link from 'next/link';
import { format, parseJSON, isEqual } from 'date-fns';

import {
  Avatar,
  Box,
  Button,
  Grid,
  Typography
} from '@mui/material';
import _get from 'lodash/get';

// Helper Functions

const stringToColor = (string: string) => {
  let hash = 0;
  let i;

  /* eslint-disable no-bitwise */
  for (i = 0; i < string.length; i += 1) {
    hash = string.charCodeAt(i) + ((hash << 5) - hash);
  }

  let color = '#';

  for (i = 0; i < 3; i += 1) {
    const value = (hash >> (i * 8)) & 0xff;
    color += `00${value.toString(16)}`.slice(-2);
  }
  /* eslint-enable no-bitwise */
  return color;
}

// TripCard
interface TripCardProps {
  trip: any
}


const TripCard: FC<TripCardProps> = (props: TripCardProps) => {

  const [isHover, setIsHover] = useState(false);

  // Event Handlers
  const handleMouseEnter = () => {
     setIsHover(true);
  };
  const handleMouseLeave = () => {
     setIsHover(false);
  };

  // Renderers
  const renderCreatorAvatar = () => {
    const stringAvatar = (name: string) => {
      return {
        sx: {
          bgcolor: stringToColor(name),
          width: 24,
          height: 24,
          fontSize: 16,
          fontWeight: 900
        },
        children: name.toUpperCase(),
      };
    }

    return (
      <Avatar
        {...stringAvatar(_get(props.trip, "creator.memberEmail.0"))}
      />
    );
  }

  const renderTripDates = () => {
    const nullDate = parseJSON("0001-01-01T00:00:00Z");
    const startDate = parseJSON(props.trip.startDate);

    if (isEqual(startDate, nullDate)) {
      return;
    }

    const endDate = parseJSON(props.trip.endDate);
    if (isEqual(endDate, nullDate)) {
      return (
        <Typography variant='body2' color="secondary">
          {format(startDate, "MMM d, yy ")}
        </Typography>
      );
    }

    return (
      <Typography variant='body2' color="secondary">
        {format(startDate, "MMM d, yy ")} - {format(endDate, "MMM d, yy ")}
      </Typography>
    );
  }

  // Styles
  const cardStyle = {
    borderBottom: isHover ? "3px solid #D28088" : "none",
    borderRadius: "6px 0 0 0",
    paddingBottom: "0.5em",
    textDecoration: "none",
    color: "#A1A5D8",
  }

  const cardImgSx = {
    width: "100%",
    borderRadius: "6px",
    objectFit: "cover" as "cover",
    boxShadow: "1px 3px 18px 0px rgba(0,0,0,0.34)",
  }

  return (
    <Link
      href={`/trips/${props.trip.id}`}
      onMouseEnter={handleMouseEnter}
      onMouseLeave={handleMouseLeave}
      style={cardStyle}
    >
      <Box>
        <img
          style={cardImgSx}
          src="https://source.unsplash.com/collection/582860/660x440"
          alt=""
        />
      </Box>
      <Box>
        <Typography gutterBottom variant='body1'><b>{props.trip.name}</b></Typography>
        <Box sx={{display: "flex", justifyContent: "space-between"}}>
          {renderCreatorAvatar()}
          {renderTripDates()}
        </Box>
      </Box>
    </Link>
  );
}


// TripsContainer

interface TripsContainerProps {
  trips: Array<any>,
  onCreateTripBtnClick: any
}

const TripsContainer: FC<TripsContainerProps> = (props: TripsContainerProps) => {

  // Event Handlers

  // Renderers
  const renderTripsTable = () => {
    const cards = props.trips.map((trip: any) => {
      return (
        <Grid item xs={12} sm={4} md={3}>
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
      <Box sx={{
        display: "flex",
        justifyContent: "space-between",
        alignItems: "center",
      }}>
        <Typography
          variant="h4"
          color="info.main"
          gutterBottom
        >
          <b>Continue Planning</b>
        </Typography>
        <Box>
          <Button
              disableElevation
              variant="contained"
              onClick={props.onCreateTripBtnClick}
              sx={{
                textTransform: 'none',
                maxHeight: "2em",
                padding: "1em 1.25em",
                borderRadius: "2em",
                fontWeight: "900"
              }}
            >
              + New trip
            </Button>
        </Box>
      </Box>
      <br />
      {renderTripsTable()}
    </>
  );
}

export default TripsContainer;

