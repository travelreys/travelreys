import React, { FC, useState } from 'react';
import _get from "lodash/get";
import _sortBy from "lodash/sortBy";
import _isEmpty from "lodash/isEmpty";
import { format, formatDuration, intervalToDuration } from 'date-fns';
import { DateRange, SelectRangeEventHandler } from 'react-day-picker';

import {
  ArrowLongRightIcon,
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
import {
  parseTimeFromZ,
  printTime,
  prettyPrintMins,
  printFromDateFromRange,
  printToDateFromRange
} from '../../utils/dates';
import { capitaliseWords } from '../../utils/strings';
import InputDatesPicker from '../InputDatesPicker';
import Alert from '../Alert';

// TripFlightCard

interface TripFlightCardProps {
  itin: any
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
        </div>
      );
    });
  }

  return (
    <div className="flex p-4 rounded shadow-md hover:shadow-lg hover:shadow-indigo-100">
      <div className="h-8 w-8 mr-4">
        <object
          className="h-8 w-8"
          data={`https://www.gstatic.com/flights/airline_logos/70px/${props.itin.marketingAirline.code}.png`}
          type="image/png"
        >
          <img
            className="h-8 w-8"
            src="https://cdn-icons-png.flaticon.com/512/4353/4353032.png"
            alt={props.itin.marketingAirline.name}
          />
        </object>
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

const cabinClasses = [
  { label: "Economy", value: "economy" },
  { label: "Premium Economy", value: "premiumeconomy" },
  { label: "Business", value: "business" },
  { label: "First Class", value: "first" },
];

interface TripFlightsSearchFormProps {
  onSearch: any
}

const TripFlightsSearchForm: FC<TripFlightsSearchFormProps> = (props: TripFlightsSearchFormProps) => {

  // Data State
  const [origIATA, setOrigIATA] = useState("");
  const [destIATA, setDestIATA] = useState("");
  const [cabinClass, setCabinClass] = useState(cabinClasses[0]);
  const [flightDates, setFlightDates] = useState<DateRange>();

  // UI State
  const [isCabinClassActive, setIsCabinClassActive] = useState(false);
  const [isOrigIATAFocus, setIsOrigIATAFocus] = useState(false);
  const [isDestIATAFocus, setIsDestIATAFocus] = useState(false);
  const [origIATAQuery, setOrigIATAQuery] = useState("");
  const [destIATAQuery, setDestIATAQuery] = useState("");

  // Event Handlers
  const searchBtnOnClick = () => {
    const departDate = printFromDateFromRange(flightDates, 'y-MM-dd');
    const arrDate = printToDateFromRange(flightDates, 'y-MM-dd');
    props.onSearch(origIATA, destIATA, departDate, arrDate, cabinClass);
  }

  const flightDatesOnSelect: SelectRangeEventHandler = (range: DateRange | undefined) => {
    setFlightDates(range);
  };

  // Renderers
  const renderCabinClassDropdown = () => {
    const cabinClassOpts = (
      <div className={FlightsModalCss.CabinClassOptCtn}>
        <ul className={FlightsModalCss.CabinClassOptList}>
          {cabinClasses.map((type: any) => (
            <li
              key={type.value}
              className={FlightsModalCss.CabinClassOpt}
              onClick={() => setCabinClass(type)}
            >
              {type.label}
            </li>
          ))}
        </ul>
      </div>
    );

    return (
      <div className='relative inline-block'>
        <button
          onClick={() => { setIsCabinClassActive(!isCabinClassActive) }}
          onBlur={() => {
            setTimeout(() => {
              setIsCabinClassActive(false);
            }, 200)
          }}
          className={FlightsModalCss.CabinClassDropdownBtn}
        >
          {cabinClass.label}&nbsp;
          {isCabinClassActive ?
            <ChevronUpIcon className={FlightsModalCss.CabinClassDropdownIcon} />
            : <ChevronDownIcon className={FlightsModalCss.CabinClassDropdownIcon} />}
        </button>
        {isCabinClassActive ? cabinClassOpts : null}
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

  const renderAirportsAutocomplete = (isFocus: any, query: string, setIATA: any, setQuery: any, setFocus: any) => {
    if (!isFocus || _isEmpty(query)) {
      return (<></>);
    }
    return (
      <div className={FlightsModalCss.AirportSearchOptCts}>
        <ul className={FlightsModalCss.AirportSearchOptList}>
          {FlightsAPI.airportAutocomplete(query).map((ap: any) => (
              <li
                key={ap.iata}
                className={FlightsModalCss.AirportSearchOpt}
                onClick={() => {
                  setIATA(ap.iata);
                  setQuery(capitaliseWords(ap.airport));
                  setFocus(false)}
                }
              >
                {capitaliseWords(ap.airport)} ({ap.iata.toUpperCase()})
              </li>
            ))}
        </ul>
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
            value={origIATAQuery}
            onChange={(e) => {setOrigIATAQuery(e.target.value)}}
            onFocus={() => {setIsOrigIATAFocus(true)}}
            onBlur={(e) => {setTimeout(() => { setIsOrigIATAFocus(false) }, 200)}}
          />
          {renderAirportsAutocomplete(isOrigIATAFocus, origIATAQuery, setOrigIATA, setOrigIATAQuery, setIsOrigIATAFocus)}
        </div>
        <div className="relative w-6/12">
          <div className={FlightsModalCss.FlightFromIconCtn}>
            <PaperAirplaneIcon className={FlightsModalCss.FlightFromIcon} />
          </div>
          <input
            type="text"
            className={FlightsModalCss.FlightFromInput}
            placeholder="to city, airport"
            value={destIATAQuery}
            onChange={(e) => { setDestIATAQuery(e.target.value)}}
            onFocus={() => {setIsDestIATAFocus(true)}}
            onBlur={(e) => {setTimeout(() => { setIsDestIATAFocus(false) }, 200)}}
          />
          {renderAirportsAutocomplete(isDestIATAFocus, destIATAQuery, setDestIATA, setDestIATAQuery, setIsDestIATAFocus)}
        </div>
      </div>
    );
  }

  return (
    <div>
      {renderSearchFilters()}
      {renderAirportsInput()}
      <InputDatesPicker
        onSelect={flightDatesOnSelect}
        dates={flightDates}
      />
    </div>
  );
}


// TripFlightsModal

interface TripFlightsModalProps {
  isOpen: boolean
  onClose: any
}

const TripFlightsModal: FC<TripFlightsModalProps> = (props: TripFlightsModalProps) => {

  const [itineraries, setItineraries] = useState([] as any);
  const [isLoading, setIsLoading] = useState(false);
  const [searchInitiated, setSearchInitiated] = useState(false);
  const [alertMsg, setAlertMsg] = useState("");


  // Event Handlers
  const onSearch = (origIATA: string, destIATA: string, departDate: string, returnDate: string | undefined, cabinClass: string) => {
    if (_isEmpty(origIATA)) {
      setAlertMsg("Please select a flight origin");
      return;
    }
    if (_isEmpty(destIATA)) {
      setAlertMsg("Please select a flight destination");
      return;
    }
    if (_isEmpty(departDate)) {
      setAlertMsg("Please select a departure date");
      return;
    }

    setSearchInitiated(true);
    setIsLoading(true);
    setAlertMsg("");

    FlightsAPI.search(origIATA, destIATA, departDate, returnDate, cabinClass)
    .then(res => {
      const itineraries = _sortBy(_get(res, "data.itineraries", []), "score")
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
        {itineraries.map((itin: any, idx: number) => <TripFlightCard key={idx} itin={itin} />)}
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
              {!_isEmpty(alertMsg) ? <Alert title={""} message={alertMsg} /> : null}
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
