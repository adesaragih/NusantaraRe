/**
 * Keadaan sesi — satu-satunya keadaan yang benar-benar GLOBAL.
 *
 * Isinya hanya "siapa yang sedang masuk, dan memegang peran apa". Data server
 * bukan urusan berkas ini, dan keadaan layar milik layarnya sendiri.
 *
 * # Kenapa `sessionStorage`, bukan `localStorage`
 *
 * Pola referensi (`REFERENSI_UI/frontend/src/store/sesi.ts`). Identitas di
 * `localStorage` bertahan melewati penutupan peramban, sehingga sesi yang
 * ditinggalkan tetap hidup pada MESIN BERSAMA — dan mesin bersama persis
 * keadaan di kantor. `sessionStorage` memendekkan umurnya menjadi seumur tab.
 *
 * Yang akan terasa: menutup tab berarti keluar; tab kedua menuntut masuk lagi;
 * memuat ulang (F5) di tab yang sama tetap masuk.
 *
 * # ⛔ Nol sandi
 *
 * Tidak ada sandi yang diminta, dikirim, maupun disimpan — selama `AUTH_STUB`
 * menyala, identitas hanyalah PENGAKUAN yang backend percaya lewat header.
 * Itu sah untuk pengembangan dan TIDAK PERNAH sah untuk produksi; gerbangnya
 * ada di `bolehMasukStub()` di bawah, dan di `stubPelaku` sisi Go.
 *
 * ⚠️ `sessionStorage` dapat MELEMPAR (mode privat, site-data diblokir) dan
 * dapat kembali kosong. Setiap sentuhan dibungkus try/catch: sesi yang tidak
 * dapat disimpan harus berarti "belum masuk", bukan layar yang pecah.
 */

import { PERAN, type KodePeran } from '../assets/labels'

/** Kunci penyimpanan — berawalan nama aplikasi supaya tidak bertabrakan. */
const KUNCI = 'rnm.sesi'

/** Identitas yang sedang masuk. */
export interface Sesi {
  /** Dikirim sebagai `X-Pelaku`. */
  akunID: string
  /**
   * Dikirim sebagai `X-Peran`, DIPISAH KOMA.
   *
   * `[terverifikasi]` `internal/handlers/pelaku.go` `pelakuDari`: header
   * dipecah pada koma, tiap potongan di-trim, potongan kosong dibuang. Jadi
   * satu pelaku MEMANG boleh memegang lebih dari satu peran, dan layar tidak
   * boleh memaksanya memilih satu.
   */
  peran: KodePeran[]
}

/** Ketiga peran yang ada — ADR-U-0002. */
export const PERAN_TERSEDIA: readonly KodePeran[] = [
  PERAN.admin,
  PERAN.medis,
  PERAN.spv,
]

/**
 * Apakah masuk lewat stub diizinkan.
 *
 * ⛔ Gerbang ini adalah satu-satunya hal yang memisahkan "pilih peran apa saja
 * dari layar" dari sistem sungguhan. Ia dibaca dari env saat membangun;
 * `VITE_AUTH_STUB` yang tidak disetel berarti TIDAK BOLEH — gagal tertutup.
 */
export function bolehMasukStub(): boolean {
  return import.meta.env.VITE_AUTH_STUB === 'true'
}

function bacaMentah(): string | null {
  try {
    return sessionStorage.getItem(KUNCI)
  } catch {
    // Penyimpanan diblokir — sama artinya dengan belum masuk.
    return null
  }
}

/** Hanya peran yang dikenal yang diterima kembali dari penyimpanan. */
function saringPeran(nilai: unknown): KodePeran[] {
  if (!Array.isArray(nilai)) return []
  const sah: KodePeran[] = []
  for (const p of nilai) {
    if (typeof p === 'string' && (PERAN_TERSEDIA as readonly string[]).includes(p)) {
      // ⛔ Disaring, tidak dipercaya apa adanya: isi sessionStorage dapat
      // disunting siapa pun yang membuka devtools. Backend tetap gerbang
      // sebenarnya, tetapi layar pun tidak boleh mengarang peran.
      sah.push(p as KodePeran)
    }
  }
  return sah
}

export const sesi = {
  /** Sesi yang tersimpan, atau `null` bila belum masuk. */
  baca(): Sesi | null {
    const mentah = bacaMentah()
    if (mentah === null || mentah === '') return null
    try {
      const isi: unknown = JSON.parse(mentah)
      if (typeof isi !== 'object' || isi === null) return null
      const o = isi as Record<string, unknown>
      const akunID = typeof o.akunID === 'string' ? o.akunID.trim() : ''
      const peran = saringPeran(o.peran)
      // ⛔ Sesi tanpa akun atau tanpa peran BUKAN sesi. Mengembalikannya
      // membuat layar mengira ada yang masuk, lalu backend menolak tiap
      // permintaan dengan 401 yang tidak dapat dijelaskan pemakai.
      if (akunID === '' || peran.length === 0) return null
      return { akunID, peran }
    } catch {
      return null
    }
  },

  /** Menyimpan sesi. Mengembalikan `false` bila penyimpanan menolak. */
  simpan(s: Sesi): boolean {
    try {
      sessionStorage.setItem(KUNCI, JSON.stringify(s))
      return true
    } catch {
      return false
    }
  },

  /** Keluar — identitasnya dibuang. */
  hapus(): void {
    try {
      sessionStorage.removeItem(KUNCI)
    } catch {
      // Tidak ada yang dapat dilakukan; sesi seumur tab akan hilang sendiri.
    }
  },
}

/**
 * Header identitas untuk satu permintaan.
 *
 * ⚠️ Mengembalikan objek KOSONG bila belum masuk - bukan header bernilai
 * kosong. Header `X-Pelaku: ` yang kosong dan header yang tidak ada sama
 * artinya bagi `pelakuDari`, tetapi yang kedua lebih jujur di panel jaringan:
 * ia memperlihatkan bahwa permintaan itu memang anonim.
 */
export function headerIdentitas(): Record<string, string> {
  const s = sesi.baca()
  if (s === null) return {}
  return {
    'X-Pelaku': s.akunID,
    // Koma, sesuai `pelakuDari`.
    'X-Peran': s.peran.join(','),
  }
}
