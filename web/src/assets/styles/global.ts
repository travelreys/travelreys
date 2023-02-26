////////////
// Common //
////////////

export const CommonCss = {
  Icon: "h-4 w-4",
  IconLarge: "h-8 w-8",
  LeftIcon: "h-4 w-4 mr-2",
  DropdownIcon: "h-4 w-4",
  DropdownBtn: "flex items-center w-full text-left hover:text-indigo-500",
  DeleteBtn: "text-red-500 flex items-center w-full text-left hover:text-red-700",
  Navbar: "p-3 font-bold text-indigo-500",
  HrShort: 'w-48 h-1 m-5 mx-auto bg-gray-300 border-0 rounded',
}

export const InputCss = {
  Ctn: "flex w-full border border-slate-200 rounded-lg mr-2",
  Icon: "inline align-bottom h-5 w-5 text-gray-500",
  Label: "inline-flex bg-gray-200 font-bold items-center px-3 text-sm text-slate-500 rounded-l-md",
  Input: "block flex-1 min-w-0 p-2.5 border-0 rounded-none rounded-r-lg text-gray-900 text-sm w-full",
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

  FlightDatesInput: "block flex-1 min-w-0 p-2.5 border-0 rounded-none rounded-r-lg text-gray-900 text-sm w-full",
  FlightDatesLabel: "inline-flex font-bold items-center px-3 text-sm text-slate-500 rounded-l-md",
  FlightPricePill: "bg-blue-100 text-blue-800 text-xs font-semibold px-2.5 py-0.5 rounded-full",

  FlightsStopHR: "w-48 h-1 mx-auto my-4 bg-gray-100 border-0 rounded md:my-10",
  FlightsStopLayoverText: "mb-1 text-sm font-normal leading-none text-red-700",
  FlightsStopTimelineText: "mb-2 text-sm font-normal leading-none text-slate-400",

  FlightsToggleBtn: "mr-2",
  FlightStopIcon: "inline h-2 w-2 text-slate-700",
  FlightStopTimelineCtn: "relative border-l border-dashed border-slate-300 mt-4 ml-2",
  FlightStopTimelineIcon: "absolute flex mt-1 text-indigo-200 items-center justify-center w-4 h-4 -left-2 bg-indigo-400 rounded-full ring-8 ring-gray-100",
  FlightStopTimelineTime: "mb-1 font-medium text-slate-900",
  FlightTransit: "flex p-4 bg-slate-50 rounded-lg shadow-md",
  FlightTransitCard: "mb-4",
  FlightTransitNumStop: "cursor-pointer border-b border-slate-400",
  FlightTransitLogoImgWrapper: "h-8 w-8 mr-4",
  FlightTransitLogoImg: "h-8 w-8",
  FlightTransitDatetime: "text-sm text-slate-800",
  FlightTransitTime: "font-medium",
  FlightTransitLongArrow: "h-6 w-8",
  FlightTransitAirportCode: "text-xs text-slate-800",
  FlightTransitDuration: "text-xs text-slate-800 block mb-1",
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
  FlightTripCard: "flex p-4 rounded shadow-md hover:shadow-lg hover:shadow-indigo-100",
  FlightPlusIcon: "h-6 w-6 text-green-500 cursor-pointer stroke-2",
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
}




/////////////////
// Trip Budget //
/////////////////

export const TripBudgetCss = {
  ProgressBarCtn: "w-full py-4 pr-2",
  ProgressBarWrapper: "bg-gray-200 rounded-full h-1.5",
  ProgressBar: "bg-indigo-600 h-1.5 rounded-full",
  SummaryCtn: "bg-gray-100 shadow rounded-lg flex p-2 divide-x-2 mb-4",
  SpendingCtn: "flex-1 p-2",
  SpendingTitle: "text-sm font-bold mb-1",
  SpendingAmount: "text-4xl",
  OptsCtn: "flex flex-col p-2",
  AddExpenseBtn: "text-indigo-500 font-bold inline-flex items-center p-2 text-sm",
  EditBudgetBtn: 'text-indigo-500 font-bold inline-flex items-center p-2 text-sm',
  SubsectionTxt: "text-lg font-bold",
  ItemCtn: "flex justify-between items-center border-b py-4 border-gray-200",
  ItemDescCtn: "flex flex-1 items-center",
  FlightItemIcon: "bg-green-200 p-2 rounded-full mr-2",
  FlightItemAirport: "text-sm text-slate-600",
  LodgingItemIcon: "bg-orange-200 p-2 rounded-full mr-2",
  LodgingDatesTxt: "text-slate-600 text-sm flex items-center cursor-pointer",
  ItinItemIcon: "flex items-center justify-center w-8 h-8 p-2 rounded-full mr-2 text-white font-bold",
  ItemNameTxt: "font-bold",
  ItemDescTxt: "text-sm text-gray-500",
  ItemPriceTxt: "font-bold",
  PriceInputCtn: "flex w-full rounded-lg mb-3 border border-gray-200",
  PriceInputLabel: "inline-flex bg-gray-200 font-bold items-center px-3 text-sm text-slate-500 rounded-l-md",
  BudgetItemIcon: "flex items-center justify-center w-8 h-8 p-2 rounded-full mr-2 text-white font-bold bg-indigo-500"
}

///////////////////
// Trip Settings //
///////////////////

export const TripSettingsCss = {
  TransportModeLabel: "block mb-2 font-semibold text-gray-900",
  TransportModeSelect: "bg-white block w-full p-2.5 border border-gray-300 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500",

  MemberSectionCtn: "mb-8",
  MemberSectionHeader: "flex justify-between items-center",
  MemberSectionTitle: "font-bold text-2xl",
  SearchMemberBtn: "font-semibold text-gray-500",
  MemberSearchHeader: "flex justify-between items-center mb-8",
  MemberSearchHeaderTxt: "text-gray-800 font-bold text-xl",
  MemberRoleSelect: "mb-4 bg-white block w-full p-2.5 border border-gray-300 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500",
  MemberSearchIconCtn: "absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none",
  MemberSearchIcon: "h-6 w-6 text-slate-700",
  MemberSearchInput: "border border-slate-200 text-gray-900 text-sm rounded-lg block w-full pl-10 p-4",
  MemberSearchItem: "flex items-center w-full p-2 mb-4 text-left rounded-lg hover:shadow hover:shadow-indigo-200",
  MemberSearchItemAvatar: "inline-block h-10 w-10 mr-4",
  MemberSearchItemName: "font-semibold",
  MemberSearchItemDesc: "text-gray-500",
  MemberCtn: "flex items-center py-4 border-b border-gray-200",
  MemberAvatarDiv: "h-12 w-12 mr-4",
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


