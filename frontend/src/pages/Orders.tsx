import { useState, useEffect } from "react";
import { useNavigate } from "react-router";
import { apiFetch } from "../services/api";
import { useAuth } from "../context/AuthContext";
import type { Order } from "../types";

function Orders() {
  const { token } = useAuth();
  const navigate = useNavigate();
  const [orders, setOrders] = useState<Order[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    if (!token) {
      navigate("/login");
      return;
    }

    const fetchOrders = async () => {
      try {
        const data = await apiFetch("/orders");
        setOrders(data);
      } catch (err) {
        setError("Failed to load orders");
      } finally {
        setLoading(false);
      }
    };

    fetchOrders();
  }, [token, navigate]);

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString("en-US", {
      year: "numeric",
      month: "long",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  if (loading) return <div className="loading">Loading orders...</div>;

  return (
    <div className="orders-container">
      <h1>Your Orders</h1>
      {error && <div className="error-message">{error}</div>}

      {orders.length === 0 ? (
        <div className="empty-orders">
          <p>You haven't placed any orders yet</p>
          <button onClick={() => navigate("/")}>Start Shopping</button>
        </div>
      ) : (
        <div className="orders-list">
          {orders.map((order) => (
            <div key={order.id} className="order-card">
              <div className="order-header">
                <h3>Order #{order.id}</h3>
                <span className={`order-status status-${order.status}`}>
                  {order.status}
                </span>
              </div>

              <div className="order-details">
                <p>Total: ${(order.total_cents / 100).toFixed(2)}</p>
                <p>Placed: {formatDate(order.created_at)}</p>
              </div>

              {order.shipping_name && (
                <div className="order-shipping">
                  <h4>Shipping To</h4>
                  <p>{order.shipping_name}</p>
                  <p>{order.shipping_address}</p>
                  <p>{order.shipping_phone}</p>
                </div>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

export default Orders;
