import React, { FC, useState } from 'react';
import _isEmpty from 'lodash/isEmpty';
import { DateRange } from 'react-day-picker';
import { CalendarDaysIcon } from '@heroicons/react/24/outline'

import DatesPicker from './DatesPicker';
import { fmt } from '../../lib/dates';


const css = {
  ctn: "flex w-full border border-slate-200 rounded-lg mr-2",
  icon: "inline align-bottom h-5 w-5 text-gray-500",
  label: "inline-flex bg-gray-200 font-bold items-center px-3 text-sm text-slate-500 rounded-l-md",
  input: "block flex-1 min-w-0 p-2.5 border-0 rounded-none rounded-r-lg text-gray-900 text-sm w-full",
}

interface InputDatesPicketProps {
  dates?: DateRange
  onSelect: any
  WrapperCss?: any
  CtnCss?: string
}

const InputDatesPicker: FC<InputDatesPicketProps> = (props: InputDatesPicketProps) => {

  const [isOpen, setIsOpen] = useState(false);

  const start = props.dates?.from ? fmt(props.dates?.from, "y-MM-dd"): undefined;
  const end = props.dates?.to ? fmt(props.dates.to, "y-MM-dd"): undefined;

  let value = start;
  if (!_isEmpty(end)) {
    value = `${start} - ${end}`
  }

  return (
    <div
      className={props.WrapperCss || "mb-4"}
      onBlur={(e) => {
        if (!e.currentTarget.contains(e.relatedTarget) && isOpen) {
          setIsOpen(false);
        }
      }}
    >
      <div className={props.CtnCss || css.ctn}>
        <span className={css.label}>
          <CalendarDaysIcon className={css.icon} />
          &nbsp;Dates
        </span>
        <input
          onClick={() => {setIsOpen(true)}}
          type="text"
          value={value}
          className={css.input}
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
