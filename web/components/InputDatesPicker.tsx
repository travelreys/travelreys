import React, { FC, useState } from 'react';
import { useMediaQuery } from 'usehooks-ts';
import _get from 'lodash/get';
import _isEmpty from 'lodash/isEmpty';

import { DayPicker, DateRange, SelectRangeEventHandler } from 'react-day-picker';
import { CalendarDaysIcon } from '@heroicons/react/24/outline'
import { InputDatesPickerCss } from '../styles/global';
import { printFromDateFromRange, printToDateFromRange } from '../utils/dates';
import DatesPicker from './DatesPicker';

interface InputDatesPicketProps {
  dates?: DateRange
  onSelect: any
}

const InputDatesPicker: FC<InputDatesPicketProps> = (props: InputDatesPicketProps) => {

  const [isOpen, setIsOpen] = useState(false);

  const startInputValue = printFromDateFromRange(props.dates, "y-MM-dd");
  const endInputValue = printToDateFromRange(props.dates, "y-MM-dd");

  let value = startInputValue;
  if (!_isEmpty(endInputValue)) {
    value = `${startInputValue} - ${endInputValue}`
  }

  // Renderers

  return (
    <div
      className="mb-4"
      onBlur={(e) => {
        if (!e.currentTarget.contains(e.relatedTarget) && isOpen) {
          setIsOpen(false);
        }
      }}
    >
      <div className={InputDatesPickerCss.Ctn}>
        <span className={InputDatesPickerCss.Label}>
          <CalendarDaysIcon className={InputDatesPickerCss.Icon} />
          &nbsp;Dates
        </span>
        <input
          onClick={() => {setIsOpen(true)}}
          type="text"
          value={value}
          onChange={() => { }}
          className={InputDatesPickerCss.Input}
        />
      </div>
      <DatesPicker
        dates={props.dates}
        onSelect={props.onSelect}
        isOpen={isOpen}
      />
    </div>
  );
}

export default InputDatesPicker;
