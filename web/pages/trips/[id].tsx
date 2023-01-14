import React, {
  ChangeEvent,
  FC,
  ReactElement,
  useEffect,
  useState,
  useRef,
} from 'react';
import _get from "lodash/get";
import { applyPatch,  } from 'json-joy/es6/json-patch';
import { useDebounce } from 'usehooks-ts';
import { useRouter } from "next/router";
import { WebsocketEvents } from 'websocket-ts/lib';
import { HeartIcon } from '@heroicons/react/24/outline'

import type { NextPageWithLayout } from '../_app'
import TripsAPI from '../../apis/trips';
import TripsSyncAPI, { JSONPatchOp } from '../../apis/tripsSync';

import { NewSyncMessageHeap } from '../../utils/heap';
import BusIcon from '../../components/icons/BusIcon';
import HotelIcon from '../../components/icons/HotelIcon';
import PlaneIcon from '../../components/icons/PlaneIcon';
import Spinner from '../../components/Spinner';
import TripsLayout from '../../components/layouts/TripsLayout';
import TripMenuJumbo from '../../components/trip/TripMenuJumbo';


// TripPageMenu

interface TripPageMenuProps {
  trip: any
  tripStateOnUpdate: any
}


const TripPageMenu: FC<TripPageMenuProps> = (props: TripPageMenuProps) => {

  // Renderers

  const renderLogistics = () => {
    const items = [
      { title: "Flights", icon: PlaneIcon },
      { title: "Transits", icon: BusIcon },
      { title: "Lodging", icon: HotelIcon },
      { title: "Insurance", icon: HeartIcon },
    ].map((item, idx) => {
      return (
        <span key={idx} className='mx-4 my-2 flex flex-col items-center '>
          <item.icon className='h-6 w-6 mb-1' />
          <span className='text-slate-400 text-sm'>{item.title}</span>
        </span>
      );
    })

    return (
      <div className="bg-white rounded-lg p-5 mx-4 mb-4">
        <h5 className="mb-4 text-md sm:text-2xl font-bold text-slate-700">
          Transportation and Lodging
        </h5>
        <div className="flex flex-row justify-around mx-2">
          {items}
        </div>
      </div>
    );
  }

  const renderTripStats = () => {
    return (
      <div className='bg-yellow-200 py-8 pb-4 mb-4'>
        {renderLogistics()}
      </div>
    )
  }

  return (
    <div className='sm:max-w-lg md:max-w-xl'>
      <TripMenuJumbo
        trip={props.trip}
        tripStateOnUpdate={props.tripStateOnUpdate}
      />
      {renderTripStats()}
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
  const [isTripLoaded, setIsTripLoaded] = useState(false);

  // Sync Session State
  const wsInstance = useRef(null as any);
  const pq = NewSyncMessageHeap();
  const nextTobCounter = useRef(1);

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
        setIsTripLoaded(true)
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
      return () => { ws.close(); }
    }
  }, [id])

  useEffect(() => {
    if (isTripLoaded) {
      wsInstance.current.addEventListener(WebsocketEvents.message, (_: any, e: any) => {
        const msg = JSON.parse(e.data);
        switch (msg.opType) {
          case "SyncOpJoinSessionBroadcast":
            nextTobCounter.current = 1
            return;
          case "SyncOpUpdateTrip":
            break
          default:
            nextTobCounter.current += 1
            return;
        }

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

          console.log(minMsg)
          // Process this message
          minMsg = pq.pop()!;
          const patch = minMsg.syncDataUpdateTrip!.ops as any;
          const patchOpts = {mutate: false} as any;
          const newTrip = applyPatch(tripRef.current, patch, patchOpts);

          nxtCtr += 1
          tripRef.current = newTrip.doc
          setTrip(tripRef.current);
          console.log(tripRef.current)

        }
      })
    }
  }, [isTripLoaded])


  // Event Handlers
  const tripStateOnUpdate = (ops: Array<JSONPatchOp>) => {
    const updateMsg = TripsSyncAPI.makeSyncMsgUpdateTrip(tripID, ops);
    wsInstance.current.send(JSON.stringify(updateMsg));
  }

  // Renderers
  if (!isTripLoaded) {
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

