// Rute modul Komite Claim Non Prop untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul.
//
// TANPA menu sendiri (perintah work owner 09-10-2026, pola Komite Claim Prop): kasus komite dibuka DI TEMPAT dari tabel
// komite inbox Claim Non Prop (`onBukaModul`, `MODUL_DIPINJAM` frontend/App.tsx) - `bukaKasus` membuka layar kasusnya
// langsung, Back / Submit kembali ke inbox Claim Non Prop (`onBeranda`).

import { BahasaUI } from '../../../inti/frontend/components/ui/bahasaUI'
import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import './komiteclaimnonprop.css'
import type { HalamanKCNP } from './menu'
import DaftarKerja from './pages/DaftarKerja'
import KasusKomite from './pages/KasusKomite'

export function RuteKCNP({ halaman, onLihatBerkas, bukaKasus, onBeranda }: PropsRute<HalamanKCNP>) {
  return (
    <BahasaUI.Provider value="en">
      <div className="komiteclaimnonprop">
        {halaman === 'komiteclaimnonprop-daftar' &&
          (bukaKasus ? (
            <KasusKomite
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

export const RUTE_MODUL: RuteModul<HalamanKCNP> = RuteKCNP
