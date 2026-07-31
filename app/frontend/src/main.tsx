import ReactDOM from 'react-dom/client'
import { createRouter } from '@tanstack/react-router'
import { routeTree } from './routeTree.gen'
import "./global.css"
import App from './App';
import ErrorBoundary from './components/error-boundary';

const router = createRouter({
  routeTree,
  defaultPreload: 'intent',
  scrollRestoration: true,
})

declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router
  }
}

const rootElement = document.getElementById('app')!


if (!rootElement.innerHTML) {
  const root = ReactDOM.createRoot(rootElement)
  root.render(
    <>
      <ErrorBoundary fallback={<p>Something went wrong</p>}>
        <App router={router} />
      </ErrorBoundary>
    </>
  )
}
