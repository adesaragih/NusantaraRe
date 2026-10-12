// Rute modul Endorsement Life.
//
// Label korpus berbahasa Inggris; `BahasaUI` "en" membuat teks bawaan komponen bersama (Memuat,
// penomoran halaman, panel backend mati) ikut Inggris di layar ini saja.
//
// ⛔ Satu halaman menu (`edm-inbox`); layar yang terbuka (kotak masuk, Create Addendum, kasus)
// disimpan di keadaan rute, supaya ia bertahan saat pemakai pindah halaman lewat cara lain lalu
// kembali. TETAPI memilih menu Endorsement Life (termasuk saat halamannya sedang tampil) SELALU
// kembali ke kotak masuk - keputusan work owner 07-10-2026 (`ketukMenu`, `menuEDMDipilihUlang`).
//
// ⚠️ Pola sama dengan PremiumList Life, tetapi fungsi pembandingnya milik modul ini sendiri: nol
// impor dari modul lain (perubahan Endorsement tidak boleh berdampak ke PremiumList dan sebaliknya).

import { useEffect, useRef, useState } from 'react'

import { BahasaUI } from '../../../inti/frontend/components/ui/bahasaUI'
import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import { HALAMAN_AWAL_EDM, type HalamanEDM } from './menu'
import BuatEndorsement from './pages/BuatEndorsement'
import InboxEndorsementLife from './pages/InboxEndorsementLife'
import InputEDMLife from './pages/InputEDMLife'

/** Layar yang tampil di halaman menu: kotak masuk, harness buat endorsement, atau kasus. */
export type Layar = { jenis: 'inbox' } | { jenis: 'buat' } | { jenis: 'kasus'; id: string }

/**
 * Menu Endorsement Life baru saja dipilih (lagi) - layar kembali ke kotak masuk.
 *
 * ⛔ `halaman` saja tidak cukup: memilih menu halaman yang SEDANG tampil tidak mengubahnya.
 * `ketukMenu` (`PropsRute`) bertambah setiap pilihan menu; tanpa `ketukMenu` = tidak pernah.
 */
export function menuEDMDipilihUlang(
  ketukLalu: number | undefined,
  ketukBaru: number | undefined,
  halaman: string,
): boolean {
  return ketukBaru !== undefined && ketukBaru !== ketukLalu && halaman === HALAMAN_AWAL_EDM
}

/** Layar sesudah ketukan menu: kotak masuk bila menu dipilih ulang, selain itu layar yang sama. */
export function layarSesudahMenu(
  layar: Layar,
  ketukLalu: number | undefined,
  ketukBaru: number | undefined,
  halaman: string,
): Layar {
  return menuEDMDipilihUlang(ketukLalu, ketukBaru, halaman) ? { jenis: 'inbox' } : layar
}

export function RuteEDM({ halaman, ketukMenu }: PropsRute<HalamanEDM>) {
  const [layar, setLayar] = useState<Layar>({ jenis: 'inbox' })
  // Klik menu (ulang) → kotak masuk (keputusan work owner 07-10-2026). Perpindahan halaman
  // lewat cara lain tidak menyentuh `ketukMenu`, jadi layar tetap bertahan.
  const ketukLalu = useRef(ketukMenu)
  useEffect(() => {
    if (menuEDMDipilihUlang(ketukLalu.current, ketukMenu, halaman)) {
      setLayar({ jenis: 'inbox' })
    }
    ketukLalu.current = ketukMenu
  }, [ketukMenu, halaman])

  if (halaman !== HALAMAN_AWAL_EDM) return null
  const keInbox = () => setLayar({ jenis: 'inbox' })
  const keKasus = (id: string) => setLayar({ jenis: 'kasus', id })
  return (
    <BahasaUI.Provider value="en">
      {layar.jenis === 'inbox' && <InboxEndorsementLife onBuka={keKasus} onBuat={() => setLayar({ jenis: 'buat' })} />}
      {layar.jenis === 'buat' && <BuatEndorsement onTutup={keInbox} onDibuat={keKasus} />}
      {layar.jenis === 'kasus' && <InputEDMLife kasusId={layar.id} onTutup={keInbox} />}
    </BahasaUI.Provider>
  )
}

/** Rute modul ini untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul. */
export const RUTE_MODUL: RuteModul<HalamanEDM> = RuteEDM
