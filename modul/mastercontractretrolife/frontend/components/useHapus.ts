// Alur hapus satu panel - popup lebih dulu (penyimpangan sadar 3), SATU tempat untuk keempat jenis.
//
// `minta` menghitung dampak SESAAT sebelum popup tampil; `ya` mengirim dampak yang dilihat. Galat
// sesudah `Yes` (mis. 409 karena data berubah) tampil DI popup dan angkanya dimuat ulang - pengguna
// melihat keadaan terbaru sebelum mencoba lagi.

import { useState } from 'react'

import { ambilDampakHapus, hapus, type Dampak, type JenisHapus } from '../api'

export interface KeadaanKonfirmasi {
  id: string
  nama: string
  dampak: Dampak | null
  galat: unknown
}

export function useHapus(jenis: JenisHapus, onTerhapus: (id: string, pesan: string) => Promise<void>) {
  const [konfirmasi, setKonfirmasi] = useState<KeadaanKonfirmasi | null>(null)
  const [sibuk, setSibuk] = useState(false)

  function muatDampak(id: string): void {
    ambilDampakHapus(jenis, id)
      .then((j) => {
        setKonfirmasi((k) => (k === null || k.id !== id ? k : { ...k, dampak: j.dampak }))
      })
      .catch((e: unknown) => {
        setKonfirmasi((k) => (k === null || k.id !== id ? k : { ...k, galat: e }))
      })
  }

  function minta(id: string, nama: string): void {
    setKonfirmasi({ id, nama, dampak: null, galat: null })
    muatDampak(id)
  }

  async function ya(): Promise<void> {
    if (konfirmasi === null || konfirmasi.dampak === null || sibuk) return
    const { id, dampak } = konfirmasi
    setSibuk(true)
    try {
      const hasil = await hapus(jenis, id, dampak)
      setKonfirmasi(null)
      await onTerhapus(id, hasil.pesan)
    } catch (e) {
      setKonfirmasi((k) => (k === null ? k : { ...k, galat: e, dampak: null }))
      muatDampak(id)
    } finally {
      setSibuk(false)
    }
  }

  return { konfirmasi, sibuk, minta, ya, batal: () => setKonfirmasi(null) }
}
