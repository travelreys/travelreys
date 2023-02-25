import React, { FC } from 'react';
import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';

import ImagesAPI from '../../apis/images';
import { fmt, isEmptyDate, parseISO } from '../../lib/dates';

interface TripCardProps {
  trip: any
}

const TripCard: FC<TripCardProps> = (props: TripCardProps) => {

  const css = {
    dateTxt: "text-slate-500 text-sm md:text-sm align-base",
    tripCardCtn: "bg-white rounded-lg shadow-md h-fit",
    tripCardImg: "rounded-t-lg",
    tripCardName: "mb-2 text-xl font-bold tracking-tight text-slate-700 truncate",
    tripCardFooter: "flex justify-between items-center",
  }

  // Renderers
  const renderTripDates = () => {
    const dateRange = {
      from: parseISO(props.trip.startDate),
      to: parseISO(props.trip.endDate)
    };
    const dateFmt = "MMM d, yy"

    if (isEmptyDate(dateRange.from)) {
      return null;
    }
    return (
      <p className={css.dateTxt}>
        {fmt(dateRange.from, dateFmt)}
        &nbsp;-&nbsp;
        {fmt(dateRange.to, dateFmt)}
      </p>
    );
  }

  return (
    <Link
      to={`/trips/${props.trip.id}`}
      className={css.tripCardCtn}
    >
      <img
        srcSet={ImagesAPI.makeSrcSet(props.trip.coverImage)}
        src={ImagesAPI.makeSrc(props.trip.coverImage)}
        className={css.tripCardImg}
        alt="cover"
        referrerPolicy='no-referrer'
      />
      <div className="p-5">
        <h5 className={css.tripCardName}>
          {props.trip.name}
        </h5>
        <div className={css.tripCardFooter}>
          {renderTripDates()}
        </div>
      </div>
    </Link>
  );
}

// TripsContainer

interface TripsContainerProps {
  trips: Array<any>,
  onCreateBtnClick: any
}

const TripsContainer: FC<TripsContainerProps> = (props: TripsContainerProps) => {

  const { t } = useTranslation();

  const css = {
    headerCtn: "flex justify-between flex-col sm:flex-row items-center mb-8",
    header: "text-3xl sm:text-5x font-bold text-slate-800",
    createBtn: "bg-indigo-400 hover:bg-indigo-800 font-bold px-5 py-2.5 mt-5 sm:mt-0 rounded-md text-center text-sm text-white",
    tableCtn: "grid grid-cols-1 sm:grid-cols-3 xl:grid-cols-4 gap-4",
  }

  return (
    <div>
      <div className={css.headerCtn}>
        <span className={css.header}>
          {t('home.upcomingTrips')}
        </span>
        <button type="button"
          className={css.createBtn}
          onClick={props.onCreateBtnClick}
        >
          + Create new trip
        </button>
      </div>
      <div className={css.tableCtn}>
        {props.trips.map((trip: any) => {
          return <TripCard trip={trip} key={trip.id} />
        })}
      </div>
    </div>
  );
}

export default TripsContainer;
