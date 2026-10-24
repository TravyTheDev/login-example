import { apiFetch } from '#/apiFetch'
import { createFileRoute } from '@tanstack/react-router'
import { useEffect, useState } from 'react'

export const Route = createFileRoute('/_authenticated/allUsers')({
  component: RouteComponent,
})

type User = {
    username: string;
    email: string;
}

function RouteComponent() {
    const { auth } = Route.useRouteContext()
    const [users, setUsers] = useState<User[]>()

    const loadAllUsers = async () => {
        try {
            const res = await apiFetch('/users/all_users')
            const users = await res.json()
            setUsers(users)
        } catch (error) {
            throw new Error('failed to load users')
        }
    }

    useEffect(() => {
        loadAllUsers()
    }, [])

    const loginAsUser = async (user: User) => {
        try {
            await auth.loginAs(user.email)
            window.location.href = '/'
        } catch (error) {
            throw new Error('failed login as')
        }
    }

  return (<div className='flex flex-col gap-1 items-start'>
    {users?.map((user) => (
        <button className='hover:cursor-pointer' key={user.email} onClick={() => loginAsUser(user)}>{user.username}</button>
    ))}
  </div>)
}
