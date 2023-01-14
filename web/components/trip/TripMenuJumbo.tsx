import React, {
  ChangeEvent,
  FC,
  useEffect,
  useState,
} from 'react';
import _get from "lodash/get";
import { parseJSON, parseISO, isEqual } from 'date-fns';
import classNames from 'classnames';
import {
  CalendarDaysIcon,
  MagnifyingGlassIcon,
  PencilIcon,
  XMarkIcon,
} from '@heroicons/react/24/outline'

import ImagesAPI, { stockImageSrc, images } from '../../apis/images';

import Spinner from '../../components/Spinner';
import { datesRenderer } from '../../utils/dates';
import TripsSyncAPI from '../../apis/tripsSync';


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
    // setIsLoading(true);
    // ImagesAPI.search(query)
    // .then(res => {
    //   const images = _get(res, "data.images");
    //   setImageList(images);
    //   setIsLoading(false);
    // });
    setImageList(images);
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
    const css = {
      figure: "relative max-w-sm transition-all rounded-lg duration-300 mb-2 group",
      figureImg: "block rounded-lg max-w-full group-hover:grayscale",
      figureBtn: "text-white m-2 py-2 px-3 rounded-full bg-green-500 hover:bg-green-700",
      figureBtnCtn: "absolute group-hover:opacity-100 opacity-0 top-2 right-0",
      figureCaption: "absolute px-1 text-sm text-white rounded-b-lg bg-slate-800/50 w-full bottom-0",
    }

    if (isLoading) {
      return <Spinner />
    }

    return (
      <div className='columns-2 md:columns-3'>
        { imageList.map((image: any) => {
            return (
              <figure className={css.figure}>
                <a href="#">
                  <img key={image.id}
                    srcSet={ImagesAPI.makeSrcSet(image)}
                    src={ImagesAPI.makeSrc(image)}
                    className={css.figureImg}
                  />
                  <div className={css.figureBtnCtn}>
                    <a
                      className={css.figureBtn}
                      onClick={() => {props.onCoverImageSelect(image)}}
                      href="#"
                    >
                      Select
                    </a>
                  </div>
                </a>
                <figcaption className={css.figureCaption}>
                  <a
                    target="_blank"
                    href={ImagesAPI.makeUserReferURL(_get(image, "user.username"))}
                  >
                    @{_get(image, "user.username")}, Unsplash
                  </a>
                </figcaption>
              </figure>
            );
          })}
      </div>
    );
  }

  if (!props.isOpen) {
    return <></>;
  }

  return (
    <div className="relative z-10" aria-labelledby="modal-title" role="dialog" aria-modal="true">
      <div className="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity"></div>
      <div className="fixed inset-0 z-10 overflow-y-auto">
        <div className="flex min-h-full flex-col p-4 text-center sm:items-center sm:p-0">
          <div className="relative transform rounded-lg bg-white text-left shadow-xl transition-all w-11/12 sm:my-8 sm:w-full sm:max-w-2xl">
            <div className="bg-white px-4 pt-5 pb-4 sm:p-8 sm:pb-4 rounded-lg">
              <div className='flex justify-between mb-6'>
                <h2 className="text-lg sm:text-2xl font-bold leading-6 text-slate-900">
                  Change cover image
                </h2>
                <button type="button" onClick={props.onClose}>
                  <XMarkIcon className='h-6 w-6 text-slate-700' />
                </button>
              </div>
              <h2 className="text-sm font-medium text-indigo-500 sm:text-xl text-slate-700 mb-2 ml-1">
                Search the web
              </h2>
              <div className="flex mb-4 justify-between">
                <input
                  type="text"
                  className={classNames(
                    "bg-gray-50",
                    "block",
                    "border-gray-300",
                    "border",
                    "focus:border-blue-500",
                    "focus:ring-blue-500",
                    "min-w-0",
                    "p-2.5",
                    "rounded-lg",
                    "text-gray-900",
                    "text-sm",
                    "w-5/6",
                    "mr-2"
                  )}
                  value={query}
                  onChange={searchImageQueryOnChange}
                  onKeyDown={searchImageQueryOnEnter}
                  placeholder="destination, theme ..."
                />
                <button
                  type='button'
                  className='flex-1 inline-flex text-white bg-indigo-500 hover:bg-indigo-800 rounded-2xl p-2.5 text-center items-center justify-around'
                  onClick={searchImage}
                >
                  <MagnifyingGlassIcon className='h-5 w-5 stroke-2 stroke-white'/>
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

const TripMenuJumbo: FC<TripMenuJumboProps> = (props: TripMenuJumboProps) => {

  console.log(props.trip)

  // UI State
  const [isCoverImageModalOpen, setIsCoverImageModalOpen] = useState(false);
  const [tripName, setTripName] = useState(props.trip.name);

  // When props.trip changes, need to update the name
  useEffect(() => {
    setTripName(props.trip.name)
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
    console.log(ops)
    props.tripStateOnUpdate(ops)
  }

  // Renderers

  const renderDatesButton = () => {
    if (!_get(props.trip, "startDate")) {
      return;
    }

    const nullDate = parseJSON("0001-01-01T00:00:00Z");
    const startDate = parseISO(props.trip.startDate);
    const endDate = parseJSON(props.trip.endDate);

    if (isEqual(startDate, nullDate)) {
      return "";
    }

    return (
      <button type="button" className="font-medium text-md text-slate-500">
        <CalendarDaysIcon className='inline h-5 w-5 align-sub' />
        &nbsp;&nbsp;
        <span>{datesRenderer(startDate, endDate)}</span>
      </button>
    );
  }

  return (
    <>
      <div className='bg-yellow-200'>
        <div className="relative">
          <img
            srcSet={ImagesAPI.makeSrcSet(props.trip.coverImage)}
            src={ImagesAPI.makeSrc(props.trip.coverImage)}
            className="block sm:max-h-96 w-full"
          />
          <button
            type='button'
            className='absolute top-4 right-4 h-10 w-10 bg-gray-800/50 p-2 text-center rounded-full'
            onClick={() => { setIsCoverImageModalOpen(true) }}
          >
            <PencilIcon className='h-6 w-6 text-white' />
          </button>
        </div>
        <div className='h-16 relative -top-24'>
          <div className="bg-white rounded-lg shadow-xl p-5 mx-4 mb-4">
            <input
              type="text"
              value={tripName}
              onChange={tripNameOnChange}
              onBlur={tripNameOnBlur}
              className={classNames(
                "mb-12",
                "text-2xl",
                "sm:text-4xl",
                "font-bold",
                "text-slate-700",
                "w-full",
                "rounded-lg",
                "p-1",
                "border-0",
                "hover:bg-slate-300",
                "hover:border-0",
                "hover:bg-slate-100",
                "focus:ring-0",
              )}
            />
            <div className='flex justify-between'>
              {renderDatesButton()}
            </div>
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
