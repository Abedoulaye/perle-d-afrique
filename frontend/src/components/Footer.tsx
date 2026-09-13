import { Link } from "react-router";
import { FaInstagram, FaFacebookF, FaPinterestP } from "react-icons/fa";
import { FaTiktok } from "react-icons/fa6";

function Footer() {
  return (
    <footer className="footer">
      <div className="footer-main">
        <Link to="/" className="footer-brand">
          <img
            src="/logo.png"
            alt="Diariata's Fabrics"
            className="footer-logo"
          />

          <div className="footer-name">
            <h2>DIARIATA'S</h2>
            <p>FABRICS</p>
          </div>
        </Link>

        <nav className="footer-links">
          <Link to="/">Home</Link>
          <Link to="/products">Shop</Link>
          <Link to="/about">About</Link>
          <Link to="/contact">Contact</Link>
          <Link to="/faq">FAQ</Link>
        </nav>

        <div className="footer-socials">
          <a href="#" aria-label="Instagram">
            <FaInstagram />
          </a>

          <a href="#" aria-label="Facebook">
            <FaFacebookF />
          </a>

          <a href="#" aria-label="Pinterest">
            <FaPinterestP />
          </a>

          <a href="#" aria-label="TikTok">
            <FaTiktok />
          </a>
        </div>
      </div>

      <div className="footer-bottom">
        <p> &copy; 2025 Diariata's Fabrics. All rights reserved.</p>
      </div>
    </footer>
  );
}

export default Footer;
