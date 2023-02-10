import React, {
  FC,
  useEffect,
  useState,
} from 'react';
import _get from "lodash/get";
import _sortBy from "lodash/sortBy";
import _isEmpty from "lodash/isEmpty";
import _findIndex from "lodash/findIndex";
import { v4 as uuidv4 } from 'uuid';
import { useDebounce } from 'usehooks-ts';
import {
  CheckIcon,
  MapPinIcon,
  PlusIcon,
  TrashIcon
} from '@heroicons/react/24/solid';
import {
  EllipsisHorizontalCircleIcon,
  GlobeAltIcon,
} from '@heroicons/react/24/outline';

import Dropdown from '../Dropdown';
import NotesEditor from '../NotesEditor';
import PlacePicturesCarousel from './PlacePicturesCarousel';
import ToggleChevron from '../ToggleChevron';

import TripsSyncAPI from '../../apis/tripsSync';
import MapsAPI, { placeFields } from '../../apis/maps';
import {
  Trips,
  LabelContentItineraryDates,
  LabelContentItineraryDatesJSONPath,
  LabelContentItineraryDatesDelimeter,
} from '../../apis/trips';
import { ActionNameSetSelectedPlace, useMap } from '../../context/maps-context';
import {
  CommonCss,
  TripContentCss,
  TripContentListCss,
  TripContentSectionCss,
} from '../../styles/global';
import { parseISO, printFmt } from '../../utils/dates';
import PlaceAutocomplete from '../maps/PlaceAutocomplete';
import { EventMarkerClickName, newEventMarkerClick } from '../maps/common';


/////////////
// Content //
/////////////

interface TripContentProps {
  content: Trips.Content
  contentListID: string
  contentIdx: number
  itinerary: Array<Trips.ItineraryList>
  tripStateOnUpdate: any
}

const ItineraryDateFmt = "eee, do MMMM";
const ItineraryBadgeDateFmt = "MMM/dd";


