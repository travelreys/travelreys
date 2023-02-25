import * as React from 'react'

export const ActionSetSelectedPlace = "setSelectedPlace";
type Action = {type: 'setSelectedPlace', value: any}
type Dispatch = (action: Action) => void
type State = {center: any, selectedPlace: any}
type MapsProviderProps = {children: React.ReactNode}


interface _MapsContext {
  state: State
  dispatch: Dispatch
}

const MapsContext = React.createContext<_MapsContext | undefined>(undefined);

const reducer = (state: any, action: any) => {
  switch (action.type) {
    case ActionSetSelectedPlace: {
      return {
        center: state.center,
        selectedPlace: action.value
      }
    }
    default: {
      throw new Error(`Unhandled action type: ${action.type}`)
    }
  }
}

export const MapsProvider = ({children}: MapsProviderProps) => {
  const [state, dispatch] = React.useReducer(
    reducer, {center: null, selectedPlace: null}
  );

  const value = {state, dispatch}
  return (
    <MapsContext.Provider value={value}>
      {children}
    </MapsContext.Provider>
  );
}

export const useMap = () => {
  const context = React.useContext(MapsContext)
  if (context === undefined) {
    throw new Error('useMap must be used within a MapsProvider')
  }
  return context;
}
