import React, { FC } from 'react';
import classNames from 'classnames';
import Link from 'next/link';
import _get from 'lodash/get';
import { parseJSON, isEqual } from 'date-fns';

import ImagesAPI from '../../apis/images';
import Avatar from '../Avatar';

import { datesRenderer } from '../../utils/dates';
interface TripCardProps {
  trip: any
}

const TripCard: FC<TripCardProps> = (props: TripCardProps) => {

  // Renderers
  const renderCreatorAvatar = () => {
    return (
      <Avatar
        name={_get(props.trip, "creator.memberEmail")}
        placement="top"
      />
    );
  }

  const renderTripDates = () => {
    const nullDate = parseJSON("0001-01-01T00:00:00Z");
    const startDate = parseJSON(props.trip.startDate);
    const endDate = parseJSON(props.trip.endDate);

    if (isEqual(startDate, nullDate)) {
      return <div className='text-slate-500'>-</div>;
    }
    return (
      <p className='text-slate-500 text-sm md:text-sm align-base'>
        {datesRenderer(startDate, endDate)}
      </p>
    );
  }

  return (
    <Link
      href={`/trips/${props.trip.id}`}
      className="bg-white rounded-lg shadow-md h-fit"
    >
      <img
        srcSet={ImagesAPI.makeSrcSet(props.trip.coverImage)}
        src={ImagesAPI.makeSrc(props.trip.coverImage)}
        className="rounded-t-lg"
      />
      <div className="p-5">
        <h5 className="mb-2 text-xl font-bold tracking-tight text-slate-700 truncate">
          {props.trip.name}
        </h5>
        <div className='flex justify-between items-center'>
          {renderTripDates()}
          {renderCreatorAvatar()}
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
    return (
      <div className="grid grid-cols-1 sm:grid-cols-3 xl:grid-cols-4 gap-4 mx-3">
        {props.trips.map((trip: any) => {
          return <TripCard trip={trip} key={trip.id} />
        })}
      </div>
    );
  }

  return (
    <div>
      <div className='flex justify-between flex-col sm:flex-row items-center mb-8'>
        <span className='text-3xl sm:text-5x font-bold text-slate-800'>Upcoming trips</span>
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

