import React, { FC, useEffect, useState } from 'react';
import _get from "lodash/get";
import { DateRange, SelectRangeEventHandler } from 'react-day-picker';
import { useTranslation } from 'react-i18next';
import {
  CalendarDaysIcon,
  MagnifyingGlassIcon,
  PencilIcon,
  XMarkIcon,
} from '@heroicons/react/24/outline'

import Modal from '../../components/common/Modal';
import DatesPicker from '../../components/common/DatesPicker';
import Spinner from '../../components/common/Spinner';

import ImagesAPI from '../../apis/images';
import {
  nullDate,
  printFromDateFromRange,
  printToDateFromRange,
  parseTripDate
} from '../../lib/dates';
import { makeReplaceOp } from '../../lib/tripsSync';
import { TripMenuJumboCss } from '../../assets/styles/global';



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
            <button type="button">
              <img
                srcSet={ImagesAPI.makeSrcSet(image)}
                src={ImagesAPI.makeSrc(image)}
                alt={"cover"}
                className={TripMenuJumboCss.FigureImg}
              />
              <div className={TripMenuJumboCss.FigureBtnCtn}>
                <button
                  type="button"
                  className={TripMenuJumboCss.FigureBtn}
                  onClick={() => {props.onCoverImageSelect(image)}}
                >
                  Select
                </button>
              </div>
            </button>
            <figcaption className={TripMenuJumboCss.FigureCaption}>
              <a
                target="_blank"
                href={ImagesAPI.makeUserURL(_get(image, "user.username"))}
                rel="noreferrer"
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


///////////////
// MenuJumbo //
///////////////
interface MenuJumboProps {
  trip: any
  tripStateOnUpdate: any
}

const MenuJumbo: FC<MenuJumboProps> = (props: MenuJumboProps) => {

  // const {t} = useTranslation();

  // State
  const [tripName, setTripName] = useState<string>("");
  const [startDt, setStartDt] = useState<Date|undefined>();
  const [endDt, setEndDt] = useState<Date|undefined>();

  // UI State
  const [isCoverImageModalOpen, setIsCoverImageModalOpen] = useState(false);
  const [isCalendarOpen, setIsCalendarOpen] = useState(false);

  // When props.trip changes, need to update the ui state
  useEffect(() => {
    setTripName(props.trip.name);
    setStartDt(parseTripDate(props.trip.startDate));
    setEndDt(parseTripDate(props.trip.endDate));
  }, [props.trip])

  // Event Handlers - Trip Name
  const tripNameOnBlur = () => {
    props.tripStateOnUpdate([makeReplaceOp("/name", tripName)])
  }

  // Event Handlers - Cover Image
  const coverImageOnSelect = (image: any) => {
    props.tripStateOnUpdate([makeReplaceOp("/coverImage", image)]);
  }

  // Event Handlers - Trip Dates
  const tripDatesOnChange: SelectRangeEventHandler = (range?: DateRange) => {
    setStartDt(range?.from);
    setEndDt(range?.to);
  }

  const tripDatesOnBlur = (e: any) => {
    const range = {from: startDt, to: endDt};
    if (!e.currentTarget.contains(e.relatedTarget) && isCalendarOpen) {
      const ops = [];
      const from = range.from || nullDate;
      const to = range.to || nullDate;
      ops.push(makeReplaceOp("/startDate", from));
      ops.push(makeReplaceOp("/endDate", to));
      props.tripStateOnUpdate(ops);
      setIsCalendarOpen(false);
      return;
    }
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
    const range = {from: startDt, to: endDt};
    const dateFmt = "MMM d, yy"
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
              {printFromDateFromRange(range, dateFmt)}
              &nbsp;-&nbsp;
              {printToDateFromRange(range, dateFmt)}
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

  const renderTripNameInput = () => {
    return (
      <div className={TripMenuJumboCss.TripNameInputCtn}>
        <div className={TripMenuJumboCss.TripNameInputWrapper}>
          <div className={TripMenuJumboCss.TripNamInputHeader}>
            <input
              type="text"
              value={tripName}
              onChange={(e) => setTripName(e.target.value)}
              onBlur={tripNameOnBlur}
              className={TripMenuJumboCss.TripNameInput}
            />
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
    </>
  );
}

export default MenuJumbo;
