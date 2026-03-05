import { cartReducer, calculateTotal, CartState, DISCOUNT_CODES } from '../CartContext';

// Helper to create a pre-populated cart state
function createCartWithItems(): CartState {
  return {
    items: [
      { id: '1', name: 'Widget', price: 25.0, quantity: 2 },  // $50
      { id: '2', name: 'Gadget', price: 50.0, quantity: 1 },  // $50
    ],
    discount: 0,
    discountCode: null,
  };
}

describe('Shopping Cart Total', () => {
  test('total is correct after applying a valid discount code', () => {
    // Arrange: cart with $100 subtotal
    let state = createCartWithItems();

    // Act: apply 10% discount
    state = cartReducer(state, {
      type: 'APPLY_DISCOUNT',
      payload: { code: 'SAVE10', percentage: DISCOUNT_CODES['SAVE10'] },
    });

    const total = calculateTotal(state);

    // Assert: $100 - 10% = $90
    // If the discount is applied twice (bug), total would be $80 or less
    expect(total).toBe(90);
  });

  test('total is correct after removing the discount code', () => {
    // Arrange: cart with discount applied
    let state = createCartWithItems();

    // Apply discount first
    state = cartReducer(state, {
      type: 'APPLY_DISCOUNT',
      payload: { code: 'SAVE20', percentage: DISCOUNT_CODES['SAVE20'] },
    });

    // Act: remove the discount
    state = cartReducer(state, { type: 'REMOVE_DISCOUNT' });

    const total = calculateTotal(state);

    // Assert: should be back to original $100
    expect(total).toBe(100);
    expect(state.discount).toBe(0);
    expect(state.discountCode).toBeNull();
  });

  test('total never goes below zero regardless of discount value', () => {
    // Arrange: cart with small subtotal
    const state: CartState = {
      items: [
        { id: '1', name: 'Small Item', price: 5.0, quantity: 1 }, // $5
      ],
      discount: 0,
      discountCode: null,
    };

    // Act: apply 100% discount (FREE100 code)
    const newState = cartReducer(state, {
      type: 'APPLY_DISCOUNT',
      payload: { code: 'FREE100', percentage: DISCOUNT_CODES['FREE100'] },
    });

    const total = calculateTotal(newState);

    // Assert: total should be exactly $0, never negative
    // If discount is subtracted twice, total would be -$5
    expect(total).toBeGreaterThanOrEqual(0);
    expect(total).toBe(0);
  });
});
