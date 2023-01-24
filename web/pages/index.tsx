import React, { useState } from 'react';
import { useRouter } from 'next/router';
import {
  DateRange,
  SelectRangeEventHandler
} from 'react-day-picker';
import _get from "lodash/get";
import _isEmpty from "lodash/isEmpty";

import TripsAPI from '../apis/trips';

import type { NextPageWithLayout } from './_app'
import Alert from '../components/Alert';
import CreateTripModal from '../components/home/CreateTripModal';
import Spinner from '../components/Spinner';
import TripsContainer from '../components/home/TripsContainer';
import TripsJumbo from '../components/home/TripsJumbo';


const HomePage: NextPageWithLayout = () => {
  const router = useRouter();

  // UI State
  const [newTripName, setNewTripName] = useState<string>("");
  const [newTripDates, setNewTripDates] = useState<DateRange>();
  const [isCreateModelOpen, setIsCreateModalOpen] = useState(false);
  const [alertProps, setAlertProps] = useState({} as any);

  let {data, error, isLoading} = TripsAPI.readTrips();

  if (isLoading) {
    return (<Spinner />);
  }

  const err = _get(error, "message") || _get(data, "error");
  if (err) {
    const props = {title: "Unexpected Error", message: err, status: "error"}
    return (<Alert {...props} />)
  }

  // Event Handlers

  const createTripModalOpenOnClick = () => {
    setIsCreateModalOpen(true);
  }

  const newTripNameOnChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setNewTripName(event.target.value);
  }

  const newTripDatesOnSelect: SelectRangeEventHandler = (range: DateRange | undefined) => {
    setNewTripDates(range);
  };

  const submitNewTripOnClick = () => {
    TripsAPI.createTrip(newTripName, newTripDates?.from, newTripDates?.to)
    .then((res: any) => {
      router.push(`/trips/${_get(res, 'data.tripPlan.id')}`);
    })
    .catch(error => {
      const errMsg = _get(error, "message") || _get(error.data, "error");
      setAlertProps({title: "Unexpected Error", message: errMsg, status: "error"});
    })
  }

  // Renderers
  const renderTrips = () => {
    const trips = _get(data, "tripPlans", []);
    if (trips.length === 0) {
      return (<TripsJumbo onCreateTripBtnClick={createTripModalOpenOnClick} />);
    }
    return (
      <TripsContainer
        trips={trips}
        onCreateTripBtnClick={createTripModalOpenOnClick}
      />
    );
  }

  return (
    <div>
      {!_isEmpty(alertProps) ? <Alert {...alertProps} /> : null}
      {renderTrips()}
      <CreateTripModal
        isOpen={isCreateModelOpen}
        onClose={() => setIsCreateModalOpen(false)}
        tripName={newTripName}
        tripNameOnChange={newTripNameOnChange}
        tripDates={newTripDates}
        tripDatesOnSelect={newTripDatesOnSelect}
        onSubmit={submitNewTripOnClick}
      />
    </div>
  );
}


export default HomePage;
