// Pembantu ANGKA layar Treaty In.
//
// ⛔ DIPINDAHKAN dari `pages/FormKontrakTreatyIn.tsx` 5 Oktober 2026 —
// pemindahan MURNI, nol perubahan perilaku. Halaman mengekspornya kembali
// supaya uji yang mengimpornya dari sana tetap berjalan tanpa disunting.

import { formatNumber, formatPersen } from '../../../../inti/frontend/lib/format'
import type { BarisLayerWarisan } from '../api'
import {
  DESIMAL_PERSEN,
  DESIMAL_PERSEN_SHARE,
  DESIMAL_UANG,
  desimalPadan,
  golongan,
  type JenisAngka,
} from '../labels'

/**
 * Isi satu sel angka — pemformat `inti` DIPANGGIL, tidak ditulis ulang.
 *
 * ⛔ Pemformatan HANYA di lapis tampilan. Nilai tersimpan tetap teks apa
 * adanya; `repository` dan `services` tidak menyentuhnya.
 */
export function selAngka(jenis: JenisAngka, nilai: string): string {
  return padankanDesimal(selAngkaDasar(jenis, nilai), desimalPadan(jenis))
}

/**
 * ⭐ PEMADANAN §24 — di lapis modul, SESUDAH `formatNumber`.
 *
 * ⛔ Ini BUKAN pemformat kedua, dan bedanya penting: ia tidak menguraikan
 * angka, tidak membulatkan, tidak mengelompokkan ribuan, dan tidak menyentuh
 * apa pun di depan koma. Ia hanya MENAMBAHKAN nol di ekor sampai panjang
 * yang kolomnya minta — pekerjaan yang `formatNumber` sengaja tidak lakukan
 * sejak ronde 69, dan yang §24 minta dikembalikan untuk layar ini saja.
 *
 * ⚠️ `desimal === null` berarti kolomnya tidak terbaca di gambar mana pun:
 * hasilnya lewat APA ADANYA, yaitu aturan ronde 69.
 *
 * ⚠️ Teks yang BUKAN angka lewat apa adanya pula. `>=30% up to < 50%` tidak
 * punya ekor desimal untuk dipadankan, dan memadankannya akan merusaknya.
 *
 * ⛔ Yang SUDAH lebih panjang TIDAK dipotong. Memotong berarti membulatkan,
 * dan membulatkan menghilangkan digit berarti — larangan tertua di jalur
 * angka ini.
 */
export function padankanDesimal(teks: string, desimal: number | null): string {
  if (desimal === null || teks === '') return teks
  // Tanda `%` dilepas lalu dipasang kembali — `formatPersen` menempelkannya
  // di ujung, dan nol harus masuk SEBELUM tanda itu.
  const persen = teks.endsWith('%')
  const inti = persen ? teks.slice(0, -1) : teks
  // Hanya angka bergaya Indonesia yang dipadankan; selain itu apa adanya.
  if (!/^-?\d{1,3}(\.\d{3})*(,\d+)?$/.test(inti)) return teks
  const koma = inti.indexOf(',')
  const adaSekarang = koma < 0 ? 0 : inti.length - koma - 1
  if (adaSekarang >= desimal) return teks
  const tambahan = '0'.repeat(desimal - adaSekarang)
  const hasil = (koma < 0 ? inti + ',' : inti) + tambahan
  return persen ? hasil + '%' : hasil
}

export function selAngkaDasar(jenis: JenisAngka, nilai: string): string {
  switch (golongan(jenis)) {
    case 'uang':
      return formatNumber(nilai, DESIMAL_UANG)
    case 'persen':
      if (!angkaMurni(nilai)) return nilai
      return formatPersen(nilai, DESIMAL_PERSEN)
    case 'persenShare':
      // ⛔ DELAPAN desimal — keputusan pemilik proses 4 Oktober 2026,
      // `KEPUTUSAN-PENYELARASAN-REPO.md` §13. Sama persis dengan batas
      // penyimpanan `NUMBER(38,8)`: menampilkan lebih berarti mengaku lebih
      // teliti daripada yang sistem simpan, menampilkan kurang menutupi
      // selisih yang orang cari ketika memeriksa.
      if (!angkaMurni(nilai)) return nilai
      return formatPersen(nilai, DESIMAL_PERSEN_SHARE)
    default:
      return nilai
  }
}

/**
 * ⛔ Nilai yang BUKAN angka tidak boleh diberi tanda `%`.
 *
 * `formatPersen` mengembalikan teks bukan-angka apa adanya lalu MENEMPELKAN
 * `%` padanya — benar untuk nilai kosong, salah untuk pita seperti
 * `>=30% up to < 50%`, yang menjadi `>=30% up to < 50%%` dengan dua tanda.
 *
 * ⚠️ Ini BUKAN pemformat kedua: ia hanya memutuskan APAKAH `formatPersen`
 * dipanggil. `format.ts` tidak disentuh — aturannya dipenuhi lewat argumen
 * dan lewat pemanggilan.
 */
export function angkaMurni(nilai: string): boolean {
  const t = nilai.trim()
  if (t === '') return false
  return /^[+-]?\d+(\.\d+)?$/.test(t)
}

/**
 * Baris grid dari larik LAYER — satu proyeksi per tab.
 *
 * ⚠️ `JENIS_*` diterapkan lewat `barisAngka`, jadi satu tempat yang
 * memutuskan kolom mana diformat sebagai apa.
 */
export function barisLayer(
  layer: readonly BarisLayerWarisan[],
  proyeksi: (b: BarisLayerWarisan) => readonly string[],
  jenis: readonly JenisAngka[],
): readonly (readonly string[])[] {
  return barisAngka(layer.map(proyeksi), jenis)
}

/**
 * Terapkan golongan angka pada tiap sel.
 *
 * ⛔ Panjang `jenis` HARUS sama dengan panjang tiap baris. Grid yang
 * larik golongannya lebih pendek akan diam-diam membiarkan kolom terakhirnya
 * tanpa format — cacat yang empat grid pernah alami sampai 4 Oktober 2026.
 */
export function barisAngka(
  baris: readonly (readonly string[])[],
  jenis: readonly JenisAngka[],
): readonly (readonly string[])[] {
  return baris.map((b) => b.map((v, i) => selAngka(jenis[i] ?? 'teks', v)))
}
