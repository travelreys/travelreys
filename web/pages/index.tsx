import React, { useState, FC } from 'react';
import {
  DateRange,
  DayPicker,
  SelectRangeEventHandler
} from 'react-day-picker';

import Container from '@mui/material/Container';
import TripsJumbo from '../components/home/TripsJumbo';
import TripsContainer from '../components/home/TripsContainer';
import CreateTripModal from '../components/home/CreateTripModal';

import TripsAPI from '../apis/trips';

// Home

const Home: FC = () => {

  // State
  const {trips, error, isLoading} = TripsAPI.readTrips();

  const [newTripName, setNewTripName] = useState<string>();
  const [newTripDates, setNewTripDates] = useState<DateRange>()

  const [isCreateModelOpen, setIsCreateModalOpen] = useState(false);

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

  // Renderers
  const renderTrips = () => {
    if (isLoading) {
      return (<div>Loading</div>)
    }
    if (error) {
      return (<div>{JSON.stringify(error)}</div>)
    }

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
      <CreateTripModal
        isOpen={isCreateModelOpen}
        onClose={createTripModalCloseOnClick}
        tripName={newTripName}
        setTripName={newTripNameOnUpdate}
        tripDates={newTripDates}
        tripDatesOnSelect={newTripDatesOnSelect}
      />
    </Container>
  );
}


export default Home;
