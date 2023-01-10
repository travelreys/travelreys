import React, { FC, ReactElement } from 'react';
import { useRouter } from "next/router";
import _get from "lodash/get";
import { parseJSON, parseISO, isEqual } from 'date-fns';

import { CalendarDaysIcon } from '@heroicons/react/24/solid'
import PlaneIcon from '../../components/icons/PlaneIcon';
import BusIcon from '../../components/icons/BusIcon';
import HotelIcon from '../../components/icons/HotelIcon';

import TripsAPI from '../../apis/trips';
import type { NextPageWithLayout } from '../_app'
import Spinner from '../../components/Spinner';
import TripsLayout from '../../components/layouts/TripsLayout';
import { datesRenderer } from '../../utils/dates';





// TripPageMenu

interface TripPageMenuProps {
  trip: any
}


const TripPageMenu: FC<TripPageMenuProps> = (props: TripPageMenuProps) => {

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
        <div className="rounded-lg shadow-md dark:bg-gray-800 dark:border-gray-700">
          <img
            className="h-auto max-w-full"
            src="https://images.unsplash.com/photo-1469854523086-cc02fe5d8800?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1421&q=80"
            alt=""
          />
        </div>
        <div className='h-16 relative -top-24'>
          <div className="bg-white rounded-lg shadow-xl p-5 mx-4 mb-4">
            <h5 className="mb-12 text-2xl sm:text-4xl font-bold text-slate-700">
              {props.trip.name}
            </h5>
            <div className='flex justify-between'>
              {renderDatesButton()}
            </div>
          </div>
        </div>
      </div>
    );
  }

  const renderLogistics = () => {
    const items = [
      {title: "Flights", icon: PlaneIcon},
      {title: "Transits", icon: BusIcon},
      {title: "Lodging", icon: HotelIcon},
    ].map((item, idx) => {
      return (
        <span key={idx} className='mx-4 my-2 flex flex-col items-center '>
          <item.icon classNames='h-6 w-6 mb-1'/>
          <span className='text-slate-400 text-sm'>{item.title}</span>
        </span>
      );
    })

    return (
      <div className="bg-white rounded-lg p-5 mx-4 mb-4">
        <h5 className="mb-4 text-md sm:text-4xl font-bold text-slate-700">
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
      <div className='bg-yellow-200 pb-4 mb-4'>
        {renderLogistics()}
      </div>
    )
  }

  return (
    <div>
      {renderTripJumbo()}
      {renderTripStats()}
    </div>
  );
}









// TripPage

const TripPage: NextPageWithLayout = () => {
  const router = useRouter();
  const { id } = router.query;

  let { data, error, isLoading } = TripsAPI.readTrip(id as string);
  const trip = _get(data, "tripPlan", {});

  console.log(data)


  // Renderers
  const renderTripMenu = () => {
    return (
      <aside className='min-h-full min-w-full'>
        <TripPageMenu trip={trip} />
      </aside>
    );
  }

  if (isLoading) {
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

