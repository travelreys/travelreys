// Common

export namespace Common {
  export interface PriceMetadata {
    amount: number
    currency: string
  }
  export interface Positioning {
    name: string
    address?: string
    continent?: string
    country?: string
    state?: string
    city?: string
    longitude?: string
    latitude?: string
  }
}

// Flights
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


export namespace Maps {
  export type Place = any
}


// Trips

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
    numGuests: number
    checkinTime: Date
    checkoutTime: Date
    priceMetadata: Common.PriceMetadata
    confirmationID?: string
    notes?: string
    place: Maps.Place
    tags: Map<string, string>
    labels: Map<string, string>
  }

  export interface TripContent {
    id: string
    title: string
    place: Maps.Place
    notes: string
    labels: Map<string, string>
  }

  export interface TripContentList {
    id: string
    name: string
    contents: Array<TripContent>
  }

  export interface TripContent {
    id: string
    title: string
    place: Maps.Place
    notes: string
    labels: Map<string, string>
  }

}






