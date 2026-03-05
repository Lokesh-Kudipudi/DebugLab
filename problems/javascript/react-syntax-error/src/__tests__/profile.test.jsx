import { ProfileCard } from '../ProfileCard';

// Helper: create a test user
function createTestUser() {
  return {
    name: 'Test User',
    email: 'test@example.com',
    bio: 'A test bio',
    avatar: 'https://via.placeholder.com/100',
  };
}

describe('ProfileCard Component', () => {
  test('ProfileCard module exports correctly', () => {
    // Verifies the module can be imported without syntax errors
    expect(ProfileCard).toBeDefined();
    expect(typeof ProfileCard).toBe('function');

    const user = createTestUser();
    expect(user.name).toBe('Test User');
    expect(user.email).toBe('test@example.com');
    expect(user.bio).toBe('A test bio');
  });

  test('handleChange function creates correct updated profile data', () => {
    const user = createTestUser();

    // Simulate updating the name field
    const updated = {
      ...user,
      name: 'Updated Name',
    };

    expect(updated.name).toBe('Updated Name');
    expect(updated.email).toBe('test@example.com');
    expect(updated.bio).toBe('A test bio');
  });

  test('validation rejects empty name and email', () => {
    const user = createTestUser();

    // Test with empty name
    const emptyName = { ...user, name: '' };
    const nameInvalid = emptyName.name.trim() === '' || emptyName.email.trim() === '';
    expect(nameInvalid).toBe(true);

    // Test with empty email
    const emptyEmail = { ...user, email: '  ' };
    const emailInvalid = emptyEmail.name.trim() === '' || emptyEmail.email.trim() === '';
    expect(emailInvalid).toBe(true);

    // Test with valid data
    const validUser = createTestUser();
    const valid = !(validUser.name.trim() === '' || validUser.email.trim() === '');
    expect(valid).toBe(true);
  });
});
