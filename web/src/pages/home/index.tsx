import React, { FC, useState, useEffect } from 'react';
import { useNavigate } from "react-router-dom";
import { useTranslation } from 'react-i18next';
import _get from "lodash/get";
import _isEmpty from "lodash/isEmpty";
import {
  DateRange,
  SelectRangeEventHandler
} from 'react-day-picker';

import TripsAPI, { CreateTripResponse } from '../../apis/trips';

import Alert from '../../components/common/Alert';
import CreateTripModal from '../../features/home/CreateTripModal';
import Spinner from '../../components/common/Spinner';
import TripsContainer from '../../features/home/TripsContainer';
import { HomeCss } from '../../assets/styles/global';




interface TripsJumboProps {
  onCreateTripBtnClick: any,
}

const TripsJumbo: FC<TripsJumboProps> = (props: TripsJumboProps) => {

  const {t} = useTranslation();

  return (
    <div>
      <h1 className={HomeCss.TripJumboTitle}>
        {t("home.tripJumbo.title")}
      </h1>
      <button type="button"
        className={HomeCss.CreateNewTripBtn}
        onClick={props.onCreateTripBtnClick}
      >
        + {t("home.tripJumbo.createBtn")}
      </button>
    </div>
  );
}

const HomePage: FC = () => {
  const history = useNavigate();
  const { t } = useTranslation();

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
    .then((res: CreateTripResponse) => {
      history(`/trips/${res.id}`);
    })
    .catch(res => {
      setAlertProps({
        title: t("errors.unexpectedError"),
        message: res.error,
        status: "error"
      });
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
