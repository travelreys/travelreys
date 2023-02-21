import React, {
  FC,
  useEffect,
  useState,
  useRef,
  useCallback,
} from 'react';
import { useParams } from "react-router-dom";
import _get from "lodash/get";
import _isEmpty from "lodash/isEmpty";
import { applyPatch, } from 'json-joy/es6/json-patch';
import { WebsocketEvents } from 'websocket-ts/lib';
import {
  BanknotesIcon,
  CalendarDaysIcon,
  HomeIcon
} from '@heroicons/react/24/outline';
import { Cog6ToothIcon } from '@heroicons/react/24/solid';

import BudgetSection from '../../features/trip/BudgetSection';
import ItinerarySection from '../../features/trip/ItinerarySection';
import NotesSection from '../../features/trip/Notes';
import Spinner from '../../components/common/Spinner';
import { NavbarLogo } from '../../components/common/Navbar';
import TripContentSection from '../../features/trip/ContentSection';
import TripLogisticsSection from '../../features/trip/LogisticsSection';
import TripMap from '../../features/maps/TripMap';
import TripMenuJumbo from '../../features/trip/MenuJumbo';
import SettingsSection from '../../features/trip/SettingsSection';

import TripsAPI, { ReadResponse } from '../../apis/trips';
import TripsSyncAPI from '../../apis/tripsSync';
import { NewSyncMessageHeap } from '../../lib/heap';
import {
  JSONPatchOp,
  makeMsgUpdateTrip,
  makeMsgJoinSession,
  OpJoinSessionBroadcast,
  OpUpdateTrip
} from '../../lib/tripsSync';
import { Auth, readAuthUser } from '../../lib/auth';
import { CommonCss, TripMenuCss } from '../../assets/styles/global';
import { MapsProvider } from '../../context/maps-context';

// TripPlanningMenu

interface TripPlanningMenuProps {
  trip: any
  tripUsers: {[key: string]: Auth.User}
  tripStateOnUpdate: any
}

