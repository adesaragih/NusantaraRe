// Dropdown master - keputusan work owner 02-10-2026 ("perubahan pada tampilan untuk semua Choose ubah jadi
// dropdown saja"): menggantikan tombol `Choose*` + popup FlowAction `Choose*` (PARITAS §4) untuk ketujuh master -
// Ceding, SOB, R/I Risk, Cause Of Loss, Policy Holder, Currency, dan R/I Rate baris `PLAN LIST`.
//
// Yang TETAP seperti XML: nilai HANYA dari daftar master - medannya tidak dapat diketik (`pyReadOnly` true b4040,
// b4428, b7362, b10626, b17062, b28105); kata `Search` dihurufbesarkan server (`SearchPolicyHolder_act` 1 b236) dan
// dicocokkan "Contains"; memilih menyalin ID + nama (`set*_DT`); kolom daftar `ID` / `Name` (`RIRate Name` untuk
// R/I Rate); mode lihat (`IsView`) tidak dapat dibuka - nama tampil sebagai teks (di form, baris `Medan` yang
// menampilkannya; dropdown hanya dirender di mode sunting).
// Yang BERUBAH: daftar dibuka dari medannya sendiri, memuat paling banyak BATAS_DROPDOWN baris (master `CLIENT`
// ratusan ribu baris) - potongan dinyatakan, sisanya dicapai lewat `Search` di dalam dropdown.

import { useCallback } from 'react'

import { cariMaster, type JenisMaster, type NilaiMaster } from '../api'
import { BATAS_DROPDOWN, potongPilihan } from '../bentuk'
import { PEMILIH_MPNL } from '../labels'
import DropdownCari from './DropdownCari'

export default function DropdownMaster({
  labelAria,
  jenis,
  nilai,
  kolomNama = PEMILIH_MPNL.kolomName,
  lihat,
  onPilih,
}: {
  /** Label VERBATIM medan / kepala kolom - untuk pembaca layar (label tampilnya milik `Medan` / kepala grid). */
  labelAria: string
  jenis: JenisMaster
  /** Nama yang tampil (`.Ceding`, `.SOBName`, ...). */
  nilai: string
  /** Kepala kolom nama (`Name`; `RIRate Name` untuk R/I Rate). */
  kolomNama?: string
  /** Mode lihat (`IsView == 'true'`): baca-saja, tidak dapat dibuka. */
  lihat: boolean
  /** Penerima `set*_DT`: menyalin ID + nama. */
  onPilih: (v: NilaiMaster) => void
}) {
  // Server diminta satu baris lebih dari yang dirender: baris lebih = daftar terpotong, dinyatakan.
  const cari = useCallback(async (kata: string) => potongPilihan((await cariMaster(jenis, kata, BATAS_DROPDOWN + 1)).daftar), [jenis])
  return (
    <DropdownCari<NilaiMaster>
      labelAria={labelAria}
      nilai={nilai}
      lihat={lihat}
      cari={cari}
      kunci={(v) => v.id}
      kolom={{ judul: [PEMILIH_MPNL.kolomId, kolomNama], isi: (v) => [v.id, v.nama] }}
      terpilih={(v) => v.nama === nilai}
      onPilih={onPilih}
    />
  )
}
