////////////
// Common //
////////////

export const CommonCss = {
  Icon: "h-4 w-4",
  LeftIcon: "h-4 w-4 mr-2",
  DropdownIcon: "h-4 w-4 mt-1",
}

//////////////////////
// InputDatesPicker //
//////////////////////

export const InputDatesPickerCss = {
  Ctn: "flex w-full border border-slate-200 rounded-lg mr-2",
  Icon: "inline align-bottom h-5 w-5 text-gray-500",
  Label: "inline-flex bg-gray-200 font-bold items-center px-3 text-sm text-slate-500 rounded-l-md",
  Input: "block flex-1 min-w-0 p-2.5 border-0 rounded-none rounded-r-lg text-gray-900 text-sm w-full",
}

///////////
// Modal //
///////////

export const ModalCss = {
  Container: "relative z-20",
  Inset: "fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity",
  Content: "fixed inset-0 z-10 overflow-y-auto",
  ContentContainer: "flex min-h-full flex-col p-4 text-center sm:items-center sm:p-0",
  ContentCard: "bg-white relative transform rounded-lg text-left shadow-xl transition-all sm:my-8 sm:w-full sm:max-w-2xl",
}

////////////////////////
// Place Autocomplete //
////////////////////////

export const PlaceAutocompleteCss = {
  AutocompleteCtn: "p-1 bg-white absolute left-0 z-30 w-full border border-slate-200 rounded-lg",
  PredictionWrapper: "flex items-center mb-4 cursor-pointer group",
  Icon: "h-6 w-6",
  IconCtn: "p-1 group-hover:text-indigo-500",
  PrimaryTxt : "text-slate-900 group-hover:text-indigo-500 text-sm font-medium",
  SecondaryTxt: "text-slate-400 group-hover:text-indigo-500 text-xs",
}


///////////
// Trips //
///////////

export const CreateTripModalCss = {
  CreateModalCard: "bg-white rounded-lg px-4 pt-5 pb-4 sm:p-8 sm:pb-4",
  CreateTripTitle: "text-2xl font-bold text-center leading-6 text-slate-900 mb-6",
  TripNameCtn: "flex mb-4 border border-slate-200 rounded-lg",
  TripNameLabel: "inline-flex font-bold items-center px-3 text-sm text-slate-500 rounded-l-md",
  TripNameInput: "block flex-1 border-0 rounded-r-lg min-w-0 p-2.5 text-gray-900 text-sm w-full",
  TripDatesCtn: "flex w-full border border-slate-200 rounded-lg",
  TripDatesIcon: "inline align-bottom h-5 w-5 text-gray-500",
  TripDatesLabel: "inline-flex font-bold items-center px-3 text-sm text-slate-500 rounded-l-md",
  TripDatesInput: "block flex-1 min-w-0 p-2.5 border-0 rounded-none rounded-r-lg text-gray-900 text-sm w-full",
  CreateTripBtnsCtn: "bg-gray-50 px-4 pt-3 pb-5 rounded-b-lg sm:flex sm:flex-row-reverse sm:px-6",
  CreateTripBtn: "inline-flex w-full justify-center rounded-md border border-transparent bg-indigo-600 px-4 py-2 text-base font-medium text-white shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 sm:ml-3 sm:w-auto sm:text-sm",
  CreateTripCancelBtn: "mt-3 inline-flex w-full justify-center rounded-md border border-gray-300 bg-white px-4 py-2 text-base font-medium text-gray-700 shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm",
}

///////////////
// Trip Page //
///////////////

export const TripMenuCss = {
  TripMenuCtn: "min-h-screen w-full z-50 sm:w-1/2  sm:max-w-lg sm:shadow-xl sm:shadow-slate-900",
  TripMenu: "pb-40 w-full",
  TripMenuNav: "p-3 font-bold text-indigo-500",
  TabsCtn: "sticky top-0 z-10 bg-indigo-100 py-8 pb-4 mb-4",
  TabsWrapper: "bg-white rounded-lg p-5 mx-4 mb-4",
  TabItemCtn: "flex flex-row justify-around mx-2",
  TabItemBtn: "mx-4 my-2 flex flex-col items-center",
  TabItemBtnTxt: "text-slate-400 text-sm",
}

