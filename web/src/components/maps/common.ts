export const EventMarkerClickName = "marker_click";
export const newEventMarkerClick = (detail: any) => {
  return new CustomEvent(EventMarkerClickName, {
    bubbles: false,
    cancelable: false,
    detail: detail,
  })
}
