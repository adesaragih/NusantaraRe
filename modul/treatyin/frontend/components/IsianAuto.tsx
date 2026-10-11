// Autocomplete Pega (`pxAutoComplete`) — ketik bebas, pilihan dari daftar.
//
// ⭐ Pega memetakan DUA properti dari baris yang dipilih (nama → medan,
// `.ID`/`.BizCode` → medan pengenal); nama yang diketik tanpa memilih
// tetap tersimpan dengan pengenal kosong — begitu pula di sini.

import { useCallback, useId } from 'react'

import { Field } from '../../../../inti/frontend/components/ui/dasar'
import type { PilihanWarisan } from '../api'
import DropdownWarisan from './DropdownWarisan'

export default function IsianAuto({
  label,
  nilai,
  pilihan,
  bisaUbah,
  onPilih,
}: {
  label: string
  nilai: string
  pilihan: readonly PilihanWarisan[]
  bisaUbah: boolean
  /** Nama yang diketik/dipilih, dan pengenalnya bila nama itu ada di daftar. */
  onPilih: (nama: string, id: string) => void
}) {
  const idDaftar = useId()
  if (!bisaUbah) return <Field label={label} value={nilai} readOnly onChange={() => undefined} />
  return (
    <div className="field">
      {label !== '' && <label className="field__label">{label}</label>}
      <input
        className="field__input"
        list={idDaftar}
        value={nilai}
        aria-label={label || undefined}
        onChange={(e) => {
          const v = e.target.value
          onPilih(v, pilihan.find((o) => o.nama === v)?.id ?? '')
        }}
      />
      <datalist id={idDaftar}>
        {pilihan.map((o, i) => (
          <option key={`${o.id}-${i}`} value={o.nama} />
        ))}
      </datalist>
    </div>
  )
}

/**
 * Pemilih berbentuk DROPDOWN — pembungkus `DropdownWarisan` atas daftar
 * yang SUDAH ada di tangan.
 *
 * ⭐ Permintaan pemilik proses 7 Oktober 2026: *"perbaiki di limits non
 * prop agar menjadi dropdown seperti dropdown ceding"*. Komponennya SAMA
 * dengan pemilih `Ceding` dan `Source of Business` — bukan tiruannya — jadi
 * ketik-saring, navigasi papan tik, dan penandaan nama kembar ikut apa
 * adanya.
 *
 * ⛔ TANDA TANGANNYA SENGAJA SAMA PERSIS dengan `IsianAuto` di atas,
 * sehingga tiap titik panggil berganti satu kata. Itu penting: ada empat
 * pemilih di tab Limits Non-Prop, dan penggantian yang menuntut penulisan
 * ulang tiap pemanggil mengundang salah satu tertinggal.
 *
 * ⚠️ `DropdownWarisan` menjadikan `ambil` TANGGUNGAN efeknya. Fungsi baru
 * tiap render akan mengambil ulang daftarnya tanpa henti, jadi ia dibekukan
 * `useCallback` di sini — SEKALI, bukan di tiap pemanggil yang mudah lupa.
 */
export function DropdownDaftar({
  label,
  nilai,
  pilihan,
  bisaUbah,
  onPilih,
}: {
  label: string
  nilai: string
  pilihan: readonly PilihanWarisan[]
  bisaUbah: boolean
  onPilih: (nama: string, id: string) => void
}) {
  const ambil = useCallback(() => Promise.resolve([...pilihan]), [pilihan])
  if (!bisaUbah) return <Field label={label} value={nilai} readOnly onChange={() => undefined} />
  return (
    <DropdownWarisan
      label={label}
      ambil={ambil}
      nilai={nilai}
      onPilih={(o) => {
        onPilih(o.nama, o.id)
      }}
    />
  )
}
