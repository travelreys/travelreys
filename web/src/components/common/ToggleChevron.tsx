import React, { FC } from 'react';
import {
  ChevronDownIcon,
  ChevronUpIcon,
} from '@heroicons/react/24/solid'

interface ToggleChevronProps {
  onClick: () => void
  isHidden: boolean
}

const ToggleChevron: FC<ToggleChevronProps> = (props: ToggleChevronProps) => {
  return (
    <button
      type="button"
      className="mr-2"
      onClick={() => {props.onClick()}}
    >
    {props.isHidden ? <ChevronUpIcon className='h-4 w-4' />
      : <ChevronDownIcon className='h-4 w-4'/>}
    </button>
  );
}

export default ToggleChevron;