const TripContent: FC<TripContentProps> = (props: TripContentProps) => {
  const [title, setTitle] = useState<string>();
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
    const ops = [
      TripsSyncAPI.newReplaceOp(
        `/contents/${props.contentListID}/contents/${props.contentIdx}/title`,
        title
      ),
    ];
    props.tripStateOnUpdate(ops);
  }

  const deleteBtnOnClick = () => {
    const ops = [];
    ops.push(
      TripsSyncAPI.makeRemoveOp(
        `/contents/${props.contentListID}/contents/${props.contentIdx}`,
        "")
    );
    props.tripStateOnUpdate(ops);
  }

  // Event Handlers - Places
  const predictionOnSelect = (placeID: string) => {
    MapsAPI.placeDetails(placeID, placeFields, sessionToken)
    .then((res) => {
      setPredictions([]);
      const place = _get(res, "data.place", {});
      const ops = [
        TripsSyncAPI.newReplaceOp(
          `/contents/${props.contentListID}/contents/${props.contentIdx}/place`,
          place),
      ];
      props.tripStateOnUpdate(ops);
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
      document.getElementById("map")?.dispatchEvent(event)
      return;
    }

    if (e.detail == 2) {
      setIsAddingPlace(true);
      return;
    }
  }

  // Event Handlers - Notes
  const notesOnChange = (content: string) => {
    const ops = [];
    ops.push(TripsSyncAPI.newReplaceOp(
      `/contents/${props.contentListID}/contents/${props.contentIdx}/notes`,
      content
    ));
    props.tripStateOnUpdate(ops);
  }

  // Event Handlers - Itinerary
  const itineraryBtnOnClick = (l: Trips.ItineraryList, itinListIdx: number) => {
    const ops = [];
    const listDt = l.date as string;

    // Update content labels, Format of itinerary dates label:
    // content.labels[LabelContentItineraryDates] = "d1|d2|d3"

    let currentItinDts = _get(props.content, LabelContentItineraryDatesJSONPath, "")
      .split(LabelContentItineraryDatesDelimeter)
      .filter((dt: string) => !_isEmpty(dt));

    let dts;
    if (currentItinDts.includes(listDt)) {
      // Remove if already exists
      dts = currentItinDts.filter((dt: string) => dt !== listDt)

      let itinCtnIdx = _findIndex(
        l.contents,
        (ct) => ct.tripContentId === props.content.id,
      );
      ops.push(TripsSyncAPI.makeRemoveOp(
        `/itinerary/${itinListIdx}/contents/${itinCtnIdx}`,
        "",
      ));

    } else {
      // Add if not exists
      dts = _sortBy(currentItinDts.concat([listDt]));

      const itinCtn: Trips.ItineraryContent = {
        id: uuidv4(),
        tripContentId: props.content.id,
        tripContentListId: props.contentListID,
        priceMetadata: {} as any,
        labels: new Map<string,string>(),
      };
      ops.push(TripsSyncAPI.makeAddOp(
        `/itinerary/${itinListIdx}/contents/-`,
        itinCtn))
    }

    if (currentItinDts) {
      ops.push(TripsSyncAPI.newReplaceOp(
        `/contents/${props.contentListID}/contents/${props.contentIdx}/labels/${LabelContentItineraryDates}`,
        dts.join(LabelContentItineraryDatesDelimeter)));
    } else {
      ops.push(TripsSyncAPI.makeAddOp(
        `/contents/${props.contentListID}/contents/${props.contentIdx}/labels/${LabelContentItineraryDates}`,
        dts.join(LabelContentItineraryDatesDelimeter)));
    }

    props.tripStateOnUpdate(ops);
  }

  // Renderers

  const renderSettingsDropdown = () => {
    const opts = [
      <button
        type='button'
        className={TripContentCss.DeleteBtn}
        onClick={deleteBtnOnClick}
      >
        <TrashIcon className={CommonCss.LeftIcon}/>
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
        onClick={() => {itineraryBtnOnClick(l, idx)}}
      >
        {printFmt(parseISO(l.date as string), ItineraryDateFmt) }
        { dates.includes(l.date as string) ?
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
            (e) => {setIsAddingPlace(true)} :
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
        <MapPinIcon className='h-4 w-4 mr-1'/>
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
        <GlobeAltIcon className={CommonCss.LeftIcon}/>
        <span className={TripContentCss.WebsiteTxt}>Website</span>
      </a>
    );
  }

  const renderPlacePicturesCarousel = () => {
    const photos = _get(props.content, "place.photos", []);
    if (_isEmpty(photos)) {
      return null;
    }
    return <PlacePicturesCarousel photos={photos}/>
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


////////////////////
// ContentSection //
////////////////////


interface ContentListProps {
  itinerary: any
  contentList: Trips.ContentList
  tripStateOnUpdate: any
}

const ContentList: FC<ContentListProps> = (props: ContentListProps) => {

  const [name, setName] = useState<string>();
  const [newContentTitle, setNewContentTitle] = useState("");
  const [isHidden, setIsHidden] = useState<boolean>(false);

  useEffect(() => {
    setName(props.contentList.name);
  }, [props.contentList])

  // Event Handlers - Content List Name

  const nameOnBlur = () => {
    const ops = [];
    ops.push(TripsSyncAPI.newReplaceOp(`/contents/${props.contentList.id}/name`, name))
    props.tripStateOnUpdate(ops);
  }

  // Event Handlers  - New Content

  const newContentBtnOnClick = () => {
    const content: Trips.Content = {
      id: uuidv4(),
      title: newContentTitle,
      notes: "",
      place: {},
      labels: new Map<string,string>(),
      comments: [],
    }
    const ops = [
      TripsSyncAPI.makeJSONPatchOp(
        "add", `/contents/${props.contentList.id}/contents/-`, content),
    ]
    props.tripStateOnUpdate(ops)
    setNewContentTitle("");
  }

  // Renderers

  const renderTripContent = () => {
    return _get(props.contentList, "contents", [])
      .map((content: any, idx: number) => (
        <TripContent
          key={idx}
          itinerary={props.itinerary}
          content={content}
          contentListID={props.contentList.id}
          contentIdx={idx}
          tripStateOnUpdate={props.tripStateOnUpdate}
        />
      ));
  }

  const renderAddNewContent = () => {
    return (
      <div className={TripContentListCss.NewContentCtn}>
        <input
          type="text"
          value={newContentTitle}
          onChange={(e) => {setNewContentTitle(e.target.value)}}
          placeholder="Add an activity..."
          className={TripContentListCss.NewContentInput}
        />
        <button
          onClick={() => {newContentBtnOnClick()}}
          className={TripContentListCss.NewContentBtn}
        >
          <PlusIcon className='h-4 w-4 stroke-2'/>
        </button>
      </div>
    );
  }

  return (
    <div className={TripContentListCss.Ctn}>
      <div className='flex'>
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
      { isHidden ? null :
        <>
          {renderTripContent()}
          {renderAddNewContent()}
        </>
      }
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

  // Event Handlers
  const addBtnOnClick = () => {
    let list: Trips.ContentList = {
      id: uuidv4(),
      contents: new Array<Trips.Content>(),
    }
    const ops = [];
    ops.push(TripsSyncAPI.makeAddOp(`/contents/${list.id}`, list))
    props.tripStateOnUpdate(ops);
  }

  // Renderers

  const renderHeader = () => {
    return (
      <div className={TripContentSectionCss.HeaderCtn}>
        <div>
          <ToggleChevron
            isHidden={isHidden}
            onClick={() => {setIsHidden(!isHidden)}}
          />
          <span className={TripContentSectionCss.Header}>
            Activities
          </span>
        </div>
        <button
          className={TripContentSectionCss.AddBtn}
          onClick={() => {addBtnOnClick()}}
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
          tripStateOnUpdate={props.tripStateOnUpdate}
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