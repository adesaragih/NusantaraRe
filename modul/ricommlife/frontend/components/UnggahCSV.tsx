// Popup "Upload CSV" (b3189 -> local action `UploadCSV_RICOMM` b3258): memilih berkas CSV. Isinya dibaca di peramban
// dan disimpan halaman; "View Upload" dan "Simpan Upload" mengirimnya ke server (server selalu mengurai ulang).

import { useState } from 'react'

import { Modal } from '../../../../inti/frontend/components/ui/dasar'
import { berkasSah } from '../aturan'
import { RC } from '../labels'

export interface BerkasCSV {
  nama: string
  ukuran: number
  isi: string
}

export default function UnggahCSV({ onTutup, onDipilih }: { onTutup: () => void; onDipilih: (b: BerkasCSV) => void }) {
  const [berkas, setBerkas] = useState<File | null>(null)
  const [pesan, setPesan] = useState<string | null>(null)
  const [sibuk, setSibuk] = useState(false)

  const kirim = () => {
    if (sibuk) return
    if (berkas === null || !berkasSah(berkas.name, berkas.size)) {
      setPesan(RC.galatBerkas)
      return
    }
    setSibuk(true)
    berkas.text().then(
      (isi) => onDipilih({ nama: berkas.name, ukuran: berkas.size, isi }),
      () => {
        setSibuk(false)
        setPesan(RC.galatBerkas)
      },
    )
  }

  return (
    <Modal
      judul={RC.unggah}
      onTutup={onTutup}
      onKirim={kirim}
      labelBatal={RC.batal}
      aksi={
        <button type="submit" className="btn btn--primary" disabled={sibuk}>
          {RC.unggah}
        </button>
      }
    >
      <div className="ricommlife__form">
        {pesan !== null && (
          <div className="alert alert--warn" role="alert">
            {pesan}
          </div>
        )}
        <label className="field">
          <span className="field__label">{RC.pilihBerkas}</span>
          <input
            className="field__input"
            type="file"
            accept=".csv,text/csv"
            onChange={(e) => {
              setPesan(null)
              setBerkas(e.target.files?.[0] ?? null)
            }}
          />
        </label>
        <p className="muted ricommlife__catatan">
          <strong>{RC.format}</strong>
        </p>
        <p className="muted ricommlife__catatan">{RC.catatanUnggah}</p>
      </div>
    </Modal>
  )
}
