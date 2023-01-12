import React, {
  ChangeEvent,
  FC,
  ReactElement,
  useEffect,
  useState,
  useRef,
} from 'react';
import { useDebounce } from 'usehooks-ts';
import { useRouter } from "next/router";
import _get from "lodash/get";
import classNames from 'classnames';
import { parseJSON, parseISO, isEqual } from 'date-fns';
import Heap from 'heap-js';
import { applyPatch } from 'json-joy/es6/json-patch';

import {
  CalendarDaysIcon,
  HeartIcon,
  MagnifyingGlassIcon,
  PencilIcon,
  XMarkIcon,
} from '@heroicons/react/24/outline'

import type { NextPageWithLayout } from '../_app'
import TripsAPI from '../../apis/trips';
import TripsSyncAPI from '../../apis/tripsSync';
import ImagesAPI, { stockImageSrc, images } from '../../apis/images';

import PlaneIcon from '../../components/icons/PlaneIcon';
import BusIcon from '../../components/icons/BusIcon';
import HotelIcon from '../../components/icons/HotelIcon';
import Spinner from '../../components/Spinner';
import TripsLayout from '../../components/layouts/TripsLayout';
import { datesRenderer } from '../../utils/dates';
import { WebsocketEvents } from 'websocket-ts/lib';

// TripPageMenu

interface TripPageMenuProps {
  trip: any
  tripStateOnUpdate: any
}


const TripPageMenu: FC<TripPageMenuProps> = (props: TripPageMenuProps) => {
  console.log(props.trip)
  console.log(props.trip.name)

  // UI State
  const [isSelectImageModalOpen, setIsSelectImageModalOpen] = useState(false);
  const [searchImageQuery, setSelectImageQuery] = useState("");
  const [searchImageList, setSearchImageList] = useState([] as any);
  const [isSearchImageLoading, setIsSearchImageLoading] = useState(false);

  const [tripName, setTripName] = useState(props.trip.name);

  useEffect(() => {
    setTripName(props.trip.name)
  }, [props.trip])


  // API
  const searchImage = () => {
    setIsSearchImageLoading(true);
    ImagesAPI.search(searchImageQuery)
    .then(res => {
      const images = _get(res, "data.images");
      setSearchImageList(images);
      setIsSearchImageLoading(false);
    });
    // setSearchImageList(images);
  }

  // Event Handlers - Cover Image
  const searchImageQueryOnChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setSelectImageQuery(event.target.value);
  }

  const searchImageQueryOnEnter = (event: React.KeyboardEvent<HTMLInputElement>) => {
    if (event.key === "Enter") {
      searchImage();
    }
  }

  const searchImageBtnOnClick = () => {
    searchImage();
  }

  // Event Handlers - Trip Name

  const tripNameOnChange = (event: ChangeEvent<HTMLInputElement>) => {
    setTripName(event.target.value)
  }

  const tripNameOnBlur = () => {
    props.tripStateOnUpdate("replace", "/name", tripName)
  }


  // Renderers
  const renderDatesButton = () => {
    if (!_get(props.trip, "startDate")) {
      return;
    }

    const nullDate = parseJSON("0001-01-01T00:00:00Z");
    const startDate = parseISO(props.trip.startDate);
    const endDate = parseJSON(props.trip.endDate);

    if (isEqual(startDate, nullDate)) {
      return "";
    }

    return (
      <button type="button" className="font-medium text-md text-slate-500">
        <CalendarDaysIcon className='inline h-5 w-5 align-sub' />
        &nbsp;&nbsp;
        <span>{datesRenderer(startDate, endDate)}</span>
      </button>
    );
  }

  const renderTripJumbo = () => {
    return (
      <div className='bg-yellow-200'>
        <div className="relative">
          <img className="block sm:max-h-96 w-full" src={stockImageSrc} />
          <button
            type='button'
            className='absolute top-4 right-4 h-10 w-10 bg-gray-800/50 p-2 text-center rounded-full'
            onClick={() => { setIsSelectImageModalOpen(true) }}
          >
            <PencilIcon className='h-6 w-6 text-white' />
          </button>
        </div>
        <div className='h-16 relative -top-24'>
          <div className="bg-white rounded-lg shadow-xl p-5 mx-4 mb-4">
            <input
              type="text"
              value={tripName}
              onChange={tripNameOnChange}
              onBlur={tripNameOnBlur}
              className="mb-12 text-2xl sm:text-4xl font-bold text-slate-700 w-full rounded-lg p-1 border-0 hover:bg-slate-300 hover:border-0 hover:bg-slate-100 focus:ring-0"
            />
            <div className='flex justify-between'>
              {renderDatesButton()}
            </div>
          </div>
        </div>
      </div>
    );
  }

  const renderSelectCoverImageModal = () => {
    if (!isSelectImageModalOpen) {
      return;
    }

    const renderImageThumbnails = () => {
      if (isSearchImageLoading) {
        return <Spinner />
      }

      const imagesThumbnails = searchImageList.map((image: any) => {
        return (
          <figure className="relative max-w-sm transition-all rounded-lg duration-300 mb-2">
            <a href="#">
              <img key={image.id}
                srcSet={ImagesAPI.makeSrcSet(image)}
                src={ImagesAPI.makeSrc(image)}
                className="block rounded-lg max-w-full"
              />
            </a>
            <figcaption className="absolute px-1 text-sm text-white rounded-b-lg bg-slate-800/50 w-full bottom-0">
              <a target="_blank"
                href={ImagesAPI.makeUserReferURL(_get(image, "user.username"))}
              >
                @{_get(image, "user.username")}, Unsplash
              </a>
            </figcaption>
          </figure>
        );
      })
      return (
        <div className='columns-2 md:columns-3'>
          {imagesThumbnails}
        </div>
      );
    }

    return (
      <div className="relative z-10" aria-labelledby="modal-title" role="dialog" aria-modal="true">
        <div className="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity"></div>
        <div className="fixed inset-0 z-10 overflow-y-auto">
          <div className="flex min-h-full items-center justify-center p-4 text-center sm:items-center sm:p-0">
            <div className="relative transform rounded-lg bg-white text-left shadow-xl transition-all w-11/12 sm:my-8 sm:w-full sm:max-w-2xl">
              <div className="bg-white px-4 pt-5 pb-4 sm:p-8 sm:pb-4 rounded-lg">
                <div className='flex justify-between mb-6'>
                  <h2 className="text-lg sm:text-2xl font-bold leading-6 text-slate-900">
                    Change cover image
                  </h2>
                  <button type="button" onClick={() => { setIsSelectImageModalOpen(false) }}>
                    <XMarkIcon className='h-6 w-6 text-slate-700' />
                  </button>
                </div>
                <h2 className="text-sm font-medium text-indigo-500 sm:text-xl text-slate-700 mb-2 ml-1">
                  Search the web
                </h2>
                <div className="flex mb-4 justify-between">
                  <input
                    type="text"
                    className={classNames(
                      "bg-gray-50",
                      "block",
                      "border-gray-300",
                      "border",
                      "focus:border-blue-500",
                      "focus:ring-blue-500",
                      "min-w-0",
                      "p-2.5",
                      "rounded-lg",
                      "text-gray-900",
                      "text-sm",
                      "w-5/6",
                      "mr-2"
                    )}
                    value={searchImageQuery}
                    onChange={searchImageQueryOnChange}
                    onKeyDown={searchImageQueryOnEnter}
                    placeholder="destination, theme ..."
                  />
                  <button
                    type='button'
                    className='flex-1 inline-flex text-white bg-indigo-500 hover:bg-indigo-800 rounded-2xl p-2.5 text-center items-center justify-around'
                    onClick={searchImageBtnOnClick}
                  >
                    <MagnifyingGlassIcon className='h-5 w-5 stroke-2 stroke-white'/>
                  </button>
                </div>
                {renderImageThumbnails()}
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }

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
    <div className='sm:max-w-xl md:max-w-2xl'>
      {renderTripJumbo()}
      {renderTripStats()}
      {renderSelectCoverImageModal()}
    </div>
  );
}


