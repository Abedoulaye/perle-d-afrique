import { useState, useEffect, type SubmitEvent } from "react";
import { useNavigate } from "react-router";
import { useAuth } from "../context/AuthContext";
import {
  getProducts,
  createProduct,
  updateProduct,
  deleteProduct,
} from "../services/productService";
import type { Product } from "../types";

function Admin() {
  const { isAdmin, token } = useAuth();
  const navigate = useNavigate();
  const [products, setProducts] = useState<Product[]>([]);
  const [editingProduct, setEditingProduct] = useState<Product | null>(null);
  const [showForm, setShowForm] = useState(false);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(true);

  // Form state
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [priceCents, setPriceCents] = useState("");
  const [stock, setStock] = useState("");
  const [image, setImage] = useState("");

  useEffect(() => {
    if (!token || !isAdmin) {
      navigate("/");
      return;
    }
    fetchProducts();
  }, [token, isAdmin, navigate]);

  const fetchProducts = async () => {
    try {
      const data = await getProducts();
      setProducts(data);
    } catch (err) {
      setError("Failed to load products");
    } finally {
      setLoading(false);
    }
  };

  const resetForm = () => {
    setName("");
    setDescription("");
    setPriceCents("");
    setStock("");
    setImage("");
    setEditingProduct(null);
    setShowForm(false);
  };

  const handleEdit = (product: Product) => {
    setEditingProduct(product);
    setName(product.name);
    setDescription(product.description);
    setPriceCents(String(product.price_cents));
    setStock(String(product.stock));
    setImage(product.image);
    setShowForm(true);
  };

  const handleDelete = async (id: number) => {
    if (!confirm("Are you sure you want to delete this product?")) return;

    try {
      await deleteProduct(id);
      fetchProducts();
    } catch (err) {
      setError("Failed to delete product");
    }
  };

  const handleSubmit = async (e: SubmitEvent) => {
    e.preventDefault();
    setError("");

    const productData = {
      name,
      description,
      price_cents: Number(priceCents),
      stock: Number(stock),
      image,
    };

    try {
      if (editingProduct) {
        await updateProduct(editingProduct.id, productData);
      } else {
        await createProduct(productData);
      }
      resetForm();
      fetchProducts();
    } catch (err) {
      setError("Failed to save product");
    }
  };

  if (loading) return <div className="loading">Loading admin panel...</div>;

  return (
    <div className="admin-container">
      <h1>Admin Panel</h1>
      {error && <div className="error-message">{error}</div>}

      <button onClick={() => setShowForm(!showForm)} className="admin-add-btn">
        {showForm ? "Cancel" : "Add New Product"}
      </button>

      {showForm && (
        <form onSubmit={handleSubmit} className="admin-form">
          <h2>{editingProduct ? "Edit Product" : "New Product"}</h2>
          <div className="form-group">
            <label>Name</label>
            <input
              value={name}
              onChange={(e) => setName(e.target.value)}
              required
            />
          </div>
          <div className="form-group">
            <label>Description</label>
            <textarea
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              required
            />
          </div>
          <div className="form-group">
            <label>Price (in cents)</label>
            <input
              type="number"
              value={priceCents}
              onChange={(e) => setPriceCents(e.target.value)}
              required
            />
          </div>
          <div className="form-group">
            <label>Stock</label>
            <input
              type="number"
              value={stock}
              onChange={(e) => setStock(e.target.value)}
              required
            />
          </div>
          <div className="form-group">
            <label>Image URL</label>
            <input
              value={image}
              onChange={(e) => setImage(e.target.value)}
              required
            />
          </div>
          <button type="submit" className="admin-submit-btn">
            {editingProduct ? "Update Product" : "Create Product"}
          </button>
        </form>
      )}

      <div className="admin-products">
        <h2>Products ({products.length})</h2>
        <table className="admin-table">
          <thead>
            <tr>
              <th>ID</th>
              <th>Image</th>
              <th>Name</th>
              <th>Price</th>
              <th>Stock</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {products.map((product) => (
              <tr key={product.id}>
                <td>{product.id}</td>
                <td>
                  <img
                    src={product.image}
                    alt={product.name}
                    className="admin-thumbnail"
                  />
                </td>
                <td>{product.name}</td>
                <td>${(product.price_cents / 100).toFixed(2)}</td>
                <td>{product.stock}</td>
                <td>
                  <button
                    onClick={() => handleEdit(product)}
                    className="admin-edit-btn"
                  >
                    Edit
                  </button>
                  <button
                    onClick={() => handleDelete(product.id)}
                    className="admin-delete-btn"
                  >
                    Delete
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}

export default Admin;
