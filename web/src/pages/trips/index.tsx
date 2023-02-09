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
import {
  BanknotesIcon,
  CalendarDaysIcon,
  FolderArrowDownIcon,
  GlobeAmericasIcon,
  HomeIcon
} from '@heroicons/react/24/outline'

import TripsAPI from '../../apis/trips';
import TripsSyncAPI, { JSONPatchOp } from '../../apis/tripsSync';

import Spinner from '../../components/Spinner';
import TripContentSection from '../../components/trip/TripContentSection';
import TripLogisticsSection from '../../components/trip/TripLogisticsSection';
import TripMap from '../../components/maps/TripMap';
import TripMenuJumbo from '../../components/trip/TripMenuJumbo';

import { TripMenuCss } from '../../styles/global';
import { MapsProvider } from '../../context/maps-context';
import { NewSyncMessageHeap } from '../../utils/heap';
import TripNotesSection from '../../components/trip/TripNotes';
import TripItinerarySection from '../../components/trip/TripItinerarySection';


// TripPlanningMenu

interface TripPlanningMenuProps {
  trip: any
  tripStateOnUpdate: any
}

const TripPlanningMenu: FC<TripPlanningMenuProps> = (props: TripPlanningMenuProps) => {

  const [tab, setTab] = useState("home");


  // Renderers
  const renderNavBar = () => {
    return (
      <nav className={TripMenuCss.TripMenuNav}>
        <Link to="/" className='block align-middle'>
          <GlobeAmericasIcon className='inline h-10 w-10'/>
          <span className='inline-block text-2xl align-middle'>
            tiinyplanet
          </span>
        </Link>
      </nav>
    );
  }

  const renderTabs = () => {
    const tabs = [
      { title: "Home", icon: HomeIcon },
      { title: "Itinerary", icon: CalendarDaysIcon },
      { title: "Budget", icon: BanknotesIcon },
      { title: "Attachments", icon: FolderArrowDownIcon },
    ];

    return (
      <div className={TripMenuCss.TabsCtn}>
        <div className={TripMenuCss.TabsWrapper}>
          <div className={TripMenuCss.TabItemCtn}>
            {tabs.map((tab: any, idx: number) => (
              <button
                key={idx} type="button"
                className={TripMenuCss.TabItemBtn}
                onClick={() => { setTab(tab.title.toLowerCase())} }
              >
                <tab.icon className='h-6 w-6 mb-1'/>
                <span className={TripMenuCss.TabItemBtnTxt}>
                  {tab.title}
                </span>
              </button>
            ))}
          </div>
        </div>
      </div>
    );
  }

  const renderHome = () => {
    return (
      <div>
        <TripNotesSection
          trip={props.trip}
          tripStateOnUpdate={props.tripStateOnUpdate}
        />
        <TripLogisticsSection
          trip={props.trip}
          tripStateOnUpdate={props.tripStateOnUpdate}
        />
        <hr className='w-48 h-1 m-5 mx-auto bg-gray-300 border-0 rounded'/>
        <TripContentSection
          trip={props.trip}
          tripStateOnUpdate={props.tripStateOnUpdate}
        />
      </div>
    );
  }

  const renderItinerary = () => {
    return (
      <div>
        <TripItinerarySection
          trip={props.trip}
          tripStateOnUpdate={props.tripStateOnUpdate}
        />
      </div>
    );
  }

  return (
    <aside className={TripMenuCss.TripMenuCtn}>
      <div className={TripMenuCss.TripMenu}>
        {renderNavBar()}
        <TripMenuJumbo
          trip={props.trip}
          tripStateOnUpdate={props.tripStateOnUpdate}
        />
        {renderTabs()}
        { tab === "home" ? renderHome() : null}
        { tab === "itinerary" ? renderItinerary() : null}
      </div>
    </aside>
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

  // Map UI State
  const [mapNode, setMapNode] = useState(null) as any;
  const [mapDivWidth, setMapDivWidth] = useState(0);
  const measuredRef = useCallback((node: any) => {
    if (node) {
      setMapNode(node)
      setMapDivWidth(node.getBoundingClientRect().width);
    }
  }, []);

  // Sync Session State
  const wsRef = useRef(null as any);
  const pq = NewSyncMessageHeap();
  const nextTobCounter = useRef(0);

  const shouldSetWs = (): boolean => {
    return (wsRef.current === null && id != null)
  }

  // Use Effects

  // Close WS when leaving page
  useEffect(() => {
    return () => {
      const ws = wsRef.current;
      if (ws.underlyingWebsocket?.readyState) {
        ws.close();
      }
    }
  }, [])


  useEffect(() => {
    // Fetch Trip
    if (id) {
      TripsAPI.readTrip(id as string).then((data) => {
        tripRef.current = _get(data, "tripPlan", {});
        setTrip(tripRef.current);
        setTripID(id as string);
      });
    }
    // Set up WS Connection
    if (shouldSetWs()) {
      const ws = TripsSyncAPI.startTripSyncSession();
      wsRef.current = ws;

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
    window.addEventListener("resize", () => {
      if(mapNode) {
        setMapDivWidth(mapNode.getBoundingClientRect().width);
      }
    })
  }, [mapNode])

  // Event Handlers
  const tripStateOnUpdate = (ops: Array<JSONPatchOp>) => {
    const updateMsg = TripsSyncAPI.makeSyncMsgUpdateTrip(tripID, ops);
    wsRef.current.send(JSON.stringify(updateMsg));
  }

  // Renderers

  if (_isEmpty(trip)) {
    return (<Spinner />);
  }

  return (
    <MapsProvider>
      <div className="flex">
        <TripPlanningMenu
          trip={tripRef.current}
          tripStateOnUpdate={tripStateOnUpdate}
        />
        <div className='flex-1' ref={measuredRef}>
          <TripMap
            trip={tripRef.current}
            width={mapDivWidth}
          />
        </div>
      </div>
    </MapsProvider>
  );
}

export default TripPage;
