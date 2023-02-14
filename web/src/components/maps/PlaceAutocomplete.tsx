import React from 'react';
import _get from "lodash/get";
import _isEmpty from "lodash/isEmpty";

import {
  MapPinIcon,
} from '@heroicons/react/24/solid'
import {
  PlaceAutocompleteCss,
} from '../../styles/global';

interface PlaceAutocompleteProps {
  predictions: Array<any>
  onSelect: (pred: string) => void
}

const PlaceAutocomplete: React.FC<PlaceAutocompleteProps> = (props: PlaceAutocompleteProps) => {
  if (_isEmpty(props.predictions)) {
    return (<></>);
  }

  return (
    <div className={PlaceAutocompleteCss.AutocompleteCtn}>
      {props.predictions.map((pre: any) => (
        <div
          className={PlaceAutocompleteCss.PredictionWrapper}
          key={pre.place_id}
          onClick={() => {props.onSelect(pre.place_id)}}
        >
          <div className={PlaceAutocompleteCss.IconCtn}>
            <MapPinIcon className={PlaceAutocompleteCss.Icon} />
          </div>
          <div className='ml-1'>
            <p className={PlaceAutocompleteCss.PrimaryTxt}>
              {_get(pre, "structured_formatting.main_text", "")}
            </p>
            <p className={PlaceAutocompleteCss.SecondaryTxt}>
              {_get(pre, "structured_formatting.secondary_text", "")}
            </p>
          </div>
        </div>
      ))}
    </div>
  );

}

export default PlaceAutocomplete;
