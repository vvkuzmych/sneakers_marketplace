import { Link, useNavigate } from 'react-router-dom';
import { useAppSelector, useAppDispatch } from '../../app/hooks';
import { logout } from '../../features/auth/authSlice';

export function Header() {
  const { isAuthenticated, user } = useAppSelector((state) => state.auth);
  const dispatch = useAppDispatch();
  const navigate = useNavigate();

  const handleLogout = () => {
    dispatch(logout());
    navigate('/login');
  };

  return (
    <header className="bg-white shadow-sm">
      <nav className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-16">
          {/* Logo */}
          <Link to="/" className="flex items-center">
            <span className="text-2xl font-bold text-blue-600">👟 Sneakers</span>
          </Link>

          {/* Navigation */}
          <div className="flex items-center space-x-4">
            <Link to="/products" className="text-gray-700 hover:text-blue-600 px-3 py-2 font-medium">
              Products
            </Link>
            
            {isAuthenticated ? (
              <>
                <Link 
                  to="/subscription/plans" 
                  className="flex items-center text-purple-600 hover:text-purple-700 px-3 py-2 font-medium"
                >
                  <span className="mr-1">💎</span>
                  Upgrade
                </Link>
                <Link to="/orders" className="text-gray-700 hover:text-blue-600 px-3 py-2">
                  Orders
                </Link>
                <Link to="/notifications" className="text-gray-700 hover:text-blue-600 px-3 py-2">
                  Notifications
                </Link>
                <Link to="/profile" className="text-gray-700 hover:text-blue-600 px-3 py-2">
                  {user?.firstName || 'Profile'}
                </Link>
                <button
                  onClick={handleLogout}
                  className="px-4 py-2 text-sm font-medium text-white bg-red-600 rounded-lg hover:bg-red-700 transition-colors"
                >
                  Logout
                </button>
              </>
            ) : (
              <>
                <Link to="/login">
                  <button className="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors">
                    Login
                  </button>
                </Link>
                <Link to="/register">
                  <button className="px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-lg hover:bg-blue-700 transition-colors">
                    Register
                  </button>
                </Link>
              </>
            )}
          </div>
        </div>
      </nav>
    </header>
  );
}
