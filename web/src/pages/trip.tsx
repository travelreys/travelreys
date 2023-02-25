import React, {
  FC,
  useEffect,
  useState,
  useRef,
  useCallback,
} from 'react';
import { useParams } from "react-router-dom";
import _isEmpty from "lodash/isEmpty";
import { applyPatch, } from 'json-joy/es6/json-patch';
import { WebsocketEvents } from 'websocket-ts/lib';
import {
  BanknotesIcon,
  CalendarDaysIcon,
  HomeIcon
} from '@heroicons/react/24/outline';
import { Cog6ToothIcon } from '@heroicons/react/24/solid';

import { NavbarLogo } from '../components/common/Navbar';
import BudgetSection from '../features/trip/BudgetSection';
import ContentSection from '../features/trip/ContentSection';
import ItinerarySection from '../features/trip/ItinerarySection';
import NotesSection from '../features/trip/Notes';
import SettingsSection from '../features/trip/SettingsSection';
import Spinner from '../components/common/Spinner';
import LogisticsSection from '../features/trip/LogisticsSection';
import TripMap from '../features/maps/TripMap';
import TripMenuJumbo from '../features/trip/MenuJumbo';

import TripsAPI, { ReadMembersResponse, ReadResponse } from '../apis/trips';
import TripsSyncAPI from '../apis/tripSync';
import { NewMessageHeap } from '../lib/heap';
import {
  makeMsgUpdateTrip,
  makeMsgJoinSession,
  OpUpdateTrip,
  OpMemberUpdate,
  UpdateTitleAddNewMember,
  Message,
} from '../lib/tripSync';
import { Op } from '../lib/jsonpatch';
import { readAuthUser, User } from '../lib/auth';
import { CommonCss, TripMenuCss } from '../assets/styles/global';
import { MapsProvider } from '../context/maps-context';
import { Member } from '../lib/trips';

// TripPlanningMenu

interface TripPlanningMenuProps {
  trip: any
  tripMembers: {[key: string]: User}
  onlineMembers: Array<Member>
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
        <LogisticsSection
          trip={props.trip}
          tripStateOnUpdate={props.tripStateOnUpdate}
        />
        <hr className={CommonCss.HrShort} />
        <ContentSection
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
          tripMembers={props.tripMembers}
          onlineMembers={props.onlineMembers}
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
            tripMembers={props.tripMembers}
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
  const [trip, setTrip] = useState(tripRef.current);
  const [tripMembers, setTripMembers] = useState({});
  const [onlineMembers, setOnlineMembers] = useState<Array<Member>>([]);

  // Sync Session State
  const wsRef = useRef(null as any);
  const pq = NewMessageHeap();
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

  useEffect(() => {
    window.addEventListener("resize", () => {
      if (mapNode) {
        setMapDivWidth(mapNode.getBoundingClientRect().width);
      }
    })
  }, [mapNode])


  /////////////////////////////
  // Use Effects - Websocket //
  /////////////////////////////

  const shouldSetWs = (id: string|undefined): boolean => {
    return (wsRef.current === null && !_isEmpty(id))
  }

  const handleOpMemberUpdate = (members: Array<Member>) => {
    setOnlineMembers(members);
  }

  const handleUpdateMsgTitle = (title?: string) => {
    if (_isEmpty(title)) {
      return;
    }
    if (title === UpdateTitleAddNewMember) {
      TripsAPI.readMembers(id as string)
        .then((data: ReadMembersResponse) => {
          setTripMembers(data.members);
        });
    }
  }

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
      TripsAPI.read(id as string)
        .then((data: ReadResponse) => {
          tripRef.current = data.trip;
          setTrip(tripRef.current);
          setTripMembers(data.members)
        });
    }
    // Set up WS Connection
    if (shouldSetWs(id)) {
      const ws = TripsSyncAPI.startSession();
      wsRef.current = ws;

      ws.addEventListener(WebsocketEvents.open, () => {
        const user = readAuthUser();
        const joinMsg = makeMsgJoinSession(id as string, user?.id || "");
        ws.send(JSON.stringify(joinMsg));
      });

      ws.addEventListener(WebsocketEvents.message, (_: any, e: any) => {
        const msg = JSON.parse(e.data) as Message;
        console.log(msg)
        switch (msg.op) {
          case OpMemberUpdate:
            handleOpMemberUpdate(msg.data.memberUpdate?.members)
            nextTobCounter.current = 1
            return;
          case OpUpdateTrip:
            break
          default:
            nextTobCounter.current += 1
            return;
        }

        // Add OpUpdateTrip message to min-heap
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
          const newTrip = applyPatch(
            tripRef.current,
            minMsg.data.updateTrip!.ops as any,
            { mutate: false },
          );
          nxtCtr += 1
          tripRef.current = newTrip.doc
          setTrip(tripRef.current);
          handleUpdateMsgTitle(minMsg.data.updateTrip?.title)
        }
      })
    }
  }, [id])

  ////////////////////
  // Event Handlers //
  ////////////////////

  const tripStateOnUpdate = (ops: Array<Op>, title: string) => {
    const msg = makeMsgUpdateTrip(id!, title, ops);
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
            tripMembers={tripMembers}
            onlineMembers={onlineMembers}
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
