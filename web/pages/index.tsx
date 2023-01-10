import React, { ReactElement, useState, FC } from 'react';
import { useRouter } from 'next/router';
import {
  DateRange,
  SelectRangeEventHandler
} from 'react-day-picker';
import _get from "lodash/get";

import TripsAPI from '../apis/trips';

import type { NextPageWithLayout } from './_app'
import Layout from "../components/layouts/Layout";
import Alert from '../components/Alert';
import CreateTripModal from '../components/home/CreateTripModal';
import TripsContainer from '../components/home/TripsContainer';
import TripsJumbo from '../components/home/TripsJumbo';


const HomePage: NextPageWithLayout = () => {
  const router = useRouter();

  // State
  const [newTripName, setNewTripName] = useState<string>("");
  const [newTripDates, setNewTripDates] = useState<DateRange>()

  // UI State
  const [isCreateModelOpen, setIsCreateModalOpen] = useState(false);
  const [alertProps, setAlertProps] = useState({} as any);

  let {data, error, isLoading} = TripsAPI.readTrips();

  // Event Handlers
  const createTripModalCloseOnClick = () => {
    setIsCreateModalOpen(false)
  }

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
    setAlertProps({});

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
  const renderAlert = () => {
    if (_get(alertProps, "status") === "error") {
      return (<Alert {...alertProps} />)
    }
  }
  const renderTrips = () => {
    if (isLoading) {
      return (<div>Loading</div>)
    }

    const err = _get(error, "message") || _get(data, "error");
    if (err) {
      const errMsg = _get(error, "message") || _get(error.data, "error");
      setAlertProps({title: "Unexpected Error", message: errMsg, status: "error"});
      return;
    }

    const trips = _get(data, "tripPlans", []);
    if (trips.length === 0) {
      return (<TripsJumbo onCreateTripBtnClick={createTripModalOpenOnClick} />);
    }
    return <TripsContainer trips={trips} onCreateTripBtnClick={createTripModalOpenOnClick} />;
  }

  return (
    <div>
      {renderAlert()}
      {renderTrips()}
      <CreateTripModal
        isOpen={isCreateModelOpen}
        onClose={createTripModalCloseOnClick}
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
