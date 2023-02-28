import _find from "lodash/find";
import _get from "lodash/get";
import _sortBy from 'lodash/sortBy';
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
export const DefaultActivityColor = "rgb(203 213 225)";

export const LabelDelimiter = "|";
export const LabelFractionalIndex = "fIndex";
export const LabelLocked = "locked";
export const LabelTransportModePref = "transportationPreference";
export const LabelItineraryDates = "itinerary|dates";
export const LabelUiColor = "ui|color";
export const LabelUiIcon = "ui|icon";

export const JSONPathPriceAmount = "price/amount";
export const JSONPathBudgetAmount = "amount/amount";
export const JSONPathLabelItineraryDates = "labels/itinerary|dates";
export const JSONPathLabelLocked = "labels/locked";
export const JSONPathLabelUiColor = "labels/ui|color";
export const JSONPathLabelUiIcon = "labels/ui|icon";


export const ActivityColorOpts = [
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
export const ActivityIconOpts = {
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


export interface Activity {
  id: string
  title: string
  place: Place
  notes: string
  labels: {[key: string]: string}
  comments: any
}

export const makeActivity = (title: string) => {
  return {
      id: uuidv4(),
      title: title,
      notes: "",
      place: {},
      labels: {},
      comments: [],
    }
}

export interface ActivityList {
  id: string
  name?: string
  activities: Array<Activity>
  labels: {[key: string]: string}
}

export const makeActivityList = () => {
  return { id: uuidv4(), name: "", activities: [], labels: {}}
}

export interface ItineraryActivity {
  id: string
  activityListId: string
  activityId: string
  price: Price
  startTime?: string | Date
  endTime?: string | Date
  labels: {[key: string]: string}
}

export const makeItineraryActivity = (actId: string, actListId: string, fIndex: string) => {
  return {
    id: uuidv4(),
    activityId: actId,
    activityListId: actListId,
    price: {} as any,
    labels: {fIndex} as any,
  }
}

export const getfIndex = (act: ItineraryActivity) => {
  return _get(act, `labels.${LabelFractionalIndex}`);
}

export const getSortedActivies = (l: ItineraryList) => {
  const activities = Object.values(_get(l, "activities", {}));
  return _sortBy(activities, (act: ItineraryActivity) => getfIndex(act))
}

export interface ItineraryList {
  id: string
  desc: string
  date: string | Date
  activities: {[key: string]: ItineraryActivity}
  routes: Array<any>,
  labels: {[key: string]: string}
}

export const getActivityColor = (l: ActivityList | ItineraryList) => {
  return _get(l, jsonPathToPath(JSONPathLabelUiColor), DefaultActivityColor);
}

export const getActivityIcon = (l: ActivityList | ItineraryList) => {
  return _get(l, jsonPathToPath(JSONPathLabelUiIcon));
}

export const isListLock = (l: ActivityList | ItineraryList) => {
  return _get(l, jsonPathToPath(JSONPathLabelLocked)) === "true";
}

export const getItineraryActivityPriceAmt = (act: ItineraryActivity) => {
  return _get(act, jsonPathToPath(JSONPathPriceAmount), 0);
}

export const getTripActivityForItineraryActivity = (trip: any, itinAct: ItineraryActivity) => {
  return trip.activities[itinAct.activityListId].activities[itinAct.activityId];
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
