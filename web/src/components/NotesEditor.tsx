import React, {
  FC,
  useEffect,
  useMemo,
  useCallback,
  useRef,
  useState
} from 'react';
import ReactDOM from 'react-dom';
import {
  Editor,
  Text,
  Transforms,
  Descendant,
  Point,
  Range,
  Element as SlateElement,
  createEditor,
} from 'slate';
import {
  Slate,
  Editable,
  ReactEditor,
  useSlate,
  useFocused,
  useSlateStatic,
  useReadOnly,
  withReact,
} from 'slate-react';
import { withHistory } from 'slate-history';
import _isEmpty from "lodash/isEmpty";
import { encode, decode } from 'js-base64';

import {  ListBulletIcon } from '@heroicons/react/24/outline';

import BoldIcon from './icons/BoldIcon';
import ItalicIcon from './icons/ItalicIcon';
import UnderlineIcon from './icons/UnderlineIcon';
import HeaderOneIcon from './icons/HeaderOneIcon';
import NumberlistIcon from './icons/NumberlistIcon';
import HeaderTwoIcon from './icons/HeaderTwoIcon';
import StrikethroughIcon from './icons/StrikethroughIcon';
import ChecklistIcon from './icons/ChecklistIcon';


const LIST_TYPES = ['numbered-list', 'bulleted-list']
const TEXT_ALIGN_TYPES = ['left', 'center', 'right', 'justify']
const defaultEditorValues = [
  { type: 'paragraph', children: [{ text: '' }] },
];


interface NotesEditorProps {
  base64Notes: string
  ctnCss?: string
  placeholder?: string
  notesOnChange: any
  readOnly?: boolean
}

const NotesEditor: FC<NotesEditorProps> = (props: NotesEditorProps) => {

  const renderElement = useCallback((props: ElementProps) => <Element {...props} />, []);
  const renderLeaf = useCallback((props: LeafProps) => <Leaf {...props} />, []);
  const editor = useMemo(() => withChecklists(withHistory(withReact(createEditor()))), []);

  const currentVal = useRef<string>(props.base64Notes);
  const initialValue = defaultEditorValues // JSON.parse(decode(props.base64Notes)));

  useEffect(() => {
    Transforms.removeNodes(editor, {
      at: {
        anchor: Editor.start(editor, []),
        focus: Editor.end(editor, []),
      },
    });
    const nodes = _isEmpty(props.base64Notes)
      ? defaultEditorValues : JSON.parse(decode(props.base64Notes))
    Transforms.insertNodes(editor, nodes, {})
  }, [props.base64Notes])

  // Event Handlers
  const editorOnBlur = () => {
    props.notesOnChange(currentVal.current);
  }

  const editorOnChange = (value: any) => {
    const isAstChange = editor.operations.some(
      op => 'set_selection' !== op.type
    )
    if (isAstChange) {
      // Save the value to Local Storage.
      const content = encode(JSON.stringify(value));
      currentVal.current = content;
    }
  }

  return (
    <div
      className={props.ctnCss || 'p-4 bg-gray-50'}
      onBlur={editorOnBlur}
    >
      <Slate
        editor={editor}
        value={initialValue}
        onChange={editorOnChange}
      >
        <HoveringToolbar />
        <Editable
          renderElement={renderElement}
          renderLeaf={renderLeaf}
          spellCheck
          readOnly={props.readOnly}
          placeholder={props.placeholder}
        />
      </Slate>
    </div>
  );
}

export default NotesEditor;

// withChecklists

const withChecklists = (editor: ReactEditor) => {
  const { deleteBackward } = editor

  editor.deleteBackward = (...args) => {
    const { selection } = editor

    if (selection && Range.isCollapsed(selection)) {
      //@ts-ignore
      const [match] = Editor.nodes(editor, {
        match: n =>
          !Editor.isEditor(n) &&
          SlateElement.isElement(n) &&
          //@ts-ignore
          n.type === 'check-list-item',
      })

      if (match) {
        const [, path] = match
        const start = Editor.start(editor, path)

        if (Point.equals(selection.anchor, start)) {
          const newProperties: Partial<SlateElement> = {
            type: 'paragraph',
          } as any;
          Transforms.setNodes(editor, newProperties, {
            match: n =>
              !Editor.isEditor(n) &&
              SlateElement.isElement(n) &&
              //@ts-ignore
              n.type === 'check-list-item',
          })
          return
        }
      }
    }

    deleteBackward(...args)
  }

  return editor
}


