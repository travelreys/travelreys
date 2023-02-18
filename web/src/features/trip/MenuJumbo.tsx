import React, { ChangeEvent, FC, useEffect, useState } from 'react';
import _get from "lodash/get";
import _isEmpty from "lodash/isEmpty";
import { DateRange, SelectRangeEventHandler } from 'react-day-picker';
import {
  CalendarDaysIcon,
  EllipsisHorizontalCircleIcon,
  MagnifyingGlassIcon,
  PencilIcon,
  XMarkIcon,
} from '@heroicons/react/24/outline'
import { Cog6ToothIcon } from '@heroicons/react/24/solid';


import Dropdown from '../../components/common/Dropdown';
import Modal from '../../components/common/Modal';
import DatesPicker from '../../components/common/DatesPicker';
import Spinner from '../../components/common/Spinner';

import TripsSyncAPI from '../../apis/tripsSync';
import ImagesAPI from '../../apis/images';
import {
  nullDate,
  printFromDateFromRange,
  printToDateFromRange,
  parseTripDate
} from '../../lib/dates';
import {
  DefaultTransportationPreference,
  LabelTransportationPreference
} from '../../lib/trips';
import { CommonCss, TripMenuJumboCss } from '../../assets/styles/global';


/////////////////////
// CoverImageModal //
/////////////////////
interface CoverImageModalProps {
  isOpen: boolean
  onClose: any
  onCoverImageSelect: any
}

const CoverImageModal: FC<CoverImageModalProps> = (props: CoverImageModalProps) => {

  const [query, setQuery] = useState("");
  const [imageList, setImageList] = useState([] as any);
  const [isLoading, setIsLoading] = useState(false);

  // API
  const searchImage = () => {
    setIsLoading(true);
    ImagesAPI.search(query)
    .then(res => {
      const images = _get(res, "data.images", []);
      setImageList(images);
      setIsLoading(false);
    });
  }

  // Event Handlers

  // Renderers
  const renderImageThumbnails = () => {
    if (isLoading) {
      return <Spinner />
    }
    return (
      <div className='columns-2 md:columns-3'>
        { imageList.map((image: any) => (
          <figure
            key={image.id}
            className={TripMenuJumboCss.Figure}
          >
            <a href="#">
              <img
                srcSet={ImagesAPI.makeSrcSet(image)}
                src={ImagesAPI.makeSrc(image)}
                className={TripMenuJumboCss.FigureImg}
              />
              <div className={TripMenuJumboCss.FigureBtnCtn}>
                <a
                  className={TripMenuJumboCss.FigureBtn}
                  onClick={() => {props.onCoverImageSelect(image)}}
                  href="#"
                >
                  Select
                </a>
              </div>
            </a>
            <figcaption className={TripMenuJumboCss.FigureCaption}>
              <a
                target="_blank"
                href={ImagesAPI.makeUserURL(_get(image, "user.username"))}
              >
                @{_get(image, "user.username")}, Unsplash
              </a>
            </figcaption>
          </figure>
        ))}
      </div>);
  }

  return (
    <Modal isOpen={props.isOpen}>
      <div className={TripMenuJumboCss.SearchImageCard}>
        <div className='flex justify-between mb-6'>
          <h2 className={TripMenuJumboCss.SearchImageTitle}>
            Change cover image
          </h2>
          <button type="button" onClick={props.onClose}>
            <XMarkIcon className='h-6 w-6 text-slate-700' />
          </button>
        </div>
        <h2 className={TripMenuJumboCss.SearchImageWebTitle}>
          Search the web
        </h2>
        <div className="flex mb-4 justify-between">
          <input
            type="text"
            className={TripMenuJumboCss.SearchImageInput}
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            onKeyDown={(e) => e.key === "Enter" ? searchImage() : ""}
            placeholder="destination, theme ..."
          />
          <button
            type='button'
            className={TripMenuJumboCss.SearchImageBtn}
            onClick={searchImage}
          >
            <MagnifyingGlassIcon className={TripMenuJumboCss.SearchImageIcon} />
          </button>
        </div>
        {renderImageThumbnails()}
      </div>
    </Modal>
  );

}


///////////////////
// SettingsModal //
///////////////////

interface SettingsModalProps {
  trip: any
  isOpen: boolean
  onClose: () => void
  onTransportationPreferenceChange: (pref: string) => void
}

const SettingsModal: FC<SettingsModalProps> = (props: SettingsModalProps) => {

  const [transportationPreference, setTransportationPreference] = useState(DefaultTransportationPreference);

  useEffect(() => {
    setTransportationPreference(_get(
      props.trip, `labels.${LabelTransportationPreference}`, DefaultTransportationPreference
    ))
  }, [props.trip])

  // Event Handlers
  const transportationPreferenceOnChange = (e: any) => {
    setTransportationPreference(e.target.value)
    props.onTransportationPreferenceChange(e)
  }

  //  Renderers

  const renderTransportationMode = () => {
    return (
      <div className='mb-2'>
        <label
          htmlFor="transportation"
          className="block mb-2 text-sm font-bold text-gray-900"
        >
          Transportation mode preference
        </label>
        <select
          id="transportation"
          value={transportationPreference}
          onChange={transportationPreferenceOnChange}
          className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 "
        >
          <option value="walk+drive">Walk short distances + Drive</option>
          <option value="walk+transit">Walk short distances + Public Transport</option>
        </select>
      </div>
    );
  }

  return (
    <Modal isOpen={props.isOpen}>
      <div className={TripMenuJumboCss.SearchImageCard}>
        <h2 className='text-xl font-bold mb-2'>Trip Settings</h2>
        {renderTransportationMode()}
        <button
          type="button"
          className='bg-indigo-500 mt-4 px-8 py-2 font-bold text-white rounded-full'
          onClick={props.onClose}
        >
          Close
        </button>
      </div>
    </Modal>
  );
}



