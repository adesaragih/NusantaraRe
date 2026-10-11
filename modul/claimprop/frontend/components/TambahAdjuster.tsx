// Popup tombol "+" Consultant / Adjuster - padanan harness Pega MstAdjusterConsultant (keputusan work owner
// 08-10-2026). Label VERBATIM section XML; master baru disimpan modul Adjuster Consultant (rute pinjaman
// `POST /api/adjuster-consultant`) yang membuat ID-nya dan menolak nama kembar. Sesudah tersimpan, layar kasus
// mengisi ID baru ke medannya (aksi SetConsultant / SetAdjsuter).

import { useState } from 'react'

import { Gagal, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { tambahAdjuster } from '../api'
import { CP } from '../labels'

export default function TambahAdjuster({
  onTersimpan,
  onTutup,
}: {
  onTersimpan: (idBaru: string) => void
  onTutup: () => void
}) {
  const [nama, setNama] = useState('')
  const [alamat, setAlamat] = useState('')
  const [telp, setTelp] = useState('')
  const [galat, setGalat] = useState<unknown>(null)
  const [sibuk, setSibuk] = useState(false)

  const simpan = () => {
    if (sibuk) return
    setSibuk(true)
    setGalat(null)
    tambahAdjuster({ name: nama.trim(), address: alamat.trim(), telpNo: telp.trim() }).then(
      (a) => {
        setSibuk(false)
        onTersimpan(a.id)
      },
      (g: unknown) => {
        setSibuk(false)
        setGalat(g)
      },
    )
  }

  return (
    <Modal
      judul={CP.popTambahAdjuster}
      onTutup={onTutup}
      onKirim={simpan}
      labelBatal={CP.cancel}
      lebar
      aksi={
        <button type="submit" className="btn btn--primary" disabled={sibuk || nama.trim() === ''}>
          {CP.save}
        </button>
      }
    >
      {galat !== null && <Gagal galat={galat} />}
      <div className="form-grid">
        <label className="field">
          <span className="field__label">
            Adjuster/Consultant Name<span className="field__req">*</span>
          </span>
          <input className="field__input" value={nama} required autoFocus onChange={(e) => setNama(e.target.value)} />
        </label>
        <label className="field">
          <span className="field__label">Telp No Adjuster/Consultant</span>
          <input className="field__input" value={telp} onChange={(e) => setTelp(e.target.value)} />
        </label>
        <label className="field field--lebar">
          <span className="field__label">Address Adjuster/Consultant</span>
          <textarea
            className="field__input claimprop__area"
            value={alamat}
            onChange={(e) => setAlamat(e.target.value)}
          />
        </label>
      </div>
    </Modal>
  )
}
