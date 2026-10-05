// Popup Add / Edit kurs tahunan (tanpa View: daftar sudah memuat seluruh kolom). Add / Edit: Currency dari view CURRENCY; Start Date dan End
// Date bawaan 1 Juli - 30 Juni tahun treaty (boleh diubah); Quarter 0 (tahunan) - 4. Backend menulis tanggal format
// Pega dan menolak Treaty Year + Currency + Quarter yang kembar.

import { useState } from 'react'

import { Gagal, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { tambah, ubah, type Isian, type Kurs, type MataUang } from '../api'
import { gantiTahun, isianDari, opsiMataUang, periksaIsian, QUARTER, rapikanIsian } from '../aturan'
import { TEY } from '../labels'

export type ModeForm = { jenis: 'tambah' } | { jenis: 'ubah'; baris: Kurs }

export default function FormKurs({
  mode,
  mataUang,
  onTutup,
  onTersimpan,
}: {
  mode: ModeForm
  mataUang: MataUang[]
  onTutup: () => void
  onTersimpan: (k: Kurs) => void
}) {
  const baris = mode.jenis === 'tambah' ? null : mode.baris
  const judul = mode.jenis === 'tambah' ? TEY.judulTambah : TEY.judulUbah(mode.baris.id)
  const [isi, setIsi] = useState<Isian>(isianDari(baris))
  const [galat, setGalat] = useState<unknown>(null)
  const [pesan, setPesan] = useState<string | null>(null)
  const [sibuk, setSibuk] = useState(false)

  const kirim = () => {
    if (sibuk) return
    const dikirim = rapikanIsian(isi)
    const salah = periksaIsian(dikirim)
    setPesan(salah)
    setGalat(null)
    if (salah !== null) return
    setSibuk(true)
    const janji = mode.jenis === 'ubah' ? ubah(dikirim) : tambah(dikirim)
    janji.then(
      (k) => {
        setSibuk(false)
        onTersimpan(k)
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
      labelBatal={TEY.batal}
      aksi={
        <button type="submit" className="btn btn--primary" disabled={sibuk}>
          {sibuk ? TEY.menyimpan : TEY.save}
        </button>
      }
    >
      <div className="treatyexchangeyearly__form">
        {galat !== null && <Gagal galat={galat} />}
        {pesan !== null && (
          <div className="alert alert--warn" role="alert">
            {pesan}
          </div>
        )}
        <label className="field">
          <span className="field__label">{TEY.id}</span>
          <input
            className="field__input field__input--readonly"
            value={baris?.id ?? ''}
            placeholder={TEY.idOtomatis}
            readOnly
          />
        </label>
        <div className="treatyexchangeyearly__baris">
          <label className="field">
            <span className="field__label">
              {TEY.tahun} <span aria-hidden="true">*</span>
            </span>
            <input
              className="field__input"
              value={isi.treatyYear}
              inputMode="numeric"
              maxLength={4}
              required
              autoFocus
              onChange={(e) => setIsi((i) => gantiTahun(i, e.target.value, mode.jenis === 'tambah'))}
            />
          </label>
          <label className="field">
            <span className="field__label">
              {TEY.mataUang} <span aria-hidden="true">*</span>
            </span>
            <select
              className="field__input"
              value={isi.idCurrency}
              required
              onChange={(e) => ubahIsian({ idCurrency: e.target.value })}
            >
              <option value="">{TEY.pilih}</option>
              {opsiMataUang(mataUang, baris?.idCurrency ?? '', baris?.currency ?? '').map((o) => (
                <option key={o.value} value={o.value}>
                  {o.label}
                </option>
              ))}
            </select>
          </label>
        </div>
        <div className="treatyexchangeyearly__baris">
          <label className="field">
            <span className="field__label">
              {TEY.mulai} <span aria-hidden="true">*</span>
            </span>
            <input
              className="field__input"
              type="date"
              value={isi.startDate}
              required
              onChange={(e) => ubahIsian({ startDate: e.target.value })}
            />
          </label>
          <label className="field">
            <span className="field__label">
              {TEY.akhir} <span aria-hidden="true">*</span>
            </span>
            <input
              className="field__input"
              type="date"
              value={isi.endDate}
              required
              onChange={(e) => ubahIsian({ endDate: e.target.value })}
            />
          </label>
        </div>
        <div className="treatyexchangeyearly__baris">
          <label className="field">
            <span className="field__label">
              {TEY.toIdr} <span aria-hidden="true">*</span>
            </span>
            <input
              className="field__input"
              value={isi.toIdr}
              inputMode="decimal"
              maxLength={30}
              required
              onChange={(e) => ubahIsian({ toIdr: e.target.value })}
            />
          </label>
          <label className="field">
            <span className="field__label">{TEY.toUsd}</span>
            <input
              className="field__input"
              value={isi.toUsd}
              inputMode="decimal"
              maxLength={30}
              onChange={(e) => ubahIsian({ toUsd: e.target.value })}
            />
          </label>
        </div>
        <label className="field">
          <span className="field__label">
            {TEY.quarter} <span aria-hidden="true">*</span>
          </span>
          <select
            className="field__input"
            value={isi.quarter}
            required
            onChange={(e) => ubahIsian({ quarter: e.target.value })}
          >
            {QUARTER.map((q) => (
              <option key={q} value={q}>
                {TEY.labelQuarter(q)}
              </option>
            ))}
          </select>
        </label>
        <p className="muted treatyexchangeyearly__catatan">{TEY.catatan}</p>
      </div>
    </Modal>
  )
}
