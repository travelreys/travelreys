import Heap from 'heap-js';
import { TripSync } from './tripsSync';

const customPriorityComparator = (a: TripSync.Message, b: TripSync.Message) => {
  return a.counter! - b.counter!;
}

export const NewSyncMessageHeap = (): Heap<TripSync.Message> => {
  return new Heap(customPriorityComparator);
}
