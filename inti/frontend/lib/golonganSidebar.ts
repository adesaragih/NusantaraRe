// Golongan sidebar yang DITUTUP pemakai - permintaan work owner 05-10-2026: "untuk yang GROUP nya buat bisa buka
// tutup". Yang dilipat hanya kepala golongan (TREATY, FACULTATIVE, KLAIM, MASTER, ...); di bawahnya tetap SATU tombol
// per modul tanpa anak (menu datar 30-09-2026).
//
// Disimpan di localStorage peramban ini: kenyamanan per pemakai, bukan data. Penyimpanan yang tidak terbaca
// (mode privat, diblokir, isi rusak) = semua golongan terbuka - bawaan yang sama dengan sebelum fitur ini.

/** Kunci localStorage - daftar KODE golongan yang ditutup (JSON). */
export const KUNCI_GOLONGAN_TERTUTUP = 'rnm.sidebar.golonganTertutup'

type Penyimpanan = Pick<Storage, 'getItem' | 'setItem'>

function penyimpanan(): Penyimpanan | null {
  try {
    return typeof window === 'undefined' ? null : window.localStorage
  } catch {
    return null
  }
}

/** Golongan yang ditutup menurut penyimpanan; kosong bila tidak terbaca. */
export function bacaGolonganTertutup(simpan: Penyimpanan | null = penyimpanan()): ReadonlySet<string> {
  try {
    const isi: unknown = JSON.parse(simpan?.getItem(KUNCI_GOLONGAN_TERTUTUP) ?? '[]')
    return new Set(Array.isArray(isi) ? isi.filter((k): k is string => typeof k === 'string') : [])
  } catch {
    return new Set()
  }
}

/** Simpan golongan yang ditutup; gagal menulis diabaikan (pilihan berlaku sampai halaman dimuat ulang). */
export function simpanGolonganTertutup(
  tertutup: ReadonlySet<string>,
  simpan: Penyimpanan | null = penyimpanan(),
): void {
  try {
    simpan?.setItem(KUNCI_GOLONGAN_TERTUTUP, JSON.stringify([...tertutup].sort()))
  } catch {
    // Penyimpanan penuh atau diblokir.
  }
}

/** Klik kepala golongan: tertutup jadi terbuka, terbuka jadi tertutup. */
export function alihkanGolongan(tertutup: ReadonlySet<string>, kode: string): ReadonlySet<string> {
  const baru = new Set(tertutup)
  if (!baru.delete(kode)) baru.add(kode)
  return baru
}

/** Golongan halaman yang sedang tampil dibuka; himpunan yang SAMA bila ia sudah terbuka (tanpa render ulang). */
export function bukaGolongan(tertutup: ReadonlySet<string>, kode: string): ReadonlySet<string> {
  if (!tertutup.has(kode)) return tertutup
  const baru = new Set(tertutup)
  baru.delete(kode)
  return baru
}