// TripPage

const customPriorityComparator = (a: any, b: any) => a.counter - b.counter;


const TripPage: NextPageWithLayout = () => {
  const router = useRouter();
  const { id } = router.query;

  // Trip State
  const [trip, setTrip] = useState(null as any);
  const [isTripLoaded, setIsTripLoaded] = useState(false);

  // Sync Session State
  const wsInstance = useRef(null as any);
  const pq = new Heap(customPriorityComparator);
  const nextTobCounter = useRef(1);


  useEffect(() => {
    if (id) {
      TripsAPI.readTrip(id as string).then((data) => {
        const t = _get(data, "tripPlan", {});
        setTrip(t);
        setIsTripLoaded(true)
      });
    }
    if (typeof window !== "undefined" && wsInstance.current === null && id) {
      const ws = TripsSyncAPI.startTripSyncSession();
      wsInstance.current = ws;

      ws.addEventListener(WebsocketEvents.open, (i, e) => {
        const joinMsg = TripsSyncAPI.makeSyncMsgJoinSession(
          id as string, "memberID", "memberEmail");
        ws.send(JSON.stringify(joinMsg));
      });
      return () => { ws.close(); }
    }
  }, [id])

  useEffect(() => {
    if (isTripLoaded ) {
      wsInstance.current.addEventListener(WebsocketEvents.message, (_: any, e: any) => {
        const msg = JSON.parse(e.data);
        console.log(msg)

        if (msg.opType === "SyncOpJoinSessionBroadcast") {
          console.log("resetting")
          nextTobCounter.current = 1
          console.log(nextTobCounter.current)

          return;
        }

        if (msg.opType !== "SyncOpUpdateTrip") {
          console.log("updating")
          nextTobCounter.current += 1
          console.log(nextTobCounter.current)
          return;
        }



        pq.push(msg);

        let counter = nextTobCounter.current
        while (true) {
          console.log(pq.length)
          if (pq.length === 0) {
            nextTobCounter.current = counter
            break;
          }
          if (pq.peek().counter !== counter) {
            console.log(pq.peek())
            console.log("break", "msg.c", pq.peek().counter, "exp.c", counter)
            nextTobCounter.current = counter
            break;
          }
          const msg = pq.pop();
          const patch = [msg.syncDataUpdateTrip];
          console.log(trip, patch)
          const newTrip = applyPatch(trip, patch, {mutate: false} as any);
          console.log(newTrip);
          counter += 1
          console.log("set new tri[", newTrip.doc)
          setTrip(newTrip.doc);
        }
      })
    }
  }, [isTripLoaded])


  // Event Handlers
  const tripStateOnUpdate = (op: string, path: string, value: string) => {
    const updateMsg = TripsSyncAPI.makeSyncMsgUpdateTrip(id as string, op, path, value,);
    wsInstance.current.send(JSON.stringify(updateMsg));
  }


  // Renderers
  const renderTripMenu = () => {
    return (
      <aside className='min-h-full min-w-full'>
        <TripPageMenu trip={trip} tripStateOnUpdate={tripStateOnUpdate} />
      </aside>
    );
  }

  if (!trip) {
    return (<Spinner />);
  }

  return (
    <div className="flex">
      {renderTripMenu()}
    </div>
  );
}

export default TripPage;

TripPage.getLayout = function getLayout(page: ReactElement) {
  return (
    <TripsLayout>{page}</TripsLayout>
  )
}

