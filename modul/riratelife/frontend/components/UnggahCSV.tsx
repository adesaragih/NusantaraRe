// Popup "Upload CSV" (b2949 -> local action `UploadCSV_RIRATE` b3019): memilih berkas CSV. Isinya dibaca di peramban
// dan disimpan halaman; "View Upload" dan "Simpan Upload" mengirimnya ke server (server selalu mengurai ulang).

import { useState } from 'react'

import { Modal } from '../../../../inti/frontend/components/ui/dasar'
import { berkasSah } from '../aturan'
import { RR } from '../labels'

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
      setPesan(RR.galatBerkas)
      return
    }
    setSibuk(true)
    berkas.text().then(
      (isi) => onDipilih({ nama: berkas.name, ukuran: berkas.size, isi }),
      () => {
        setSibuk(false)
        setPesan(RR.galatBerkas)
      },
    )
  }

  return (
    <Modal
      judul={RR.unggah}
      onTutup={onTutup}
      onKirim={kirim}
      labelBatal={RR.batal}
      aksi={
        <button type="submit" className="btn btn--primary" disabled={sibuk}>
          {RR.unggah}
        </button>
      }
    >
      <div className="riratelife__form">
        {pesan !== null && (
          <div className="alert alert--warn" role="alert">
            {pesan}
          </div>
        )}
        <label className="field">
          <span className="field__label">{RR.pilihBerkas}</span>
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
        <p className="muted riratelife__catatan">
          <strong>{RR.format}</strong>
        </p>
        <p className="muted riratelife__catatan">{RR.catatanUnggah}</p>
      </div>
    </Modal>
  )
}
