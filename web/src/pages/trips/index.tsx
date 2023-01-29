import React, {
  FC,
  useEffect,
  useState,
  useRef,
  useCallback,
} from 'react';
import { Link, useParams } from "react-router-dom";
import _get from "lodash/get";
import _isEmpty from "lodash/isEmpty";
import { applyPatch,  } from 'json-joy/es6/json-patch';
import { WebsocketEvents } from 'websocket-ts/lib';
import { GlobeAmericasIcon } from '@heroicons/react/24/outline'

import TripsAPI from '../../apis/trips';
import TripsSyncAPI, { JSONPatchOp } from '../../apis/tripsSync';

import { NewSyncMessageHeap } from '../../utils/heap';
import Spinner from '../../components/Spinner';
import TripMenuJumbo from '../../components/trip/TripMenuJumbo';
import TripLogisticsSection from '../../components/trip/TripLogisticsSection';

import { TripMenuCss } from '../../styles/global';
import TripMap from '../../components/maps/TripMap';


// TripPageMenu

interface TripPageMenuProps {
  trip: any
  tripStateOnUpdate: any
}

const TripPageMenu: FC<TripPageMenuProps> = (props: TripPageMenuProps) => {
  return (
    <div className={TripMenuCss.TripMenu}>
      <nav className={TripMenuCss.TripMenuNav}>
        <Link to="/" className='block align-middle'>
          <GlobeAmericasIcon className='inline h-10 w-10'/>
          <span className='inline-block text-2xl align-middle'>
            tiinyplanet
          </span>
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

const TripPage: FC = () => {
  const routerParams = useParams();
  const { id } = routerParams;

  // Trip State
  const tripRef = useRef(null as any);
  const [tripID, setTripID] = useState("");
  const [trip, setTrip] = useState(tripRef.current);

  const [mapNode, setMapNode] = useState(null) as any;
  const [width, setWidth] = useState(0);
  const measuredRef = useCallback((node: any) => {
    if (node) {
      setMapNode(node)
      // console.log(node.getBoundingClientRect().width)
      setWidth(node.getBoundingClientRect().width);
    }
  }, []);

  // Sync Session State
  const wsInstance = useRef(null as any);
  const pq = NewSyncMessageHeap();
  const nextTobCounter = useRef(0);

  const shouldSetWs = (): boolean => {
    return (wsInstance.current === null && id != null)
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
    }
  }, [id])

  useEffect(() => {
    return () => {
      const ws = wsInstance.current;
      if (ws.underlyingWebsocket?.readyState) {
        ws.close();
      }
    }
  }, [])

  useEffect(() => {
    window.addEventListener("resize", () => {
      if(mapNode) {
        setWidth(mapNode.getBoundingClientRect().width);
      }
    })
  }, [mapNode])

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
      <aside className={TripMenuCss.TripMenuCtn}>
        <TripPageMenu
          trip={tripRef.current}
          tripStateOnUpdate={tripStateOnUpdate}
        />
      </aside>
      <div className='flex-1' ref={measuredRef}>
        <div className="absolute bg-green-200 w-96 h-screen">
          <TripMap trip={tripRef.current} width={width} />
        </div>
      </div>
    </div>
  );
}

export default TripPage;
