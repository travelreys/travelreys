import React, {
  FC,
  useEffect,
  useState,
} from 'react';
import _find from "lodash/find";
import _findIndex from "lodash/findIndex";
import _get from "lodash/get";
import _isEmpty from "lodash/isEmpty";
import _last from "lodash/last";
import _sortBy from "lodash/sortBy";
import { v4 as uuidv4 } from 'uuid';
import { useDebounce } from 'usehooks-ts';
import {
  CheckIcon,
  MapPinIcon,
  PlusIcon,
  SwatchIcon,
  TrashIcon
} from '@heroicons/react/24/solid';
import {
  EllipsisHorizontalCircleIcon,
  GlobeAltIcon,
} from '@heroicons/react/24/outline';

import Dropdown from '../../components/common/Dropdown';
import NotesEditor from '../../components/common/NotesEditor';
import PlaceAutocomplete from '../maps/PlaceAutocomplete';
import PlacePicturesCarousel from './PlacePicturesCarousel';
import ToggleChevron from '../../components/common/ToggleChevron';
import ColorIconModal from './ColorIconModal';
import ActivityListPin from '../maps/ListPin';

import MapsAPI, { AutocompleteResponse, PlaceDetailsResponse, placeFields } from '../../apis/maps';
import {
  Activity,
  ActivityColorOpts,
  ActivityIconOpts,
  ActivityList,
  DefaultActivityColor,
  getActivityColor,
  getActivityIcon,
  getfIndex,
  ItineraryActivity,
  ItineraryList,
  JSONPathLabelUiColor,
  JSONPathLabelUiIcon,
  LabelDelimiter,
  LabelItineraryDates,
  LabelUiColor,
  LabelUiIcon,
  makeActivity,
  makeActivityList,
  makeItineraryActivity,
} from '../../lib/trips';
import { ActionSetSelectedPlace, useMap } from '../../context/maps-context';
import { CommonCss, } from '../../assets/styles/global';
import { parseISO, fmt } from '../../lib/dates';
import { MapElementID, newEventMarkerClick } from '../../lib/maps';
import { makeAddOp, makeRemoveOp, makeRepOp } from '../../lib/jsonpatch';
import { generateKeyBetween } from '../../lib/fractional';
import { useTranslation } from 'react-i18next';

interface TripActivityProps {
  activity: Activity
  itinerary: Array<ItineraryList>

  onUpdateActivityName: (title: string, id: string) => void
  onDeleteActivity: (id: string) => void
  onUpdateActivityPlace: (id: string, place: any) => void
  onUpdateActivityNotes: (id: string, notes: string) => void
  onUpdateActivityItineraryDate: (id: string, itinListIdx: number) => void
}

const ItineraryDateFmt = "eee, do MMM";
const ItineraryBadgeDateFmt = "MMM/dd";


