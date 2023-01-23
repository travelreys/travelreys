import React, { FC, useState } from 'react';
import _get from "lodash/get";
import _sortBy from "lodash/sortBy";
import _isEmpty from "lodash/isEmpty";
import { useMediaQuery } from 'usehooks-ts';
import {
  formatDuration,
  format,
  intervalToDuration,
  parseISO,
} from 'date-fns';
import {
  DayPicker,
  DateRange,
  SelectRangeEventHandler
} from 'react-day-picker';

import {
  ArrowLongRightIcon,
  CalendarDaysIcon,
  ChevronDownIcon,
  ChevronUpIcon,
  MapPinIcon,
  PaperAirplaneIcon,
  PlusCircleIcon,
  XMarkIcon,
} from '@heroicons/react/24/outline'
import { ModalCss, FlightsModalCss } from '../../styles/global';

import FlightsAPI, { flights } from '../../apis/flights';
import Spinner from '../../components/Spinner';
import { parseTimeFromZ, printTime, prettyPrintMins } from '../../utils/dates';

// TripFlightCard

interface TripFlightCardProps {
  itin: any
}

const TripFlightCard: FC<TripFlightCardProps> = (props: TripFlightCardProps) => {

  const [isExpanded, setIsExpanded] = useState(false);

  const numStops = _get(props.itin, "legs").length === 1
    ? "Non-stop" : `${_get(props.itin, "legs").length} stops`;


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
        <>
          <ol className="relative border-l border-dashed border-slate-300 mt-4 ml-2">
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
        </>
      );
    });
  }

  return (
    <div className="flex p-4 rounded shadow-md hover:shadow-lg hover:shadow-indigo-100">
      <div className="h-8 w-8 mr-4">
        <img
          alt={props.itin.marketingAirline.name}
          className="h-8 w-8"
          src={`https://www.gstatic.com/flights/airline_logos/70px/${props.itin.marketingAirline.code}.png`}
        />
      </div>
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
          <PlusCircleIcon className={FlightsModalCss.FlightPlusIcon} />
        </div>
        <a href={props.itin.bookingURL} target="_blank">
          <span className={FlightsModalCss.FlightBookLink}>SGD {props.itin.price}</span>
        </a>
      </div>
    </div>
  );
}


// TripFlightsSearchForm

interface TripFlightsSearchFormProps {
  onSearch: any
}

