import React, { FC, useEffect, useState } from 'react';
import _get from "lodash/get";
import _sortBy from "lodash/sortBy";
import _minBy from "lodash/minBy";

import _isEmpty from "lodash/isEmpty";
import { DateRange, SelectRangeEventHandler } from 'react-day-picker';
import { v4 as uuidv4 } from 'uuid';
import { formatDuration, intervalToDuration } from 'date-fns';

import {
  ArrowLongRightIcon,
  ChevronDownIcon,
  ChevronUpIcon,
  MapPinIcon,
  PaperAirplaneIcon,
  PlusCircleIcon,
  XMarkIcon,
} from '@heroicons/react/24/outline'

import Alert from '../../components/common/Alert';
import Modal from '../../components/common/Modal';
import InputDatesPicker from '../../components/common/InputDatesPicker';
import Spinner from '../../components/common/Spinner';

import FlightsAPI from '../../apis/flights';
import { Trips } from '../../lib/trips';
import {
  parseISO,
  printFmt,
  prettyPrintMins,
  printFromDateFromRange,
  printToDateFromRange,
  parseTripDate
} from '../../lib/dates';
import { capitaliseWords } from '../../lib/strings';
import { CommonCss, FlightsModalCss } from '../../assets/styles/global';

////////////////////////////
// OnewayFlightsContainer //
////////////////////////////

interface OnewayFlightsContainerProps {
  oneways: any
  onSelect: any
}

const OnewayFlightsContainer: FC<OnewayFlightsContainerProps> = (props: OnewayFlightsContainerProps) => {

  return (
    <div>
      <p className={FlightsModalCss.FlightSearchResultsTitle}>
        One-way Flights
      </p>
      {props.oneways.map((flight: any, idx: number) =>
        <FlightCard
          key={idx}
          flight={flight.depart}
          bookingMetadata={flight.bookingMetadata}
          onSelect={props.onSelect}
        />
      )}
    </div>
  );
}

///////////////////////////////
// RoundtripFlightsContainer //
///////////////////////////////

interface RoundtripFlightsContainerProps {
  roundtrips: any
  onSelect: any
}

const RoundtripFlightsContainer: FC<RoundtripFlightsContainerProps> = (props: RoundtripFlightsContainerProps) => {

  const [stepperStep, setStepperStep] = useState(0);
  const [roundtrip, setRoundtrip] = useState(null as any);

  let roundtrips = Object.values(props.roundtrips).map((rt: any) => {
    const minScore = _minBy(rt.bookingMetadata, (bm: any) => bm.score)
    return Object.assign(rt, {score: minScore});
  });
  roundtrips = _sortBy(roundtrips, (rt: any) => rt.score.score);

  // Event Handlers
  const selectDepartingFlightsStepperOnClick = () => {
    setStepperStep(0);
    setRoundtrip(null);
  }

  const departFlightOnSelect = (flight: any, _: any) => {
    setRoundtrip(_get(props.roundtrips, flight.id));
    setStepperStep(1);
  }

  const returnFlightOnSelect = (flight: any, bookingMetadata: any) => {
    props.onSelect(roundtrip.depart, flight, bookingMetadata)
  }

  // Renderers

  const renderStepper = () => {
    const texts = [
      <span
        className='cursor-pointer'
        onClick={selectDepartingFlightsStepperOnClick}
      >
        Select Departing Flights&nbsp;&nbsp;&gt;
      </span>,
      <span>Select Return Flights</span>
    ]
    return (
      <ol className={FlightsModalCss.RoundTripStepperCtn}>
        {texts.map((text: any, idx: number) => {
          const css = idx === stepperStep
            ? FlightsModalCss.RoundTripStepperActive: FlightsModalCss.RoundTripStepper
          return (<li className={css}>{text}</li>);
        })}
      </ol>
    );
  }

  const renderDepartingFlights = () => {
    const flights = roundtrips.map((rt: any) => rt.depart);
    return (
      <div>
        {flights.map((flight: any, idx: number) =>
          <FlightCard
            key={idx}
            flight={flight}
            onSelect={departFlightOnSelect}
            bookingMetadata={null}
          />
        )}
      </div>
    );
  }

  const renderReturnFlights = () => {
    roundtrip.returns.forEach((rt: any, idx: number) => {
      rt.bookingMetadata = roundtrip.bookingMetadata[idx];
    })
    const returnFlights = _sortBy(roundtrip.returns,
      (rt: any) => rt.bookingMetadata.score);

    return (
      <div>
        {returnFlights.map((flight: any, idx: number) =>
          <FlightCard
            key={idx}
            flight={flight}
            onSelect={returnFlightOnSelect}
            bookingMetadata={flight.bookingMetadata}
          />
        )}
      </div>
    );
  }

  return (
    <div>
      {renderStepper()}
      { stepperStep === 0 ? renderDepartingFlights(): null }
      { stepperStep === 1 ? renderReturnFlights(): null }
    </div>
  );
}


////////////////
// FlightCard //
////////////////

