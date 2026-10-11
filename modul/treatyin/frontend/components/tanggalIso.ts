// Tanggal medan kepala, dan aturan `TreatyInSetTreatyYear`.
//
// ---------------------------------------------------------------------
// ⛔ SATU BENTUK KABEL, DAN ITU `DD-MM-YYYY`
// ---------------------------------------------------------------------
// `FieldTanggal` (inti) memakai `DD-MM-YYYY` sebagai bentuk nilainya:
// `keInputTanggal` mengubahnya menjadi `yyyy-mm-dd` untuk `<input type="date">`,
// dan `dariInputTanggal` mengubahnya KEMBALI menjadi `DD-MM-YYYY` saat
// pemakai memilih tanggal.
//
// ⛔ DUA CACAT LAHIR DARI TIDAK MENGHORMATI ITU, dan keduanya DIAM:
//
//  1. Medan `Commencement`/`Termination` dahulu diisi `01/01/2025` — bentuk
//     BACA. `<input type="date">` menolaknya tanpa bersuara, dan medan yang
//     menolak nilainya terlihat persis seperti medan yang memang kosong.
//
//  2. Perbaikan pertamanya memakai `yyyy-mm-dd`, sementara `FieldTanggal`
//     mengembalikan `DD-MM-YYYY`. Akibatnya terlihat di tangkapan layar
//     pemilik proses 6 Oktober 2026: memilih Commencement membuat
//     `Treaty Year` DAN `Termination` KOSONG — sebab kedua turunan menolak
//     bentuk yang mereka terima.
//
// ⭐ Karena itu berkas ini menerima KEDUA bentuk dan selalu MEMANCARKAN yang
// satu: `DD-MM-YYYY`. Fungsi yang benar hanya bila pemanggilnya ingat bentuk
// mana yang dipakai adalah fungsi yang akan salah pada pemanggil kedua.

/** `DD-MM-YYYY` atau `YYYY-MM-DD` → `[tahun, bulan, hari]`; kosong bila bukan. */
function pecah(nilai: string): [string, string, string] | null {
  const s = nilai.trim()
  let m = /^(\d{2})-(\d{2})-(\d{4})$/.exec(s)
  if (m) return [m[3] ?? '', m[2] ?? '', m[1] ?? '']
  m = /^(\d{4})-(\d{2})-(\d{2})$/.exec(s)
  if (m) return [m[1] ?? '', m[2] ?? '', m[3] ?? '']
  return null
}

/**
 * `YYYYMMDD` (bentuk tersimpan) → `DD-MM-YYYY` (bentuk kabel `FieldTanggal`).
 *
 * ⛔ Yang BUKAN delapan angka menjadi teks KOSONG, bukan diteruskan apa
 * adanya: medan tanggal menolak teks sembarang dan menampilkannya sebagai
 * kosong juga — meneruskannya hanya memindahkan kebingungan ke tempat lain.
 */
export function keKabel(yyyymmdd: string): string {
  const s = yyyymmdd.trim()
  if (!/^\d{8}$/.test(s)) return ''
  return `${s.slice(6, 8)}-${s.slice(4, 6)}-${s.slice(0, 4)}`
}

/**
 * `TreatyIn.TreatyYear := @substring(TreatyIn.Commencement, 0, 4)`
 *
 * ⭐ DISALIN dari `DataTransform/TreatyInSetTreatyYear.xml` langkah 1, yang
 * berjalan pada event `change` milik `TreatyIn.Commencement`.
 */
export function tahunDariMulai(nilai: string): string {
  const p = pecah(nilai)
  return p === null ? '' : p[0]
}

/**
 * `TreatyIn.Termination := @addCalendar(TreatyIn.Commencement, "1", 0,…)`
 *
 * ⭐ DISALIN dari langkah 2 transform yang sama: tambah SATU TAHUN.
 * Memancarkan bentuk kabel, apa pun bentuk masukannya.
 *
 * ⚠️ 29 Februari: `2024-02-29` menjadi `29-02-2025`, tanggal yang tidak ada.
 * Pega `addCalendar` memberi 28 Februari. Bedanya SATU hari pada tahun
 * kabisat, dan ia dinyatakan di sini alih-alih ditambal diam-diam — nol dari
 * 1.854 kontrak mulai 29 Februari, dan menambal yang nol kejadian menambah
 * kode yang tidak pernah diuji.
 */
export function akhirSetahunSesudah(nilai: string): string {
  const p = pecah(nilai)
  if (p === null) return ''
  const tahun = Number(p[0])
  if (!Number.isFinite(tahun)) return ''
  return `${p[2]}-${p[1]}-${String(tahun + 1)}`
}

/**
 * Bentuk kabel (`DD-MM-YYYY`) atau `YYYY-MM-DD` → `YYYYMMDD`, bentuk
 * tersimpan — yang RD spreading bandingkan dengan `TREATYYEAR.STARTDATE`.
 * Yang tak terbaca menjadi teks KOSONG.
 */
export function keSimpan(nilai: string): string {
  const p = pecah(nilai)
  return p === null ? '' : `${p[0]}${p[1]}${p[2]}`
}

/**
 * `YYYYMMDD` (bentuk tersimpan) → `dd/mm/yy`, bentuk BACA grid — sama dengan
 * `TanggalTampil` services (`warisan_daftar.go`), yang dipakai baris yang
 * datang dari kontrak. Baris di penampung halaman (`halaman.tsx`) menyimpan
 * bentuk tersimpan, jadi terjemahan tampilnya dikerjakan di sini.
 * Yang bukan delapan angka tampil APA ADANYA.
 */
export function tanggalTampil(yyyymmdd: string): string {
  const s = yyyymmdd.trim()
  if (!/^\d{8}$/.test(s)) return yyyymmdd
  return `${s.slice(6, 8)}/${s.slice(4, 6)}/${s.slice(2, 4)}`
}
