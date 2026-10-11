// Rute modul Komite Claim Life - modul TANPA MENU (perintah work owner 09-10-2026, `layar.ts`; pola Komite Claim Prop):
// kasus komite dibuka DI TEMPAT dari tabel komite inbox Claim Life (`onBukaModul`, `MODUL_DIPINJAM` frontend/App.tsx) -
// `bukaKasus` membuka layar kasusnya langsung, Kembali ke inbox Claim Life (`onBeranda`). Tanpa `bukaKasus` modul ini
// tidak merender apa pun (Inbox Komite sendiri dihapus; daftarnya tabel komite inbox Claim Life).

import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import type { HalamanKomite } from './layar'
import KasusKomite from './pages/KasusKomite'
import './komiteclaimlife.css'

export function RuteKomite({ halaman, masuk, bukaKasus, onBeranda }: PropsRute<HalamanKomite>) {
  return (
    // Akar gaya modul: semua aturan `komiteclaimlife.css` diawali `.komiteclaimlife` (`display: contents`).
    <div className="komiteclaimlife">
      {halaman === 'komiteclaimlife-kasus' && bukaKasus && (
        <KasusKomite key={bukaKasus.ketuk} kasusID={bukaKasus.id} peran={masuk.peran} onKembali={() => onBeranda?.()} />
      )}
    </div>
  )
}

/** Rute modul ini untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul. */
export const RUTE_MODUL: RuteModul<HalamanKomite> = RuteKomite
