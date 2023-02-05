import React, {
  FC,
  Ref,
  PropsWithChildren,
  useEffect,
  useMemo,
  useCallback,
  useRef
} from 'react';
import ReactDOM from 'react-dom';
import isHotkey from 'is-hotkey';
import {
  Editor,
  Text,
  Transforms,
  createEditor,
  Descendant,
  Range,
  Element as SlateElement,
} from 'slate'
import {
  Slate,
  Editable,
  useSlate,
  useFocused,
  withReact
} from 'slate-react'
import { withHistory } from 'slate-history'
import BoldIcon from './icons/BoldIcon';
import ItalicIcon from './icons/ItalicIcon';
import UnderlineIcon from './icons/UnderlineIcon';
import HeaderOneIcon from './icons/HeaderOneIcon';
import NumberlistIcon from './icons/NumberlistIcon';
import { ListBulletIcon } from '@heroicons/react/24/outline';
import HeaderTwoIcon from './icons/HeaderTwoIcon';


const HOTKEYS = {
  'mod+b': 'bold',
  'mod+i': 'italic',
  'mod+u': 'underline',
} as any;

const LIST_TYPES = ['numbered-list', 'bulleted-list']
const TEXT_ALIGN_TYPES = ['left', 'center', 'right', 'justify']

interface NotesEditorProps {
  placeholder?: string
}


const NotesEditor: FC<NotesEditorProps> = (props: NotesEditorProps) => {

  const renderElement = useCallback((props: ElementProps) => <Element {...props} />, [])
  const renderLeaf = useCallback((props: LeafProps) => <Leaf {...props} />, [])
  const editor = useMemo(() => withHistory(withReact(createEditor())), [])

  const initialValue = [
    {
      type: 'paragraph',
      children: [{ text: 'A line of text in a paragraph.' }],
    },
  ]

  // Helpers
  const isMarkActive = (editor: any, format: any) => {
    const marks = Editor.marks(editor) as any;
    return marks ? marks[format] === true : false
  }

  const toggleMark = (editor: any, format: any) => {
    if (isMarkActive(editor, format)) {
      Editor.removeMark(editor, format)
    } else {
      Editor.addMark(editor, format, true)
    }
  }

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
      editor,
      { [format]: isActive ? null : true },
      { match: Text.isText, split: true }
    )
  }

  // Event Handlers

  const onKeyDown = (e :any) => {
    for (const hotkey in HOTKEYS) {
      if (isHotkey(hotkey, e as any)) {
        e.preventDefault()
        const mark = HOTKEYS[hotkey];
        toggleMark(editor, mark)
      }
    }
  }

  return (
    <div className='p-4 bg-gray-50'>
      <Slate editor={editor} value={initialValue}>
        <HoveringToolbar />
        <Editable
          renderElement={renderElement}
          renderLeaf={renderLeaf}
          placeholder={props.placeholder || "Enter notes..."}
          spellCheck
          onKeyDown={onKeyDown}
          onDOMBeforeInput={(event: InputEvent) => {
            switch (event.inputType) {
              case 'formatBold':
                event.preventDefault()
                return toggleFormat(editor, 'bold')
              case 'formatItalic':
                event.preventDefault()
                return toggleFormat(editor, 'italic')
              case 'formatUnderline':
                event.preventDefault()
                return toggleFormat(editor, 'underlined')
            }
          }}
        />
      </Slate>

    </div>
  );
}

export default NotesEditor;


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
      )
    case 'heading-one':
      return (
        <h1 style={style} className="text-xl font-bold" {...attributes}>
          {children}
        </h1>
      )
    case 'heading-two':
      return (
        <h2 style={style} className="text-lg font-bold" {...attributes}>
          {children}
        </h2>
      )
    case 'list-item':
      return (
        <li style={style} {...attributes}>
          {children}
        </li>
      )
    case 'numbered-list':
      return (
        <ol style={style} className="list-decimal" {...attributes}>
          {children}
        </ol>
      )
    default:
      return (
        <p style={style} {...attributes}>
          {children}
        </p>
      )
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
        className={`inline-block p-2 absolute z-50 -mt-3 bg-gray-800 rounded-lg`}
        onMouseDown={onMouseDown}
      >
        {renderFormatButtons()}
        {renderBlockButtons()}
      </div>
    </Portal>
  )
}


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
      className={'cursor-pointer p-1 mr-1 rounded' + (isActive ? " bg-indigo-500" : "")}
      onClick={onMouseDown}
    >
      <props.icon
        className='h-4 w-4  stroke-white'
      />
    </button>
  )

}


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
      className={'cursor-pointer p-1 mr-1 rounded' + (isActive ? " bg-indigo-500" : "")}
      onClick={() => toggleFormat(editor, props.format)}
    >
      <props.icon
        className='h-4 w-4 rounded stroke-white'
      />
    </button>
  )
}
