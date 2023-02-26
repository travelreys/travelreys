import React, { FC, useState } from 'react';
import { TripNotesCss } from '../../assets/styles/global';
import NotesEditor from '../../components/common/NotesEditor';
import ToggleChevron from '../../components/common/ToggleChevron';
import { makeRepOp } from '../../lib/jsonpatch';

interface NotesSectionProps {
  trip: any
  tripOnUpdate: any
}

const NotesSection: FC<NotesSectionProps> = (props: NotesSectionProps) => {

  const [isHidden, setIsHidden] = useState(false);

  // Event Handlers
  const notesOnChange = (content: string) => {
    props.tripOnUpdate([makeRepOp(`/notes`, content)]);
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
