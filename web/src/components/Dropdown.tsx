import React, { FC, useState } from 'react';
import _get from "lodash/get";
import _sortBy from "lodash/sortBy";
import _isEmpty from "lodash/isEmpty";

interface DropdownProps {
  menu: any
  opts: any
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

  return (
    <div className='relative'>
      <button
        type="button"
        onClick={() => { setIsActive(!isActive) }}
        onBlur={() => {
          setTimeout(() => {
            setIsActive(false);
          }, 150)
        }}
      >
        {props.menu}
      </button>
      {isActive ? opts : null}
    </div>
  );

}

export default Dropdown;
