import { useState } from "react";

const faqs = [
  {
    question: "What types of products do you sell?",
    answer:
      "We offer a collection of fabrics, women's clothing, accessories, and other styles inspired by African and global fashion.",
  },
  {
    question: "Where do you ship?",
    answer:
      "[Add your shipping locations here. For example: We currently ship throughout the United States and internationally to select countries.]",
  },
  {
    question: "How long does shipping take?",
    answer:
      "[Add your shipping timeframe here. For example: Orders are typically processed within 1–3 business days and delivered within 3–7 business days.]",
  },
  {
    question: "What is your return policy?",
    answer:
      "[Add your return policy here. For example: We accept returns within 14 days of delivery as long as the item meets our return requirements.]",
  },
  {
    question: "How can I track my order?",
    answer:
      "[Add your order tracking information here. You can later connect this section to your order/shipping system.]",
  },
  {
    question: "What payment methods do you accept?",
    answer:
      "[Add your accepted payment methods here. For example: We accept major credit and debit cards through our secure payment system.]",
  },
  {
    question: "How can I contact you about an order?",
    answer:
      "You can contact us through our Contact page. Please include your order number if your question is about an existing order.",
  },
];

function FAQ() {
  const [openIndex, setOpenIndex] = useState<number | null>(null);

  const toggleFAQ = (index: number) => {
    setOpenIndex(openIndex === index ? null : index);
  };

  return (
    <div className="info-page">
      <section className="info-hero">
        <p className="info-eyebrow">HELP CENTER</p>
        <h1>Frequently Asked Questions</h1>
        <p>
          Find answers to some of the most common questions about Diariata's
          Fabrics.
        </p>
      </section>

      <section className="faq-section">
        <div className="faq-container">
          {faqs.map((faq, index) => (
            <div
              className={`faq-item ${
                openIndex === index ? "faq-item-open" : ""
              }`}
              key={faq.question}
            >
              <button
                className="faq-question"
                onClick={() => toggleFAQ(index)}
                aria-expanded={openIndex === index}
              >
                <span>{faq.question}</span>
                <span className="faq-icon">
                  {openIndex === index ? "−" : "+"}
                </span>
              </button>

              {openIndex === index && (
                <div className="faq-answer">
                  <p>{faq.answer}</p>
                </div>
              )}
            </div>
          ))}
        </div>
      </section>
    </div>
  );
}

export default FAQ;
