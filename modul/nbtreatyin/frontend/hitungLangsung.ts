// Hitung langsung saat mengetik - work owner 10-10-2026 ("udah di ubah premi OGP, tapi deduction A tidak berubah
// otomatis, harus di triger dulu"). Isian angka ber-aksi dulu hanya menghitung saat fokus ditinggalkan (blur); kini
// juga JEDA_HITUNG_MS sesudah ketikan terakhir. Dua pengaman:
//
//   - nilai kosong / nol TIDAK dihitung saat mengetik: `CountOGPONP_Act` langkah 1-2 mengenolkan persen dan hasil OGP /
//     ONP bila premi kosong atau "0" - menghapus angka lalu berhenti sejenak tidak boleh menghapus % Deduction. Nilai
//     kosong / nol tetap dihitung saat blur, seperti sebelumnya.
//   - jawaban server tidak menimpa ketikan yang lebih baru: medan yang berubah sesudah permintaan dikirim
//     dipertahankan (`pertahankanKetikan`), dan jawaban permintaan yang sudah disusul diabaikan (layar).

import type { Halaman } from './api'
import { nilaiNol } from './sajian'

/** Jeda sesudah ketikan terakhir sebelum hitung dikirim. */
export const JEDA_HITUNG_MS = 600

/** Ketikan ini dihitung tanpa menunggu blur: medan punya aksi dan nilainya bukan kosong / nol. */
export function hitungSaatKetik(punyaAksi: boolean, nilai: string): boolean {
  return punyaAksi && !nilaiNol(nilai)
}

/**
 * Halaman hasil server, dengan medan yang diubah pengguna SESUDAH permintaan dikirim (`kini` berbeda dari `dikirim`)
 * dipertahankan - hitung susulan untuk ketikan itu menyusul lewat jedanya sendiri.
 */
export function pertahankanKetikan(server: Halaman, dikirim: Halaman, kini: Halaman): Halaman {
  const nilai = { ...server.nilai }
  let ada = false
  for (const [jalur, v] of Object.entries(kini.nilai)) {
    if (v !== dikirim.nilai[jalur]) {
      nilai[jalur] = v
      ada = true
    }
  }
  return ada ? { ...server, nilai } : server
}
