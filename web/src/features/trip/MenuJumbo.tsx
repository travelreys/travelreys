import React, { FC, useEffect, useState } from 'react';
import _get from "lodash/get";
import { DateRange, SelectRangeEventHandler } from 'react-day-picker';
import { CalendarDaysIcon, PencilIcon } from '@heroicons/react/24/outline'

import Avatar from '../../components/common/Avatar';
import DatesPicker from '../../components/common/DatesPicker';
import CoverImageModal from './CoverImageModal';

import ImagesAPI from '../../apis/images';
import {
  fmt,
  nullDate,
  parseTripDate
} from '../../lib/dates';
import { makeRepOp } from '../../lib/jsonpatch';
import { Member, userFromMember } from '../../lib/trips';
import { LabelUserGoogleImage, User } from '../../lib/auth';
import { CommonCss } from '../../assets/styles/global';



interface MenuJumboProps {
  trip: any
  tripMembers: {[key: string]: User}
  onlineMembers: Array<Member>
  tripOnUpdate: any
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
    props.tripOnUpdate([makeRepOp("/name", tripName)])
  }

  // Event Handlers - Cover Image
  const coverImageOnSelect = (image: any) => {
    props.tripOnUpdate([makeRepOp("/coverImage", image)]);
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
      ops.push(makeRepOp("/startDate", from));
      ops.push(makeRepOp("/endDate", to));
      props.tripOnUpdate(ops);
      setIsCalendarOpen(false);
      return;
    }
  }

  // Renderers
  const css = {
      datesBtn: "font-medium text-md text-slate-500",
      datesBtnIcon: "inline h-5 w-5 align-sub mr-2",
      coverImage: "block sm:max-h-96 w-full",
      imageEditIconCtn: "absolute top-4 right-4 h-10 w-10 bg-gray-800/70 p-2 text-center rounded-full",
      imageEditIcon: "h-6 w-6 text-white",
      nameInputCtn: "h-16 relative -top-24",
      nameInputWrapper: "bg-white rounded-lg shadow p-5 mx-4 mb-4",
      namInputHeader: "flex items-center mb-12",
      nameInput: "text-2xl sm:text-4xl font-bold text-slate-700 w-full rounded-lg p-1 border-0 hover:bg-slate-300 hover:border-0 hover:bg-slate-100 focus:ring-0",
      onlineBtn: "flex items-center justify-center w-8 h-8 text-xs font-medium text-white bg-gray-700 border-2 border-white rounded-full hover:bg-gray-600",
  }


  const renderCoverImage = () => {
    return (
      <div className="relative">
        <img
          srcSet={ImagesAPI.makeSrcSet(props.trip.coverImage)}
          src={ImagesAPI.makeSrc(props.trip.coverImage)}
          className={css.coverImage}
        />
        <button
          type='button'
          className={css.imageEditIconCtn}
          onClick={() => { setIsCoverImageModalOpen(true) }}
        >
          <PencilIcon className={css.imageEditIcon} />
        </button>
      </div>
    );
  }

  const renderDatesButton = () => {
    const dateFmt = "MMM d, yy"
    return (
      <div className='flex-1' onBlur={tripDatesOnBlur}>
        <button
          type="button"
          className={css.datesBtn}
          onClick={() => { setIsCalendarOpen(true) }}
        >
          <CalendarDaysIcon className={css.datesBtnIcon} />
          {startDt ? <span>{fmt(startDt, dateFmt)}</span> : null }
          {endDt ? <span>&nbsp;-&nbsp;{fmt(endDt, dateFmt)}</span> : null}
        </button>
        <DatesPicker
          onSelect={tripDatesOnChange}
          isOpen={isCalendarOpen}
          dates={{from: startDt, to: endDt}}
        />
      </div>
    );
  }

  const renderOnlineMembers = () => {
    const imgs = [];
    props.onlineMembers.slice(0, 5).forEach((om: Member) => {
      const usr = userFromMember(om, props.tripMembers);
      imgs.push(
        <div className={CommonCss.IconLarge}>
          <Avatar
            key={om.id}
            placement="top"
            imgurl={_get(usr, `labels.${LabelUserGoogleImage}`)}
            name={_get(usr, "name", "")}
          />
        </div>
      );
    });

    if (props.onlineMembers.length > 5) {
      imgs.push(
        <button type='button' className={css.onlineBtn}>
          {props.onlineMembers.length - 5}
        </button>)
    }
    return (
      <div className="flex -space-x-3">
        {imgs}
      </div>
    );
  }

  const renderTripNameInput = () => {
    return (
      <div className={css.nameInputCtn}>
        <div className={css.nameInputWrapper}>
          <div className={css.namInputHeader}>
            <input
              type="text"
              value={tripName}
              onChange={(e) => setTripName(e.target.value)}
              onBlur={tripNameOnBlur}
              className={css.nameInput}
            />
          </div>
          <div className='flex items-center'>
            {renderDatesButton()}
            {renderOnlineMembers()}
          </div>
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
