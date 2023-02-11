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

import Dropdown from '../Dropdown';
import NotesEditor from '../NotesEditor';
import PlaceAutocomplete from '../maps/PlaceAutocomplete';
import PlacePicturesCarousel from './PlacePicturesCarousel';
import ToggleChevron from '../ToggleChevron';
import ColorIconModal from './ColorIconModal';

import TripsSyncAPI from '../../apis/tripsSync';
import MapsAPI, { ModeDriving, placeFields } from '../../apis/maps';
import {
  Trips,
  LabelContentItineraryDates,
  LabelContentItineraryDatesJSONPath,
  LabelContentItineraryDatesDelimeter,
  ContentColorOpts,
  ContentIconOpts,
  LabelContentListColor,
  LabelContentListColorJSONPath,
  LabelContentListIconJSONPath,
  LabelContentListIcon,
} from '../../apis/trips';
import {
  ActionNameSetSelectedPlace,
  useMap,
} from '../../context/maps-context';
import {
  CommonCss,
  TripContentCss,
  TripContentListCss,
  TripContentSectionCss,
} from '../../styles/global';
import { parseISO, printFmt } from '../../utils/dates';
import { MapElementID, newEventMarkerClick } from '../maps/common';
import isEmpty from 'lodash/isEmpty';


/////////////
// Content //
/////////////

interface TripContentProps {
  content: Trips.Content
  contentIdx: number
  itinerary: Array<Trips.ItineraryList>

  onUpdateContentName: (title: string, idx: number) => void
  onDeleteContent: (idx: number) => void
  onUpdateContentPlace: (idx: number, place: any) => void
  onUpdateContentNotes: (idx: number, notes: string) => void
  onUpdateContentItineraryDate: (idx: number, itinListIdx: number) => void
}

const ItineraryDateFmt = "eee, do MMM";
const ItineraryBadgeDateFmt = "MMM/dd";


const TripContent: FC<TripContentProps> = (props: TripContentProps) => {
  const [title, setTitle] = useState<string>("");
  const [isAddingPlace, setIsAddingPlace] = useState<boolean>(false);
  const [searchPlaceQuery, setSearchPlaceQuery] = useState<string>("");
  const [predictions, setPredictions] = useState([] as any);
  const [sessionToken, setSessionToken] = useState<string>("");

  const debouncedValue = useDebounce<string>(searchPlaceQuery, 500);

  const { dispatch } = useMap();

  useEffect(() => {
    setTitle(props.content.title);
  }, [props.content])

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
      .then((res) => {
        setPredictions(_get(res, "data.predictions", []))
      });
  }

  // Event Handlers - Header
  const titleInputOnBlur = () => {
    props.onUpdateContentName(title, props.contentIdx);
  }

  const deleteBtnOnClick = () => {
    props.onDeleteContent(props.contentIdx);
  }

  // Event Handlers - Places
  const predictionOnSelect = (placeID: string) => {
    MapsAPI.placeDetails(placeID, placeFields, sessionToken)
      .then((res) => {
        setPredictions([]);
        const place = _get(res, "data.place", {});
        props.onUpdateContentPlace(props.contentIdx, place)
      })
      .finally(() => {
        setSessionToken("");
      });
  }

  const placeOnClick = (e: React.MouseEvent) => {
    if (e.detail == 1) {
      dispatch({
        type: ActionNameSetSelectedPlace,
        value: props.content.place
      });
      const event = newEventMarkerClick(props.content.place);
      document.getElementById(MapElementID)?.dispatchEvent(event)
      return;
    }

    if (e.detail == 2) {
      setIsAddingPlace(true);
      return;
    }
  }

  // Event Handlers - Notes
  const notesOnChange = (content: string) => {
    props.onUpdateContentNotes(props.contentIdx, content)
  }

  // Event Handlers - Itinerary
  const itinOptOnClick = (itinListIdx: number) => {
    props.onUpdateContentItineraryDate(props.contentIdx, itinListIdx);
  }

  // Renderers

  const renderSettingsDropdown = () => {
    const opts = [
      <button
        type='button'
        className={TripContentCss.DeleteBtn}
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
    // content.labels[LabelContentItineraryDatesJSONPath] = "d1|d2|d3"

    const dates = _get(props.content, LabelContentItineraryDatesJSONPath, "")
      .split(LabelContentItineraryDatesDelimeter)
      .filter((dt: string) => !_isEmpty(dt));

    const opts = props.itinerary.map((l: Trips.ItineraryList, idx: number) => (
      <button
        type='button'
        className={TripContentCss.ItineraryDateBtn}
        onClick={() => { itinOptOnClick(idx) }}
      >
        {printFmt(parseISO(l.date as string), ItineraryDateFmt)}
        {dates.includes(l.date as string) ?
          <CheckIcon className={CommonCss.Icon} /> : null}
      </button>
    ));

    const datesBadges = dates
      .map((dt: string) => (
        <span key={dt} className={TripContentCss.ItineraryBadge}>
          {printFmt(parseISO(dt), ItineraryBadgeDateFmt)}
        </span>
      ));

    const emptyBtn = (
      <span className={TripContentCss.AddItineraryBtn}>
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
          className={TripContentCss.TitleInput}
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
          className={TripContentCss.PlaceInput}
        />
      );
    } else {
      const addr = _get(props.content, "place.name", "");
      placeNode = (
        <button
          type='button'
          onClick={
            _isEmpty(addr) ?
              (e) => { setIsAddingPlace(true) } :
              (e) => { placeOnClick(e) }
          }
        >
          {_isEmpty(addr) ?
            "Click here to add a location..." : addr
          }
        </button>
      );
    }

    return (
      <div className={TripContentCss.PlaceCtn}>
        <MapPinIcon className='h-4 w-4 mr-1' />
        {placeNode}
      </div>
    );
  }

  const renderWebsite = () => {
    const website = _get(props.content, "place.website", "")
    if (_isEmpty(website)) {
      return null
    }
    return (
      <a
        className={TripContentCss.WebsiteLink}
        href={website}
        target="_blank"
      >
        <GlobeAltIcon className={CommonCss.LeftIcon} />
        <span className={TripContentCss.WebsiteTxt}>Website</span>
      </a>
    );
  }

  const renderPlacePicturesCarousel = () => {
    const photos = _get(props.content, "place.photos", []);
    if (_isEmpty(photos)) {
      return null;
    }
    return <PlacePicturesCarousel photos={photos} />
  }

  return (
    <div className={TripContentCss.Ctn}>
      {renderHeader()}
      {renderPlace()}
      <PlaceAutocomplete
        predictions={predictions}
        onSelect={predictionOnSelect}
      />
      {renderWebsite()}
      <NotesEditor
        ctnCss='p-0 mb-2'
        base64Notes={props.content.notes}
        notesOnChange={notesOnChange}
        placeholder={"Notes..."}
      />
      {renderPlacePicturesCarousel()}
    </div>
  );
}


