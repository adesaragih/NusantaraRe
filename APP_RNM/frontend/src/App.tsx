import { useState } from 'react'

import KlaimLife from './pages/KlaimLife'
import Masuk from './pages/Masuk'
import RegisterKlaim from './pages/RegisterKlaim'
import { Shell, type Halaman } from './components/Shell'
import { sesi, type Sesi } from './store/sesi'

// App = gerbang sesi + Shell.
//
// ⛔ Sesudah F0.3, NOL halaman dirender di luar Shell. Cabang "sudah masuk"
// yang F0.2 tinggalkan - dua halaman ditumpuk di bawah satu bilah - dibuang.
//
// ⚠️ Halaman `inbox` lahir di F0.4; sampai itu ia menyatakan dirinya belum ada
// DI DALAM Shell, bukan lewat butir menu `BelumTersedia`. Menu tetap dua, dan
// keduanya berbukti korpus.
export default function App() {
  // Sesi dibaca SEKALI saat menyala: memuat ulang (F5) di tab yang sama tetap
  // masuk, sebab `sessionStorage` bertahan selama tab hidup.
  const [masuk, setMasuk] = useState<Sesi | null>(() => sesi.baca())
  const [halaman, setHalaman] = useState<Halaman>('inbox')

  if (masuk === null) {
    return <Masuk onMasuk={setMasuk} />
  }

  return (
    <Shell
      masuk={masuk}
      halaman={halaman}
      onPindah={setHalaman}
      onKeluar={() => {
        setMasuk(null)
      }}
    >
      {halaman === 'register' && <RegisterKlaim />}
      {halaman === 'detail' && <KlaimLife />}
      {halaman === 'inbox' && <KlaimLife />}
    </Shell>
  )
}
