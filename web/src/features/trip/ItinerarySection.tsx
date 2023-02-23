import React, {
  FC,
  useCallback,
  useEffect,
  useState,
} from 'react';
import _get from "lodash/get";
import _sortBy from "lodash/sortBy";
import _isEmpty from "lodash/isEmpty";
import _find from "lodash/find";
import { DndProvider, useDrag, useDrop } from 'react-dnd';
import { HTML5Backend } from 'react-dnd-html5-backend';
import {
  ArrowLongRightIcon,
  CurrencyDollarIcon,
  MapPinIcon,
  SwatchIcon,
} from '@heroicons/react/24/solid';
import {
  EllipsisHorizontalCircleIcon,
} from '@heroicons/react/24/outline'

import {
  MapElementID,
  newEventMarkerClick
} from '../maps/common';

import ColorIconModal from './ColorIconModal';
import Dropdown from '../../components/common/Dropdown';
import HotelIcon from '../../components/icons/HotelIcon';
import NotesEditor from '../../components/common/NotesEditor';
import PlaneIcon from '../../components/icons/PlaneIcon';
import ToggleChevron from '../../components/common/ToggleChevron';

import TripsSyncAPI from '../../apis/tripsSync';
import {
  ContentColorOpts,
  ContentIconOpts,
  DefaultContentColor,
  LabelContentListColor,
  LabelContentListIcon,
  PriceAmountJSONPath,
  LabelContentListColorJSONPath,
  LabelContentListIconJSONPath,
  Trips
} from '../../lib/trips';
import {
  areYMDEqual,
  isEmptyDate,
  parseFlightDateZ,
  parseISO,
  printFmt
} from '../../lib/dates'
import {
  flightArrivalTime,
  flightDepartureTime,
  FlightItineraryTypeRoundtrip,
  Flights
} from '../../lib/flights';

import {
  ActionSetSelectedPlace,
  useMap
} from '../../context/maps-context';
import {
  InputDatesPickerCss,
  TripItinerarySectionCss,
  TripItineraryListCss,
  TripItineraryCss,
  CommonCss,
  TripLogisticsCss,
} from '../../assets/styles/global';
import { makeAddOp, makeRemoveOp, makeReplaceOp } from '../../lib/tripsSync';


const ItineraryDateFmt = "eeee, do MMMM"
const ItineraryContentCard = "ItineraryContentCard";


// ItineraryContent
interface ItineraryContentProps {
  content: Trips.Content
  itineraryListIdx: number
  itineraryList: Trips.ItineraryList
  itineraryContentIdx: number
  itineraryContent: Trips.ItineraryContent
  tripStateOnUpdate: any

  moveCard: (id: string, to: number) => void
  findCard: (id: string) => { index: number }
  dropCard: (id: string) => void
}

