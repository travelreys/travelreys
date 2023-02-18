import React, { FC } from 'react';
import _get from "lodash/get";
import _isEmpty from "lodash/isEmpty";

import { CreateTripModalCss, ModalCss } from '../../assets/styles/global';
import InputDatesPicker from '../../components/common/InputDatesPicker';


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
    <div className={ModalCss.Container}>
      <div className={ModalCss.Inset}></div>
      <div className={ModalCss.Content}>
        <div className={ModalCss.ContentContainer}>
          <div className={ModalCss.ContentCard}>
            <div className={CreateTripModalCss.CreateModalCard}>
              <h2 className={CreateTripModalCss.CreateTripTitle}>
                Create New Trip
              </h2>
              {renderTripNameInput()}
              <InputDatesPicker
                onSelect={props.tripDatesOnSelect}
                dates={props.tripDates}
              />
            </div>
            <div className={CreateTripModalCss.CreateTripBtnsCtn}>
              <button type="button"
                className={CreateTripModalCss.CreateTripBtn}
                onClick={props.onSubmit}
              >
                Let's Go
              </button>
              <button type="button"
                className={CreateTripModalCss.CreateTripCancelBtn}
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
