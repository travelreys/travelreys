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

import { Trips } from '../../apis/types';
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
      <TripFlightsModal
        trip={props.trip}
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
      <TripLodgingsModal
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
    const ops = [];
    ops.push(TripsSyncAPI.makeAddOp(`/flights/${flight.id}`, flight));
    props.tripStateOnUpdate(ops);
  }

  const flightOnDelete = (flight: Trips.Flight) => {
    const ops = [];
    ops.push(TripsSyncAPI.makeRemoveOp(`/flights/${flight.id}`, flight))
    props.tripStateOnUpdate(ops);
  }

  // Event Handlers - Lodging

  const lodgingOnSelect = (lodging: Trips.Lodging) => {
    const ops = [];
    ops.push(TripsSyncAPI.makeAddOp(`/lodgings/${lodging.id}`, lodging))
    props.tripStateOnUpdate(ops);
  }

  const lodgingOnUpdate = (lodging: any, updates: any) => {
    const ops = [] as any;
    Object.entries(updates).forEach(([key, value]) => {
      const fullpath = `/lodgings/${lodging.id}/${key}`;
      ops.push(TripsSyncAPI.makeReplaceOp(fullpath, value));
    });
    props.tripStateOnUpdate(ops);
  }

  const lodgingOnDelete = (lodging: Trips.Lodging) => {
    const ops = [];
    ops.push(TripsSyncAPI.makeRemoveOp(`/lodgings/${lodging.id}`, lodging))
    props.tripStateOnUpdate(ops);
  }

  // Renderers

  const renderTabs = () => {
    const tabs = [
      { title: "Flights", icon: PlaneIcon },
      { title: "Transits", icon: BusIcon },
      { title: "Lodging", icon: HotelIcon },
      { title: "Attachments", icon: FolderArrowDownIcon },
    ];

    return (
      <div className={TripLogisticsCss.TabsCtn}>
        <div className={TripLogisticsCss.TabsWrapper}>
          <h5 className={TripLogisticsCss.TabsCtnHeader}>
            Logistics
          </h5>
          <div className={TripLogisticsCss.TabItemCtn}>
            {tabs.map((tab: any, idx: number) => (
              <button
                key={idx} type="button"
                className={TripLogisticsCss.TabItemBtn}
              >
                <tab.icon className='h-6 w-6 mb-1'/>
                <span className={TripLogisticsCss.TabItemBtnTxt}>
                  {tab.title}
                </span>
              </button>
            ))}
          </div>
        </div>
      </div>
    );
  }

  return (
    <div>
      {/* {renderTabs()} */}
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
