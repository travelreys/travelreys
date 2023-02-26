import React, {
  FC,
  useCallback,
  useEffect,
  useState,
} from 'react';
import _get from "lodash/get";
import _isEmpty from "lodash/isEmpty";
import _find from "lodash/find";
import _sortBy from "lodash/sortBy";
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

import ColorIconModal from './ColorIconModal';
import Dropdown from '../../components/common/Dropdown';
import HotelIcon from '../../components/icons/fill/HotelIcon';
import NotesEditor from '../../components/common/NotesEditor';
import FlightIcon from '../../components/icons/fill/FlightIcon';
import ToggleChevron from '../../components/common/ToggleChevron';

import {
  ActivityColorOpts,
  ActivityIconOpts,
  DefaultActivityColor,
  LabelUiColor,
  LabelUiIcon,
  JSONPathPriceAmount,
  JSONPathLabelUiColor,
  JSONPathLabelUiIcon,
  Activity,
  ItineraryList,
  ItineraryActivity,
  Lodging,
  getActivityColor,
  getfIndex,
  LabelFractionalIndex,
  getTripActivityForItineraryActivity
} from '../../lib/trips';
import {
  areYMDEqual,
  fmt,
  isEmptyDate,
  parseFlightDateZ,
  parseISO,
} from '../../lib/dates'
import {
  getArrivalTime,
  getDepartureTime,
  FlightItineraryTypeRoundtrip,
  Flight,
} from '../../lib/flights';
import {
  makeAddOp,
  makeRemoveOp,
  makeRepOp
} from '../../lib/jsonpatch';
import {
  MapElementID,
  newEventMarkerClick
} from '../../lib/maps';
import {
  ActionSetSelectedPlace,
  useMap
} from '../../context/maps-context';
import {
  InputCss,
  CommonCss,
  TripLogisticsCss,
} from '../../assets/styles/global';
import { generateKeyBetween } from '../../lib/fractional';


const ItineraryDateFmt = "eeee, do MMMM"
const DnDName = "ItineraryActivity";


interface ItineraryActivityCardProps {
  activity: Activity
  itinListIdx: number
  itinList: ItineraryList
  itinActIdx: number
  itinAct: ItineraryActivity
  tripOnUpdate: any

  moveCard: (id: string, to: number) => void
  findCard: (id: string) => { index: number }
  dropCard: (id: string) => void
}