const ItineraryContent: FC<ItineraryContentProps> = (props: ItineraryContentProps) => {

  const [isUpdatingPrice, setIsUpdatingPrice] = useState<Boolean>(false);
  const [priceAmount, setPriceAmount] = useState<Number>();
  const origCardIdx = props.findCard(props.itineraryContent.id).index;

  const { dispatch } = useMap();

  useEffect(() => {
    setPriceAmount(props.itineraryContent.price.amount);
  }, [props.itineraryContent])

  // Event Handles - DnD

  const [{ }, drag] = useDrag(
    () => ({
      type: ItineraryContentCard,
      item: { id: props.itineraryContent.id, origCardIdx },
      collect: (monitor) => ({
        isDragging: monitor.isDragging(),
      }),
      end: (item, monitor) => {
        const { id: droppedId, origCardIdx } = item
        if (monitor.didDrop()) {
          props.dropCard(droppedId);
        } else {
          props.moveCard(droppedId, origCardIdx)
        }
      },
    }),
    [props.itineraryContent.id, origCardIdx, props.moveCard],
  )

  const [, drop] = useDrop(
    () => ({
      accept: ItineraryContentCard,
      hover({ id: draggedId }: any) {
        if (draggedId !== props.itineraryContent.id) {
          const { index: overIndex } = props.findCard(props.itineraryContent.id)
          props.moveCard(draggedId, overIndex)
        }
      },
    }),
    [props.findCard, props.moveCard],
  )

  // Event Handlers - Places

  const placeOnClick = (e: React.MouseEvent) => {
    dispatch({ type: ActionSetSelectedPlace, value: props.content.place })
    const event = newEventMarkerClick(props.content.place);
    document.getElementById(MapElementID)?.dispatchEvent(event)
    return;
  }

  // Event Handlers - Price
  const priceOnClick = (e: any) => {
    if (e.detail <= 1) {
      return;
    }
    setIsUpdatingPrice(true)
  }

  const priceOnChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setPriceAmount(e.target.value ? Number(e.target.value) : undefined);
  }

  const priceOnBlur = () => {
    const { itineraryListIdx, itineraryContentIdx } = props;
    props.tripStateOnUpdate([
      makeReplaceOp(
        `/itinerary/${itineraryListIdx}/contents/${itineraryContentIdx}/${PriceAmountJSONPath}`,
        priceAmount,
      )
    ]);
    setIsUpdatingPrice(false);
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
        <MapPinIcon className='h-4 w-4 mr-1' />
        {placeNode}
      </p>
    );
  }

  const renderPricePill = () => {
    if (isUpdatingPrice) {
      return (
        <div className={TripItineraryCss.PriceInputCtn}>
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
      <p className={TripItineraryCss.PricePill} onClick={priceOnClick}>
        $ {priceAmount ? String(priceAmount) : "Add cost"}
      </p>
    );
  }

  const renderDirectionsDropdown = () => {
    return null;
  }

  return (
    <div>
      <div
        className={TripItineraryCss.Ctn}
        ref={(node) => drag(drop(node))}
      >
        {renderTitleInput()}
        {renderPlace()}
        <NotesEditor
          ctnCss='p-0 mb-2'
          base64Notes={props.content.notes}
          notesOnChange={() => { }}
          placeholder={"Notes..."}
          readOnly
        />
        {renderPricePill()}
      </div>
      {renderDirectionsDropdown()}
    </div>
  );
}

// TripItineraryList

interface TripItineraryListProps {
  trip: any
  itineraryListIdx: number
  itineraryList: Trips.ItineraryList
  tripStateOnUpdate: any

  onUpdateColorIcon: (itinListIdx: number, color?: string, icon?: string) => void
}

