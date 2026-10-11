// 10 baris per halaman di daftar lokasi / cargo tab Spreading dan Deduction (work owner 11-10-2026 "BUATIN 10 PER PAGE, DI
// SPREADING, DEDUCTION JUGA"). Pemotongan di layar saja - data case tetap dimuat utuh (Save / validasi atas seluruh daftar).
// Indeks ASLI dipertahankan: penyunting baris dan kunci buka-tutup memakai posisi di seluruh daftar, bukan di halaman.

import PagerHalaman from './PagerHalaman'

/** Baris lokasi / cargo per halaman. */
export const UKURAN_HALAMAN_LOKASI = 10

/** Halaman dijepit ke 1..jumlah halaman (daftar bisa memendek sesudah muat ulang). */
export function halamanAman(total: number, halaman: number, ukuran = UKURAN_HALAMAN_LOKASI): number {
  return Math.min(Math.max(1, halaman), Math.max(1, Math.ceil(total / ukuran)))
}

/** Isi satu halaman sebagai pasangan [baris, indeks asli]. */
export function potongHalaman<T>(daftar: readonly T[], halaman: number, ukuran = UKURAN_HALAMAN_LOKASI): [T, number][] {
  const awal = (halamanAman(daftar.length, halaman, ukuran) - 1) * ukuran
  return daftar.slice(awal, awal + ukuran).map((x, i): [T, number] => [x, awal + i])
}

/** Pager di atas tabel; tidak tampil bila semua baris muat satu halaman. */
export default function PagerLokasi({ total, halaman, onPindah }: { total: number; halaman: number; onPindah: (h: number) => void }) {
  if (total <= UKURAN_HALAMAN_LOKASI) return null
  return <PagerHalaman total={total} halaman={halamanAman(total, halaman)} ukuran={UKURAN_HALAMAN_LOKASI} onPindah={onPindah} />
}
