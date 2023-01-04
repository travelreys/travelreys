import React, { useState, FC } from 'react';
import {
  DateRange,
  DayPicker,
  SelectRangeEventHandler
} from 'react-day-picker';

import Box from '@mui/material/Box';
import Button from '@mui/material/Button';
import Container from '@mui/material/Container';
import Modal from '@mui/material/Modal';
import Typography from '@mui/material/Typography';
import TextField from '@mui/material/TextField';
import InputAdornment from '@mui/material/InputAdornment';


interface CreateTripModalProps {
  isOpen: boolean,
  onClose: any,
  tripName: string | undefined,
  setTripName: any,
  tripDates: any,
  tripDatesOnSelect: any,
}

const CreateTripModal: FC<CreateTripModalProps> = (props: CreateTripModalProps) => {

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
    <Modal
      open={props.isOpen}
      onClose={props.onClose}
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
          onChange={props.setTripName}
          value={props.tripName}
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
          selected={props.tripDates}
          onSelect={props.tripDatesOnSelect}
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
  );
}


export default CreateTripModal;
