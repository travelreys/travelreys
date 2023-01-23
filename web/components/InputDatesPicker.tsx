import React, { FC, useState } from 'react';
import { useMediaQuery } from 'usehooks-ts';
import _get from 'lodash/get';
import _isEmpty from 'lodash/isEmpty';
import { format } from 'date-fns';
import { DayPicker, DateRange, SelectRangeEventHandler } from 'react-day-picker';
import { CalendarDaysIcon } from '@heroicons/react/24/outline'
import { InputDatesPickerCss } from '../styles/global';
import { printFromDateFromRange, printToDateFromRange } from '../utils/dates';

interface InputDatesPicketProps {
  dates?: DateRange
  onSelect: any
}

const InputDatesPicker: FC<InputDatesPicketProps> = (props: InputDatesPicketProps) => {

  const [isOpen, setIsOpen] = useState(false);
  const matches = useMediaQuery('(min-width: 768px)');

  const startInputValue = printFromDateFromRange(props.dates, "y-MM-dd");
  const endInputValue = printToDateFromRange(props.dates, "y-MM-dd");

  let value = startInputValue;
  if (!_isEmpty(endInputValue)) {
    value = `${startInputValue} - ${endInputValue}`
  }

  // Event Handlers
  const datesOnSelect: SelectRangeEventHandler = (range: DateRange | undefined) => {
    props.onSelect(range);
  };

  // Renderers

  const renderDayPicker = () => {
    if (!isOpen) {
      return;
    }
    return (
      <div className='relative'>
        <div className='absolute bg-white border border-slate-200'>
          <DayPicker
            mode="range"
            numberOfMonths={matches ? 2 : 1}
            pagedNavigation
            styles={{ months: { margin: "0", display: "flex", justifyContent: "space-around" } }}
            modifiersStyles={{
              selected: { background: "#AC8AC3" }
            }}
            selected={props.dates}
            onSelect={datesOnSelect}
          />
        </div>
      </div>
    );
  }

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
      {renderDayPicker()}
    </div>
  );
}

export default InputDatesPicker;