// Element

interface ElementProps {
  attributes: any
  children: any
  element: any
}


const Element = (props: ElementProps) => {
  const { attributes, children, element } = props;
  const style = { textAlign: element.align }

  switch (element.type) {
    case 'bulleted-list':
      return (
        <ul style={style} className="list-disc list-inside" {...attributes}>
          {children}
        </ul>
      );
    case 'heading-one':
      return (
        <h1 style={style} className="text-xl font-bold" {...attributes}>
          {children}
        </h1>
      );
    case 'heading-two':
      return (
        <h2 style={style} className="text-lg font-bold" {...attributes}>
          {children}
        </h2>
      );
    case 'list-item':
      return (
        <li style={style} {...attributes}>
          {children}
        </li>
      );
    case 'numbered-list':
      return (
        <ol style={style} className="list-decimal" {...attributes}>
          {children}
        </ol>
      );
    case 'check-list-item':
      return <CheckListItemElement {...props} />;
    default:
      return (
        <p style={style} {...attributes}>
          {children}
        </p>
      );
  }
}

// Leaf

interface LeafProps {
  attributes: any
  children: any
  leaf: any
}

const Leaf = (props: LeafProps) => {
  let { attributes, children, leaf } = props;
  if (leaf.bold) {
    children = <strong>{children}</strong>
  }
  if (leaf.italic) {
    children = <em>{children}</em>
  }
  if (leaf.underlined) {
    children = <u>{children}</u>
  }
  if (leaf.strikethrough) {
    children = <s>{children}</s>
  }
  return <span {...attributes}>{children}</span>
}

// Portal

interface PortalProps {
  children: React.ReactNode
}

export const Portal: FC<PortalProps> = (props: PortalProps) => {
  const { children } = props;
  return typeof document === 'object'
    ? ReactDOM.createPortal(children, document.body)
    : null
}

// HoveringToolbar

interface HoveringToolbarProps {}

const HoveringToolbar: FC<HoveringToolbarProps> = (props: HoveringToolbarProps) => {
  const ref = useRef() as any;
  const editor = useSlate();
  const inFocus = useFocused();

  useEffect(() => {
    const el = ref.current;
    const { selection } = editor;

    if (!el) {
      return;
    }

    if (
      !selection ||
      !inFocus ||
      Range.isCollapsed(selection) ||
      Editor.string(editor, selection) === ''
    ) {
      el.removeAttribute('style');
      return;
    }

    const domSelection = window.getSelection();
    if (domSelection) {
      const domRange = domSelection.getRangeAt(0);
      const rect = domRange.getBoundingClientRect();
      el.style.opacity = '1'
      el.style.top = `${rect.top + window.pageYOffset - el.offsetHeight}px`
      let left = rect.left + window.pageXOffset - el.offsetWidth / 2 + rect.width / 2;
      el.style.left = `${left <= 0 ? 0 : left}px`
    }
  })

  // Event Handlers
  const onMouseDown = (e: any) => {
    // prevent toolbar from taking focus away from editor
    e.preventDefault();
  }

  // Renderers
  const renderFormatButtons = () => {
    return [
      {format: "bold", icon: BoldIcon},
      {format: "italic", icon: ItalicIcon},
      {format: "underlined", icon: UnderlineIcon},
      {format: "strikethrough", icon: StrikethroughIcon}
    ].map((btn: any) => (
      <FormatButton
        key={btn.format}
        format={btn.format}
        icon={btn.icon}
      />
    ));
  }

  const renderBlockButtons = () => {
    return [
      {format: "heading-one", icon: HeaderOneIcon},
      {format: "heading-two", icon: HeaderTwoIcon},
      {format: "bulleted-list", icon: ListBulletIcon},
      {format: "numbered-list", icon: NumberlistIcon},
      {format: "check-list-item", icon: ChecklistIcon},
    ].map((btn: any) => (
      <BlockButton
        key={btn.format}
        format={btn.format}
        icon={btn.icon}
      />
    ));
  }

  return (
    <Portal>
      <div
        ref={ref}
        className="inline-block p-2 absolute z-50 opacity-0 -top-[10000px] -left-[10000px] -mt-3 bg-gray-800 rounded-lg"
        onMouseDown={onMouseDown}
      >
        {renderFormatButtons()}
        <span className='border-l border-white'>&nbsp;</span>
        {renderBlockButtons()}
      </div>
    </Portal>
  )
}


// Block Buttons

