import React, { FC } from 'react';
import _get from "lodash/get";
import _isEmpty from "lodash/isEmpty";

import { CreateModalCss, ModalCss } from '../../assets/styles/global';
import InputDatesPicker from '../../components/common/InputDatesPicker';


interface CreateModalProps {
  isOpen: boolean,
  onClose: any,
  tripName: string | undefined,
  tripNameOnChange: any,
  tripDates: any,
  tripDatesOnSelect: any,
  onSubmit: any
}

const CreateModal: FC<CreateModalProps> = (props: CreateModalProps) => {

  // Event Handlers
  const onClose = () => {
    props.onClose();
  }

  // Renderers
  const renderTripNameInput = () => {
    return (
      <div className={CreateModalCss.TripNameCtn}>
        <span className={CreateModalCss.TripNameLabel}>
          Where to?
        </span>
        <input
          type="text"
          className={CreateModalCss.TripNameInput}
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
            <div className={CreateModalCss.CreateModalCard}>
              <h2 className={CreateModalCss.CreateTitle}>
                Create New Trip
              </h2>
              {renderTripNameInput()}
              <InputDatesPicker
                onSelect={props.tripDatesOnSelect}
                dates={props.tripDates}
              />
            </div>
            <div className={CreateModalCss.CreateBtnsCtn}>
              <button type="button"
                className={CreateModalCss.CreateBtn}
                onClick={props.onSubmit}
              >
                Let's Go
              </button>
              <button type="button"
                className={CreateModalCss.CreateCancelBtn}
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


export default CreateModal;
