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
import ActivitySection from '../features/trip/ActivitySection';
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
  OpJoinSession,
  OpLeaveSession,
  MsgUpdateTripTitleAddNewMember,
  Message,
} from '../lib/tripSync';
import { Op } from '../lib/jsonpatch';
import { readAuthUser, User } from '../lib/auth';
import { Member } from '../lib/trips';
import { CommonCss } from '../assets/styles/global';
import { MapsProvider } from '../context/maps-context';

// TripPlanningMenu

interface TripPlanningMenuProps {
  trip: any
  tripMembers: {[key: string]: User}
  onlineMembers: Array<Member>
  tripOnUpdate: any
}

const TripPlanningMenu: FC<TripPlanningMenuProps> = (props: TripPlanningMenuProps) => {
  const [tab, setTab] = useState("home");

  // Renderers
  const css = {
    ctn: "min-h-screen w-full z-50 sm:w-1/2  sm:max-w-lg sm:shadow-xl sm:shadow-slate-900",
    wrapper: "pb-40 w-full",
    tabsCtn: "sticky top-0 z-10 bg-indigo-100 py-8 pb-4 mb-4",
    tabsWrapper: "bg-white rounded-lg p-5 mx-4 mb-4",
    tabItemCtn: "flex flex-row justify-around mx-2",
    tabItemBtn: "mx-4 my-2 flex flex-col items-center",
    tabItemBtnTxt: "text-slate-400 text-sm",
  }

  const renderTabs = () => {
    const tabs = [
      { title: "Home", icon: HomeIcon },
      { title: "Itinerary", icon: CalendarDaysIcon },
      { title: "Budget", icon: BanknotesIcon },
      // { title: "Attachments", icon: FolderArrowDownIcon },
      { title: "Settings", icon: Cog6ToothIcon },
    ];

    return (
      <div className={css.tabsCtn}>
        <div className={css.tabsWrapper}>
          <div className={css.tabItemCtn}>
            {tabs.map((tab: any, idx: number) => (
              <button
                key={idx} type="button"
                className={css.tabItemBtn}
                onClick={() => { setTab(tab.title.toLowerCase()) }}
              >
                <tab.icon className='h-6 w-6 mb-1' />
                <span className={css.tabItemBtnTxt}>
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
          tripOnUpdate={props.tripOnUpdate}
        />
        <LogisticsSection
          trip={props.trip}
          tripOnUpdate={props.tripOnUpdate}
        />
        <hr className={CommonCss.HrShort} />
        <ActivitySection
          trip={props.trip}
          tripOnUpdate={props.tripOnUpdate}
        />
      </div>
    );
  }

  return (
    <aside className={css.ctn}>
      <div className={css.wrapper}>
        <nav className={CommonCss.Navbar}>
          <NavbarLogo href='/home' />
        </nav>
        <TripMenuJumbo
          trip={props.trip}
          tripMembers={props.tripMembers}
          onlineMembers={props.onlineMembers}
          tripOnUpdate={props.tripOnUpdate}
        />
        {renderTabs()}
        {tab === "home" ? renderHome() : null}
        {tab === "itinerary"
          ?
          <ItinerarySection
            trip={props.trip}
            tripOnUpdate={props.tripOnUpdate}
          />
          : null}
        {tab === "budget"
          ?
          <BudgetSection
            trip={props.trip}
            tripOnUpdate={props.tripOnUpdate}
          />
          : null}
        {tab === "settings"
          ?
          <SettingsSection
            trip={props.trip}
            tripMembers={props.tripMembers}
            tripOnUpdate={props.tripOnUpdate}
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
  const tobCounter = useRef(0);

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

  const handleJoinSessionMsg = (msg: Message) => {
    setOnlineMembers(msg.data.joinSession?.members || []);
  }

  const handleLeaveSessionMsg = (msg: Message) => {
    setOnlineMembers(msg.data.leaveSession?.members || []);
  }

  const handleUpdateMsg = (msg: Message) => {
    const title = msg.data.updateTrip?.title;
    if (_isEmpty(title)) {
      return;
    }
    if (title === MsgUpdateTripTitleAddNewMember) {
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

        if (msg.op === OpJoinSession) {
          if (msg.data.joinSession?.id === readAuthUser()?.id) {
            tobCounter.current = msg.counter!
            handleJoinSessionMsg(msg)
            return;
          }
        }
        pq.push(msg);

        let nxtCtr = tobCounter.current + 1;
        while (true) {
          if (pq.length === 0) {
            break;
          }
          let minMsg = pq.peek()!;
          // console.log("min", minMsg.counter, "nxt", nxtCtr)
          // Messages in the heap are further up in the TOB,
          // waiting for the next in-order TOB msg
          if (minMsg.counter !== nxtCtr) {
            break;
          }

          // Process this message
          minMsg = pq.pop()!;

          switch (minMsg.op) {
            case OpJoinSession:
              handleJoinSessionMsg(minMsg);
              break;
            case OpLeaveSession:
              handleLeaveSessionMsg(minMsg);
              break;
            case OpUpdateTrip:
              console.log(minMsg.data.updateTrip!.ops)
              const newTrip = applyPatch(
                tripRef.current,
                minMsg.data.updateTrip!.ops as any,
                { mutate: false },
              );
              tripRef.current = newTrip.doc
              setTrip(tripRef.current);
              handleUpdateMsg(minMsg)
              break;
          }
          tobCounter.current = nxtCtr
          nxtCtr = nxtCtr + 1
        }
      })
    }
  }, [id])

  ////////////////////
  // Event Handlers //
  ////////////////////

  const tripOnUpdate = (ops: Array<Op>, title: string) => {
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
            tripOnUpdate={tripOnUpdate}
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
