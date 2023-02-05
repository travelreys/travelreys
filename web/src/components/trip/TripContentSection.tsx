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
import { MapPinIcon, PlusIcon } from '@heroicons/react/24/outline'

import PlacePicturesCarousel from './PlacePicturesCarousel';
import TripsSyncAPI from '../../apis/tripsSync';
import MapsAPI, { placeFields } from '../../apis/maps';
import { Trips } from '../../apis/types';

import { useMap } from '../../context/maps-context';
import { TripContentSectionCss, TripContentListCss, TripContentCss } from '../../styles/global';


// TripContent
interface TripContentProps {
  contentListID: string
  contentIdx: number
  content: Trips.Content
  tripStateOnUpdate: any
}

const TripContent: FC<TripContentProps> = (props: TripContentProps) => {
  const [title, setTitle] = useState<string>();
  const [isAddingPlace, setIsAddingPlace] = useState<boolean>(false);
  const [searchPlaceQuery, setSearchPlaceQuery] = useState<string>("");
  const [predictions, setPredictions] = useState([] as any);
  const [sessionToken, setSessionToken] = useState<string>("");

  const debouncedValue = useDebounce<string>(searchPlaceQuery, 500);

  const {dispatch} = useMap();

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

  const placeOnClick = () => {
    dispatch({type:"setSelectedPlace", value: props.content.place})
  }

  // Renderers
  const renderTitleInput = () => {
    return (
      <input
        type="text"
        value={title}
        onChange={(e) => setTitle(e.target.value)}
        onBlur={titleInputOnBlur}
        placeholder="Add a name"
        className={TripContentCss.TitleInput}
      />
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
      const addr = _get(props.content, "place.formatted_address", "");
      if (_isEmpty(addr)) {
        placeNode = (
          <button type='button' onClick={() => {setIsAddingPlace(true)}}>
            Click here to add a location...
          </button>
        );
      } else {
        placeNode = (
          <button type='button'
            className='hover:text-indigo-500'
            onClick={placeOnClick}
          >
            {addr}
          </button>
        );
      }
    }

    return (
      <p className='text-slate-600 text-sm flex items-center mb-2'>
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

  const renderPlacePicturesCarousel = () => {
    const photos = _get(props.content, "place.photos", []);
    if (_isEmpty(photos)) {
      return null;
    }
    return <PlacePicturesCarousel photos={photos}/>
  }

  const renderNotes = () => {
    return (null);
  }

  return (
    <div className={TripContentCss.Ctn}>
      {renderTitleInput()}
      {renderPlace()}
      {renderPlacesAutocomplete()}
      {renderPlacePicturesCarousel()}
    </div>
  );
}


// TripContentList

interface TripContentListProps {
  contentList: Trips.ContentList
  tripStateOnUpdate: any
}

const TripContentList: FC<TripContentListProps> = (props: TripContentListProps) => {

  const [name, setName] = useState<string>();
  const [newContentTitle, setNewContentTitle] = useState("");

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

  const renderTripContent = () => {
    const contents = _get(props.contentList, "contents", []);
    return (
      <div>
        {contents.map((content: any, idx: number) => (
          <TripContent
            key={idx}
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
      <input
        type="text"
        value={name}
        onChange={(e) => setName(e.target.value)}
        onBlur={nameOnBlur}
        placeholder={`Add a title (e.g, "Food to try")`}
        className={TripContentListCss.NameInput}
      />
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

  const renderContentLists = () => {
    const contentLists = Object.values(_get(props.trip, "contents", {}));
    return contentLists.map((contentList: any) => {
      return (
      <div key={contentList.id}>
        <TripContentList
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
        <h3 className={TripContentSectionCss.Header}>
          Activities
        </h3>
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