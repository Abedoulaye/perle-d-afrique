import { useState, useEffect } from "react";
import { useNavigate } from "react-router";
import { apiFetch } from "../services/api";
import { useAuth } from "../context/AuthContext";
import type { CartDetail } from "../types";

function Cart() {
  const { token } = useAuth();
  const navigate = useNavigate();
  const [items, setItems] = useState<CartDetail[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  const fetchCart = async () => {
    try {
      const data = await apiFetch("/cart");
      setItems(data);
    } catch (err) {
      setError("Failed to load cart");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (!token) {
      navigate("/login");
      return;
    }
    fetchCart();
  }, [token, navigate]);

  const updateQuantity = async (productId: number, newQuantity: number) => {
    if (newQuantity <= 0) {
      await removeItem(productId);
      return;
    }
    try {
      await apiFetch("/cart", {
        method: "PUT",
        body: JSON.stringify({ product_id: productId, quantity: newQuantity }),
      });
      fetchCart();
    } catch (err) {
      setError("Failed to update quantity");
    }
  };

  const removeItem = async (productId: number) => {
    try {
      await apiFetch(`/cart/${productId}`, {
        method: "DELETE",
      });
      fetchCart();
    } catch (err) {
      setError("Failed to remove item");
    }
  };

  const clearCart = async () => {
    try {
      await apiFetch("/cart", {
        method: "DELETE",
      });
      fetchCart();
    } catch (err) {
      setError("Failed to clear cart");
    }
  };

  const subtotal = items.reduce((sum, item) => {
    return sum + (item.price_cents || 0) * item.quantity;
  }, 0);

  if (loading) return <div className="loading">Loading cart...</div>;

  return (
    <div className="cart-container">
      <h1>Your Cart</h1>
      {error && <div className="error-message">{error}</div>}

      {items.length === 0 ? (
        <div className="empty-cart">
          <p>Your cart is empty</p>
          <button onClick={() => navigate("/")}>Continue Shopping</button>
        </div>
      ) : (
        <>
          <div className="cart-items">
            {items.map((item) => (
              <div key={item.id} className="cart-item">
                <img
                  src={item.image}
                  alt={item.name}
                  className="cart-item-image"
                />
                <div className="cart-item-info">
                  <h3>{item.name}</h3>
                  <p>${((item.price_cents || 0) / 100).toFixed(2)}</p>
                </div>
                <div className="cart-item-quantity">
                  <button
                    onClick={() =>
                      updateQuantity(item.product_id, item.quantity - 1)
                    }
                  >
                    -
                  </button>
                  <span>{item.quantity}</span>
                  <button
                    onClick={() =>
                      updateQuantity(item.product_id, item.quantity + 1)
                    }
                  >
                    +
                  </button>
                </div>
                <div className="cart-item-total">
                  $
                  {(((item.price_cents || 0) * item.quantity) / 100).toFixed(2)}
                </div>
                <button
                  className="remove-item"
                  onClick={() => removeItem(item.product_id)}
                >
                  Remove
                </button>
              </div>
            ))}
          </div>

          <div className="cart-summary">
            <h2>Subtotal: ${(subtotal / 100).toFixed(2)}</h2>
            <div className="cart-actions">
              <button onClick={clearCart} className="clear-cart">
                Clear Cart
              </button>
              <button
                onClick={() => navigate("/checkout")}
                className="checkout-btn"
              >
                Proceed to Checkout
              </button>
            </div>
          </div>
        </>
      )}
    </div>
  );
}

export default Cart;
