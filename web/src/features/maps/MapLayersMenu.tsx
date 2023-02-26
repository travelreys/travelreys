import React, { FC } from 'react';
import _get from "lodash/get";
import _flatten from "lodash/flatten";
import _isEmpty from "lodash/isEmpty";
import _find from "lodash/find";

import ListPin from './ListPin';
import {
  ActivityIconOpts,
  DefaultActivityColor,
  getActivityColor,
  getActivityIcon,
  LabelUiColor,

} from '../../lib/trips';
import { parseISO, fmt } from '../../lib/dates'

interface MapLayersMenuProps {
  trip: any
  layersViz: any
  onSelect: (layersSelected: any) => void
}

const MapLayersMenu: FC<MapLayersMenuProps> = (props: MapLayersMenuProps) => {
  // Event Handlers

  const selectAllOnClick = () => {
    const newLayersViz = Object.assign({}, props.layersViz);
    Object.keys(newLayersViz).forEach((k: string) => {
      newLayersViz[k] = true;
    });
    props.onSelect(newLayersViz);
  }

  const deselectAllOnClick = () => {
    const newLayersViz = Object.assign({}, props.layersViz);
    Object.keys(newLayersViz).forEach((k: string) => {
      newLayersViz[k] = false;
    });
    props.onSelect(newLayersViz);
  }

  const onSelectLayer = (id: string) => {
    const newLayersViz = Object.assign({}, props.layersViz);
    newLayersViz[id] = !newLayersViz[id];
    props.onSelect(newLayersViz);
  }

  // Renderers

  const renderHomeOpts = () => {
    const activities = [
      {
        id: "lodgings",
        name: "Hotels and Lodgings",
        labels: { "ui|color": "rgb(249 115 22)", "ui|icon": "hotel"}
      }
    ].concat(Object.values(_get(props.trip, "activities", [])))

    const layers = activities
      .map((l: any) => {
        const color = getActivityColor(l) || DefaultActivityColor;
        const icon = ActivityIconOpts[getActivityIcon(l)];
        return (
          <div key={l.id} className='flex justify-between items-center mb-1'>
            <div className='flex items-center flex-1 text-gray-700'>
              <ListPin icon={icon} color={color}/>
              {l.name}
            </div>
            <input
              type="checkbox"
              checked={props.layersViz[l.id]}
              className="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 rounded focus:ring-blue-500"
              onChange={() => {onSelectLayer(l.id)}}
            />
          </div>
        )
      });
    return (
      <div className='w-full'>
        <div className='mb-1'>
          <h4 className='text-sm font-bold mb-1'>Overview</h4>
          {layers}
        </div>
      </div>
    );
  }

  const renderItineraryOpts = () => {
    const itinLayers = _get(props.trip, "itinerary", [])
      .map((l: any) => {
        const color = _get(l, `labels.${LabelUiColor}`, DefaultActivityColor);
        return (
          <div key={l.id} className='flex justify-between items-center mb-1'>
            <div className='flex items-center flex-1 text-gray-700'>
              <ListPin icon={""} color={color}/>
              { fmt(parseISO(l.date), "eee, MM/dd")  }
            </div>
            <input
              type="checkbox"
              checked={props.layersViz[l.id]}
              className="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 rounded focus:ring-blue-500"
              onChange={() => {onSelectLayer(l.id)}}
            />
          </div>
        )
      });
    return (
      <div className='w-full'>
        <div className='mb-1'>
          <h4 className='text-sm font-bold mb-1'>Itinerary</h4>
          {itinLayers}
        </div>
      </div>
    );
  }

  return (
    <div className='absolute mt-8 ml-12 bg-white p-4 rounded-xl z-10 w-96'>
      <h4 className='font-bold text-base mb-2'>Map Drawings</h4>
      <div className='flex mb-1 text-sm'>
        <button
          className='font-semibold text-indigo-500 mr-2 hover:text-indigo-700'
          onClick={selectAllOnClick}
        >
          Select all
        </button>
        <button
          className='font-semibold text-gray-500 hover:text-indigo-700'
          onClick={deselectAllOnClick}
        >
          Deselect all
        </button>
      </div>
      <hr className='my-2' />
      {renderHomeOpts()}
      {renderItineraryOpts()}
    </div>
  );
}

export default MapLayersMenu;
