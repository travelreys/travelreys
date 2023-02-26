import React, { FC, useRef, useState } from 'react';
import { ChevronUpIcon } from '@heroicons/react/24/solid';
import { ChevronDownIcon } from '@heroicons/react/24/outline';

import useOutsideAlerter from '../../hooks/useOutsideAlerter';


const css = {
  optsCtn: "z-10 w-44 rounded-lg bg-white shadow block absolute right-0",
  optsList: "z-10 w-44 rounded-lg bg-white shadow",
  optItem: "block rounded-lg py-2 px-4",
  btn: "flex items-center",
  chevron: "h-4 w-4 text-slate-700",
}

interface DropdownProps {
  menu: any
  opts: any
  displayChevron?: boolean
  placement?: string
}

const Dropdown: FC<DropdownProps> = (props: DropdownProps) => {

  const [isActive, setIsActive] = useState(false);
  const wrapperRef = useRef(null);
  useOutsideAlerter(wrapperRef, () => {setIsActive(false)});

  const renderOpts = () => {
    return (
      <div className={css.optsCtn}>
        <ul className={css.optsList}>
          {props.opts.map((opt: any, idx: number) => (
            <li key={idx} className={css.optItem}>{opt}</li>
          ))}
        </ul>
      </div>
    );
  }

  // Renderers

  const renderChevrons = () => {
    if (props.displayChevron) {
      return isActive
        ? <ChevronUpIcon className={css.chevron} />
        : <ChevronDownIcon className={css.chevron} />
    }
    return null;
  }

  return (
    <div ref={wrapperRef} className='relative'>
      <button
        type="button"
        className={css.btn}
        onClick={() => { setIsActive(!isActive) }}
      >
        {props.menu}
        &nbsp;
        {renderChevrons()}
      </button>
      {isActive ? renderOpts() : null}
    </div>
  );

}

export default Dropdown;
