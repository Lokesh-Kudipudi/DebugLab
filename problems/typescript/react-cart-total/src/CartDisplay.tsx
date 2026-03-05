import React from 'react';
import { useCart, DISCOUNT_CODES } from './CartContext';

export function CartDisplay() {
  const { state, dispatch, total, subtotal } = useCart();
  const [codeInput, setCodeInput] = React.useState('');
  const [error, setError] = React.useState('');

  const handleApplyDiscount = () => {
    const percentage = DISCOUNT_CODES[codeInput.toUpperCase()];
    if (percentage !== undefined) {
      dispatch({
        type: 'APPLY_DISCOUNT',
        payload: { code: codeInput.toUpperCase(), percentage },
      });
      setError('');
    } else {
      setError('Invalid discount code');
    }
  };

  const handleRemoveDiscount = () => {
    dispatch({ type: 'REMOVE_DISCOUNT' });
    setCodeInput('');
  };

  // BUG: Discount is subtracted AGAIN here — it was already subtracted
  // in calculateTotal() inside CartContext.tsx
  const discountAmount = subtotal * (state.discount / 100);
  const displayTotal = total - discountAmount;

  return (
    <div className="cart">
      <h2>Shopping Cart</h2>

      {state.items.length === 0 ? (
        <p>Your cart is empty</p>
      ) : (
        <>
          <ul>
            {state.items.map(item => (
              <li key={item.id}>
                {item.name} × {item.quantity} — ${(item.price * item.quantity).toFixed(2)}
                <button onClick={() => dispatch({ type: 'REMOVE_ITEM', payload: item.id })}>
                  Remove
                </button>
              </li>
            ))}
          </ul>

          <div className="cart-summary">
            <p>Subtotal: ${subtotal.toFixed(2)}</p>

            {state.discountCode && (
              <p>
                Discount ({state.discountCode} — {state.discount}% off): -${discountAmount.toFixed(2)}
                <button onClick={handleRemoveDiscount}>Remove</button>
              </p>
            )}

            <p className="total" data-testid="cart-total">
              Total: ${displayTotal.toFixed(2)}
            </p>

            <button
              className="checkout-btn"
              data-testid="checkout-button"
              disabled={displayTotal <= 0}
            >
              Checkout — ${displayTotal.toFixed(2)}
            </button>
          </div>

          {!state.discountCode && (
            <div className="discount-form">
              <input
                type="text"
                value={codeInput}
                onChange={e => setCodeInput(e.target.value)}
                placeholder="Discount code"
                data-testid="discount-input"
              />
              <button onClick={handleApplyDiscount} data-testid="apply-discount">
                Apply
              </button>
              {error && <p className="error">{error}</p>}
            </div>
          )}
        </>
      )}
    </div>
  );
}
