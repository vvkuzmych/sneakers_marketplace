import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { useAppDispatch } from '../../app/hooks';
import { useLoginMutation } from './authApi';
import { setCredentials } from './authSlice';
import { Input } from '../../components/ui/Input';
import styles from './Login.module.css';

export function Login() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [login, { isLoading, error }] = useLoginMutation();
  const dispatch = useAppDispatch();
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    try {
      const result = await login({ email, password }).unwrap();
      
      console.log('🔍 Raw login response:', result);
      console.log('🔍 result keys:', Object.keys(result));
      
      // Backend returns snake_case, convert to camelCase
      const credentials = {
        user: result.user,
        accessToken: (result as any).access_token || (result as any).accessToken,
        refreshToken: (result as any).refresh_token || (result as any).refreshToken,
      };
      
      console.log('✅ Login successful, tokens:', {
        accessToken: credentials.accessToken?.substring(0, 30) + '...',
        refreshToken: credentials.refreshToken?.substring(0, 30) + '...',
      });
      
      if (!credentials.accessToken || !credentials.refreshToken) {
        console.error('❌ NO TOKENS IN RESPONSE!');
        alert('Error: No tokens received from server!');
        return;
      }
      
      dispatch(setCredentials(credentials));
      navigate('/products');
    } catch (err) {
      console.error('❌ Login failed:', err);
    }
  };

  return (
    <div className={styles.container}>
      <div className={styles.formWrapper}>
        <div className={styles.header}>
          <h2 className={styles.title}>
            Sign in to your account
          </h2>
          <p className={styles.subtitle}>
            Or{' '}
            <Link to="/register" className={styles.link}>
              create a new account
            </Link>
          </p>
        </div>
        
        <form className={styles.form} onSubmit={handleSubmit}>
          {error && (
            <div className={styles.errorBox}>
              <p className={styles.errorText}>
                {('data' in error && typeof error.data === 'object' && error.data !== null && 'error' in error.data)
                  ? String(error.data.error)
                  : 'Login failed. Please try again.'}
              </p>
            </div>
          )}

          <div className={styles.inputGroup}>
            <Input
              label="Email address"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
              autoComplete="email"
            />

            <Input
              label="Password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              autoComplete="current-password"
            />
          </div>

          <div className={styles.buttonWrapper}>
            <button
              type="submit"
              disabled={isLoading}
              className={styles.button}
            >
              {isLoading ? 'Loading...' : 'Sign in'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
