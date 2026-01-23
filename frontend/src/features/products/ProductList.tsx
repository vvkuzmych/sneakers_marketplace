import { Link } from 'react-router-dom';
import { useGetProductsQuery } from './productsApi';
import { Typography } from '../../components/ui/Typography';
import { Box } from '../../components/ui/Box';
import { Card } from '../../components/ui/Card';
import { Badge } from '../../components/ui/Badge';
import CurrentSubscriptionWidget from '../subscription/CurrentSubscriptionWidget';
import { useAppSelector } from '../../app/hooks';

export function ProductList() {
  const { data, isLoading, error } = useGetProductsQuery({ page: 1, pageSize: 12 });
  const { isAuthenticated } = useAppSelector((state) => state.auth);

  if (isLoading) {
    return (
      <Box className="min-h-screen" flex alignItems="center" justifyContent="center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
      </Box>
    );
  }

  if (error) {
    return (
      <Box className="max-w-7xl mx-auto px-4" p={8}>
        <Box className="bg-red-50 rounded-lg" p={4}>
          <Typography variant="body" color="error">
            Failed to load products. Please try again later.
          </Typography>
        </Box>
      </Box>
    );
  }

  return (
    <Box className="max-w-7xl mx-auto px-4" p={8}>
      <Typography variant="h1" className="mb-8">Sneakers Catalog</Typography>
      
      {/* Subscription Widget - Only for authenticated users */}
      {isAuthenticated && (
        <Box className="mb-8">
          <CurrentSubscriptionWidget />
        </Box>
      )}
      
      <Box className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4" gap={6}>
        {data?.products.map((product) => (
          <Link
            key={product.id}
            to={`/bidding/${product.id}`}
            className="block"
          >
            <Card padding="none" hover>
              <Box className="aspect-square bg-gray-200" flex alignItems="center" justifyContent="center">
                <span className="text-6xl">👟</span>
              </Box>
              <Box p={4}>
                <Typography variant="h4" className="truncate">
                  {product.name}
                </Typography>
                <Typography variant="caption" color="secondary" className="mt-1">
                  {product.brand} - {product.model}
                </Typography>
                <Box className="mt-4" flex alignItems="center" justifyContent="between">
                  <Typography variant="h3" className="text-green-600">
                    ${product.retailPrice}
                  </Typography>
                  <Badge variant="secondary" size="sm">
                    {product.category}
                  </Badge>
                </Box>
                <button
                  className="mt-4 w-full px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors font-semibold text-sm"
                  onClick={(e) => {
                    e.preventDefault();
                    window.location.href = `/bidding/${product.id}`;
                  }}
                >
                  📊 View Bidding
                </button>
              </Box>
            </Card>
          </Link>
        ))}
      </Box>

      {data && data.products.length === 0 && (
        <Box className="text-center" p={12}>
          <Typography variant="body" color="secondary">
            No products found.
          </Typography>
        </Box>
      )}
    </Box>
  );
}
