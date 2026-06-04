import { useEffect, useRef, useState } from 'react';
import { Link } from 'react-router-dom';
import LeafScrollAnime from '../../components/LeafScrollAnime';
import './Landing.css';

/* ── Scroll-reveal hook ── */
function useReveal(threshold = 0.15) {
  const ref = useRef(null);
  useEffect(() => {
    const el = ref.current;
    if (!el) return;
    const obs = new IntersectionObserver(
      ([entry]) => { if (entry.isIntersecting) { el.classList.add('is-visible'); obs.unobserve(el); } },
      { threshold }
    );
    obs.observe(el);
    return () => obs.disconnect();
  }, [threshold]);
  return ref;
}

import LogoIcon from '../../components/LogoIcon';

/* ── SVG Icons ── */
const IconUpload = () => (
  <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">
    <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
    <polyline points="17 8 12 3 7 8"/>
    <line x1="12" y1="3" x2="12" y2="15"/>
  </svg>
);

const IconScan = () => (
  <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">
    <circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/>
  </svg>
);

const IconChart = () => (
  <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">
    <line x1="18" y1="20" x2="18" y2="10"/>
    <line x1="12" y1="20" x2="12" y2="4"/>
    <line x1="6" y1="20" x2="6" y2="14"/>
  </svg>
);

const IconLeaf = () => (
  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">
    <path d="M11 20A7 7 0 0 1 9.8 6.1C15.5 5 17 4.48 19 2c1 2 2 4.18 2 8 0 5.5-4.78 10-10 10z"/>
    <path d="M2 21c0-3 1.85-5.36 5.08-6C9.5 14.52 12 13 13 12"/>
  </svg>
);

const IconBolt = () => (
  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">
    <polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"/>
  </svg>
);

const IconShield = () => (
  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">
    <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
  </svg>
);

const IconPhone = () => (
  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">
    <rect x="5" y="2" width="14" height="20" rx="2" ry="2"/>
    <line x1="12" y1="18" x2="12.01" y2="18"/>
  </svg>
);

const IconGitHub = () => (
  <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
    <path d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0 1 12 6.844a9.59 9.59 0 0 1 2.504.337c1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.02 10.02 0 0 0 22 12.017C22 6.484 17.522 2 12 2z"/>
  </svg>
);

/* ── Navbar ── */
function Navbar() {
  const [scrolled, setScrolled] = useState(false);
  const [menuOpen, setMenuOpen] = useState(false);

  useEffect(() => {
    const onScroll = () => setScrolled(window.scrollY > 48);
    window.addEventListener('scroll', onScroll, { passive: true });
    return () => window.removeEventListener('scroll', onScroll);
  }, []);

  const scrollTo = (id) => (e) => {
    e.preventDefault();
    setMenuOpen(false);
    document.getElementById(id)?.scrollIntoView({ behavior: 'smooth' });
  };

  return (
    <nav className={`land-nav ${scrolled ? 'land-nav--scrolled' : ''}`} id="land-navbar">
      <a href="/" className="land-nav__logo" aria-label="Bunna — inicio">
        <LogoIcon className="land-nav__logo-icon" imgClassName="land-nav__logo-img" />
        <span className="land-nav__logo-text">Bunna</span>
      </a>

      <ul className="land-nav__links">
        <li><a href="#como-funciona" onClick={scrollTo('como-funciona')}>Cómo funciona</a></li>
        <li><a href="#caracteristicas" onClick={scrollTo('caracteristicas')}>Características</a></li>
        <li><a href="#nosotros" onClick={scrollTo('nosotros')}>Nosotros</a></li>
      </ul>

      <Link to="/login" className="land-nav__cta" id="nav-login-btn">
        Iniciar sesión
      </Link>

      {/* Mobile hamburger */}
      <button
        className={`land-nav__burger ${menuOpen ? 'land-nav__burger--open' : ''}`}
        onClick={() => setMenuOpen(!menuOpen)}
        aria-label="Abrir menú"
        id="nav-burger-btn"
      >
        <span /><span /><span />
      </button>

      {/* Mobile menu */}
      <div className={`land-nav__mobile ${menuOpen ? 'land-nav__mobile--open' : ''}`}>
        <a href="#como-funciona" onClick={scrollTo('como-funciona')}>Cómo funciona</a>
        <a href="#caracteristicas" onClick={scrollTo('caracteristicas')}>Características</a>
        <a href="#nosotros" onClick={scrollTo('nosotros')}>Nosotros</a>
        <Link to="/login" className="land-nav__mobile-cta" onClick={() => setMenuOpen(false)}>
          Iniciar sesión
        </Link>
      </div>
    </nav>
  );
}

