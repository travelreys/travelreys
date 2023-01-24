import React, { FC, useState } from 'react';
import _get from "lodash/get";
import _sortBy from "lodash/sortBy";
import _isEmpty from "lodash/isEmpty";
import { formatDuration, intervalToDuration } from 'date-fns';

import {
  ArrowLongRightIcon,
  ChevronDownIcon,
  ChevronUpIcon,
  PlusCircleIcon,
} from '@heroicons/react/24/outline'
import { FlightsModalCss } from '../../styles/global';

import {
  parseTimeFromZ,
  printTime,
  prettyPrintMins,
} from '../../utils/dates';


interface TripFlightCardProps {
  itin: any
  onSelect: any
}

const TripFlightCard: FC<TripFlightCardProps> = (props: TripFlightCardProps) => {

  const [isExpanded, setIsExpanded] = useState(false);

  const numStops = _get(props.itin, "legs").length === 1
    ? "Non-stop" : `${_get(props.itin, "legs").length - 1} stops`;


  // Renderers
  const renderNumStops = () => {
    return (
      <span
        className='cursor-pointer border-b border-slate-400'
        onClick={() => { setIsExpanded(!isExpanded) }}
      >
        {numStops}&nbsp;
        {isExpanded
          ? <ChevronDownIcon className={FlightsModalCss.FlightStopIcon} />
          : <ChevronUpIcon className={FlightsModalCss.FlightStopIcon} />}
      </span>
    );
  }

  const renderStopsInfo = () => {
    return props.itin.legs.map((leg: any, idx: number) => {
      let layoverDuration = null;
      if (idx !== props.itin.legs.length - 1) {
        layoverDuration = intervalToDuration({
          start: parseTimeFromZ(props.itin.legs[idx + 1].departure.datetime),
          end: parseTimeFromZ(leg.arrival.datetime),
        });
      }

      return (
        <div key={idx}>
          <ol className={FlightsModalCss.FlightStopTimelineCtn}>
            <li className="mb-4 ml-6">
              <div className={FlightsModalCss.FlightStopTimelineIcon} />
              <h3 className={FlightsModalCss.FlightStopTimelineTime}>
                {printTime(parseTimeFromZ(leg.departure.datetime))}
              </h3>
              <p className={FlightsModalCss.FlightsStopTimelineText}>
                {leg.departure.airport.code} ({leg.departure.airport.name})
              </p>
              <p className={FlightsModalCss.FlightsStopTimelineText}>
                Travel time: {prettyPrintMins(leg.duration)}
              </p>
              <p className={FlightsModalCss.FlightsStopTimelineText}>
                {leg.operatingAirline.name} {leg.operatingAirline.code} {leg.flightNo}
              </p>
            </li>
            <li className="mb-4 ml-6">
              <div className={FlightsModalCss.FlightStopTimelineIcon} />
              <h3 className={FlightsModalCss.FlightStopTimelineTime}>
                {printTime(parseTimeFromZ(leg.arrival.datetime))}
              </h3>
              <p className={FlightsModalCss.FlightsStopTimelineText}>
                {leg.arrival.airport.code} ({leg.arrival.airport.name})
              </p>
              {layoverDuration ?
                  <p className={FlightsModalCss.FlightsStopLayoverText}>
                    {formatDuration(layoverDuration)} layover </p> : null
              }
            </li>
          </ol>
          <hr className={FlightsModalCss.FlightsStopHR} />
        </div>
      );
    });
  }

  const renderAirlineLogo = () => {
    const imgUrl = `https://www.gstatic.com/flights/airline_logos/70px/${props.itin.marketingAirline.code}.png`;
    const fallbackUrl = "https://cdn-icons-png.flaticon.com/512/4353/4353032.png";
    return (
      <div className="h-8 w-8 mr-4">
        <object className="h-8 w-8" data={imgUrl} type="image/png">
          <img className="h-8 w-8" src={fallbackUrl} alt={props.itin.marketingAirline.name} />
        </object>
      </div>
    );
  }

  const renderPricePill = () => {
    const pill = (
      <span className={FlightsModalCss.FlightPricePill}>
        {props.itin.priceWithCurrency.currency} {props.itin.priceWithCurrency.amount}
      </span>
    );
    return (
      <a href={props.itin.bookingURL} target="_blank">
        {pill}
      </a>
    );
  }

  return (
    <div className={FlightsModalCss.FlightTripCard}>
      {renderAirlineLogo()}
      <div className='flex-1'>
        <div className="flex">
          <span className=''>
            <p className='font-medium'>{printTime(parseTimeFromZ(props.itin.departure.datetime))}</p>
            <p className="text-xs text-slate-400">{props.itin.departure.airport.code}</p>
          </span>
          <ArrowLongRightIcon className='h-6 w-8' />
          <span className='mb-1'>
            <p className='font-medium'>{printTime(parseTimeFromZ(props.itin.arrival.datetime))}</p>
            <p className="text-xs text-slate-400">{props.itin.arrival.airport.code}</p>
          </span>
        </div>
        <span className="text-xs text-slate-400 block mb-1">
          {prettyPrintMins(props.itin.duration)}
        </span>
        <span className="text-xs text-slate-400 block mb-1">
          {props.itin.marketingAirline.name} | {renderNumStops()}
        </span>
        {isExpanded ? renderStopsInfo() : null}
      </div>
      <div className='flex flex-col text-right justify-between'>
        <div className='flex flex-row-reverse'>
          <PlusCircleIcon
            onClick={() => {props.onSelect(props.itin)}}
            className={FlightsModalCss.FlightPlusIcon}
          />
        </div>
        {renderPricePill()}
      </div>
    </div>
  );
}

export default TripFlightCard;
