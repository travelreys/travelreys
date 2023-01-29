import React, { FC, useEffect, useState, useRef} from 'react';
import { Wrapper, Status } from "@googlemaps/react-wrapper";
import { PLACE_IMAGE_APIKEY } from '../../apis/maps';
import { gMapsMarkerFactory  } from './GMapsPinIcon';


import Spinner from '../Spinner';

interface InnerMap {}

const InnerMap: FC = () => {

  class Popup extends google.maps.OverlayView {
    position: google.maps.LatLng;
    containerDiv: HTMLDivElement;

    constructor(position: google.maps.LatLng, content: HTMLElement) {
      super();
      this.position = position;

      content.classList.add("popup-bubble");

      // This zero-height div is positioned at the bottom of the bubble.
      const bubbleAnchor = document.createElement("div");

      // bubbleAnchor.classList.add("popup-bubble-anchor");
      bubbleAnchor.appendChild(content);

      // This zero-height div is positioned at the bottom of the tip.
      this.containerDiv = document.createElement("div");
      this.containerDiv.classList.add("popup-container");
      this.containerDiv.appendChild(bubbleAnchor);

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

  const ref = useRef() as any;
  const [text, setText] = useState("")

  const contentString =
    '<div id="content" class="text-red-500">' +
    "</div>" +
      "<p><b>Uluru</b> </p>" +
    "</div>" +
    "</div>";

  useEffect(() => {
    const map = new window.google.maps.Map(ref.current, {
      center: { lat: 37.42, lng: -122.1 },
      zoom: 14,
      mapTypeControl: false,
      gestureHandling: "greedy"
    })

    const infowindow = new google.maps.InfoWindow({
      content: contentString,
      ariaLabel: "Uluru",
    });


    const marker =  gMapsMarkerFactory("blue");
    // marker.setMap(map);
    // marker.setPosition(map.getCenter());
    // window.google.maps.event.addListener(marker, "mouseover", () => {
    //   const symbol = marker.getIcon() as google.maps.Symbol;
    //   symbol.fillOpacity = 1
    //   marker.setIcon(symbol);
    //   infowindow.open({anchor: marker, map})
    //   setText("jello")
    // });

    // window.google.maps.event.addListener(marker, "mouseout", () => {
    //   const symbol = marker.getIcon() as google.maps.Symbol;
    //   symbol.fillOpacity = 0.8;
    //   marker.setIcon(symbol);
    //   infowindow.close()
    //   setText("1123")
    // });

    const priceTag = document.createElement("p");
    priceTag.textContent = "123333333"
    priceTag.addEventListener("mouseover", () => {
      priceTag.classList.add("text-red-500")
      priceTag.textContent = 'over'
    })

    const popup = new Popup(
      map.getCenter()!,
      priceTag
    );
    popup.setMap(map);

  }, [])

  return (
    <>
      <div ref={ref} id="map" className='w-full h-screen'/>
      <p>{text}</p>
    </>
  );
}

interface MapComponentProps {}

const MapComponent: FC<MapComponentProps> = (props: MapComponentProps) => {

  // Renderer
  const render = (status: string) => {
    switch (status) {
      case Status.FAILURE:
        return <div>error</div>;
      case Status.LOADING:
        return <Spinner />;
      default:
        return <></>;
    }
  }

  return (
    <Wrapper
      apiKey={PLACE_IMAGE_APIKEY}
      render={render}
      libraries={["marker"]}
    >
      <InnerMap />
    </Wrapper>
  );
}

export default MapComponent;
