import React, { FC } from 'react';
import classNames from 'classnames';
import { Link } from 'react-router-dom';
import _get from 'lodash/get';
import { parseJSON, isEqual } from 'date-fns';
import { useTranslation } from 'react-i18next';

import ImagesAPI from '../../apis/images';
import Avatar from '../../components/common/Avatar';
import { printFromDateFromRange, printToDateFromRange } from '../../lib/dates';


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
    const dateRange = {
      from: parseJSON(props.trip.startDate),
      to: parseJSON(props.trip.endDate)
    };

    if (isEqual(dateRange.from, nullDate)) {
      return <div className='text-slate-500'>-</div>;
    }
    return (
      <p className='text-slate-500 text-sm md:text-sm align-base'>
        {printFromDateFromRange(dateRange, "MMM d, yy ")}
        &nbsp;-&nbsp;
        {printToDateFromRange(dateRange, "MMM d, yy ")}
      </p>
    );
  }

  return (
    <Link
      to={`/trips/${props.trip.id}`}
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

  const { t } = useTranslation();

  // Renderers
  const renderTripsTable = () => {
    return (
      <div className="grid grid-cols-1 sm:grid-cols-3 xl:grid-cols-4 gap-4">
        {props.trips.map((trip: any) => {
          return <TripCard trip={trip} key={trip.id} />
        })}
      </div>
    );
  }

  return (
    <>
      <div className='flex justify-between flex-col sm:flex-row items-center mb-8'>
        <span className='text-3xl sm:text-5x font-bold text-slate-800'>
          {t('title.upcomingTrips')}
        </span>
        <button type="button"
          className={classNames(
            "bg-indigo-400",
            "hover:bg-indigo-800",
            "font-bold",
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
    </>
  );
}

export default TripsContainer;
