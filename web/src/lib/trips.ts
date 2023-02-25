import _find from "lodash/find";
import _get from "lodash/get";

import { Auth } from "./auth";
import { Common } from './common';
import { Flights } from './flights';
import { Maps } from './maps';

import CameraIcon from '../components/icons/fill/CameraIcon';
import CupIcon from '../components/icons/fill/CupIcon';
import DiningIcon from '../components/icons/fill/DiningIcon';
import FlightIcon from '../components/icons/fill/FlightIcon';
import HotelIcon from '../components/icons/fill/HotelIcon';
import MapPinIcon from '../components/icons/fill/MapPinIcon';
import NatureIcon from '../components/icons/fill/NatureIcon';
import ShoppingIcon from '../components/icons/fill/ShoppingIcon';


export namespace Trips {
  export interface Member {
    id: string
    role: "creator" | "collaborator" | "participant"
    labels: {[key: string]: string}
  }

  interface BaseTransit {
    id: string
    type: string

    confirmationID?: string
    notes?: string
    price?: Common.Price
    tags: {[key: string]: string}
    labels: {[key: string]: string}
  }

  export interface Flight extends BaseTransit {
    itineraryType: string
    depart: Flights.Flight
    return: Flights.Flight
  }

  export interface Lodging {
    id: string
    numGuests?: number
    checkinTime?: Date | string
    checkoutTime?: Date | string
    price: Common.Price
    confirmationID?: string
    notes?: string
    place: Maps.Place
    tags: {[key: string]: string}
    labels: {[key: string]: string}
  }

  export interface Content {
    id: string
    title: string
    place: Maps.Place
    notes: string
    labels: {[key: string]: string}
    comments: any
  }

  export interface ContentList {
    id: string
    name?: string
    contents: Array<Content>
    labels: {[key: string]: string}
  }

  export interface ItineraryContent {
    id: string
    tripContentListId: string
    tripContentId: string
    price: Common.Price
    startTime?: string | Date
    endTime?: string | Date
    labels: {[key: string]: string}
  }

  export interface ItineraryList {
    id: string
    desc: string
    date: string | Date
    contents: Array<ItineraryContent>
    routes: Array<any>,
    labels: {[key: string]: string}
  }

  export interface Budget {
    amount: Common.Price
    items: Array<BudgetItem>
    labels: {[key: string]: string}
    tags: {[key: string]: string}
  }

  export interface BudgetItem {
    id: string
    title: string
    desc: string
    price: Common.Price
    labels: {[key: string]: string}
    tags: {[key: string]: string}
  }
}

export const MemberRoleCreator = "creator";
export const MemberRoleCollaborator = "collaborator";
export const MemberRoleParticipant = "participant";
export const LabelTransportModePref = "transportationPreference";
export const DefaultTransportModePref = "walk+drive";
export const PriceAmountPath = "price.amount";
export const BudgetAmountJSONPath = "amount/amount";
export const PriceAmountJSONPath = "price/amount";
export const LabelContentItineraryDates = "itinerary|dates";
export const LabelContentItineraryDatesJSONPath = "labels/itinerary|dates";
export const LabelDelimiter = "|";
export const LabelContentListColor = "ui|color";
export const LabelContentListColorJSONPath = "labels/ui|color";
export const LabelContentListIcon = "ui|icon";
export const LabelContentListIconJSONPath = "labels/ui|icon";
export const LabelFractionalIndex = "fIndex";
export const DefaultContentColor = "rgb(203 213 225)";
export const ContentColorOpts = ["rgb(74 222 128)", "rgb(34 211 238)",  "rgb(96 165 250)","rgb(129 140 248)",  "rgb(232 121 249)","rgb(244 114 182)", "rgb(248 113 113)", "rgb(251 146 60)", "rgb(253 224 71)", "rgb(161 98 7)"];
export const ContentIconOpts = {
  "flight": FlightIcon,
  "hotel": HotelIcon,
  "camera": CameraIcon,
  "coffee": CupIcon,
  "forkspoon": DiningIcon,
  "nature": NatureIcon,
  "pin": MapPinIcon,
  "shopping": ShoppingIcon,
} as {[key: string]: any}


// Member Helpers
export const userFromMemberID = (member: Trips.Member, userMap: any): Auth.User | undefined => {
  const memberID = member.id;
  return _get(userMap, memberID);
}

// Flights Helpers

export const flightItineraryType = (flight: Trips.Flight) => {
  return flight.itineraryType;
}

export const flilghtPriceAmt = (l: Trips.Flight) => {
  return _get(l, PriceAmountPath, 0)
}

// Lodging Helpers

export const lodgingPriceAmt = (l: Trips.Lodging) => {
  return _get(l, PriceAmountPath, 0)
}

// Trip Content Helpers

export const tripContentColor = (l: Trips.ContentList| Trips.ItineraryList) => {
  return _get(l, LabelContentListIconJSONPath, DefaultContentColor);
}


// Itinerary Content Helpers
export const itineraryContentPriceAmt = (ctnt: Trips.ItineraryContent) => {
  return _get(ctnt, PriceAmountPath, 0);
}

export const tripContentForItineraryContent = (trip: any, contentListID: string, contentID: string) => {
  return _find(
    trip.contents[contentListID].contents,
    (c: any) => c.id === contentID);
}

// Budget Helpers

export const budgetAmt = (budget: Trips.Budget) => {
  return _get(budget, "amount.amount", 0)
}

export const budgetItemPriceAmt = (bi: Trips.BudgetItem) => {
  return _get(bi, PriceAmountPath, 0)
}
