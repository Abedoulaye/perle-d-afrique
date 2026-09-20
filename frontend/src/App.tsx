import { Routes, Route } from "react-router";
import Navbar from "./components/Navbar";
import ProductList from "./pages/ProductList";
import ProductDetail from "./pages/ProductDetail";
import Login from "./pages/Login";
import Register from "./pages/Register";
import Cart from "./pages/Cart";
import Checkout from "./pages/Checkout";
import Orders from "./pages/Orders";
import Admin from "./pages/Admin";
import Hero from "./components/Hero";
import StoreBenefits from "./components/StoreBenefits";
import Footer from "./components/Footer";
import FeaturedProducts from "./components/FeaturedProducts";
import VerifyEmail from "./pages/VerifyEmail";
import About from "./pages/About";
import Contact from "./pages/Contact";
import FAQ from "./pages/FAQ";

function App() {
  return (
    <>
      <Routes>
        <Route
          path="/"
          element={
            <>
              <Navbar />
              <Hero />
              <FeaturedProducts />
              <StoreBenefits />
              <Footer />
            </>
          }
        />
        <Route
          path="/products"
          element={
            <>
              <ProductList />
              <Footer />
            </>
          }
        />
        <Route
          path="/products/:id"
          element={
            <>
              <Navbar />
              <ProductDetail />
              <Footer />
            </>
          }
        />
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />
        <Route
          path="/cart"
          element={
            <>
              <Navbar />
              <Cart />
            </>
          }
        />
        <Route
          path="/checkout"
          element={
            <>
              <Navbar />
              <Checkout />
            </>
          }
        />
        <Route
          path="/orders"
          element={
            <>
              <Navbar />
              <Orders />
              <Footer />
            </>
          }
        />
        <Route
          path="/admin"
          element={
            <>
              <Navbar />
              <Admin />
              <Footer />
            </>
          }
        />
        <Route path="/verify" element={<VerifyEmail />} />
        <Route
          path="/about"
          element={
            <>
              <About />
              <Footer />
            </>
          }
        />
        <Route
          path="/contact"
          element={
            <>
              <Contact />
              <Footer />
            </>
          }
        />
        <Route
          path="/faq"
          element={
            <>
              <FAQ />
              <Footer />
            </>
          }
        />
      </Routes>
    </>
  );
}

export default App;
