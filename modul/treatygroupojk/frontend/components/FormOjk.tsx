// Popup Add / Edit Treaty Group OJK (tanpa View: daftar sudah memuat seluruh kolom). Name dan Name (IDN) ditulis huruf
// besar oleh backend; ID dibuat backend (nomor tertinggi + 1).

import { useState } from 'react'

import { Gagal, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { tambah, ubah, type Isian, type Ojk } from '../api'
import { isianDari, periksaIsian, rapikanIsian } from '../aturan'
import { TGO } from '../labels'

export type ModeForm = { jenis: 'tambah' } | { jenis: 'ubah'; baris: Ojk }

export default function FormOjk({
  mode,
  onTutup,
  onTersimpan,
}: {
  mode: ModeForm
  onTutup: () => void
  onTersimpan: (o: Ojk) => void
}) {
  const baris = mode.jenis === 'tambah' ? null : mode.baris
  const [isi, setIsi] = useState<Isian>(isianDari(baris))
  const [galat, setGalat] = useState<unknown>(null)
  const [pesan, setPesan] = useState<string | null>(null)
  const [sibuk, setSibuk] = useState(false)

  const judul = mode.jenis === 'tambah' ? TGO.judulTambah : TGO.judulUbah(mode.baris.id)

  const kirim = () => {
    if (sibuk) return
    const dikirim = rapikanIsian(isi)
    const salah = periksaIsian(dikirim)
    setPesan(salah)
    setGalat(null)
    if (salah !== null) return
    setSibuk(true)
    const janji = mode.jenis === 'ubah' ? ubah(mode.baris.id, dikirim) : tambah(dikirim)
    janji.then(
      (o) => {
        setSibuk(false)
        onTersimpan(o)
      },
      (g: unknown) => {
        setSibuk(false)
        setGalat(g)
      },
    )
  }

  const ubahIsian = (s: Partial<Isian>) => setIsi((i) => ({ ...i, ...s }))

  return (
    <Modal
      judul={judul}
      onTutup={onTutup}
      onKirim={kirim}
      labelBatal={TGO.batal}
      aksi={
        <button type="submit" className="btn btn--primary" disabled={sibuk}>
          {sibuk ? TGO.menyimpan : TGO.save}
        </button>
      }
    >
      <div className="treatygroupojk__form">
        {galat !== null && <Gagal galat={galat} />}
        {pesan !== null && (
          <div className="alert alert--warn" role="alert">
            {pesan}
          </div>
        )}
        <label className="field">
          <span className="field__label">{TGO.id}</span>
          <input
            className="field__input field__input--readonly"
            value={baris?.id ?? ''}
            placeholder={TGO.idOtomatis}
            readOnly
          />
        </label>
        <label className="field">
          <span className="field__label">
            {TGO.nama} <span aria-hidden="true">*</span>
          </span>
          <input
            className="field__input"
            value={isi.name}
            maxLength={50}
            required
            autoFocus
            onChange={(e) => ubahIsian({ name: e.target.value.toUpperCase() })}
          />
        </label>
        <label className="field">
          <span className="field__label">
            {TGO.namaIdn} <span aria-hidden="true">*</span>
          </span>
          <input
            className="field__input"
            value={isi.nameIdn}
            maxLength={50}
            required
            onChange={(e) => ubahIsian({ nameIdn: e.target.value.toUpperCase() })}
          />
        </label>
        <p className="muted treatygroupojk__catatan">{TGO.catatanHurufBesar}</p>
      </div>
    </Modal>
  )
}
