// Popup Add / Edit Treaty Group (tanpa View: perintah work owner 05-10-2026, daftar sudah cukup). OJK Business dari
// dropdown TREATYGROUPOJK; COA tidak diketik - tampil dan berganti begitu OJK dipilih ("ubah pas pilih OJK Business"),
// disimpan backend dengan aturan yang sama. Order No tidak tampil ("hide aja dari tampilan") - backend tetap
// menyalinnya dari OJK.

import { useState } from 'react'

import { Gagal, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { tambah, ubah, type Detail, type Grup, type Isian, type Ojk } from '../api'
import { coaTampil, isianDari, opsiOjk, periksaIsian, rapikanIsian } from '../aturan'
import { TG } from '../labels'

export type ModeForm = { jenis: 'tambah' } | { jenis: 'ubah'; baris: Grup }

export default function FormGrup({
  mode,
  ojk,
  onTutup,
  onTersimpan,
}: {
  mode: ModeForm
  ojk: Ojk[]
  onTutup: () => void
  onTersimpan: (d: Detail) => void
}) {
  const baris = mode.jenis === 'tambah' ? null : mode.baris
  const [isi, setIsi] = useState<Isian>(isianDari(baris))
  const [galat, setGalat] = useState<unknown>(null)
  const [pesan, setPesan] = useState<string | null>(null)
  const [sibuk, setSibuk] = useState(false)

  const judul = mode.jenis === 'tambah' ? TG.judulTambah : TG.judulUbah(mode.baris.id)

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
      (d) => {
        setSibuk(false)
        onTersimpan(d)
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
      labelBatal={TG.batal}
      aksi={
        <button type="submit" className="btn btn--primary" disabled={sibuk}>
          {sibuk ? TG.menyimpan : TG.save}
        </button>
      }
    >
      <div className="treatygroup__form">
        {galat !== null && <Gagal galat={galat} />}
        {pesan !== null && (
          <div className="alert alert--warn" role="alert">
            {pesan}
          </div>
        )}
        <label className="field">
          <span className="field__label">{TG.id}</span>
          <input
            className="field__input field__input--readonly"
            value={baris?.id ?? ''}
            placeholder={TG.idOtomatis}
            readOnly
          />
        </label>
        <label className="field">
          <span className="field__label">
            {TG.ojk} <span aria-hidden="true">*</span>
          </span>
          <select
            className="field__input"
            value={isi.ojkId}
            required
            autoFocus
            onChange={(e) => ubahIsian({ ojkId: e.target.value })}
          >
            <option value="">{TG.pilih}</option>
            {opsiOjk(ojk, baris?.ojkId ?? '', baris?.ojkName ?? '').map((o) => (
              <option key={o.value} value={o.value}>
                {o.label}
              </option>
            ))}
          </select>
        </label>
        <label className="field">
          <span className="field__label">
            {TG.nama} <span aria-hidden="true">*</span>
          </span>
          <input
            className="field__input"
            value={isi.name}
            maxLength={1000}
            required
            onChange={(e) => ubahIsian({ name: e.target.value.toUpperCase() })}
          />
        </label>
        <label className="field">
          <span className="field__label">{TG.soa}</span>
          <input
            className="field__input"
            value={isi.soaName}
            maxLength={1000}
            onChange={(e) => ubahIsian({ soaName: e.target.value.toUpperCase() })}
          />
        </label>
        <label className="field">
          <span className="field__label">{TG.coa}</span>
          <input
            className="field__input field__input--readonly"
            value={coaTampil(ojk, isi.ojkId, baris)}
            placeholder={isi.ojkId === '' ? TG.coaOtomatis : TG.coaKosong}
            readOnly
          />
        </label>
        <p className="muted treatygroup__catatan">{TG.catatan}</p>
      </div>
    </Modal>
  )
}