/* ── Hero ── */
function Hero() {
  const imgRef = useRef(null);

  useEffect(() => {
    const img = imgRef.current;
    if (!img) return;
    if (img.complete) { img.classList.add('land-hero__bg--loaded'); return; }
    img.addEventListener('load', () => img.classList.add('land-hero__bg--loaded'));
  }, []);

  const scrollTo = (id) => (e) => {
    e.preventDefault();
    document.getElementById(id)?.scrollIntoView({ behavior: 'smooth' });
  };

  return (
    <section className="land-hero" aria-label="Sección principal">
      <img
        ref={imgRef}
        src="/landing-hero.webp"
        alt="Plantación de café en las montañas"
        className="land-hero__bg"
      />
      <div className="land-hero__overlay" aria-hidden="true" />

      <div className="land-hero__content">
        <span className="land-hero__tag">Diagnóstico visual de cultivos</span>
        <h1 className="land-hero__title">
          Tu cultivo merece<br />
          la ciencia más <em>precisa</em>
        </h1>
        <p className="land-hero__desc">
          Bunna detecta deficiencias de nitrógeno en hojas de café con una sola foto.
          Resultados en segundos, desde cualquier lugar del mundo.
        </p>
        <div className="land-hero__actions">
          <Link to="/login" className="land-btn-primary" id="hero-login-btn">
            Acceder a mi finca →
          </Link>
          <a href="#como-funciona" className="land-btn-ghost" onClick={scrollTo('como-funciona')}>
            Ver cómo funciona
          </a>
        </div>
      </div>

      <div className="land-hero__scroll" aria-hidden="true">
        <span className="land-hero__scroll-text">Descubre</span>
        <div className="land-hero__scroll-line" />
      </div>
    </section>
  );
}



