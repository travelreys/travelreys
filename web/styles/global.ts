export const InputDatesPickerCss = {
  Ctn: "flex w-full border border-slate-200 rounded-lg",
  Icon: "inline align-bottom h-5 w-5 text-gray-500",
  Label: "inline-flex font-bold items-center px-3 text-sm text-slate-500 rounded-l-md",
  Input: "block flex-1 min-w-0 p-2.5 border-0 rounded-none rounded-r-lg text-gray-900 text-sm w-full",
}

export const ModalCss = {
  Container: "relative z-10",
  Inset: "fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity",
  Content: "fixed inset-0 z-10 overflow-y-auto",
  ContentContainer: "flex min-h-full flex-col p-4 text-center sm:items-center sm:p-0",
  ContentCard: "bg-white relative transform rounded-lg text-left shadow-xl transition-all pb-5 sm:my-8 sm:w-full sm:max-w-2xl",
}

export const CreateTripModalCss = {
  CreateModalCard: "bg-white px-4 pt-5 pb-4 sm:p-8 sm:pb-4",
  CreateTripTitle: "text-2xl text-center font-medium leading-6 text-slate-900 mb-6",
  TripNameCtn: "flex mb-4 border border-slate-200 rounded-lg",
  TripNameLabel: "inline-flex font-bold items-center px-3 text-sm text-slate-500 rounded-l-md",
  TripNameInput: "block flex-1 border-0 rounded-r-lg min-w-0 p-2.5 text-gray-900 text-sm w-full",
  TripDatesCtn: "flex w-full border border-slate-200 rounded-lg",
  TripDatesIcon: "inline align-bottom h-5 w-5 text-gray-500",
  TripDatesLabel: "inline-flex font-bold items-center px-3 text-sm text-slate-500 rounded-l-md",
  TripDatesInput: "block flex-1 min-w-0 p-2.5 border-0 rounded-none rounded-r-lg text-gray-900 text-sm w-full",
  CreateTripBtnsCtn: "bg-gray-50 px-4 py-3 sm:flex sm:flex-row-reverse sm:px-6",
  CreateTripBtn: "inline-flex w-full justify-center rounded-md border border-transparent bg-indigo-600 px-4 py-2 text-base font-medium text-white shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 sm:ml-3 sm:w-auto sm:text-sm",
  CreateTripCancelBtn: "mt-3 inline-flex w-full justify-center rounded-md border border-gray-300 bg-white px-4 py-2 text-base font-medium text-gray-700 shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm",
}

export const TripMenuJumboCss = {
  TripDatesBtn: "font-medium text-md text-slate-500",
  TripDatesBtnIcon: "inline h-5 w-5 align-sub",
  TripCoverImage: "block sm:max-h-96 w-full",
  TripImageEditIconCtn: "absolute top-4 right-4 h-10 w-10 bg-gray-800/70 p-2 text-center rounded-full",
  TripImageEditIcon: "h-6 w-6 text-white",
  TripNameInputCtn: "h-16 relative -top-24",
  TripNameInputWrapper: "bg-white rounded-lg shadow-xl p-5 mx-4 mb-4",
  TripNameInput: "mb-12 text-2xl sm:text-4xl font-bold text-slate-700 w-full rounded-lg p-1 border-0 hover:bg-slate-300 hover:border-0 hover:bg-slate-100 focus:ring-0",
  Figure: "relative max-w-sm transition-all rounded-lg duration-300 mb-2 group",
  FigureImg: "block rounded-lg max-w-full group-hover:grayscale",
  FigureBtn: "text-white m-2 py-2 px-3 rounded-full bg-green-500 hover:bg-green-700",
  FigureBtnCtn: "absolute group-hover:opacity-100 opacity-0 top-2 right-0",
  FigureCaption: "absolute px-1 text-sm text-white rounded-b-lg bg-slate-800/50 w-full bottom-0",
  SearchImageCard: "bg-white px-4 pt-5 pb-4 sm:p-8 sm:pb-4 rounded-lg",
  SearchImageTitle: "text-lg sm:text-2xl font-bold leading-6 text-slate-900",
  SearchImageWebTitle: "text-sm font-medium text-indigo-500 sm:text-xl text-slate-700 mb-2 ml-1",
  SearchImageInput: "bg-gray-50 block border-gray-300 border focus:border-blue-500 focus:ring-blue-500 min-w-0 p-2.5 rounded-lg text-gray-900 text-sm w-5/6 mr-2",
  SearchImageBtn: "flex-1 inline-flex text-white bg-indigo-500 hover:bg-indigo-800 rounded-2xl p-2.5 text-center items-center justify-around",
  SearchImageIcon: "h-5 w-5 stroke-2 stroke-white",
}


