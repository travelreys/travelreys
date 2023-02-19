import React, { FC } from 'react';
import { Link } from 'react-router-dom';
import _get from 'lodash/get';
import { useTranslation } from 'react-i18next';

import ImagesAPI from '../../apis/images';
import {
  isEmptyDate,
  parseISO,
  printFromDateFromRange,
  printToDateFromRange,
} from '../../lib/dates';
import { TripContainerCss } from '../../assets/styles/global';


interface TripCardProps {
  trip: any
}

const TripCard: FC<TripCardProps> = (props: TripCardProps) => {

  // Renderers
  const renderTripDates = () => {
    const dateRange = {
      from: parseISO(props.trip.startDate),
      to: parseISO(props.trip.endDate)
    };
    const dateFmt = "MMM d, yy"

    if (isEmptyDate(dateRange.from)) {
      return <div className='text-slate-500'>-</div>;
    }
    return (
      <p className={TripContainerCss.DateTxt}>
        {printFromDateFromRange(dateRange, dateFmt)}
        &nbsp;-&nbsp;
        {printToDateFromRange(dateRange, dateFmt)}
      </p>
    );
  }

  return (
    <Link
      to={`/trips/${props.trip.id}`}
      className={TripContainerCss.TripCardCtn}
    >
      <img
        srcSet={ImagesAPI.makeSrcSet(props.trip.coverImage)}
        src={ImagesAPI.makeSrc(props.trip.coverImage)}
        className={TripContainerCss.TripCardImg}
      />
      <div className="p-5">
        <h5 className={TripContainerCss.TripCardName}>
          {props.trip.name}
        </h5>
        <div className={TripContainerCss.TripCardFooter}>
          {renderTripDates()}
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
  return (
    <>
      <div className={TripContainerCss.HeaderCtn}>
        <span className={TripContainerCss.Header}>
          {t('home.upcomingTrips')}
        </span>
        <button type="button"
          className={TripContainerCss.CreateBtn}
          onClick={props.onCreateTripBtnClick}
        >
          + Create new trip
        </button>
      </div>
      <div className={TripContainerCss.TableCtn}>
        {props.trips.map((trip: any) => {
          return <TripCard trip={trip} key={trip.id} />
        })}
      </div>
    </>
  );
}

export default TripsContainer;
