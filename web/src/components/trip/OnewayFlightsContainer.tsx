import React, { FC } from 'react';
import _get from "lodash/get";
import _sortBy from "lodash/sortBy";
import _isEmpty from "lodash/isEmpty";

import { FlightsModalCss } from '../../styles/global';

import TripFlightCard from './FlightCard';

// OnewayFlightsContainer

interface OnewayFlightsContainerProps {
  oneways: any
  onSelect: any
}

const OnewayFlightsContainer: FC<OnewayFlightsContainerProps> = (props: OnewayFlightsContainerProps) => {

  return (
    <div>
      <p className={FlightsModalCss.FlightSearchResultsTitle}>
        One-way Flights
      </p>
      {props.oneways.map((flight: any, idx: number) =>
        <TripFlightCard
          key={idx}
          flight={flight.depart}
          bookingMetadata={flight.bookingMetadata}
          onSelect={props.onSelect}
        />
      )}
    </div>
  );
}

export default OnewayFlightsContainer;
