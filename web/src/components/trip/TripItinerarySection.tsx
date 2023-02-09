import React, {
  FC,
  useEffect,
  useState,
} from 'react';
import _get from "lodash/get";
import _sortBy from "lodash/sortBy";
import _isEmpty from "lodash/isEmpty";
import _find from "lodash/find";
import {
  ChevronDownIcon,
  ChevronUpIcon,
  MapPinIcon,
  TrashIcon
} from '@heroicons/react/24/solid'
import {
  EllipsisHorizontalCircleIcon,
  GlobeAltIcon,
} from '@heroicons/react/24/outline'


import PlacePicturesCarousel from './PlacePicturesCarousel';
import Dropdown from '../Dropdown';
import NotesEditor from '../NotesEditor';

import TripsSyncAPI from '../../apis/tripsSync';
import { Trips } from '../../apis/types';
import { useMap } from '../../context/maps-context';
import {
  TripItinerarySectionCss,
  TripItineraryListCss,
  TripItineraryCss
} from '../../styles/global';
import {
  areYMDEqual,
  isEmptyDate,
  parseTimeFromZ,
  printTime
} from '../../utils/dates'
import PlaneIcon from '../icons/PlaneIcon';
import HotelIcon from '../icons/HotelIcon';


// ItineraryContent
interface ItineraryContentProps {
  content: Trips.Content
  itineraryContentIdx: number
  itineraryContent: Trips.ItineraryContent
  tripStateOnUpdate: any
}

const ItineraryContent: FC<ItineraryContentProps> = (props: ItineraryContentProps) => {

  const { dispatch } = useMap();

  // Event Handlers - Title

  const deleteBtnOnClick = () => {
    const ops = [] as any;
    // ops.push(
    //   TripsSyncAPI.makeRemoveOp(
    //     `/contents/${props.contentListID}/contents/${props.contentIdx}`,
    //     "")
    // );
    props.tripStateOnUpdate(ops);
  }

  // Event Handlers - Places

  const placeOnClick = (e: React.MouseEvent) => {
    dispatch({type:"setSelectedPlace", value: props.content.place})
    const event = new CustomEvent('marker_click', {
      bubbles: false,
      cancelable: false,
      detail: props.content.place,
    });
    document.getElementById("map")!.dispatchEvent(event)
    return;
  }

  // Renderers
  const renderTitleInput = () => {
    return (
      <div className='flex justify-between'>
        <p className={TripItineraryCss.TitleInput}>
          {props.content.title}
        </p>
      </div>
    );
  }

  const renderPlace = () => {
    let placeNode;
    const addr = _get(props.content, "place.name", "");
    if (_isEmpty(addr)) {
      placeNode = null;
    } else {
      placeNode = (
        <button type='button' onClick={placeOnClick}>
          {addr}
        </button>
      );
    }
    return (
      <p className='text-slate-600 text-sm flex items-center mb-1 hover:text-indigo-500'>
        <MapPinIcon className='h-4 w-4 mr-1'/>
        {placeNode}
      </p>
    );
  }

  return (
    <div className={TripItineraryCss.Ctn}>
      {renderTitleInput()}
      {renderPlace()}
      <NotesEditor
        ctnCss='p-0 mb-2'
        base64Notes={props.content.notes}
        notesOnChange={() => {}}
        placeholder={"Notes..."}
        readOnly
      />
    </div>
  );
}


// TripItineraryList

interface TripItineraryListProps {
  trip: any
  itineraryList: Trips.ItineraryList
  tripStateOnUpdate: any
}

