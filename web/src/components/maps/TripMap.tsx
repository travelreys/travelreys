import React, {
  FC,
  useEffect,
  useState,
  useRef
} from 'react';
import _get from "lodash/get";
import _flatten from "lodash/flatten";
import _isEmpty from "lodash/isEmpty";
import _find from "lodash/find";
import { Wrapper, Status } from "@googlemaps/react-wrapper";
import {
  ClockIcon,
  StarIcon,
  MapPinIcon,
  PhoneIcon,
  GlobeAltIcon,
  Square3Stack3DIcon,
} from '@heroicons/react/24/solid'
import {
  MagnifyingGlassCircleIcon,
  XMarkIcon
} from '@heroicons/react/24/outline'

import MapsAPI, {
  placeAtmosphereFields,
  PLACE_IMAGE_APIKEY
} from '../../apis/maps';
import {
  ContentIconOpts,
  DefaultContentColor,
  LabelContentListColor,
  LabelContentListIcon,
  Trips
} from '../../apis/trips';

import Spinner from '../Spinner';
import ContentListPin from './ContentListPin';

import {
  EventMarkerClickName,
  EventZoomMarkerClick,
  MapElementID,
  newZoomMarkerClick,
} from './common';
import GoogleIcon from '../icons/GoogleIcon';
import { makeNumberPin, makePinWithTooltip } from './mapsPinIcons';
import { ActionNameSetSelectedPlace, useMap } from '../../context/maps-context';
import { CommonCss, TripMapCss } from '../../styles/global';
import {
  parseISO,
  printFmt
} from '../../utils/dates'


