import React, { FC, useEffect, useState } from 'react';
import _get from "lodash/get";
import _sortBy from "lodash/sortBy";
import _isEmpty from "lodash/isEmpty";
import {
  SelectRangeEventHandler,
  DateRange
} from 'react-day-picker';
import {
  CalendarDaysIcon,
  PhoneIcon,
  TrashIcon,
  MapPinIcon,
} from '@heroicons/react/24/solid';
import {
  CurrencyDollarIcon,
  GlobeAltIcon
} from '@heroicons/react/24/outline';

import InputDatesPicker from '../InputDatesPicker';
import PlacePicturesCarousel from './PlacePicturesCarousel';
import { Trips } from '../../apis/types';
import { InputDatesPickerCss, LodgingCardCss } from '../../styles/global';
import {
  printTime,
  isEmptyDate,
  parseTimeFromZ,
  parseTripDate
} from '../../utils/dates';
import { capitaliseWords } from '../../utils/strings';
import { useMap } from '../../context/maps-context';

// TripLodgingCard

interface TripLodgingCardProps {
  lodging: Trips.Lodging
  onUpdate: any
  onDelete: any
}

const TripLodgingCard: FC<TripLodgingCardProps> = (props: TripLodgingCardProps) => {
  const place = _get(props.lodging, "place");
  const {dispatch} = useMap();


  // UI State
  const [isShowEdit, setIsShowEdit] = useState<Boolean>(false);
  const [checkinDates, setCheckinDates] = useState<DateRange>();
  const [priceAmount, setPriceAmount] = useState<Number>();
  const [updatedPaths, setUpdatedPaths] = useState({} as any);

  useEffect(() => {
    setCheckinDates({
      from: parseTripDate(props.lodging.checkinTime as (string|undefined)),
      to: parseTripDate(props.lodging.checkoutTime as (string|undefined)),
    });
    setPriceAmount(props.lodging.priceMetadata.amount);
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

  const editFormOnClick = (e: React.MouseEvent) => {
    e.stopPropagation();
  }

  const datesOnChange: SelectRangeEventHandler = (range: DateRange | undefined) => {
    setCheckinDates(range);
    setUpdatedPaths(Object.assign(updatedPaths, {"checkinTime": range?.from,}));
    setUpdatedPaths(Object.assign(updatedPaths, {"checkoutTime": range?.to,}));
  }

  const priceOnChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setPriceAmount(e.target.value ? Number(e.target.value) : undefined);
    setUpdatedPaths(Object.assign(updatedPaths, {
      "priceMetadata/amount": Number(e.target.value),
    }));
  }

  const deleteBtnOnClick = (e: React.MouseEvent) => {
    props.onDelete(props.lodging)
  }

  const placeOnClick = () => {
    // dispatch({type:"setSelectedPlace", value: props.lodging.place})
    // const event = new CustomEvent('marker_click', {
    //   bubbles: false,
    //   cancelable: false,
    //   detail: props.lodging.place,
    // });
    // document.getElementById("map")!.dispatchEvent(event)
  }

  // Renderers
  const renderPriceMetadata = () => {
    if (priceAmount === undefined || priceAmount === 0) {
      return null;
    }
    return (
      <p className={LodgingCardCss.PricePill}>
        $ {String(priceAmount)}
      </p>
    );
  }

  const renderNonEdit = () => {
    const dateFmt = "eee, MMM dd"
    return (
      <div>
        <button type='button'
          className={LodgingCardCss.AddrTxt}
          onClick={placeOnClick}
        >
          <MapPinIcon className='h-4 w-4 mr-1'/>
          {place.formatted_address}
        </button>
        { _isEmpty(place.website) ? null
          : <a
              className='flex items-center'
              href={place.website}
              target="_blank">
                <GlobeAltIcon className='h-4 w-4' />&nbsp;
                <span className={LodgingCardCss.WebsiteTxt}>Website</span>
              </a>
        }
        <p className={LodgingCardCss.PhoneTxt}>
          <PhoneIcon className='h-4 w-4' />&nbsp;
          {place.international_phone_number}
        </p>
        <p className={LodgingCardCss.DatesTxt}>
          <CalendarDaysIcon className='h-4 w-4' />&nbsp;
          {isEmptyDate(props.lodging.checkinTime) ? null
            : printTime(parseTimeFromZ(props.lodging.checkinTime as string), dateFmt)}
          {isEmptyDate(props.lodging.checkoutTime) ? null :
            " - " + printTime(parseTimeFromZ(props.lodging.checkoutTime as string), dateFmt)}
        </p>
        {renderPriceMetadata()}
        <PlacePicturesCarousel photos={props.lodging.place.photos} />
      </div>
    );
  }

  const renderEditForm = () => {
    return (
      <div
        className='mt-2'
        onClick={editFormOnClick}
      >
        <InputDatesPicker
          WrapperCss={"mb-2"}
          CtnCss={LodgingCardCss.DatesPickerCtn}
          onSelect={datesOnChange}
          dates={checkinDates}
        />
        <div className={LodgingCardCss.PriceInputCtn}>
          <span className={InputDatesPickerCss.Label}>
            <CurrencyDollarIcon className={InputDatesPickerCss.Icon} />
            &nbsp;Amount
          </span>
          <input
            type="number"
            value={priceAmount as any}
            onChange={priceOnChange}
            className={InputDatesPickerCss.Input}
          />
        </div>
        <button
          type='button'
          className={LodgingCardCss.DeleteBtn}
          onClick={deleteBtnOnClick}
        >
          <TrashIcon className='h-4 w-4' />
        </button>
      </div>);
  }

  // Renderers
  return (
    <div className={LodgingCardCss.Ctn} onClick={cardOnDoubleClick}>
      <h4 className={LodgingCardCss.Header}>
        {capitaliseWords(place.name)}
      </h4>
      {isShowEdit ? renderEditForm() : renderNonEdit()}
    </div>
  );
}

export default TripLodgingCard;
