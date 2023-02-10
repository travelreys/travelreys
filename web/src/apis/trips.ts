import axios from 'axios'

import { Common, BASE_URL } from './common';
import { Flights } from './flights';
import { Maps } from './maps';

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
    labels: Map<string, string>
  }
}

export const LabelContentItineraryDates = "itinerary|dates";
export const LabelContentItineraryDatesJSONPath = "labels.itinerary|dates";
export const LabelContentItineraryDatesDelimeter = "|";



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







