import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { Provider } from 'react-redux';
import { store } from './app/store';
import { Header } from './components/layout/Header';
import { Login } from './features/auth/Login';
import { Register } from './features/auth/Register';
import { ProtectedRoute } from './features/auth/ProtectedRoute';
import { ProductList } from './features/products/ProductList';
import BiddingPage from './features/bidding/BiddingPage';

function App() {
  return (
    <Provider store={store}>
      <BrowserRouter>
        <div className="min-h-screen bg-gray-50">
          <Header />
          <main>
            <Routes>
              {/* Public routes */}
              <Route path="/" element={<Navigate to="/products" replace />} />
              <Route path="/login" element={<Login />} />
              <Route path="/register" element={<Register />} />
              <Route path="/products" element={<ProductList />} />
              <Route path="/bidding" element={<Navigate to="/products" replace />} />
              <Route
                path="/bidding/:productId"
                element={
                  <ProtectedRoute>
                    <BiddingPage />
                  </ProtectedRoute>
                }
              />
              
              {/* Protected routes */}
              <Route
                path="/orders"
                element={
                  <ProtectedRoute>
                    <div className="max-w-7xl mx-auto px-4 py-8">
                      <h1 className="text-3xl font-bold">My Orders</h1>
                      <p className="mt-4 text-gray-600">Orders page coming soon...</p>
                    </div>
                  </ProtectedRoute>
                }
              />
              <Route
                path="/notifications"
                element={
                  <ProtectedRoute>
                    <div className="max-w-7xl mx-auto px-4 py-8">
                      <h1 className="text-3xl font-bold">Notifications</h1>
                      <p className="mt-4 text-gray-600">Notifications page coming soon...</p>
                    </div>
                  </ProtectedRoute>
                }
              />
              <Route
                path="/profile"
                element={
                  <ProtectedRoute>
                    <div className="max-w-7xl mx-auto px-4 py-8">
                      <h1 className="text-3xl font-bold">Profile</h1>
                      <p className="mt-4 text-gray-600">Profile page coming soon...</p>
                    </div>
                  </ProtectedRoute>
                }
              />
              
              {/* 404 */}
              <Route path="*" element={<div className="max-w-7xl mx-auto px-4 py-8 text-center"><h1 className="text-2xl font-bold">404 - Page Not Found</h1></div>} />
            </Routes>
          </main>
        </div>
      </BrowserRouter>
    </Provider>
  );
}

export default App;
