import Hero from "../components/Hero/Hero";
import Portfolio from "../components/Portfolio/Portfolio";
import Reviews from "../components/Reviews/Reviews";
import BookingForm from "../components/BookingForm/BookingForm";
import Contacts from "../components/Contacts/Contacts";

export default function Home() {
  return (
    <main>
      <Hero />
      <Portfolio />
      <Reviews />
      <BookingForm />
      <Contacts />
    </main>
  );
}
