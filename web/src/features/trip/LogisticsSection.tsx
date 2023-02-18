import React, { FC, useState } from 'react';
import _get from "lodash/get";
import _sortBy from "lodash/sortBy";
import _isEmpty from "lodash/isEmpty";
import { v4 as uuidv4 } from 'uuid';

import {
  ChevronDownIcon,
  ChevronUpIcon,
} from '@heroicons/react/24/outline'

import FlightsModal from './FlightsModal';
import LodgingsModal from './LodgingsModal';
import LodgingCard from './LodgingCard';

import { Trips } from '../../lib/trips';
import TripsSyncAPI from '../../apis/tripsSync';
import TransitFlightCard from './TransitFlightCard';

import { TripLogisticsCss } from '../../assets/styles/global';



// FlightsSection

interface FlightsSectionProps {
  trip: any
  onFlightSelect: any
  onFlightDelete: any
}

const FlightsSection: FC<FlightsSectionProps> = (props: FlightsSectionProps) => {

  const [isTripFlightsModalOpen, setIsTripFlightsModalOpen] = useState(false);
  const [isHidden, setIsHidden] = useState(false);

  // Event Handler
  const onFlightSelect = (transit: any) => {
    props.onFlightSelect(transit);
    setIsTripFlightsModalOpen(false);
  }

  // Renderers
  const renderHiddenToggle = () => {
    return (
      <button
        type="button"
        className={TripLogisticsCss.FlightsToggleBtn}
        onClick={() => {setIsHidden(!isHidden)}}
      >
      {isHidden ? <ChevronUpIcon className='h-4 w-4' />
        : <ChevronDownIcon className='h-4 w-4'/> }
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
        <div className={TripLogisticsCss.FlightsHeaderCtn}>
          {renderHiddenToggle()}
          <span>Flights</span>
        </div>
        <button
          className={TripLogisticsCss.SearchFlightBtn}
          onClick={() => {setIsTripFlightsModalOpen(true)}}
        >
          +&nbsp;&nbsp;Search for a flight
        </button>
      </div>
      {renderItineraries()}
      <FlightsModal
        trip={props.trip}
        isOpen={isTripFlightsModalOpen}
        onFlightSelect={onFlightSelect}
        onClose={() => { setIsTripFlightsModalOpen(false)}}
      />
    </div>
  );
}


// LodgingSection

interface LodgingSectionProps {
  trip: any
  onLodgingSelect: any
  onLodgingUpdate: any
  onLodgingDelete: any
}

const LodgingSection: FC<LodgingSectionProps> = (props: LodgingSectionProps) => {

  const [isLodgingModalOpen, setIsLogdingModalOpen] = useState(false);
  const [isHidden, setIsHidden] = useState(false);

  // Event Handlers
  const onLodgingSelect = (lodging: any) => {
    props.onLodgingSelect(lodging);
    setIsLogdingModalOpen(false);
  }

  // Renderers
  const renderHiddenToggle = () => {
    return (
      <button
        type="button"
        className={TripLogisticsCss.FlightsToggleBtn}
        onClick={() => {setIsHidden(!isHidden)}}
      >
      {isHidden ? <ChevronUpIcon className='h-4 w-4' />
        : <ChevronDownIcon className='h-4 w-4'/>}
      </button>
    );
  }

  const renderLodgings = () => {
    if (isHidden) {
      return null;
    }
    const lodgings = Object.values(props.trip.lodgings);
    return lodgings.map((lodge: any) => (
      <LodgingCard
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
        <div className={TripLogisticsCss.FlightsHeaderCtn}>
          {renderHiddenToggle()}
          <span>Hotels and Lodgings</span>
        </div>

        <button
          className={TripLogisticsCss.SearchFlightBtn}
          onClick={() => {setIsLogdingModalOpen(true)}}
        >
          +&nbsp;&nbsp;Add a lodging
        </button>
      </div>
      {renderLodgings()}
      <LodgingsModal
        trip={props.trip}
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

  const flightOnSelect = (flight: Trips.Flight) => {
    props.tripStateOnUpdate([
      TripsSyncAPI.makeAddOp(`/flights/${flight.id}`, flight)
    ]);
  }

  const flightOnDelete = (flight: Trips.Flight) => {
    props.tripStateOnUpdate([
      TripsSyncAPI.makeRemoveOp(`/flights/${flight.id}`, flight)
    ]);
  }

  // Event Handlers - Lodging

  const lodgingOnSelect = (lodging: Trips.Lodging) => {
    props.tripStateOnUpdate([
      TripsSyncAPI.makeAddOp(`/lodgings/${lodging.id}`, lodging)
    ]);
  }

  const lodgingOnUpdate = (lodging: any, updates: any) => {
    const ops = [] as any;
    Object.entries(updates).forEach(([key, value]) => {
      const fullpath = `/lodgings/${lodging.id}/${key}`;
      ops.push(TripsSyncAPI.newReplaceOp(fullpath, value));
    });
    props.tripStateOnUpdate(ops);
  }

  const lodgingOnDelete = (lodging: Trips.Lodging) => {
    props.tripStateOnUpdate([
      TripsSyncAPI.makeRemoveOp(`/lodgings/${lodging.id}`, lodging)
    ]);
  }

  // Renderers

  return (
    <div>
      <FlightsSection
        trip={props.trip}
        onFlightSelect={flightOnSelect}
        onFlightDelete={flightOnDelete}
      />
      <LodgingSection
        trip={props.trip}
        onLodgingSelect={lodgingOnSelect}
        onLodgingUpdate={lodgingOnUpdate}
        onLodgingDelete={lodgingOnDelete}
      />
    </div>
  );

}

export default TripLogisticsSection;
