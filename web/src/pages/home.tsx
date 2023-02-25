import React, { FC, useState, useEffect } from 'react';
import { useNavigate } from "react-router-dom";
import { useTranslation } from 'react-i18next';
import _isEmpty from "lodash/isEmpty";
import {
  DateRange,
  SelectRangeEventHandler
} from 'react-day-picker';

import TripsAPI, { CreateResponse, ReadsResponse } from '../apis/trips';

import Alert from '../components/common/Alert';
import CreateModal from '../features/home/CreateModal';
import Spinner from '../components/common/Spinner';
import TripsContainer from '../features/home/TripsContainer';

interface TripsJumboProps {
  onCreateBtnClick: any,
}

const TripsJumbo: FC<TripsJumboProps> = (props: TripsJumboProps) => {
  const {t} = useTranslation();
  const css = {
    title: "text-4xl font-bold text-slate-700 mb-5",
    createBtn: "bg-indigo-500 font-bold px-5 py-2.5 rounded-md text-white hover:bg-indigo-70",
  }
  return (
    <div>
      <h1 className={css.title}>
        {t("home.tripJumbo.title")}
      </h1>
      <button type="button"
        className={css.createBtn}
        onClick={props.onCreateBtnClick}
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
    TripsAPI.list()
      .then((res: ReadsResponse) => {
        setTrips(res.trips);
        setAlertProps({});
        setIsLoading(false);
      })
      .catch(res => {
        setAlertProps({
          title: t("errors.unexpectedError"),
          message: res.error,
          status: "error"
        });
      })
  }, [t])

  // Event Handlers

  const createModalOpenOnClick = () => {
    setIsCreateModalOpen(true);
  }

  const newTripNameOnChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setNewTripName(event.target.value);
  }

  const newTripDatesOnSelect: SelectRangeEventHandler = (range?: DateRange) => {
    setNewTripDates(range);
  };

  const submitNewTripOnClick = () => {
    TripsAPI.create(newTripName, newTripDates?.from, newTripDates?.to)
    .then((res: CreateResponse) => {
      history(`/trips/${res.id}`);
    })
    .catch(res => {
      setAlertProps({
        title: t("errors.unexpectedError"),
        message: res.error,
        status: "error"
      });
    });
  }

  // Renderers
  const renderTrips = () => {
    if (trips.length === 0) {
      return (<TripsJumbo onCreateBtnClick={createModalOpenOnClick} />);
    }
    return (
      <TripsContainer
        trips={trips}
        onCreateBtnClick={createModalOpenOnClick}
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
      <CreateModal
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
