export type Place = any

export const MapElementID = "map";
export const EventMarkerClickName = "marker_click";
export const EventZoomMarkerClick = "zoom_marker_click";

export const newEventMarkerClick = (detail: any) => {
  return new CustomEvent(EventMarkerClickName, {
    bubbles: false,
    cancelable: false,
    detail: detail,
  })
}

export const newZoomMarkerClick = (detail: any) => {
  return new CustomEvent(EventZoomMarkerClick, {
    bubbles: false,
    cancelable: false,
    detail: detail,
  })
}
