import React, { FC, useState } from 'react';
import _get from "lodash/get";
import _sortBy from "lodash/sortBy";
import _isEmpty from "lodash/isEmpty";
import { formatDuration, intervalToDuration } from 'date-fns';

import {
  ArrowLongRightIcon,
  ChevronDownIcon,
  ChevronUpIcon,
  EllipsisHorizontalCircleIcon,
} from '@heroicons/react/24/outline'
import { TrashIcon } from '@heroicons/react/24/solid';

import Dropdown from '../../components/common/Dropdown';
import { CommonCss, TripLogisticsCss } from '../../assets/styles/global';

import {
  flightArrivalTime,
  flightDepartureTime,
  FlightDirectionDepart,
  FlightItineraryTypeRoundtrip,
  flightLegOpAirline,
  flightLegs,
  flightLogoUrl,
  Flights,
  logoFallbackImg
} from '../../lib/flights';
import { flightItineraryType } from '../../lib/trips';
import {
  parseISO,
  printFmt,
  prettyPrintMins,
} from '../../lib/dates';
import {capitaliseWords} from '../../lib/strings';



interface TransitFlightCardProps {
  flight: any
  onDelete: any
}

const TransitFlightCard: FC<TransitFlightCardProps> = (props: TransitFlightCardProps) => {
  const [isExpanded, setIsExpanded] = useState(false);


  // Renderers
  const renderNumStops = () => {
    const departFlight = _get(props.flight, FlightDirectionDepart, {});
    const numStops =  flightLegs(departFlight).length === 1
    ? "Non-stop" : `${flightLegs(departFlight).length - 1} stops`;

    return (
      <span
        className={TripLogisticsCss.FlightTransitNumStop}
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
    const timeFmt = "hh:mm aa";

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
                {printFmt(parseISO(leg.departure.datetime), timeFmt)}
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
                {printFmt(parseISO(leg.arrival.datetime), timeFmt)}
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
    const airline: any = flightLegOpAirline(_get(flight.legs, "0", {}));
    return (
      <div className={TripLogisticsCss.FlightTransitLogoImgWrapper}>
        <object
          className={TripLogisticsCss.FlightTransitLogoImg}
          data={flightLogoUrl(airline.code)}
          type="image/png"
        >
          <img
            className={TripLogisticsCss.FlightTransitLogoImg}
            src={logoFallbackImg}
            alt={airline.name}
          />
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

  const renderSettingsDropdown = () => {
    const opts = [
      <button
        type='button'
        className={CommonCss.DeleteBtn}
        onClick={() => props.onDelete(props.flight)}
      >
        <TrashIcon className={CommonCss.LeftIcon} />
        Delete
      </button>
    ];
    const menu = (
      <EllipsisHorizontalCircleIcon
        className={CommonCss.DropdownIcon} />
    );
    return (
      <div className="flex flex-row-reverse">
        <Dropdown menu={menu} opts={opts} />
      </div>
    )
  }

  const renderFlight = (flight: Flights.Flight, direction: string) => {
    const airline: any = flightLegOpAirline(flight.legs[0]);
    const timeFmt = "hh:mm aa";
    const dateFmt = "eee, MMM d";

    const departTime = flightDepartureTime(flight) as string;
    const arrTime = flightArrivalTime(flight) as string;

    return (
      <div className={TripLogisticsCss.FlightTransit}>
        {renderAirlineLogo(flight)}
        <div className='flex-1'>
          <p className={TripLogisticsCss.FlightTransitDatetime}>
            {capitaliseWords(direction)} Flight&nbsp;&#x2022;&nbsp;
            {printFmt(parseISO(departTime), dateFmt)}
          </p>
          <div className="flex">
            <span>
              <p className={TripLogisticsCss.FlightTransitTime}>
                {printFmt(parseISO(departTime), timeFmt)}
              </p>
              <p className={TripLogisticsCss.FlightTransitAirportCode}>
                {flight.departure.airport.code}
              </p>
            </span>
            <ArrowLongRightIcon
              className={TripLogisticsCss.FlightTransitLongArrow}
            />
            <span className='mb-1'>
              <p className={TripLogisticsCss.FlightTransitTime}>
                {printFmt(parseISO(arrTime), timeFmt)}
              </p>
              <p className={TripLogisticsCss.FlightTransitAirportCode}>
                {flight.arrival.airport.code}
              </p>
            </span>
          </div>
          <span className={TripLogisticsCss.FlightTransitDuration}>
            {prettyPrintMins(flight.duration)}
          </span>
          <span className={TripLogisticsCss.FlightTransitDuration}>
            {airline.name} | {renderNumStops()}
          </span>
          {isExpanded ? renderStopsInfo(flight) : null}
        </div>
        { direction !== "departing" ? null :
          <div className='flex flex-col justify-between'>
            {renderSettingsDropdown()}
            {renderPricePill()}
          </div>
        }
      </div>
    );
  }


  const renderDepartFlight = () => {
    return renderFlight(_get(props.flight, "depart", {}), "departing");
  }

  const renderReturnFlight = () => {
    if (flightItineraryType(props.flight) !== FlightItineraryTypeRoundtrip) {
      return null;
    }
    return renderFlight(_get(props.flight, "return", {}), "returning");
  }

  return (
    <div className={TripLogisticsCss.FlightTransitCard}>
      {renderDepartFlight()}
      {renderReturnFlight()}
    </div>
  );
}

export default TransitFlightCard;
