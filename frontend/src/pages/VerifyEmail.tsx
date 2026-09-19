import { useEffect, useState } from "react";
import { useSearchParams, Link } from "react-router";
import { apiFetch } from "../services/api";

function VerifyEmail() {
  const [searchParams] = useSearchParams();
  const [status, setStatus] = useState<"loading" | "success" | "error">(
    "loading",
  );
  const [message, setMessage] = useState("");

  useEffect(() => {
    const token = searchParams.get("token");

    if (!token) {
      setStatus("error");
      setMessage("No verification token provided");
      return;
    }

    apiFetch("/verify-email", {
      method: "POST",
      body: JSON.stringify({ token }),
    })
      .then(() => {
        setStatus("success");
        setMessage("Email verified! You can now log in.");
      })
      .catch((err) => {
        setStatus("error");
        setMessage(err.message || "Verification failed");
      });
  }, []);

  return (
    <div className="verify-container">
      {status === "loading" && <p>Verifying...</p>}
      {status === "success" && (
        <>
          <h1>✅ {message}</h1>
          <Link to="/login">Go to Login</Link>
        </>
      )}
      {status === "error" && (
        <>
          <h1>❌ {message}</h1>
          <Link to="/">Go Home</Link>
        </>
      )}
    </div>
  );
}

export default VerifyEmail;
