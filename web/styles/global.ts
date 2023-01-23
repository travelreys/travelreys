export const ModalCss = {
  Container: "relative z-10",
  Inset: "fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity",
  Content: "fixed inset-0 z-10 overflow-y-auto",
  ContentContainer: "flex min-h-full flex-col p-4 text-center sm:items-center sm:p-0",
  ContentCard: "bg-white relative transform rounded-lg text-left shadow-xl transition-all pb-5 sm:my-8 sm:w-full sm:max-w-2xl",
}

export const CreateTripModalCss = {
  TripNameCtn: "flex mb-4 border border-slate-200 rounded-lg",
  TripNameLabel: "inline-flex font-bold items-center px-3 text-sm text-slate-500 rounded-l-md",
  TripNameInput: "block flex-1 border-0 rounded-r-lg min-w-0 p-2.5 text-gray-900 text-sm w-full",
  TripDatesCtn: "flex w-full border border-slate-200 rounded-lg",
  TripDatesIcon: "inline align-bottom h-5 w-5 text-gray-500",
  TripDatesLabel: "inline-flex font-bold items-center px-3 text-sm text-slate-500 rounded-l-md",
  TripDatesInput: "block flex-1 min-w-0 p-2.5 border-0 rounded-none rounded-r-lg text-gray-900 text-sm w-full",
}

export const TripStatsCss = {

};

export const FlightsModalCss = {
  ItineraryDropdownBtn: "hover:bg-indigo-100 font-medium rounded-lg text-sm px-4 py-2.5 text-center inline-flex items-center",
  ItineraryDropdownIcon: "h-4 w-4 text-slate-700",
  FlightSearchBtn: "bg-indigo-500 font-medium rounded-lg text-sm text-white px-4 py-2.5 text-center inline-flex items-center",
  FlightFromIconCtn: "absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none mb-2",
  FlightFromIcon: "h-6 w-6 text-slate-700",
  FlightFromInput: "border border-slate-200 text-gray-900 text-sm rounded block w-full pl-10 p-4",
  FlightSearchHR: "h-px my-8 bg-gray-200 border-0",
  FlightSearchResultsTitle: "text-lg sm:text-2xl mb-2 font-medium text-slate-900",
  FlightPlusIcon: "h-6 w-6 text-green-500 cursor-pointer",
  FlightBookLink: "underline",
  FlightDatesCtn: "flex w-full border border-slate-200 rounded-lg",
  FlightDatesIcon: "inline align-bottom h-5 w-5 text-gray-500",
  FlightDatesLabel: "inline-flex font-bold items-center px-3 text-sm text-slate-500 rounded-l-md",
  FlightDatesInput: "block flex-1 min-w-0 p-2.5 border-0 rounded-none rounded-r-lg text-gray-900 text-sm w-full",
  FlightStopIcon: "inline h-2 w-2 text-slate-700",
  FlightStopTimelineIcon: "absolute flex mt-1 text-indigo-200 items-center justify-center w-4 h-4 -left-2 bg-indigo-200 rounded-full ring-8 ring-white",
  FlightStopTimelineTime: "mb-1 font-medium text-slate-900",
  FlightsStopTimelineText: "mb-2 text-sm font-normal leading-none text-slate-400",
  FlightsStopLayoverText: "mb-1 text-sm font-normal leading-none text-red-700",
  FlightsStopHR: "w-48 h-1 mx-auto my-4 bg-gray-100 border-0 rounded md:my-10 dark:bg-gray-700",
}