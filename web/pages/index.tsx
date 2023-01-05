import React, { useState, FC } from 'react';
import {
  DateRange,
  SelectRangeEventHandler
} from 'react-day-picker';

import _get from "lodash/get";

import Alert from '@mui/material/Alert';
import AlertTitle from '@mui/material/AlertTitle';
import Container from '@mui/material/Container';
import TripsJumbo from '../components/home/TripsJumbo';
import TripsContainer from '../components/home/TripsContainer';
import CreateTripModal from '../components/home/CreateTripModal';

import TripsAPI from '../apis/trips';

// Home

const Home: FC = () => {

  // State
  const [newTripName, setNewTripName] = useState<string>("");
  const [newTripDates, setNewTripDates] = useState<DateRange>()

  // UI State
  const [errorMsg, setErrorMsg] = useState<string>();
  const [isCreateModelOpen, setIsCreateModalOpen] = useState(false);

  let {data, error, isLoading} = TripsAPI.readTrips();

  // Event Handlers
  const createTripModalCloseOnClick = () => {
    setIsCreateModalOpen(false)
  }

  const createTripModalOpenOnClick = () => {
    setIsCreateModalOpen(true);
  }

  const newTripNameOnUpdate = (event: React.ChangeEvent<HTMLInputElement>) => {
    setNewTripName(event.target.value);
  }

  const newTripDatesOnSelect: SelectRangeEventHandler = (range: DateRange | undefined) => {
    setNewTripDates(range);
  };

  const submitNewTripOnClick = () => {
    TripsAPI.createTrip(newTripName, newTripDates?.from, newTripDates?.to)
    .then(res => {
      // go to trips page
    })
    .catch(error => {
      const errMsg = _get(error, "message") || _get(error.data, "error");
      setErrorMsg(errMsg);
    })
  }

  // Renderers
  const renderTrips = () => {
    if (isLoading) {
      return (<div>Loading</div>)
    }

    const err = _get(error, "message") || _get(data, "error");
    if (err) {
      setErrorMsg(err);
      return;
    }

    const trips = _get(data, "tripPlans", []);
    if (trips.length === 0) {
      return (<TripsJumbo onCreateTripBtnClick={createTripModalOpenOnClick} />);
    }
    return <TripsContainer trips={trips} onCreateTripBtnClick={createTripModalOpenOnClick} />;
  }

  const renderErrorMsg = () => {
    if (!errorMsg) {
      return;
    }
    return (
      <Alert severity="error">
        <AlertTitle><b>Error</b></AlertTitle>
        {errorMsg}
      </Alert>
    );
  }

  return (
    <Container maxWidth="lg" sx={{
      marginTop: "1em",
      padding: "0 48px!important"
    }}>
      {renderErrorMsg()}
      {renderTrips()}
      <CreateTripModal
        isOpen={isCreateModelOpen}
        onClose={createTripModalCloseOnClick}
        tripName={newTripName}
        setTripName={newTripNameOnUpdate}
        tripDates={newTripDates}
        tripDatesOnSelect={newTripDatesOnSelect}
        onSubmit={submitNewTripOnClick}
      />
    </Container>
  );
}


export default Home;
