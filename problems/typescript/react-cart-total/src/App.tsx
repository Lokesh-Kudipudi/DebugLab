import React from 'react';
import { CartProvider } from './CartContext';
import { CartDisplay } from './CartDisplay';

export default function App() {
  return (
    <CartProvider>
      <div className="app">
        <h1>My Store</h1>
        <CartDisplay />
      </div>
    </CartProvider>
  );
}