const TripActivity: FC<TripActivityProps> = (props: TripActivityProps) => {
  const [title, setTitle] = useState<string>("");
  const [isAddingPlace, setIsAddingPlace] = useState<boolean>(false);
  const [searchPlaceQuery, setSearchPlaceQuery] = useState<string>("");
  const [predictions, setPredictions] = useState([] as any);
  const [sessionToken, setSessionToken] = useState<string>("");

  const debouncedValue = useDebounce<string>(searchPlaceQuery, 500);
  const { dispatch } = useMap();


  useEffect(() => {
    setTitle(props.activity.title);
  }, [props.activity])

  useEffect(() => {
    if (!_isEmpty(debouncedValue)) {
      autocomplete(debouncedValue);
    }
  }, [debouncedValue])

  // API
  const autocomplete = (query: string) => {
    let token = sessionToken;
    if (_isEmpty(token)) {
      token = uuidv4();
      setSessionToken(token);
    }
    MapsAPI.placeAutocomplete(query, [], token)
    .then((res: AutocompleteResponse) => {
      setPredictions(res.predictions)
    });
  }

  // Event Handlers

  const titleInputOnBlur = () => {
    props.onUpdateActivityName(title, props.activity.id);
  }

  const deleteBtnOnClick = () => {
    props.onDeleteActivity(props.activity.id);
  }

  const predictionOnSelect = (placeID: string) => {
    MapsAPI.placeDetails(placeID, placeFields, sessionToken)
    .then((res: PlaceDetailsResponse) => {
      setPredictions([]);
      props.onUpdateActivityPlace(props.activity.id, res.place)
    })
    .finally(() => {
      setSessionToken("");
    });
  }

  const placeOnClick = (e: React.MouseEvent) => {
    if (e.detail === 1) {
      dispatch({
        type: ActionSetSelectedPlace,
        value: props.activity.place
      });
      const event = newEventMarkerClick(props.activity.place);
      document.getElementById(MapElementID)?.dispatchEvent(event)
      return;
    }

    if (e.detail === 2) {
      setIsAddingPlace(true);
      return;
    }
  }

  const notesOnChange = (activity: string) => {
    props.onUpdateActivityNotes(props.activity.id, activity)
  }

  const itinOptOnClick = (itinListIdx: number) => {
    props.onUpdateActivityItineraryDate(props.activity.id, itinListIdx);
  }

  // Renderers

  const css = {
    autocompleteCtn: "p-1 bg-white absolute left-0 z-30 w-full border border-slate-200 rounded-lg",
    ctn: "bg-slate-50 rounded-lg shadow-xs mb-4 p-4 relative",
    itineraryDateBtn: "flex items-center w-full justify-between hover:text-indigo-500 text-align-right",
    titleInput: "p-0 mb-1 font-bold text-gray-800 bg-transparent placeholder:text-gray-400 rounded border-0 hover:border-0 focus:ring-0 duration-400",
    websiteLink: "flex items-center mb-1",
    websiteTxt: "text-indigo-500 text-sm flex items-center",
    addItineraryBtn: "text-xs text-gray-800 font-bold bg-indigo-200 rounded-full px-2 py-1 hover:bg-indigo-400",
    itineraryBadge: "bg-indigo-100 text-indigo-800 text-xs font-medium mr-2 px-2.5 py-0.5 rounded",
    placeCtn: "text-slate-600 text-sm flex items-center mb-1 hover:text-indigo-500",
    placeInput: "p-0 mb-1 text-sm text-gray-600 bg-transparent placeholder:text-gray-400 rounded border-0 hover:border-0 focus:ring-0 duration-400",
   }

  const renderSettingsDropdown = () => {
    const opts = [
      <button
        type='button'
        className={CommonCss.DeleteBtn}
        onClick={deleteBtnOnClick}
      >
        <TrashIcon className={CommonCss.LeftIcon} />
        Delete
      </button>
    ];
    const menu = (
      <EllipsisHorizontalCircleIcon
        className={CommonCss.DropdownIcon} />
    );
    return <Dropdown menu={menu} opts={opts} />
  }

  const renderItineraryDropdown = () => {
    // Format of itinerary dates label:
    // activity.labels[JSONPathLabelItineraryDates] = "d1|d2|d3"
    const dates = _get(props.activity, `labels.${LabelItineraryDates}`, "")
      .split(LabelDelimiter)
      .filter((dt: string) => !_isEmpty(dt));

    const opts = props.itinerary.map((l: ItineraryList, idx: number) => (
      <button
        type='button'
        className={css.itineraryDateBtn}
        onClick={() => { itinOptOnClick(idx) }}
      >
        {fmt(parseISO(l.date as string), ItineraryDateFmt)}
        {dates.includes(l.date as string) ?
          <CheckIcon className={CommonCss.Icon} /> : null}
      </button>
    ));

    const datesBadges = dates
      .map((dt: string) => (
        <span key={dt} className={css.itineraryBadge}>
          {fmt(parseISO(dt), ItineraryBadgeDateFmt)}
        </span>
      ));

    const emptyBtn = (
      <span className={css.addItineraryBtn}>
        Add to Itinerary
      </span>
    );
    const menu = dates.length === 0 ? emptyBtn : datesBadges;
    return <Dropdown menu={menu} opts={opts} />
  }

  const renderHeader = () => {
    return (
      <div className='flex justify-between'>
        <input
          type="text"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          onBlur={titleInputOnBlur}
          placeholder="Add a name"
          className={css.titleInput}
        />
        <div className='flex items-center'>
          {renderItineraryDropdown()}
          &nbsp;&nbsp;
          {renderSettingsDropdown()}
        </div>
      </div>
    );
  }

  const renderPlace = () => {
    let placeNode;
    if (isAddingPlace) {
      placeNode = (
        <input
          type="text"
          autoFocus
          value={searchPlaceQuery}
          onChange={(e) => setSearchPlaceQuery(e.target.value)}
          onBlur={() => { setIsAddingPlace(false); }}
          placeholder="name, address..."
          className={css.placeInput}
        />
      );
    } else {
      const addr = _get(props.activity, "place.name", "");
      placeNode = (
        <button
          type='button'
          onClick={_isEmpty(addr)
            ? (e) => { setIsAddingPlace(true) }
            : (e) => { placeOnClick(e) }}
        >
          {_isEmpty(addr) ? "Click here to add a location..." : addr}
        </button>
      );
    }

    return (
      <div className={css.placeCtn}>
        <MapPinIcon className='h-4 w-4 mr-1' />
        {placeNode}
      </div>
    );
  }

  const renderWebsite = () => {
    const website = _get(props.activity, "place.website", "")
    if (_isEmpty(website)) {
      return null
    }
    return (
      <a
        className={css.websiteLink}
        href={website}
        target="_blank"
        rel="noreferrer"
      >
        <GlobeAltIcon className={CommonCss.LeftIcon} />
        <span className={css.websiteTxt}>Website</span>
      </a>
    );
  }

  const renderPlacePicturesCarousel = () => {
    const photos = _get(props.activity, "place.photos", []);
    if (_isEmpty(photos)) {
      return null;
    }
    return <PlacePicturesCarousel photos={photos} />
  }

  return (
    <div className={css.ctn}>
      {renderHeader()}
      {renderPlace()}
      <PlaceAutocomplete
        predictions={predictions}
        onSelect={predictionOnSelect}
      />
      {renderWebsite()}
      <NotesEditor
        ctnCss='p-0 mb-2'
        base64Notes={props.activity.notes}
        notesOnChange={notesOnChange}
        placeholder={"Notes..."}
      />
      {renderPlacePicturesCarousel()}
    </div>
  );
}


