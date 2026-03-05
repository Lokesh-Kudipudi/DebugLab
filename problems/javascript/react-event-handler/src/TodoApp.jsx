import React, { useState } from 'react';

let nextId = 1;

export default function TodoApp() {
  const [todos, setTodos] = useState([]);
  const [input, setInput] = useState('');

  // BUG 1: Missing e.preventDefault() — form submission reloads the page
  const handleSubmit = (e) => {
    if (input.trim() === '') return;
    setTodos([...todos, { id: nextId++, text: input, done: false }]);
    setInput('');
  };

  // BUG 2: Using index instead of todo.id — toggles the wrong item
  const handleToggle = (index) => {
    setTodos(todos.map((todo, i) =>
      i === index ? { ...todo, done: !todo.done } : todo
    ));
  };

  // BUG 3: Using splice which mutates state directly, and wrong index logic
  const handleDelete = (index) => {
    const newTodos = todos;
    newTodos.splice(index, 1);
    setTodos(newTodos);
  };

  return (
    <div className="todo-app">
      <h1>Todo List</h1>
      <form onSubmit={handleSubmit}>
        <input
          type="text"
          value={input}
          onChange={e => setInput(e.target.value)}
          placeholder="Add a todo..."
          data-testid="todo-input"
        />
        <button type="submit" data-testid="add-button">Add</button>
      </form>
      <ul data-testid="todo-list">
        {todos.map((todo, index) => (
          <li key={todo.id} data-testid={`todo-${todo.id}`}>
            <span
              style={{ textDecoration: todo.done ? 'line-through' : 'none' }}
              onClick={() => handleToggle(index)}
              data-testid={`toggle-${todo.id}`}
            >
              {todo.text}
            </span>
            <button
              onClick={() => handleDelete(index)}
              data-testid={`delete-${todo.id}`}
            >
              Delete
            </button>
          </li>
        ))}
      </ul>
    </div>
  );
}