const TripItineraryList: FC<TripItineraryListProps> = (props: TripItineraryListProps) => {

  const [isHidden, setIsHidden] = useState<boolean>(false);
  const [itinContents, setItinContents] = useState(props.itineraryList.contents);
  const [isColorIconModalOpen, setIsColorIconModalOpen] = useState<boolean>(false);

  useEffect(() => {
    setItinContents(props.itineraryList.contents)
  }, [props.itineraryList])


  // Event Handlers

  const colorIconOnSubmit = (color?: string, icon?: string) => {
    props.onUpdateColorIcon(props.itineraryListIdx, color, icon)
  }

  const updateItinContents = (newItinContents: Array<Trips.ItineraryContent>) => {
    props.tripStateOnUpdate([
      makeReplaceOp(`/itinerary/${props.itineraryListIdx}/contents`, newItinContents),
    ]);
  }

  // DnD Helpers

  const findCard = useCallback((id: string) => {
    const itinContent = _find(itinContents, (cont: Trips.ItineraryContent) => cont.id === id);
    return { itinContent, index: itinContents.indexOf(itinContent!) }
  }, [itinContents]);

  const moveCard = useCallback((id: string, atIndex: number) => {
    const { itinContent, index } = findCard(id);
    itinContents.splice(index, 1);
    itinContents.splice(atIndex, 0, itinContent!);
    const newItinContents = itinContents.map((x) => x);
    setItinContents(newItinContents);
  }, [findCard, itinContents, setItinContents])

  const dropCard = useCallback((id: string) => {
    updateItinContents(itinContents);
  }, [findCard, itinContents, setItinContents])

  const [, drop] = useDrop(() => ({ accept: ItineraryContentCard }))

  // Renderers
  const renderSettingsDropdown = () => {
    const opts = [
      <button
        type='button'
        className={CommonCss.DropdownBtn}
        onClick={() => setIsColorIconModalOpen(true)}
      >
        <SwatchIcon className={CommonCss.LeftIcon} />
        Change Color & Icon
      </button>,
    ];
    const menu = (
      <EllipsisHorizontalCircleIcon
        className={CommonCss.DropdownIcon} />
    );
    return <Dropdown menu={menu} opts={opts} />
  }

  const renderHeader = () => {
    return (
      <div className='flex mb-2 justify-between'>
        <div className='flex flex-1'>
          <ToggleChevron
            isHidden={isHidden}
            onClick={() => { setIsHidden(!isHidden) }}
          />
          <p className='text-xl font-bold'>
            {printFmt(parseISO(props.itineraryList.date as string), ItineraryDateFmt)}
          </p>
        </div>
        {renderSettingsDropdown()}
      </div>

    );
  }

  const renderFlights = () => {
    const flights = Object.values(_get(props.trip, "flights", {}));
    const today = parseISO(props.itineraryList.date as string);

    const render = (flight: Flights.Flight) => {
      const departTime = flightDepartureTime(flight) as string;
      const arrTime = flightArrivalTime(flight) as string;
      const timeFmt = "hh:mm aa";

      return (
        <div className="flex items-center w-full p-3 space-x-4 text-gray-800 divide-x divide-gray-200 rounded-lg shadow">
          <span className='bg-green-200 p-2 rounded-full'>
            <PlaneIcon className='w-4 h-4' />
          </span>
          <div className="flex pl-4">
            <span>
              <p className={TripLogisticsCss.FlightTransitTime}>
                {printFmt(parseFlightDateZ(departTime), timeFmt)}
              </p>
              <p className={TripLogisticsCss.FlightTransitAirportCode}>
                {flight.departure.airport.code}
              </p>
            </span>
            <ArrowLongRightIcon
              className={TripLogisticsCss.FlightTransitLongArrow}
            />
            <span className='mb-1'>
              <p className={TripLogisticsCss.FlightTransitTime}>
                {printFmt(parseFlightDateZ(arrTime), timeFmt)}
              </p>
              <p className={TripLogisticsCss.FlightTransitAirportCode}>
                {flight.arrival.airport.code}
              </p>
            </span>
          </div>
        </div>
      );
    }

    const departs = [] as Array<Flights.Flight>;

    flights.forEach((item: any) => {
      let flight = item as Trips.Flight;

      const onewayDepartFlightDt = parseFlightDateZ(
        flightDepartureTime(flight.depart) as any);
      if (!isEmptyDate(onewayDepartFlightDt)
        && areYMDEqual(today, onewayDepartFlightDt)) {
        departs.push(flight.depart);
      }

      if (flight.itineraryType === FlightItineraryTypeRoundtrip) {
        const returnDepartFlightDt = parseFlightDateZ(
          flightDepartureTime(flight.return) as any);
        if (!isEmptyDate(returnDepartFlightDt)
          && areYMDEqual(today, returnDepartFlightDt)) {
          departs.push(flight.return);
        }
      }
    });

    return (
      <div className='w-full mb-2'>
        {departs.map((flight: Flights.Flight) => render(flight))}
      </div>
    );
  }

  const renderLodgings = () => {
    const lodgings = Object.values(_get(props.trip, "lodgings", {}));
    const today = parseISO(props.itineraryList.date as string);

    const render = (idx: number, place: any, checkin: boolean) => {
      return (
        <div key={idx} className={TripItineraryListCss.LodgingWrapper}>
          <span className={TripItineraryListCss.LodgingIconWrapper}>
            <HotelIcon className={CommonCss.Icon} />
          </span>
          <div className={TripItineraryListCss.LodgingName}>{place.name}</div>
          <span className={TripItineraryListCss.LodgingStatus}>
            {checkin ? "Check in" : "Check out"}
          </span>
        </div>
      );
    }

    const checkins = [] as Array<Trips.Lodging>;
    const checkouts = [] as Array<Trips.Lodging>;

    lodgings.forEach((item: any) => {
      let lod = item as Trips.Lodging;
      if (!_isEmpty(lod.checkinTime)) {
        const checkinTime = parseISO(lod.checkinTime as string);
        if (areYMDEqual(today, checkinTime)) {
          checkins.push(item);
        }
        const checkoutTime = parseISO(lod.checkoutTime as string);
        if (areYMDEqual(today, checkoutTime)) {
          checkouts.push(item);
        }
      }
    });

    return (
      <div className={TripItineraryListCss.LodgingCtn}>
        {checkouts.map((item: any, idx: number) => render(idx, item.place, false))}
        {checkins.map((item: any, idx: number) => render(idx, item.place, true))}
      </div>
    );
  }

  const renderContents = () => {
    if (_isEmpty(itinContents)) {
      return (
        <p className='text-gray-500'>
          No activites added for today.
        </p>
      );
    }
    const listItems = itinContents.map((itinCtn: Trips.ItineraryContent, idx: number) => {
      const content = _find(
        _get(props.trip, `contents.${itinCtn.tripContentListId}.contents`, []),
        (ctn: Trips.Content) => ctn.id === itinCtn.tripContentId
      );
      const color = _get(props.itineraryList, `labels.${LabelContentListColor}`, DefaultContentColor);
      return (
        <li key={idx} className={TripItineraryListCss.ItinItem}>
          <span
            style={{backgroundColor: color}}
            className={TripItineraryListCss.ItinContentIcon}
          >
            {idx + 1}
          </span>
          <ItineraryContent
            content={content}
            itineraryListIdx={props.itineraryListIdx}
            itineraryList={props.itineraryList}
            itineraryContentIdx={idx}
            itineraryContent={itinCtn}
            tripStateOnUpdate={props.tripStateOnUpdate}
            findCard={findCard}
            moveCard={moveCard}
            dropCard={dropCard}
          />
        </li>
      );
    })

    return (
      <div className={TripItineraryListCss.ContentsCtn}>
        <ol
          ref={drop}
          className={TripItineraryListCss.ContentsWrapper}
        >
          {listItems}
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
      <ColorIconModal
        isOpen={isColorIconModalOpen}
        colors={ContentColorOpts}
        icons={Object.keys(ContentIconOpts)}
        onClose={() => setIsColorIconModalOpen(false)}
        onSubmit={colorIconOnSubmit}
      />
    </div>
  );
}

//////////////////////////
// TripItinerarySection //
//////////////////////////

interface ItinerarySectionProps {
  trip: any
  tripStateOnUpdate: any
}

const ItinerarySection: FC<ItinerarySectionProps> = (props: ItinerarySectionProps) => {

  const updateItineraryListColorIcon = (itinListIdx: number, color?: string, icon?: string) => {
    const ctntList = _get(props.trip, `itinerary.${itinListIdx}`);
    const colorLabel = _get(ctntList, `labels.${LabelContentListColor}`);
    const iconLabel = _get(ctntList, `labels.${LabelContentListIcon}`);

    const ops = [];
    if (_isEmpty(color) && !_isEmpty(colorLabel)) {
      ops.push(makeRemoveOp(`/itinerary/${itinListIdx}/${LabelContentListColorJSONPath}`, ""));
    }
    if (!_isEmpty(color)) {
      if (_isEmpty(colorLabel)) {
        ops.push(makeAddOp(`/itinerary/${itinListIdx}/${LabelContentListColorJSONPath}`, color));
      } else {
        ops.push(makeReplaceOp(`/itinerary/${itinListIdx}/${LabelContentListColorJSONPath}`, color));
      }
    }

    if (_isEmpty(icon) && !_isEmpty(iconLabel)) {
      ops.push(makeRemoveOp(`/itinerary/${itinListIdx}/${LabelContentListIconJSONPath}`, ""));
    }
    if (!_isEmpty(icon)) {
      if (_isEmpty(colorLabel)) {
        ops.push(makeAddOp(`/itinerary/${itinListIdx}/${LabelContentListIconJSONPath}`, icon));
      } else {
        ops.push(makeReplaceOp(`/itinerary/${itinListIdx}/${LabelContentListIconJSONPath}`, icon));
      }
    }
    props.tripStateOnUpdate(ops);
  }

  return (
    <div className='p-5'>
      {
        _get(props.trip, "itinerary", [])
          .map((l: any, idx: number) => (
            <DndProvider key={l.id} backend={HTML5Backend}>
              <TripItineraryList
                trip={props.trip}
                itineraryListIdx={idx}
                itineraryList={l}
                tripStateOnUpdate={props.tripStateOnUpdate}
                onUpdateColorIcon={updateItineraryListColorIcon}
              />
              <hr className={TripItinerarySectionCss.Hr} />
            </DndProvider>
          ))
      }
    </div>
  );
}

export default ItinerarySection;
