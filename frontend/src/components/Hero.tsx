import { useState } from "react";
import carousel1 from "../assets/carousel1.png";
import carousel2 from "../assets/carousel2.png";
import carousel3 from "../assets/carousel3.png";

const slides = [
  {
    image: carousel1,
    eyebrow: "TRADITION MEETS TIMELESS STYLE",
    title: "Beautiful Fabrics. Unique Styles.",
    description:
      "Explore our collection of authentic African, Middle Eastern, and global-inspired clothing and fabrics designed for confidence, culture, and you.",
  },
  {
    image: carousel2,
    eyebrow: "STYLE WITH CULTURE",
    title: "Wear Something Beautiful.",
    description:
      "Discover distinctive clothing and fabrics that bring timeless traditions and modern style together.",
  },
  {
    image: carousel3,
    eyebrow: "CULTURE IN EVERY THREAD",
    title: "Find Your Next Favorite.",
    description:
      "Explore our collection of unique fabrics, clothing, and accessories.",
  },
];

function Hero() {
  const [currentSlide, setCurrentSlide] = useState(0);

  const nextSlide = () => {
    setCurrentSlide((current) =>
      current === slides.length - 1 ? 0 : current + 1,
    );
  };

  const previousSlide = () => {
    setCurrentSlide((current) =>
      current === 0 ? slides.length - 1 : current - 1,
    );
  };

  const slide = slides[currentSlide];

  return (
    <section
      className="hero"
      style={{ backgroundImage: `url(${slide.image})` }}
    >
      <button
        className="hero-arrow hero-arrow-left"
        onClick={previousSlide}
        aria-label="Previous slide"
      >
        &#10094;
      </button>

      <div className="hero-content">
        <p className="hero-eyebrow">{slide.eyebrow}</p>

        <h1>{slide.title}</h1>

        <p className="hero-description">{slide.description}</p>

        <a href="/products" className="hero-button">
          <span>Shop Now </span>
          <svg
            xmlns="http://www.w3.org/2000/svg"
            height="19px"
            viewBox="0 -960 960 960"
            width="24px"
            fill="#000"
            className="arrow"
          >
            <path d="M647-440H160v-80h487L423-744l57-56 320 320-320 320-57-56 224-224Z" />
          </svg>
        </a>
      </div>

      <button
        className="hero-arrow hero-arrow-right"
        onClick={nextSlide}
        aria-label="Next slide"
      >
        &#10095;
      </button>

      <div className="hero-dots">
        {slides.map((_, index) => (
          <button
            key={index}
            className={currentSlide === index ? "active" : ""}
            onClick={() => setCurrentSlide(index)}
            aria-label={`Go to slide ${index + 1}`}
          />
        ))}
      </div>
    </section>
  );
}

export default Hero;
