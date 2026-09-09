import { useState, useEffect } from "react";
import { useParams, useNavigate } from "react-router";
import { getProduct } from "../services/productService";
import { apiFetch } from "../services/api";
import { useAuth } from "../context/AuthContext";
import type { Product } from "../types";

function ProductDetail() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { token } = useAuth();
  const [product, setProduct] = useState<Product | null>(null);
  const [quantity, setQuantity] = useState(1);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [adding, setAdding] = useState(false);

  useEffect(() => {
    const fetchProduct = async () => {
      try {
        const data = await getProduct(Number(id));
        setProduct(data);
      } catch (err) {
        setError("Failed to load product");
      } finally {
        setLoading(false);
      }
    };
    fetchProduct();
  }, [id]);

  const addToCart = async () => {
    if (!token) {
      navigate("/login");
      return;
    }

    setAdding(true);
    try {
      await apiFetch("/cart", {
        method: "POST",
        body: JSON.stringify({
          product_id: Number(id),
          quantity,
        }),
      });
      alert("Added to cart!");
    } catch (err) {
      setError("Failed to add to cart");
    } finally {
      setAdding(false);
    }
  };

  if (loading) return <div className="loading">Loading product...</div>;
  if (error) return <div className="error">{error}</div>;
  if (!product) return <div className="error">Product not found</div>;

  return (
    <div className="product-detail">
      <img
        src={product.image}
        alt={product.name}
        className="product-detail-image"
      />
      <div className="product-detail-info">
        <h1>{product.name}</h1>
        <p className="product-price">
          ${(product.price_cents / 100).toFixed(2)}
        </p>
        <p className="product-description">{product.description}</p>
        <p
          className={`stock ${product.stock > 0 ? "in-stock" : "out-of-stock"}`}
        >
          {product.stock > 0 ? `${product.stock} in stock` : "Out of stock"}
        </p>

        {product.stock > 0 && (
          <div className="add-to-cart">
            <label>
              Quantity:
              <input
                type="number"
                min="1"
                max={product.stock}
                value={quantity}
                onChange={(e) => setQuantity(Number(e.target.value))}
              />
            </label>
            <button onClick={addToCart} disabled={adding}>
              {adding ? "Adding..." : "Add to Cart"}
            </button>
          </div>
        )}
      </div>
    </div>
  );
}

export default ProductDetail;
