import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { RouterProvider } from 'react-router-dom'
import './index.css'
import router from './router'
import { ReactQueryProvider } from './features/layouts/ReactQuery'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <ReactQueryProvider>
    <RouterProvider router={router} />
    </ReactQueryProvider>
  </StrictMode>
)
