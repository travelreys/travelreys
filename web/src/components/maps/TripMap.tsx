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
  MagnifyingGlassCircleIcon,
  XMarkIcon
} from '@heroicons/react/24/outline'

import MapsAPI, {
  placeAtmosphereFields,
  PLACE_IMAGE_APIKEY
} from '../../apis/maps';

import Spinner from '../Spinner';
import {
  EventMarkerClickName,
  EventZoomMarkerClick,
  MapElementID,
  newZoomMarkerClick,
} from './common';
import GoogleIcon from '../icons/GoogleIcon';
import { makeActivityPin, makeHotelPin } from './mapsPinIcons';
import { ActionNameSetSelectedPlace, useMap } from '../../context/maps-context';
import { CommonCss, TripMapCss } from '../../styles/global';
import { Trips } from '../../apis/trips';


const defaultMapCenter = { lat: 33.3960897, lng: 126.264522 }
const defaultMapOpts = {
  center: defaultMapCenter,
  zoom: 10,
  mapTypeControl: false,
  streetViewControl: false,
  fullscreenControl: false,
  rotateControl: false,
  keyboardShortcuts: true,
  gestureHandling: "greedy",
}

interface PlaceMarker {
  elem: HTMLElement
  place: any
}


interface InnerMapProps {
  markers: Array<PlaceMarker>
  itineraryPolylines: Array<Array<string>>
  width: any
}

