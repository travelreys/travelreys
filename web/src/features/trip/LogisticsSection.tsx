import React, { FC, useState } from 'react';
import {
  ChevronDownIcon,
  ChevronUpIcon,
} from '@heroicons/react/24/outline'

import FlightsModal from './FlightsModal';
import LodgingsModal from './LodgingsModal';
import LodgingCard from './LodgingCard';
import TransitFlightCard from './TransitFlightCard';

import {
  makeAddOp,
  makeRemoveOp,
  makeReplaceOp
} from '../../lib/jsonpatch';
import {  CommonCss } from '../../assets/styles/global';
import { Flight, Lodging } from '../../lib/trips';
import ToggleChevron from '../../components/common/ToggleChevron';

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

  const css = {
    titleCtn: "flex justify-between mb-4",
    headerCtn: "text-2xl sm:text-3xl font-bold text-slate-700",
    searchBtn: "text-slate-500 text-sm mt-1 font-bold",
  };

  const renderFlights = () => {
    if (isHidden) {
      return null;
    }
    return Object.values(props.trip.flights)
      .map((flight: any, idx: number) =>
        <TransitFlightCard
          key={idx}
          flight={flight}
          onDelete={props.onFlightDelete}
        />
      );
  }

  return (
    <div className='p-5'>
      <div className={css.titleCtn}>
        <div className={css.headerCtn}>
          <ToggleChevron isHidden={isHidden} onClick={() => setIsHidden(!isHidden)} />
          <span>Flights</span>
        </div>
        <button
          className={css.searchBtn}
          onClick={() => {setIsTripFlightsModalOpen(true)}}
        >
          +&nbsp;&nbsp;Search for a flight
        </button>
      </div>
      {renderFlights()}
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

  const css = {
    titleCtn: "flex justify-between mb-4",
    headerCtn: "text-2xl sm:text-3xl font-bold text-slate-700",
    searchBtn: "text-slate-500 text-sm mt-1 font-bold",
  };

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
      <div className={css.titleCtn}>
        <div className={css.headerCtn}>
          <ToggleChevron isHidden={isHidden} onClick={() => setIsHidden(!isHidden)} />
          <span>Hotels and Lodgings</span>
        </div>

        <button
          className={css.searchBtn}
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

interface LogisticsSectionProps {
  trip: any
  tripStateOnUpdate: any
}

const LogisticsSection: FC<LogisticsSectionProps> = (props: LogisticsSectionProps) => {

  // Event Handlers

  const flightOnSelect = (flight: Flight) => {
    props.tripStateOnUpdate([makeAddOp(`/flights/${flight.id}`, flight)]);
  }

  const flightOnDelete = (flight: Flight) => {
    props.tripStateOnUpdate([makeRemoveOp(`/flights/${flight.id}`, flight)]);
  }

  const lodgingOnSelect = (lodging: Lodging) => {
    props.tripStateOnUpdate([makeAddOp(`/lodgings/${lodging.id}`, lodging)]);
  }

  const lodgingOnUpdate = (lodging: any, updates: any) => {
    const ops = [] as any;
    Object.entries(updates).forEach(([key, value]) => {
      const fullpath = `/lodgings/${lodging.id}/${key}`;
      ops.push(makeReplaceOp(fullpath, value));
    });
    props.tripStateOnUpdate(ops);
  }

  const lodgingOnDelete = (lodging: Lodging) => {
    props.tripStateOnUpdate([makeRemoveOp(`/lodgings/${lodging.id}`, lodging)]);
  }

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

export default LogisticsSection;
