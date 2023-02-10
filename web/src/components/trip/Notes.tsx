import React, { FC, useState } from 'react';
import _get from "lodash/get";
import _sortBy from "lodash/sortBy";
import _isEmpty from "lodash/isEmpty";

import {
  ChevronDownIcon,
  ChevronUpIcon,
  FolderArrowDownIcon
} from '@heroicons/react/24/outline'

import TripsSyncAPI from '../../apis/tripsSync';
import { TripNotesCss } from '../../styles/global';
import NotesEditor from '../NotesEditor';
import ToggleChevron from '../ToggleChevron';

interface NotesSectionProps {
  trip: any
  tripStateOnUpdate: any
}

const NotesSection: FC<NotesSectionProps> = (props: NotesSectionProps) => {

  const [isHidden, setIsHidden] = useState(false);

  // Event Handlers
  const notesOnChange = (content: string) => {
    const ops = [];
    ops.push(TripsSyncAPI.newReplaceOp(`/notes`, content));
    props.tripStateOnUpdate(ops);
  }

  // Renderers

  return (
    <div className='p-5'>
      <div className={TripNotesCss.TitleCtn}>
        <div className={TripNotesCss.HeaderCtn}>
          <ToggleChevron
            isHidden={isHidden}
            onClick={() => setIsHidden(!isHidden)}
          />
          <span>Notes</span>
        </div>
      </div>
      {isHidden ? null :
        <NotesEditor
          base64Notes={props.trip.notes}
          notesOnChange={notesOnChange}
        />
      }
    </div>
  );

}

export default NotesSection;
