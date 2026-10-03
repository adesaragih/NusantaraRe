// Satu medan layar realisasi, menurut definisinya (`medan.ts`).
//
// ⛔ Angka uang ditampilkan APA ADANYA, berpasangan dengan kode mata uangnya
// (AC 85). Format penyajian - jumlah desimal, pemisah ribuan - belum ditetapkan
// (AC 86, `[terbuka]`), jadi tidak ada format yang ditebak di sini.

import { Area, Field, Pilih, type Opsi } from '../../../../inti/frontend/components/ui/dasar'
import { nilai, type Halaman } from '../api'
import type { Medan } from '../medan'

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

export default function KotakMedan({ medan, halaman, wajib, hanyaBaca, opsiMataUang, opsiMO, onUbah, onSelesai }: PropsKotakMedan) {
  const v = nilai(halaman, medan.jalur)
  const kunci = hanyaBaca || medan.kunci || medan.jenis === 'tampil'
  const kode = medan.mataUang ? nilai(halaman, medan.mataUang) : ''
  const pesan = halaman.pesan?.[medan.jalur]?.join('; ')

  if (kunci && medan.jenis !== 'centang') {
    return (
      <div className="field nbti__medan">
        <span className="field__label">
          {medan.label}
          {wajib && <span className="field__req">*</span>}
        </span>
        <span className="nbti__nilai" data-jalur={medan.jalur}>
          {kode && v !== '' && <span className="nbti__kode">{kode}</span>}
          {v}
        </span>
        {pesan && <div className="field__error">{pesan}</div>}
      </div>
    )
  }
  switch (medan.jenis) {
    case 'centang':
      return (
        <label className="nbti__centang">
          <input
            type="checkbox"
            checked={v === 'true'}
            disabled={kunci && hanyaBaca}
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
    default:
      return (
        <div className="nbti__dengan-kode" onBlur={() => onSelesai(medan, nilai(halaman, medan.jalur))}>
          {kode && <span className="nbti__kode">{kode}</span>}
          <Field
            label={medan.label}
            value={v}
            required={wajib}
            error={pesan}
            type="text"
            onChange={(x) => onUbah(medan.jalur, x)}
          />
        </div>
      )
  }
}
