import { useAuth } from '#/auth'
import { createFileRoute, Link } from '@tanstack/react-router'

export const Route = createFileRoute('/_authenticated/mypage')({
  component: RouteComponent,
})

function RouteComponent() {
  const auth = useAuth()
  const navigate = Route.useNavigate()
  const logout = () => {
    try {
      if (auth.user?.impersonatedBy) {
        auth.logoutFromLoginAs()
        navigate({ to: '/' })
      } else {
        auth.logout()
        navigate({ to: '/' })
      }
    } catch (error) {
      throw new Error('failed to logout')
    }
  }
  return <div className="p-6">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold">My Page</h1>
        <button
          onClick={logout}
          className="bg-red-600 text-white px-4 py-2 rounded hover:bg-red-700"
        >
          Sign Out
        </button>
      </div>

      <div className="bg-white p-6 rounded-lg shadow">
        <h2 className="text-xl font-semibold mb-2">Welcome back!</h2>
        <p className="text-gray-600">
          Hello, <strong>{auth.user?.username}</strong>! You are successfully
          authenticated.
        </p>
        <p className="text-sm text-gray-500 mt-2">Email: {auth.user?.email}</p>
      </div>
      <div>
        <Link to='/'>Home</Link>
      </div>
      <Link to='/allUsers'>Login in as other users</Link>
    </div>
}
