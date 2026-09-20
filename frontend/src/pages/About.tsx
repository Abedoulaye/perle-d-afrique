import { Link } from "react-router";

function About() {
  return (
    <div className="info-page">
      <section className="info-hero">
        <p className="info-eyebrow">ABOUT US</p>
        <h1>About Diariata's Fabrics</h1>
        <p>
          Discover the story behind our collection of beautiful fabrics,
          clothing, and timeless styles.
        </p>
      </section>

      <section className="info-section">
        <div className="info-section-content">
          <h2>Our Story</h2>
          <p>
            Diariata's Fabrics was created with a love for beautiful fabrics,
            clothing, culture, and timeless style.
          </p>

          <p>
            At Diariata's Fabrics, we bring together beautiful fabrics, elegant
            clothing, and distinctive accessories for customers who appreciate
            culture, creativity, and timeless style.
          </p>
        </div>
      </section>

      <section className="info-section info-section-alt">
        <div className="info-section-content">
          <h2>What We Offer</h2>

          <div className="info-grid">
            <div>
              <h3>Fabrics</h3>
              <p>
                Explore a collection of distinctive fabrics for creating
                beautiful looks.
              </p>
            </div>

            <div>
              <h3>Women's Clothing</h3>
              <p>
                Discover elegant clothing inspired by traditional and modern
                styles.
              </p>
            </div>

            <div>
              <h3>Accessories</h3>
              <p>
                Find pieces that complete your look and add personality to every
                outfit.
              </p>
            </div>
          </div>
        </div>
      </section>

      <section className="info-section">
        <div className="info-section-content">
          <h2>Our Mission</h2>
          <p>
            Our goal is to make beautiful, distinctive styles accessible while
            celebrating culture, craftsmanship, and individuality.
          </p>

          <Link to="/products" className="info-button">
            Shop Our Collection →
          </Link>
        </div>
      </section>
    </div>
  );
}

export default About;
