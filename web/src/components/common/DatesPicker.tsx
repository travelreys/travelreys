import React, { FC } from 'react';
import { useMediaQuery } from 'usehooks-ts';


import {
  DayPicker,
  DateRange,
  SelectRangeEventHandler
} from 'react-day-picker';

interface InputDatesPickerProps {
  dates?: DateRange
  onSelect: any
  isOpen: boolean
}

const DatesPicker: FC<InputDatesPickerProps> = (props: InputDatesPickerProps) => {

  const matches = useMediaQuery('(min-width: 768px)');

  // Event Handlers
  const datesOnSelect: SelectRangeEventHandler = (range?: DateRange) => {
    props.onSelect(range);
  };

  // Renderers

  if (!props.isOpen) {
    return <></>;
  }

  return (
    <div className='relative z-50'>
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

export default DatesPicker;