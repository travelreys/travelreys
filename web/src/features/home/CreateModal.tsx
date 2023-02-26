import React, { FC } from 'react';

import InputDatesPicker from '../../components/common/InputDatesPicker';
import Modal from '../../components/common/Modal';

interface CreateModalProps {
  isOpen: boolean,
  onClose: any,
  tripName?: string,
  tripNameOnChange: any,
  tripDates: any,
  tripDatesOnSelect: any,
  onSubmit: any
}

const CreateModal: FC<CreateModalProps> = (props: CreateModalProps) => {

  const css = {
    createModalCard: "bg-white rounded-lg px-4 pt-5 pb-4 sm:p-8 sm:pb-4",
    createTitle: "text-2xl font-bold text-center leading-6 text-slate-900 mb-6",
    tripNameCtn: "flex mb-4 border border-slate-200 rounded-lg",
    tripNameLabel: "inline-flex font-bold items-center px-3 text-sm text-slate-500 rounded-l-md",
    tripNameInput: "block flex-1 border-0 rounded-r-lg min-w-0 p-2.5 text-gray-900 text-sm w-full",
    tripDatesCtn: "flex w-full border border-slate-200 rounded-lg",
    tripDatesIcon: "inline align-bottom h-5 w-5 text-gray-500",
    tripDatesLabel: "inline-flex font-bold items-center px-3 text-sm text-slate-500 rounded-l-md",
    tripDatesInput: "block flex-1 min-w-0 p-2.5 border-0 rounded-none rounded-r-lg text-gray-900 text-sm w-full",
    createBtnsCtn: "bg-gray-50 px-4 pt-3 pb-5 rounded-b-lg sm:flex sm:flex-row-reverse sm:px-6",
    createBtn: "inline-flex w-full justify-center rounded-md border border-transparent bg-indigo-600 px-4 py-2 text-base font-medium text-white shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 sm:ml-3 sm:w-auto sm:text-sm",
    createCancelBtn: "mt-3 inline-flex w-full justify-center rounded-md border border-gray-300 bg-white px-4 py-2 text-base font-medium text-gray-700 shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm",
  }

  return (
    <Modal isOpen={props.isOpen}>
      <div className={css.createModalCard}>
        <h2 className={css.createTitle}>
          Create New Trip
        </h2>
        <div className={css.tripNameCtn}>
          <span className={css.tripNameLabel}>
            Where to?
          </span>
          <input
            type="text"
            className={css.tripNameInput}
            value={props.tripName}
            onChange={props.tripNameOnChange}
            placeholder="annual hiking vacation, business trip at new york ..."
          />
        </div>
        <InputDatesPicker
          onSelect={props.tripDatesOnSelect}
          dates={props.tripDates}
        />
      </div>
      <div className={css.createBtnsCtn}>
        <button type="button"
          className={css.createBtn}
          onClick={props.onSubmit}
        >
          Let's Go
        </button>
        <button type="button"
          className={css.createCancelBtn}
          onClick={() => { props.onClose() }}
        >
          Cancel
        </button>
      </div>
    </Modal>
  );
}

export default CreateModal;
