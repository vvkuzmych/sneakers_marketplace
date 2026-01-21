import { Link } from 'react-router-dom';
import { useGetProductsQuery } from './productsApi';

export function ProductList() {
  const { data, isLoading, error } = useGetProductsQuery({ page: 1, pageSize: 12 });

  if (isLoading) {
    return (
      <div className="flex justify-center items-center min-h-screen">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="max-w-7xl mx-auto px-4 py-8">
        <div className="bg-red-50 p-4 rounded-lg">
          <p className="text-red-800">Failed to load products. Please try again later.</p>
        </div>
      </div>
    );
  }

  return (
    <div className="max-w-7xl mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold text-gray-900 mb-8">Sneakers Catalog</h1>
      
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
        {data?.products.map((product) => (
          <Link
            key={product.id}
            to={`/products/${product.id}`}
            className="bg-white rounded-lg shadow-md overflow-hidden hover:shadow-xl transition-shadow"
          >
            <div className="aspect-square bg-gray-200 flex items-center justify-center">
              <span className="text-6xl">👟</span>
            </div>
            <div className="p-4">
              <h3 className="text-lg font-semibold text-gray-900 truncate">
                {product.name}
              </h3>
              <p className="text-sm text-gray-500 mt-1">
                {product.brand} - {product.model}
              </p>
              <div className="mt-4 flex items-center justify-between">
                <span className="text-xl font-bold text-green-600">
                  ${product.retailPrice}
                </span>
                <span className="text-sm text-gray-500">{product.category}</span>
              </div>
            </div>
          </Link>
        ))}
      </div>

      {data && data.products.length === 0 && (
        <div className="text-center py-12">
          <p className="text-gray-500">No products found.</p>
        </div>
      )}
    </div>
  );
}
