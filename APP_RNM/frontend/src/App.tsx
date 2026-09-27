import { useState } from 'react'

import InboxClaimLife from './pages/InboxClaimLife'
import KlaimLife from './pages/KlaimLife'
import Masuk from './pages/Masuk'
import RegisterKlaim from './pages/RegisterKlaim'
import { Shell, type Halaman } from './components/Shell'
import { sesi, type Sesi } from './store/sesi'

// App = gerbang sesi + Shell.
//
// ⛔ NOL halaman dirender di luar Shell, kecuali layar Masuk - dan ada ujinya.
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
      {halaman === 'inbox' && (
        <InboxClaimLife
          peran={masuk.peran}
          // ⚠️ Klik baris membuka layar TAHAP kasus itu. Sampai kelompok A3
          // masing-masing selesai, ia membuka halaman Detail yang ada -
          // keterangannya DI DALAM halaman, bukan sebagai butir menu.
          onBuka={() => {
            setHalaman('detail')
          }}
          onRegister={() => {
            setHalaman('register')
          }}
        />
      )}
      {halaman === 'register' && <RegisterKlaim />}
      {halaman === 'detail' && <KlaimLife />}
    </Shell>
  )
}
