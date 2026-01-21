import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { useAppDispatch } from '../../app/hooks';
import { useRegisterMutation } from './authApi';
import { setCredentials } from './authSlice';
import { Input } from '../../components/ui/Input';
import styles from './Register.module.css';

export function Register() {
  const [formData, setFormData] = useState({
    email: '',
    password: '',
    firstName: '',
    lastName: '',
    phone: '',
  });
  const [register, { isLoading, error }] = useRegisterMutation();
  const dispatch = useAppDispatch();
  const navigate = useNavigate();

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    try {
      const result = await register(formData).unwrap();
      
      // Backend returns snake_case, convert to camelCase
      const credentials = {
        user: result.user,
        accessToken: (result as any).access_token || result.accessToken,
        refreshToken: (result as any).refresh_token || result.refreshToken,
      };
      
      console.log('✅ Registration successful, tokens:', {
        accessToken: credentials.accessToken?.substring(0, 20) + '...',
        refreshToken: credentials.refreshToken?.substring(0, 20) + '...',
      });
      
      dispatch(setCredentials(credentials));
      navigate('/products');
    } catch (err) {
      console.error('Registration failed:', err);
    }
  };

  return (
    <div className={styles.container}>
      <div className={styles.formWrapper}>
        <div className={styles.header}>
          <h2 className={styles.title}>
            Create your account
          </h2>
          <p className={styles.subtitle}>
            Or{' '}
            <Link to="/login" className={styles.link}>
              sign in to existing account
            </Link>
          </p>
        </div>
        
        <form className={styles.form} onSubmit={handleSubmit}>
          {error && (
            <div className={styles.errorBox}>
              <p className={styles.errorText}>
                {('data' in error && typeof error.data === 'object' && error.data !== null && 'error' in error.data)
                  ? String(error.data.error)
                  : 'Registration failed. Please try again.'}
              </p>
            </div>
          )}

          <div className={styles.inputGroup}>
            <div className={styles.nameGroup}>
              <Input
                label="First Name"
                name="firstName"
                value={formData.firstName}
                onChange={handleChange}
                required
              />
              <Input
                label="Last Name"
                name="lastName"
                value={formData.lastName}
                onChange={handleChange}
                required
              />
            </div>

            <Input
              label="Email address"
              type="email"
              name="email"
              value={formData.email}
              onChange={handleChange}
              required
              autoComplete="email"
            />

            <Input
              label="Phone (optional)"
              type="tel"
              name="phone"
              value={formData.phone}
              onChange={handleChange}
              autoComplete="tel"
            />

            <Input
              label="Password"
              type="password"
              name="password"
              value={formData.password}
              onChange={handleChange}
              required
              autoComplete="new-password"
            />
          </div>

          <div className={styles.buttonWrapper}>
            <button
              type="submit"
              disabled={isLoading}
              className={styles.button}
            >
              {isLoading ? 'Loading...' : 'Create account'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
