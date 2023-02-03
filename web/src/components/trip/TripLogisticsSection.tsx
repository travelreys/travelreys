import React, { FC, useState } from 'react';
import _get from "lodash/get";
import _sortBy from "lodash/sortBy";
import _isEmpty from "lodash/isEmpty";
import { v4 as uuidv4 } from 'uuid';

import {
  ChevronDownIcon,
  ChevronUpIcon,
  FolderArrowDownIcon
} from '@heroicons/react/24/outline'

import BusIcon from '../icons/BusIcon';
import HotelIcon from '../icons/HotelIcon';
import PlaneIcon from '../icons/PlaneIcon';
import TripFlightsModal from './TripFlightsModal';

import TripsSyncAPI from '../../apis/tripsSync';
import TransitFlightCard from './TransitFlightCard';

import { TripLogisticsCss } from '../../styles/global';
import TripLodgingsModal from './TripLodgingsModal';
import TripLodgingCard from './TripLodgingCard';


// TripFlightsSection

interface TripFlightsSectionProps {
  trip: any
  onFlightSelect: any
  onFlightDelete: any
}

const TripFlightsSection: FC<TripFlightsSectionProps> = (props: TripFlightsSectionProps) => {

  const [isTripFlightsModalOpen, setIsTripFlightsModalOpen] = useState(false);
  const [isHidden, setIsHidden] = useState(false);

  // Event Handler
  const onFlightSelect = (transit: any) => {
    props.onFlightSelect(transit);
    setIsTripFlightsModalOpen(false);
  }

  // Renderers
  const renderHiddenToggle = () => {
    let icon = <ChevronDownIcon className='h-4 w-4'/>;
    if (isHidden) {
      icon = <ChevronUpIcon className='h-4 w-4' />;
    }
    return (
      <button
        type="button"
        onClick={() => {setIsHidden(!isHidden)}}
      >
      {icon}
      </button>
    );
  }

  const renderItineraries = () => {
    if (isHidden) {
      return null;
    }

    const flights = Object.values(props.trip.flights);
    return flights.map((flight: any, idx: number) =>
      <TransitFlightCard
        key={idx}
        flight={flight}
        onDelete={props.onFlightDelete}
      />
    )
  }

  return (
    <div className='p-5'>
      <div className={TripLogisticsCss.FlightsTitleCtn}>
        <div className='text-2xl sm:text-3xl font-bold text-slate-700'>
          <span>Flights&nbsp;&nbsp;</span>
          {renderHiddenToggle()}
        </div>
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
        onFlightSelect={onFlightSelect}
        onClose={() => { setIsTripFlightsModalOpen(false)}}
      />
    </div>
  );
}


// TripLodgingSection

interface TripLodgingSectionProps {
  trip: any
  onLodgingSelect: any
  onLodgingUpdate: any
  onLodgingDelete: any
}

const TripLodgingSection: FC<TripLodgingSectionProps> = (props: TripLodgingSectionProps) => {

  const [isLodgingModalOpen, setIsLogdingModalOpen] = useState(false);
  const [isHidden, setIsHidden] = useState(false);

  // Event Handlers
  const onLodgingSelect = (lodging: any) => {
    props.onLodgingSelect(lodging);
    setIsLogdingModalOpen(false);
  }

  // Renderers
  const renderHiddenToggle = () => {
    let icon = <ChevronDownIcon className='h-4 w-4'/>;
    if (isHidden) {
      icon = <ChevronUpIcon className='h-4 w-4' />;
    }
    return (
      <button
        type="button"
        onClick={() => {setIsHidden(!isHidden)}}
      >
      {icon}
      </button>
    );
  }

  const renderLodgings = () => {
    if (isHidden) {
      return null;
    }
    const lodgings = Object.values(props.trip.lodgings);
    return lodgings.map((lodge: any) => (
      <TripLodgingCard
        key={lodge.id}
        lodging={lodge}
        onDelete={props.onLodgingDelete}
        onUpdate={props.onLodgingUpdate}
      />
    ));
  }

  return (
    <div className='p-5'>
      <div className={TripLogisticsCss.FlightsTitleCtn}>
        <div className='text-2xl sm:text-3xl font-bold text-slate-700'>
          <span>Hotels and Lodgings&nbsp;&nbsp;</span>
          {renderHiddenToggle()}
        </div>

        <button
          className='text-slate-500 text-sm mt-1 font-bold'
          onClick={() => {setIsLogdingModalOpen(true)}}
        >
          +&nbsp;&nbsp;Add a lodging
        </button>
      </div>
      {renderLodgings()}
      <TripLodgingsModal
        isOpen={isLodgingModalOpen}
        onLodgingSelect={onLodgingSelect}
        onClose={() => { setIsLogdingModalOpen(false) }}
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

  // Event Handlers - Flights

  const flightOnSelect = (flight: any) => {
    let transit = { id: uuidv4(), type: "flight" }
    transit = Object.assign(transit, flight)

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

  // Event Handlers - Lodging

  const lodgingOnSelect = (lodging: any) => {
    lodging = Object.assign(lodging, {id: uuidv4()});
    const ops = [
      TripsSyncAPI.makeJSONPatchOp(
        "add", `/lodgings/${lodging.id}`, lodging)
    ];
    props.tripStateOnUpdate(ops);
  }

  const lodgingOnUpdate = (lodging: any, updates: any) => {
    const ops = [] as any;
    Object.entries(updates).forEach(([key, value]) => {
      const fullpath = `/lodgings/${lodging.id}/${key}`;
      ops.push(TripsSyncAPI.makeJSONPatchOp("replace", fullpath, value));
    });
    props.tripStateOnUpdate(ops);
  }

  const lodgingOnDelete = (lodging: any) => {
    const ops = [
      TripsSyncAPI.makeJSONPatchOp(
        "remove", `/lodgings/${lodging.id}`, lodging)
    ];
    props.tripStateOnUpdate(ops);
  }

  // Renderers

  const renderTabs = () => {
    const items = [
      { title: "Flights", icon: PlaneIcon },
      { title: "Transits", icon: BusIcon },
      { title: "Lodging", icon: HotelIcon },
      { title: "Attachments", icon: FolderArrowDownIcon },
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
            Logistics
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
      <TripLodgingSection
        trip={props.trip}
        onLodgingSelect={lodgingOnSelect}
        onLodgingUpdate={lodgingOnUpdate}
        onLodgingDelete={lodgingOnDelete}
      />
    </div>
  );

}

export default TripLogisticsSection;
