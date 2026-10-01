// Pencarian master untuk tiga autocomplete Pega (`pxAutoComplete`) - SATU tempat:
//
//   `REINSURER NAME` / `SECURITY REINSURER NAME` - RD `BrowseCedingCoLife_RD` (`InputSecurityLifeReinsurers.xml`
//       b3840, `InputSecurityReinsurerLife.xml` b3922): tampil `.ClientName`, isi `REINSURERID ← .ID`;
//   `BUSINESS NAME` - RD `BrowseBusinessLife_RD` (`InputBusinessLifeReinsurers.xml` b4022): tampil `.Note`,
//       isi `BIZCODE ← .ID`, kolom tambahan `.OLDID`;
//   `R/I RATE` - RD `BrowseRateLifeSummary` (b4471): tampil `.USEDBY`, isi `RIRATEID ← .ID` - ⚠️ OQ-MCRL-13.
//
// Setiap ketikan dicari ke server (`PilihSaring.onCari`) karena server memotong jawaban. Hanya
// jawaban TERAKHIR yang dipakai - sukses MAUPUN gagal - dan `reset` (form lain dibuka) membatalkan
// jawaban yang masih di jalan. Pola dari `modul/treatycontractout` (ditiru, tidak diimpor).

import { useCallback, useRef, useState } from 'react'

import type { OpsiSaring } from '../../../../inti/frontend/components/ui/pilihSaring'
import {
  cariMasterBusiness,
  cariMasterReinsurer,
  cariRingkasanRate,
  type Daftar,
  type MasterBusiness,
  type MasterReinsurer,
} from '../api'

/** Butir `REINSURER NAME`: nama, ID sebagai keterangan (nama kembar dapat dibedakan). */
export function opsiReinsurer(p: MasterReinsurer): OpsiSaring {
  return { value: p.id, label: p.clientName.trim() === '' ? p.id : p.clientName, keterangan: p.id }
}

/** Butir `BUSINESS NAME`: `.Note`, keterangan `.OLDID` (kolom tambahan RD). */
export function opsiBusiness(b: MasterBusiness): OpsiSaring {
  return { value: b.id, label: b.note.trim() === '' ? b.id : b.note, keterangan: b.oldId }
}

/** Butir `R/I RATE`: `.USEDBY`, isi `.ID`. */
export function opsiRate(r: { id: string; usedBy: string }): OpsiSaring {
  return { value: r.id, label: r.usedBy.trim() === '' ? r.id : r.usedBy, keterangan: r.id }
}

function useCari<T>(cariKe: (kata: string) => Promise<Daftar<T>>, keOpsi: (t: T) => OpsiSaring, onGalat: (e: unknown) => void) {
  const [pilihan, setPilihan] = useState<OpsiSaring[]>([])
  const [memuat, setMemuat] = useState(false)
  const urutan = useRef(0)

  const cari = useCallback(
    (kata: string) => {
      const ke = ++urutan.current
      setMemuat(true)
      cariKe(kata)
        .then((d) => {
          if (ke !== urutan.current) return
          setPilihan(d.daftar.map(keOpsi))
          setMemuat(false)
        })
        .catch((e: unknown) => {
          if (ke !== urutan.current) return
          setPilihan([])
          setMemuat(false)
          onGalat(e)
        })
    },
    [cariKe, keOpsi, onGalat],
  )

  const reset = useCallback(() => {
    urutan.current++
    setPilihan([])
    setMemuat(false)
  }, [])

  return { pilihan, memuat, cari, reset }
}

export function useCariReinsurer(onGalat: (e: unknown) => void) {
  return useCari(cariMasterReinsurer, opsiReinsurer, onGalat)
}

export function useCariBusiness(onGalat: (e: unknown) => void) {
  return useCari(cariMasterBusiness, opsiBusiness, onGalat)
}

export function useCariRate(onGalat: (e: unknown) => void) {
  return useCari(cariRingkasanRate, opsiRate, onGalat)
}