/* ── How it works ── */
function HowItWorks() {
  const headerRef = useReveal(0.1);
  const stepsRef  = useReveal(0.1);

  const steps = [
    {
      icon: <IconUpload />,
      title: 'Sube una foto',
      desc:  'Toma una fotografía de la hoja de café desde tu celular o tablet y cárgala en la plataforma.',
    },
    {
      icon: <IconScan />,
      title: 'El sistema analiza',
      desc:  'El motor de visión detecta patrones visuales asociados a deficiencias de nitrógeno en cuestión de segundos.',
    },
    {
      icon: <IconChart />,
      title: 'Recibe el diagnóstico',
      desc:  'Obtienes el nivel de nitrógeno, la confianza del análisis y una recomendación de acción específica.',
    },
  ];

  return (
    <section className="land-section land-section--dark" id="como-funciona" aria-labelledby="how-title">
      <div className="land-how">
        <div className="land-how__header land-reveal" ref={headerRef}>
          <div>
            <p className="land-how__label">Proceso</p>
            <h2 className="land-how__title" id="how-title">
              Tres pasos para<br />conocer tu cultivo
            </h2>
          </div>
          <p className="land-how__desc">
            Sin equipos costosos, sin laboratorios, sin esperas. Solo tu teléfono y Bunna.
          </p>
        </div>

        <div className="land-steps land-reveal" ref={stepsRef}>
          {steps.map((s, i) => (
            <div className="land-step" key={i}>
              <span className="land-step__num" aria-hidden="true">0{i + 1}</span>
              <div className="land-step__icon">{s.icon}</div>
              <h3 className="land-step__title">{s.title}</h3>
              <p className="land-step__desc">{s.desc}</p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}

/* ── Features ── */
function Features() {
  const featureRef = useRef(null);

  useEffect(() => {
    const items = featureRef.current?.querySelectorAll('.land-feature-item');
    if (!items) return;
    const obs = new IntersectionObserver(
      ([entry]) => {
        if (entry.isIntersecting) {
          items.forEach(el => el.classList.add('is-visible'));
          obs.unobserve(entry.target);
        }
      },
      { threshold: 0.2 }
    );
    if (featureRef.current) obs.observe(featureRef.current);
    return () => obs.disconnect();
  }, []);

  const features = [
    { icon: <IconLeaf />,   title: 'Detección precisa',       sub: 'Análisis visual especializado en hojas de café.' },
    { icon: <IconBolt />,   title: 'Resultados instantáneos', sub: 'Diagnóstico completo en menos de 3 segundos.' },
    { icon: <IconShield />, title: 'Datos seguros',            sub: 'Tus imágenes y resultados están protegidos en la nube.' },
    { icon: <IconPhone />,  title: 'Desde cualquier lugar',   sub: 'Funciona en cualquier dispositivo con cámara y conexión.' },
  ];

  return (
    <section className="land-section land-section--light" id="caracteristicas" aria-labelledby="feat-title">
      <div className="land-features">
        <div className="land-features__visual">
          <img src="/landing-leaf.webp" alt="Hoja de café analizada por Bunna" className="land-features__img" />
          <div className="land-features__badge">
            <span className="land-features__badge-dot" aria-hidden="true" />
            <span className="land-features__badge-text">Análisis en tiempo real</span>
          </div>
        </div>

        <div className="land-features__body">
          <p className="land-features__label">Características</p>
          <h2 className="land-features__title" id="feat-title">
            Tecnología pensada<br />para el <em>caficultor</em>
          </h2>
          <p className="land-features__desc">
            Una herramienta diseñada para ser tan útil dentro de la finca como en la oficina.
            Sin curva de aprendizaje, sin instalaciones complejas.
          </p>

          <div className="land-feature-list" ref={featureRef}>
            {features.map((f, i) => (
              <div className="land-feature-item" key={i}>
                <div className="land-feature-item__icon">{f.icon}</div>
                <div className="land-feature-item__text">
                  <strong>{f.title}</strong>
                  <span>{f.sub}</span>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </section>
  );
}

/* ── About ── */
const TEAM = [
  { username: 'DavosJar',       url: 'https://github.com/DavosJar' },
  { username: 'IvanFernandez02', url: 'https://github.com/IvanFernandez02' },
  { username: 'CesarSTF',       url: 'https://github.com/CesarSTF' },
  { username: 'cesar050',       url: 'https://github.com/cesar050' },
  { username: 'Anthony-ggg',    url: 'https://github.com/Anthony-ggg' },
];

function About() {
  const leftRef  = useReveal(0.1);
  const rightRef = useReveal(0.1);

  return (
    <section className="land-section land-section--mid" id="nosotros" aria-labelledby="about-title">
      <div className="land-about">
        <div className="land-about__left land-reveal" ref={leftRef}>
          <h2 className="land-about__title" id="about-title">
            Construido por quienes<br />
            <em>viven el campo</em>
          </h2>
          <p className="land-about__desc">
            Bunna nació de la necesidad real de los caficultores de acceder a diagnósticos
            agronómicos rápidos y confiables, sin depender de laboratorios costosos ni de
            tiempos de espera largos.
          </p>
          <p className="land-about__desc">
            Somos un equipo multidisciplinario apasionado por la tecnología y el agro,
            comprometido con llevar herramientas de diagnóstico accesibles a quienes
            cuidan los cultivos cada día.
          </p>
          <div className="land-about__quote">
            <p>
              "La fertilización correcta en el momento correcto puede marcar la diferencia
              entre una cosecha ordinaria y una extraordinaria."
            </p>
            <cite>— Equipo Bunna</cite>
          </div>
        </div>

        <div className="land-about__visual land-reveal" ref={rightRef}>
          <p className="land-about__team-label">Equipo de desarrollo</p>
          <div className="land-team-grid">
            {TEAM.map((member) => (
              <a
                key={member.username}
                href={member.url}
                target="_blank"
                rel="noopener noreferrer"
                className="land-team-card"
                aria-label={`Perfil de GitHub de ${member.username}`}
              >
                <div className="land-team-card__avatar">
                  <img
                    src={`${member.url}.png?size=80`}
                    alt={member.username}
                    loading="lazy"
                    onError={(e) => { e.target.style.display = 'none'; }}
                  />
                </div>
                <div className="land-team-card__info">
                  <span className="land-team-card__name">{member.username}</span>
                  <span className="land-team-card__github">
                    <IconGitHub />
                    GitHub
                  </span>
                </div>
              </a>
            ))}
          </div>
        </div>
      </div>
    </section>
  );
}

/* ── CTA strip ── */
function CTAStrip() {
  const ref = useReveal(0.15);
  return (
    <section className="land-cta" aria-labelledby="cta-title">
      <div className="land-cta__content land-reveal" ref={ref}>
        <h2 className="land-cta__title" id="cta-title">
          Tu cultivo merece un<br />
          diagnóstico <em>inteligente</em>
        </h2>
        <p className="land-cta__desc">
          Accede a tu cuenta y comienza a analizar tus hojas hoy mismo. Sin costos adicionales.
        </p>
        <Link to="/login" className="land-btn-primary" id="cta-login-btn">
          Entrar a Bunna →
        </Link>
      </div>
    </section>
  );
}

/* ── Footer ── */
function Footer() {
  return (
    <footer className="land-footer">
      <a href="/" className="land-footer__logo" aria-label="Bunna">
        <LogoIcon className="land-nav__logo-icon" imgClassName="land-nav__logo-img" />
        <span className="land-footer__logo-text">Bunna</span>
      </a>
      <p className="land-footer__copy">© {new Date().getFullYear()} Bunna. Todos los derechos reservados.</p>
      <nav className="land-footer__links" aria-label="Links del footer">
        <a href="#como-funciona" onClick={(e) => { e.preventDefault(); document.getElementById('como-funciona')?.scrollIntoView({ behavior: 'smooth' }); }}>Cómo funciona</a>
        <a href="#nosotros" onClick={(e) => { e.preventDefault(); document.getElementById('nosotros')?.scrollIntoView({ behavior: 'smooth' }); }}>Nosotros</a>
        <Link to="/login">Ingresar</Link>
      </nav>
    </footer>
  );
}

/* ── Page ── */
export default function LandingPage() {
  return (
    <>
      <Navbar />
      <main>
        <Hero />
        <HowItWorks />
        <Features />
        <LeafScrollAnime />
        <About />
        <CTAStrip />
      </main>
      <Footer />
    </>
  );
}
