import ReactDOM from 'react-dom/client'
import { RouterProvider } from '@tanstack/react-router'
import { AuthProvider, useAuth } from './auth'
import {router} from './router'
import './styles.css'
const rootElement = document.getElementById('app')!

function InnerApp() {
  const auth = useAuth()
  return <RouterProvider router={router} context={{ auth }} />
}

function App() {
  return (
    <AuthProvider>
      <InnerApp />
    </AuthProvider>
  )
}

if (!rootElement.innerHTML) {
  const root = ReactDOM.createRoot(rootElement)
  root.render(<App />)
}
