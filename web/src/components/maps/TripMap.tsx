import React, {
  FC,
  useEffect,
  useState,
  useRef
} from 'react';
import _get from "lodash/get";
import _flatten from "lodash/flatten";
import _isEmpty from "lodash/isEmpty";
import { Wrapper, Status } from "@googlemaps/react-wrapper";
import {
  ClockIcon,
  StarIcon,
  MapPinIcon,
  PhoneIcon,
  GlobeAltIcon,
} from '@heroicons/react/24/solid'
import {
  XMarkIcon
} from '@heroicons/react/24/outline'

import MapsAPI, {
  placeAtmosphereFields,
  PLACE_IMAGE_APIKEY
} from '../../apis/maps';

import Spinner from '../Spinner';
import { makeActivityPin, makeHotelPin } from './GMapsPinIcon';
import { useMap } from '../../context/maps-context';
import GoogleIcon from '../icons/GoogleIcon';


const defaultMapCenter = { lat: 33.3960897, lng: 126.264522 }
const defaultMapOpts = {
  center: defaultMapCenter,
  zoom: 10,
  mapTypeControl: false,
  streetViewControl: false,
  fullscreenControl: false,
  rotateControl: false,
  keyboardShortcuts: true,
  gestureHandling: "greedy"
}

interface PlaceMarker {
  elem: HTMLElement
  place: any
}


interface InnerMapProps {
  markers: Array<PlaceMarker>
  width: any
}

const InnerMap: FC<InnerMapProps> = (props: InnerMapProps) => {
  const ref = useRef() as any;
  const map = useRef() as any;
  const {state, dispatch} = useMap();

  const currentPopups = useRef([]) as any;
  const currentMapCenter = useRef(null) as any;

  class Popup extends google.maps.OverlayView {
    position: google.maps.LatLng;
    containerDiv: HTMLDivElement;

    constructor(position: google.maps.LatLng, content: HTMLElement) {
      super();
      this.position = position;

      // This zero-height div is positioned at the bottom of the bubble.
      const bubbleAnchor = document.createElement("div");

      // bubbleAnchor.classList.add("popup-bubble-anchor");
      bubbleAnchor.appendChild(content);

      // This zero-height div is positioned at the bottom of the tip.
      this.containerDiv = document.createElement("div");
      this.containerDiv.classList.add("popup-container");
      this.containerDiv.appendChild(content);

      // Optionally stop clicks, etc., from bubbling up to the map.
      Popup.preventMapHitsAndGesturesFrom(this.containerDiv);
    }

    /** Called when the popup is added to the map. */
    onAdd() {
      this.getPanes()!.floatPane.appendChild(this.containerDiv);
      this.getPanes()!.overlayMouseTarget.appendChild(this.containerDiv);

      // set this as locally scoped var so event does not get confused
      var me = this;

      // Add a listener - we'll accept clicks anywhere on this div, but you may want
      // to validate the click i.e. verify it occurred in some portion of your overlay.
      this.containerDiv.addEventListener('click', function() {
        google.maps.event.trigger(me, 'click');
      });
    }

    /** Called when the popup is removed from the map. */
    onRemove() {
      if (this.containerDiv.parentElement) {
        this.containerDiv.parentElement.removeChild(this.containerDiv);
      }
    }

    /** Called each frame when the popup needs to draw itself. */
    draw() {
      const divPosition = this.getProjection().fromLatLngToDivPixel(
        this.position
      )!;

      // Hide the popup when it is far out of view.
      const display =
        Math.abs(divPosition.x) < 4000 && Math.abs(divPosition.y) < 4000
          ? "block"
          : "none";

      if (display === "block") {
        this.containerDiv.style.left = divPosition.x + "px";
        this.containerDiv.style.top = divPosition.y + "px";
      }

      if (this.containerDiv.style.display !== display) {
        this.containerDiv.style.display = display;
      }
    }
  }

  useEffect(() => {
    map.current = new window.google.maps.Map(ref.current, defaultMapOpts)
    map.current.addListener("center_changed", () => {
      currentMapCenter.current = map.current.getCenter();
    });
    ref.current.addEventListener("marker_click", (e: any) => {
      const center = _get(e.detail, "geometry.location", defaultMapCenter);
      map.current.setCenter(center);
    })
  }, [])

  useEffect(() => {
    if (!_isEmpty(state.selectedPlace)) {
      const center = _get(state.selectedPlace, "geometry.location", {});
      map.current.setCenter(center);
    }
  }, [state.selectedPlace])

  useEffect(() => {
    // Clear all markers from the previous render
    currentPopups.current.forEach((pp: any) => {
      pp.setMap(null);
    })
    currentPopups.current = [];

    // Make new markers
    const bounds = new google.maps.LatLngBounds();
    props.markers.forEach((marker: PlaceMarker) => {
      const latlng = _get(marker, "place.geometry.location") as any;
      const popup = new Popup(latlng, marker.elem);
      popup.setMap(map.current);
      popup.addListener("click", () => {
        map.current.setCenter(popup.position);
        currentMapCenter.current = popup.position;
        dispatch({type: "setSelectedPlace", value: marker.place})
      })

      currentPopups.current.push(popup);
      bounds.extend(popup.position);
    });

    if (_isEmpty(currentMapCenter.current)) {
      map.current.setCenter(bounds.getCenter());
      currentMapCenter.current = bounds.getCenter();
    }
    // map.current.fitBounds(bounds);
  }, [props.markers])

  return (
    <div ref={ref}
      id="map"
      className='h-full'
      style={{width: props.width}}
    />
  );
}

