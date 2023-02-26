import React, { FC } from 'react';
import _get from "lodash/get";
import _flatten from "lodash/flatten";
import _isEmpty from "lodash/isEmpty";
import _find from "lodash/find";
import {
  ClockIcon,
  StarIcon,
  MapPinIcon,
  PhoneIcon,
  GlobeAltIcon,
} from '@heroicons/react/24/solid'
import {
  MagnifyingGlassCircleIcon,
  XMarkIcon
} from '@heroicons/react/24/outline'

import GoogleIcon from '../../components/icons/GoogleIcon';
import {
  MapElementID,
  newZoomMarkerClick,
} from '../../lib/maps';
import { CommonCss, TripMapCss } from '../../assets/styles/global';

interface PlaceDetailsCardProps {
  placeDetails: any
  width: string
  onClose: () => void
}

const PlaceDetailsCard: FC<PlaceDetailsCardProps> = (props: PlaceDetailsCardProps) => {
  const { placeDetails } = props;

  const renderHeader = () => {
    return (
      <p className={TripMapCss.HeaderCtn}>
        <span className={TripMapCss.TitleCtn}>
          <button
            type="button"
            onClick={() => {
              const event = newZoomMarkerClick(placeDetails);
              document.getElementById(MapElementID)?.dispatchEvent(event)
            }}
          >
            <MagnifyingGlassCircleIcon className={CommonCss.LeftIcon} />
          </button>
          {placeDetails.name}
        </span>
        <button type="button" onClick={props.onClose}>
          <XMarkIcon className={CommonCss.Icon} />
        </button>
      </p>
    );
  }

  const renderSummary = () => {
    return (
      <p className={TripMapCss.SummaryTxt}>
        {_get(placeDetails, "editorial_summary.overview", "")}
      </p>
    );
  }

  const renderAddr = () => {
    return (
      <p className={TripMapCss.AddrTxt}>
        <MapPinIcon className={CommonCss.LeftIcon} />
        {placeDetails.formatted_address}
      </p>
    );
  }

  const renderRatings = () => {
    if (placeDetails.user_ratings_total === 0) {
      return null;
    }
    return (
      <p className={TripMapCss.RatingsStar}>
        <StarIcon className={CommonCss.LeftIcon} />
        {placeDetails.rating}&nbsp;&nbsp;
        <span className={TripMapCss.RatingsTxt}>
          ({placeDetails.user_ratings_total})
        </span>
        &nbsp;&nbsp;
        <GoogleIcon className={CommonCss.DropdownIcon} />
      </p>
    );
  }

  const renderOpeningHours = () => {
    const weekdayTexts = _get(placeDetails, "opening_hours.weekday_text", []);
    if (_isEmpty(weekdayTexts)) {
      return null;
    }
    return (
      <div>
        <p className={TripMapCss.OpeningHrsTxt}>
          <ClockIcon className={CommonCss.LeftIcon} />Opening hours
        </p>
        {weekdayTexts.map((txt: string, idx: number) =>
          (<p key={idx} className={TripMapCss.WeekdayTxt}>{txt}</p>)
        )}
      </div>
    );
  }

  const renderPhone = () => {
    return placeDetails.international_phone_number
      ?
      <a
        href={`tel:${placeDetails.international_phone_number.replace(/\s/, "-")}`}
        target="_blank"
        rel='noreferrer'
        className={TripMapCss.PhoneBtn}
      >
        <PhoneIcon className={TripMapCss.PhoneIcon} /> Call
      </a>
      : null
  }

  const renderWebsite = () => {
    return placeDetails.website
      ?
      <a
        href={placeDetails.website}
        target="_blank"
        rel='noreferrer'
        className={TripMapCss.PhoneBtn}
      >
        <GlobeAltIcon className={TripMapCss.PhoneIcon} /> Web
      </a>
      : null
  }

  const renderGmapBtn = () => {
    return (
      <a
        href={placeDetails.url}
        className={TripMapCss.GmapBtn}
      >
        <GoogleIcon className={CommonCss.LeftIcon} /> Google Maps
      </a>
    );
  }

  if (placeDetails === null) {
    return null;
  }

  return (
    <div
      className={TripMapCss.DetailsWrapper}
      style={{ width: props.width }}
    >
      <div className={TripMapCss.DetailsCard}>
        {renderHeader()}
        {renderSummary()}
        {renderAddr()}
        {renderRatings()}
        {renderOpeningHours()}
        <div className={TripMapCss.BtnCtn}>
          {renderPhone()}
          {renderWebsite()}
          {renderGmapBtn()}
        </div>
      </div>
    </div>
  );
}

export default PlaceDetailsCard;
