import React, { FC, useState } from 'react';
import _get from "lodash/get";
import _sortBy from "lodash/sortBy";
import _isEmpty from "lodash/isEmpty";
import { ChevronUpIcon } from '@heroicons/react/24/solid';
import { ChevronDownIcon } from '@heroicons/react/24/outline';

interface DropdownProps {
  menu: any
  opts: any
  displayChevron?: boolean
  placement?: string
}

const Dropdown: FC<DropdownProps> = (props: DropdownProps) => {

  const [isActive, setIsActive] = useState(false);

  const opts = (
    <div className={"z-10 w-44 rounded-lg bg-white shadow block absolute right-0"}>
      <ul className={"z-10 w-44 rounded-lg bg-white shadow"}>
        {props.opts.map((opt: any, idx: number) => (
          <li
            key={idx}
            className={"block rounded-lg py-2 px-4"}>
            {opt}
          </li>
        ))}
      </ul>
    </div>
  )

  // Renderers

  const renderChevrons = () => {
    if (props.displayChevron) {
      return isActive ?
      <ChevronUpIcon className={"h-4 w-4 text-slate-700"} />
      : <ChevronDownIcon className={"h-4 w-4 text-slate-700"} />
    }
    return null;
  }

  return (
    <div className='relative'>
      <button
        type="button"
        className='flex items-center'
        onClick={() => { setIsActive(!isActive) }}
        onBlur={() => {
          setTimeout(() => {
            setIsActive(false);
          }, 150)
        }}
      >
        {props.menu}
        &nbsp;
        {renderChevrons()}
      </button>
      {isActive ? opts : null}
    </div>
  );

}

export default Dropdown;
