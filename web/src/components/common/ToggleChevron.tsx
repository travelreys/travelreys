import React, { FC } from 'react';
import {
  ChevronDownIcon,
  ChevronUpIcon,
} from '@heroicons/react/24/solid'
import { CommonCss } from '../../assets/styles/global';

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
    {props.isHidden ? <ChevronUpIcon className={CommonCss.Icon} />
      : <ChevronDownIcon className={CommonCss.Icon}/>}
    </button>
  );
}

export default ToggleChevron;
