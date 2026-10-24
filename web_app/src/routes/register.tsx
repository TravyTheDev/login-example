import { apiFetch } from '#/apiFetch';
import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react';

export const Route = createFileRoute('/register')({
  component: RouteComponent,
})

type RegisterData = {
    username: string;
    email: string;
    password: string;
    passwordConfirm: string;
}

function RouteComponent() {
    const navigate = Route.useNavigate()
    const [registerData, setRegisterData] = useState<RegisterData>({
    username: "",
    email: "",
    password: "",
    passwordConfirm: "",
  });

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;

    setRegisterData((prev) => ({
      ...prev,
      [name]: value,
    }));
  };

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    if (registerData.password !== registerData.passwordConfirm) {
      alert("Passwords do not match!");
      return;
    }

    try {
        await apiFetch('/auth/register', {
            method: 'POST',
            body: JSON.stringify({ ...registerData }),
        })
        navigate({ to: '/login' })
    } catch (error) {
        throw new Error('registration failed')
    }
  };

  return (
    <div>
      <h1>Register New User</h1>

      <form onSubmit={handleSubmit}>
        <div>
          <label htmlFor="username">Username</label>
          <input
            id="username"
            name="username"
            type="text"
            value={registerData.username}
            onChange={handleChange}
            className='border border-gray-300 rounded-md'
          />
        </div>

        <div>
          <label htmlFor="email">Email</label>
          <input
            id="email"
            name="email"
            type="email"
            value={registerData.email}
            onChange={handleChange}
            className='border border-gray-300 rounded-md'
          />
        </div>

        <div>
          <label htmlFor="password">Password</label>
          <input
            id="password"
            name="password"
            type="password"
            value={registerData.password}
            onChange={handleChange}
            className='border border-gray-300 rounded-md'
          />
        </div>

        <div>
          <label htmlFor="passwordConfirm">Confirm Password</label>
          <input
            id="passwordConfirm"
            name="passwordConfirm"
            type="password"
            value={registerData.passwordConfirm}
            onChange={handleChange}
            className='border border-gray-300 rounded-md'
          />
        </div>

        <button className='px-4 py-2 border hover:cursor-pointer' type="submit">Register</button>
      </form>
    </div>
  );
}
