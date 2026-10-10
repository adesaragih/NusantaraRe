// Rute modul Komite Claim Prop untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul.
//
// Sejak 09-10-2026 Komite Claim Prop TANPA menu sendiri (keputusan work owner): kasus komite dibuka DI TEMPAT dari tabel
// komite inbox Claim Prop (`onBukaModul`, `MODUL_DIPINJAM` frontend/App.tsx; "jangan pop up, langsung buka komitenya") -
// `bukaKasus` membuka layar kasusnya langsung, Back / Submit kembali ke inbox Claim Prop (`onBeranda`).

import { useState } from 'react'

import { BahasaUI } from '../../../inti/frontend/components/ui/bahasaUI'
import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import './komiteclaimprop.css'
import { KCP } from './labels'
import type { HalamanKCP } from './menu'
import DaftarKerja from './pages/DaftarKerja'
import KasusKomite from './pages/KasusKomite'

/** Kasus yang dibuka lewat `bukaKasus`: pesan sesudah Submit (PDF akseptasi gagal) tampil dulu sebelum kembali. */
function KasusTerbuka({
  id,
  onLihatBerkas,
  onKembali,
}: {
  id: string
  onLihatBerkas?: (modul: string, id: string) => boolean
  onKembali: () => void
}) {
  const [pesan, setPesan] = useState<string[]>([])
  if (pesan.length > 0) {
    return (
      <section className="inbox komiteclaimprop__akar">
        {pesan.map((p) => (
          <div key={p} className="alert alert--error">
            {p}
          </div>
        ))}
        <div className="komiteclaimprop__aksi">
          <button type="button" className="btn btn--ghost" onClick={onKembali}>
            {KCP.kembali}
          </button>
        </div>
      </section>
    )
  }
  return (
    <KasusKomite
      id={id}
      onLihatBerkas={onLihatBerkas}
      onKembali={(p) => {
        if (p && p.length > 0) setPesan(p)
        else onKembali()
      }}
    />
  )
}

export function RuteKCP({ halaman, onLihatBerkas, bukaKasus, onBeranda }: PropsRute<HalamanKCP>) {
  return (
    <BahasaUI.Provider value="en">
      <div className="komiteclaimprop">
        {halaman === 'komiteclaimprop-daftar' &&
          (bukaKasus ? (
            <KasusTerbuka
              key={bukaKasus.ketuk}
              id={bukaKasus.id}
              onLihatBerkas={onLihatBerkas}
              onKembali={() => onBeranda?.()}
            />
          ) : (
            <DaftarKerja onLihatBerkas={onLihatBerkas} />
          ))}
      </div>
    </BahasaUI.Provider>
  )
}

export const RUTE_MODUL: RuteModul<HalamanKCP> = RuteKCP
