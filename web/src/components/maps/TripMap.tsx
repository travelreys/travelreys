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
import { PLACE_IMAGE_APIKEY } from '../../apis/maps';

import Spinner from '../Spinner';
import { makeActivityPin, makeHotelPin } from './GMapsPinIcon';
import { useMap } from '../../context/maps-context';


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
      if (state.selectedPlace !== null) {
        dispatch({type:"setSelectedPlace", value: null})
      }
    });
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
  const render = (status: Status): React.ReactElement => {
    if (status === Status.FAILURE) return <Spinner />;
    return <Spinner />;
  };

  return (
    <Wrapper
      apiKey={PLACE_IMAGE_APIKEY}
      render={render}
      libraries={["marker"]}
    >
      <InnerMap markers={makeMarkers()} width={props.width} />
    </Wrapper>
  );
}

export default TripMap;
