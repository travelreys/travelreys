import _get from 'lodash/get';
import { makeCommonAxios } from './common';

export const placeFields = [
  "address_component",
  "adr_address",
  "business_status",
  "formatted_address",
  "geometry",
  "name",
  "photos",
  "place_id",
  "types",
  "utc_offset",
  "opening_hours",
  "formatted_phone_number",
  "international_phone_number",
  "website",
  "url",
];
export const placeAtmosphereFields = placeFields.concat([
  "editorial_summary",
  "price_level",
  "rating",
  "reviews",
  "user_ratings_total"
]);
export const ModeDriving = "driving";

export const EMBED_MAPS_APIKEY = "AIzaSyBaqenQ0nQVtkhnXBn-oWBtlPDL5uHmvNU";
export const PLACE_IMAGE_APIKEY = "AIzaSyBgNwirAT6TSS208emcC0Lbgex6i3EwhR0";

export interface AutocompleteResponse {
  predictions: any
  error?: string
}

export interface PlaceDetailsResponse {
  place: any
  errors?: string
}

const placeAutocomplete = (query: string, types: Array<string>, sessiontoken: string): Promise<AutocompleteResponse> => {
  const typesParam = types.join(",");
  return makeCommonAxios().get("/api/v1/maps/place/autocomplete", {
    params: { query, types: typesParam, sessiontoken }
  })
  .then((res) => {
    const predictions = _get(res, "data.predictions", []);
    return {predictions};
  })
  .catch((err) => {
    return {predictions: [], error: err.message};
  });
}

const placeDetails = (placeID: string, fields: Array<string>, sessiontoken?: string): Promise<PlaceDetailsResponse> => {
  const fieldsParam = fields.join(",");
  return makeCommonAxios().get("/api/v1/maps/place/details", {
    params: { placeID, fields: fieldsParam, sessiontoken }
  })
  .then((res) => {
    const place = _get(res, "data.place", []);
    return {place};
  })
  .catch((err) => {
    return {place: {}, error: err.message};
  });
}

const directions = (originPlaceID: string, destPlaceID: string, mode: string) => {
  return makeCommonAxios().get("/api/v1/maps/place/directions", {
    params: { originPlaceID, destPlaceID, mode }
  });
}

const optimizeRoute = (originPlaceID: string, destPlaceID: string, waypointsPlaceID: string) => {
  return makeCommonAxios().get("/api/v1/maps/place/optimize-route", {
    params: { originPlaceID, destPlaceID, waypointsPlaceID }
  });
}

export default {
  placeAutocomplete,
  placeDetails,
  directions,
  optimizeRoute,
};