interface TripActivityListProps {
  itinerary: any
  list: ActivityList

  onDeleteList: (actListId: string) => void
  onUpdateName: (name: string, actListId: string) => void
  onAddActivity: (title: string, actListId: string) => void
  onUpdateColorIcon: (actListId: string, color?: string, icon?: string) => void
  onUpdateActivityName: (title: string, id: string, listId: string) => void
  onDeleteActivity: (id: string, actListId: string) => void
  onUpdateActivityPlace: (id: string, place: any, actListId: string) => void
  onUpdateActivityNotes: (id: string, notes: string, actListId: string) => void
  onUpdateActivityItineraryDate: (id: string, itinListIdx: number, actListId: string) => void
}

const TripActivityList: FC<TripActivityListProps> = (props: TripActivityListProps) => {
  const [name, setName] = useState<string>("");
  const [newActivityTitle, setNewActivityTitle] = useState("");
  const [isHidden, setIsHidden] = useState<boolean>(false);
  const [isColorIconModalOpen, setIsColorIconModalOpen] = useState<boolean>(false);

  useEffect(() => {
    setName(_get(props.list, "name", ""));
  }, [props.list])

  const nameOnBlur = () => {
    props.onUpdateName(name, props.list.id);
  }

  const deleteBtnOnClick = () => {
    props.onDeleteList(props.list.id)
  }

  const newActivityBtnOnClick = () => {
    props.onAddActivity(newActivityTitle, props.list.id)
    setNewActivityTitle("");
  }

  const colorIconOnSubmit = (color?: string, icon?: string) => {
    props.onUpdateColorIcon(props.list.id, color, icon)
  }

  const onUpdateActivityName = (title: string, id: string) => {
    props.onUpdateActivityName(title, id, props.list.id);
  }

  const onDeleteActivity = (id: string) => {
    props.onDeleteActivity(id, props.list.id);
  }

  const onUpdateActivityPlace = (id: string, place: any) => {
    props.onUpdateActivityPlace(id, place, props.list.id);
  }

  const onUpdateActivityNotes = (id: string, notes: string) => {
    props.onUpdateActivityNotes(id, notes, props.list.id);
  }

  const onUpdateActivityItineraryDate = (id: string, itinListIdx: number) => {
    props.onUpdateActivityItineraryDate(id, itinListIdx, props.list.id)
  }

  // Renderers
  const css = {
    ctn: "rounded-lg shadow-xs mb-4",
    headerAct: "flex mb-2 w-full justify-between items-center",
    headerInputAct: "flex flex-1",
    nameInput: "p-0 w-full text-xl mb-1 sm:text-2xl font-bold text-gray-800 placeholder:text-gray-400 rounded border-0 hover:bg-gray-300 hover:border-0 focus:ring-0 focus:p-1 duration-500",
    newActivityAct: "flex my-4 w-full",
    newActivityInput: "flex-1 mr-1 text-md sm:text-md font-bold text-gray-800 placeholder:font-normal placeholder:text-gray-300 placeholder:italic rounded border-0 bg-gray-100 hover:border-0 focus:ring-0",
    newActivityBtn: "text-green-600 w-1/12 hover:bg-green-50 rounded-lg text-sm font-bold inline-flex justify-around items-center",
  }

  const renderSettingsDropdown = () => {
    const opts = [
      (
        <button
          type='button'
          className={CommonCss.DropdownBtn}
          onClick={() => setIsColorIconModalOpen(true)}
        >
          <SwatchIcon className={CommonCss.LeftIcon} />
          Change Color & Icon
        </button>
      ),
      (
        <button
          type='button'
          className={CommonCss.DeleteBtn}
          onClick={deleteBtnOnClick}
        >
          <TrashIcon className={CommonCss.LeftIcon} />
          Delete
        </button>
      )
    ];
    const menu = (
      <EllipsisHorizontalCircleIcon className={CommonCss.DropdownIcon} />
    );
    return <Dropdown menu={menu} opts={opts} />
  }

  const renderHeader = () => {
    const color = getActivityColor(props.list) || DefaultActivityColor;
    const icon = getActivityIcon(props.list) || DefaultActivityColor;
    return (
      <div className={css.headerAct}>
        <div className={css.headerInputAct}>
          <ToggleChevron
            onClick={() => setIsHidden(!isHidden)}
            isHidden={isHidden}
          />
          <ActivityListPin color={color} icon={ActivityIconOpts[icon]} />
          <input
            type="text"
            className={css.nameInput}
            value={name}
            onChange={(e) => setName(e.target.value)}
            onBlur={nameOnBlur}
            placeholder='Add a title (e.g, "Food to try")'
          />
        </div>
        {renderSettingsDropdown()}
      </div>
    );
  }

  const renderTripActivity = () => {
    return Object.values(_get(props.list, "activities", {}))
      .map((activity: any, idx: number) => (
        <TripActivity
          key={idx}
          itinerary={props.itinerary}
          activity={activity}
          onUpdateActivityName={onUpdateActivityName}
          onDeleteActivity={onDeleteActivity}
          onUpdateActivityPlace={onUpdateActivityPlace}
          onUpdateActivityNotes={onUpdateActivityNotes}
          onUpdateActivityItineraryDate={onUpdateActivityItineraryDate}
        />
      ));
  }

  const renderAddNewActivity = () => {
    return (
      <div className={css.newActivityAct}>
        <input
          type="text"
          value={newActivityTitle}
          onChange={(e) => { setNewActivityTitle(e.target.value) }}
          placeholder="Add an activity..."
          className={css.newActivityInput}
        />
        <button
          onClick={newActivityBtnOnClick}
          className={css.newActivityBtn}
        >
          <PlusIcon className={CommonCss.Icon} />
        </button>
      </div>
    );
  }

  return (
    <div className={css.ctn}>
      {renderHeader()}
      {isHidden ? null :
        <>
          {renderTripActivity()}
          {renderAddNewActivity()}
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

interface ActivitySectionProps {
  trip: any
  tripOnUpdate: any
}

const ActivitySection: FC<ActivitySectionProps> = (props: ActivitySectionProps) => {
  const [isHidden, setIsHidden] = useState(false);
  const { t } = useTranslation();

  // Event Handlers
  const addListBtnOnClick = () => {
    let list =  makeActivityList();
    props.tripOnUpdate([makeAddOp(`/activities/${list.id}`, list)]);
  }

  const activityListUpdateName = (name: string, id: string) => {
    props.tripOnUpdate([makeRepOp(`/activities/${id}/name`, name)]);
  }

  const activityListAddActivity = (title: string, id: string) => {
    const act = makeActivity(title);
    props.tripOnUpdate([makeAddOp(`/activities/${id}/activities/${act.id}`, act)])
  }

  const deleteActivityList = (id: string) => {
    const ops = [makeRemoveOp(`/activities/${id}`, "")]
    _get(props.trip, "itinerary", [])
      .forEach((itinList: ItineraryList) => {
        itinList.activities
          .filter((itinAct: ItineraryActivity) => itinAct.activityListId === id)
          .forEach((itinAct: ItineraryActivity) => {
            ops.unshift(
              makeRemoveOp(`/itinerary/${itinList.id}/activities/${itinAct.id}`, "")
            );
          });
      });
    props.tripOnUpdate(ops);
  }

  const updateActivityListColorIcon = (id: string, color?: string, icon?: string) => {
    const actList = _get(props.trip, `activities.${id}`);
    const currColor = _get(actList, `labels.${LabelUiColor}`);
    const currIcon = _get(actList, `labels.${LabelUiIcon}`);

    const ops = [];
    if (_isEmpty(color) && !_isEmpty(currColor)) {
      ops.push(makeRemoveOp(`/activities/${id}/${JSONPathLabelUiColor}`, ""));
    }
    if (!_isEmpty(color)) {
      const op = _isEmpty(currColor) ? makeAddOp : makeRepOp;
      ops.push(op(`/activities/${id}/${JSONPathLabelUiColor}`, color));
      ops.push(op(`/activities/${id}/${JSONPathLabelUiColor}`, color));
    }

    if (_isEmpty(icon) && !_isEmpty(currIcon)) {
      ops.push(makeRemoveOp(`/activities/${id}/${JSONPathLabelUiIcon}`, ""));
    }
    if (!_isEmpty(icon)) {
      const op = _isEmpty(currColor) ? makeAddOp : makeRepOp;
      ops.push(op(`/activities/${id}/${JSONPathLabelUiIcon}`, icon));
    }
    props.tripOnUpdate(ops);
  }

  const updateActivityName = (title: string, id: string, listId: string) => {
    props.tripOnUpdate([makeRepOp( `/activities/${listId}/activities/${id}/title`, title)]);
  }

  const deleteActivity = (id: string, listId: string) => {
    props.tripOnUpdate([makeRemoveOp(`/activities/${listId}/activities/${id}`, "")]);
  }

  const updateActivityPlace = (id: string, place: any, listId: string) => {
    props.tripOnUpdate([makeRepOp(`/activities/${listId}/activities/${id}/place`, place),]);
  }

  const updateActivityNotes = (id: string, notes: string, listId: string) => {
    props.tripOnUpdate([
      makeRepOp(`/activities/${listId}/activities/${id}/notes`, notes)
    ]);
  }

  const updateActivityItineraryDate = async (id: string, itinListId: number, actListId: string) => {
    const activity = _get(props.trip, `activities.${actListId}.activities.${id}`, {}) as Activity;
    const itinList = _get(props.trip, `itinerary.${itinListId}`, {}) as ItineraryList;
    const itinListDt = itinList.date as string;

    const ops = [];

    // Update activity labels, Format of itinerary dates label:
    // activity.labels[LabelItineraryDates] = "d1|d2|d3"
    let currItinDts = _get(activity, `labels.${LabelItineraryDates}`, "")
      .split(LabelDelimiter)
      .filter((dt: string) => !_isEmpty(dt));

    const isRemove = currItinDts.includes(itinListDt);

    let newItinDts;
    if (isRemove) {
      // 1. Remove from activity label if it exists
      newItinDts = currItinDts.filter((dt: string) => dt !== itinListDt);

      // 2. Remove ItineraryActivity from ItineraryList
      let itinActId = _findIndex(itinList.activities, (act) => act.activityId === activity.id);
      ops.push(makeRemoveOp(`/itinerary/${itinListId}/activities/${itinActId}`, "",));

    } else {
      // 1. Add to activity label if its a new itinerary date
      newItinDts = _sortBy(currItinDts.concat([itinListDt]));

      // 2. Add ItineraryActivity to ItineraryList
      const sortedActivites = _sortBy(itinList.activities, (act) => getfIndex(act))
      const start = _get(sortedActivites.slice(-1), "0.labels.fIndex", null);
      const fIndex = generateKeyBetween(start, null)

      const itinAct = makeItineraryActivity(activity.id, actListId, fIndex);
      ops.push(makeAddOp(`/itinerary/${itinListId}/activities/${itinAct.id}`, itinAct))
    }

    // Update activity's itinerary dates
    if (currItinDts.length !== 0) {
      ops.unshift(makeRepOp(
        `/activities/${actListId}/activities/${id}/labels/${LabelItineraryDates}`,
        newItinDts.join(LabelDelimiter)));
    } else {
      ops.unshift(makeAddOp(
        `/activities/${actListId}/activities/${id}/labels/${LabelItineraryDates}`,
        newItinDts.join(LabelDelimiter)));
    }

    props.tripOnUpdate(ops);
  }

  // Renderers
  const css = {
    headerAct: "flex justify-between mb-4",
    header: "text-2xl sm:text-3xl font-bold text-slate-700",
    addBtn: "text-white py-2 px-4 bg-indigo-500 rounded-lg text-sm font-semibold",
    toggleBtn: "mr-2",
  }

  const renderHeader = () => {
    return (
      <div className={css.headerAct}>
        <div>
          <ToggleChevron
            isHidden={isHidden}
            onClick={() => {setIsHidden(!isHidden)}}
          />
          <span className={css.header}>
            {t('tripPage.activitySection.header')}
          </span>
        </div>
        <button
          className={css.addBtn}
          onClick={addListBtnOnClick}
        >
          {t('tripPage.activitySection.addBtn')}
        </button>
      </div>
    );
  }

  const renderActivityLists = () => {
    if (isHidden) {
      return null;
    }
    const lists = Object.values(_get(props.trip, "activities", {}));
    return lists.map((l: any) => (
      <TripActivityList
        key={l.id}
        itinerary={props.trip.itinerary}
        list={l}
        onDeleteList={deleteActivityList}
        onUpdateName={activityListUpdateName}
        onAddActivity={activityListAddActivity}
        onUpdateColorIcon={updateActivityListColorIcon}
        onUpdateActivityName={updateActivityName}
        onDeleteActivity={deleteActivity}
        onUpdateActivityPlace={updateActivityPlace}
        onUpdateActivityNotes={updateActivityNotes}
        onUpdateActivityItineraryDate={updateActivityItineraryDate}
      />
    ));
  }

  return (
    <div className='p-5'>
      {renderHeader()}
      {renderActivityLists()}
    </div>
  );
}


export default ActivitySection;


// // 3. Add ItineraryActivityRoute to ItineraryList
// if (itinListCtnt.length > 0) {
//   const lastItinCtn = _last(itinListCtnt);
//   const lastCtnt = _find(
//     _get(props.trip, `activities[${lastItinCtn?.activityListId}].activities`),
//     (ctn: Activity) => ctn.id === lastItinCtn?.activityId,
//   );

//   const lastCtntPlaceID = _get(lastCtnt, "place.place_id");
//   const ctnPlaceID = _get(activity, "place.place_id");

//   if (lastCtntPlaceID && ctnPlaceID) {
//     const resp = await MapsAPI.directions(lastCtntPlaceID, ctnPlaceID, ModeDriving);
//     if (resp.data.routeList.length > 0) {
//       ops.push(makeAddOp(`/itinerary/${itinListIdx}/routes/-`, resp.data.routeList[0]));
//     }
//   }
// }

  // 3. Remove ItineraryActivityRoute from ItineraryList
  // ops.push(makeRepOp(`/itinerary/${itinListIdx}/routes`, []));
  // if (itinActIdx === 0) {
  //   const routeIdx = itinActIdx + 1;
  //   if (routeIdx < itinList.routes.length) {
  //     ops.push(makeRemoveOp(`/itinerary/${itinListIdx}/routes/${routeIdx}`, ""));
  //   }
  // } else if (itinActIdx === itinListCtnt.length - 1) {
  //   const routeIdx = itinActIdx - 1;
  //   ops.push(makeRemoveOp(`/itinerary/${itinListIdx}/routes/${routeIdx}`, ""));
  // } else {
  //   ops.push(makeRemoveOp(`/itinerary/${itinListIdx}/routes/${itinActIdx}`, ""));
  //   ops.push(makeRemoveOp(`/itinerary/${itinListIdx}/routes/${itinActIdx - 1}`, ""));
  // }