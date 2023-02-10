import * as React from 'react'

type Action = {type: 'setSelectedPlace', value: any}

export const ActionNameSetSelectedPlace = "setSelectedPlace";

type Dispatch = (action: Action) => void
type State = {center: any, selectedPlace: any}
type MapsProviderProps = {children: React.ReactNode}


const MapsContext = React.createContext<
  {state: State; dispatch: Dispatch} | undefined
>(undefined);

const mapsReducer = (state: any, action: any) => {
  switch (action.type) {
    case 'setSelectedPlace': {
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
    mapsReducer,
    {center: null, selectedPlace: null}
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
  return context
}
