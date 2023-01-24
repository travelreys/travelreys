import React, { FC, useState } from 'react';
import _get from "lodash/get";
import _sortBy from "lodash/sortBy";
import _isEmpty from "lodash/isEmpty";
import { formatDuration, intervalToDuration } from 'date-fns';

import {
  ArrowLongRightIcon,
  ChevronDownIcon,
  ChevronUpIcon,
  XCircleIcon,
} from '@heroicons/react/24/outline'
import { TripLogisticsCss } from '../../styles/global';

import {
  parseTimeFromZ,
  printTime,
  prettyPrintMins,
} from '../../utils/dates';


interface TransitFlightCardProps {
  flight: any
  onDelete: any
}

const TransitFlightCard: FC<TransitFlightCardProps> = (props: TransitFlightCardProps) => {

  const [isExpanded, setIsExpanded] = useState(false);

  const numStops = _get(props.flight, "legs").length === 1
    ? "Non-stop" : `${_get(props.flight, "legs").length - 1} stops`;


  // Renderers
  const renderNumStops = () => {
    return (
      <span
        className='cursor-pointer border-b border-slate-400'
        onClick={() => { setIsExpanded(!isExpanded) }}
      >
        {numStops}&nbsp;
        {isExpanded
          ? <ChevronDownIcon className={TripLogisticsCss.FlightStopIcon} />
          : <ChevronUpIcon className={TripLogisticsCss.FlightStopIcon} />}
      </span>
    );
  }

  const renderStopsInfo = () => {
    return props.flight.legs.map((leg: any, idx: number) => {
      let layoverDuration = null;
      if (idx !== props.flight.legs.length - 1) {
        layoverDuration = intervalToDuration({
          start: parseTimeFromZ(props.flight.legs[idx + 1].departure.datetime),
          end: parseTimeFromZ(leg.arrival.datetime),
        });
      }

      return (
        <div key={idx}>
          <ol className={TripLogisticsCss.FlightStopTimelineCtn}>
            <li className="mb-4 ml-6">
              <div className={TripLogisticsCss.FlightStopTimelineIcon} />
              <h3 className={TripLogisticsCss.FlightStopTimelineTime}>
                {printTime(parseTimeFromZ(leg.departure.datetime))}
              </h3>
              <p className={TripLogisticsCss.FlightsStopTimelineText}>
                {leg.departure.airport.code} ({leg.departure.airport.name})
              </p>
              <p className={TripLogisticsCss.FlightsStopTimelineText}>
                Travel time: {prettyPrintMins(leg.duration)}
              </p>
              <p className={TripLogisticsCss.FlightsStopTimelineText}>
                {leg.operatingAirline.name} {leg.operatingAirline.code} {leg.flightNo}
              </p>
            </li>
            <li className="mb-4 ml-6">
              <div className={TripLogisticsCss.FlightStopTimelineIcon} />
              <h3 className={TripLogisticsCss.FlightStopTimelineTime}>
                {printTime(parseTimeFromZ(leg.arrival.datetime))}
              </h3>
              <p className={TripLogisticsCss.FlightsStopTimelineText}>
                {leg.arrival.airport.code} ({leg.arrival.airport.name})
              </p>
              {layoverDuration ?
                  <p className={TripLogisticsCss.FlightsStopLayoverText}>
                    {formatDuration(layoverDuration)} layover </p> : null
              }
            </li>
          </ol>
          <hr className={TripLogisticsCss.FlightsStopHR} />
        </div>
      );
    });
  }

  const renderAirlineLogo = () => {
    const imgUrl = `https://www.gstatic.com/flights/airline_logos/70px/${props.flight.marketingAirline.code}.png`;
    const fallbackUrl = "https://cdn-icons-png.flaticon.com/512/4353/4353032.png";
    return (
      <div className="h-8 w-8 mr-4">
        <object className="h-8 w-8" data={imgUrl} type="image/png">
          <img className="h-8 w-8" src={fallbackUrl} alt={props.flight.marketingAirline.name} />
        </object>
      </div>
    );
  }

  const renderPricePill = () => {
    return (
      <span className={TripLogisticsCss.FlightPricePill}>
        {props.flight.priceWithCurrency.currency} {props.flight.priceWithCurrency.amount}
      </span>
    );
  }

  return (
    <div className={TripLogisticsCss.FlightTripCard}>
      {renderAirlineLogo()}
      <div className='flex-1'>
        <div className="flex">
          <span className=''>
            <p className='font-medium'>{printTime(parseTimeFromZ(props.flight.departure.datetime))}</p>
            <p className="text-xs text-slate-800">{props.flight.departure.airport.code}</p>
          </span>
          <ArrowLongRightIcon className='h-6 w-8' />
          <span className='mb-1'>
            <p className='font-medium'>{printTime(parseTimeFromZ(props.flight.arrival.datetime))}</p>
            <p className="text-xs text-slate-800">{props.flight.arrival.airport.code}</p>
          </span>
        </div>
        <span className="text-xs text-slate-800 block mb-1">
          {prettyPrintMins(props.flight.duration)}
        </span>
        <span className="text-xs text-slate-800 block mb-1">
          {props.flight.marketingAirline.name} | {renderNumStops()}
        </span>
        {isExpanded ? renderStopsInfo() : null}
      </div>
      <div className='flex flex-col text-right justify-between'>
        <div className='flex flex-row-reverse'>
          <XCircleIcon
            onClick={() => {props.onDelete(props.flight)}}
            className={TripLogisticsCss.FlightTransitIcon}
          />
        </div>
        {renderPricePill()}
      </div>
    </div>
  );
}

export default TransitFlightCard;
