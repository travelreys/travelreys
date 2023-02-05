import React, { FC, useState } from 'react';
import _get from "lodash/get";
import _sortBy from "lodash/sortBy";
import _isEmpty from "lodash/isEmpty";

import {
  ChevronDownIcon,
  ChevronUpIcon,
  FolderArrowDownIcon
} from '@heroicons/react/24/outline'

import { Trips } from '../../apis/types';
import TripsSyncAPI from '../../apis/tripsSync';
import { TripNodesCss } from '../../styles/global';
import NotesEditor from '../NotesEditor';

interface TripNotesSectionProps {
  trip: any
  tripStateOnUpdate: any
}

const TripNotesSection: FC<TripNotesSectionProps> = (props: TripNotesSectionProps) => {

  const [isHidden, setIsHidden] = useState(false);

  // Event Handlers
  const notesOnChange = (content: string) => {
    console.log(content);
    const ops = [];
    ops.push(TripsSyncAPI.makeReplaceOp(`/notes`, content));
    props.tripStateOnUpdate(ops);
  }

  // Renderers
  const renderHiddenToggle = () => {
    return (
      <button
        type="button"
        className={TripNodesCss.ToggleBtn}
        onClick={() => {setIsHidden(!isHidden)}}
      >
      {isHidden ? <ChevronUpIcon className='h-4 w-4' />
        : <ChevronDownIcon className='h-4 w-4'/> }
      </button>
    );
  }

  return (
    <div className='p-5'>
      <div className={TripNodesCss.TitleCtn}>
        <div className={TripNodesCss.HeaderCtn}>
          {renderHiddenToggle()}
          <span>Notes </span>
        </div>
      </div>
      <NotesEditor
        base64Notes={props.trip.notes}
        notesOnChange={notesOnChange}
      />
    </div>
  );

}

export default TripNotesSection;
