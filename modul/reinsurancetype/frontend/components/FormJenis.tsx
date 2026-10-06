// Popup Add / Edit Reinsurance Type (tanpa View: daftar sudah cukup). Type 1-4 (Prompt value Pega), Flag active / inactive, Group
// Type OR / QS / RI / SPL. Nilai warisan di luar pilihan (Flag `1`, kosong) tetap tampil sebagai nilai lama supaya baris
// lama dapat disimpan tanpa mengubahnya. Nama ditulis huruf besar oleh backend; ID dibuat backend.

import { useState } from 'react'

import { Gagal, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { tambah, ubah, type Isian, type Jenis } from '../api'
import { FLAG, GROUP_TYPE, isianDari, opsiDengan, periksaIsian, rapikanIsian, TYPE } from '../aturan'
import { RT } from '../labels'

export type ModeForm = { jenis: 'tambah' } | { jenis: 'ubah'; baris: Jenis }

export default function FormJenis({
  mode,
  onTutup,
  onTersimpan,
}: {
  mode: ModeForm
  onTutup: () => void
  onTersimpan: (j: Jenis) => void
}) {
  const baris = mode.jenis === 'tambah' ? null : mode.baris
  const [isi, setIsi] = useState<Isian>(isianDari(baris))
  const [galat, setGalat] = useState<unknown>(null)
  const [pesan, setPesan] = useState<string | null>(null)
  const [sibuk, setSibuk] = useState(false)

  const judul = mode.jenis === 'tambah' ? RT.judulTambah : RT.judulUbah(mode.baris.id)

  const kirim = () => {
    if (sibuk) return
    const dikirim = rapikanIsian(isi)
    const salah = periksaIsian(dikirim, baris)
    setPesan(salah)
    setGalat(null)
    if (salah !== null) return
    setSibuk(true)
    const janji = mode.jenis === 'ubah' ? ubah(mode.baris.id, dikirim) : tambah(dikirim)
    janji.then(
      (j) => {
        setSibuk(false)
        onTersimpan(j)
      },
      (g: unknown) => {
        setSibuk(false)
        setGalat(g)
      },
    )
  }

  const ubahIsian = (s: Partial<Isian>) => setIsi((i) => ({ ...i, ...s }))
  const flagKosongLama = baris !== null && baris.flag === ''

  return (
    <Modal
      judul={judul}
      onTutup={onTutup}
      onKirim={kirim}
      labelBatal={RT.batal}
      lebar
      aksi={
        <button type="submit" className="btn btn--primary" disabled={sibuk}>
          {sibuk ? RT.menyimpan : RT.save}
        </button>
      }
    >
      <div className="reinsurancetype__form">
        {galat !== null && <Gagal galat={galat} />}
        {pesan !== null && (
          <div className="alert alert--warn" role="alert">
            {pesan}
          </div>
        )}
        <div className="reinsurancetype__baris">
          <label className="field">
            <span className="field__label">{RT.id}</span>
            <input
              className="field__input field__input--readonly"
              value={baris?.id ?? ''}
              placeholder={RT.idOtomatis}
              readOnly
            />
          </label>
          <label className="field">
            <span className="field__label">
              {RT.nama} <span aria-hidden="true">*</span>
            </span>
            <input
              className="field__input"
              value={isi.name}
              maxLength={100}
              required
              autoFocus
              onChange={(e) => ubahIsian({ name: e.target.value.toUpperCase() })}
            />
          </label>
        </div>
        <div className="reinsurancetype__baris">
          <label className="field">
            <span className="field__label">
              {RT.type} <span aria-hidden="true">*</span>
            </span>
            <select className="field__input" value={isi.type} onChange={(e) => ubahIsian({ type: e.target.value })}>
              <option value="">{RT.pilih}</option>
              {opsiDengan(TYPE, baris?.type ?? '', RT.labelType).map((o) => (
                <option key={o.value} value={o.value}>
                  {o.label}
                </option>
              ))}
            </select>
          </label>
          <label className="field">
            <span className="field__label">
              {RT.flag} <span aria-hidden="true">*</span>
            </span>
            <select className="field__input" value={isi.flag} onChange={(e) => ubahIsian({ flag: e.target.value })}>
              {flagKosongLama && <option value="">{RT.nilaiLama(RT.kosongNilai)}</option>}
              {opsiDengan(FLAG, baris?.flag ?? '').map((o) => (
                <option key={o.value} value={o.value}>
                  {o.label}
                </option>
              ))}
            </select>
          </label>
        </div>
        <label className="field">
          <span className="field__label">{RT.soa}</span>
          <input
            className="field__input"
            value={isi.soaName}
            maxLength={100}
            onChange={(e) => ubahIsian({ soaName: e.target.value.toUpperCase() })}
          />
        </label>
        {/* Code tidak tampil (perintah work owner 05-10-2026: "Code HAPUS"): nilai tersimpan ikut terkirim apa adanya,
            baris baru `00` (isianDari). */}
        <label className="field">
          <span className="field__label">{RT.noUrut}</span>
          <input
            className="field__input"
            value={isi.noUrut}
            inputMode="numeric"
            maxLength={10}
            onChange={(e) => ubahIsian({ noUrut: e.target.value })}
          />
        </label>
        <label className="field">
          <span className="field__label">{RT.groupType}</span>
          <select
            className="field__input"
            value={isi.groupType}
            onChange={(e) => ubahIsian({ groupType: e.target.value })}
          >
            <option value="">{RT.tanpaGroup}</option>
            {opsiDengan(GROUP_TYPE, baris?.groupType ?? '').map((o) => (
              <option key={o.value} value={o.value}>
                {o.label}
              </option>
            ))}
          </select>
        </label>
        <p className="muted reinsurancetype__catatan">{RT.catatan}</p>
      </div>
    </Modal>
  )
}
