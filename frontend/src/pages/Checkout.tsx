import { useState, useEffect } from "react";
import { useNavigate } from "react-router";
import { apiFetch } from "../services/api";
import { useAuth } from "../context/AuthContext";
import { loadStripe } from "@stripe/stripe-js";
import {
  Elements,
  CardElement,
  useStripe,
  useElements,
} from "@stripe/react-stripe-js";

const stripePromise = loadStripe(import.meta.env.VITE_STRIPE_PUBLISHABLE_KEY);

function PaymentForm({
  orderId,
  shippingName,
  shippingAddress,
}: {
  orderId: number;
  shippingName: string;
  shippingAddress: string;
}) {
  const stripe = useStripe();
  const elements = useElements();
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();

    if (!stripe || !elements) return;

    setLoading(true);
    setError("");

    try {
      const { client_secret } = await apiFetch("/create-payment-intent", {
        method: "POST",
        body: JSON.stringify({ order_id: orderId }),
      });

      const { error: stripeError, paymentIntent } =
        await stripe.confirmCardPayment(client_secret, {
          payment_method: {
            card: elements.getElement(CardElement)!,
          },
        });

      if (stripeError) {
        setError(stripeError.message || "Payment failed");
      } else if (paymentIntent.status === "succeeded") {
        navigate("/orders");
      }
    } catch (err) {
      setError("Payment failed");
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="payment-form">
      <div className="shipping-summary">
        <p>
          <strong>Delivering to:</strong> {shippingName} — {shippingAddress}
        </p>
      </div>

      <div className="card-element-container">
        <CardElement
          options={{
            style: {
              base: {
                fontSize: "16px",
                color: "#424770",
                "::placeholder": { color: "#aab7c4" },
              },
              invalid: { color: "#9e2146" },
            },
          }}
        />
      </div>

      <button
        type="submit"
        disabled={!stripe || loading}
        className="checkout-btn"
      >
        {loading ? "Processing..." : "Pay Now"}
      </button>

      {error && <div className="error-message">{error}</div>}
    </form>
  );
}

function Checkout() {
  const { token } = useAuth();
  const navigate = useNavigate();

  const [shippingName, setShippingName] = useState("");
  const [shippingAddress, setShippingAddress] = useState("");
  const [shippingPhone, setShippingPhone] = useState("");

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [orderId, setOrderId] = useState<number | null>(null);

  useEffect(() => {
    if (!token) {
      navigate("/login");
    }
  }, [token, navigate]);

  const createOrder = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError("");

    try {
      const order = await apiFetch("/orders", {
        method: "POST",
        body: JSON.stringify({
          shipping_name: shippingName,
          shipping_address: shippingAddress,
          shipping_phone: shippingPhone,
        }),
      });
      setOrderId(order.order_id);
    } catch (err) {
      setError("Failed to create order. Is your cart empty?");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="checkout-container">
      <h1>Checkout</h1>
      {error && <div className="error-message">{error}</div>}

      {!orderId ? (
        <form onSubmit={createOrder} className="shipping-form">
          <h2>Shipping Information</h2>

          <div className="form-group">
            <label htmlFor="shippingName">Full Name</label>
            <input
              id="shippingName"
              type="text"
              value={shippingName}
              onChange={(e) => setShippingName(e.target.value)}
              required
            />
          </div>

          <div className="form-group">
            <label htmlFor="shippingAddress">Delivery Address</label>
            <input
              id="shippingAddress"
              type="text"
              value={shippingAddress}
              onChange={(e) => setShippingAddress(e.target.value)}
              required
            />
          </div>

          <div className="form-group">
            <label htmlFor="shippingPhone">Phone Number</label>
            <input
              id="shippingPhone"
              type="tel"
              value={shippingPhone}
              onChange={(e) => setShippingPhone(e.target.value)}
              required
            />
          </div>

          <button type="submit" disabled={loading} className="checkout-btn">
            {loading ? "Creating order..." : "Continue to Payment"}
          </button>
        </form>
      ) : (
        <Elements stripe={stripePromise}>
          <PaymentForm
            orderId={orderId}
            shippingName={shippingName}
            shippingAddress={shippingAddress}
          />
        </Elements>
      )}
    </div>
  );
}

export default Checkout;
