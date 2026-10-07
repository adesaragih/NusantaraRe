// Asal: salinan `modul/nbtreatyin/frontend/components/KotakMedan.tsx` (06-10-2026), kelas `nbti__` -> `edmt__`.
// Disesuaikan EDM: dropdown hanya-baca ber-RD (EDM Type, Currency, Marketing Officer) menampilkan teks pilihan
// acuannya (`Medan.acuan`); jenis medan yang tidak dipakai layar EDM dibuang.
//
// Satu medan layar endorsemen, menurut definisinya (`medan.ts`). Hanya-baca selalu berformat; isian angka
// (`InputAngka`) hanya angka, ribuan otomatis, rata kanan, nol = placeholder 0. Nilai tersimpan tidak pernah
// diubah oleh format.

import { useId } from 'react'

import { Area, Pilih, type Opsi } from '../../../../inti/frontend/components/ui/dasar'
import { nilai, type Halaman } from '../api'
import { teksPilihan, type Medan, type SumberAcuan } from '../medan'
import { nilaiNol, sajikan } from '../sajian'
import InputAngka from './InputAngka'

export interface PropsKotakMedan {
  medan: Medan
  halaman: Halaman
  wajib: boolean
  hanyaBaca: boolean
  /** Daftar acuan (`GET /acuan`): dropdown Marketing Officer dan teks dropdown hanya-baca. */
  opsi: Record<SumberAcuan, Opsi[]>
  onUbah: (jalur: string, v: string) => void
  onSelesai: (medan: Medan, v: string) => void
}

export default function KotakMedan({ medan, halaman, wajib, hanyaBaca, opsi, onUbah, onSelesai }: PropsKotakMedan) {
  const v = nilai(halaman, medan.jalur)
  const kunci = hanyaBaca || medan.kunci === true || medan.jenis === 'tampil'
  const pesan = halaman.pesan?.[medan.jalur]?.join('; ')
  const sajian = medan.sajian
  const idIsian = useId()
  const daftarAcuan = medan.acuan ?? (medan.jenis === 'mo' ? 'mo' : undefined)

  if (kunci && medan.jenis !== 'centang') {
    // angka rata kanan; nol di bagian uang tampil "0" redup (sama dengan placeholder isian)
    const angka = sajian !== undefined && sajian !== 'tanggal' && !medan.opsi
    const nol = angka && sajian.nolPolos === true && nilaiNol(v)
    const teks = medan.opsi
      ? teksPilihan(medan.opsi, v)
      : daftarAcuan
        ? teksPilihan(opsi[daftarAcuan], v)
        : sajikan(v, sajian)
    return (
      <div className="field edmt__medan">
        <span className="field__label">
          {medan.label}
          {wajib && <span className="field__req">*</span>}
        </span>
        <span
          className={'edmt__nilai' + (angka ? ' edmt__nilai--angka' : '') + (nol ? ' edmt__nilai--nol' : '')}
          data-jalur={medan.jalur}
        >
          {teks}
        </span>
        {pesan && <div className="field__error">{pesan}</div>}
      </div>
    )
  }
  switch (medan.jenis) {
    case 'centang':
      // hanya-baca: kotak tampil, tidak dapat diubah
      return (
        <label className="edmt__centang">
          <input
            type="checkbox"
            checked={v === 'true'}
            disabled={kunci}
            onChange={(e) => onSelesai(medan, e.target.checked ? 'true' : 'false')}
          />
          {medan.label}
        </label>
      )
    case 'area':
      return (
        <Area label={medan.label} value={v} required={wajib} error={pesan} onChange={(x) => onUbah(medan.jalur, x)} />
      )
    case 'radio': {
      // pxRadioButtons: nilai tersimpan di luar daftar ditambahkan sebagai pilihan sendiri (tidak dibuang)
      const pilihan = medan.opsi ?? []
      const semua = v !== '' && !pilihan.some((o) => o.value === v) ? [...pilihan, { value: v, label: v }] : pilihan
      return (
        <div className="field">
          <span className="field__label">
            {medan.label}
            {wajib && <span className="field__req">*</span>}
          </span>
          <div className="edmt__radio" role="radiogroup" aria-label={medan.label}>
            {semua.map((o) => (
              <label key={o.value}>
                <input
                  type="radio"
                  name={`edmt-${medan.jalur}`}
                  checked={v === o.value}
                  onChange={() => onSelesai(medan, o.value)}
                />
                {o.label}
              </label>
            ))}
          </div>
          {pesan && <div className="field__error">{pesan}</div>}
        </div>
      )
    }
    case 'mo':
      return (
        <Pilih
          label={medan.label}
          value={v}
          required={wajib}
          error={pesan}
          opsi={opsi.mo}
          onChange={(x) => onSelesai(medan, x)}
        />
      )
    case 'angka':
      // isian angka (permintaan work owner 06-10-2026): hanya angka, ribuan otomatis, rata kanan, placeholder "0"
      return (
        <div className="edmt__dengan-kode">
          <div className="field">
            <label className="field__label" htmlFor={idIsian}>
              {medan.label}
              {wajib && <span className="field__req">*</span>}
            </label>
            <InputAngka
              id={idIsian}
              label={medan.label}
              value={v}
              sajian={sajian ?? {}}
              onChange={(x) => onUbah(medan.jalur, x)}
              onBlur={() => onSelesai(medan, nilai(halaman, medan.jalur))}
            />
            {pesan && <div className="field__error">{pesan}</div>}
          </div>
        </div>
      )
  }
  return null
}
