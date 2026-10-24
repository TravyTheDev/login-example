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
    const navigate = Route.useNavigate()
    const [users, setUsers] = useState<User[]>()

    const loadAllUsers = async () => {
        try {
            const res = await apiFetch('/users/all_users')
            const users = await res.json()
            setUsers(users)
        } catch (error) {
            console.log(error)
        }
    }

    useEffect(() => {
        loadAllUsers()
    }, [])

    const loginAsUser = async (user: User) => {
        try {
            auth.loginAs(user.email)
            navigate({ to: '/' })
        } catch (error) {
            console.log(error)
        }
    }

  return (<div>
    {users?.map((user) => (
        <button key={user.email} onClick={() => loginAsUser(user)}>{user.username}</button>
    ))}
  </div>)
}
