import React, { FC } from 'react';
import classNames from 'classnames';
import Link from 'next/link';
import _get from 'lodash/get';
import { parseJSON, isEqual } from 'date-fns';

import { datesRenderer } from '../../utils/dates';
import Avatar from '../Avatar';


interface TripCardProps {
  trip: any
}

const TripCard: FC<TripCardProps> = (props: TripCardProps) => {

  // Renderers
  const renderCreatorAvatar = () => {

    return (<Avatar name={_get(props.trip, "creator.memberEmail")} placement="top" />);
  }

  const renderTripDates = () => {
    const nullDate = parseJSON("0001-01-01T00:00:00Z");
    const startDate = parseJSON(props.trip.startDate);
    const endDate = parseJSON(props.trip.endDate);


    if (isEqual(startDate, nullDate)) {
      return <div className='text-slate-500'>-</div>;
    }
    return (
      <p className='text-slate-500'>
        {datesRenderer(startDate, endDate)}
      </p>
    );
  }

  return (
    <Link href={`/trips/${props.trip.id}`}>
      <div className="bg-white rounded-lg shadow-md dark:bg-gray-800 dark:border-gray-700">
        <img
          className="rounded-t-lg"
          src="https://images.unsplash.com/photo-1469854523086-cc02fe5d8800?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1421&q=80"
          alt=""
        />
        <div className="p-5">
          <h5 className="mb-2 text-xl font-bold tracking-tight text-gray-900 dark:text-white">
            {props.trip.name}
          </h5>
          <div className='flex justify-between'>
            {renderTripDates()}
            {renderCreatorAvatar()}
          </div>
        </div>
      </div>
    </Link>
  );
}


// TripsContainer

interface TripsContainerProps {
  trips: Array<any>,
  onCreateTripBtnClick: any
}

const TripsContainer: FC<TripsContainerProps> = (props: TripsContainerProps) => {

  // Event Handlers

  // Renderers
  const renderTripsTable = () => {
    const cards = props.trips.map((trip: any) => {
      return <TripCard trip={trip} key={trip.id} />
    })
    return (
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-3">
        {cards}
      </div>
    );
  }

  return (
    <div>
      <div className='flex justify-between flex-col sm:flex-row items-center mb-8'>
        <span className='text-5xl font-bold text-slate-800'>Upcoming trips</span>
        <button type="button"
          className={classNames(
            "bg-indigo-400",
            "font-medium",
            "hover:bg-indigo-800",
            "px-5",
            "py-2.5",
            "mt-5",
            "sm:mt-0",
            "rounded-md",
            "text-center",
            "text-sm",
            "text-white",
          )}
          onClick={props.onCreateTripBtnClick}
        >
          + Create new trip
        </button>
      </div>
      {renderTripsTable()}
    </div>
  );
}

export default TripsContainer;

