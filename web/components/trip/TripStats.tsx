import React, { FC, useState } from 'react';
import _get from "lodash/get";
import _sortBy from "lodash/sortBy";
import _isEmpty from "lodash/isEmpty";

import BusIcon from '../../components/icons/BusIcon';
import HotelIcon from '../../components/icons/HotelIcon';
import PlaneIcon from '../../components/icons/PlaneIcon';
import TripFlightsModal from './TripFlightsModal';
import {
  HeartIcon,
} from '@heroicons/react/24/outline'



interface TripFlightsSectionProps {
  trip: any
}

const TripFlightsSection: FC<TripFlightsSectionProps> = (props: TripFlightsSectionProps) => {

  const [isTripFlightsModalOpen, setIsTripFlightsModalOpen] = useState(false);
  const [selectedFlight, setSelectedFlight] = useState(null);


  return (
    <div className='p-5'>
      <h3 className='text-2xl sm:text-5xl font-bold text-slate-700'>
        Flights
      </h3>
      <button
        className='text-slate-500 text-sm mt-1 font-bold'
        onClick={() => {setIsTripFlightsModalOpen(true)}}
      >
        +&nbsp;&nbsp;Search for a flight
      </button>
      <TripFlightsModal
        isOpen={isTripFlightsModalOpen}
        onClose={() => { setIsTripFlightsModalOpen(false)}}
      />
    </div>
  );
}


// Trip Stats

interface TripStatsProps {
  trip: any
}

const TripStats: FC<TripStatsProps> = (props: TripStatsProps) => {

  // Renderers

  const renderLogisticsTabs = () => {
    const items = [
      { title: "Flights", icon: PlaneIcon },
      { title: "Transits", icon: BusIcon },
      { title: "Lodging", icon: HotelIcon },
      { title: "Insurance", icon: HeartIcon },
    ].map((item, idx) => {
      return (
        <button
          key={idx} type="button"
          className='mx-4 my-2 flex flex-col items-center'
        >
          <item.icon className='h-6 w-6 mb-1' />
          <span className='text-slate-400 text-sm'>{item.title}</span>
        </button>
      );
    })

    return (
      <div className='bg-indigo-100 py-8 pb-4 mb-4'>
        <div className="bg-white rounded-lg p-5 mx-4 mb-4">
          <h5 className="mb-4 text-md sm:text-2xl font-bold text-slate-700">
            Transportation and Lodging
          </h5>
          <div className="flex flex-row justify-around mx-2">
            {items}
          </div>
        </div>
      </div>
    );
  }

  return (
    <div>
      {renderLogisticsTabs()}
      <TripFlightsSection trip={props.trip} />
    </div>
  );



}

export default TripStats;
