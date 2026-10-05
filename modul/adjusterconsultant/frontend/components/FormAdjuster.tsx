// Popup Add / Edit / View Adjuster Consultant - padanan section `MstAdjusterConsultant` mode isian (`NewAdjustConsult_Act`
// / `EditMstConsultant_Act` lalu Save `SaveAdjusterConsultant_Act`). View = teks saja (bukan isian dimatikan): kepala
// bernama + lencana ID dan status, kartu kontak, jejak ubah di kaki (seperti View Treaty Exchange Yearly; "buat rapih",
// 05-10-2026). Nama dan alamat ditulis huruf besar oleh backend, seperti Pega.

import { useState } from 'react'

import { Gagal, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { tambah, ubah, type Adjuster, type Isian } from '../api'
import { isianDari, periksaIsian, rapikanIsian } from '../aturan'
import { ADJ } from '../labels'

export type ModeForm = { jenis: 'tambah' } | { jenis: 'ubah'; baris: Adjuster } | { jenis: 'lihat'; baris: Adjuster }

function Pasangan({ isi }: { isi: [string, string][] }) {
  return (
    <dl className="adjusterconsultant__pasangan">
      {isi.map(([label, nilai]) => (
        <div key={label}>
          <dt>{label}</dt>
          <dd>{nilai === '' ? '—' : nilai}</dd>
        </div>
      ))}
    </dl>
  )
}

function LihatAdjuster({ adjuster: a, onTutup }: { adjuster: Adjuster; onTutup: () => void }) {
  return (
    <Modal judul={ADJ.judulLihat(a.id)} onTutup={onTutup} labelBatal={ADJ.tutup} lebar>
      <div className="adjusterconsultant__lihat">
        <div className="adjusterconsultant__lihat-kepala">
          <div className="adjusterconsultant__lihat-mata">
            <p className="adjusterconsultant__lihat-kode">{a.name === '' ? '—' : a.name}</p>
          </div>
          <span className="adjusterconsultant__lencana">
            {ADJ.id} {a.id}
          </span>
          <span
            className={
              a.active
                ? 'adjusterconsultant__status'
                : 'adjusterconsultant__status adjusterconsultant__status--nonaktif'
            }
          >
            {a.active ? ADJ.aktif : ADJ.nonaktif}
          </span>
        </div>
        <div className="adjusterconsultant__lihat-kartu">
          <section className="adjusterconsultant__kartu adjusterconsultant__kartu--penuh">
            <h3 className="adjusterconsultant__judul-kartu">{ADJ.kontak}</h3>
            <Pasangan
              isi={[
                [ADJ.telp, a.telpNo],
                [ADJ.alamat, a.address],
              ]}
            />
          </section>
        </div>
        <p className="muted adjusterconsultant__jejak">{ADJ.jejak(a.id, a.username, a.editDate)}</p>
      </div>
    </Modal>
  )
}

export default function FormAdjuster({
  mode,
  onTutup,
  onTersimpan,
}: {
  mode: ModeForm
  onTutup: () => void
  onTersimpan: (a: Adjuster) => void
}) {
  const baris = mode.jenis === 'tambah' ? null : mode.baris
  const [isi, setIsi] = useState<Isian>(isianDari(baris))
  const [galat, setGalat] = useState<unknown>(null)
  const [pesan, setPesan] = useState<string | null>(null)
  const [sibuk, setSibuk] = useState(false)

  const judul =
    mode.jenis === 'tambah'
      ? ADJ.judulTambah
      : mode.jenis === 'ubah'
        ? ADJ.judulUbah(mode.baris.id)
        : ADJ.judulLihat(mode.baris.id)

  if (mode.jenis === 'lihat') return <LihatAdjuster adjuster={mode.baris} onTutup={onTutup} />

  const kirim = () => {
    if (sibuk) return
    const dikirim = rapikanIsian(isi)
    const salah = periksaIsian(dikirim, ADJ.galatNama)
    setPesan(salah)
    setGalat(null)
    if (salah !== null) return
    setSibuk(true)
    const janji = mode.jenis === 'ubah' ? ubah(mode.baris.id, dikirim) : tambah(dikirim)
    janji.then(
      (a) => {
        setSibuk(false)
        onTersimpan(a)
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
      labelBatal={ADJ.batal}
      aksi={
        <button type="submit" className="btn btn--primary" disabled={sibuk}>
          {sibuk ? ADJ.menyimpan : ADJ.save}
        </button>
      }
    >
      <div className="adjusterconsultant__form">
        {galat !== null && <Gagal galat={galat} />}
        {pesan !== null && (
          <div className="alert alert--warn" role="alert">
            {pesan}
          </div>
        )}
        <label className="field">
          <span className="field__label">{ADJ.id}</span>
          <input className="field__input" value={baris?.id ?? ''} placeholder={ADJ.idOtomatis} readOnly />
        </label>
        <label className="field">
          <span className="field__label">
            {ADJ.nama} <span aria-hidden="true">*</span>
          </span>
          <input
            className="field__input"
            value={isi.name}
            maxLength={500}
            required
            autoFocus
            onChange={(e) => ubahIsian({ name: e.target.value })}
          />
        </label>
        <label className="field">
          <span className="field__label">{ADJ.alamat}</span>
          <textarea
            className="field__input"
            value={isi.address}
            maxLength={1000}
            onChange={(e) => ubahIsian({ address: e.target.value })}
          />
        </label>
        <label className="field">
          <span className="field__label">{ADJ.telp}</span>
          <input
            className="field__input"
            value={isi.telpNo}
            maxLength={50}
            onChange={(e) => ubahIsian({ telpNo: e.target.value })}
          />
        </label>
        <p className="muted adjusterconsultant__catatan">{ADJ.catatanHurufBesar}</p>
      </div>
    </Modal>
  )
}