const TripFlightsSearchForm: FC<TripFlightsSearchFormProps> = (props: TripFlightsSearchFormProps) => {

  const [itineraryClass, setItineraryClass] = useState(itineraryClasses[0]);
  const [isCabinClassDropdownActive, setIsCabinClassDropdownActive] = useState(false);
  const [flightDates, setFlightDates] = useState<DateRange>();
  const [isCalendarOpen, setIsCalendarOpen] = useState(false);

  const matches = useMediaQuery('(min-width: 768px)');

  // Event Handlers
  const searchBtnOnClick = () => {
    props.onSearch();
  }

  const dateInputOnClick = (event: React.MouseEvent<HTMLInputElement>) => {
    setIsCalendarOpen(!isCalendarOpen)
  }

  const flightDatesOnSelect: SelectRangeEventHandler = (range: DateRange | undefined) => {
    setFlightDates(range);
  };

  // Renderers
  const renderCabinClassDropdown = () => {
    const cabinClassOpts = (
      <div className="z-10 w-44 rounded-lg bg-white shadow block absolute">
        <ul className="z-10 w-44 rounded-lg bg-white shadow">
          {itineraryClasses.map((type: any) => (
            <li
              key={type.value}
              className="block rounded-lg py-2 px-4 cursor-pointer hover:bg-indigo-100"
              onClick={() => setItineraryClass(type)}
            >
              {type.label}
            </li>
          ))
          }
        </ul>
      </div>
    );

    return (
      <div className='relative inline-block'>
        <button
          onClick={() => { setIsCabinClassDropdownActive(!isCabinClassDropdownActive) }}
          onBlur={() => {
            setTimeout(() => {
              setIsCabinClassDropdownActive(false);
            }, 200)
          }}
          className={FlightsModalCss.ItineraryDropdownBtn}
        >
          {itineraryClass.label}&nbsp;
          {isCabinClassDropdownActive ?
            <ChevronUpIcon className={FlightsModalCss.ItineraryDropdownIcon} />
            : <ChevronDownIcon className={FlightsModalCss.ItineraryDropdownIcon} />}
        </button>
        {isCabinClassDropdownActive ? cabinClassOpts : null}
      </div>
    );
  }

  const renderSearchFilters = () => {
    return (
      <div className="flex pb-2 justify-between">
        <div className='flex-1'>
          {renderCabinClassDropdown()}
        </div>
        <button
          className={FlightsModalCss.FlightSearchBtn}
          onClick={searchBtnOnClick}
        >
          Search
        </button>
      </div>
    );
  }

  const renderAirportsInput = () => {
    return (
      <div className='flex justify-between mb-4'>
        <div className="relative w-6/12">
          <div className={FlightsModalCss.FlightFromIconCtn}>
            <MapPinIcon className={FlightsModalCss.FlightFromIcon} />
          </div>
          <input
            type="text"
            className={FlightsModalCss.FlightFromInput}
            placeholder="from city, airport"
          />
        </div>
        <div className="relative w-6/12">
          <div className={FlightsModalCss.FlightFromIconCtn}>
            <PaperAirplaneIcon className={FlightsModalCss.FlightFromIcon} />
          </div>
          <input
            type="text"
            className={FlightsModalCss.FlightFromInput}
            placeholder="to city, airport"
          />
        </div>
      </div>
    );
  }

  const renderDatesInputs = () => {
    const startInputValue = _get(flightDates, "from") ? format(_get(flightDates, "from", new Date()), 'y-MM-dd') : "";
    const endInputValue = _get(flightDates, "to") ? format(_get(flightDates, "to", new Date()), 'y-MM-dd') : "";
    let value = startInputValue;
    if (!_isEmpty(endInputValue)) {
      value = `${startInputValue} - ${endInputValue}`
    }

    return (
      <div className={FlightsModalCss.FlightDatesCtn}>
        <span className={FlightsModalCss.FlightDatesLabel}>
          <CalendarDaysIcon className={FlightsModalCss.FlightDatesIcon} />
          &nbsp;Dates
        </span>
        <input
          type="text"
          value={value}
          onChange={() => { }}
          onClick={dateInputOnClick}
          className={FlightsModalCss.FlightDatesInput}
        />
      </div>
    );
  }

  const renderDayPicker = () => {
    if (!isCalendarOpen) {
      return;
    }
    return (
      <div className='relative'>
        <div className='absolute bg-white mt-2 border border-slate-200'>
          <DayPicker
            mode="range"
            numberOfMonths={matches ? 2 : 1}
            pagedNavigation
            styles={{ months: { margin: "0", display: "flex", justifyContent: "space-around" } }}
            modifiersStyles={{
              selected: { background: "#AC8AC3" }
            }}
            selected={flightDates}
            onSelect={flightDatesOnSelect}
          />
        </div>
      </div>
    );
  }

  return (
    <div>
      {renderSearchFilters()}
      {renderAirportsInput()}
      <div className="mb-4">
        {renderDatesInputs()}
        {renderDayPicker()}
      </div>
    </div>
  );
}


// TripFlightsModal

interface TripFlightsModalProps {
  isOpen: boolean
  onClose: any
}

const itineraryClasses = [
  { label: "Economy", value: "ECO" },
  { label: "Premium Economy", value: "PEC" },
  { label: "Business", value: "BUS" },
  { label: "First Class", value: "FST" },
];


const TripFlightsModal: FC<TripFlightsModalProps> = (props: TripFlightsModalProps) => {

  const [itineraries, setItineraries] = useState([] as any);

  const [isLoading, setIsLoading] = useState(false);
  const [searchInitiated, setSearchInitiated] = useState(false);


  // Event Handlers
  const onSearch = (origIATA: string, destIATA: string, departDate: string) => {
    setSearchInitiated(true);
    setIsLoading(true);
    FlightsAPI.search("")
    .then(res => {
      const itineraries = _sortBy(_get(res, "data.itineraries", []), "price")
      setItineraries(itineraries);
    })
    .finally(() => {
      setIsLoading(false);
    });
    // const itineraries = _sortBy(_get(flights, "itineraries", []), "price")
    // setItineraries(itineraries);
  }

  // Renderers
  const renderItineraries = () => {
    if (!searchInitiated) {
      return (<></>);
    }

    if (isLoading) {
      return <Spinner />
    }

    return (
      <div>
        <p className={FlightsModalCss.FlightSearchResultsTitle}>Flights</p>
        {itineraries.map((itin: any) => <TripFlightCard itin={itin} />)}
      </div>
    );
  }


  if (!props.isOpen) {
    return (<></>);
  }

  return (
    <div className={ModalCss.Container}>
      <div className={ModalCss.Inset}></div>
      <div className={ModalCss.Content}>
        <div className={ModalCss.ContentContainer}>
          <div className={ModalCss.ContentCard}>
            <div className="px-4 pt-5 sm:p-8 sm:pb-2 rounded-t-lg">
              <div className='flex justify-between mb-6'>
                <h2 className="text-lg sm:text-2xl font-bold text-center text-slate-900">
                  Search flights
                </h2>
                <button type="button" onClick={props.onClose}>
                  <XMarkIcon className='h-6 w-6 text-slate-700' />
                </button>
              </div>
              <TripFlightsSearchForm onSearch={onSearch} />
              {renderItineraries()}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export default TripFlightsModal;