///////////
// Jumbo //
///////////

export const TripMenuJumboCss = {
  TripDatesBtn: "font-medium text-md text-slate-500",
  TripDatesBtnIcon: "inline h-5 w-5 align-sub mr-2",
  TripCoverImage: "block sm:max-h-96 w-full",
  TripImageEditIconCtn: "absolute top-4 right-4 h-10 w-10 bg-gray-800/70 p-2 text-center rounded-full",
  TripImageEditIcon: "h-6 w-6 text-white",
  TripNameInputCtn: "h-16 relative -top-24",
  TripNameInputWrapper: "bg-white rounded-lg shadow p-5 mx-4 mb-4",
  TripNameInput: "text-2xl sm:text-4xl font-bold text-slate-700 w-full rounded-lg p-1 border-0 hover:bg-slate-300 hover:border-0 hover:bg-slate-100 focus:ring-0",
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
  SettingsBtn: "flex items-center w-full text-left hover:text-indigo-500",
}

///////////
// Notes //
///////////

export const TripNotesCss = {
  TitleCtn: "flex justify-between mb-4",
  HeaderCtn: "text-2xl sm:text-3xl font-bold text-slate-700",
  ToggleBtn: "mr-2"
}

///////////////
// Logistics //
///////////////

export const TripLogisticsCss = {
  FlightDatesCtn: "flex w-full border border-slate-400 rounded-lg",
  FlightDatesIcon: "inline align-bottom h-5 w-5 text-gray-500",
  SearchFlightBtn: "text-slate-500 text-sm mt-1 font-bold",
  FlightDatesInput: "block flex-1 min-w-0 p-2.5 border-0 rounded-none rounded-r-lg text-gray-900 text-sm w-full",
  FlightDatesLabel: "inline-flex font-bold items-center px-3 text-sm text-slate-500 rounded-l-md",
  FlightPricePill: "bg-blue-100 text-blue-800 text-xs font-medium px-2.5 py-0.5 rounded-full",
  FlightsHeaderCtn: "text-2xl sm:text-3xl font-bold text-slate-700",
  FlightsStopHR: "w-48 h-1 mx-auto my-4 bg-gray-100 border-0 rounded md:my-10",
  FlightsStopLayoverText: "mb-1 text-sm font-normal leading-none text-red-700",
  FlightsStopTimelineText: "mb-2 text-sm font-normal leading-none text-slate-400",
  FlightsTitleCtn: "flex justify-between mb-4",
  FlightsToggleBtn: "mr-2",
  FlightStopIcon: "inline h-2 w-2 text-slate-700",
  FlightStopTimelineCtn: "relative border-l border-dashed border-slate-300 mt-4 ml-2",
  FlightStopTimelineIcon: "absolute flex mt-1 text-indigo-200 items-center justify-center w-4 h-4 -left-2 bg-indigo-400 rounded-full ring-8 ring-gray-100",
  FlightStopTimelineTime: "mb-1 font-medium text-slate-900",
  FlightTransit: "flex p-4 bg-slate-50 rounded-lg shadow-md",
  FlightTransitCard: "mb-4",
  FlightTransitIcon: "h-6 w-6 text-red-500 cursor-pointer",
};

/////////////
// Flights //
/////////////

