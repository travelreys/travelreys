import React, { useState, FC } from 'react';
import classNames from 'classnames';

interface TripsJumboProps {
  onCreateTripBtnClick: any,
}

const TripsJumbo: FC<TripsJumboProps> = (props: TripsJumboProps) => {
  return (
    <div>
      <h1 className='text-4xl font-bold text-slate-700 mb-5'>
        Plan your next adventure!
      </h1>
      <button type="button"
        className={classNames(
          "bg-gradient-to-r",
          "font-medium",
          "from-indigo-500",
          "hover:bg-gradient-to-br",
          "mb-2",
          "mr-2",
          "px-5",
          "py-2.5",
          "rounded-md",
          "text-center",
          "text-sm",
          "text-white",
          "to-indigo-700",
          "via-indigo-600",
        )}
        onClick={props.onCreateTripBtnClick}
      >
        + Create new trip
      </button>
    </div>
  );
}

export default TripsJumbo;
