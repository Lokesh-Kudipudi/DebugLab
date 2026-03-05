import React, { createContext, useContext, useReducer, ReactNode } from 'react';

// Types
export interface CartItem {
  id: string;
  name: string;
  price: number;
  quantity: number;
}

export interface CartState {
  items: CartItem[];
  discount: number; // percentage, e.g. 10 means 10%
  discountCode: string | null;
}

export type CartAction =
  | { type: 'ADD_ITEM'; payload: CartItem }
  | { type: 'REMOVE_ITEM'; payload: string }
  | { type: 'APPLY_DISCOUNT'; payload: { code: string; percentage: number } }
  | { type: 'REMOVE_DISCOUNT' };

// Valid discount codes
export const DISCOUNT_CODES: Record<string, number> = {
  SAVE10: 10,
  SAVE20: 20,
  HALF: 50,
  FREE100: 100,
};

const initialState: CartState = {
  items: [],
  discount: 0,
  discountCode: null,
};

// BUG 1: Direct state mutation in reducer — the reducer mutates the state
// object directly instead of returning a new object, causing React to miss
// re-renders and produce stale state.
export function cartReducer(state: CartState, action: CartAction): CartState {
  switch (action.type) {
    case 'ADD_ITEM': {
      const existing = state.items.find(i => i.id === action.payload.id);
      if (existing) {
        // BUG: Mutating state directly instead of creating new array
        state.items = state.items.map(i =>
          i.id === action.payload.id
            ? { ...i, quantity: i.quantity + action.payload.quantity }
            : i
        );
        return state; // BUG: Returning same reference
      }
      state.items = [...state.items, action.payload]; // BUG: Mutating state
      return state; // BUG: Returning same reference
    }

    case 'REMOVE_ITEM': {
      state.items = state.items.filter(i => i.id !== action.payload); // BUG: Mutating state
      return state; // BUG: Returning same reference
    }

    case 'APPLY_DISCOUNT': {
      state.discount = action.payload.percentage; // BUG: Mutating state
      state.discountCode = action.payload.code;   // BUG: Mutating state
      return state; // BUG: Returning same reference
    }

    case 'REMOVE_DISCOUNT': {
      state.discount = 0;        // BUG: Mutating state
      state.discountCode = null;  // BUG: Mutating state
      return state; // BUG: Returning same reference
    }

    default:
      return state;
  }
}

// BUG 2: Discount is subtracted twice — once here in calculateTotal,
// and it will be subtracted again in the CartDisplay component.
export function calculateTotal(state: CartState): number {
  const subtotal = state.items.reduce(
    (sum, item) => sum + item.price * item.quantity,
    0
  );

  // BUG: Applying discount here...
  const discountAmount = subtotal * (state.discount / 100);
  const total = subtotal - discountAmount;

  return total;
}

// Context
interface CartContextType {
  state: CartState;
  dispatch: React.Dispatch<CartAction>;
  total: number;
  subtotal: number;
}

const CartContext = createContext<CartContextType | undefined>(undefined);

export function CartProvider({ children }: { children: ReactNode }) {
  const [state, dispatch] = useReducer(cartReducer, initialState);

  const subtotal = state.items.reduce(
    (sum, item) => sum + item.price * item.quantity,
    0
  );

  const total = calculateTotal(state);

  return (
    <CartContext.Provider value={{ state, dispatch, total, subtotal }}>
      {children}
    </CartContext.Provider>
  );
}

export function useCart() {
  const context = useContext(CartContext);
  if (!context) {
    throw new Error('useCart must be used within a CartProvider');
  }
  return context;
}
