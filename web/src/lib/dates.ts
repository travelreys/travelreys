import _get from "lodash/get";
import _isEmpty from "lodash/isEmpty";
import {
  format,
  formatDuration,
  intervalToDuration,
  isEqual,
  parseISO as _parseISO,
  parseJSON,
} from 'date-fns';
import * as Locales from 'date-fns/locale';
import { readUserLocale } from './auth';


export const nullDate = parseJSON("0001-01-01T00:00:00Z");

export const isEmptyDate = (date: Date | string | undefined) => {
  if (date === undefined || date === "") {
    return true;
  }
  if (typeof date === "object") {
    return isEqual(date, nullDate);
  }
  return false
}

const dateLocaleFromLng = () => {
  const lngTkns = readUserLocale().split("-");
  if (lngTkns.length > 1) {
    lngTkns[1] = lngTkns[1].toUpperCase()
  }
  //@ts-ignore
  return _get(Locales, lngTkns.join(""), Locales['enUS']);
}

export const parseFlightDateZ = (date: string) => {
  return _parseISO(date.substring(0, date.length - 1))
}

export const parseISO = (date: string) => {
  return _parseISO(date);
}

export const parseTripDate = (tripDate?: string) => {
  if (_isEmpty(tripDate)) {
    return undefined;
  }
  const td = tripDate!
  return isEmptyDate(parseISO(td)) ? undefined : parseISO(td);
}

export const fmt = (date: Date, fmt: string) => {
  return format(date, fmt, { locale: dateLocaleFromLng() });
}

export const fmtMins = (mins: number) => {
  const duration = intervalToDuration({
    start: 0,
    end: mins * 60 * 1000
  });
  return formatDuration(duration);
}

export const areYMDEqual = (date1: Date, date2: Date) => {
  return isEqual(
    new Date(date1.getFullYear(), date1.getMonth(), date1.getDate()),
    new Date(date2.getFullYear(), date2.getMonth(), date2.getDate()),
  )
}
