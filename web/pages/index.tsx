import React, { useState, FC } from 'react';
import {
  DateRange,
  DayPicker,
  SelectRangeEventHandler
} from 'react-day-picker';
import { format, isAfter, isBefore, isValid, parse } from 'date-fns';

import Box from '@mui/material/Box';
import Button from '@mui/material/Button';
import Container from '@mui/material/Container';
import Modal from '@mui/material/Modal';
import Typography from '@mui/material/Typography';
import TextField from '@mui/material/TextField';
import InputAdornment from '@mui/material/InputAdornment';



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

// Trips Jumbo

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

// Home

const Home = () => {

  // State
  const [trips, setTrips] = useState([] as any);
  const [newTripDates, setNewTripDates] = useState<DateRange>()

  const [isCreateModelOpen, setIsCreateModalOpen] = useState(false);



  // Event Handlers
  const createTripModalCloseOnClick = () => {
    setIsCreateModalOpen(false)
  }

  const createTripModalOpenOnClick = () => {
    setIsCreateModalOpen(true);
  }

  const newTripDatesOnSelect: SelectRangeEventHandler = (range: DateRange | undefined) => {
    setNewTripDates(range);
  };








  // Renderers
  const renderTrips = () => {
    if (trips.length === 0) {
      return (<TripsJumbo onCreateTripBtnClick={createTripModalOpenOnClick} />);
    }
    return <TripsContainer trips={trips} />;
  }

  // Styles
  const modalStyle = {
    display: "flex",
    justifyContent: "space-around",
    flexDirection: "column",
    position: 'absolute' as 'absolute',
    top: '50%',
    left: '50%',
    transform: 'translate(-50%, -50%)',
    width: "60vw",
    bgcolor: 'background.paper',
    borderRadius: "0.5em",
    padding: "2em 5em"
  };

  return (
    <Container maxWidth="lg" sx={{
      marginTop: "1em",
      padding: "0 48px!important"
    }}>
      {renderTrips()}
      <Modal
        open={isCreateModelOpen}
        onClose={createTripModalCloseOnClick}
        aria-labelledby="create new trip modal"
        aria-describedby="modal for creating a new trip"
      >
        <Box sx={modalStyle}>
          <Typography variant="h4" align="center" color="text.primary">
            <b>Create A New Trip</b>
          </Typography>
          <br />
          <TextField
            id="input-with-icon-textfield"
            placeholder='e.g annual family vacation, honeymoon to paris ...'
            InputProps={{
              sx: { borderRadius: "0.75em" },
              startAdornment: (
                <InputAdornment position="start">
                  <Typography variant="button" color="common.black">
                    <b>Where to? </b>
                  </Typography>
                </InputAdornment>
              ),
            }}
          />
          <DayPicker
            mode="range"
            numberOfMonths={2}
            pagedNavigation
            styles={{months: { margin: "0" }}}
            modifiersStyles={{
              selected: { background: "#AC8AC3" }
            }}
            selected={newTripDates}
            onSelect={newTripDatesOnSelect}
          />
          <br />
          <Box sx={{ display: "flex", width:"100%", justifyContent: "space-around" }}>
            <Button
              disableElevation
              variant="contained"
              sx={{
                textTransform: 'none',
                padding: "0.5em 2em",
                borderRadius: "2em",

              }}
            >
              <b>Lets Go !</b>
            </Button>
          </Box>
        </Box>
      </Modal>
    </Container>
  );
}


export default Home;