const TripItineraryList: FC<TripItineraryListProps> = (props: TripItineraryListProps) => {

  const [isHidden, setIsHidden] = useState<boolean>(false);

  // Renderers
  const renderHeader = () => {
    return (
      <div className='flex mb-2'>
        <button
          type="button"
          className={TripItinerarySectionCss.ToggleBtn}
          onClick={() => {setIsHidden(!isHidden)}}
        >
        {isHidden ? <ChevronUpIcon className='h-4 w-4' />
          : <ChevronDownIcon className='h-4 w-4'/>}
        </button>
        <p className='text-xl font-bold'>
          {printTime(parseTimeFromZ(props.itineraryList.date as string), "eeee, do MMMM") }
        </p>
      </div>

    );
  }

  const renderFlights = () => {
    const flights = Object.values(_get(props.trip, "flights", {}));
    const today = parseTimeFromZ(props.itineraryList.date as string);

    const render = (place: any, depart: boolean) => {
      return (
        <div className="flex items-center w-full p-3 space-x-4 text-gray-800 divide-x divide-gray-200 rounded-lg shadow">
          <span className='bg-green-200 p-2 rounded-full'><PlaneIcon className='w-4 h-4' /></span>
          <div className="flex-1 pl-4 text-sm font-normal">{place.name}</div>
          <span className='pl-2 font-semibold text-sm'>{depart ? "Depar": "Arrival"}</span>
        </div>
      );
    }

    const departs = [] as Array<Trips.Flight>;
    const arrivals = [] as Array<Trips.Flight>;

    flights.forEach((item: any) => {
      let flight = item as Trips.Flight;

      const onewayDepartFlightDt = flight.depart.departure.datetime;
      if (!isEmptyDate(onewayDepartFlightDt)
        && areYMDEqual(today, onewayDepartFlightDt)) {
        departs.push(item);
      }

      const onewayArrFlightDt = flight.depart.arrival.datetime;
      if (!isEmptyDate(onewayArrFlightDt)
        && areYMDEqual(today, onewayArrFlightDt)) {
        arrivals.push(item);
      }

      const returnDepartFlightDt = flight.return.departure.datetime;
      if (!isEmptyDate(returnDepartFlightDt)
        && areYMDEqual(today, returnDepartFlightDt)) {
        departs.push(item);
      }

      const returnArrFlightDt = flight.return.arrival.datetime
      if (!isEmptyDate(returnArrFlightDt)
        && areYMDEqual(today, returnArrFlightDt)) {
        arrivals.push(item);
      }
    });

    return (
      <div className='w-full mb-2'>
        {departs.map((item: any) => render(item.place, false))}
        {arrivals.map((item: any) => render(item.place, true))}
      </div>
    );
  }

  const renderLodgings = () => {
    const lodgings = Object.values(_get(props.trip, "lodgings", {}));
    const today = parseTimeFromZ(props.itineraryList.date as string);

    const render = (place: any, checkin: boolean) => {
      return (
        <div className="flex items-center w-full p-3 space-x-4 text-gray-800 divide-x divide-gray-200 rounded-lg shadow">
          <span className='bg-indigo-200 p-2 rounded-full'><HotelIcon className='w-4 h-4' /></span>
          <div className="flex-1 pl-4 text-sm font-normal">{place.name}</div>
          <span className='pl-2 font-semibold text-sm'>{checkin ? "Check in": "Check out"}</span>
        </div>
      );
    }

    const checkins = [] as Array<Trips.Lodging>;
    const checkouts = [] as Array<Trips.Lodging>;

    lodgings.forEach((item: any) => {
      let lod = item as Trips.Lodging;
      if (!_isEmpty(lod.checkinTime)) {
        const checkinTime = parseTimeFromZ(lod.checkinTime as string);
        if(areYMDEqual(today, checkinTime)) {
          checkins.push(item);
        }
        const checkoutTime = parseTimeFromZ(lod.checkoutTime as string);
        if(areYMDEqual(today, checkoutTime)) {
          checkouts.push(item);
        }
      }
    });

    return (
      <div className='w-full mb-2'>
        {checkouts.map((item: any) => render(item.place, false))}
        {checkins.map((item: any) => render(item.place, true))}
      </div>
    );
  }

  const renderContents = () => {
    const itinConents = _get(props.itineraryList, "contents", []);
    if (_isEmpty(itinConents)) {
      return (
        <p className='text-gray-500'>
          No activites added for today.
        </p>
      )
    }
    return (
      <div className='pl-6 mt-4'>
        <ol className='relative border-l border-gray-200'>
          {itinConents.map((itinCtn: Trips.ItineraryContent, idx: number) => {
            const contentList = _get(props.trip, `contents.${itinCtn.tripContentListId}`);
            const content = _find(
              contentList.contents,
              (ctn: Trips.Content) => ctn.id === itinCtn.tripContentId);

            return (
              <li
                key={idx}
                className="mb-10 ml-6"

              >
                <span className="absolute flex items-center justify-center w-6 h-6 bg-yellow-200 rounded-full -left-3 ring-8 ring-white font-bold text-gray-500 text-sm">
                  {idx + 1}
                </span>
                <ItineraryContent
                  content={content}
                  itineraryContentIdx={idx}
                  itineraryContent={itinCtn}
                  tripStateOnUpdate={props.tripStateOnUpdate}
                />
              </li>
            );
          })}
        </ol>
      </div>
    );
  }

  return (
    <div className={TripItineraryListCss.Ctn}>
      {renderHeader()}
      {isHidden ? null :
        <>
          {renderFlights()}
          {renderLodgings()}
          {renderContents()}
        </>
      }
    </div>
  );
}

// TripItinerarySection


interface TripItinerarySectionProps {
  trip: any
  tripStateOnUpdate: any
}

const TripItinerarySection: FC<TripItinerarySectionProps> = (props: TripItinerarySectionProps) => {

  return (
    <div className='p-5'>
      {
        _get(props.trip, "itinerary", []).map((l: any) => (
        <div key={l.id}>
          <TripItineraryList
            trip={props.trip}
            itineraryList={l}
            tripStateOnUpdate={props.tripStateOnUpdate}
          />
          <hr className={TripItinerarySectionCss.Hr} />
        </div>
      ))
    }
    </div>
  );

}

export default TripItinerarySection;
