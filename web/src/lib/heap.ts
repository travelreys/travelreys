import Heap from 'heap-js';
import { Message } from './tripSync';

const customPriorityComparator = (a: Message, b: Message) => {
  return a.counter! - b.counter!;
}

export const NewMessageHeap = (): Heap<Message> => {
  return new Heap(customPriorityComparator);
}
