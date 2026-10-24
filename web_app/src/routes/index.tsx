import { useAuth } from '#/auth'
import { createFileRoute, Link } from '@tanstack/react-router'

export const Route = createFileRoute('/')({ component: Home })

function Home() {
  const auth  = useAuth()
  
  return (
    <div className="p-8">
      <h1 className="text-4xl font-bold">Welcome</h1>
      {auth.isAuthenticated ? 
        <div className="whitespace-pre-wrap">
          <span>Hello </span>
          <Link to='/mypage'>{auth.user?.username}</Link>
        </div>
        :
      <div className="flex flex-col gap-1 items-start font-semibold">
        <Link to='/login' search={{redirect: '/login'}}>Login</Link>
        <Link to='/register' search={{redirect: '/register'}}>Register</Link>
      </div>
      }
    </div>
  )
}
