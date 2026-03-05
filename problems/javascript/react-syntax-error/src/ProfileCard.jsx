import React, { useState } from 'react';

// BUG 1: Arrow function uses `=` instead of `=>`
// BUG 2: Missing closing curly brace for the if-block
// BUG 3: Missing closing angle bracket on a JSX div tag
export function ProfileCard({ user, onSave }) {
  const [isEditing, setIsEditing] = useState(false);
  const [formData, setFormData] = useState(user);

  // BUG: Arrow function uses `=` instead of `=>`
  const handleChange = (field, value) = {
    setFormData(prev => ({
      ...prev,
      [field]: value,
    }));
  };

  const handleSave = () => {
    // BUG: Missing closing brace for the if-block
    if (formData.name.trim() === '' || formData.email.trim() === '') {
      alert('Name and email are required!');

    onSave(formData);
    setIsEditing(false);
  };

  // BUG: Missing closing angle bracket on a JSX div tag
  return (
    <div className="profile-card">
      <div className="profile-header"
        <img src={user.avatar} alt={user.name} className="avatar" />
        <h2>{user.name}</h2>
      </div>

      {isEditing ? (
        <div className="profile-form">
          <label>
            Name:
            <input
              type="text"
              value={formData.name}
              onChange={e => handleChange('name', e.target.value)}
              data-testid="name-input"
            />
          </label>
          <label>
            Email:
            <input
              type="email"
              value={formData.email}
              onChange={e => handleChange('email', e.target.value)}
              data-testid="email-input"
            />
          </label>
          <label>
            Bio:
            <textarea
              value={formData.bio}
              onChange={e => handleChange('bio', e.target.value)}
              data-testid="bio-input"
            />
          </label>
          <button onClick={handleSave} data-testid="save-button">Save</button>
          <button onClick={() => setIsEditing(false)}>Cancel</button>
        </div>
      ) : (
        <div className="profile-info">
          <p data-testid="user-email">{user.email}</p>
          <p data-testid="user-bio">{user.bio}</p>
          <button onClick={() => setIsEditing(true)} data-testid="edit-button">
            Edit Profile
          </button>
        </div>
      )}
    </div>
  );
}
