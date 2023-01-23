import {
  format,
  formatDuration,
  intervalToDuration,
  isEqual,
  parseISO,
  parseJSON,
} from 'date-fns';

export const datesRenderer = (startDate: Date, endDate: Date) => {
  const nullDate = parseJSON("0001-01-01T00:00:00Z");
  if (isEqual(startDate, nullDate)) {
    return "";
  }
  if (isEqual(endDate, nullDate)) {
    return format(startDate, "MMM d, yy ");
  }
  return `${format(startDate, "MMM d, yy ")} - ${format(endDate, "MMM d, yy ")}`;
}


export const parseTimeFromZ = (date: string) => {
  const d = date.substring(0, date.length - 1);
  return parseISO(d);
}

export const printTime = (date: Date) => {
  return format(date, "hh:mm aa");
}


export const prettyPrintMins = (mins: number) => {
  const duration = intervalToDuration({
    start: 0,
    end: mins * 60 * 1000
  });
  return formatDuration(duration);
}