/////////////////
// ContentList //
/////////////////


interface ContentListProps {
  itinerary: any
  contentList: Trips.ContentList

  onDeleteList: (contentListID: string) => void
  onUpdateName: (name: string, contentListID: string) => void
  onAddContent: (title: string, contentListID: string) => void
  onUpdateColorIcon: (color: string|undefined, icon: string|undefined, contentListID: string) => void
  onUpdateContentName: (title: string, idx: number, contentListID: string) => void
  onDeleteContent: (idx: number, contentListID: string) => void
  onUpdateContentPlace: (idx: number, place: any, contentListID: string) => void
  onUpdateContentNotes: (idx: number, notes: string, contentListID: string) => void
  onUpdateContentItineraryDate: (idx: number, itinListIdx: number, contentListID: string) => void
}

const ContentList: FC<ContentListProps> = (props: ContentListProps) => {

  const [name, setName] = useState<string>("");
  const [newContentTitle, setNewContentTitle] = useState("");
  const [isHidden, setIsHidden] = useState<boolean>(false);
  const [isColorIconModalOpen, setIsColorIconModalOpen] = useState<boolean>(false);

  useEffect(() => {
    setName(_get(props.contentList, "name", ""));
  }, [props.contentList])

  // Event Handlers - Content List

  const nameOnBlur = () => {
    props.onUpdateName(name, props.contentList.id);
  }

  const deleteBtnOnClick = () => {
    props.onDeleteList(props.contentList.id)
  }

  const newContentBtnOnClick = () => {
    props.onAddContent(newContentTitle, props.contentList.id)
    setNewContentTitle("");
  }

  const colorIconOnSubmit = (color: string | undefined, icon: string | undefined) => {
    props.onUpdateColorIcon(color, icon, props.contentList.id)
  }

  // Event Handlers - Content

  const onUpdateContentName = (title: string, idx: number) => {
    props.onUpdateContentName(title, idx, props.contentList.id);
  }

  const onDeleteContent = (idx: number) => {
    props.onDeleteContent(idx, props.contentList.id);
  }

  const onUpdateContentPlace = (idx: number, place: any) => {
    props.onUpdateContentPlace(idx, place, props.contentList.id);
  }

  const onUpdateContentNotes = (idx: number, notes: string) => {
    props.onUpdateContentNotes(idx, notes, props.contentList.id);
  }

  const onUpdateContentItineraryDate = (idx: number, itinListIdx: number) => {
    props.onUpdateContentItineraryDate(idx, itinListIdx, props.contentList.id)
  }


  // Renderers
  const renderSettingsDropdown = () => {
    const opts = [
      (<button
        type='button'
        className={TripContentListCss.ChooseColorBtn}
        onClick={() => setIsColorIconModalOpen(true)}
      >
        <SwatchIcon className={CommonCss.LeftIcon} />
        Change Color & Icon
      </button>),
      (<button
        type='button'
        className={TripContentListCss.DeleteBtn}
        onClick={deleteBtnOnClick}
      >
        <TrashIcon className={CommonCss.LeftIcon} />
        Delete
      </button>)
    ];
    const menu = (
      <EllipsisHorizontalCircleIcon
        className={CommonCss.DropdownIcon} />
    );
    return <Dropdown menu={menu} opts={opts} />
  }

  const renderHeader = () => {
    return (
      <div className='flex mb-2 w-full justify-between'>
        <div className='flex flex-1'>
          <ToggleChevron
            onClick={() => setIsHidden(!isHidden)}
            isHidden={isHidden}
          />
          <input
            type="text"
            value={name}
            onChange={(e) => setName(e.target.value)}
            onBlur={nameOnBlur}
            placeholder={`Add a title (e.g, "Food to try")`}
            className={TripContentListCss.NameInput}
          />
        </div>
        {renderSettingsDropdown()}
      </div>
    );
  }

  const renderTripContent = () => {
    return _get(props.contentList, "contents", [])
      .map((content: any, idx: number) => (
        <TripContent
          key={idx}
          itinerary={props.itinerary}
          content={content}
          contentIdx={idx}
          onUpdateContentName={onUpdateContentName}
          onDeleteContent={onDeleteContent}
          onUpdateContentPlace={onUpdateContentPlace}
          onUpdateContentNotes={onUpdateContentNotes}
          onUpdateContentItineraryDate={onUpdateContentItineraryDate}
        />
      ));
  }

  const renderAddNewContent = () => {
    return (
      <div className={TripContentListCss.NewContentCtn}>
        <input
          type="text"
          value={newContentTitle}
          onChange={(e) => { setNewContentTitle(e.target.value) }}
          placeholder="Add an activity..."
          className={TripContentListCss.NewContentInput}
        />
        <button
          onClick={() => { newContentBtnOnClick() }}
          className={TripContentListCss.NewContentBtn}
        >
          <PlusIcon className={CommonCss.Icon} />
        </button>
      </div>
    );
  }

  return (
    <div className={TripContentListCss.Ctn}>
      {renderHeader()}
      {isHidden ? null :
        <>
          {renderTripContent()}
          {renderAddNewContent()}
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


////////////////////
// ContentSection //
////////////////////

interface ContentSectionProps {
  trip: any
  tripStateOnUpdate: any
}

const ContentSection: FC<ContentSectionProps> = (props: ContentSectionProps) => {

  const [isHidden, setIsHidden] = useState(false);

  // Event Handlers -- Content List
  const addContentListBtnOnClick = () => {
    let list: Trips.ContentList = {
      id: uuidv4(),
      contents: new Array<Trips.Content>(),
      labels: new Map<string, string>(),
    }
    props.tripStateOnUpdate([TripsSyncAPI.makeAddOp(`/contents/${list.id}`, list)]);
  }

  const contentListUpdateName = (name: string, contentListID: string) => {
    props.tripStateOnUpdate([
      TripsSyncAPI.newReplaceOp(`/contents/${contentListID}/name`, name)
    ]);
  }

  const contentListAddContent = (title: string, contentListID: string) => {
    const content: Trips.Content = {
      id: uuidv4(),
      title: title,
      notes: "",
      place: {},
      labels: new Map<string, string>(),
      comments: [],
    }
    props.tripStateOnUpdate([
      TripsSyncAPI.makeAddOp(`/contents/${contentListID}/contents/-`, content),
    ])
  }

  const deleteContentList = (contentListID: string) => {
    props.tripStateOnUpdate([
      TripsSyncAPI.makeRemoveOp(`/contents/${contentListID}`, "")
    ]);
  }

  const updateContentListColorIcon = (color: string | undefined, icon: string | undefined, contentListID: string) => {
    const ctntList = _get(props.trip, `contents.${contentListID}`);
    const colorLabel = _get(ctntList, `labels.${LabelContentListColor}`);
    const iconLabel = _get(ctntList, `labels.${LabelContentListIcon}`);

    const ops = [];
    if (_isEmpty(color) && !_isEmpty(colorLabel)) {
      ops.push(TripsSyncAPI.makeRemoveOp(`/contents/${contentListID}/${LabelContentListColorJSONPath}`, ""));
    }
    if (!_isEmpty(color)) {
      if (isEmpty(colorLabel)) {
        ops.push(TripsSyncAPI.makeAddOp(`/contents/${contentListID}/${LabelContentListColorJSONPath}`, color));
      } else {
        ops.push(TripsSyncAPI.newReplaceOp(`/contents/${contentListID}/${LabelContentListColorJSONPath}`, color));
      }
    }

    if (_isEmpty(icon) && !_isEmpty(iconLabel)) {
      ops.push(TripsSyncAPI.makeRemoveOp(`/contents/${contentListID}/${LabelContentListIconJSONPath}`, ""));
    }
    if (!_isEmpty(icon)) {
      if (isEmpty(colorLabel)) {
        ops.push(TripsSyncAPI.makeAddOp(`/contents/${contentListID}/${LabelContentListIconJSONPath}`, icon));
      } else {
        ops.push(TripsSyncAPI.newReplaceOp(`/contents/${contentListID}/${LabelContentListIconJSONPath}`, icon));
      }
    }
    props.tripStateOnUpdate(ops);
  }

  // Event Handlers -- Content
  const onUpdateContentName = (title: string, idx: number, contentListID: string) => {
    props.tripStateOnUpdate([
      TripsSyncAPI.newReplaceOp(
        `/contents/${contentListID}/contents/${idx}/title`,
        title
      ),
    ]);
  }

  const onDeleteContent = (idx: number, contentListID: string) => {
    props.tripStateOnUpdate([
      TripsSyncAPI.makeRemoveOp(`/contents/${contentListID}/contents/${idx}`, "")
    ]);
  }

  const onUpdateContentPlace = (idx: number, place: any, contentListID: string) => {
    props.tripStateOnUpdate([
      TripsSyncAPI.newReplaceOp(`/contents/${contentListID}/contents/${idx}/place`, place),
    ]);
  }

  const onUpdateContentNotes = (idx: number, notes: string, contentListID: string) => {
    props.tripStateOnUpdate([
      TripsSyncAPI.newReplaceOp(`/contents/${contentListID}/contents/${idx}/notes`, notes)
    ]);
  }

  const onUpdateContentItineraryDate = async (idx: number, itinListIdx: number, contentListID: string) => {
    const content = _get(props.trip, `contents.${contentListID}.contents[${idx}]`, {}) as Trips.Content;
    const itinList = _get(props.trip, `itinerary[${itinListIdx}]`, {}) as Trips.ItineraryList;
    const itinListCtnt = itinList.contents;
    const listDt = itinList.date as string;

    const ops = [];
    // Update content labels, Format of itinerary dates label:
    // content.labels[LabelContentItineraryDates] = "d1|d2|d3"

    let currentItinDts = _get(content, LabelContentItineraryDatesJSONPath, "")
      .split(LabelContentItineraryDatesDelimeter)
      .filter((dt: string) => !_isEmpty(dt));

    let newItinDts;
    if (currentItinDts.includes(listDt)) {
      // 1. Remove from content label if it exists
      newItinDts = currentItinDts.filter((dt: string) => dt !== listDt);

      // 2. Remove ItineraryContent from ItineraryList
      let itinCtnIdx = _findIndex(itinListCtnt, (ct) => ct.tripContentId === content.id);
      ops.push(TripsSyncAPI.makeRemoveOp(`/itinerary/${itinListIdx}/contents/${itinCtnIdx}`, "",));

      // 3. Remove ItineraryContentRoute from ItineraryList
      ops.push(TripsSyncAPI.newReplaceOp(`/itinerary/${itinListIdx}/routes`, []));
      // if (itinCtnIdx === 0) {
      //   const routeIdx = itinCtnIdx + 1;
      //   if (routeIdx < itinList.routes.length) {
      //     ops.push(TripsSyncAPI.makeRemoveOp(`/itinerary/${itinListIdx}/routes/${routeIdx}`, ""));
      //   }
      // } else if (itinCtnIdx === itinListCtnt.length - 1) {
      //   const routeIdx = itinCtnIdx - 1;
      //   ops.push(TripsSyncAPI.makeRemoveOp(`/itinerary/${itinListIdx}/routes/${routeIdx}`, ""));
      // } else {
      //   ops.push(TripsSyncAPI.makeRemoveOp(`/itinerary/${itinListIdx}/routes/${itinCtnIdx}`, ""));
      //   ops.push(TripsSyncAPI.makeRemoveOp(`/itinerary/${itinListIdx}/routes/${itinCtnIdx - 1}`, ""));
      // }
    } else {
      // 1. Add to content label if its a new date
      newItinDts = _sortBy(currentItinDts.concat([listDt]));

      // 2. Add ItineraryContent to ItineraryList
      const itinCtn: Trips.ItineraryContent = {
        id: uuidv4(),
        tripContentId: content.id,
        tripContentListId: contentListID,
        priceMetadata: {} as any,
        labels: new Map<string, string>(),
      };
      ops.push(TripsSyncAPI.makeAddOp(`/itinerary/${itinListIdx}/contents/-`, itinCtn))


      // 3. Add ItineraryContentRoute to ItineraryList
      if (itinListCtnt.length > 0) {
        const lastItinCtn = _last(itinListCtnt);
        const lastCtnt = _find(
          _get(props.trip, `contents[${lastItinCtn?.tripContentListId}].contents`),
          (ctnt: Trips.Content) => ctnt.id == lastItinCtn?.tripContentId,
        );

        const lastCtntPlaceID = _get(lastCtnt, "place.place_id");
        const ctntPlaceID = _get(content, "place.place_id");

        if (lastCtntPlaceID && ctntPlaceID) {
          const resp = await MapsAPI.directions(lastCtntPlaceID, ctntPlaceID, ModeDriving);
          if (resp.data.routeList.length > 0) {
            ops.push(TripsSyncAPI.makeAddOp(`/itinerary/${itinListIdx}/routes/-`, resp.data.routeList[0]));
          }
        }
      }
    }

    if (currentItinDts.length !== 0) {
      ops.unshift(TripsSyncAPI.newReplaceOp(
        `/contents/${contentListID}/contents/${idx}/labels/${LabelContentItineraryDates}`,
        newItinDts.join(LabelContentItineraryDatesDelimeter)));
    } else {
      ops.unshift(TripsSyncAPI.makeAddOp(
        `/contents/${contentListID}/contents/${idx}/labels/${LabelContentItineraryDates}`,
        newItinDts.join(LabelContentItineraryDatesDelimeter)));
    }

    props.tripStateOnUpdate(ops);
  }


  // Renderers

  const renderHeader = () => {
    return (
      <div className={TripContentSectionCss.HeaderCtn}>
        <div>
          <ToggleChevron
            isHidden={isHidden}
            onClick={() => { setIsHidden(!isHidden) }}
          />
          <span className={TripContentSectionCss.Header}>
            Activities
          </span>
        </div>
        <button
          className={TripContentSectionCss.AddBtn}
          onClick={() => { addContentListBtnOnClick() }}
        >
          +&nbsp;&nbsp;New List&nbsp;
        </button>
      </div>
    );
  }

  const renderContentLists = () => {
    const contentLists = Object.values(_get(props.trip, "contents", {}));
    return contentLists.map((contentList: any) => (
      <div key={contentList.id}>
        <ContentList
          itinerary={props.trip.itinerary}
          contentList={contentList}
          onDeleteList={deleteContentList}
          onUpdateName={contentListUpdateName}
          onAddContent={contentListAddContent}
          onUpdateColorIcon={updateContentListColorIcon}
          onUpdateContentName={onUpdateContentName}
          onDeleteContent={onDeleteContent}
          onUpdateContentPlace={onUpdateContentPlace}
          onUpdateContentNotes={onUpdateContentNotes}
          onUpdateContentItineraryDate={onUpdateContentItineraryDate}
        />
        <hr className={TripContentSectionCss.Hr} />
      </div>
    ));
  }

  return (
    <div className='p-5'>
      {renderHeader()}
      {renderContentLists()}
    </div>
  );

}

export default ContentSection;