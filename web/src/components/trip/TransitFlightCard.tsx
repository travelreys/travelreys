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
  parseISO,
  printFmt,
  prettyPrintMins,
} from '../../utils/dates';
import {capitaliseWords} from '../../utils/strings';


interface TransitFlightCardProps {
  flight: any
  onDelete: any
}

const TransitFlightCard: FC<TransitFlightCardProps> = (props: TransitFlightCardProps) => {
  const [isExpanded, setIsExpanded] = useState(false);

  const departFlight = _get(props.flight, "depart", {});
  const returnFlight = _get(props.flight, "return");
  const numStops = _get(departFlight, "legs", []).length === 1
    ? "Non-stop" : `${_get(departFlight, "legs", []).length - 1} stops`;

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

  const renderStopsInfo = (flight: any) => {
    return flight.legs.map((leg: any, idx: number) => {
      let layoverDuration = null;
      if (idx !== flight.legs.length - 1) {
        layoverDuration = intervalToDuration({
          start: parseISO(flight.legs[idx + 1].departure.datetime),
          end: parseISO(leg.arrival.datetime),
        });
      }


      return (
        <div key={idx}>
          <ol className={TripLogisticsCss.FlightStopTimelineCtn}>
            <li className="mb-4 ml-6">
              <div className={TripLogisticsCss.FlightStopTimelineIcon} />
              <h3 className={TripLogisticsCss.FlightStopTimelineTime}>
                {printFmt(parseISO(leg.departure.datetime), "hh:mm aa")}
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
                {printFmt(parseISO(leg.arrival.datetime), "hh:mm aa")}
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

  const renderAirlineLogo = (flight: any) => {
    const airline = _get(flight.legs, "0.operatingAirline", {});
    const imgUrl = `https://www.gstatic.com/flights/airline_logos/70px/${airline.code}.png`;
    const fallbackUrl = "https://cdn-icons-png.flaticon.com/512/4353/4353032.png";
    return (
      <div className="h-8 w-8 mr-4">
        <object className="h-8 w-8" data={imgUrl} type="image/png">
          <img className="h-8 w-8" src={fallbackUrl} alt={airline.name} />
        </object>
      </div>
    );
  }

  const renderPricePill = () => {
    return (
      <span className={TripLogisticsCss.FlightPricePill}>
        {props.flight.priceMetadata.currency} {props.flight.priceMetadata.amount}
      </span>
    );
  }

  const renderFlight = (flight: any, direction: string) => {
    const airline = _get(flight.legs, "0.operatingAirline", {});
    return (
      <div className={TripLogisticsCss.FlightTransit}>
        {renderAirlineLogo(flight)}
        <div className='flex-1'>
          <p className='text-sm text-slate-800'>
            {capitaliseWords(direction)} Flight&nbsp;&#x2022;&nbsp;
            {printFmt(parseISO(flight.departure.datetime), "eee, MMM d")}
          </p>
          <div className="flex">
            <span className=''>
              <p className='font-medium'>
                {printFmt(parseISO(flight.departure.datetime), "hh:mm aa")}
              </p>
              <p className="text-xs text-slate-800">{flight.departure.airport.code}</p>
            </span>
            <ArrowLongRightIcon className='h-6 w-8' />
            <span className='mb-1'>
              <p className='font-medium'>
                {printFmt(parseISO(flight.arrival.datetime), "hh:mm aa")}
              </p>
              <p className="text-xs text-slate-800">{flight.arrival.airport.code}</p>
            </span>
          </div>
          <span className="text-xs text-slate-800 block mb-1">
            {prettyPrintMins(flight.duration)}
          </span>
          <span className="text-xs text-slate-800 block mb-1">
            {airline.name} | {renderNumStops()}
          </span>
          {isExpanded ? renderStopsInfo(flight) : null}
        </div>
        { direction !== "departing" ? null :
          <div className='flex flex-col text-right justify-between'>
            <div className='flex flex-row-reverse'>
              <XCircleIcon
                onClick={() => {props.onDelete(props.flight)}}
                className={TripLogisticsCss.FlightTransitIcon}
              />
            </div>
            {renderPricePill()}
          </div>
        }
      </div>
    );
  }

  return (
    <div className={TripLogisticsCss.FlightTransitCard}>
      {renderFlight(_get(props.flight, "depart", {}), "departing")}
      { props.flight.itineraryType === "roundtrip"
        ? renderFlight(_get(props.flight, "return", {}), "returning") : null }
    </div>
  );
}

export default TransitFlightCard;
