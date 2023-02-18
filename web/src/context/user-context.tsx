import * as React from 'react'
import { Auth, persistAuthUser } from '../lib/auth';

export const ActionSetUser = "setUser";
type Action = {type: 'setUser', value: Auth.User | null}
type Dispatch = (action: Action) => void
type State = {user: Auth.User | null }
type UserProviderProps = {children: React.ReactNode}


export const makeSetUserAction = (user: Auth.User|null): Action => {
  return {type: 'setUser', value: user}
}

interface _UserContext {
  state: State
  dispatch: Dispatch
}

const UserContext = React.createContext<_UserContext | undefined>(undefined);

const reducer = (state: State, action: Action) => {
  switch (action.type) {
    case ActionSetUser: {
      if (action.value) {
        persistAuthUser(action.value);
      }
      return { user: action.value };
    }
    default: {
      throw new Error(`Unhandled action type: ${action.type}`);
    }
  }
}

export const UserProvider = ({children}: UserProviderProps) => {
  const [state, dispatch] = React.useReducer(reducer, {user: null});

  const value = {state, dispatch}
  return (
    <UserContext.Provider value={value}>
      {children}
    </UserContext.Provider>
  );
}

export const useUser = () => {
  const context = React.useContext(UserContext)
  if (context === undefined) {
    throw new Error('useUser must be used within a UserProvider')
  }
  return context;
}