const ItineraryActivityCard: FC<ItineraryActivityCardProps> = (props: ItineraryActivityCardProps) => {

  const [isUpdatingPrice, setIsUpdatingPrice] = useState<Boolean>(false);
  const [priceAmount, setPriceAmount] = useState<Number>();
  const origCardIdx = props.findCard(props.itinAct.id).index;

  const { dispatch } = useMap();

  useEffect(() => {
    setPriceAmount(props.itinAct.price.amount);
  }, [props.itinAct]);

  const [_, drag] = useDrag(
    () => ({
      type: DnDName,
      item: { id: props.itinAct.id, origCardIdx },
      collect: (monitor) => ({
        isDragging: monitor.isDragging(),
      }),
      end: (item, monitor) => {
        const { id, origCardIdx } = item
        if (monitor.didDrop()) {
          props.dropCard(id);
        } else {
          props.moveCard(id, origCardIdx)
        }
      },
    }),
    [props.itinAct.id, origCardIdx, props.moveCard],
  )

  const [, drop] = useDrop(
    () => ({
      accept: DnDName,
      hover({ id }: any) {
        if (id !== props.itinAct.id) {
          const { index: overIndex } = props.findCard(props.itinAct.id)
          props.moveCard(id, overIndex)
        }
      },
    }),
    [props.findCard, props.moveCard],
  )


  // Event Handles

  const placeOnClick = (e: React.MouseEvent) => {
    dispatch({ type: ActionSetSelectedPlace, value: props.activity.place })
    const event = newEventMarkerClick(props.activity.place);
    document.getElementById(MapElementID)?.dispatchEvent(event)
    return;
  }

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
    const { itinListIdx, itinAct } = props;
    props.tripOnUpdate([
      makeRepOp(`/itinerary/${itinListIdx}/activities/${itinAct.id}/${JSONPathPriceAmount}`, priceAmount),
    ]);
    setIsUpdatingPrice(false);
  }

  // Renderers
  const css = {
    ctn: "bg-slate-50 rounded-lg shadow-xs px-4 py-2 relative shadow",
    titleInput: "p-0 mb-1 font-bold text-gray-800 bg-transparent placeholder:text-gray-400 rounded border-0 hover:border-0 focus:ring-0 duration-400",
    autocompleteCtn: "p-1 bg-white absolute left-0 z-30 w-full border border-slate-200 rounded-lg",
    predictionWrapper: "flex items-center mb-4 cursor-pointer group",
    websiteTxt: "text-indigo-500 text-sm flex items-center",
    placeTxt: 'text-slate-600 text-sm flex items-center mb-1 hover:text-indigo-500',
    priceInputCtn: "flex w-full rounded mb-2",
    pricePill: "bg-blue-100 text-blue-800 text-xs font-semibold px-2.5 py-0.5 rounded-full mb-2 w-fit cursor-pointer",
  }

  const renderPlace = () => {
    const addr = _get(props.activity, "place.name", "");
    let placeNode = _isEmpty(addr) ? null : (
      <button type='button' onClick={placeOnClick}>
        {addr}
      </button>
    );
    return (
      <p className={css.placeTxt}>
        <MapPinIcon className={CommonCss.Icon} />
        {placeNode}
      </p>
    );
  }

  const renderPricePill = () => {
    if (isUpdatingPrice) {
      return (
        <div className={css.priceInputCtn}>
          <span className={InputCss.Label}>
            <CurrencyDollarIcon className={InputCss.Icon} />
            &nbsp;Amount
          </span>
          <input
            type="number"
            autoFocus
            className={InputCss.Input}
            value={`${priceAmount}`}
            onChange={priceOnChange}
            onBlur={priceOnBlur}
          />
        </div>
      );
    }
    return (
      <p className={css.pricePill} onClick={priceOnClick}>
        $ {priceAmount ? String(priceAmount) : "Add cost"}
      </p>
    );
  }

  return (
    <div
      className={css.ctn}
      ref={(node) => drag(drop(node))}
    >
      <div className='flex justify-between'>
        <p className={css.titleInput}>
          {props.activity.title}
        </p>
      </div>
      {renderPlace()}
      <NotesEditor
        ctnCss='p-0 mb-2'
        base64Notes={props.activity.notes}
        notesOnChange={() => { }}
        placeholder={"Notes..."}
        readOnly
      />
      {renderPricePill()}
    </div>
  );
}


interface TripItineraryListProps {
  idx: number
  list: ItineraryList
  trip: any
  tripOnUpdate: any

  onUpdateColorIcon: (idx: number, color?: string, icon?: string) => void
}

