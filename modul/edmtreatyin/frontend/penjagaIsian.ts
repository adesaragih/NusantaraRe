// Penjaga isian layar kasus (tinjauan kode 06-10-2026): setiap jawaban server (`POST /hitung`, `PUT`, pilih bisnis)
// MENGGANTI seluruh halaman kerja. Tanpa penjaga, isian yang diketik selama permintaan berjalan hilang, dan jawaban
// yang tiba tak berurutan menimpa jawaban yang lebih baru.
//
//   - permintaan dijalankan BERURUTAN; masing-masing dikirim dengan halaman TERBARU saat benar-benar berangkat
//     (sudah memuat jawaban sebelumnya dan isian sesudahnya);
//   - isian yang diubah SESUDAH permintaan berangkat dicatat dan diterapkan ulang di atas jawaban server.

import type { Halaman } from './api'

export type UbahHalaman = (h: Halaman) => Halaman

export interface PenjagaIsian {
  /** Halaman terbaru: jawaban server terakhir + isian pengguna sesudahnya. */
  kini(): Halaman | null
  /** Halaman dari luar antrean (muat awal). Membuang catatan isian. */
  pasang(h: Halaman): void
  /** Isian pengguna; dicatat bila sebuah permintaan sedang berjalan. */
  ubah(f: UbahHalaman): Halaman | null
  /**
   * Antrekan satu permintaan. `minta` menerima halaman terbaru saat berangkat; `halamanDari` mengambil halaman
   * jawaban (null = jawaban tanpa halaman) yang dipasang - beserta isian susulan - SEBELUM permintaan berikut berangkat.
   */
  kirim<T>(minta: (h: Halaman) => Promise<T>, halamanDari: (x: T) => Halaman | null): Promise<T>
}

export function buatPenjagaIsian(awal: Halaman | null): PenjagaIsian {
  let h = awal
  let berjalan = false
  let susulan: UbahHalaman[] = []
  let antrean: Promise<unknown> = Promise.resolve()

  const kirim = <T>(minta: (h: Halaman) => Promise<T>, halamanDari: (x: T) => Halaman | null): Promise<T> => {
    const langkah = antrean.then(async () => {
      if (h === null) throw new Error('halaman kerja belum dimuat')
      susulan = []
      berjalan = true
      try {
        const x = await minta(h)
        const baru = halamanDari(x)
        if (baru !== null) {
          let nh = baru
          for (const f of susulan) nh = f(nh)
          h = nh
        }
        return x
      } finally {
        berjalan = false
        susulan = []
      }
    })
    antrean = langkah.catch(() => undefined)
    return langkah
  }

  return {
    kini: () => h,
    pasang: (x) => {
      h = x
      susulan = []
    },
    ubah: (f) => {
      if (h === null) return h
      h = f(h)
      if (berjalan) susulan.push(f)
      return h
    },
    kirim,
  }
}
