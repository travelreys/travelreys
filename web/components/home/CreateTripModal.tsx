import React, { useState, FC } from 'react';
import classNames from 'classnames';
import _get from "lodash/get";
import { format } from 'date-fns';

import { DayPicker } from 'react-day-picker';
import { CalendarDaysIcon } from '@heroicons/react/24/solid'


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
      <div className="flex mb-4">
        <span className="inline-flex font-bold items-center px-3 text-sm bg-gray-50 text-gray-900 border border-r-0 border-gray-300 rounded-l-md dark:bg-gray-600 dark:text-gray-400 dark:border-gray-600">
          Where to?
        </span>
        <input
          type="text"
          className={classNames(
            "bg-gray-50",
            "block",
            "border-gray-300",
            "border",
            "flex-1",
            "focus:border-blue-500",
            "focus:ring-blue-500",
            "min-w-0",
            "p-2.5",
            "rounded-none",
            "rounded-r-lg",
            "text-gray-900",
            "text-sm",
            "w-full",
          )}
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

    // HTML Classes
    const inputClasses = classNames(
      "bg-gray-50",
      "block",
      "border-gray-300",
      "border",
      "flex-1",
      "focus:border-blue-500",
      "focus:ring-blue-500",
      "min-w-0",
      "p-2.5",
      "rounded-none",
      "rounded-r-lg",
      "text-gray-900",
      "text-sm",
      "w-full",
    );
    return (
      <div className="flex justify-between">
        <div className="flex w-72">
          <span className="inline-flex font-bold items-center px-3 text-sm bg-gray-50 text-gray-900 border border-r-0 border-gray-300 rounded-l-md dark:bg-gray-600 dark:text-gray-400 dark:border-gray-600">
            <CalendarDaysIcon className='inline align-bottom h-5 w-5 text-gray-500' />
            &nbsp;
            Start
          </span>
          <input
            type="text"
            value={startInputValue}
            onChange={() => {}}
            onClick={dateInputOnClick}
            className={inputClasses}
            placeholder="start date ..."
          />
        </div>
        <div className="flex w-72">
          <span className="inline-flex font-bold items-center px-3 text-sm bg-gray-50 text-gray-900 border border-r-0 border-gray-300 rounded-l-md dark:bg-gray-600 dark:text-gray-400 dark:border-gray-600">
            <CalendarDaysIcon className='inline align-bottom h-5 w-5 text-gray-500' />
            &nbsp;
            End
          </span>
          <input
            type="text"
            value={endInputValue}
            onChange={() => {}}
            onClick={dateInputOnClick}
            className={inputClasses}
            placeholder="end date ..."
          />
        </div>
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
            numberOfMonths={2}
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