///////////////
// MenuJumbo //
///////////////
interface MenuJumboProps {
  trip: any
  tripStateOnUpdate: any
}

const MenuJumbo: FC<MenuJumboProps> = (props: MenuJumboProps) => {

  // State
  const [tripName, setTripName] = useState<string>("");
  const [startDt, setStartDt] = useState<Date|undefined>();
  const [endDt, setEndDt] = useState<Date|undefined>();

  // UI State
  const [isCoverImageModalOpen, setIsCoverImageModalOpen] = useState(false);
  const [isCalendarOpen, setIsCalendarOpen] = useState(false);
  const [isSettingsModalOpen, setIsSettingsModalOpen] = useState(false);

  // When props.trip changes, need to update the ui state
  useEffect(() => {
    setTripName(props.trip.name);
    setStartDt(parseTripDate(props.trip.startDate));
    setEndDt(parseTripDate(props.trip.endDate));
  }, [props.trip])


  // Event Handlers - Trip Name
  const tripNameOnBlur = () => {
    props.tripStateOnUpdate([TripsSyncAPI.newReplaceOp("/name", tripName)])
  }

  // Event Handlers - Cover Image
  const coverImageOnSelect = (image: any) => {
    props.tripStateOnUpdate([TripsSyncAPI.newReplaceOp("/coverImage", image)]);
  }

  // Event Handlers - Trip Dates
  const tripDatesOnChange: SelectRangeEventHandler = (range: DateRange | undefined) => {
    setStartDt(range?.from);
    setEndDt(range?.to);
  }

  const tripDatesOnBlur = (e: any) => {
    const range = {from: startDt, to: endDt};
    if (!e.currentTarget.contains(e.relatedTarget) && isCalendarOpen) {
      const ops = [];
      const from = range.from || nullDate;
      const to = range.to || nullDate;
      ops.push(TripsSyncAPI.newReplaceOp("/startDate", from));
      ops.push(TripsSyncAPI.newReplaceOp("/endDate", to));
      props.tripStateOnUpdate(ops);
      setIsCalendarOpen(false);
      return;
    }
  }

  // Event Handlers - Trip Settings
  const transportationPreferenceOnChange = (e: any) => {
    const op = _get(props.trip, `/labels.${LabelTransportationPreference}`)
      ? TripsSyncAPI.newReplaceOp: TripsSyncAPI.makeAddOp;
    props.tripStateOnUpdate([
      op(`/labels/${LabelTransportationPreference}`, e.target.value)
    ]);
  }

  // Renderers
  const renderCoverImage = () => {
    return (
      <div className="relative">
        <img
          srcSet={ImagesAPI.makeSrcSet(props.trip.coverImage)}
          src={ImagesAPI.makeSrc(props.trip.coverImage)}
          className={TripMenuJumboCss.TripCoverImage}
        />
        <button
          type='button'
          className={TripMenuJumboCss.TripImageEditIconCtn}
          onClick={() => { setIsCoverImageModalOpen(true) }}
        >
          <PencilIcon className={TripMenuJumboCss.TripImageEditIcon} />
        </button>
      </div>
    );
  }

  const renderDatesButton = () => {
    const range = {from: startDt, to: endDt}
    return (
      <div onBlur={tripDatesOnBlur}>
        <button
          type="button"
          className={TripMenuJumboCss.TripDatesBtn}
          onClick={() => { setIsCalendarOpen(true) }}
        >
          <CalendarDaysIcon className={TripMenuJumboCss.TripDatesBtnIcon} />
          {startDt ?
            <span>
              {printFromDateFromRange(range, "MMM d, yy ")}
              &nbsp;-&nbsp;
              {printToDateFromRange(range, "MMM d, yy ")}
            </span> : null}
        </button>
        <DatesPicker
          onSelect={tripDatesOnChange}
          isOpen={isCalendarOpen}
          dates={{from: startDt, to: endDt}}
        />
      </div>
    );
  }

  const renderSettingsDropdown = () => {
    const opts = [
      (<button
        type='button'
        className={TripMenuJumboCss.SettingsBtn}
        onClick={() => setIsSettingsModalOpen(true)}
      >
        <Cog6ToothIcon className={CommonCss.LeftIcon} />
        Settings
      </button>),
    ];
    const menu = (
      <EllipsisHorizontalCircleIcon
        className={CommonCss.DropdownIcon} />
    );
    return <Dropdown menu={menu} opts={opts} />
  }

  const renderTripNameInput = () => {
    return (
      <div className={TripMenuJumboCss.TripNameInputCtn}>
        <div className={TripMenuJumboCss.TripNameInputWrapper}>
          <div className='flex items-center mb-12'>
            <input
              type="text"
              value={tripName}
              onChange={(e) => setTripName(e.target.value)}
              onBlur={tripNameOnBlur}
              className={TripMenuJumboCss.TripNameInput}
            />
            {renderSettingsDropdown()}
          </div>
          {renderDatesButton()}
        </div>
      </div>
    );
  }

  return (
    <>
      <div className='bg-indigo-100'>
        {renderCoverImage()}
        {renderTripNameInput()}
      </div>
      <CoverImageModal
        isOpen={isCoverImageModalOpen}
        onClose={() => {setIsCoverImageModalOpen(false)}}
        onCoverImageSelect={coverImageOnSelect}
      />
      <SettingsModal
        trip={props.trip}
        isOpen={isSettingsModalOpen}
        onClose={() => {setIsSettingsModalOpen(false)}}
        onTransportationPreferenceChange={transportationPreferenceOnChange}
      />
    </>
  );
}

export default MenuJumbo;
