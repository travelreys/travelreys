import React, {
  FC,
  useEffect,
  useState,
  useRef
} from 'react';
import _get from "lodash/get";
import _isEmpty from "lodash/isEmpty";

import { Wrapper, Status } from "@googlemaps/react-wrapper";
import { PLACE_IMAGE_APIKEY } from '../../apis/maps';

import Spinner from '../Spinner';
import { makeHotelPin } from './GMapsPinIcon';

interface InnerMapProps {
  markers: any
  width: any
}

const InnerMap: FC<InnerMapProps> = (props: InnerMapProps) => {
  const ref = useRef() as any;
  const map = useRef() as any;

  const currentPopups = useRef([]) as any;

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
    console.log("rendering")
    map.current = new window.google.maps.Map(ref.current, {
      center: { lat: 33.3960897, lng: 126.264522 },
      zoom: 10,
      mapTypeControl: false,
      streetViewControl: false,
      fullscreenControl: false,
      rotateControl: false,
      keyboardShortcuts: true,
      gestureHandling: "greedy"
    })
  }, [])

  useEffect(() => {
    currentPopups.current.forEach((pp: any) => {
      pp.setMap(null);
    })
    currentPopups.current = [];

    const bounds = new google.maps.LatLngBounds();
    props.markers.forEach((marker: any) => {
      const popup = new Popup(marker.latlng, marker.elem);
      popup.setMap(map.current);
      currentPopups.current.push(popup)
      bounds.extend(popup.position);
    });
    map.current.setCenter(bounds.getCenter());
    // map.current.fitBounds(bounds);

  }, [props.markers])

  return (
    <div ref={ref} id="map" className='h-full' style={{width: props.width}}/>
  );
}

interface TripMapComponentProps {
  trip: any
  width: any
}

const TripMap: FC<TripMapComponentProps> = (props: TripMapComponentProps) => {

  // Map Markers
  const lodgingToMapMarkers = () => {
    const lodgings = _get(props.trip, "lodgings", []);
    if (_isEmpty(lodgings)) {
      return [];
    }

    return Object.values(lodgings).map((lodge: any) => {
      const elem = makeHotelPin(lodge.place.name);
      const latlng = {
        lat: lodge.place.geometry.location.lat,
        lng: lodge.place.geometry.location.lng,
      } as any;
      return {latlng, elem}
    });
  }

  const makeMarkers = () => {
    let markers = [] as any;
    markers = markers.concat(lodgingToMapMarkers());
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
