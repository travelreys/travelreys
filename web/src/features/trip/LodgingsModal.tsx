import React, { FC, useState, useEffect } from 'react';
import _get from "lodash/get";
import _sortBy from "lodash/sortBy";
import _isEmpty from "lodash/isEmpty";
import { useDebounce } from 'usehooks-ts';
import { v4 as uuidv4 } from 'uuid';
import { DateRange } from 'react-day-picker';

import {
  MagnifyingGlassIcon,
  XMarkIcon
} from '@heroicons/react/24/outline'

import Alert from '../../components/common/Alert';
import InputDatesPicker from '../../components/common/InputDatesPicker';
import Modal from '../../components/common/Modal';

import MapsAPI, { EMBED_MAPS_APIKEY, placeFields } from '../../apis/maps';
import { LodgingsModalCss } from '../../assets/styles/global';
import { Trips } from '../../lib/trips';
import { parseTripDate } from '../../lib/dates';
import PlaceAutocomplete from '../maps/PlaceAutocomplete';

// TripLodgingsModal

interface TripLodgingsModalProps {
  readonly trip: any
  isOpen: boolean
  onClose: any
  onLodgingSelect: any
}

const TripLodgingsModal: FC<TripLodgingsModalProps> = (props: TripLodgingsModalProps) => {

  const [sessionToken, setSessionToken] = useState("");
  const [searchQuery, setSearchQuery] = useState("");
  const [predictions, setPredictions] = useState([] as any);
  const [selectedPlaceID, setSelectedPlaceID] = useState("");
  const [selectedPlace, setSelectedPlace] = useState(null as any);
  const [checkinDates, setCheckinDates] = useState<DateRange>();

  const [alertMsg, setAlertMsg] = useState("");
  const debouncedValue = useDebounce<string>(searchQuery, 500);

  useEffect(() => {
    if (!_isEmpty(debouncedValue)) {
      autocomplete(debouncedValue);
    }
  }, [debouncedValue]);

  useEffect(() => {
    setCheckinDates({
      from: parseTripDate(props.trip.startDate || undefined),
      to: parseTripDate(props.trip.endDate || undefined),
    });
  }, [props.trip])


  // API
  const autocomplete = (query: string) => {
    let token = sessionToken;
    if (_isEmpty(token)) {
      token = uuidv4();
      setSessionToken(token);
    }
    MapsAPI.placeAutocomplete(query, ["lodging"], token)
    .then((res) => {
      setPredictions(_get(res, "data.predictions", []))
    });
  }

  const getPlaceDetails = (placeID: string) => {
    MapsAPI.placeDetails(placeID, placeFields, sessionToken)
    .then((res) => {
      setPredictions([]);
      setSelectedPlaceID(placeID);
      setSelectedPlace(_get(res, "data.place", null));
    })
    .finally(() => {
      setSessionToken("");
    });
  }

  // Event Handlers

  const predictionOnSelect = (placeID: string) => {
    getPlaceDetails(placeID);
  }

  const addLodgingBtnOnClick = () => {
    const tripLodging: Trips.Lodging = {
      id: uuidv4(),
      checkinTime: checkinDates?.from,
      checkoutTime: checkinDates?.to,
      place: selectedPlace,
      priceMetadata: {} as any,
      tags: new Map<string,string>(),
      labels: new Map<string,string>(),
    };
    props.onLodgingSelect(tripLodging);
  }

  // Renderers
  const renderAlert = () => {
    return !_isEmpty(alertMsg) ? <Alert title={""} message={alertMsg} /> : null
  }

  const renderSearchInput = () => {
    return (
      <div className="relative mb-2">
        <div className={LodgingsModalCss.SearchIconCtn}>
          <MagnifyingGlassIcon className={LodgingsModalCss.SearchIcon} />
        </div>
        <input
          type="text"
          className={LodgingsModalCss.SearchInput}
          placeholder="name, address"
          value={searchQuery}
          onChange={(e) => {setSearchQuery(e.target.value)}}
        />
      </div>
    );
  }

  const renderMapForSelectedPlace = () => {
    if (_isEmpty(selectedPlaceID)) {
      return (<></>);
    }
    const src = `https://www.google.com/maps/embed/v1/place?key=${EMBED_MAPS_APIKEY}&q=place_id:${selectedPlaceID}`;

    return (
      <div className='flex flex-col pb-2'>
        <iframe
          className="w-full mb-2"
          style={{border: "0"}}
          referrerPolicy="no-referrer-when-downgrade"
          src={src}
          allowFullScreen
        ></iframe>
        <div className='flex justify-around'>
          <button
            type="button"
            className={LodgingsModalCss.AddBtn}
            onClick={addLodgingBtnOnClick}
          >
            Add Lodging
          </button>
        </div>
      </div>
    );
  }

  return (
    <Modal isOpen={props.isOpen}>
      <div className={LodgingsModalCss.Ctn}>
        <div className={LodgingsModalCss.Wrapper}>
          <h2 className={LodgingsModalCss.Header}>
            Search hotels or lodgings
          </h2>
          <button type="button" onClick={props.onClose}>
            <XMarkIcon className={LodgingsModalCss.CloseIcon} />
          </button>
        </div>
        {renderAlert()}
        {renderSearchInput()}
        <InputDatesPicker
          onSelect={(range: DateRange) => { setCheckinDates(range); }}
          dates={checkinDates}
          WrapperCss={"mb-2"}
          CtnCss={LodgingsModalCss.InputDatesCtn}
        />
        <PlaceAutocomplete
          predictions={predictions}
          onSelect={predictionOnSelect}
        />
        {renderMapForSelectedPlace()}
      </div>
    </Modal>
  );
}

export default TripLodgingsModal;
