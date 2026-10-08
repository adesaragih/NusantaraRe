// Tautan yang membuka satu berkas modul di tab baru (keputusan work owner 08-10-2026: tombol View polis Claim Prop
// membuka berkas NB / EDM Treaty In di tab baru). Pengirim memakai form GET tersembunyi `target="_blank"` dengan
// medan `PARAM_MODUL` dan `PARAM_KASUS` ke jalur halaman aplikasi (bukan /api/ - penjaga navigasi
// `claimlife/frontend/unduhdokumen.test.ts` melarang window.open dan href dinamis). App membacanya sekali saat dimuat
// lalu membuka berkasnya lewat jalur kotak masuk Beranda (`PropsRute.bukaKasus`); modul yang menunya tidak dipegang
// akun tidak dipasang - App menampilkan `PESAN_TAUTAN_TANPA_MODUL`.

export const PARAM_MODUL = 'modul'
export const PARAM_KASUS = 'kasus'

export interface TautanKasus {
  modul: string
  id: string
}

/** Tautan berkas dari `location.search`; `null` bila tidak lengkap. */
export function bacaTautanKasus(search: string): TautanKasus | null {
  const p = new URLSearchParams(search)
  const modul = (p.get(PARAM_MODUL) ?? '').trim()
  const id = (p.get(PARAM_KASUS) ?? '').trim()
  return modul !== '' && id !== '' ? { modul, id } : null
}

export const PESAN_TAUTAN_TANPA_MODUL =
  'This file cannot be opened here: the menu of its module is not assigned to your account.'
