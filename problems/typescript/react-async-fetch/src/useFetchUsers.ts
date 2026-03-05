import { useState, useEffect } from 'react';

export interface User {
  id: number;
  name: string;
  email: string;
}

interface FetchState<T> {
  data: T | null;
  loading: boolean;
  error: string | null;
}

// Simulated API — returns users after a delay
export async function fetchUsers(query: string): Promise<User[]> {
  return new Promise(resolve => {
    setTimeout(() => {
      const allUsers: User[] = [
        { id: 1, name: 'Alice Johnson', email: 'alice@example.com' },
        { id: 2, name: 'Bob Smith', email: 'bob@example.com' },
        { id: 3, name: 'Charlie Brown', email: 'charlie@example.com' },
        { id: 4, name: 'Diana Prince', email: 'diana@example.com' },
      ];
      if (!query) {
        resolve(allUsers);
        return;
      }
      resolve(allUsers.filter(u =>
        u.name.toLowerCase().includes(query.toLowerCase())
      ));
    }, 100);
  });
}

// Custom hook with bugs
// BUG 1: Loading state is never set to false after fetch completes
// BUG 2: No cleanup — setting state after unmount causes crash
// BUG 3: Missing dependency in useEffect — stale query causes wrong results
export function useFetchUsers(query: string): FetchState<User[]> {
  const [state, setState] = useState<FetchState<User[]>>({
    data: null,
    loading: true,
    error: null,
  });

  useEffect(() => {
    setState(prev => ({ ...prev, loading: true, error: null }));

    // BUG: No abort/cancelled flag — state update after unmount
    fetchUsers(query)
      .then(users => {
        // BUG: Never sets loading to false
        setState(prev => ({ ...prev, data: users }));
      })
      .catch(err => {
        setState(prev => ({ ...prev, error: err.message }));
      });

    // BUG: No cleanup function to prevent state update after unmount
  }, []); // BUG: Empty dependency array — doesn't refetch when query changes

  return state;
}
