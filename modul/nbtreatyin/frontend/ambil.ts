// Efek ambil-dengan-batal - SATU tempat untuk pola "muat data saat kunci
// berubah, buang hasil permintaan yang sudah usang" (komponen dilepas atau
// kuncinya berubah sebelum jawaban tiba).

import { useEffect, useState, type DependencyList } from 'react'

/** Jalankan `ambil`; hasil/galat diteruskan kecuali sudah dibatalkan. Kembalian = pembatal. */
export function ambilBatal<T>(ambil: () => Promise<T>, berhasil: (x: T) => void, gagal: (e: unknown) => void): () => void {
  let dibuang = false
  ambil().then(
    (x) => {
      if (!dibuang) berhasil(x)
    },
    (e: unknown) => {
      if (!dibuang) gagal(e)
    },
  )
  return () => {
    dibuang = true
  }
}

/** Efek `ambilBatal` setiap `kunci` berubah. */
export function useAmbilBatal<T>(
  ambil: () => Promise<T>,
  berhasil: (x: T) => void,
  gagal: (e: unknown) => void,
  kunci: DependencyList,
): void {
  // `kunci` adalah daftar ketergantungan yang ditentukan pemanggil.
  // eslint-disable-next-line react-hooks/exhaustive-deps
  useEffect(() => ambilBatal(ambil, berhasil, gagal), kunci)
}

/** Data yang dimuat ulang setiap `kunci` berubah: `null` selama memuat. */
export function useAmbil<T>(ambil: () => Promise<T>, kunci: DependencyList): { data: T | null; galat: unknown } {
  const [data, setData] = useState<T | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  useAmbilBatal(
    () => {
      setData(null)
      setGalat(null)
      return ambil()
    },
    setData,
    setGalat,
    kunci,
  )
  return { data, galat }
}
