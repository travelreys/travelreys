import axios from 'axios';
import _find from "lodash/find";
import _get from "lodash/get";

import { Common, BASE_URL } from './common';
import { Flights } from './flights';
import { Maps } from './maps';

import NatureIconFill from '../components/icons/NatureIconFill';
import CupIconFill from '../components/icons/CupIconFill';
import DiningIconFill from '../components/icons/DiningIconFill';
import ShoppingIconFill from '../components/icons/ShoppingIconFill';
import CameraIconFill from '../components/icons/CameraIconFill';
import MapPinIconFill from '../components/icons/MapPinIconFill';
import HotelIcon from '../components/icons/HotelIcon';

export namespace Trips {
  interface BaseTransit {
    id: string
    type: string

    confirmationID?: string
    notes?: string
    price?: Common.PriceMetadata
    tags: Map<string, string>
    labels: Map<string, string>
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
    priceMetadata: Common.PriceMetadata
    confirmationID?: string
    notes?: string
    place: Maps.Place
    tags: Map<string, string>
    labels: Map<string, string>
  }

  export interface Content {
    id: string
    title: string
    place: Maps.Place
    notes: string
    labels: Map<string, string>
    comments: any
  }

  export interface ContentList {
    id: string
    name?: string
    contents: Array<Content>
    labels: Map<string, string>
  }

  export interface ItineraryContent {
    id: string
    tripContentListId: string
    tripContentId: string
    priceMetadata: Common.PriceMetadata
    startTime?: string | Date
    endTime?: string | Date
    labels: Map<string, string>
  }

  export interface ItineraryList {
    id: string
    desc: string
    date: string | Date
    contents: Array<ItineraryContent>
    routes: Array<any>,
    labels: Map<string, string>
  }

  export interface Budget {
    amount: Common.PriceMetadata
    items: Array<BudgetItem>
    labels: Map<string, string>
    tags: Map<string, string>
  }

  export interface BudgetItem {
    id: string
    title: string
    desc: string
    priceMetadata: Common.PriceMetadata
    labels: Map<string, string>
    tags: Map<string, string>
  }
}

export const LabelTransportationPreference = "transportationPreference";
export const DefaultTransportationPreference = "walk+drive";
export const PriceMetadataAmountPath = "priceMetadata.amount";
export const PriceMetadataAmountJSONPath = "priceMetadata/amount";
export const LabelContentItineraryDates = "itinerary|dates";
export const LabelContentItineraryDatesJSONPath = "labels/itinerary|dates";
export const LabelContentItineraryDatesDelimeter = "|";
export const LabelContentListColor = "ui|color";
export const LabelContentListColorJSONPath = "labels/ui|color";
export const LabelContentListIcon = "ui|icon";
export const LabelContentListIconJSONPath = "labels/ui|icon";
export const DefaultContentColor = "rgb(203 213 225)";
export const ContentColorOpts = ["rgb(74 222 128)", "rgb(34 211 238)",  "rgb(96 165 250)","rgb(129 140 248)",  "rgb(232 121 249)","rgb(244 114 182)", "rgb(248 113 113)", "rgb(251 146 60)", "rgb(253 224 71)", "rgb(161 98 7)"];
export const ContentIconOpts = {
  "hotel": HotelIcon,
  "camera": CameraIconFill,
  "coffee": CupIconFill,
  "forkspoon": DiningIconFill,
  "nature": NatureIconFill,
  "pin": MapPinIconFill,
  "shopping": ShoppingIconFill,
} as {[key: string]: any}



const TripsAPI = {
  createTrip: (name: string, startDate: Date | undefined, endDate: Date | undefined) => {
    const url = `${BASE_URL}/api/v1/trips`;
    return axios.post(url, {name, startDate, endDate})
  },

  readTrips: () => {
    const url = `${BASE_URL}/api/v1/trips`;
    return axios.get(url).then(res => res.data);
  },

  readTrip: (id: string | undefined) => {
    const url = `${BASE_URL}/api/v1/trips/${id}`;
    return axios.get(url).then(res => res.data);
  },
};

export default TripsAPI;

export const lodgingPriceAmt = (l: Trips.Lodging) => {
  return _get(l, PriceMetadataAmountPath, 0)
}

export const itineraryContentPriceAmt = (ctnt: Trips.ItineraryContent) => {
  return _get(ctnt, PriceMetadataAmountPath, 0);
}


export const tripContentForItineraryContent = (trip: any, contentListID: string, contentID: string) => {
  return _find(
    trip.contents[contentListID].contents,
    (c: any) => c.id === contentID);
}

export const tripContentColor = (l: Trips.ContentList| Trips.ItineraryList) => {
  return _get(l, LabelContentListIconJSONPath, DefaultContentColor);
}

export const budgetItemPriceAmt = (bi: Trips.BudgetItem) => {
  return _get(bi, PriceMetadataAmountPath, 0)
}




