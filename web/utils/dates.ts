import _get from "lodash/get";
import {
  format,
  formatDuration,
  intervalToDuration,
  isEqual,
  parseISO,
  parseJSON,
} from 'date-fns';
import { DateRange } from 'react-day-picker';

export const isNullDate = (date: Date) => {
  const nullDate = parseJSON("0001-01-01T00:00:00Z");
  return isEqual(date, nullDate);
}

export const printFromDateFromRange = (range: DateRange | undefined, fmt: string) => {
  const date = _get(range, "from");
  if (date) {
    return format(date, fmt);
  }
  return undefined;
}

export const printToDateFromRange = (range: DateRange | undefined, fmt: string) => {
  const date = _get(range, "to");
  if (date) {
    return format(date, fmt);
  }
  return undefined;
}

export const parseTimeFromZ = (date: string) => {
  const d = date.substring(0, date.length - 1);
  return parseISO(d);
}

export const printTime = (date: Date, fmt: string) => {
  return format(date, fmt);
}

export const prettyPrintMins = (mins: number) => {
  const duration = intervalToDuration({
    start: 0,
    end: mins * 60 * 1000
  });
  return formatDuration(duration);
}


