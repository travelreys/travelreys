import { format, parseJSON, isEqual } from 'date-fns';

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