export const TripLogisticsCss = {
  FlightTripCard: "flex p-4 bg-gray-100 rounded shadow-md mb-2",
  FlightsTitleCtn: "flex justify-between mb-4",
  FlightTransitIcon: "h-6 w-6 text-red-500 cursor-pointer",
  FlightPricePill: "bg-blue-100 text-blue-800 text-xs font-medium px-2.5 py-0.5 rounded-full",
  FlightDatesCtn: "flex w-full border border-slate-400 rounded-lg",
  FlightDatesIcon: "inline align-bottom h-5 w-5 text-gray-500",
  FlightDatesLabel: "inline-flex font-bold items-center px-3 text-sm text-slate-500 rounded-l-md",
  FlightDatesInput: "block flex-1 min-w-0 p-2.5 border-0 rounded-none rounded-r-lg text-gray-900 text-sm w-full",
  FlightStopIcon: "inline h-2 w-2 text-slate-700",
  FlightStopTimelineCtn: "relative border-l border-dashed border-slate-300 mt-4 ml-2",
  FlightStopTimelineIcon: "absolute flex mt-1 text-indigo-200 items-center justify-center w-4 h-4 -left-2 bg-indigo-400 rounded-full ring-8 ring-gray-100",
  FlightStopTimelineTime: "mb-1 font-medium text-slate-900",
  FlightsStopTimelineText: "mb-2 text-sm font-normal leading-none text-slate-400",
  FlightsStopLayoverText: "mb-1 text-sm font-normal leading-none text-red-700",
  FlightsStopHR: "w-48 h-1 mx-auto my-4 bg-gray-100 border-0 rounded md:my-10",
};

export const FlightsModalCss = {
  CabinClassDropdownBtn: "hover:bg-indigo-100 font-medium rounded-lg text-sm px-4 py-2.5 text-center inline-flex items-center",
  CabinClassDropdownIcon: "h-4 w-4 text-slate-700",
  CabinClassOptCtn: "z-10 w-44 rounded-lg bg-white shadow block absolute",
  CabinClassOptList: "z-10 w-44 rounded-lg bg-white shadow",
  CabinClassOpt: "block rounded-lg py-2 px-4 cursor-pointer hover:bg-indigo-100",
  AirportSearchOptCts: "z-10 rounded-lg bg-white shadow block absolute",
  AirportSearchOptList: "z-10 rounded-lg bg-white shadow",
  AirportSearchOpt: "block rounded-lg py-2 px-4 cursor-pointer hover:bg-indigo-100",
  FlightSearchBtn: "bg-indigo-500 font-medium rounded-lg text-sm text-white px-4 py-2.5 text-center inline-flex items-center",
  FlightFromIconCtn: "absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none",
  FlightFromIcon: "h-6 w-6 text-slate-700",
  FlightFromInput: "border border-slate-200 text-gray-900 text-sm rounded block w-full pl-10 p-4",
  FlightSearchHR: "h-px my-8 bg-gray-200 border-0",
  FlightSearchResultsTitle: "text-lg sm:text-2xl mb-2 font-medium text-slate-900",
  FlightTripCard: "flex p-4 rounded shadow-md hover:shadow-lg hover:shadow-indigo-100",
  FlightPlusIcon: "h-6 w-6 text-green-500 cursor-pointer",
  FlightPricePill: "bg-blue-100 text-blue-800 text-xs font-medium px-2.5 py-0.5 rounded-full",
  FlightDatesCtn: "flex w-full border border-slate-200 rounded-lg",
  FlightDatesIcon: "inline align-bottom h-5 w-5 text-gray-500",
  FlightDatesLabel: "inline-flex font-bold items-center px-3 text-sm text-slate-500 rounded-l-md",
  FlightDatesInput: "block flex-1 min-w-0 p-2.5 border-0 rounded-none rounded-r-lg text-gray-900 text-sm w-full",
  FlightStopIcon: "inline h-2 w-2 text-slate-700",
  FlightStopTimelineCtn: "relative border-l border-dashed border-slate-300 mt-4 ml-2",
  FlightStopTimelineIcon: "absolute flex mt-1 text-indigo-200 items-center justify-center w-4 h-4 -left-2 bg-indigo-200 rounded-full ring-8 ring-white",
  FlightStopTimelineTime: "mb-1 font-medium text-slate-900",
  FlightsStopTimelineText: "mb-2 text-sm font-normal leading-none text-slate-800",
  FlightsStopLayoverText: "mb-1 text-sm font-normal leading-none text-red-800",
  FlightsStopHR: "w-48 h-1 mx-auto my-4 bg-gray-100 border-0 rounded md:my-10",
}