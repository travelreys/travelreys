import React, {
  ChangeEvent,
  FC,
  useEffect,
  useState,
} from 'react';
import _get from "lodash/get";
import _isEmpty from "lodash/isEmpty";
import { parseJSON, isEqual } from 'date-fns';
import { DateRange, SelectRangeEventHandler } from 'react-day-picker';
import {
  CalendarDaysIcon,
  MagnifyingGlassIcon,
  PencilIcon,
  XMarkIcon,
} from '@heroicons/react/24/outline'

import Spinner from '../../components/Spinner';
import { printFromDateFromRange, printToDateFromRange } from '../../utils/dates';
import TripsSyncAPI from '../../apis/tripsSync';
import ImagesAPI, { stockImageSrc } from '../../apis/images';

import { ModalCss, TripMenuJumboCss } from '../../styles/global';
import DatesPicker from '../DatesPicker';


// CoverImageModal

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
      const images = _get(res, "data.images");
      setImageList(images);
      setIsLoading(false);
    });
  }

  // Event Handlers
  const searchImageQueryOnChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setQuery(event.target.value);
  }

  const searchImageQueryOnEnter = (event: React.KeyboardEvent<HTMLInputElement>) => {
    if (event.key === "Enter") {
      searchImage();
    }
  }

  // Renderers
  const renderImageThumbnails = () => {
    if (isLoading) {
      return <Spinner />
    }

    return (
      <div className='columns-2 md:columns-3'>
        { imageList.map((image: any) => (
          <figure className={TripMenuJumboCss.Figure}>
            <a href="#">
              <img key={image.id}
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
              <a target="_blank" href={ImagesAPI.makeUserReferURL(_get(image, "user.username"))}>
                @{_get(image, "user.username")}, Unsplash
              </a>
            </figcaption>
          </figure>
        ))}
      </div>
    );
  }

  if (!props.isOpen) {
    return <></>;
  }

  return (
    <div className={ModalCss.Container}>
      <div className={ModalCss.Inset}></div>
      <div className={ModalCss.Content}>
        <div className={ModalCss.ContentContainer}>
          <div className={ModalCss.ContentCard}>
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
                  onChange={searchImageQueryOnChange}
                  onKeyDown={searchImageQueryOnEnter}
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
          </div>
        </div>
      </div>
    </div>
  );

}


// TripMenuJumboProps


interface TripMenuJumboProps {
  trip: any
  tripStateOnUpdate: any
}

const nullDate = parseJSON("0001-01-01T00:00:00Z");
const parseTripDate = (tripDate: string | undefined) => {
  if (_isEmpty(tripDate)) {
    return undefined;
  }
  return isEqual(parseJSON(tripDate!), nullDate) ? undefined : parseJSON(tripDate!);
}

const TripMenuJumbo: FC<TripMenuJumboProps> = (props: TripMenuJumboProps) => {

  // State
  const [tripName, setTripName] = useState(props.trip.name);

  const [tripStartDate, setTripStartDate] = useState(parseTripDate(props.trip.startDate));
  const [tripEndDate, setTripEndDate] = useState(parseTripDate(props.trip.endDate));

  // UI State
  const [isCoverImageModalOpen, setIsCoverImageModalOpen] = useState(false);
  const [isCalendarOpen, setIsCalendarOpen] = useState(false);

  // When props.trip changes, need to update the name
  useEffect(() => {
    setTripName(props.trip.name);
    setTripStartDate(parseTripDate(props.trip.startDate));
    setTripEndDate(parseTripDate(props.trip.endDate));
  }, [props.trip])



  // Event Handlers - Trip Name
  const tripNameOnChange = (event: ChangeEvent<HTMLInputElement>) => {
    setTripName(event.target.value)
  }

  const tripNameOnBlur = () => {
    const ops = [TripsSyncAPI.makeJSONPatchOp("replace", "/name", tripName)]
    props.tripStateOnUpdate(ops)
  }

  // Event Handlers - Cover Image
  const coverImageOnSelect = (image: any) => {
    const ops = [
      TripsSyncAPI.makeJSONPatchOp(
        "replace", "/coverImage", image)
    ];
    props.tripStateOnUpdate(ops);
  }

  // Event Handlers - Trip Dates
  const tripDatesOnChange: SelectRangeEventHandler = (range: DateRange | undefined) => {
    setTripStartDate(range?.from);
    setTripEndDate(range?.to);
  }

  const tripDatesOnBlur = (e: any) => {
    const range = {from: tripStartDate, to: tripEndDate};
    if (!e.currentTarget.contains(e.relatedTarget) && isCalendarOpen) {
      const ops = [];

      if (range.from !== undefined) {
        ops.push(TripsSyncAPI.makeJSONPatchOp("replace", "/startDate", range.from));
      } else {
        ops.push(TripsSyncAPI.makeJSONPatchOp("replace", "/startDate", nullDate));
      }
      if (range.to !== undefined) {
        ops.push(TripsSyncAPI.makeJSONPatchOp("replace", "/endDate", range.to));
      } else {
        ops.push(TripsSyncAPI.makeJSONPatchOp("replace", "/endDate", nullDate));
      }
      props.tripStateOnUpdate(ops);
      setIsCalendarOpen(false);
      return;
    }
  }

  // Renderers

  const renderDatesButton = () => {
    const range = {from: tripStartDate, to: tripEndDate}
    return (
      <div onBlur={tripDatesOnBlur}>
        <button
          type="button"
          className={TripMenuJumboCss.TripDatesBtn}
          onClick={() => { setIsCalendarOpen(true) }}
        >
          <CalendarDaysIcon className={TripMenuJumboCss.TripDatesBtnIcon} />
          &nbsp;&nbsp;
          {tripStartDate ?
            <span>
              {printFromDateFromRange(range, "MMM d, yy ")}
              &nbsp;-&nbsp;
              {printToDateFromRange(range, "MMM d, yy ")}
            </span> : null}
        </button>
        <DatesPicker
          onSelect={tripDatesOnChange}
          isOpen={isCalendarOpen}
          dates={{from: tripStartDate, to: tripEndDate}}
        />
      </div>
    );
  }

  return (
    <>
      <div className='bg-indigo-100'>
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
        <div className={TripMenuJumboCss.TripNameInputCtn}>
          <div className={TripMenuJumboCss.TripNameInputWrapper}>
            <input
              type="text"
              value={tripName}
              onChange={tripNameOnChange}
              onBlur={tripNameOnBlur}
              className={TripMenuJumboCss.TripNameInput}
            />
            {renderDatesButton()}
          </div>
        </div>
      </div>
      <CoverImageModal
        onClose={() => {setIsCoverImageModalOpen(false)}}
        onCoverImageSelect={coverImageOnSelect}
        isOpen={isCoverImageModalOpen}
      />
    </>
  );
}

export default TripMenuJumbo;
