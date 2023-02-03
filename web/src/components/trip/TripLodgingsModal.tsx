import React, { FC, useState, useEffect } from 'react';
import _get from "lodash/get";
import _sortBy from "lodash/sortBy";
import _isEmpty from "lodash/isEmpty";
import { useDebounce } from 'usehooks-ts';
import { v4 as uuidv4 } from 'uuid';

import {
  MagnifyingGlassIcon,
  MapPinIcon,
  XMarkIcon
} from '@heroicons/react/24/outline'
import { LodgingsModalCss, ModalCss } from '../../styles/global';

import Alert from '../Alert';
import Spinner from '../Spinner';

import MapsAPI, { EMBED_MAPS_APIKEY } from '../../apis/maps';
import { DateRange } from 'react-day-picker';
import InputDatesPicker from '../InputDatesPicker';


// TripLodgingsModal

interface TripLodgingsModalProps {
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
    const fields = [
      "address_component", "adr_address", "business_status", "formatted_address", "geometry",  "name", "photos", "place_id", "types", "utc_offset",
      "opening_hours", "formatted_phone_number", "international_phone_number", "website",
    ];
    MapsAPI.placeDetails(placeID, fields, sessionToken)
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
  useEffect(() => {
    if (!_isEmpty(debouncedValue)) {
      autocomplete(debouncedValue);
    }
  }, [debouncedValue])


  const predictionOnSelect = (placeID: string) => {
    getPlaceDetails(placeID);
  }

  const addLodgingBtnOnClick = () => {
    const tripLodging = {
      checkinTime: checkinDates?.from,
      checkoutTime: checkinDates?.to,
      place: selectedPlace,
      priceMetadata: {},
    };
    props.onLodgingSelect(tripLodging);
  }

  // Renderers
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

  const renderAutocomplete = () => {
    if (_isEmpty(predictions)) {
      return (<></>);
    }
    return (
      <div className='p-1'>
        {predictions.map((pre: any) => (
          <div
            className='flex items-center mb-4 cursor-pointer group'
            key={pre.place_id}
            onClick={() => {predictionOnSelect(pre.place_id)}}
          >
            <div className='p-1 group-hover:text-indigo-500'>
              <MapPinIcon className='h-6 w-6' />
            </div>
            <div className='ml-1'>
              <p className='text-slate-900 group-hover:text-indigo-500 text-sm font-medium'>
                {_get(pre, "structured_formatting.main_text", "")}
              </p>
              <p className="text-slate-400 group-hover:text-indigo-500 text-xs">
                {_get(pre, "structured_formatting.secondary_text", "")}
              </p>
            </div>
          </div>
        ))}
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

  if (!props.isOpen) {
    return (<></>);
  }

  return (
    <div className={ModalCss.Container}>
      <div className={ModalCss.Inset}></div>
      <div className={ModalCss.Content}>
        <div className={ModalCss.ContentContainer}>
          <div className={ModalCss.ContentCard}>
            <div className="px-4 pt-5 sm:p-8 sm:pb-2 rounded-t-lg mb-4">
              <div className='flex justify-between mb-6'>
                <h2 className="text-xl sm:text-2xl font-bold text-center text-slate-900">
                  Search hotels or lodgings
                </h2>
                <button type="button" onClick={props.onClose}>
                  <XMarkIcon className='h-6 w-6 text-slate-700' />
                </button>
              </div>
              {!_isEmpty(alertMsg) ? <Alert title={""} message={alertMsg} /> : null}
              <InputDatesPicker
                onSelect={(range: DateRange) => { setCheckinDates(range); }}
                dates={checkinDates}
                WrapperCss={"mb-2"}
                CtnCss={"flex w-full border border-slate-200 rounded"}
              />
              {renderSearchInput()}
              {renderAutocomplete()}
              {renderMapForSelectedPlace()}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export default TripLodgingsModal;
