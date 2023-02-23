import React, { FC, useEffect, useState } from 'react';
import _get from "lodash/get";
import _sortBy from "lodash/sortBy";
import _isEmpty from "lodash/isEmpty";
import {
  SelectRangeEventHandler,
  DateRange,
} from 'react-day-picker';
import {
  CalendarDaysIcon,
  PhoneIcon,
  TrashIcon,
  MapPinIcon,
} from '@heroicons/react/24/solid';
import {
  CurrencyDollarIcon,
  EllipsisHorizontalCircleIcon,
  GlobeAltIcon
} from '@heroicons/react/24/outline';

import Dropdown from '../../components/common/Dropdown';
import InputDatesPicker from '../../components/common/InputDatesPicker';
import PlacePicturesCarousel from './PlacePicturesCarousel';
import { Trips } from '../../lib/trips';
import { CommonCss, InputDatesPickerCss, LodgingCardCss } from '../../assets/styles/global';
import {
  printFmt,
  isEmptyDate,
  parseISO,
  parseTripDate
} from '../../lib/dates';
import { capitaliseWords } from '../../lib/strings';
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
  const [isUpdatingDates, setIsUpdatingDates] = useState<Boolean>(false);
  const [isUpdatingPrice, setIsUpdatingPrice] = useState<Boolean>(false);
  const [checkinDates, setCheckinDates] = useState<DateRange>();
  const [priceAmount, setPriceAmount] = useState<Number>();


  useEffect(() => {
    setCheckinDates({
      from: parseTripDate(props.lodging.checkinTime as (string|undefined)),
      to: parseTripDate(props.lodging.checkoutTime as (string|undefined)),
    });
    setPriceAmount(props.lodging.price.amount);
  }, [props.lodging])

  // Event Handlers - Dates

  const datesOnClick = (e: any) => {
    if (e.detail <= 1) {
      return;
    }
    setIsUpdatingDates(true)
  }

  const datesOnChange: SelectRangeEventHandler = (range?: DateRange) => {
    setCheckinDates(range);
  }

  const datesOnBlur = (e: any) => {
    if (!e.currentTarget.contains(e.relatedTarget) && isUpdatingDates) {
      setIsUpdatingDates(false);
      props.onUpdate(props.lodging, {
        "checkinTime": checkinDates?.from,
        "checkoutTime": checkinDates?.to
      })
    }
  }

  // Event Handlers - Price
  const priceOnClick = (e:  any) => {
    if (e.detail <= 1) {
      return;
    }
    setIsUpdatingPrice(true)
  }

  const priceOnChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setPriceAmount(e.target.value ? Number(e.target.value) : undefined);
  }

  const priceOnBlur = () => {
    props.onUpdate(props.lodging, {
      "price/amount": priceAmount,
    });
    setIsUpdatingPrice(false);
  }

  // Event Handlers - Delete

  const deleteBtnOnClick = (e: React.MouseEvent) => {
    props.onDelete(props.lodging)
  }

  // Event Handlers - Place

  const placeOnClick = () => {
    dispatch({type:"setSelectedPlace", value: props.lodging.place})
    const event = new CustomEvent('marker_click', {
      bubbles: false,
      cancelable: false,
      detail: props.lodging.place,
    });
    document.getElementById("map")!.dispatchEvent(event)
  }

  // Renderers
  const renderHeaders = () => {
    const opts = [
      <button
          type='button'
          className={CommonCss.DeleteBtn}
          onClick={deleteBtnOnClick}
        >
        <TrashIcon className='h-4 w-4 mr-2' />Delete
      </button>
    ];
    const menu = <EllipsisHorizontalCircleIcon className='h-4 w-4' />;
    return (
      <div className='flex justify-between'>
        <h4 className={LodgingCardCss.Header}>
          {capitaliseWords(place.name)}
        </h4>
        <Dropdown menu={menu} opts={opts} />
      </div>
    );
  }

  const renderWebsite = () => {
    return _isEmpty(place.website) ? null :
      <a
        className='flex items-center'
        href={place.website}
        target="_blank"
      >
        <GlobeAltIcon className='h-4 w-4' />&nbsp;
        <span className={LodgingCardCss.WebsiteTxt}>Website</span>
      </a>
  }

  const renderPrice = () => {
    if (isUpdatingPrice) {
      return (
        <div className={LodgingCardCss.PriceInputCtn}>
          <span className={InputDatesPickerCss.Label}>
            <CurrencyDollarIcon className={InputDatesPickerCss.Icon} />
            &nbsp;Amount
          </span>
          <input
            type="number"
            autoFocus
            value={priceAmount as any}
            onChange={priceOnChange}
            onBlur={priceOnBlur}
            className={InputDatesPickerCss.Input}
          />
        </div>
      );
    }
    return (
      <p className={LodgingCardCss.PricePill} onClick={priceOnClick}>
        $ {priceAmount ? String(priceAmount): "Add cost"}
      </p>
    );
  }

  const renderDates = () => {
    if (isUpdatingDates) {
      return (
        <div onBlur={datesOnBlur}>
          <InputDatesPicker
            WrapperCss={"mb-2"}
            CtnCss={LodgingCardCss.DatesPickerCtn}
            onSelect={datesOnChange}
            dates={checkinDates}
          />
        </div>
      );
    }

    const dateFmt = "eee, MMM dd";
    return (
      <p
        className={LodgingCardCss.DatesTxt}
        onClick={datesOnClick}
      >
        <CalendarDaysIcon className={CommonCss.Icon} />&nbsp;
        {isEmptyDate(props.lodging.checkinTime) ? null
          : printFmt(parseISO(props.lodging.checkinTime as string), dateFmt)}
        {isEmptyDate(props.lodging.checkoutTime) ? null :
          " - " + printFmt(parseISO(props.lodging.checkoutTime as string), dateFmt)}
      </p>
    );
  }

  return (
    <div className={LodgingCardCss.Ctn}>
      {renderHeaders()}
      <button type='button'
        className={LodgingCardCss.AddrTxt}
        onClick={placeOnClick}
      >
        <MapPinIcon className='h-4 w-4 mr-1'/>
        {place.formatted_address}
      </button>
      {renderWebsite()}
      <p className={LodgingCardCss.PhoneTxt}>
        <PhoneIcon className='h-4 w-4' />&nbsp;
        {place.international_phone_number}
      </p>
      {renderDates()}
      {renderPrice()}
      <PlacePicturesCarousel photos={props.lodging.place.photos} />
    </div>
  );
}

export default TripLodgingCard;