export const FlightsModalCss = {
  Ctn: "px-4 pt-5 sm:p-8 sm:pb-2 rounded-t-lg mb-4",
  Wrapper: "flex justify-between mb-6",
  Header: "text-xl sm:text-2xl font-bold text-center text-slate-900",
  CloseIcon: "h-6 w-6 text-slate-700",
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
  RoundTripStepperCtn: "flex items-center w-full mb-4 space-x-2 text-sm font-medium text-center text-gray-500 bg-white rounded-lg sm:p-4 sm:space-x-4",
  RoundTripStepperActive: "flex items-center text-blue-600 dark:text-blue-500",
  RoundTripStepper: "flex items-center",
  RoundTripStepperText: "flex items-center",
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

/////////////
// Lodging //
/////////////

export const LodgingsModalCss = {
  Ctn: "px-4 pt-5 sm:p-8 sm:pb-2 rounded-t-lg mb-4",
  Wrapper: "flex justify-between mb-6",
  Header: "text-xl sm:text-2xl font-bold text-center text-slate-900",
  CloseIcon: "h-6 w-6 text-slate-700",
  SearchIconCtn: "absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none",
  SearchIcon: "h-6 w-6 text-slate-700",
  SearchInput: "border border-slate-200 text-gray-900 text-sm rounded block w-full pl-10 p-4",
  InputDatesCtn: "flex w-full border border-slate-200 rounded",
  AddBtn: "mt-2 bg-indigo-500 font-medium rounded-lg text-sm text-white px-4 py-2.5",
  PredictionsCtn: "flex items-center mb-4 cursor-pointer group",
  PredictionIconCtn: "p-1 group-hover:text-indigo-500",
  PredictionMain: "text-slate-900 group-hover:text-indigo-500 text-sm font-medium",
  PredictionSecondary: "text-slate-400 group-hover:text-indigo-500 text-xs",
}

export const LodgingCardCss = {
  Ctn: "p-4 bg-slate-50 rounded-lg shadow-md mb-4",
  Header: "font-bold text-sm mb-1",
  AddrTxt: "text-slate-600 text-sm text-left flex items-center mb-1 hover:text-indigo-500",
  WebsiteTxt: "text-indigo-500 text-sm flex items-center",
  PhoneTxt: "text-slate-600 text-sm flex items-center mb-1",
  DatesTxt: "text-slate-600 text-sm flex items-center mb-2 cursor-pointer",
  PricePill: "bg-blue-100 text-blue-800 text-xs font-semibold px-2.5 py-0.5 rounded-full mb-2 w-fit cursor-pointer",
  DatesPickerCtn: "flex w-full rounded",
  PriceInputCtn: "flex w-full rounded mb-2",
  DeleteBtn: "text-red-500 flex items-center",
}

//////////////////
// Trip Content //
//////////////////

export const TripContentSectionCss = {
  HeaderCtn: "flex justify-between mb-4",
  Header: "text-2xl sm:text-3xl font-bold text-slate-700",
  AddBtn: "text-white py-2 px-4 bg-indigo-500 rounded-lg text-sm font-semibold",
  Hr: "w-48 h-1 m-5 mx-auto bg-gray-300 border-0 rounded",
  ToggleBtn: "mr-2",
}

export const TripContentListCss = {
  Ctn: "rounded-lg shadow-xs mb-4",
  ChooseColorBtn: "flex items-center w-full text-left",
  DeleteBtn: "text-red-500 flex items-center w-full text-left",
  NameInput: "p-0 w-full text-xl mb-1 sm:text-2xl font-bold text-gray-800 placeholder:text-gray-400 rounded border-0 hover:bg-gray-300 hover:border-0 focus:ring-0 focus:p-1 duration-500",
  NewContentCtn: "flex my-4 w-full",
  NewContentInput: "flex-1 mr-1 text-md sm:text-md font-bold text-gray-800 placeholder:font-normal placeholder:text-gray-300 placeholder:italic rounded border-0 bg-gray-100 hover:border-0 focus:ring-0",
  NewContentBtn: "text-green-600 w-1/12 hover:bg-green-50 rounded-lg text-sm font-bold inline-flex justify-around items-center",
}

export const TripContentCss = {
  AutocompleteCtn: "p-1 bg-white absolute left-0 z-30 w-full border border-slate-200 rounded-lg",
  Ctn: "bg-slate-50 rounded-lg shadow-xs mb-4 p-4 relative",
  DeleteBtn: "text-red-500 flex items-center w-full text-left",
  ItineraryDateBtn: "flex items-center w-full justify-between hover:text-indigo-500 text-align-right",
  TitleInput: "p-0 mb-1 font-bold text-gray-800 bg-transparent placeholder:text-gray-400 rounded border-0 hover:border-0 focus:ring-0 duration-400",
  WebsiteLink: "flex items-center mb-1",
  WebsiteTxt: "text-indigo-500 text-sm flex items-center",
  AddItineraryBtn: "text-xs text-gray-800 font-bold bg-indigo-200 rounded-full px-2 py-1 hover:bg-indigo-400",
  ItineraryBadge: "bg-indigo-100 text-indigo-800 text-xs font-medium mr-2 px-2.5 py-0.5 rounded",
  PlaceCtn: "text-slate-600 text-sm flex items-center mb-1 hover:text-indigo-500",
  PlaceInput: "p-0 mb-1 text-sm text-gray-600 bg-transparent placeholder:text-gray-400 rounded border-0 hover:border-0 focus:ring-0 duration-400",
}

////////////////////
// Trip Itinerary //
////////////////////

export const TripItinerarySectionCss = {
  Hr: "w-48 h-1 m-5 mb-8 mx-auto bg-gray-300 border-0 rounded",
  ToggleBtn: "mr-2",
}

export const TripItineraryListCss = {
  Ctn: "rounded-lg shadow-xs",
  NameInput: "p-0 w-full text-xl mb-1 sm:text-2xl font-bold text-gray-800 placeholder:text-gray-400 rounded border-0 hover:bg-gray-300 hover:border-0 focus:ring-0 focus:p-1 duration-500",
  ContentsCtn: "pl-6 py-4",
  ContentsWrapper: "relative border-l border-gray-200",
  LodgingCtn: "w-full mb-2",
  LodgingWrapper: "flex items-center w-full p-3 space-x-4 text-gray-800 divide-x divide-gray-200 rounded-lg shadow",
  LodgingIconWrapper: "bg-indigo-200 p-2 rounded-full",
  LodgingName: "flex-1 pl-4 text-sm font-normal",
  LodgingStatus: "pl-2 font-semibold text-sm",
  ItinItem: "mb-8 ml-6",
  ItinContentIcon: "absolute flex items-center justify-center w-6 h-6 rounded-full -left-3 ring-8 ring-white font-bold text-white text-sm",
  ChooseColorBtn: "flex items-center w-full text-left",

}

export const TripItineraryCss = {
  Ctn: "bg-slate-50 rounded-lg shadow-xs px-4 py-2 relative shadow",
  TitleInput: "p-0 mb-1 font-bold text-gray-800 bg-transparent placeholder:text-gray-400 rounded border-0 hover:border-0 focus:ring-0 duration-400",
  AutocompleteCtn: "p-1 bg-white absolute left-0 z-30 w-full border border-slate-200 rounded-lg",
  PredictionWrapper: "flex items-center mb-4 cursor-pointer group",
  WebsiteTxt: "text-indigo-500 text-sm flex items-center",
  DeleteBtn: "text-red-500 flex items-center w-full",
  PriceInputCtn: "flex w-full rounded mb-2",
  PricePill: "bg-blue-100 text-blue-800 text-xs font-semibold px-2.5 py-0.5 rounded-full mb-2 w-fit cursor-pointer",
}

/////////////
// TripMap //
/////////////

export const TripMapCss = {
  AddrTxt: "text-gray-600 flex items-center mb-1",
  BtnCtn: "flex items-center mt-6",
  Ctn: "fixed h-screen",
  DetailsCard: "bg-white p-4 mx-4 h-11/12 w-11/12 max-w-3xl rounded-xl pointer-events-auto",
  DetailsWrapper: "absolute bottom-0 mb-8 z-10 pointer-events-none",
  GmapBtn: "flex w-fit rounded-full py-2 px-6 items-center border border-gray-200 font-semibold text-gray-500",
  HeaderCtn: "flex justify-between items-center",
  OpeningHrsTxt: "flex text-gray-600 items-center truncate",
  PhoneBtn: "flex w-fit rounded-full py-2 px-6 mr-2 items-center border border-gray-200 font-semibold text-gray-500",
  PhoneIcon: "h-4 w-4 text-indigo-500 mr-2",
  RatingsStar: "text-yellow-500 flex items-center mb-1",
  RatingsTxt: "text-gray-600",
  SummaryTxt: "text-gray-600 mb-1",
  TitleCtn: "font-bold text-lg flex items-center",
  WeekdayTxt: "text-slate-600 ml-6",
}


