import { useEffect, useRef, useState } from 'react';
import anime from 'animejs';
import './LeafScrollAnime.css';

export default function LeafScrollAnime() {
  const containerRef = useRef(null);
  const leafBodyRef = useRef(null);
  const leafTimelineRef = useRef(null);
  const veinsTimelineRef = useRef(null);
  const [scrollProgress, setScrollProgress] = useState(0);

  const isHealthy = scrollProgress < 0.5;

  useEffect(() => {
    if (!leafBodyRef.current) return;

    /* ── Timeline principal: color de la hoja ── */
    leafTimelineRef.current = anime({
      targets: leafBodyRef.current,
      fill: ['#2d6a3f', '#7a9a30', '#c4a820', '#e8b840'],
      easing: 'linear',
      autoplay: false,
      duration: 1000,
    });

    /* ── Timeline secundaria: nervaduras se desvanecen ── */
    veinsTimelineRef.current = anime({
      targets: '.anime-leaf__vein',
      opacity: [0.35, 0.12],
      strokeWidth: ['2.5px', '1.8px'],
      easing: 'linear',
      autoplay: false,
      duration: 1000,
    });

    /* ── Scroll listener ── */
    const handleScroll = () => {
      const el = containerRef.current;
      if (!el) return;

      const rect = el.getBoundingClientRect();
      const vh = window.innerHeight;
      
      // Solo calcular el progreso mientras el contenedor está "pegado" (rect.top <= 0)
      // La distancia scrolleable es la altura total del contenedor menos el viewport
      const scrollableDistance = rect.height - vh;
      const progress = Math.min(Math.max(-rect.top / scrollableDistance, 0), 1);

      setScrollProgress(progress);

      if (leafTimelineRef.current) {
        leafTimelineRef.current.seek(progress * leafTimelineRef.current.duration);
      }
      if (veinsTimelineRef.current) {
        veinsTimelineRef.current.seek(progress * veinsTimelineRef.current.duration);
      }
    };

    window.addEventListener('scroll', handleScroll, { passive: true });
    handleScroll();

    return () => window.removeEventListener('scroll', handleScroll);
  }, []);

  /* ── Derived visual values ── */
  let rotation = scrollProgress * 15;
  let leafScale = 1 + scrollProgress * 0.06;
  
  // Lógica de caída (desprendimiento) a partir del 75%
  let fallY = 0;
  let fallX = 0;
  
  if (scrollProgress > 0.75) {
    const fallProgress = (scrollProgress - 0.75) / 0.25; // De 0 a 1
    fallY = Math.pow(fallProgress, 2) * 220; // Aceleración de caída simulada
    fallX = fallProgress * 30; // Vuelo hacia la derecha
    rotation += fallProgress * 45; // Giro adicional al caer
  }

  const glowColor = isHealthy
    ? 'rgba(45, 106, 63, 0.3)'
    : 'rgba(232, 184, 64, 0.3)';

  return (
    <section className="anime-showcase" ref={containerRef} aria-label="Simulador visual de deficiencia">
      {/* Ambient glow */}
      <div
        className="anime-showcase__glow"
        style={{
          background: `radial-gradient(ellipse 60% 50% at 65% 50%, ${glowColor} 0%, transparent 70%)`,
        }}
      />

      <div className="anime-showcase__sticky">
        <div className="anime-showcase__grid">

          {/* ─── Left: Info panel ─── */}
          <div className="anime-showcase__info">
            <span className="anime-showcase__eyebrow">Simulador de diagnóstico</span>

            <div className="anime-showcase__cards">
              {/* Healthy card */}
              <div className={`anime-card ${isHealthy ? 'anime-card--active' : ''}`}>
                <div className="anime-card__dot anime-card__dot--green" />
                <h3 className="anime-card__title">Hoja Saludable</h3>
                <p className="anime-card__desc">
                  El color verde intenso y uniforme indica una concentración óptima
                  de clorofila. La planta absorbe nitrógeno correctamente, garantizando
                  un crecimiento vigoroso y buena producción.
                </p>
              </div>

              {/* Deficient card */}
              <div className={`anime-card ${!isHealthy ? 'anime-card--active' : ''}`}>
                <div className="anime-card__dot anime-card__dot--yellow" />
                <h3 className="anime-card__title">Deficiencia de Nitrógeno</h3>
                <p className="anime-card__desc">
                  Aparece una clorosis generalizada que inicia en las hojas más viejas.
                  El color se torna amarillento por degradación de la clorofila, reduciendo
                  la fotosíntesis y el rendimiento del cultivo.
                </p>
              </div>
            </div>

            {/* Progress gauge */}
            <div className="anime-gauge">
              <div className="anime-gauge__labels">
                <span className={isHealthy ? 'is-active' : ''}>Saludable</span>
                <span className={!isHealthy ? 'is-active' : ''}>Deficiente</span>
              </div>
              <div className="anime-gauge__track">
                <div
                  className="anime-gauge__fill"
                  style={{ width: `${scrollProgress * 100}%` }}
                />
              </div>
            </div>
          </div>

          {/* ─── Right: Leaf SVG ─── */}
          <div className="anime-showcase__visual">
            <div className="anime-leaf__wrapper">
              <svg
                viewBox="0 0 240 340"
                className="anime-leaf__svg"
                aria-hidden="true"
                style={{
                  filter: `drop-shadow(0 24px 48px ${glowColor})`,
                  overflow: 'visible'
                }}
              >
                {/* ── Branch (Static) ── */}
                <g className="anime-leaf__branch">
                  <path 
                    d="M10,20 C60,40 100,45 120,45 C150,45 200,10 230,0" 
                    fill="none" 
                    stroke="#4a3721" 
                    strokeWidth="8" 
                    strokeLinecap="round" 
                  />
                  <path 
                    d="M100,38 C108,30 115,20 125,10" 
                    fill="none" 
                    stroke="#4a3721" 
                    strokeWidth="5" 
                    strokeLinecap="round" 
                  />
                </g>

                {/* ── Leaf (Animated & Detaching) ── */}
                {/* Usamos transformOrigin en la base del tallo para que gire desde ahí */}
                <g 
                  style={{ 
                    transform: `translate(${fallX}px, ${fallY}px) rotate(${rotation}deg) scale(${leafScale})`,
                    transformOrigin: '120px 45px',
                    transition: 'transform 0.1s linear'
                  }}
                >
                  {/* Main leaf shape */}
                  <path
                    ref={leafBodyRef}
                    className="anime-leaf__body"
                    d="M120,45 C165,95 195,150 195,205 C195,265 160,295 120,310 C80,295 45,265 45,205 C45,150 75,95 120,45 Z"
                    fill="#2d6a3f"
                  />

                  {/* Shading overlay */}
                  <defs>
                    <linearGradient id="leafShading" x1="0%" y1="0%" x2="100%" y2="100%">
                      <stop offset="0%" stopColor="#ffffff" stopOpacity="0.12" />
                      <stop offset="50%" stopColor="#ffffff" stopOpacity="0" />
                      <stop offset="100%" stopColor="#000000" stopOpacity="0.18" />
                    </linearGradient>
                  </defs>
                  <path
                    d="M120,45 C165,95 195,150 195,205 C195,265 160,295 120,310 C80,295 45,265 45,205 C45,150 75,95 120,45 Z"
                    fill="url(#leafShading)"
                  />

                  {/* Central rib */}
                  <line
                    x1="120" y1="45" x2="120" y2="310"
                    stroke="#ffffff" strokeWidth="3" strokeLinecap="round"
                    opacity="0.35"
                  />

                  {/* Lateral veins */}
                  <g className="anime-leaf__veins">
                    <path className="anime-leaf__vein" d="M120,85   C145,98   170,125  180,150"  stroke="#fff" fill="none" strokeLinecap="round" />
                    <path className="anime-leaf__vein" d="M120,85   C95,98    70,125   60,150"   stroke="#fff" fill="none" strokeLinecap="round" />

                    <path className="anime-leaf__vein" d="M120,130  C150,148  178,178  188,210"  stroke="#fff" fill="none" strokeLinecap="round" />
                    <path className="anime-leaf__vein" d="M120,130  C90,148   62,178   52,210"   stroke="#fff" fill="none" strokeLinecap="round" />

                    <path className="anime-leaf__vein" d="M120,180  C148,202  175,230  182,260"  stroke="#fff" fill="none" strokeLinecap="round" />
                    <path className="anime-leaf__vein" d="M120,180  C92,202   65,230   58,260"   stroke="#fff" fill="none" strokeLinecap="round" />

                    <path className="anime-leaf__vein" d="M120,230  C140,248  162,270  168,288"  stroke="#fff" fill="none" strokeLinecap="round" />
                    <path className="anime-leaf__vein" d="M120,230  C100,248  78,270   72,288"   stroke="#fff" fill="none" strokeLinecap="round" />
                  </g>

                  {/* Subtle leaf outline */}
                  <path
                    d="M120,45 C165,95 195,150 195,205 C195,265 160,295 120,310 C80,295 45,265 45,205 C45,150 75,95 120,45 Z"
                    fill="none"
                    stroke={isHealthy ? '#1a4a2c' : '#a88018'}
                    strokeWidth="1.2"
                    opacity="0.35"
                  />
                </g>
              </svg>
            </div>

            {/* Floating badge */}
            <div className="anime-leaf__badge">
              <span className={`anime-leaf__badge-dot ${isHealthy ? 'is-green' : 'is-yellow'}`} />
              <span className="anime-leaf__badge-text">
                {isHealthy ? 'Hoja Saludable' : 'Deficiencia de N'}
              </span>
            </div>
          </div>

        </div>
      </div>
    </section>
  );
}
