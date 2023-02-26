import React, { FC } from 'react';

////////////////////
// ListPin //
////////////////////


interface ListPinProps {
  icon?: any
  color: string
}

const ListPin: FC<ListPinProps> = (props: ListPinProps) => {
  return (
    <div className='relative -ml-2'>
      <span className='absolute right-2.5 top-2 text-base font-bold pointer-events-none'>
        {props.icon ? <props.icon className="fill-white stroke-2 w-3 h-3" /> : null}
      </span>
      <svg
        className="h-8 w-8 stroke-white stroke-2"
        style={{fill: props.color}}
        viewBox="0 0 24 24"
      >
        <path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z" />
      </svg>
    </div>
  );
}

export default ListPin;
