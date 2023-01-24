import React, {
  FC,
  ReactElement,
  useEffect,
  useState,
  useRef,
} from 'react';
import Link from 'next/link';
import _get from "lodash/get";
import _isEmpty from "lodash/isEmpty";
import { applyPatch,  } from 'json-joy/es6/json-patch';
import { useDebounce } from 'usehooks-ts';
import { useRouter } from "next/router";
import { WebsocketEvents } from 'websocket-ts/lib';
import { GlobeAmericasIcon } from '@heroicons/react/24/outline'

import type { NextPageWithLayout } from '../_app'
import TripsAPI from '../../apis/trips';
import TripsSyncAPI, { JSONPatchOp } from '../../apis/tripsSync';

import { NewSyncMessageHeap } from '../../utils/heap';
import Spinner from '../../components/Spinner';
import TripsLayout from '../../components/layouts/TripsLayout';
import TripMenuJumbo from '../../components/trip/TripMenuJumbo';
import TripLogisticsSection from '../../components/trip/TripLogisticsSection';


// TripPageMenu

interface TripPageMenuProps {
  trip: any
  tripStateOnUpdate: any
}


const TripPageMenu: FC<TripPageMenuProps> = (props: TripPageMenuProps) => {
  return (
    <div className='sm:max-w-lg md:max-w-xl'>
      <nav className="p-3 font-bold text-indigo-500" >
        <Link href="/" className='block align-middle'>
          <GlobeAmericasIcon className='inline h-10 w-10'/>
          <span className='inline-block text-2xl align-middle'>tiinyplanet</span>
        </Link>
      </nav>
      <TripMenuJumbo
        trip={props.trip}
        tripStateOnUpdate={props.tripStateOnUpdate}
      />
      <TripLogisticsSection
        trip={props.trip}
        tripStateOnUpdate={props.tripStateOnUpdate}
      />
    </div>
  );
}

// TripPage

const TripPage: NextPageWithLayout = () => {
  const router = useRouter();
  const { id } = router.query;

  // Trip State
  const [tripID, setTripID] = useState("");
  const tripRef = useRef(null as any);
  const [trip, setTrip] = useState(tripRef.current);
  // Sync Session State
  const wsInstance = useRef(null as any);
  const pq = NewSyncMessageHeap();
  const nextTobCounter = useRef(0);

  const shouldSetWs = (): boolean => {
    return (typeof window !== "undefined"
       && wsInstance.current === null
       && id != null)
  }

  useEffect(() => {
    if (id) {
      TripsAPI.readTrip(id as string).then((data) => {
        tripRef.current = _get(data, "tripPlan", {});
        setTrip(tripRef.current);
        setTripID(id as string);
      });
    }
    if (shouldSetWs()) {
      const ws = TripsSyncAPI.startTripSyncSession();
      wsInstance.current = ws;

      ws.addEventListener(WebsocketEvents.open, () => {
        const joinMsg = TripsSyncAPI.makeSyncMsgJoinSession(
          id as string,
          "memberID",
          "memberEmail");
        ws.send(JSON.stringify(joinMsg));
      });

      ws.addEventListener(WebsocketEvents.message, (_: any, e: any) => {
        const msg = JSON.parse(e.data);
        switch (msg.opType) {
          case "SyncOpJoinSessionBroadcast":
            console.log("bcast")
            nextTobCounter.current = 1
            return;
          case "SyncOpUpdateTrip":
            break
          default:
            nextTobCounter.current += 1
            return;
          }

        console.log(msg, nextTobCounter.current)
        // Add message to min-heap
        pq.push(msg);

        let nxtCtr = nextTobCounter.current
        while (true) {
          if (pq.length === 0) {
            nextTobCounter.current = nxtCtr
            break;
          }
          let minMsg = pq.peek()!;

          // Messages in the heap are further up in the TOB,
          // waiting for the next in-order TOB msg
          if (minMsg.counter !== nxtCtr) {
            nextTobCounter.current = nxtCtr;
            break;
          }

          // Process this message
          minMsg = pq.pop()!;
          const patch = minMsg.syncDataUpdateTrip!.ops as any;
          const patchOpts = {mutate: false} as any;
          const newTrip = applyPatch(tripRef.current, patch, patchOpts);

          nxtCtr += 1
          tripRef.current = newTrip.doc
          setTrip(tripRef.current);
        }
      })

      return () => { ws.close(); }
    }
  }, [id])

  // Event Handlers
  const tripStateOnUpdate = (ops: Array<JSONPatchOp>) => {
    const updateMsg = TripsSyncAPI.makeSyncMsgUpdateTrip(tripID, ops);
    wsInstance.current.send(JSON.stringify(updateMsg));
  }

  // Renderers
  if (_isEmpty(trip)) {
    return (<Spinner />);
  }

  return (
    <div className="flex">
      <aside className='min-h-full min-w-full'>
        <TripPageMenu
          trip={trip}
          tripStateOnUpdate={tripStateOnUpdate}
        />
      </aside>
    </div>
  );
}

export default TripPage;

TripPage.getLayout = function getLayout(page: ReactElement) {
  return (
    <TripsLayout>{page}</TripsLayout>
  )
}

