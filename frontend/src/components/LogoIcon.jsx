import { useState } from 'react';

export default function LogoIcon({ className, imgClassName }) {
  const [imgError, setImgError] = useState(false);
  
  return imgError ? (
    <div className={className || "logo-fallback-icon"}>☕</div>
  ) : (
    <img 
      src="/logo.png" 
      alt="Bunna Logo" 
      className={imgClassName || className || "logo-img"} 
      onError={() => setImgError(true)} 
    />
  );
}