const defaultMapCenter = { lat: 1.290969, lng: 103.8560011 }
const defaultMapOpts = {
  center: defaultMapCenter,
  zoom: 14,
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

interface ContentListMapDrawings {
  id: string
  markers: Array<PlaceMarker>
  visible: boolean
}

interface ItineraryMapDrawings {
  id: string
  markers: Array<PlaceMarker>
  visible: boolean
  polylines: Array<any>
}

interface InnerMapProps {
  contentListMapDrawings: Array<ContentListMapDrawings>
  itineraryMapDrawings: Array<ItineraryMapDrawings>
  width: any
}

const InnerMap: FC<InnerMapProps> = (props: InnerMapProps) => {
  const ref = useRef() as any;
  const map = useRef() as any;
  const { state, dispatch } = useMap();

  const currentMapCenter = useRef(null) as any;

  // Map Drawings
  const contentListMapPopups = useRef({}) as any;
  const itineraryListMapPopups = useRef({}) as any;
  const itineraryListMapPolylines = useRef({}) as any;


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

  // Helpers
  const clearMap = () => {
    // Clear all markers from the previous render
    [contentListMapPopups, itineraryListMapPopups]
      .forEach((popupList: any) => {
        Object.values(popupList.current)
          .forEach((puList: any) => {
            puList.forEach((pu: Popup) => pu.setMap(null))
        });
      });
    Object.values(itineraryListMapPolylines.current)
      .forEach((pl:any) => pl.setMap(null));

    contentListMapPopups.current = {}
    itineraryListMapPopups.current = {}
    itineraryListMapPolylines.current = {}
  }

  const makeDrawingsMarkers = (drawings: any, popupList: any, bounds: any) => {
    drawings.forEach((draw: ContentListMapDrawings) => {
      if (!draw.visible) {
        return;
      }
      draw.markers.forEach((marker: PlaceMarker) => {
        const latlng = _get(marker, "place.geometry.location") as any;
        const popup = new Popup(latlng, marker.elem);
        popup.setMap(map.current);
        popup.addListener("click", () => { popupOnClick(popup, marker)})
        const arr = _get(popupList.current, draw.id, [])
        arr.push(popup)
        popupList.current[draw.id] = arr;
        bounds.extend(popup.position);
      });
    });
  }

  const makeDrawingsPolylines = (drawings: any, polylineList: any) => {
    drawings.forEach((draw: ItineraryMapDrawings) => {
      if (!draw.visible) {
        return;
      }
      const latlng = _flatten(
        _get(draw.polylines, "polylines", [])
          .map((pl: any) => google.maps.geometry.encoding.decodePath(pl))
      ) as any;
      const path = new google.maps.Polyline({
        path: latlng,
        geodesic: true,
        strokeColor: _get(draw.polylines, "color"),
        strokeOpacity: 0.6,
        strokeWeight: 8,
      });
      path.setMap(map.current);
      polylineList.current[draw.id] = path;
    });
  }

  // Event Handlers - Map
  const mapOnCenterChange = () => {
    currentMapCenter.current = map.current.getCenter();
  }

  const mapOnClick = (e: any) => {
    if (e.placeId) {
      // Call event.stop() on the event to prevent the default info window from showing.
      e.stop();
      dispatch({
        type: ActionNameSetSelectedPlace,
        value: {
          place_id: e.placeId,
          geometry: { location: e.latLng }
        }
      });
      map.current.panTo(e.latLng);
    }
  }

  // Event Handlers - Map Wrapper Ref
  const customMarkerOnClick = (e: any) => {
    const center = _get(e.detail, "geometry.location", defaultMapCenter);
    map.current.setCenter(center);
  }

  const customMarkOnZoom = (e: any) => {
    const center = _get(e.detail, "geometry.location", defaultMapCenter);
    map.current.setCenter(center);
    map.current.setZoom(17);
  }

  // Event Handlers - Popup
  const popupOnClick = (popup: Popup, marker: PlaceMarker) => {
    map.current.setCenter(popup.position);
    currentMapCenter.current = popup.position;
    dispatch({ type: ActionNameSetSelectedPlace, value: marker.place })
  }

  // useEffects

  useEffect(() => {
    map.current = new window.google.maps.Map(ref.current, defaultMapOpts)
    map.current.addListener("center_changed", mapOnCenterChange);
    map.current.addListener("click", mapOnClick)
    ref.current.addEventListener(EventMarkerClickName, customMarkerOnClick);
    ref.current.addEventListener(EventZoomMarkerClick, customMarkOnZoom);
  }, [])

  useEffect(() => {
    if (!_isEmpty(state.selectedPlace)) {
      const center = _get(state.selectedPlace, "geometry.location");
      map.current.panTo(center);
    }
  }, [state.selectedPlace])

  useEffect(() => {
    clearMap();

    const bounds = new google.maps.LatLngBounds();

    // Content List Markers
    makeDrawingsMarkers(props.contentListMapDrawings, contentListMapPopups, bounds);
    makeDrawingsMarkers(props.itineraryMapDrawings, itineraryListMapPopups, bounds);
    makeDrawingsPolylines(props.itineraryMapDrawings, itineraryListMapPolylines)

    // map.current.setCenter(currentMapCenter.current);
    if (_isEmpty(currentMapCenter.current)) {
      map.current.fitBounds(bounds);
      // map.current.setCenter(bounds.getCenter());
    }

  }, [props.contentListMapDrawings, props.itineraryMapDrawings])



  // Renderers

  return (
    <div ref={ref}
      id={MapElementID}
      className='h-full'
      style={{ width: props.width }}
    />
  );
}


//////////////////////
// PlaceDetailsCard //
//////////////////////

interface PlaceDetailsCardProps {
  placeDetails: any
  width: string
  onClose: () => void
}

const PlaceDetailsCard: FC<PlaceDetailsCardProps> = (props: PlaceDetailsCardProps) => {
  const { placeDetails } = props;

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
        <button type="button" onClick={props.onClose}>
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

  if (placeDetails === null) {
    return null;
  }

  return (
    <div
      className={TripMapCss.DetailsWrapper}
      style={{ width: props.width }}
    >
      <div className={TripMapCss.DetailsCard}>
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


///////////////////
// MapLayersMenu //
///////////////////

interface MapLayersMenuProps {
  trip: any
  layersViz: any
  onSelect: (layersSelected: any) => void
}

const MapLayersMenu: FC<MapLayersMenuProps> = (props: MapLayersMenuProps) => {

  const selectAllOnClick = () => {
    const newLayersViz = Object.assign({}, props.layersViz);
    Object.keys(newLayersViz).forEach((k: string) => {
      newLayersViz[k] = true;
    });
    props.onSelect(newLayersViz);
  }

  const deselectAllOnClick = () => {
    const newLayersViz = Object.assign({}, props.layersViz);
    Object.keys(newLayersViz).forEach((k: string) => {
      newLayersViz[k] = false;
    });
    props.onSelect(newLayersViz);
  }

  const onSelectLayer = (id: string) => {
    const newLayersViz = Object.assign({}, props.layersViz);
    newLayersViz[id] = !newLayersViz[id];
    props.onSelect(newLayersViz);
  }

  // Renderers

  const renderHomeOpts = () => {
    const contents = [
      {
        id: "lodgings",
        name: "Hotels and Lodgings",
        labels: { color: "black", icon: "hotel"}
      }
    ].concat(Object.values(_get(props.trip, "contents", [])))

    const contentLayers = contents
      .map((l: any) => {
        const color = _get(l, `labels.${LabelContentListColor}`, DefaultContentColor);
        const icon = ContentIconOpts[_get(l, `labels.${LabelContentListIcon}`, "")];
        return (
          <div key={l.id} className='flex justify-between items-center mb-1'>
            <div className='flex items-center flex-1 text-gray-700'>
              <ContentListPin icon={icon} color={color}/>
              {l.name}
            </div>
            <input
              type="checkbox"
              checked={props.layersViz[l.id]}
              className="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 rounded focus:ring-blue-500"
              onChange={() => {onSelectLayer(l.id)}}
            />
          </div>
        )
      });
    return (
      <div className='w-full'>
        <div className='mb-1'>
          <h4 className='text-sm font-bold mb-1'>Overview</h4>
          {contentLayers}
        </div>
      </div>
    );
  }

  const renderItineraryOpts = () => {
    const itinLayers = _get(props.trip, "itinerary", [])
      .map((l: any) => {
        const color = _get(l, `labels.${LabelContentListColor}`, DefaultContentColor);
        return (
          <div key={l.id} className='flex justify-between items-center mb-1'>
            <div className='flex items-center flex-1 text-gray-700'>
              <ContentListPin icon={""} color={color}/>
              { printFmt(parseISO(l.date), "eee, MM/dd")  }
            </div>
            <input
              type="checkbox"
              checked={props.layersViz[l.id]}
              className="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 rounded focus:ring-blue-500"
              onChange={() => {onSelectLayer(l.id)}}
            />
          </div>
        )
      });
    return (
      <div className='w-full'>
        <div className='mb-1'>
          <h4 className='text-sm font-bold mb-1'>Itinerary</h4>
          {itinLayers}
        </div>
      </div>
    );
  }


  return (
    <div className='absolute mt-8 ml-12 bg-white p-4 rounded-xl z-10 w-96'>
      <h4 className='font-bold text-base mb-2'>Map Drawings</h4>
      <div className='flex mb-1 text-sm'>
        <button
          className='font-semibold text-indigo-500 mr-2 hover:text-indigo-700'
          onClick={selectAllOnClick}
        >
          Select all
        </button>
        <button
          className='font-semibold text-gray-500 hover:text-indigo-700'
          onClick={deselectAllOnClick}
        >
          Deselect all
        </button>
      </div>
      <hr className='my-2' />
      {renderHomeOpts()}
      {renderItineraryOpts()}
    </div>
  );
}



/////////////
// TripMap //
/////////////
interface TripMapProps {
  trip: any
  width: any
}

const TripMap: FC<TripMapProps> = (props: TripMapProps) => {

  const { state } = useMap();
  const [isLayersMenuOpen, setIsLayersMenuOpen] = useState(false);
  const [placeDetails, setPlaceDetails] = useState(null) as any;
  const [layersViz, setLayersViz] = useState() as any;

  // Helpers
  useEffect(() => {
    if (!_isEmpty(layersViz)) {
      return;
    }
    const vis = {"lodgings": true} as any;
    Object.values(_get(props.trip, "contents", {}))
      .forEach((ctntList: any) =>  {vis[ctntList.id] = true});
    Object.values(_get(props.trip, "itinerary", {}))
      .forEach((itinList: any) =>  {vis[itinList.id] = true});
    setLayersViz(vis);
  }, [props.trip])


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


  // Events Handlers
  const onSelectLayersViz = (selection: any) => {
    setLayersViz(selection);
  }

  // Map Markers

  const makeContentListsMapDrawings = () => {
    // Lodging
    const lodgings = _get(props.trip, "lodgings", {});
    const markers = Object.values(lodgings).map((lodge: any) => ({
      elem: makePinWithTooltip(lodge.place.name, "black", "hotel"),
      place: lodge.place
    }));
    let drawings = [{
      id: "lodgings",
      markers,
      visible: _get(layersViz, "lodgings", true),
    } as ContentListMapDrawings]

    // Content Lists
    const ctntLists = Object.values(_get(props.trip, "contents", {}));
    drawings = drawings.concat(
      ctntLists.map((l: any) => {
        const color = _get(l, `labels.${LabelContentListColor}`, DefaultContentColor)
        const icon = _get(l, `labels.${LabelContentListIcon}`, "")
        const markers = l.contents
          .filter((ct: any) => {
            const latlng = _get(ct, "place.geometry.location")
            return latlng !== undefined && !(latlng.lat === 0 && latlng.lng === 0)
          })
          .map((ct: any) => ({
            elem: makePinWithTooltip(ct.place.name, color, icon),
            place: ct.place
          }));
        return {
          id: l.id,
          markers,
          visible: _get(layersViz, l.id, true),
        } as ContentListMapDrawings
      })
    );
    return drawings;
  }

  const makeItineraryListsMapDrawings = () => {
    const itntCtntList = _get(props.trip, "itinerary", []);

    return itntCtntList.map((l: any) => {
      const color = _get(l, `labels.${LabelContentListColor}`, DefaultContentColor)

      const markers = l.contents.map((itinCtnt: any, idx: number) => {
        const ctnt = _find(
          _get(props.trip, `contents.${itinCtnt.tripContentListId}.contents`, []),
          (ctn: Trips.Content) => ctn.id === itinCtnt.tripContentId
        );
        const latlng = _get(ctnt, "place.geometry.location")
        if (latlng === undefined || (latlng.lat === 0 && latlng.lng === 0)) {
          return null;
        }
        return {
          elem: makeNumberPin(ctnt.place.name, color, `${idx}`),
          place: ctnt.place
        };
      })
      .filter((item: any) => item !== null);
      const polylines = l.routes.map((r: any) => r.overview_polyline.points);
      return {
        id: l.id,
        markers,
        visible: _get(layersViz, l.id, true),
        polylines: {polylines, color: color}
      }
    });
  }


  // Renderers
  const renderLayersBtn = () => {
    return (
      <div className='relative'>
        <button
          type="button"
          onClick={() => {setIsLayersMenuOpen(!isLayersMenuOpen)}}
          className='absolute mt-4 ml-4 p-2 bg-white/75 rounded-full z-10'
        >
          <Square3Stack3DIcon className='h-6 w-6' />
        </button>
        {isLayersMenuOpen
          ? <MapLayersMenu
              trip={props.trip}
              layersViz={layersViz}
              onSelect={onSelectLayersViz}
            />
          : null
        }
      </div>
    );
  }

  const renderMap = (status: Status): React.ReactElement => {
    if (status === Status.FAILURE) return <Spinner />;
    return <Spinner />;
  };

  return (
    <div
      className={TripMapCss.Ctn}
      style={{width: props.width}}
    >
      {renderLayersBtn()}
      <PlaceDetailsCard
        placeDetails={placeDetails}
        width={props.width}
        onClose={() => { setPlaceDetails(null) }}
      />
      <Wrapper
        apiKey={PLACE_IMAGE_APIKEY}
        render={renderMap}
        libraries={["marker", "geometry"]}
      >
        <InnerMap
          width={props.width}
          contentListMapDrawings={makeContentListsMapDrawings()}
          itineraryMapDrawings={makeItineraryListsMapDrawings()}
        />
      </Wrapper>
    </div>
  );
}

export default TripMap;
