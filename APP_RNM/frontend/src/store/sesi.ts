/**
 * Identitas pelaku — satu-satunya keadaan yang benar-benar GLOBAL.
 *
 * # ⛔ TIDAK ADA FORM LOGIN `[perintah work owner 27-09-2026]`
 *
 * Form masuk dibuang. Selama `VITE_AUTH_STUB=true`, identitas dibentuk SAAT
 * APLIKASI MENYALA dari env Vite - bukan dari layar, bukan dari
 * `sessionStorage`. Alasannya sederhana: layar masuk tanpa sandi bukan
 * autentikasi, dan mempertahankannya hanya menambah langkah yang tidak
 * memutuskan apa pun.
 *
 * ⛔ GAGAL TERTUTUP. `VITE_AUTH_STUB` yang bukan `'true'` - termasuk TIDAK
 * DISETEL - berarti TIDAK ADA identitas: header tidak dikirim, dan backend
 * menolak setiap jalur beridentitas. Itu keadaan yang benar sampai IAM ada
 * (tiket 07 / ADR-U-0030).
 *
 * ⛔ Nol sandi diminta, dikirim, atau disimpan - tidak ada tempatnya lagi.
 *
 * ⚠️ Bentuk `Sesi` DIPERTAHANKAN. Kelak IAM mengisinya; yang berubah hanya
 * SUMBERnya, bukan bentuk yang dibaca seluruh layar.
 */

import { PERAN, type KodePeran } from '../assets/labels'

/** Identitas yang sedang dipakai. */
export interface Sesi {
  /** Dikirim sebagai `X-Pelaku`. */
  akunID: string
  /**
   * Dikirim sebagai `X-Peran`, DIPISAH KOMA.
   *
   * `[terverifikasi]` `internal/handlers/pelaku.go` `pelakuDari`: header
   * dipecah pada koma, tiap potongan di-trim, potongan kosong dibuang. Jadi
   * satu pelaku MEMANG boleh memegang lebih dari satu peran.
   */
  peran: KodePeran[]
}

/** Ketiga peran yang ada — ADR-U-0002. */
export const PERAN_TERSEDIA: readonly KodePeran[] = [
  PERAN.admin,
  PERAN.medis,
  PERAN.spv,
]

/** Akun bawaan bila `VITE_STUB_PELAKU` tidak disetel. */
const AKUN_BAWAAN = 'UJI-ADMIN'

/**
 * Apakah identitas stub diizinkan.
 *
 * ⛔ Gerbang ini satu-satunya yang memisahkan "identitas dari env" dari
 * sistem sungguhan, dan ia gagal tertutup.
 */
export function bolehMasukStub(): boolean {
  return import.meta.env.VITE_AUTH_STUB === 'true'
}

/** Hanya peran yang dikenal yang diterima. */
function saringPeran(nilai: readonly string[]): KodePeran[] {
  const sah: KodePeran[] = []
  for (const p of nilai) {
    const bersih = p.trim()
    if (bersih !== '' && (PERAN_TERSEDIA as readonly string[]).includes(bersih)) {
      // ⛔ Disaring, tidak dipercaya apa adanya: env dapat salah ketik, dan
      // peran yang dikarang akan ditolak backend dengan 403 yang tidak dapat
      // dijelaskan pemakai.
      sah.push(bersih as KodePeran)
    }
  }
  return sah
}

/**
 * Identitas pelaku, atau `null` bila stub-nya mati.
 *
 * ⚠️ Dibaca dari env pada SETIAP panggilan, bukan disimpan di modul: nilai
 * env dibekukan saat Vite membangun, jadi membacanya ulang tidak mahal - dan
 * fungsi tanpa keadaan tersembunyi jauh lebih mudah diuji.
 */
export function pelakuStub(): Sesi | null {
  if (!bolehMasukStub()) return null

  const akun = (import.meta.env.VITE_STUB_PELAKU ?? '').trim()
  const daftar = (import.meta.env.VITE_STUB_PERAN ?? '').split(',')
  const peran = saringPeran(daftar)

  return {
    akunID: akun === '' ? AKUN_BAWAAN : akun,
    // ⚠️ Bawaannya KETIGA peran, bukan satu: pengembang yang menyalakan stub
    // hampir selalu ingin melihat seluruh antrian. Menyempitkannya ke satu
    // peran membuat tiga dari empat tab menghilang tanpa sebab yang terlihat.
    peran: peran.length === 0 ? [...PERAN_TERSEDIA] : peran,
  }
}

/**
 * Header identitas untuk satu permintaan.
 *
 * ⚠️ Objek KOSONG bila stub mati - bukan header bernilai kosong. Keduanya
 * sama artinya bagi `pelakuDari`, tetapi yang pertama lebih jujur di panel
 * jaringan: ia memperlihatkan bahwa permintaan itu memang anonim.
 */
export function headerIdentitas(): Record<string, string> {
  const s = pelakuStub()
  if (s === null) return {}
  return {
    'X-Pelaku': s.akunID,
    // Koma, sesuai `pelakuDari`.
    'X-Peran': s.peran.join(','),
  }
}