interface TripMapComponentProps {
  trip: any
  width: any
}

const TripMap: FC<TripMapComponentProps> = (props: TripMapComponentProps) => {

  const {state} = useMap();
  const [placeDetails, setPlaceDetails] = useState(null) as any;

  // API
  useEffect(() => {
    if (state.selectedPlace === null) {
      return;
    }
    const placeID = state.selectedPlace.place_id;
    MapsAPI.placeDetails(placeID, placeAtmosphereFields, "")
    .then((res) => {
      setPlaceDetails(_get(res, "data.place", null));
    })
  }, [state.selectedPlace]);

  // Map Markers
  const lodgingToMapMarkers = () => {
    const lodgings = _get(props.trip, "lodgings", {});
    return Object.values(lodgings).map((lodge: any) => ({
      elem: makeHotelPin(lodge.place.name),
      place: lodge.place
    }));
  }

  const contentToMapMarkers = () => {
    return _flatten(
      Object.values(_get(props.trip, "contents", {}))
      .map((list: any) => list.contents)
    )
    .filter((ct: any) => !_isEmpty(ct.place))
    .map((ct: any) => ({
      elem: makeActivityPin(ct.place.name),
      place: ct.place
    }));
  }

  const makeMarkers = () => {
    let markers = [] as any;
    markers = markers.concat(lodgingToMapMarkers());
    markers = markers.concat(contentToMapMarkers());
    return markers;
  }

  // Renderer
  const renderMap = (status: Status): React.ReactElement => {
    if (status === Status.FAILURE) return <Spinner />;
    return <Spinner />;
  };

  const renderPlaceDetailsCard = () => {
    if (placeDetails === null) {
      return null;
    }

    const renderHeader = () => {
      return (
        <p className='font-bold text-lg flex justify-between items-center'>
          {placeDetails.name}
          <button type="button" onClick={() => {setPlaceDetails(null)}}>
            <XMarkIcon className='h-6 w-6' />
          </button>
        </p>
      );
    }

    const renderRatings = () => {
      return (
        <p className='text-yellow-500 flex items-center mb-1'>
          <StarIcon className="h-4 w-4" />&nbsp;
          {placeDetails.rating}&nbsp;&nbsp;
          <span className='text-gray-600'>
            ({placeDetails.user_ratings_total})
          </span>
          &nbsp;&nbsp;
          <GoogleIcon className="h-4 w-4 mt-1" />
        </p>
      );
    }

    const renderPhone = () => {
      return placeDetails.international_phone_number ?
      <p className='text-gray-600 flex items-center mb-1'>
        <PhoneIcon className='h-4 w-4' />&nbsp;
        {placeDetails.international_phone_number}
      </p>
      : null
    }

    const renderWebsite = () => {
      return placeDetails.website ?
      <a
        href={placeDetails.website}
        target="_blank"
        className='text-indigo-500 flex items-center mb-1'>
        <GlobeAltIcon className='h-4 w-4 text-gray-600 mr-2' />
        Website
      </a>
      : null
    }

    const renderOpeningHours = () => {
      const weekdayTexts = _get(placeDetails, "opening_hours.weekday_text", []);
      if (_isEmpty(weekdayTexts)) {
        return null;
      }
      const text = weekdayTexts.map((txt: string, idx: number) =>
        (<p key={idx} className="text-slate-600 ml-6">{txt}</p>)
      )
      return (
        <div>
          <p className='flex text-gray-600 items-center truncate'>
            <ClockIcon className='h-4 w-4 mr-2'/>Opening hours
          </p>
          {text}
        </div>
      );
    }

    const renderGmapBtn = () => {
      return (
        <a
          href={placeDetails.url}
          className='flex w-fit rounded-full py-2 px-6 mt-8 items-center border border-gray-200 font-semibold text-gray-500'>
          <GoogleIcon className='h-4 w-4 mr-2'/> Google Maps
        </a>
      );
    }

    return (
      <div className='bg-white p-4 mx-4 h-11/12 w-11/12 max-w-3xl rounded-xl pointer-events-auto'>
        {renderHeader()}
        <p className='text-gray-600 mb-1'>
          {_get(placeDetails, "editorial_summary.overview", "")}
        </p>
        <p className='text-gray-600 flex items-center mb-1'>
          <MapPinIcon className="h-4 w-4" />&nbsp;
          {placeDetails.formatted_address}
        </p>
        {renderRatings()}
        {renderPhone()}
        {renderWebsite()}
        {renderOpeningHours()}
        {renderGmapBtn()}
      </div>
    );
  }



  return (
    <div className='fixed w-screen h-screen'>
      <div
        className='absolute bottom-0 mb-8 z-10 pointer-events-none'
        style={{width: props.width}}
      >
        {renderPlaceDetailsCard()}
      </div>
      <Wrapper
        apiKey={PLACE_IMAGE_APIKEY}
        render={renderMap}
        libraries={["marker"]}
      >
        <InnerMap markers={makeMarkers()} width={props.width} />
      </Wrapper>
    </div>
  );
}

export default TripMap;
