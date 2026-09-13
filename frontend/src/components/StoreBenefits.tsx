import { Truck, ShieldCheck, RotateCcw, Heart } from "lucide-react";

const benefits = [
  {
    icon: Truck,
    title: "Fast & Reliable Shipping",
    description: "Get your order delivered with care.",
  },
  {
    icon: ShieldCheck,
    title: "Secure Payments",
    description: "Shop with confidence.",
  },
  {
    icon: RotateCcw,
    title: "Easy Returns",
    description: "Hassle-free within 14 days.",
  },
  {
    icon: Heart,
    title: "A Culture of Style",
    description: "Authentic. Elegant. Timeless.",
  },
];

function StoreBenefits() {
  return (
    <section className="store-benefits">
      <div className="benefits-container">
        {benefits.map((benefit, index) => {
          const Icon = benefit.icon;

          return (
            <div className="benefit" key={benefit.title}>
              <Icon className="benefit-icon" />

              <h3>{benefit.title}</h3>

              <p>{benefit.description}</p>

              {index < benefits.length - 1 && (
                <div className="benefit-divider" />
              )}
            </div>
          );
        })}
      </div>
    </section>
  );
}

export default StoreBenefits;
