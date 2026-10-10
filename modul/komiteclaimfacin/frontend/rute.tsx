// Rute modul Komite Claim Fac In untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul.
//
// TANPA menu (perintah work owner 09-10-2026, `layar.ts`; pola Komite Claim Prop / Non Prop): kasus komite dibuka DI
// TEMPAT dari tabel komite inbox Claim Fac In (`onBukaModul`, `MODUL_DIPINJAM` frontend/App.tsx) - `bukaKasus` membuka
// layar kasusnya langsung, Back / Submit kembali ke inbox Claim Fac In (`onBeranda`). Tanpa `bukaKasus` modul ini tidak
// merender apa pun (tidak ada halaman daftar kerja sendiri).

import { useState } from 'react'

import { BahasaUI } from '../../../inti/frontend/components/ui/bahasaUI'
import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import './komiteclaimfacin.css'
import { KCFI } from './labels'
import type { HalamanKCFI } from './layar'
import KasusKomite from './pages/KasusKomite'

/** Kasus yang dibuka lewat `bukaKasus`: info sesudah Submit (`HasilKeputusan.info`) tampil dulu sebelum kembali. */
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
      <section className="inbox komiteclaimfacin__akar">
        {pesan.map((p) => (
          <div key={p} className="alert alert--info">
            {p}
          </div>
        ))}
        <div className="komiteclaimfacin__aksi">
          <button type="button" className="btn btn--ghost" onClick={onKembali}>
            {KCFI.kembali}
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

export function RuteKCFI({ halaman, onLihatBerkas, bukaKasus, onBeranda }: PropsRute<HalamanKCFI>) {
  return (
    <BahasaUI.Provider value="en">
      <div className="komiteclaimfacin">
        {halaman === 'komiteclaimfacin-kasus' && bukaKasus && (
          <KasusTerbuka
            key={bukaKasus.ketuk}
            id={bukaKasus.id}
            onLihatBerkas={onLihatBerkas}
            onKembali={() => onBeranda?.()}
          />
        )}
      </div>
    </BahasaUI.Provider>
  )
}

export const RUTE_MODUL: RuteModul<HalamanKCFI> = RuteKCFI
