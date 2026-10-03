// Satu medan layar realisasi, menurut definisinya (`medan.ts`).
//
// Angka uang ditampilkan berpasangan dengan kode mata uangnya (AC 85), dalam
// format sel Section Pega (K14, `sajian.ts`): hanya-baca selalu berformat;
// tersunting berformat selama tidak difokus bila sel ber-
// `pyShowReadonlyFormatting=true`, dan mentah saat diketik. Nilai tersimpan
// tidak pernah diubah oleh format (AC 24).

import { useState } from 'react'

import { Area, Field, Pilih, type Opsi } from '../../../../inti/frontend/components/ui/dasar'
import { nilai, type Halaman } from '../api'
import type { Medan } from '../medan'
import { sajikan, type Sajian } from '../sajian'

export interface PropsKotakMedan {
  medan: Medan
  halaman: Halaman
  wajib: boolean
  hanyaBaca: boolean
  opsiMataUang: Opsi[]
  opsiMO: Opsi[]
  onUbah: (jalur: string, v: string) => void
  onSelesai: (medan: Medan, v: string) => void
}

/** Sajian sel: dari definisinya; medan tanggal tanpa sajian tetap tanggal (AC 33). */
function sajianMedan(m: Medan): Sajian | undefined {
  return m.sajian ?? (m.jenis === 'tanggal' ? 'tanggal' : undefined)
}

export default function KotakMedan({ medan, halaman, wajib, hanyaBaca, opsiMataUang, opsiMO, onUbah, onSelesai }: PropsKotakMedan) {
  const [fokus, setFokus] = useState(false)
  const v = nilai(halaman, medan.jalur)
  const kunci = hanyaBaca || medan.kunci === true || medan.jenis === 'tampil'
  const kode = medan.mataUang ? nilai(halaman, medan.mataUang) : ''
  const pesan = halaman.pesan?.[medan.jalur]?.join('; ')
  const sajian = sajianMedan(medan)

  if (kunci && medan.jenis !== 'centang') {
    return (
      <div className="field nbti__medan">
        <span className="field__label">
          {medan.label}
          {wajib && <span className="field__req">*</span>}
        </span>
        <span className="nbti__nilai" data-jalur={medan.jalur}>
          {kode && v !== '' && <span className="nbti__kode">{kode}</span>}
          {sajikan(v, sajian)}
        </span>
        {pesan && <div className="field__error">{pesan}</div>}
      </div>
    )
  }
  switch (medan.jenis) {
    case 'centang':
      // `pyDisabled` / hanya-baca: kotak tampil, tidak dapat diubah.
      return (
        <label className="nbti__centang">
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
      return <Area label={medan.label} value={v} required={wajib} error={pesan} onChange={(x) => onUbah(medan.jalur, x)} />
    case 'mataUang':
    case 'mo':
      return (
        <Pilih
          label={medan.label}
          value={v}
          required={wajib}
          error={pesan}
          opsi={medan.jenis === 'mo' ? opsiMO : opsiMataUang}
          onChange={(x) => onSelesai(medan, x)}
        />
      )
    case 'tanggal':
      return (
        <div onBlur={() => onSelesai(medan, nilai(halaman, medan.jalur))}>
          <Field label={medan.label} type="date" value={v.slice(0, 10)} required={wajib} error={pesan} onChange={(x) => onUbah(medan.jalur, x)} />
        </div>
      )
    default: {
      const berformat = !fokus && sajian !== undefined && sajian !== 'tanggal' && sajian.formatSaatSunting === true
      return (
        <div
          className="nbti__dengan-kode"
          onFocus={() => setFokus(true)}
          onBlur={() => {
            setFokus(false)
            onSelesai(medan, nilai(halaman, medan.jalur))
          }}
        >
          {kode && <span className="nbti__kode">{kode}</span>}
          <Field
            label={medan.label}
            value={berformat ? sajikan(v, sajian) : v}
            required={wajib}
            error={pesan}
            type="text"
            onChange={(x) => onUbah(medan.jalur, x)}
          />
        </div>
      )
    }
  }
}
