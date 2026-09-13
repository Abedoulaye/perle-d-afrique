import { useState, useEffect } from "react";
import { Link } from "react-router";
import { getProducts } from "../services/productService";
import type { Product } from "../types";

function ProductList() {
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    const fetchProducts = async () => {
      try {
        const data = await getProducts();
        setProducts(data);
      } catch (err) {
        console.error(err);
        setError("Failed to load products");
      } finally {
        setLoading(false);
      }
    };
    fetchProducts();
  }, []);

  if (loading) return <div className="loading">Loading products...</div>;
  if (error) return <div className="error">{error}</div>;
  const featuredProducts = products.slice(0, 6);
  return (
    <>
      <div className="products-header">
        <div className="products-title">
          <h2>Featured Products</h2>
          <span className="title-dash"></span>
        </div>

        <Link to="/products" className="view-all">
          View All <span>→</span>
        </Link>
      </div>
      <div className="product-grid">
        {featuredProducts.map((product) => (
          <Link
            to={`/products/${product.id}`}
            key={product.id}
            className="product-card"
          >
            <img
              src={product.image}
              alt={product.name}
              className="product-image"
            />
            <div className="product-info">
              <h3>{product.name}</h3>
              <p className="product-price">
                ${(product.price_cents / 100).toFixed(2)}
              </p>
            </div>
          </Link>
        ))}
      </div>
    </>
  );
}

export default ProductList;
