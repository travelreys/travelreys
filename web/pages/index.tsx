import React, { useState, FC } from 'react';
import {
  DateRange,
  SelectRangeEventHandler
} from 'react-day-picker';

import _get from "lodash/get";

import TripsAPI from '../apis/trips';


const HomePage: FC = () => {

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
    return
    // if (trips.length === 0) {
    //   return (<TripsJumbo onCreateTripBtnClick={createTripModalOpenOnClick} />);
    // }
    // return <TripsContainer trips={trips} onCreateTripBtnClick={createTripModalOpenOnClick} />;
  }

  return (
    <div>

      <p className="font-mono">The quick brown fox ...</p>
    </div>
  );
}


export default HomePage;
