import React, { FC, useState } from 'react';
import _get from "lodash/get";
import _sortBy from "lodash/sortBy";
import _isEmpty from "lodash/isEmpty";
import { v4 as uuidv4 } from 'uuid';

import BusIcon from '../icons/BusIcon';
import HotelIcon from '../icons/HotelIcon';
import PlaneIcon from '../icons/PlaneIcon';
import TripFlightsModal from './TripFlightsModal';
import { HeartIcon } from '@heroicons/react/24/outline'

import TripsSyncAPI from '../../apis/tripsSync';
import TransitFlightCard from './TransitFlightCard';

import { TripLogisticsCss } from '../../styles/global';


// TripFlightsSection

interface TripFlightsSectionProps {
  trip: any
  onFlightSelect: any
  onFlightDelete: any
}

const TripFlightsSection: FC<TripFlightsSectionProps> = (props: TripFlightsSectionProps) => {

  const [isTripFlightsModalOpen, setIsTripFlightsModalOpen] = useState(false);

  // Renderers
  const renderItineraries = () => {
    return (
      <div>
        {Object.values(props.trip.flights).map((flight: any, idx: number) =>
          <TransitFlightCard
            key={idx}
            flight={flight}
            onDelete={props.onFlightDelete}
          />
        )}
      </div>
    );
  }

  return (
    <div className='p-5'>
      <div className={TripLogisticsCss.FlightsTitleCtn}>
        <h3 className='text-2xl sm:text-5xl font-bold text-slate-700'>
          Flights
        </h3>
        <button
          className='text-slate-500 text-sm mt-1 font-bold'
          onClick={() => {setIsTripFlightsModalOpen(true)}}
        >
          +&nbsp;&nbsp;Search for a flight
        </button>
      </div>
      {renderItineraries()}
      <TripFlightsModal
        isOpen={isTripFlightsModalOpen}
        onFlightSelect={props.onFlightSelect}
        onClose={() => { setIsTripFlightsModalOpen(false)}}
      />
    </div>
  );
}


// Trip Logistics

interface TripLogisticsSectionProps {
  trip: any
  tripStateOnUpdate: any
}

const TripLogisticsSection: FC<TripLogisticsSectionProps> = (props: TripLogisticsSectionProps) => {

  // Event Handlers

  const flightOnSelect = (itin: any) => {
    let transit = { id: uuidv4(), type: "flight" }
    transit = Object.assign(transit, itin)

    const ops = [
      TripsSyncAPI.makeJSONPatchOp(
        "add", `/flights/${transit.id}`, transit)
    ];
    props.tripStateOnUpdate(ops);
  }

  const flightOnDelete = (transit: any) => {
    const ops = [
      TripsSyncAPI.makeJSONPatchOp(
        "remove", `/flights/${transit.id}`, transit)
    ];
    props.tripStateOnUpdate(ops);
  }

  // Renderers

  const renderTabs = () => {
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
      {renderTabs()}
      <TripFlightsSection
        trip={props.trip}
        onFlightSelect={flightOnSelect}
        onFlightDelete={flightOnDelete}
      />
    </div>
  );

}

export default TripLogisticsSection;
