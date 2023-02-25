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
  ContentColorOpts,
  ContentIconOpts,
  DefaultContentColor,
  LabelUiColor,
  LabelUiIcon,
  JSONPathPriceAmount,
  JSONPathLabelUiColor,
  JSONPathLabelUiIcon,
  Content,
  ItineraryList,
  ItineraryContent,
  Lodging,
  getContentColor,
  getfIndex
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
  TripItinerarySectionCss,
  TripItineraryListCss,
  TripItineraryCss,
  CommonCss,
  TripLogisticsCss,
} from '../../assets/styles/global';
import { generateKeyBetween } from '../../lib/fractional';


const ItineraryDateFmt = "eeee, do MMMM"
const DropdownItineraryContentCard = "ItineraryContentCard";


// ItineraryContentCard
interface ItineraryContentCardProps {
  content: Content
  itinListIdx: number
  itineraryList: ItineraryList
  itinCtntIdx: number
  itineraryContent: ItineraryContent
  tripOnUpdate: any

  moveCard: (id: string, to: number) => void
  findCard: (id: string) => { index: number }
  dropCard: (id: string) => void
}

const ItineraryContentCard: FC<ItineraryContentCardProps> = (props: ItineraryContentCardProps) => {

  const [isUpdatingPrice, setIsUpdatingPrice] = useState<Boolean>(false);
  const [priceAmount, setPriceAmount] = useState<Number>();
  const origCardIdx = props.findCard(props.itineraryContent.id).index;

  const { dispatch } = useMap();

  useEffect(() => {
    setPriceAmount(props.itineraryContent.price.amount);
  }, [props.itineraryContent])

  // Event Handles - DnD

  const [_, drag] = useDrag(
    () => ({
      type: DropdownItineraryContentCard,
      item: { id: props.itineraryContent.id, origCardIdx },
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
    [props.itineraryContent.id, origCardIdx, props.moveCard],
  )

  const [, drop] = useDrop(
    () => ({
      accept: DropdownItineraryContentCard,
      hover({ id }: any) {
        if (id !== props.itineraryContent.id) {
          const { index: overIndex } = props.findCard(props.itineraryContent.id)
          props.moveCard(id, overIndex)
        }
      },
    }),
    [props.findCard, props.moveCard],
  )

  const placeOnClick = (e: React.MouseEvent) => {
    dispatch({ type: ActionSetSelectedPlace, value: props.content.place })
    const event = newEventMarkerClick(props.content.place);
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
    const { itinListIdx, itinCtntIdx } = props;
    props.tripOnUpdate([
      makeRepOp(`/itinerary/${itinListIdx}/contents/${itinCtntIdx}/${JSONPathPriceAmount}`, priceAmount),
    ]);
    setIsUpdatingPrice(false);
  }

  // Renderers
  const css = {
    placeTxt: 'text-slate-600 text-sm flex items-center mb-1 hover:text-indigo-500',
    priceInputCtn: "flex w-full rounded mb-2",
    pricePill: "bg-blue-100 text-blue-800 text-xs font-semibold px-2.5 py-0.5 rounded-full mb-2 w-fit cursor-pointer",
  }

  const renderPlace = () => {
    const addr = _get(props.content, "place.name", "");
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
      className={TripItineraryCss.Ctn}
      ref={(node) => drag(drop(node))}
    >
      <div className='flex justify-between'>
        <p className={TripItineraryCss.TitleInput}>
          {props.content.title}
        </p>
      </div>
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
  );
}

// TripItineraryList

interface TripItineraryListProps {
  trip: any
  itinListIdx: number
  itineraryList: ItineraryList
  tripOnUpdate: any

  onUpdateColorIcon: (itinListIdx: number, color?: string, icon?: string) => void
}

const TripItineraryList: FC<TripItineraryListProps> = (props: TripItineraryListProps) => {

  const [isHidden, setIsHidden] = useState<boolean>(false);
  const [itinCtns, setItinContents] = useState(props.itineraryList.contents);
  const [isColorIconModalOpen, setIsColorIconModalOpen] = useState<boolean>(false);

  useEffect(() => {
    setItinContents(props.itineraryList.contents)
  }, [props.itineraryList])

  // Event Handlers

  const colorIconOnSubmit = (color?: string, icon?: string) => {
    props.onUpdateColorIcon(props.itinListIdx, color, icon)
  }

  const updateItinContents = (newItinContents: Array<ItineraryContent>) => {

  }

  // DnD Helpers

  const findCard = useCallback((id: string) => {
    const itinCtnt = _find(itinCtns, (ct: ItineraryContent) => ct.id === id);
    return { itinCtnt, index: itinCtns.indexOf(itinCtnt!) }
  }, [itinCtns]);

  const moveCard = useCallback((id: string, origIdx: number) => {
    const { itinCtnt, index } = findCard(id);
    itinCtns.splice(index, 1);
    itinCtns.splice(origIdx, 0, itinCtnt!);
    const newItinContents = itinCtns.map((x) => x);
    setItinContents(newItinContents);
  }, [findCard, itinCtns, setItinContents])

  const dropCard = useCallback((id: string) => {
    const newIdx = findCard(id).index;
    const start = _get(itinCtns[newIdx-1], "labels.fIndex", null);
    const end = _get(itinCtns[newIdx+1], "labels.fIndex", null);
    const fIndex = generateKeyBetween(start, end);


  }, [findCard, itinCtns, setItinContents])

  const [, drop] = useDrop(() => ({ accept: DropdownItineraryContentCard }))

  // Renderers
  const css = {
    ctn: "w-full mb-2",
    flightCtn: "flex items-center w-full p-3 space-x-4 text-gray-800 divide-x divide-gray-200 rounded-lg shadow",
    iconCtn: "bg-green-200 p-2 rounded-full",
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
            {fmt(parseISO(props.itineraryList.date as string), ItineraryDateFmt)}
          </p>
        </div>
        {renderSettingsDropdown()}
      </div>
    );
  }

  const renderFlights = () => {
    const flights = Object.values(_get(props.trip, "flights", {}));
    const today = parseISO(props.itineraryList.date as string);
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
      <div className={TripItineraryListCss.LodgingCtn}>
        {checkouts.map((item: any, idx: number) => render(idx, item.place, false))}
        {checkins.map((item: any, idx: number) => render(idx, item.place, true))}
      </div>
    );
  }

  const renderContents = () => {
    if (_isEmpty(itinCtns)) {
      return (
        <p className={css.noActivityTxt}>
          No activites added for today.
        </p>
      );
    }
    const sortedItinCtns = _sortBy(itinCtns, (c: ItineraryContent) => getfIndex(c));
    const listItems = sortedItinCtns.map((itinCtn: ItineraryContent, idx: number) => {
      const contents = _get(props.trip, `contents.${itinCtn.tripContentListId}.contents`, [])
      const ctnt = _find(contents, (ctn: Content) => ctn.id === itinCtn.tripContentId);
      const color = getContentColor(props.itineraryList) || DefaultContentColor;
      return (
        <li key={idx} className={TripItineraryListCss.ItinItem}>
          <span
            style={{ backgroundColor: color }}
            className={TripItineraryListCss.ItinContentIcon}
          >
            {idx + 1}
          </span>
          <ItineraryContentCard
            content={ctnt}
            itinListIdx={props.itinListIdx}
            itineraryList={props.itineraryList}
            itinCtntIdx={idx}
            itineraryContent={itinCtn}
            tripOnUpdate={props.tripOnUpdate}
            findCard={findCard}
            moveCard={moveCard}
            dropCard={dropCard}
          />
        </li>
      );
    })

    return (
      <div className={TripItineraryListCss.ContentsCtn}>
        <ol ref={drop} className={TripItineraryListCss.ContentsWrapper}>
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


interface ItinerarySectionProps {
  trip: any
  tripOnUpdate: any
}

const ItinerarySection: FC<ItinerarySectionProps> = (props: ItinerarySectionProps) => {

  const updateItineraryListColorIcon = (itinListIdx: number, color?: string, icon?: string) => {
    const ctntList = _get(props.trip, `itinerary.${itinListIdx}`);
    const currColor = _get(ctntList, `labels.${LabelUiColor}`);
    const currIcon = _get(ctntList, `labels.${LabelUiIcon}`);

    const ops = [];
    if (_isEmpty(color) && !_isEmpty(currColor)) {
      ops.push(makeRemoveOp(`/itinerary/${itinListIdx}/${JSONPathLabelUiColor}`, ""));
    }
    if (!_isEmpty(color)) {
      const op = _isEmpty(currColor) ? makeAddOp : makeRepOp;
      ops.push(op(`/itinerary/${itinListIdx}/${JSONPathLabelUiColor}`, color));
    }
    if (_isEmpty(icon) && !_isEmpty(currIcon)) {
      ops.push(makeRemoveOp(`/itinerary/${itinListIdx}/${JSONPathLabelUiIcon}`, ""));
    }
    if (!_isEmpty(icon)) {
      const op = _isEmpty(currColor) ? makeAddOp : makeRepOp;
      ops.push(op(`/itinerary/${itinListIdx}/${JSONPathLabelUiIcon}`, icon));
    }
    props.tripOnUpdate(ops);
  }

  // const updateItinContent = () => {
  //   props.tripOnUpdate([
  //     makeReplaceOp(`/itinerary/${props.itineraryListIdx}/contents`, newItinContents),
  //   ]);
  // }

  return (
    <div className='p-5'>
      {
        _get(props.trip, "itinerary", [])
          .map((l: any, idx: number) => (
            <DndProvider key={l.id} backend={HTML5Backend}>
              <TripItineraryList
                trip={props.trip}
                itinListIdx={idx}
                itineraryList={l}
                tripOnUpdate={props.tripOnUpdate}
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
