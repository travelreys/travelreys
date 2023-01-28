import Heap from 'heap-js';
import { SyncMessage } from '../apis/tripsSync';

const customPriorityComparator = (a: SyncMessage, b: SyncMessage) => {
  return a.counter - b.counter;
}

export const NewSyncMessageHeap = (): Heap<SyncMessage> => {
  return new Heap(customPriorityComparator);
}
