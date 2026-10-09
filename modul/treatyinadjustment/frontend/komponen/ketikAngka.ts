// Penyaring ketikan kotak ANGKA berbentuk TEKS — permintaan pemakai
// 9 Oktober 2026: *"jangan di inputan tanggal bisa masukkan huruf begitu
// juga tempat inputan angka, cek keseluruhan"*.
//
// ⭐ Kontrol `pxNumber` sudah memakai `FieldAngka` (inti), yang menyaring
// sendiri. Yang disaring DI SINI adalah medan yang ekspornya `pxTextInput`
// padahal isinya angka (jumlah uang, persen, jumlah hari) — daftarnya
// dihimpun dari kerangka bangkitan, satu per satu, bukan ditebak dari nama.
//
// ⛔ SALINAN `saringAngka` modul Treaty In
// (`modul/treatyin/frontend/components/saringAngka.ts`) — modul tidak boleh
// saling impor. Perubahan perilaku di salah satunya disalin ke yang lain.
//
// ⛔ TIDAK disaring (teks sungguhan): `Comment` (Description deduksi),
// `CoInShare` (`>=30% up to < 50%`), dan `Layer` grid Reinsurer /
// Facultative Reinsurers (`ShareReins`, `ShareFacultativeReinsurers` —
// sama dengan layar Treaty In, placeholder `Text`).

import { segmenAkhir } from '../ekspor/golongan'

/** Pecahan: digit, SATU pemisah desimal (koma → titik), `-` di depan. Bulat: digit saja. */
export function saringAngka(teks: string, bulat = false): string {
  const t = teks.trim()
  if (bulat) {
    const batas = t.search(/[.,]/)
    return (batas >= 0 ? t.slice(0, batas) : t).replace(/\D/g, '')
  }
  const minus = t.startsWith('-') ? '-' : ''
  const bersih = t.replace(/,/g, '.').replace(/[^\d.]/g, '')
  const i = bersih.indexOf('.')
  return minus + (i < 0 ? bersih : bersih.slice(0, i + 1) + bersih.slice(i + 1).replace(/\./g, ''))
}

/** `pxTextInput` berisi angka pecahan (uang, persen) — kunci segmen akhir. */
const PECAHAN: ReadonlySet<string> = new Set([
  'Amount', 'RSMDLimit', 'Earthquake', 'FloodJab', 'FloodNation',
  'Limit', 'Limit2', 'AgregateLimit', 'AgregateLimit2',
  'RNMShare', 'BrokeragePercent', 'FacultativeShare', 'FacultativeShareBrokerage', 'FacShare', 'FacShareBrokerage',
  'RIOGR', 'RIONR', 'PremiumReservePct', 'ProfitCommision', 'ProfitME', 'ProfitYDCF',
  'Layer', 'LayerPart',
  // `Value to IDR` grid Rate of Exchange (`GRID_KURS`, di luar kerangka tab).
  'Conversion',
])
/** `pxTextInput` berisi bilangan bulat (jumlah hari, cacah). */
const BULAT: ReadonlySet<string> = new Set([
  'InstallmentNo', 'WPC', 'ReportingInterval', 'ReportingSubmission', 'ReportingConfirmation', 'ReportingSettlement', 'SubDays',
])
/** Larik yang kolom `Layer`-nya TEKS (grid reasuradur). */
const LARIK_LAYER_TEKS: ReadonlySet<string> = new Set(['ShareReins', 'ShareFacultativeReinsurers'])

/**
 * Penyaring untuk satu medan/sel teks, atau `undefined` bila teks bebas.
 * `larik` = larik grid sel itu (kosong untuk medan).
 */
export function penyaringKetik(kunci: string, larik = ''): ((v: string) => string) | undefined {
  const k = segmenAkhir(kunci)
  if (k === 'Layer' && LARIK_LAYER_TEKS.has(segmenAkhir(larik))) return undefined
  if (BULAT.has(k)) return (v) => saringAngka(v, true)
  if (PECAHAN.has(k)) return (v) => saringAngka(v)
  return undefined
}

/**
 * Teks berisi angka PECAHAN yang ditampilkan BERFORMAT (`FieldAngka`, titik
 * ribuan, dua desimal) — sama dengan kotak angka lain layar ini. Permintaan
 * pemakai 9 Oktober 2026: Amount panel New tab Maximum Retention tampil
 * mentah. `Layer`/`LayerPart` DIKECUALIKAN: nomor urut, bukan uang — tampil
 * `1,00` dilaporkan pemilik proses (layar Treaty In, 7 Oktober 2026).
 */
export function angkaBerformat(kunci: string, larik = ''): boolean {
  const k = segmenAkhir(kunci)
  if (k === 'Layer' || k === 'LayerPart') return false
  return penyaringKetik(kunci, larik) !== undefined && PECAHAN.has(k)
}

/** Ketikan kotak tanggal: digit dan pemisah `/ - .` serta spasi, maksimal 10 aksara. */
export function saringTanggal(teks: string): string {
  return teks.replace(/[^\d/.\- ]/g, '').slice(0, 10)
}
