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

import TripsSyncAPI from '../../apis/tripsSync';
import MapsAPI from '../../apis/maps';
import PlacePicturesCarousel from './PlacePicturesCarousel';
import { useMap } from '../../context/maps-context';


// TripContent
interface TripContentProps {
  contentListID: string
  contentIdx: number
  content: any
  tripStateOnUpdate: any
}

const TripContent: FC<TripContentProps> = (props: TripContentProps) => {

  const {dispatch} = useMap();

  const [title, setTitle] = useState(props.content.title)
  const [isAddingPlace, setIsAddingPlace] = useState(false);
  const [searchPlaceQuery, setSearchPlaceQuery] = useState("");
  const [predictions, setPredictions] = useState([] as any);
  const [sessionToken, setSessionToken] = useState("");

  const debouncedValue = useDebounce<string>(searchPlaceQuery, 500);

  useEffect(() => {
    setTitle(props.content.title)
  }, [props.content])

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
    // setPredictions(preds.predictions);
  }

  const getPlaceDetails = (placeID: string) => {
    const fields = [
      "address_component", "adr_address", "business_status", "formatted_address", "geometry",  "name", "photos", "place_id", "types", "utc_offset",
      "opening_hours", "formatted_phone_number", "international_phone_number", "website",
    ];
    MapsAPI.placeDetails(placeID, fields, sessionToken)
    .then((res) => {
      setPredictions([]);
      const place = _get(res, "data.place", {});
      const ops = [
        TripsSyncAPI.makeJSONPatchOp(
          "replace",
          `/contents/${props.contentListID}/contents/${props.contentIdx}/place`,
          place),
      ];
      props.tripStateOnUpdate(ops);
    })
    .finally(() => {
      setSessionToken("");
    });
  }

  // Event Handlers - Title
  const titleInputOnBlur = () => {
    const ops = [
      TripsSyncAPI.makeJSONPatchOp(
        "replace",
        `/contents/${props.contentListID}/contents/${props.contentIdx}/title`,
        title),
    ];
    props.tripStateOnUpdate(ops);
  }

  // Event Handlers - Places
  useEffect(() => {
    if (!_isEmpty(debouncedValue)) {
      autocomplete(debouncedValue);
    }
  }, [debouncedValue])

  const predictionOnSelect = (placeID: string) => {
    getPlaceDetails(placeID);
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
        placeholder={`Add a name (e.g, "Scenic Spots")`}
        className="p-0 mb-1 font-bold text-gray-800 bg-transparent placeholder:text-gray-400 rounded border-0 hover:border-0 focus:ring-0 duration-400"
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
      <div className='p-1 bg-white absolute left-0 w-full border border-slate-200 rounded-lg'>
        {predictions.map((pre: any) => (
          <div
            className='flex items-center mb-4 cursor-pointer group'
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
    <div className='bg-slate-50 rounded-lg shadow-xs mb-4 p-4 relative' >
      {renderTitleInput()}
      {renderPlace()}
      {renderPlacesAutocomplete()}
      {renderPlacePicturesCarousel()}
    </div>
  );
}


// TripContentList

interface TripContentListProps {
  contentList: any
  tripStateOnUpdate: any
}

const TripContentList: FC<TripContentListProps> = (props: TripContentListProps) => {

  const [name, setName] = useState(props.contentList.name);
  const [newContentTitle, setNewContentTitle] = useState("");

  useEffect(() => {
    setName(props.contentList.name);
  }, [props.contentList])


  // Event Handlers - Content List Name

  const nameOnBlur = () => {
    const ops = [
      TripsSyncAPI.makeJSONPatchOp(
        "replace", `/contents/${props.contentList.id}/name`, name),
    ];
    props.tripStateOnUpdate(ops);
  }

  // Event Handlers  - New Content

  const newContentBtnOnClick = () => {
    const newContent = {
      id: uuidv4(),
      title: newContentTitle,
      comments: [],
      labels: {},
    }
    const ops = [
      TripsSyncAPI.makeJSONPatchOp(
        "add", `/contents/${props.contentList.id}/contents/-`, newContent),
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

  const renderAddNewTripContent = () => {
    return (
      <div className='flex my-4 w-full'>
        <input
          type="text"
          value={newContentTitle}
          onChange={(e) => {setNewContentTitle(e.target.value)}}
          placeholder={`Add an activity...`}
          className="flex-1 mr-1 text-md sm:text-md font-bold text-gray-800 placeholder:font-normal placeholder:text-gray-300 placeholder:italic rounded border-0 bg-gray-100 hover:border-0 focus:ring-0"
        />
        <button
          className='text-green-600 w-1/12 hover:bg-green-50 rounded-lg text-sm font-bold inline-flex justify-around items-center'
          onClick={() => {newContentBtnOnClick()}}
        >
          <PlusIcon className='h-4 w-4 stroke-2'/>
        </button>
      </div>
    );
  }

  return (
    <div className='rounded-lg shadow-xs mb-4' >
      <input
        type="text"
        value={name}
        onChange={(e) => setName(e.target.value)}
        onBlur={nameOnBlur}
        placeholder={`Add a title (e.g, "Scenic Spots")`}
        className="p-0 text-xl mb-1 sm:text-2xl font-bold text-gray-800 placeholder:text-gray-400 rounded border-0 hover:bg-gray-300 hover:border-0 focus:ring-0 focus:p-1 duration-500"
      />
      {renderTripContent()}
      {renderAddNewTripContent()}
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
  const addNewListOnClick = () => {
    let contentList = { id: uuidv4(), contents: [] }
    const ops = [
      TripsSyncAPI.makeJSONPatchOp(
        "add", `/contents/${contentList.id}`, contentList)
    ];
    props.tripStateOnUpdate(ops);
  }

  // Renderers

  const renderActivitiesLists = () => {
    const contentLists = Object.values(_get(props.trip, "contents") || {});
    return contentLists.map((contentList: any) => {
      return (
      <>
        <TripContentList
          key={contentList.id}
          contentList={contentList}
          tripStateOnUpdate={props.tripStateOnUpdate}
        />
        <hr className='w-48 h-1 m-5 mx-auto bg-gray-300 border-0 rounded'/>
      </>
      );
    });
  }

  return (
    <div className='p-5'>
      <div className="flex justify-between mb-4">
        <h3 className='text-2xl sm:text-3xl font-bold text-slate-700'>
          Activities
        </h3>
        <button
          className='text-white py-2 px-4 bg-indigo-500 rounded-lg text-sm font-semibold'
          onClick={() => {addNewListOnClick()}}
        >
          +&nbsp;&nbsp;New List&nbsp;
        </button>
      </div>
      <div>
        {renderActivitiesLists()}
      </div>
    </div>
  );

}

export default TripContentSection;