const TripPlanningMenu: FC<TripPlanningMenuProps> = (props: TripPlanningMenuProps) => {

  const [tab, setTab] = useState("home");

  // Renderers
  const renderTabs = () => {
    const tabs = [
      { title: "Home", icon: HomeIcon },
      { title: "Itinerary", icon: CalendarDaysIcon },
      { title: "Budget", icon: BanknotesIcon },
      // { title: "Attachments", icon: FolderArrowDownIcon },
      { title: "Settings", icon: Cog6ToothIcon },
    ];

    return (
      <div className={TripMenuCss.TabsCtn}>
        <div className={TripMenuCss.TabsWrapper}>
          <div className={TripMenuCss.TabItemCtn}>
            {tabs.map((tab: any, idx: number) => (
              <button
                key={idx} type="button"
                className={TripMenuCss.TabItemBtn}
                onClick={() => { setTab(tab.title.toLowerCase()) }}
              >
                <tab.icon className='h-6 w-6 mb-1' />
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
        <NotesSection
          trip={props.trip}
          tripStateOnUpdate={props.tripStateOnUpdate}
        />
        <TripLogisticsSection
          trip={props.trip}
          tripStateOnUpdate={props.tripStateOnUpdate}
        />
        <hr className={CommonCss.HrShort} />
        <TripContentSection
          trip={props.trip}
          tripStateOnUpdate={props.tripStateOnUpdate}
        />
      </div>
    );
  }

  return (
    <aside className={TripMenuCss.TripMenuCtn}>
      <div className={TripMenuCss.TripMenu}>
        <nav className={CommonCss.Navbar}>
          <NavbarLogo href='/home' />
        </nav>
        <TripMenuJumbo
          trip={props.trip}
          tripStateOnUpdate={props.tripStateOnUpdate}
        />
        {renderTabs()}
        {tab === "home" ? renderHome() : null}
        {tab === "itinerary"
          ?
          <ItinerarySection
            trip={props.trip}
            tripStateOnUpdate={props.tripStateOnUpdate}
          />
          : null}
        {tab === "budget"
          ?
          <BudgetSection
            trip={props.trip}
            tripStateOnUpdate={props.tripStateOnUpdate}
          />
          : null}
        {tab === "settings"
          ?
          <SettingsSection
            trip={props.trip}
            tripUsers={props.tripUsers}
            tripStateOnUpdate={props.tripStateOnUpdate}
          />
          : null}
      </div>
    </aside>
  );
}

// TripPage

const TripPage: FC = () => {
  const { id } = useParams();

  // Trip State
  const tripRef = useRef(null as any);
  const [tripID, setTripID] = useState("");
  const [trip, setTrip] = useState(tripRef.current);
  const [tripUsers, setTripUsers] = useState({});


  // Sync Session State
  const wsRef = useRef(null as any);
  const pq = NewSyncMessageHeap();
  const nextTobCounter = useRef(0);


  // Map UI State
  const [mapNode, setMapNode] = useState(null) as any;
  const [mapDivWidth, setMapDivWidth] = useState(0);
  const measuredRef = useCallback((node: any) => {
    if (node) {
      setMapNode(node)
      setMapDivWidth(node.getBoundingClientRect().width);
    }
  }, []);

  ///////////////////////
  // Use Effects - Map //
  ///////////////////////

  useEffect(() => {
    window.addEventListener("resize", () => {
      if (!mapNode) {
        return;
      }
      setMapDivWidth(mapNode.getBoundingClientRect().width);
    })
  }, [mapNode])


  /////////////////////////////
  // Use Effects - Websocker //
  /////////////////////////////

  // Close WS when leaving page
  useEffect(() => {
    return () => {
      const ws = wsRef.current;
      if (ws.underlyingWebsocket?.readyState) {
        ws.close();
      }
    }
  }, [])

  const shouldSetWs = (id: string|undefined): boolean => {
    return (wsRef.current === null && !_isEmpty(id))
  }

  useEffect(() => {
    // Fetch Trip
    if (id) {
      TripsAPI.readTrip(id as string)
        .then((data: ReadResponse) => {
          tripRef.current = data.trip;
          setTrip(tripRef.current);
          setTripUsers(data.users)
          setTripID(id as string);
        });
    }
    // Set up WS Connection
    if (shouldSetWs(id)) {
      const ws = TripsSyncAPI.startTripSyncSession();
      wsRef.current = ws;

      ws.addEventListener(WebsocketEvents.open, () => {
        const user = readAuthUser();
        const joinMsg = makeMsgJoinSession(id as string, user?.id || "");
        ws.send(JSON.stringify(joinMsg));
      });

      ws.addEventListener(WebsocketEvents.message, (_: any, e: any) => {
        const msg = JSON.parse(e.data);
        switch (msg.opType) {
          case OpJoinSessionBroadcast:
            nextTobCounter.current = 1
            return;
          case OpUpdateTrip:
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
          const patch = minMsg.data.updateTrip!.ops as any;
          const patchOpts = { mutate: false } as any;
          const newTrip = applyPatch(tripRef.current, patch, patchOpts);

          nxtCtr += 1
          tripRef.current = newTrip.doc
          setTrip(tripRef.current);
        }
      })
    }
  }, [id])

  ////////////////////
  // Event Handlers //
  ////////////////////

  const tripStateOnUpdate = (ops: Array<JSONPatchOp>, title: string) => {
    const msg = makeMsgUpdateTrip(tripID, title, ops);
    wsRef.current.send(JSON.stringify(msg));
  }

  // Renderers

  if (_isEmpty(trip)) {
    return (<Spinner />);
  }

  return (
    <MapsProvider>
      <main className='min-h-screen'>
        <div className="flex">
          <TripPlanningMenu
            trip={tripRef.current}
            tripUsers={tripUsers}
            tripStateOnUpdate={tripStateOnUpdate}
          />
          <div className='flex-1' ref={measuredRef}>
            <TripMap
              trip={tripRef.current}
              width={mapDivWidth}
            />
          </div>
        </div>
      </main>
    </MapsProvider>
  );
}

export default TripPage;
