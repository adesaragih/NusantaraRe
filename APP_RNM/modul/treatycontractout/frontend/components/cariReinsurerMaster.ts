// Pencarian master reinsurer (`BrowseAgentReinsSOA_RD`) untuk pemilih
// `Reinsurer` (form reinsurer) dan `Security Name` (form security) — SATU
// tempat untuk keduanya [keputusan work owner 30-09-2026: dropdown yang dapat
// difilter, tanpa kotak cari terpisah].
//
// Master aktif lebih dari 100 nama dan server memotong di 100, jadi setiap
// ketikan dicari ke server (`PilihSaring.onCari`). Hanya jawaban TERAKHIR
// yang dipakai — sukses MAUPUN gagal — dan `reset` (form lain dibuka)
// membatalkan jawaban yang masih di jalan.

import { useCallback, useRef, useState } from 'react'

import type { OpsiSaring } from '../../../../inti/frontend/components/ui/pilihSaring'
import { cariReinsurerMaster, type ReinsurerMaster } from '../api'

/** Satu butir pilihan: nama, dengan ID sebagai keterangan (nama kembar dapat dibedakan). */
export function opsiReinsurerMaster(p: ReinsurerMaster): OpsiSaring {
  return { value: p.id, label: p.clientName.trim() === '' ? p.id : p.clientName, keterangan: p.id }
}

export function useCariReinsurerMaster(onGalat: (e: unknown) => void) {
  const [pilihan, setPilihan] = useState<OpsiSaring[]>([])
  const [memuat, setMemuat] = useState(false)
  const urutan = useRef(0)

  const cari = useCallback(
    (kata: string) => {
      const ke = ++urutan.current
      setMemuat(true)
      cariReinsurerMaster(kata)
        .then((d) => {
          if (ke !== urutan.current) return
          setPilihan(d.map(opsiReinsurerMaster))
          setMemuat(false)
        })
        .catch((e: unknown) => {
          if (ke !== urutan.current) return
          setMemuat(false)
          onGalat(e)
        })
    },
    [onGalat],
  )

  const reset = useCallback(() => {
    urutan.current++
    setPilihan([])
    setMemuat(false)
  }, [])

  return { pilihan, memuat, cari, reset }
}
