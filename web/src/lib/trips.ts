import _find from "lodash/find";
import _get from "lodash/get";
import { v4 as uuidv4 } from 'uuid';

import { jsonPathToPath } from "./jsonpatch";

import CameraIcon from '../components/icons/fill/CameraIcon';
import CupIcon from '../components/icons/fill/CupIcon';
import DiningIcon from '../components/icons/fill/DiningIcon';
import FlightIcon from '../components/icons/fill/FlightIcon';
import HotelIcon from '../components/icons/fill/HotelIcon';
import MapPinIcon from '../components/icons/fill/MapPinIcon';
import NatureIcon from '../components/icons/fill/NatureIcon';
import ShoppingIcon from '../components/icons/fill/ShoppingIcon';
import { User } from "./auth";
import { Price } from "./common";
import { Place } from "./maps";


export const MemberRoleCreator = "creator";
export const MemberRoleCollaborator = "collaborator";
export const MemberRoleParticipant = "participant";

export const DefaultTransportModePref = "walk+drive";
export const DefaultContentColor = "rgb(203 213 225)";

export const LabelDelimiter = "|";
export const LabelTransportModePref = "transportationPreference";
export const LabelItineraryDates = "itinerary|dates";
export const LabelUiColor = "ui|color";
export const LabelUiIcon = "ui|icon";
export const LabelFractionalIndex = "fIndex";

export const JSONPathPriceAmount = "price/amount";
export const JSONPathBudgetAmount = "amount/amount";
export const JSONPathLabelItineraryDates = "labels/itinerary|dates";
export const JSONPathLabelUiColor = "labels/ui|color";
export const JSONPathLabelUiIcon = "labels/ui|icon";


export const ContentColorOpts = [
  "rgb(74 222 128)",
  "rgb(34 211 238)",
  "rgb(96 165 250)",
  "rgb(129 140 248)",
  "rgb(232 121 249)",
  "rgb(244 114 182)",
  "rgb(248 113 113)",
  "rgb(251 146 60)",
  "rgb(253 224 71)",
  "rgb(161 98 7)"
];
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


export type Role = "creator" | "collaborator" | "participant";

export interface Member {
  id: string
  role: Role
  labels: {[key: string]: string}
};


export const userFromMember = (member: Member, userMap: any): User | undefined => {
  return _get(userMap, member.id);
}

interface BaseTransit {
  id: string
  type: string

  confirmationID?: string
  notes?: string
  price?: Price
  tags: {[key: string]: string}
  labels: {[key: string]: string}
}

export interface Flight extends BaseTransit {
  itineraryType: string
  depart: Flight
  return: Flight
}

export const getFlightItineraryType = (flight: Flight) => {
  return flight.itineraryType;
}

export const getFlilghtPriceAmt = (l: Flight) => {
  return _get(l, jsonPathToPath(JSONPathPriceAmount), 0)
}

export const makeOnewayFlight = (depart: any, bookingMetadata: any) => {
  return {
    id: uuidv4(),
    type: "flight",
    tags: {},
    labels: {},
    itineraryType: "oneway",
    depart,
    return: {} as any,
    price: bookingMetadata.price,
  }
}

export const makeRoundTripFlight = (departFlight: any, returnFlight: any, bookingMetadata: any) => {
  return {
    id: uuidv4(),
    type: "flight",
    tags: {},
    labels: {},
    itineraryType: "roundtrip",
    depart: departFlight,
    return: returnFlight,
    price: bookingMetadata.price,
  }
}


export interface Lodging {
  id: string
  numGuests?: number
  checkinTime?: Date | string
  checkoutTime?: Date | string
  price: Price
  confirmationID?: string
  notes?: string
  place: Place
  tags: {[key: string]: string}
  labels: {[key: string]: string}
}


export const getLodgingPriceAmt = (l: Lodging) => {
  return _get(l, jsonPathToPath(JSONPathPriceAmount), 0)
}


export interface Content {
  id: string
  title: string
  place: Place
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
  price: Price
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

export const getContentColor = (l: ContentList | ItineraryList) => {
  return _get(l, jsonPathToPath(JSONPathLabelUiIcon), DefaultContentColor);
}

export const getItineraryContentPriceAmt = (ctnt: ItineraryContent) => {
  return _get(ctnt, jsonPathToPath(JSONPathPriceAmount), 0);
}

export const getTripContentForItineraryContent = (trip: any, contentListID: string, contentID: string) => {
  return _find(trip.contents[contentListID].contents, (c: any) => c.id === contentID);
}

export interface Budget {
  amount: Price
  items: Array<BudgetItem>
  labels: {[key: string]: string}
  tags: {[key: string]: string}
}

export interface BudgetItem {
  id: string
  title: string
  desc: string
  price: Price
  labels: {[key: string]: string}
  tags: {[key: string]: string}
}

export const getBudgetAmt = (budget: Budget) => {
  return _get(budget, jsonPathToPath(JSONPathBudgetAmount), 0)
}

export const getBudgetItemPriceAmt = (bi: BudgetItem) => {
  return _get(bi, jsonPathToPath(JSONPathPriceAmount), 0)
}
