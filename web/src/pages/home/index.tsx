import React, { FC, useState, useEffect } from 'react';
import { useNavigate } from "react-router-dom";
import _get from "lodash/get";
import _isEmpty from "lodash/isEmpty";
import {
  DateRange,
  SelectRangeEventHandler
} from 'react-day-picker';

import TripsAPI from '../../apis/trips';

import Alert from '../../components/Alert';
import CreateTripModal from '../../components/home/CreateTripModal';
import Spinner from '../../components/Spinner';
import TripsContainer from '../../components/home/TripsContainer';
import TripsJumbo from '../../components/home/TripsJumbo';


const HomePage: FC = () => {
  const history = useNavigate();

  // UI State
  const [isLoading, setIsLoading] = useState(false);
  const [trips, setTrips] = useState([] as any);

  const [newTripName, setNewTripName] = useState<string>("");
  const [newTripDates, setNewTripDates] = useState<DateRange>();
  const [isCreateModelOpen, setIsCreateModalOpen] = useState(false);
  const [alertProps, setAlertProps] = useState({} as any);

  useEffect(() => {
    setIsLoading(true);
    TripsAPI.readTrips()
    .then((res) => {
      setTrips(_get(res, "tripPlans", []));
      setIsLoading(false);
    })
  }, [])

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
      history(`/trips/${_get(res, 'data.tripPlan.id')}`);
    })
    .catch(error => {
      const errMsg = _get(error, "message") || _get(error.data, "error");
      setAlertProps({title: "Unexpected Error", message: errMsg, status: "error"});
    })
  }

  // Renderers
  const renderTrips = () => {
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

  if (isLoading) {
    return (<Spinner />);
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
