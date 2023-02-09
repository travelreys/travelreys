import React, { FC, useEffect, useState } from 'react';
import _get from "lodash/get";
import _sortBy from "lodash/sortBy";
import _isEmpty from "lodash/isEmpty";
import { DateRange, SelectRangeEventHandler } from 'react-day-picker';
import { v4 as uuidv4 } from 'uuid';

import {
  ChevronDownIcon,
  ChevronUpIcon,
  MapPinIcon,
  PaperAirplaneIcon,
  XMarkIcon,
} from '@heroicons/react/24/outline'
import { FlightsModalCss } from '../../styles/global';

import FlightsAPI from '../../apis/flights';
import { Trips } from '../../apis/types';
import Alert from '../Alert';
import Modal from '../Modal';
import InputDatesPicker from '../InputDatesPicker';
import Spinner from '../../components/Spinner';
import OnewayFlightsContainer from './OnewayFlightsContainer';
import RoundtripFlightsContainer from './RoundtripFlightsContainer';
import {
  printFromDateFromRange,
  printToDateFromRange,
  parseTripDate
} from '../../utils/dates';
import { capitaliseWords } from '../../utils/strings';

// TripFlightsSearchForm

const cabinClasses = [
  { label: "Economy", value: "economy" },
  { label: "Premium Economy", value: "premiumeconomy" },
  { label: "Business", value: "business" },
  { label: "First Class", value: "first" },
];

interface TripFlightsSearchFormProps {
  readonly trip: any
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
    props.onSearch(origIATA, destIATA, departDate, arrDate, cabinClass.value);
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
  readonly trip: any
  isOpen: boolean
  onClose: any
  onFlightSelect: any
}

const TripFlightsModal: FC<TripFlightsModalProps> = (props: TripFlightsModalProps) => {

  const [oneways, setOneways] = useState([] as any);
  const [roundtrips, setRoundtrips] = useState({} as any);

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
      price: bookingMetadata.priceMetadata,
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
      price: bookingMetadata.priceMetadata,
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
        <TripFlightsSearchForm trip={props.trip} onSearch={onSearch} />
        {renderItineraries()}
      </div>
    </Modal>
  );
}

export default TripFlightsModal;
