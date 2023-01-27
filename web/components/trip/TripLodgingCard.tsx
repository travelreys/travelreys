import React, { FC, useEffect, useState } from 'react';
import _get from "lodash/get";
import _sortBy from "lodash/sortBy";
import _isEmpty from "lodash/isEmpty";
import { parseJSON } from 'date-fns';
import { SelectRangeEventHandler, DateRange } from 'react-day-picker';
import {
  CalendarDaysIcon,
  PhoneIcon,
  MapPinIcon,
  TrashIcon,
} from '@heroicons/react/24/solid';
import { CurrencyDollarIcon } from '@heroicons/react/24/outline';

import InputDatesPicker from '../InputDatesPicker';
import { printTime, isNullDate } from '../../utils/dates';
import { InputDatesPickerCss } from '../../styles/global';

import { Navigation, Pagination } from 'swiper';

import { Swiper, SwiperSlide } from 'swiper/react';

// Import Swiper styles
import 'swiper/css';
import 'swiper/css/navigation';
import 'swiper/css/pagination';


// TripLodgingPhotosCarosel

interface TripLodgingPhotosCaroselProps {
  photos: any
}

const TripLodgingPhotosCarosel: FC<TripLodgingPhotosCaroselProps> = (props: TripLodgingPhotosCaroselProps) => {
  console.log(props.photos)
  const renderLodgingImgs = () => {
    return props.photos.map((photo: any) => (
      <SwiperSlide key={photo.photo_reference}>
        <img src={gMapsPlaceImageSrcURL(photo.photo_reference)} className="rounded" />
      </SwiperSlide>
    ));
  }

  if (props.photos.length === 0) {
    return (<></>);
  }

  return (
    <Swiper
      modules={[Navigation, Pagination]}

      slidesPerView={1}
      navigation
      pagination={{ clickable: true }}
      scrollbar={{ draggable: true }}
      onSwiper={(swiper) => console.log(swiper)}
      onSlideChange={() => console.log('slide change')}
    >
      {renderLodgingImgs()}
    </Swiper>
  );
}







// TripLogdingCardProps

interface TripLodgingCardProps {
  lodging: any
  onUpdate: any
  onDelete: any
}

const gMapsPlaceImageSrcURL = (ref: string) => {
  return [
    "https://maps.googleapis.com/maps/api/place/photo",
    "?maxwidth=800",
    `&photo_reference=${ref}`,
    `&key=${"AIzaSyBgNwirAT6TSS208emcC0Lbgex6i3EwhR0"}`,
  ].join("");
}


const TripLodgingCard: FC<TripLodgingCardProps> = (props: TripLodgingCardProps) => {
  const place = _get(props.lodging, "place");
  const checkinTime = parseJSON(_get(props.lodging, "checkinTime"));
  const checkoutTime = parseJSON(_get(props.lodging, "checkoutTime"));

  // UI State
  const [isShowEdit, setIsShowEdit] = useState(false);
  const [checkinDates, setCheckinDates] = useState({
    from: checkinTime,
    to: checkoutTime,
  } as any);
  const [priceMetadata, setPriceMetadata] = useState(props.lodging.priceMetadata);
  const [updatedPaths, setUpdatedPaths] = useState({} as any);

  useEffect(() => {
    setCheckinDates({
      from: parseJSON(_get(props.lodging, "checkinTime")),
      to: parseJSON(_get(props.lodging, "checkoutTime")),
    });
    setPriceMetadata(props.lodging.priceMetadata);
  }, [props.lodging])

  // Event Handlers
  const cardOnDoubleClick = (event: React.MouseEvent) => {
    if (event.detail <= 1) {
      return;
    }
    if (!isShowEdit) {
      setIsShowEdit(true);
      return;
    }

    if (_isEmpty(updatedPaths)) {
      setIsShowEdit(false);
      return;
    }

    props.onUpdate(props.lodging, updatedPaths);
    setIsShowEdit(false);
  }

  const datesOnChange: SelectRangeEventHandler = (range: DateRange | undefined) => {
    console.log(range)
    setCheckinDates(range);
    setUpdatedPaths(Object.assign(updatedPaths, {
      "checkinTime": range?.from,
    }));
    setUpdatedPaths(Object.assign(updatedPaths, {
      "checkoutTime": range?.to,
    }));
  }

  // Renderers


  const renderPriceMetadata = () => {
    const amount = _get(props.lodging, "priceMetadata.amount");
    if (amount === 0 || amount === undefined) {
      return null;
    }
    return (
      <p className='mb-2'>
        <span className="bg-blue-100 text-blue-800 text-xs font-medium px-2.5 py-0.5 rounded-full">
          $ {props.lodging.priceMetadata.amount}
        </span>
      </p>
    );
  }

  const renderNonEdit = () => {
    return (
      <div>
        <p className='text-slate-600 text-sm flex items-center mb-1'>
          <MapPinIcon className='h-4 w-4' />&nbsp;
          {place.formatted_address}
        </p>
        <p className='text-slate-600 text-sm flex items-center mb-1'>
          <PhoneIcon className='h-4 w-4' />&nbsp;
          {place.international_phone_number}
        </p>
        <p className='text-slate-600 text-sm flex items-center mb-1'>
          <CalendarDaysIcon className='h-4 w-4' />&nbsp;
          {isNullDate(checkinTime) ? null : printTime(checkinTime, "eee, MMM dd")}
          {isNullDate(checkinTime) ? null : " - " + printTime(checkoutTime, "eee, MMM dd")}
        </p>
        {renderPriceMetadata()}
        <TripLodgingPhotosCarosel photos={props.lodging.place.photos} />
      </div>
    );
  }

  const renderEditForm = () => {
    return (
      <div
        className='mt-2'
        onClick={(e) => { e.stopPropagation() }}
      >
        <InputDatesPicker
          WrapperCss={"mb-2"}
          CtnCss={"flex w-full rounded"}
          onSelect={datesOnChange}
          dates={checkinDates}
        />
        <div className="flex w-full rounded mb-2">
          <span className={InputDatesPickerCss.Label}>
            <CurrencyDollarIcon className={InputDatesPickerCss.Icon} />
            &nbsp;Amount
          </span>
          <input
            type="number"
            value={priceMetadata.amount || undefined}
            onChange={(e) => {
              setPriceMetadata({ amount: e.target.value });
              setUpdatedPaths(Object.assign(updatedPaths, {
                "priceMetadata/amount": Number(e.target.value),
              }));
            }}
            className={InputDatesPickerCss.Input}
          />
        </div>
        <button
          type='button'
          className='bg-red-500 py-2 px-4 rounded-lg text-white'
          onClick={() => { props.onDelete(props.lodging) }}
        >
          <TrashIcon className='h-4 w-4' />
        </button>
      </div>);
  }

  // Renderers
  return (
    <div
      className='p-4 bg-slate-50 rounded-lg shadow-md mb-4 cursor-pointer'
      onClick={cardOnDoubleClick}
    >
      <div className='flex justify-between items-center'>
        <h4 className='font-bold mb-1'>{place.name}</h4>
      </div>
      {isShowEdit ? renderEditForm() : renderNonEdit()}
    </div>
  );

}

export default TripLodgingCard;
