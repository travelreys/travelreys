import axios from 'axios';
import _get from 'lodash/get';
import _filter from 'lodash/filter';

import { Common, BASE_URL } from './common';

import airports from '../assets/airports.json';
import { Trips } from './trips';

export namespace Flights {
  export interface Airline {
    name: string
    code: string
    websiteURL?: string
    phoneNumber?: string
  }
  export interface Airport extends Common.Positioning {
    code: string
  }


  export type ItineraryType = "roundtrip" | "oneway"
  export type CabinClass = "economy" | "premiumeconomy" | "business" | "first"

  export interface Itineraries {
    type: ItineraryType
    oneways?: OnewayList
    roundtrips?: RoundTripMap
  }
  export interface Oneway {
    depart: Flight
    bookingMetadata: BookingMetadata
  }
  export type OnewayList = Array<Oneway>
  export interface RoundTrip {
    depart: Flight
    returns: FlightsList
    bookingMetadata: BookingMetadataList
  }
  export type RoundTripMap = Map<string, RoundTrip>
  export interface BookingMetadata {
    score: number
    price: Common.PriceMetadata
    bookingURL: string
    bookingDeeplinkURL: string
  }
  export type BookingMetadataList = Array<BookingMetadata>
  export interface Flight {
    id: string
    departure: Departure
    arrival: Arrival
    numStops: number
    duration: number
    legs: LegsList
  }
  export type FlightsList = Array<Flight>

  export interface Leg {
    flightNo: string
    departure: Departure
    arrival: Arrival
    duration: number
    operatingAirline: Airline
  }
  export type LegsList = Array<Leg>

  export interface Departure {
    airport: Airport
    datetime: Date
  }
  export interface Arrival {
    airport: Airport
    datetime: Date
  }
}

const FlightsAPI = {
  search: (
    origIATA: string,
    destIATA: string,
    departDate: string,
    returnDate: string | undefined,
    cabinClass: string
  ) => {
    const url = `${BASE_URL}/api/v1/flights/search`;
    return axios.get(url, {
      params: {
        numAdults: '1',
        currency: 'SGD',
        origIATA: origIATA.toUpperCase(),
        destIATA: destIATA.toUpperCase(),
        departDate,
        returnDate,
        cabinClass,
       }
    });
  },

  airportAutocomplete: (q: string) => {
    return _filter(airports, (a) => a.airport.includes(q.toLowerCase())).slice(0, 10);
  }
};

export default FlightsAPI;

export const logoFallbackImg = "https://cdn-icons-png.flaticon.com/512/4353/4353032.png";

export const FlightDirectionDepart = "depart";
export const FlightDirectionReturn = "return";
export const FlightItineraryTypeOneway = "oneway";
export const FlightItineraryTypeRoundtrip = "roundtrip";

export const flightLogoUrl = (iata: string) => {
  return `https://www.gstatic.com/flights/airline_logos/70px/${iata}.png`;
}

export const flightDepartureTime = (flight: Flights.Flight) => {
  return _get(flight, "departure.datetime", "");
}

export const flightArrivalTime = (flight: Flights.Flight) => {
  return _get(flight, "arrival.datetime", "");
}

export const flightLegs = (flight: Flights.Flight) => {
  return _get(flight, "legs", []);
}

export const flightLegOpAirline = (leg: Flights.Leg) => {
  return _get(leg, "operatingAirline", {});
}