const InnerMap: FC<InnerMapProps> = (props: InnerMapProps) => {
  const ref = useRef() as any;
  const map = useRef() as any;
  const { state, dispatch } = useMap();

  const currentPopups = useRef([]) as any;
  const currentPolylines = useRef([]) as any;
  const currentMapCenter = useRef(null) as any;

  // Popup
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
      this.containerDiv.addEventListener('click', function () {
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

  // useEffects

  useEffect(() => {
    map.current = new window.google.maps.Map(ref.current, defaultMapOpts)
    map.current.addListener("center_changed", () => {
      currentMapCenter.current = map.current.getCenter();
    });
    map.current.addListener("click", (e: any) => {
      if (e.placeId) {
        // Call event.stop() on the event to prevent the default info window from showing.
        e.stop();
        // do any other stuff you want to do
        dispatch({
          type: ActionNameSetSelectedPlace,
          value: {
            place_id: e.placeId,
            geometry: {
              location: { lat: e.latLng.lat(), lng: e.latLng.lng() }
            }
          }
        });
        map.current.panTo({ lat: e.latLng.lat(), lng: e.latLng.lng() });
      }
    })
    ref.current.addEventListener(EventMarkerClickName, (e: any) => {
      const center = _get(e.detail, "geometry.location", defaultMapCenter);
      map.current.setCenter(center);
    });
    ref.current.addEventListener(EventZoomMarkerClick, (e: any) => {
      const center = _get(e.detail, "geometry.location", defaultMapCenter);
      map.current.setCenter(center);
      map.current.setZoom(17);
    });
  }, [])

  useEffect(() => {
    if (!_isEmpty(state.selectedPlace)) {
      const center = _get(state.selectedPlace, "geometry.location");
      map.current.panTo(center);
    }
  }, [state.selectedPlace])

  useEffect(() => {
    // Clear all markers from the previous render
    currentPopups.current.forEach((pp: any) => {
      pp.setMap(null);
    })
    currentPopups.current = [];

    currentPolylines.current.forEach((pl: any) => {
      pl.setMap(null);
    })
    currentPolylines.current = [];


    // Make new markers
    const bounds = new google.maps.LatLngBounds();
    props.markers.forEach((marker: PlaceMarker) => {
      const latlng = _get(marker, "place.geometry.location") as any;
      const popup = new Popup(latlng, marker.elem);
      popup.setMap(map.current);
      popup.addListener("click", () => {
        map.current.setCenter(popup.position);
        currentMapCenter.current = popup.position;
        dispatch({ type: ActionNameSetSelectedPlace, value: marker.place })
      })
      currentPopups.current.push(popup);
      bounds.extend(popup.position);
    });

    if (_isEmpty(currentMapCenter.current)) {
      map.current.setCenter(bounds.getCenter());
      currentMapCenter.current = bounds.getCenter();
    }

    props.itineraryPolylines.forEach((polylines: Array<string>) => {
      const latlng = _flatten(polylines.map((pl) => google.maps.geometry.encoding.decodePath(pl)));
      const path = new google.maps.Polyline({
        path: latlng,
        geodesic: true,
        strokeColor: "#FF0000",
        strokeOpacity: 0.6,
        strokeWeight: 8,
      });
      path.setMap(map.current);
      currentPolylines.current.push(path);
    })
  }, [props.markers])


  // Renderers

  return (
    <div ref={ref}
      id={MapElementID}
      className='h-full'
      style={{ width: props.width }}
    />
  );
}

interface TripMapComponentProps {
  trip: any
  width: any
}

const TripMap: FC<TripMapComponentProps> = (props: TripMapComponentProps) => {

  const { state } = useMap();
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

  // Direction polylines
  const makeItineraryDirectionsPolylines = () => {
    return _get(props.trip, "itinerary", []).map((itin: Trips.ItineraryList) => {
      return _flatten(itin.routes.map((r) => r.overview_polyline.points));
    });
  }

  // Renderers

  const renderPlaceDetailsCard = () => {
    if (placeDetails === null) {
      return null;
    }

    const renderHeader = () => {
      return (
        <p className={TripMapCss.HeaderCtn}>
          <span className={TripMapCss.TitleCtn}>
            <button type="button" onClick={() => {
              const event = newZoomMarkerClick(placeDetails);
              document.getElementById(MapElementID)?.dispatchEvent(event)
            }}>
              <MagnifyingGlassCircleIcon className={CommonCss.LeftIcon} />
            </button>
            {placeDetails.name}
          </span>
          <button type="button" onClick={() => { setPlaceDetails(null) }}>
            <XMarkIcon className={CommonCss.Icon} />
          </button>
        </p>
      );
    }

    const renderSummary = () => {
      return (
        <p className={TripMapCss.SummaryTxt}>
          {_get(placeDetails, "editorial_summary.overview", "")}
        </p>
      );
    }

    const renderAddr = () => {
      return (
        <p className={TripMapCss.AddrTxt}>
          <MapPinIcon className={CommonCss.LeftIcon} />
          {placeDetails.formatted_address}
        </p>
      );
    }

    const renderRatings = () => {
      if (placeDetails.user_ratings_total === 0) {
        return null;
      }
      return (
        <p className={TripMapCss.RatingsStar}>
          <StarIcon className={CommonCss.LeftIcon} />
          {placeDetails.rating}&nbsp;&nbsp;
          <span className={TripMapCss.RatingsTxt}>
            ({placeDetails.user_ratings_total})
          </span>
          &nbsp;&nbsp;
          <GoogleIcon className={CommonCss.DropdownIcon} />
        </p>
      );
    }

    const renderOpeningHours = () => {
      const weekdayTexts = _get(placeDetails, "opening_hours.weekday_text", []);
      if (_isEmpty(weekdayTexts)) {
        return null;
      }
      return (
        <div>
          <p className={TripMapCss.OpeningHrsTxt}>
            <ClockIcon className={CommonCss.LeftIcon} />Opening hours
          </p>
          {weekdayTexts.map((txt: string, idx: number) =>
            (<p key={idx} className={TripMapCss.WeekdayTxt}>{txt}</p>)
          )}
        </div>
      );
    }

    const renderPhone = () => {
      return placeDetails.international_phone_number
        ?
        <a
          href={`tel:${placeDetails.international_phone_number.replace(/\s/, "-")}`}
          target="_blank"
          className={TripMapCss.PhoneBtn}
        >
          <PhoneIcon className={TripMapCss.PhoneIcon} />
          Call
        </a>
        : null
    }

    const renderWebsite = () => {
      return placeDetails.website
        ?
        <a
          href={placeDetails.website}
          target="_blank"
          className={TripMapCss.PhoneBtn}
        >
          <GlobeAltIcon className={TripMapCss.PhoneIcon} />
          Web
        </a>
        : null
    }

    const renderGmapBtn = () => {
      return (
        <a
          href={placeDetails.url}
          className={TripMapCss.GmapBtn}
        >
          <GoogleIcon className={CommonCss.LeftIcon} /> Google Maps
        </a>
      );
    }

    return (
      <div
        className={TripMapCss.DetailsWrapper}
        style={{ width: props.width }}
      >
        <div className={TripMapCss.DetailsCard}>
          {placeDetails.place_id}
          {renderHeader()}
          {renderSummary()}
          {renderAddr()}
          {renderRatings()}
          {renderOpeningHours()}
          <div className={TripMapCss.BtnCtn}>
            {renderPhone()}
            {renderWebsite()}
            {renderGmapBtn()}
          </div>
        </div>
      </div>
    );
  }

  const renderMap = (status: Status): React.ReactElement => {
    if (status === Status.FAILURE) return <Spinner />;
    return <Spinner />;
  };

  return (
    <div className={TripMapCss.Ctn}>
      {renderPlaceDetailsCard()}
      <Wrapper
        apiKey={PLACE_IMAGE_APIKEY}
        render={renderMap}
        libraries={["marker", "geometry"]}
      >
        <InnerMap
          markers={makeMarkers()}
          itineraryPolylines={makeItineraryDirectionsPolylines()}
          width={props.width}
        />
      </Wrapper>
    </div>
  );
}

export default TripMap;