const TripItineraryList: FC<TripItineraryListProps> = (props: TripItineraryListProps) => {
  const [isHidden, setIsHidden] = useState<boolean>(false);
  const [sortedActivites, setSortedActivies] = useState([] as any);

  const [isColorIconModalOpen, setIsColorIconModalOpen] = useState<boolean>(false);

  useEffect(() => {
    const activities = Object.values(_get(props.list, "activities", {}));
    setSortedActivies(_sortBy(activities, (act) => getfIndex(act)));
  }, [props.list])


  // Event Handlers
  const colorIconOnSubmit = (color?: string, icon?: string) => {
    props.onUpdateColorIcon(props.idx, color, icon)
  }

  // DnD Helpers

  const findCard = useCallback((id: string) => {
    const act = _get(props.list, `activities.${id}`, {})
    return { act, index: sortedActivites.indexOf(act) }
  }, [sortedActivites]);

  const moveCard = useCallback((id: string, origIdx: number) => {
    const { act, index } = findCard(id);
    let newSortedList = sortedActivites.splice(index, 1);
    newSortedList = newSortedList.concat(sortedActivites.splice(origIdx, 0, act))
    setSortedActivies(newSortedList);
  }, [findCard, sortedActivites, setSortedActivies])

  const dropCard = useCallback((id: string) => {
    const newIdx = findCard(id).index;
    const start = _get(sortedActivites[newIdx-1], "labels.fIndex", null);
    const end = _get(sortedActivites[newIdx+1], "labels.fIndex", null);
    const fIndex = generateKeyBetween(start, end);
    props.tripOnUpdate([
      makeRepOp(`/itinerary/${props.idx}/activities/${id}/labels/${LabelFractionalIndex}`, fIndex)
    ]);
  }, [findCard, sortedActivites, setSortedActivies])

  const [, drop] = useDrop(() => ({ accept: DnDName }))


  // Renderers
  const css = {
    activitiesCtn: "pl-6 py-4",
    activitiesWrapper: "relative border-l border-gray-200",
    ctn: "w-full mb-2",
    flightCtn: "flex items-center w-full p-3 space-x-4 text-gray-800 divide-x divide-gray-200 rounded-lg shadow",
    iconCtn: "bg-green-200 p-2 rounded-full",
    itinActivityIcon: "absolute flex items-center justify-center w-6 h-6 rounded-full -left-3 ring-8 ring-white font-bold text-white text-sm",
    itinItem: "mb-8 ml-6",
    lodgingCtn: "w-full mb-2",
    lodgingIconWrapper: "bg-orange-200 p-2 rounded-full",
    lodgingName: "flex-1 pl-4 text-sm font-normal",
    lodgingStatus: "pl-2 font-semibold text-sm",
    lodgingWrapper: "flex items-center w-full p-3 space-x-4 text-gray-800 divide-x divide-gray-200 rounded-lg shadow",
    noActivityTxt: "text-gray-500",
  }

  const renderSettingsDropdown = () => {
    const opts = [
      <button
        type='button'
        className={CommonCss.DropdownBtn}
        onClick={() => setIsColorIconModalOpen(true)}
      >
        <SwatchIcon className={CommonCss.LeftIcon} />
        Change Color & Icon
      </button>
    ];
    const menu = (
      <EllipsisHorizontalCircleIcon className={CommonCss.DropdownIcon} />
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
            {fmt(parseISO(props.list.date as string), ItineraryDateFmt)}
          </p>
        </div>
        {renderSettingsDropdown()}
      </div>
    );
  }

  const renderFlights = () => {
    const flights = Object.values(_get(props.trip, "flights", {}));
    const today = parseISO(props.list.date as string);
    const departs = [] as Array<Flight>;

    flights.forEach((item: any) => {
      const departFlightDtb = parseFlightDateZ(getDepartureTime(item) as string);
      if (!isEmptyDate(departFlightDtb)
        && areYMDEqual(today, departFlightDtb)) {
        departs.push(item);
      }
      if (item.itineraryType === FlightItineraryTypeRoundtrip) {
        const returnFlightDt = parseFlightDateZ(getDepartureTime(item.return) as any);
        if (!isEmptyDate(returnFlightDt)
          && areYMDEqual(today, returnFlightDt)) {
          departs.push(item.return);
        }
      }
    });

    const render = (flight: Flight) => {
      const departTime = getDepartureTime(flight) as string;
      const arrTime = getArrivalTime(flight) as string;
      const timeFmt = "hh:mm aa";

      return (
        <div className={css.flightCtn}>
          <span className={css.iconCtn}>
            <FlightIcon className={CommonCss.Icon} />
          </span>
          <div className="flex pl-4">
            <span>
              <p className={TripLogisticsCss.FlightTransitTime}>
                {fmt(parseFlightDateZ(departTime), timeFmt)}
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
                {fmt(parseFlightDateZ(arrTime), timeFmt)}
              </p>
              <p className={TripLogisticsCss.FlightTransitAirportCode}>
                {flight.arrival.airport.code}
              </p>
            </span>
          </div>
        </div>
      );
    }

    return (
      <div className={css.ctn}>
        {departs.map((flight: Flight) => render(flight))}
      </div>
    );
  }

  const renderLodgings = () => {
    const lodgings = Object.values(_get(props.trip, "lodgings", {}));
    const today = parseISO(props.list.date as string);

    const render = (idx: number, place: any, checkin: boolean) => {
      return (
        <div key={idx} className={css.lodgingWrapper}>
          <span className={css.lodgingIconWrapper}>
            <HotelIcon className={CommonCss.Icon} />
          </span>
          <div className={css.lodgingName}>{place.name}</div>
          <span className={css.lodgingStatus}>
            {checkin ? "Check in" : "Check out"}
          </span>
        </div>
      );
    }

    const checkins = [] as Array<Lodging>;
    const checkouts = [] as Array<Lodging>;

    lodgings.forEach((item: any) => {
      let lod = item as Lodging;
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
      <div className={css.lodgingCtn}>
        {checkouts.map((item: any, idx: number) => render(idx, item.place, false))}
        {checkins.map((item: any, idx: number) => render(idx, item.place, true))}
      </div>
    );
  }

  const renderActivities = () => {
    if (_isEmpty(sortedActivites)) {
      return (
        <p className={css.noActivityTxt}>
          No activites added for today.
        </p>
      );
    }

    const items = sortedActivites.map((itinAct: ItineraryActivity, idx: number) => {
      const act = getTripActivityForItineraryActivity(props.trip, itinAct)
      const color = getActivityColor(props.list) || DefaultActivityColor;
      return (
        <li key={idx} className={css.itinItem}>
          <span
            style={{ backgroundColor: color }}
            className={css.itinActivityIcon}
          >
            {idx + 1}
          </span>
          <ItineraryActivityCard
            activity={act}
            itinListIdx={props.idx}
            itinList={props.list}
            itinActIdx={idx}
            itinAct={itinAct}
            tripOnUpdate={props.tripOnUpdate}
            findCard={findCard}
            moveCard={moveCard}
            dropCard={dropCard}
          />
        </li>
      );
    })

    return (
      <div className={css.activitiesCtn}>
        <ol ref={drop} className={css.activitiesWrapper}>
          {items}
        </ol>
      </div>
    );
  }

  return (
    <div className={css.ctn}>
      {renderHeader()}
      {isHidden ? null :
        <>
          {renderFlights()}
          {renderLodgings()}
          {renderActivities()}
        </>
      }
      <ColorIconModal
        isOpen={isColorIconModalOpen}
        colors={ActivityColorOpts}
        icons={Object.keys(ActivityIconOpts)}
        onClose={() => setIsColorIconModalOpen(false)}
        onSubmit={colorIconOnSubmit}
      />
    </div>
  );
}










interface ItinerarySectionProps {
  trip: any
  tripOnUpdate: any
}

const ItinerarySection: FC<ItinerarySectionProps> = (props: ItinerarySectionProps) => {

  const updateItineraryListColorIcon = (idx: number, color?: string, icon?: string) => {
    const itinList = _get(props.trip, `itinerary.${idx}`);
    const currColor = _get(itinList, `labels.${LabelUiColor}`);
    const currIcon = _get(itinList, `labels.${LabelUiIcon}`);

    const ops = [];
    if (_isEmpty(color) && !_isEmpty(currColor)) {
      ops.push(makeRemoveOp(`/itinerary/${idx}/${JSONPathLabelUiColor}`, ""));
    }
    if (!_isEmpty(color)) {
      const op = _isEmpty(currColor) ? makeAddOp : makeRepOp;
      ops.push(op(`/itinerary/${idx}/${JSONPathLabelUiColor}`, color));
    }
    if (_isEmpty(icon) && !_isEmpty(currIcon)) {
      ops.push(makeRemoveOp(`/itinerary/${idx}/${JSONPathLabelUiIcon}`, ""));
    }
    if (!_isEmpty(icon)) {
      const op = _isEmpty(currColor) ? makeAddOp : makeRepOp;
      ops.push(op(`/itinerary/${idx}/${JSONPathLabelUiIcon}`, icon));
    }
    props.tripOnUpdate(ops);
  }

  // const updateItinActivity = () => {
  //   props.tripOnUpdate([
  //     makeReplaceOp(`/itinerary/${props.itineraryListIdx}/activities`, newItinActivities),
  //   ]);
  // }


  // Renderers
  const renderItineray = () => {
    return _get(props.trip, "itinerary", {})
    .map((l: any, idx: number) => (
      <DndProvider key={l.id} backend={HTML5Backend}>
        <TripItineraryList
          idx={idx}
          list={l}
          trip={props.trip}
          tripOnUpdate={props.tripOnUpdate}
          onUpdateColorIcon={updateItineraryListColorIcon}
        />
      </DndProvider>
    ));
  }

  return (
    <div className='p-5'>
      {renderItineray()}
    </div>
  );
}

export default ItinerarySection;
