import React, { FC, useState } from 'react';
import { useMediaQuery } from 'usehooks-ts';
import _get from 'lodash/get';
import _isEmpty from 'lodash/isEmpty';
import { DayPicker, DateRange, SelectRangeEventHandler } from 'react-day-picker';

interface InputDatesPicketProps {
  dates?: DateRange
  onSelect: any
  isOpen: boolean
}

const DatesPicker: FC<InputDatesPicketProps> = (props: InputDatesPicketProps) => {

  const matches = useMediaQuery('(min-width: 768px)');

  // Event Handlers
  const datesOnSelect: SelectRangeEventHandler = (range: DateRange | undefined) => {
    props.onSelect(range);
  };

  // Renderers

  if (!props.isOpen) {
    return <></>;
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

export default DatesPicker;
