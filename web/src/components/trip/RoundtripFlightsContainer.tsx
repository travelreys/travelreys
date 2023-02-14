import React, { FC, useState } from 'react';
import _get from "lodash/get";
import _sortBy from "lodash/sortBy";
import _minBy from "lodash/minBy";
import _isEmpty from "lodash/isEmpty";

import { FlightsModalCss } from '../../styles/global';

import TripFlightCard from './FlightCard';

// RoundtripFlightsContainer

interface RoundtripFlightsContainerProps {
  roundtrips: any
  onSelect: any
}

const RoundtripFlightsContainer: FC<RoundtripFlightsContainerProps> = (props: RoundtripFlightsContainerProps) => {

  const [stepperStep, setStepperStep] = useState(0);
  const [roundtrip, setRoundtrip] = useState(null as any);

  let roundtrips = Object.values(props.roundtrips).map((rt: any) => {
    const minScore = _minBy(rt.bookingMetadata, (bm: any) => bm.score)
    return Object.assign(rt, {score: minScore});
  });
  roundtrips = _sortBy(roundtrips, (rt: any) => rt.score.score);

  // Event Handlers
  const selectDepartingFlightsStepperOnClick = () => {
    setStepperStep(0);
    setRoundtrip(null);
  }

  const departFlightOnSelect = (flight: any, _: any) => {
    setRoundtrip(_get(props.roundtrips, flight.id));
    setStepperStep(1);
  }

  const returnFlightOnSelect = (flight: any, bookingMetadata: any) => {
    props.onSelect(roundtrip.depart, flight, bookingMetadata)
  }

  // Renderers

  const renderStepper = () => {
    const texts = [
      <span
        className='cursor-pointer'
        onClick={selectDepartingFlightsStepperOnClick}
      >
        Select Departing Flights&nbsp;&nbsp;&gt;
      </span>,
      <span>Select Return Flights</span>
    ]
    return (
      <ol className={FlightsModalCss.RoundTripStepperCtn}>
        {texts.map((text: any, idx: number) => {
          const css = idx === stepperStep
            ? FlightsModalCss.RoundTripStepperActive: FlightsModalCss.RoundTripStepper
          return (<li className={css}>{text}</li>);
        })}
      </ol>
    );
  }

  const renderDepartingFlights = () => {
    const flights = roundtrips.map((rt: any) => rt.depart);
    return (
      <div>
        {flights.map((flight: any, idx: number) =>
          <TripFlightCard
            key={idx}
            flight={flight}
            onSelect={departFlightOnSelect}
            bookingMetadata={null}
          />
        )}
      </div>
    );
  }

  const renderReturnFlights = () => {
    roundtrip.returns.forEach((rt: any, idx: number) => {
      rt.bookingMetadata = roundtrip.bookingMetadata[idx];
    })
    const returnFlights = _sortBy(roundtrip.returns, (rt: any) => rt.bookingMetadata.score);

    return (
      <div>
        {returnFlights.map((flight: any, idx: number) =>
          <TripFlightCard
            key={idx}
            flight={flight}
            onSelect={returnFlightOnSelect}
            bookingMetadata={flight.bookingMetadata}
          />
        )}
      </div>
    );
  }

  return (
    <div>
      {renderStepper()}
      { stepperStep === 0 ? renderDepartingFlights(): null }
      { stepperStep === 1 ? renderReturnFlights(): null }
    </div>
  );
}

export default RoundtripFlightsContainer;
