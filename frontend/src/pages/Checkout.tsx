import { useState, useEffect } from "react";
import { useNavigate } from "react-router";
import { apiFetch } from "../services/api";
import { useAuth } from "../context/AuthContext";
import { loadStripe } from "@stripe/stripe-js";

const stripePromise = loadStripe(import.meta.env.VITE_STRIPE_PUBLISHABLE_KEY);

function Checkout() {
  const { token } = useAuth();
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [orderId, setOrderId] = useState<number | null>(null);

  const createOrder = async () => {
    setLoading(true);
    setError("");
    try {
      const order = await apiFetch("/orders", {
        method: "POST",
      });
      setOrderId(order.order_id);
      return order;
    } catch (err) {
      setError("Failed to create order. Is your cart empty?");
      throw err;
    } finally {
      setLoading(false);
    }
  };

  const handlePayment = async () => {
    if (!orderId) {
      setError("Please create an order first");
      return;
    }

    setLoading(true);
    setError("");

    try {
      // Create PaymentIntent
      const { client_secret } = await apiFetch("/create-payment-intent", {
        method: "POST",
        body: JSON.stringify({ order_id: orderId }),
      });

      // Load Stripe
      const stripe = await stripePromise;
      if (!stripe) {
        throw new Error("Stripe failed to load");
      }

      // Confirm payment with test card
      const { error: stripeError } = await stripe.confirmCardPayment(
        client_secret,
        {
          payment_method: {
            card: {
              number: "4242424242424242",
              exp_month: 12,
              exp_year: 2030,
              cvc: "123",
            } as any,
          },
        },
      );

      if (stripeError) {
        setError(stripeError.message || "Payment failed");
      } else {
        alert("Payment successful!");
        navigate("/orders");
      }
    } catch (err) {
      setError("Payment failed");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="checkout-container">
      <h1>Checkout</h1>
      {error && <div className="error-message">{error}</div>}

      {!orderId ? (
        <button
          onClick={createOrder}
          disabled={loading}
          className="checkout-btn"
        >
          {loading ? "Creating order..." : "Create Order"}
        </button>
      ) : (
        <div className="checkout-payment">
          <p>Order #{orderId} created!</p>
          <button
            onClick={handlePayment}
            disabled={loading}
            className="checkout-btn"
          >
            {loading ? "Processing..." : "Pay Now"}
          </button>
        </div>
      )}
    </div>
  );
}

export default Checkout;
/* Production ready verion

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

function CheckoutForm({ orderId }: { orderId: number }) {
  const stripe = useStripe();
  const elements = useElements();
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();

    if (!stripe || !elements) {
      return;
    }

    setLoading(true);
    setError("");

    try {
      // Create PaymentIntent
      const { client_secret } = await apiFetch("/create-payment-intent", {
        method: "POST",
        body: JSON.stringify({ order_id: orderId }),
      });

      // Confirm payment with the card element
      const { error: stripeError, paymentIntent } = await stripe.confirmCardPayment(
        client_secret,
        {
          payment_method: {
            card: elements.getElement(CardElement)!,
          },
        }
      );

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
    <form onSubmit={handleSubmit}>
      <div className="card-element-container">
        <CardElement
          options={{
            style: {
              base: {
                fontSize: "16px",
                color: "#424770",
                "::placeholder": {
                  color: "#aab7c4",
                },
              },
              invalid: {
                color: "#9e2146",
              },
            },
          }}
        />
      </div>
      <button type="submit" disabled={!stripe || loading} className="checkout-btn">
        {loading ? "Processing..." : "Pay Now"}
      </button>
      {error && <div className="error-message">{error}</div>}
    </form>
  );
}

function Checkout() {
  const { token } = useAuth();
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [orderId, setOrderId] = useState<number | null>(null);

  useEffect(() => {
    if (!token) {
      navigate("/login");
    }
  }, [token, navigate]);

  const createOrder = async () => {
    setLoading(true);
    setError("");
    try {
      const order = await apiFetch("/orders", {
        method: "POST",
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
        <button onClick={createOrder} disabled={loading} className="checkout-btn">
          {loading ? "Creating order..." : "Create Order"}
        </button>
      ) : (
        <Elements stripe={stripePromise}>
          <CheckoutForm orderId={orderId} />
        </Elements>
      )}
    </div>
  );
}

export default Checkout;
*/
