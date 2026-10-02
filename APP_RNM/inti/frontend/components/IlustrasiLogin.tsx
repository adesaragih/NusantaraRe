// Ilustrasi meja kerja halaman login - dari loginbaru.html (lampiran work owner
// 01-10-2026). Atribut SVG ditulis ulang ke bentuk JSX; id gradient/filter
// berawalan `masuk-`. Hiasan saja: `aria-hidden`.

export default function IlustrasiLogin() {
  return (
    <svg viewBox="0 0 400 420" aria-hidden="true" focusable="false">
    <defs>
    <linearGradient id="masuk-deskTop" x1="0" y1="0" x2="0" y2="1"><stop offset="0" stopColor="#ffe7c2" /><stop offset="1" stopColor="#f5c98c" /></linearGradient>
    <linearGradient id="masuk-deskFront" x1="0" y1="0" x2="0" y2="1"><stop offset="0" stopColor="#eab676" /><stop offset="1" stopColor="#d89a57" /></linearGradient>
    <linearGradient id="masuk-leg" x1="0" y1="0" x2="1" y2="0"><stop offset="0" stopColor="#d89a57" /><stop offset="1" stopColor="#b97a3c" /></linearGradient>
    <linearGradient id="masuk-laptop" x1="0" y1="0" x2="0" y2="1"><stop offset="0" stopColor="#f4f5fb" /><stop offset="1" stopColor="#c9cde2" /></linearGradient>
    <linearGradient id="masuk-screen" x1="0" y1="0" x2="1" y2="1"><stop offset="0" stopColor="#3a3470" /><stop offset="1" stopColor="#231f4d" /></linearGradient>
    <linearGradient id="masuk-mug" x1="0" y1="0" x2="1" y2="0"><stop offset="0" stopColor="#ff8f94" /><stop offset="1" stopColor="#e5484d" /></linearGradient>
    <linearGradient id="masuk-pot" x1="0" y1="0" x2="1" y2="0"><stop offset="0" stopColor="#ffffff" /><stop offset="1" stopColor="#dcd3f5" /></linearGradient>
    <linearGradient id="masuk-leaf" x1="0" y1="0" x2="0" y2="1"><stop offset="0" stopColor="#8ee3b3" /><stop offset="1" stopColor="#3aa872" /></linearGradient>
    <filter id="masuk-shadow" x="-30%" y="-30%" width="160%" height="160%">
    <feDropShadow dx="0" dy="10" stdDeviation="10" floodColor="#4c1d95" floodOpacity="0.25" />
    </filter>
    </defs>

    <ellipse cx="200" cy="394" rx="160" ry="16" fill="#6d3fc4" opacity="0.22" />

    <rect x="82" y="276" width="14" height="114" rx="7" fill="url(#masuk-leg)" />
    <rect x="304" y="276" width="14" height="114" rx="7" fill="url(#masuk-leg)" />
    <rect x="96" y="330" width="208" height="8" rx="4" fill="#c98a49" opacity="0.7" />
    <rect x="50" y="262" width="300" height="18" rx="8" fill="url(#masuk-deskFront)" />
    <rect x="50" y="252" width="300" height="16" rx="8" fill="url(#masuk-deskTop)" />

    <ellipse cx="60" cy="192" rx="9" ry="26" transform="rotate(-28 60 192)" fill="url(#masuk-leaf)" />
    <ellipse cx="86" cy="190" rx="9" ry="28" transform="rotate(22 86 190)" fill="url(#masuk-leaf)" />
    <ellipse cx="55" cy="205" rx="7" ry="17" transform="rotate(-58 55 205)" fill="url(#masuk-leaf)" />
    <ellipse cx="92" cy="204" rx="7" ry="17" transform="rotate(56 92 204)" fill="url(#masuk-leaf)" />
    <ellipse cx="73" cy="186" rx="8" ry="32" fill="url(#masuk-leaf)" />
    <path d="M57 218h32l-4 31a4 4 0 0 1-4 3H65a4 4 0 0 1-4-3z" fill="url(#masuk-pot)" />

    <g filter="url(#masuk-shadow)">
    <rect x="122" y="128" width="160" height="112" rx="12" fill="#e7e9f4" />
    <rect x="131" y="137" width="142" height="94" rx="7" fill="url(#masuk-screen)" />
    <rect x="143" y="149" width="44" height="7" rx="3.5" fill="#b98df3" />
    <rect x="143" y="162" width="70" height="5" rx="2.5" fill="#4a4486" />
    <rect x="222" y="149" width="38" height="7" rx="3.5" fill="#4a4486" />
    <rect x="149" y="196" width="12" height="24" rx="3" fill="#ff9aa0" />
    <rect x="167" y="182" width="12" height="38" rx="3" fill="#b98df3" />
    <rect x="185" y="190" width="12" height="30" rx="3" fill="#8ee3b3" />
    <circle cx="240" cy="198" r="17" fill="none" stroke="#4a4486" strokeWidth="7" />
    <circle cx="240" cy="198" r="17" fill="none" stroke="#ffd18a" strokeWidth="7" strokeLinecap="round" strokeDasharray="70 107" transform="rotate(-90 240 198)" />
    <path d="M104 240h196l14 10a3 3 0 0 1-2 5H92a3 3 0 0 1-2-5z" fill="url(#masuk-laptop)" />
    </g>

    <g filter="url(#masuk-shadow)">
    <path d="M346 226h4a9 9 0 0 1 0 18h-4" fill="none" stroke="#e5484d" strokeWidth="5" />
    <rect x="320" y="218" width="28" height="34" rx="8" fill="url(#masuk-mug)" />
    </g>
    <path d="M328 208c-4-5 4-8 0-14M338 208c-4-5 4-8 0-14" fill="none" stroke="#fff" strokeWidth="3" strokeLinecap="round" opacity="0.75" />

    <g>
    <g filter="url(#masuk-shadow)">
    <rect x="44" y="60" width="104" height="48" rx="16" fill="#fff" />
    <path d="M70 106l-6 14 18-12z" fill="#fff" />
    </g>
    <circle cx="76" cy="84" r="5" fill="#b98df3" />
    <circle cx="96" cy="84" r="5" fill="#d7bcfb" />
    <circle cx="116" cy="84" r="5" fill="#ff9aa0" />
    </g>
    <g>
    <circle cx="322" cy="92" r="28" fill="#fff" filter="url(#masuk-shadow)" />
    <path d="M309 92l9 9 17-18" fill="none" stroke="#3aa872" strokeWidth="6" strokeLinecap="round" strokeLinejoin="round" />
    </g>
    <circle cx="200" cy="56" r="5" fill="#fff" opacity="0.8" />
    <circle cx="362" cy="170" r="4" fill="#fff" opacity="0.7" />
    <circle cx="28" cy="250" r="3" fill="#fff" opacity="0.7" />
    </svg>
  )
}
