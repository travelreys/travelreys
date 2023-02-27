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
import { Square3Stack3DIcon } from '@heroicons/react/24/solid'

import { makeNumberPin, makePinWithTooltip } from './mapsPinIcons';
import Spinner from '../../components/common/Spinner';

import MapsAPI, {
  placeAtmosphereFields,
  PlaceDetailsResponse,
  PLACE_IMAGE_APIKEY
} from '../../apis/maps';
import {
  Activity,
  DefaultActivityColor,
  getActivityColor,
  getActivityIcon,
  getSortedActivies,
  LabelUiColor,
} from '../../lib/trips';
import {
  EventMarkerClickName,
  EventZoomMarkerClick,
  MapElementID,
} from '../../lib/maps';
import { ActionSetSelectedPlace, useMap } from '../../context/maps-context';
import { TripMapCss } from '../../assets/styles/global';
import MapLayersMenu from './MapLayersMenu';
import PlaceDetailsCard from './PlaceDetailsCard';


const defaultMapCenter = { lat: 1.290969, lng: 103.8560011 }
const defaultMapOpts = {
  center: defaultMapCenter,
  zoom: 12,
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

interface ActivityListMapDrawings {
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
  contentListMapDrawings: Array<ActivityListMapDrawings>
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
    drawings.forEach((draw: ActivityListMapDrawings) => {
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
        type: ActionSetSelectedPlace,
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
    dispatch({ type: ActionSetSelectedPlace, value: marker.place })
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

    // Activity List Markers
    makeDrawingsMarkers(props.contentListMapDrawings, contentListMapPopups, bounds);
    makeDrawingsMarkers(props.itineraryMapDrawings, itineraryListMapPopups, bounds);
    makeDrawingsPolylines(props.itineraryMapDrawings, itineraryListMapPolylines)

    // map.current.setCenter(currentMapCenter.current);
    if (_isEmpty(currentMapCenter.current)) {
      map.current.fitBounds(bounds);
      // map.current.setCenter(bounds.getCenter());
    }

  }, [props.contentListMapDrawings, props.itineraryMapDrawings])

  return (
    <div ref={ref}
      id={MapElementID}
      className='h-full'
      style={{ width: props.width }}
    />
  );
}

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
    Object.values(_get(props.trip, "activities", {}))
      .forEach((l: any) =>  {vis[l.id] = true});
    Object.values(_get(props.trip, "itinerary", {}))
      .forEach((l: any) =>  {vis[l.id] = true});
    setLayersViz(vis);
  }, [props.trip, layersViz])


  useEffect(() => {
    if (state.selectedPlace === null) {
      return;
    }
    const placeID = state.selectedPlace.place_id;
    MapsAPI.placeDetails(placeID, placeAtmosphereFields, "")
    .then((res: PlaceDetailsResponse) => {
      setPlaceDetails(res.place);
    })
  }, [state.selectedPlace]);


  // Events Handlers
  const onSelectLayersViz = (selection: any) => {
    setLayersViz(selection);
  }

  // Map Markers

  const makeActivityListsMapDrawings = () => {
    // Lodging
    const lodgings = _get(props.trip, "lodgings", {});
    const markers = Object.values(lodgings).map((lodge: any) => ({
      elem: makePinWithTooltip(lodge.place.name, "rgb(249 115 22)", "hotel"),
      place: lodge.place
    }));
    let drawings = [{
      id: "lodgings",
      markers,
      visible: _get(layersViz, "lodgings", true),
    } as ActivityListMapDrawings]

    // Activity Lists
    const actLists = Object.values(_get(props.trip, "activities", {}));
    return drawings.concat(
      actLists.map((l: any) => {
        const markers = Object.values(l.activities)
          .filter((ct: any) => {
            const latlng = _get(ct, "place.geometry.location")
            return latlng !== undefined && !(latlng.lat === 0 && latlng.lng === 0)
          })
          .map((ct: any) => ({
            elem: makePinWithTooltip(
              ct.place.name,
              getActivityColor(l) || DefaultActivityColor,
              getActivityIcon(l)),
            place: ct.place
          }));
        return {
          id: l.id,
          markers,
          visible: _get(layersViz, l.id, true),
        } as ActivityListMapDrawings
      })
    );
  }

  const makeItineraryListsMapDrawings = () => {
    const itinList = _get(props.trip, "itinerary", []);
    return itinList.map((l: any) => {
      const color = _get(l, `labels.${LabelUiColor}`, DefaultActivityColor);

      const sortedActivites = getSortedActivies(l);
      const markers = sortedActivites
      .map((itinAct: any, idx: number) => {
        const actList = _get(props.trip, `activities.${itinAct.activityListId}.activities`, []);
        const act = _find(actList, (act: Activity) => act.id === itinAct.activityId);
        const latlng = _get(act, "place.geometry.location")
        if (latlng === undefined || (latlng.lat === 0 && latlng.lng === 0)) {
          return null;
        }
        return {
          elem: makeNumberPin(act.place.name, color, `${idx}`),
          place: act.place
        };
      })
      .filter((item: any) => item !== null);
      const polylines = Object.values(l.routes)
        .map((pairing: any) => (pairing[0].overview_polyline.points))
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
          contentListMapDrawings={makeActivityListsMapDrawings()}
          itineraryMapDrawings={makeItineraryListsMapDrawings()}
        />
      </Wrapper>
    </div>
  );
}

export default TripMap;