interface FlightCardProps {
  flight: any
  bookingMetadata: any
  onSelect: any
}

const FlightCard: FC<FlightCardProps> = (props: FlightCardProps) => {

  const [isExpanded, setIsExpanded] = useState(false);

  const airline = _get(props.flight.legs, "0.operatingAirline", {});
  const numStops = _get(props.flight, "legs", []).length === 1
    ? "Non-stop" : `${_get(props.flight, "legs", []).length - 1} stops`;

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
    return props.flight.legs.map((leg: any, idx: number) => {
      let layoverDuration = null;
      if (idx !== props.flight.legs.length - 1) {
        layoverDuration = intervalToDuration({
          start: parseISO(props.flight.legs[idx + 1].departure.datetime),
          end: parseISO(leg.arrival.datetime),
        });
      }

      return (
        <div key={idx}>
          <ol className={FlightsModalCss.FlightStopTimelineCtn}>
            <li className="mb-4 ml-6">
              <div className={FlightsModalCss.FlightStopTimelineIcon} />
              <h3 className={FlightsModalCss.FlightStopTimelineTime}>
                {printFmt(parseISO(leg.departure.datetime), "hh:mm aa")}
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
                {printFmt(parseISO(leg.arrival.datetime), "hh:mm aa")}
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
    const airline = _get(props.flight.legs, "0.operatingAirline", {});
    const imgUrl = `https://www.gstatic.com/flights/airline_logos/70px/${airline.code}.png`;
    const fallbackUrl = "https://cdn-icons-png.flaticon.com/512/4353/4353032.png";
    return (
      <div className="h-8 w-8 mr-4">
        <object className={CommonCss.IconLarge} data={imgUrl} type="image/png">
          <img className={CommonCss.IconLarge} src={fallbackUrl} alt={airline.name} />
        </object>
      </div>
    );
  }

  const renderPricePill = () => {
    if (_isEmpty(props.bookingMetadata)) {
      return (<></>);
    }
    const pill = (
      <span className={FlightsModalCss.FlightPricePill}>
        {props.bookingMetadata.price.currency}
        &nbsp;
        {props.bookingMetadata.price.amount}
      </span>
    );
    return (
      <a
        href={props.bookingMetadata.bookingURL}
        target="_blank"
        referrerPolicy="no-referrer"
      >
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
            <p className='font-medium'>
              {printFmt(parseISO(props.flight.departure.datetime), "hh:mm aa")}
            </p>
            <p className="text-xs text-slate-400">{props.flight.departure.airport.code}</p>
          </span>
          <ArrowLongRightIcon className='h-6 w-8' />
          <span className='mb-1'>
            <p className='font-medium'>
              {printFmt(parseISO(props.flight.arrival.datetime), "hh:mm aa")}
            </p>
            <p className="text-xs text-slate-400">{props.flight.arrival.airport.code}</p>
          </span>
        </div>
        <span className="text-xs text-slate-400 block mb-1">
          {prettyPrintMins(props.flight.duration)}
        </span>
        <span className="text-xs text-slate-400 block mb-1">
          {airline.name} | {renderNumStops()}
        </span>
        {isExpanded ? renderStopsInfo() : null}
      </div>
      <div className='flex flex-col text-right justify-between'>
        <div className='flex flex-row-reverse'>
          <PlusCircleIcon
            onClick={() => {props.onSelect(props.flight, props.bookingMetadata)}}
            className={FlightsModalCss.FlightPlusIcon}
          />
        </div>
        {renderPricePill()}
      </div>
    </div>
  );
}


///////////////////////////
// FlightsSearchForm //
///////////////////////////

const cabinClasses = [
  { label: "Economy", value: "economy" },
  { label: "Premium Economy", value: "premiumeconomy" },
  { label: "Business", value: "business" },
  { label: "First Class", value: "first" },
];

interface FlightsSearchFormProps {
  readonly trip: any
  onSearch: any
}

const FlightsSearchForm: FC<FlightsSearchFormProps> = (props: FlightsSearchFormProps) => {

  // Data State
  const [origin, setOrigIATA] = useState("");
  const [destination, setDestIATA] = useState("");
  const [cabinClass, setCabinClass] = useState(cabinClasses[0]);
  const [flightDates, setFlightDates] = useState<DateRange>();

  // UI State
  const [isCabinClassActive, setIsCabinClassActive] = useState(false);
  const [isOrigIATAFocus, setIsOrigIATAFocus] = useState(false);
  const [isDestIATAFocus, setIsDestIATAFocus] = useState(false);
  const [originQuery, setOrigIATAQuery] = useState("");
  const [destinationQuery, setDestIATAQuery] = useState("");


  useEffect(() => {
    setFlightDates({
      from: parseTripDate(props.trip.startDate || undefined),
      to: parseTripDate(props.trip.endDate || undefined),
    });
  }, [props.trip])

  // Event Handlers
  const searchBtnOnClick = () => {
    const departDate = printFromDateFromRange(flightDates, 'y-MM-dd');
    const arrDate = printToDateFromRange(flightDates, 'y-MM-dd');
    props.onSearch(origin, destination, departDate, arrDate, cabinClass.value);
  }

  const flightDatesOnSelect: SelectRangeEventHandler = (range?: DateRange) => {
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
            value={originQuery}
            onChange={(e) => {setOrigIATAQuery(e.target.value)}}
            onFocus={() => {setIsOrigIATAFocus(true)}}
            onBlur={(e) => {setTimeout(() => { setIsOrigIATAFocus(false) }, 200)}}
          />
          {renderAirportsAutocomplete(isOrigIATAFocus, originQuery, setOrigIATA, setOrigIATAQuery, setIsOrigIATAFocus)}
        </div>
        <div className="relative w-6/12">
          <div className={FlightsModalCss.FlightFromIconCtn}>
            <PaperAirplaneIcon className={FlightsModalCss.FlightFromIcon} />
          </div>
          <input
            type="text"
            className={FlightsModalCss.FlightFromInput}
            placeholder="to city, airport"
            value={destinationQuery}
            onChange={(e) => { setDestIATAQuery(e.target.value)}}
            onFocus={() => {setIsDestIATAFocus(true)}}
            onBlur={(e) => {setTimeout(() => { setIsDestIATAFocus(false) }, 200)}}
          />
          {renderAirportsAutocomplete(isDestIATAFocus, destinationQuery, setDestIATA, setDestIATAQuery, setIsDestIATAFocus)}
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


/////////////////
// FlightModal //
/////////////////

interface FlightsModalProps {
  readonly trip: any
  isOpen: boolean
  onClose: any
  onFlightSelect: any
}

const FlightsModal: FC<FlightsModalProps> = (props: FlightsModalProps) => {

  const [oneways, setOneways] = useState([] as any);
  const [roundtrips, setRoundtrips] = useState({} as any);

  const [isLoading, setIsLoading] = useState(false);
  const [searchInitiated, setSearchInitiated] = useState(false);
  const [alertMsg, setAlertMsg] = useState("");


  // Event Handlers
  const onSearch = (origin: string, destination: string, departDate: string, returnDate: string | undefined, cabinClass: string) => {
    if (_isEmpty(origin)) {
      setAlertMsg("Please select a flight origin");
      return;
    }
    if (_isEmpty(destination)) {
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

    FlightsAPI.search(origin, destination, departDate, returnDate, cabinClass)
    .then(res => {
      if (_isEmpty(returnDate)) {
        const oneways = _sortBy(_get(res, "data.itineraries.oneways", []), "bookingMetadata.score")
        setOneways(oneways);
        setRoundtrips({});
      } else {
        const roundtrips = _get(res, "data.itineraries.roundtrips", {})
        setRoundtrips(roundtrips);
        setOneways([]);
      }
    })
    .finally(() => {
      setIsLoading(false);
    });
  }

  const onSelectOnewayFlight = (depart: any, bookingMetadata: any) => {
    const tripFlight: Trips.Flight = {
      id: uuidv4(),
      type: "flight",
      tags: new Map<string, string>(),
      labels: new Map<string, string>(),
      itineraryType: "oneway",
      depart,
      return: {} as any,
      price: bookingMetadata.price,
    };
    props.onFlightSelect(tripFlight);
  }

  const onSelectRoundTripFlight = (departFlight: any, returnFlight: any, bookingMetadata: any) => {
    const tripFlight: Trips.Flight = {
      id: uuidv4(),
      type: "flight",
      tags: new Map<string, string>(),
      labels: new Map<string, string>(),
      itineraryType: "roundtrip",
      depart: departFlight,
      return: returnFlight,
      price: bookingMetadata.price,
    };
    props.onFlightSelect(tripFlight);
  }

  // Renderers
  const renderAlert = () => {
    return !_isEmpty(alertMsg) ? <Alert title={""} message={alertMsg} /> : null
  }

  const renderItineraries = () => {
    if (!searchInitiated) {
      return (<></>);
    }

    if (isLoading) {
      return <Spinner />
    }

    if (!_isEmpty(oneways)) {
      return (
        <OnewayFlightsContainer
          oneways={oneways}
          onSelect={onSelectOnewayFlight} />
      );
    }
    return (
      <RoundtripFlightsContainer
        roundtrips={roundtrips}
        onSelect={onSelectRoundTripFlight}/>
    );
  }

  if (!props.isOpen) {
    return (<></>);
  }

  return (
    <Modal isOpen={props.isOpen}>
      <div className={FlightsModalCss.Ctn}>
        <div className={FlightsModalCss.Wrapper}>
          <h2 className={FlightsModalCss.Header}>
            Search flights
          </h2>
          <button type="button" onClick={props.onClose}>
            <XMarkIcon className={FlightsModalCss.CloseIcon} />
          </button>
        </div>
        {renderAlert()}
        <FlightsSearchForm trip={props.trip} onSearch={onSearch} />
        {renderItineraries()}
      </div>
    </Modal>
  );
}

export default FlightsModal;
