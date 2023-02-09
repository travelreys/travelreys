import React, {
  FC,
  useEffect,
  useState,
} from 'react';
import _get from "lodash/get";
import _sortBy from "lodash/sortBy";
import _isEmpty from "lodash/isEmpty";
import { v4 as uuidv4 } from 'uuid';
import { useDebounce } from 'usehooks-ts';
import {
  CheckIcon,
  ChevronDownIcon,
  ChevronUpIcon,
  MapPinIcon,
  PlusIcon,
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
import MapsAPI, { placeFields } from '../../apis/maps';
import { Trips } from '../../apis/types';
import { useMap } from '../../context/maps-context';
import {
  TripContentSectionCss,
  TripContentListCss,
  TripContentCss
} from '../../styles/global';
import { areYMDEqual, parseTimeFromZ, printTime } from '../../utils/dates';


// TripContent
interface TripContentProps {
  content: Trips.Content
  contentListID: string
  contentIdx: number
  itinerary: Array<Trips.ItineraryList>
  tripStateOnUpdate: any
}

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

  // Event Handlers - Title
  const titleInputOnBlur = () => {
    const ops = [
      TripsSyncAPI.makeReplaceOp(
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
        TripsSyncAPI.makeReplaceOp(
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
      dispatch({type:"setSelectedPlace", value: props.content.place})
      const event = new CustomEvent('marker_click', {
        bubbles: false,
        cancelable: false,
        detail: props.content.place,
      });
      document.getElementById("map")!.dispatchEvent(event)
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
    ops.push(TripsSyncAPI.makeReplaceOp(
      `/contents/${props.contentListID}/contents/${props.contentIdx}/notes`,
      content
    ));
    props.tripStateOnUpdate(ops);
  }

  // Event Handlers - Itinerary
  const itineraryBtnOnClick = (l: Trips.ItineraryList, lIdx: number) => {
    const ops = [];
    const _dt = l.date as string;

    // Update content labels
    let currentDts = _get(props.content, "labels.itinerary|dates", "")
      .split("|")
      .filter((dt: string) => !_isEmpty(dt));

    let dts;
    if (currentDts.includes(_dt)) {
      dts = currentDts.filter((dt: string) => dt !== _dt)
    } else {
      currentDts.push(_dt);
      dts = _sortBy(currentDts);
    }

    if (_get(props.content, "labels.itinerary|dates")) {
      ops.push(TripsSyncAPI.makeReplaceOp(
        `/contents/${props.contentListID}/contents/${props.contentIdx}/labels/itinerary|dates`,
        dts.join("|")));
    } else {
      ops.push(TripsSyncAPI.makeAddOp(
        `/contents/${props.contentListID}/contents/${props.contentIdx}/labels/itinerary|dates`,
        dts.join("|")));
    }

    // Update itinerary
    currentDts = _get(props.content, "labels.itinerary|dates", "").split("|");
    if (currentDts.includes(_dt)) {
      let itinCtnIdx;
      l.contents.forEach((ct: Trips.ItineraryContent, idx: number) => {
        if (ct.tripContentId === props.content.id) {
          itinCtnIdx = idx;
        }
      });
      ops.push(TripsSyncAPI.makeRemoveOp(
        `/itinerary/${lIdx}/contents/${itinCtnIdx}`,
        ""));
    } else {
      const itinCtn: Trips.ItineraryContent = {
        tripContentId: props.content.id,
        tripContentListId: props.contentListID,
        price: {} as any,
        labels: new Map<string,string>(),
      };
      ops.push(TripsSyncAPI.makeAddOp(
        `/itinerary/${lIdx}/contents/-`,
        itinCtn))
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
        <TrashIcon className='h-4 w-4 mr-2' />Delete
      </button>
    ];
    const menu = (
      <EllipsisHorizontalCircleIcon className='h-4 w-4 mt-1' />
    );
    return <Dropdown menu={menu} opts={opts} />
  }

  const renderItineraryDropdown = () => {
    const dates = _get(props.content, "labels.itinerary|dates", "")
      .split("|").filter((dt: string) => !_isEmpty(dt));

    const opts = props.itinerary.map((l: Trips.ItineraryList, idx: number) => {
      const isAdded = dates.includes(l.date as string);

      return (
        <button
          type='button'
          className={TripContentCss.ItineraryDateBtn}
          onClick={() => {itineraryBtnOnClick(l, idx)}}
        >
          {printTime(parseTimeFromZ(l.date as string), "eee, do MMMM") }
          { isAdded ? <CheckIcon className='w-4 h-4' /> : null}
        </button>
      );
    });

    const datesBadge = dates
      .map((dt: string) => (
        <span className="bg-indigo-100 text-indigo-800 text-xs font-medium mr-2 px-2.5 py-0.5 rounded">
          {printTime(parseTimeFromZ(dt), "MMM/dd")}
        </span>
      ));

    const emptyBtn = (
      <span className='text-xs text-gray-800 font-bold bg-indigo-200 rounded-full px-2 py-1 hover:bg-indigo-400'>
        Add to Itinerary
      </span>
    );
    const menu = (<div>{dates.length === 0 ? emptyBtn : datesBadge}</div>);

    return <Dropdown menu={menu} opts={opts} />
  }


  const renderTitleInput = () => {
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
          placeholder={`name, address...`}
          className="p-0 mb-1 text-sm text-gray-600 bg-transparent placeholder:text-gray-400 rounded border-0 hover:border-0 focus:ring-0 duration-400"
        />
      );
    } else {
      const addr = _get(props.content, "place.name", "");
      if (_isEmpty(addr)) {
        placeNode = (
          <button type='button' onClick={() => {setIsAddingPlace(true)}}>
            Click here to add a location...
          </button>
        );
      } else {
        placeNode = (
          <button type='button' onClick={placeOnClick}>
            {addr}
          </button>
        );
      }
    }

    return (
      <p className='text-slate-600 text-sm flex items-center mb-1 hover:text-indigo-500'>
        <MapPinIcon className='h-4 w-4 mr-1'/>
        {placeNode}
      </p>
    );
  }

  const renderPlacesAutocomplete = () => {
    if (_isEmpty(predictions)) {
      return (<></>);
    }
    return (
      <div className={TripContentCss.AutocompleteCtn}>
        {predictions.map((pre: any) => (
          <div
            className={TripContentCss.PredictionWrapper}
            key={pre.place_id}
            onClick={() => {predictionOnSelect(pre.place_id)}}
          >
            <div className='p-1 group-hover:text-indigo-500'>
              <MapPinIcon className='h-6 w-6' />
            </div>
            <div className='ml-1'>
              <p className='text-slate-900 group-hover:text-indigo-500 text-sm font-medium'>
                {_get(pre, "structured_formatting.main_text", "")}
              </p>
              <p className="text-slate-400 group-hover:text-indigo-500 text-xs">
                {_get(pre, "structured_formatting.secondary_text", "")}
              </p>
            </div>
          </div>
        ))}
      </div>
    );
  }

  const renderWebsite = () => {
    return _isEmpty(_get(props.content, "place.website", "")) ? null
    : <a
        className='flex items-center mb-1'
        href={_get(props.content, "place.website")}
        target="_blank"
      >
        <GlobeAltIcon className='h-4 w-4' />&nbsp;
        <span className={TripContentCss.WebsiteTxt}>Website</span>
      </a>
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
      {renderTitleInput()}
      {renderPlace()}
      {renderPlacesAutocomplete()}
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


// TripContentList

interface TripContentListProps {
  itinerary: any
  contentList: Trips.ContentList
  tripStateOnUpdate: any
}

const TripContentList: FC<TripContentListProps> = (props: TripContentListProps) => {

  const [name, setName] = useState<string>();
  const [newContentTitle, setNewContentTitle] = useState("");
  const [isHidden, setIsHidden] = useState<boolean>(false);

  useEffect(() => {
    setName(props.contentList.name);
  }, [props.contentList])

  // Event Handlers - Content List Name

  const nameOnBlur = () => {
    const ops = [];
    ops.push(TripsSyncAPI.makeReplaceOp(`/contents/${props.contentList.id}/name`, name))
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
  const renderHiddenToggle = () => {
    return (
      <button
        type="button"
        className={TripContentSectionCss.ToggleBtn}
        onClick={() => {setIsHidden(!isHidden)}}
      >
      {isHidden ? <ChevronUpIcon className='h-4 w-4' />
        : <ChevronDownIcon className='h-4 w-4'/>}
      </button>
    );
  }

  const renderTripContent = () => {
    if (isHidden) {
      return null;
    }
    const contents = _get(props.contentList, "contents", []);
    return (
      <div>
        {contents.map((content: any, idx: number) => (
          <TripContent
            key={idx}
            itinerary={props.itinerary}
            content={content}
            contentListID={props.contentList.id}
            contentIdx={idx}
            tripStateOnUpdate={props.tripStateOnUpdate}
          />
        ))}
      </div>
    );
  }

  const renderAddNewContent = () => {
    if (isHidden) {
      return null;
    }

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
        {renderHiddenToggle()}
        <input
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
          onBlur={nameOnBlur}
          placeholder={`Add a title (e.g, "Food to try")`}
          className={TripContentListCss.NameInput}
        />
      </div>
      {renderTripContent()}
      {renderAddNewContent()}
    </div>
  );
}

// TripContentSection


interface TripContentSectionProps {
  trip: any
  tripStateOnUpdate: any
}

const TripContentSection: FC<TripContentSectionProps> = (props: TripContentSectionProps) => {

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

  const renderHiddenToggle = () => {
    return (
      <button
        type="button"
        className={TripContentSectionCss.ToggleBtn}
        onClick={() => {setIsHidden(!isHidden)}}
      >
      {isHidden ? <ChevronUpIcon className='h-4 w-4' />
        : <ChevronDownIcon className='h-4 w-4'/>}
      </button>
    );
  }

  const renderContentLists = () => {
    if (isHidden) {
      return null;
    }
    const contentLists = Object.values(_get(props.trip, "contents", {}));
    return contentLists.map((contentList: any) => {
      return (
      <div key={contentList.id}>
        <TripContentList
          itinerary={props.trip.itinerary}
          contentList={contentList}
          tripStateOnUpdate={props.tripStateOnUpdate}
        />
        <hr className={TripContentSectionCss.Hr} />
      </div>
      );
    });
  }

  return (
    <div className='p-5'>
      <div className={TripContentSectionCss.HeaderCtn}>
        <div>
          {renderHiddenToggle()}
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
      {renderContentLists()}
    </div>
  );

}

export default TripContentSection;