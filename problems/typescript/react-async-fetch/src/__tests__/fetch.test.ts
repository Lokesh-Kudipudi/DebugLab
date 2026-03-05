import { useFetchUsers, fetchUsers, User } from '../useFetchUsers';

describe('useFetchUsers Hook Logic', () => {
  test('fetchUsers returns all users when query is empty', async () => {
    const users = await fetchUsers('');
    expect(users.length).toBe(4);
    expect(users[0].name).toBe('Alice Johnson');
  });

  test('fetchUsers filters users by query correctly', async () => {
    const users = await fetchUsers('alice');
    expect(users.length).toBe(1);
    expect(users[0].name).toBe('Alice Johnson');
    expect(users[0].email).toBe('alice@example.com');
  });

  test('fetchUsers returns empty array for non-matching query', async () => {
    const users = await fetchUsers('zzzznonexistent');
    expect(users.length).toBe(0);
  });
});
