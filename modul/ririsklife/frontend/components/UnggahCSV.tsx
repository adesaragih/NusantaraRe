// Popup "Upload CSV" (`InboxSummaryRIRisk` b3156, local action): memilih berkas CSV. Isinya dibaca di peramban
// dan disimpan halaman; "View Upload" dan "Simpan Upload" mengirimnya ke server (server selalu mengurai ulang).

import { useState } from 'react'

import { Modal } from '../../../../inti/frontend/components/ui/dasar'
import { berkasSah } from '../aturan'
import { RK } from '../labels'

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
      setPesan(RK.galatBerkas)
      return
    }
    setSibuk(true)
    berkas.text().then(
      (isi) => onDipilih({ nama: berkas.name, ukuran: berkas.size, isi }),
      () => {
        setSibuk(false)
        setPesan(RK.galatBerkas)
      },
    )
  }

  return (
    <Modal
      judul={RK.unggah}
      onTutup={onTutup}
      onKirim={kirim}
      labelBatal={RK.batal}
      aksi={
        <button type="submit" className="btn btn--primary" disabled={sibuk}>
          {RK.unggah}
        </button>
      }
    >
      <div className="ririsklife__form">
        {pesan !== null && (
          <div className="alert alert--warn" role="alert">
            {pesan}
          </div>
        )}
        <label className="field">
          <span className="field__label">{RK.pilihBerkas}</span>
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
        {/* Label `Format excel : ...` b5495 tersembunyi di Pega (`1=2` b5263) - tidak dirender; kepala berkas dijelaskan catatan ini. */}
        <p className="muted ririsklife__catatan">{RK.catatanUnggah}</p>
      </div>
    </Modal>
  )
}
