import React, { useState, FC } from 'react';
import { useMediaQuery } from 'usehooks-ts';
import { format } from 'date-fns';
import _get from "lodash/get";
import _isEmpty from "lodash/isEmpty";

import { DayPicker } from 'react-day-picker';
import { CalendarDaysIcon } from '@heroicons/react/24/solid'
import { CreateTripModalCss } from '../../styles/global';


interface CreateTripModalProps {
  isOpen: boolean,
  onClose: any,
  tripName: string | undefined,
  tripNameOnChange: any,
  tripDates: any,
  tripDatesOnSelect: any,
  onSubmit: any
}

const CreateTripModal: FC<CreateTripModalProps> = (props: CreateTripModalProps) => {

  const [isCalendarOpen, setIsCalendarOpen] = useState(false);
  const matches = useMediaQuery('(min-width: 768px)')


  // Event Handlers
  const onClose = () => {
    props.onClose();
  }

  const dateInputOnClick = (event: React.MouseEvent<HTMLInputElement>) => {
    setIsCalendarOpen(!isCalendarOpen)
  }

  // Renderers
  const renderTripNameInput = () => {
    return (
      <div className={CreateTripModalCss.TripNameCtn}>
        <span className={CreateTripModalCss.TripNameLabel}>
          Where to?
        </span>
        <input
          type="text"
          className={CreateTripModalCss.TripNameInput}
          value={props.tripName}
          onChange={props.tripNameOnChange}
          placeholder="annual hiking vacation, business trip at new york ..."
        />
      </div>
    );
  }

  const renderDatesInputs = () => {
    const startInputValue = _get(props.tripDates, "from") ? format(_get(props.tripDates, "from"), 'y-MM-dd') : "";
    const endInputValue = _get(props.tripDates, "to") ? format(_get(props.tripDates, "to"), 'y-MM-dd') : "";
    let value = startInputValue;
    if (!_isEmpty(endInputValue)) {
      value = `${startInputValue} - ${endInputValue}`
    }

    return (
      <div className={CreateTripModalCss.TripDatesCtn}>
        <span className={CreateTripModalCss.TripDatesLabel}>
          <CalendarDaysIcon className={CreateTripModalCss.TripDatesIcon} />
          &nbsp;
          Dates
        </span>
        <input
          type="text"
          value={value}
          onChange={() => {}}
          onClick={dateInputOnClick}
          className={CreateTripModalCss.TripDatesInput}
        />
      </div>
    );
  }

  const renderDayPicker = () => {
    if (!isCalendarOpen) {
      return;
    }
    return (
      <div className='relative'>
        <div className='absolute bg-white mt-2 border border-slate-200'>
          <DayPicker
            mode="range"
            numberOfMonths={matches ? 2 : 1}
            pagedNavigation
            styles={{ months: { margin: "0", display: "flex", justifyContent: "space-around" } }}
            modifiersStyles={{
              selected: { background: "#AC8AC3" }
            }}
          selected={props.tripDates}
          onSelect={props.tripDatesOnSelect}
          />
        </div>
      </div>
    );
  }


  if (!props.isOpen) {
    return <></>;
  }

  return (
    <div className="relative z-10" aria-labelledby="modal-title" role="dialog" aria-modal="true">
      <div className="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity"></div>
      <div className="fixed inset-0 z-10 overflow-y-auto">
        <div className="flex min-h-full items-center justify-center p-4 text-center sm:items-center sm:p-0">
          <div className="relative transform rounded-lg bg-white text-left shadow-xl transition-all w-11/12 sm:my-8 sm:w-full sm:max-w-2xl">
            <div className="bg-white px-4 pt-5 pb-4 sm:p-8 sm:pb-4">
              <h2 className="text-2xl text-center font-medium leading-6 text-slate-900 mb-6">Create New Trip</h2>
              {renderTripNameInput()}
              {renderDatesInputs()}
              {renderDayPicker()}
            </div>
            <div className="bg-gray-50 px-4 py-3 sm:flex sm:flex-row-reverse sm:px-6">
              <button type="button"
                className="inline-flex w-full justify-center rounded-md border border-transparent bg-indigo-600 px-4 py-2 text-base font-medium text-white shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 sm:ml-3 sm:w-auto sm:text-sm"
                onClick={props.onSubmit}
              >
                Let's Go
              </button>
              <button type="button"
                className="mt-3 inline-flex w-full justify-center rounded-md border border-gray-300 bg-white px-4 py-2 text-base font-medium text-gray-700 shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm"
                onClick={onClose}
              >
                Cancel
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}


export default CreateTripModal;
