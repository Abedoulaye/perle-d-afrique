function Contact() {
  return (
    <div className="info-page">
      <section className="info-hero">
        <p className="info-eyebrow">GET IN TOUCH</p>
        <h1>Contact Us</h1>
        <p>
          Have a question about an order, product, or anything else? We'd love
          to hear from you.
        </p>
      </section>

      <section className="contact-section">
        <div className="contact-container">
          <div className="contact-info">
            <h2>Let's Talk</h2>

            <p>
              If you have a question about our products, orders, shipping, or
              anything else, feel free to reach out.
            </p>

            <div className="contact-detail">
              <h3>Questions About Your Order?</h3>
              <p>
                If you're contacting us about an existing order, please include
                your order number so we can help you more quickly.
              </p>
            </div>

            <div className="contact-detail">
              <h3>Looking for Something Specific?</h3>
              <p>
                Feel free to ask about our fabrics, clothing, availability,
                sizing, or any other product questions.
              </p>
            </div>

            <div className="contact-detail">
              <h3>Location</h3>
              <p>[Philadelphia, PA / United States]</p>
            </div>

            <div className="contact-detail">
              <h3>Business Hours</h3>
              <p>[8:00 A.M to 8:00 P.M]</p>
            </div>
          </div>

          <form className="contact-form">
            <div className="form-group">
              <label htmlFor="contact-name">Name</label>
              <input
                id="contact-name"
                type="text"
                placeholder="Your name"
                required
              />
            </div>

            <div className="form-group">
              <label htmlFor="contact-email">Email</label>
              <input
                id="contact-email"
                type="email"
                placeholder="you@example.com"
                required
              />
            </div>

            <div className="form-group">
              <label htmlFor="contact-subject">Subject</label>
              <input
                id="contact-subject"
                type="text"
                placeholder="How can we help?"
                required
              />
            </div>

            <div className="form-group">
              <label htmlFor="contact-message">Message</label>
              <textarea
                id="contact-message"
                placeholder="Write your message..."
                required
              />
            </div>

            <button type="button" className="contact-submit-btn">
              Send Message
            </button>
          </form>
        </div>
      </section>
    </div>
  );
}

export default Contact;
