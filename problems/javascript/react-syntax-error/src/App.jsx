import React, { useState } from 'react';
import { ProfileCard } from './ProfileCard';

export default function App() {
  const [users, setUsers] = useState([
    {
      name: 'Alice Johnson',
      email: 'alice@example.com',
      bio: 'Full-stack developer who loves React',
      avatar: 'https://via.placeholder.com/100',
    },
    {
      name: 'Bob Smith',
      email: 'bob@example.com',
      bio: 'DevOps engineer and cloud enthusiast',
      avatar: 'https://via.placeholder.com/100',
    },
  ]);

  // BUG: Incorrect arrow function syntax — `=` instead of `=>`
  const handleSave = (index) => (updated) = {
    setUsers(prev => {
      const copy = [...prev];
      copy[index] = updated;
      return copy;
    });
  };

  return (
    <div className="app">
      <h1>User Profiles</h1>
      <div className="profiles-grid">
        {users.map((user, index) => (
          <ProfileCard
            key={user.email}
            user={user}
            onSave={handleSave(index)}
          />
        ))}
      </div>
    </div>
  );
}