interface BlockButtonProps {
  format: string
  icon: any
}

const BlockButton: FC<BlockButtonProps> = (props: BlockButtonProps) => {
  const editor = useSlate();

  const isBlockActive = (edt: any, format: string, blockType: any = 'type') => {
    const { selection } = edt
    if (!selection) return false

    const [match] = Array.from(
      Editor.nodes(edt, {
        at: Editor.unhangRange(edt, selection),
        match: (n: any) =>
          !Editor.isEditor(n) &&
          SlateElement.isElement(n) &&
          //@ts-ignore
          n[blockType] === format,
      })
    )
    return !!match
  }


  const toggleBlock = (edt: any, format: string) => {
    const isActive = isBlockActive(
      edt,
      format,
      TEXT_ALIGN_TYPES.includes(format) ? 'align' : 'type'
    )
    const isList = LIST_TYPES.includes(format)

    Transforms.unwrapNodes(edt, {
      match: (n: any) =>
        !Editor.isEditor(n) &&
        SlateElement.isElement(n) &&
        //@ts-ignore
        LIST_TYPES.includes(n.type) &&
        !TEXT_ALIGN_TYPES.includes(format),
      split: true,
    })
    let newProperties: Partial<SlateElement>
    if (TEXT_ALIGN_TYPES.includes(format)) {
      //@ts-ignore
      newProperties = { align: isActive ? undefined : format }
    } else {
      newProperties = {
        //@ts-ignore
        type: isActive ? 'paragraph' : isList ? 'list-item' : format,
      }
    }
    Transforms.setNodes<SlateElement>(edt, newProperties)

    if (!isActive && isList) {
      const block = { type: format, children: [] }
      Transforms.wrapNodes(edt, block)
    }
  }


  // Event Handlers
  const onMouseDown = (e: any) => {
    e.preventDefault()
    toggleBlock(editor, props.format)
  }

  const isActive = isBlockActive(
    editor,
    props.format,
    TEXT_ALIGN_TYPES.includes(props.format) ? 'align' : 'type'
  )

  return (
    <button
      type="button"
      className={'cursor-pointer p-1 mr-1 rounded'
        + (isActive ? " bg-indigo-500" : "")}
      onClick={onMouseDown}
    >
      <props.icon
        className='h-4 w-4  stroke-white'
      />
    </button>
  )

}


// FormatButton

interface FormatButtonProps {
  format: string
  icon: any
}

const FormatButton: FC<FormatButtonProps> = (props: FormatButtonProps) => {
  const editor = useSlate();

  const isFormatActive = (edt: any, format: any) => {
    const [match]: any = Editor.nodes(edt, {
      match: (n: any) => n[format] === true,
      mode: 'all',
    })
    return !!match
  }

  const toggleFormat = (edt: any, format: any) => {
    const isActive = isFormatActive(edt, format)
    Transforms.setNodes(
      edt,
      { [format]: isActive ? null : true },
      { match: Text.isText, split: true }
    )
  }

  const isActive = isFormatActive(editor, props.format);

  return (
    <button
      type="button"
      className={
        'cursor-pointer p-1 mr-1 rounded'
          + (isActive ? " bg-indigo-500" : "")}
      onClick={() => toggleFormat(editor, props.format)}
    >
      <props.icon
        className='h-4 w-4 rounded stroke-white'
      />
    </button>
  )
}

// ChecklistElement

interface ChecklistElementProps {
  attributes: any
  children: any
  element: any
}

const CheckListItemElement: FC<ChecklistElementProps> = (props: ChecklistElementProps) => {
  const editor = useSlateStatic() as any;
  const readOnly = useReadOnly();
  const { checked } = props.element;

  const spanCss = checked ? "flex-1 text-gray-500 line-through" : "flex-1";

  return (
    <div
      {...props.attributes}
      className="flex items-center mt-2"
    >
      <span
        contentEditable={false}
        className="mr-2"
      >
        <input
          type="checkbox"
          checked={checked}
          className="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 rounded focus:ring-blue-500"
          onChange={event => {
            const path = ReactEditor.findPath(editor, props.element)
            const newProperties: Partial<SlateElement> = {
              checked: event.target.checked,
            } as any;
            Transforms.setNodes(editor, newProperties, { at: path })
          }}
        />
      </span>
      <span
        contentEditable={!readOnly}
        suppressContentEditableWarning
        className={spanCss}
      >
        {props.children}
      </span>
    </div>
  )
}