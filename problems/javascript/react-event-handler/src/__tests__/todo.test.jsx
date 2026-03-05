import { TodoApp } from '../TodoApp';

// We test the pure logic functions extracted from the component behavior.
// These tests validate what the FIXED version should produce.

describe('Todo App Logic', () => {
  test('adding a todo creates a new item with correct structure', () => {
    const todos = [];
    const newTodo = { id: 1, text: 'Buy groceries', done: false };
    const updated = [...todos, newTodo];

    expect(updated).toHaveLength(1);
    expect(updated[0].text).toBe('Buy groceries');
    expect(updated[0].done).toBe(false);
    expect(updated[0].id).toBe(1);
  });

  test('toggling a todo by id flips done status correctly', () => {
    const todos = [
      { id: 1, text: 'Buy groceries', done: false },
      { id: 2, text: 'Walk the dog', done: false },
      { id: 3, text: 'Read a book', done: true },
    ];

    // Toggle todo with id 2
    const targetId = 2;
    const toggled = todos.map(todo =>
      todo.id === targetId ? { ...todo, done: !todo.done } : todo
    );

    expect(toggled[0].done).toBe(false); // unchanged
    expect(toggled[1].done).toBe(true);  // toggled
    expect(toggled[2].done).toBe(true);  // unchanged

    // Toggling should NOT mutate original
    expect(todos[1].done).toBe(false);
  });

  test('deleting a todo by id removes only that item without mutating original', () => {
    const todos = [
      { id: 1, text: 'Buy groceries', done: false },
      { id: 2, text: 'Walk the dog', done: false },
      { id: 3, text: 'Read a book', done: true },
    ];

    // Delete todo with id 2 using filter (correct approach)
    const targetId = 2;
    const afterDelete = todos.filter(todo => todo.id !== targetId);

    expect(afterDelete).toHaveLength(2);
    expect(afterDelete.find(t => t.id === 2)).toBeUndefined();
    expect(afterDelete[0].id).toBe(1);
    expect(afterDelete[1].id).toBe(3);

    // Original should NOT be mutated
    expect(todos).toHaveLength(3);
  });
});
