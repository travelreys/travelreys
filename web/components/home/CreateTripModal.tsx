import React, { useState, FC } from 'react';
import _get from "lodash/get";
import _isEmpty from "lodash/isEmpty";

import { CreateTripModalCss } from '../../styles/global';
import InputDatesPicker from '../InputDatesPicker';

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

  // Event Handlers
  const onClose = () => {
    props.onClose();
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
              <InputDatesPicker
                onSelect={props.tripDatesOnSelect}
                dates={props.tripDates}
              />
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
