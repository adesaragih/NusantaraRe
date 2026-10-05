// Popup Add / Edit / View Business Group. View = kepala bernama + lencana ID dan satu kartu rincian (seperti View
// Treaty Exchange Yearly; "buat rapih", 05-10-2026). Add / Edit: Treaty Group dari dropdown TREATYGROUP
// (namanya disalin backend ke TREATYNAME); Name dan Alias Name ditulis huruf besar oleh backend; Alias kosong = Name.

import { useState } from 'react'

import { Gagal, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { tambah, ubah, type BisnisGrup, type Isian, type TreatyGroup } from '../api'
import { isianDari, kodeNama, opsiTreaty, periksaIsian, rapikanIsian } from '../aturan'
import { BG } from '../labels'

export type ModeForm =
  { jenis: 'tambah' } | { jenis: 'ubah'; baris: BisnisGrup } | { jenis: 'lihat'; baris: BisnisGrup }

function Pasangan({ isi }: { isi: [string, string][] }) {
  return (
    <dl className="businessgroup__pasangan">
      {isi.map(([label, nilai]) => (
        <div key={label}>
          <dt>{label}</dt>
          <dd>{nilai === '' ? '—' : nilai}</dd>
        </div>
      ))}
    </dl>
  )
}

function LihatBisnisGrup({ bisnis: b, onTutup }: { bisnis: BisnisGrup; onTutup: () => void }) {
  return (
    <Modal judul={BG.judulLihat(b.id)} onTutup={onTutup} labelBatal={BG.tutup} lebar>
      <div className="businessgroup__lihat">
        <div className="businessgroup__lihat-kepala">
          <div className="businessgroup__lihat-mata">
            <p className="businessgroup__lihat-kode">{b.name === '' ? '—' : b.name}</p>
          </div>
          <span className="businessgroup__lencana">
            {BG.id} {b.id}
          </span>
        </div>
        <div className="businessgroup__lihat-kartu">
          <section className="businessgroup__kartu businessgroup__kartu--penuh">
            <h3 className="businessgroup__judul-kartu">{BG.rincian}</h3>
            <Pasangan
              isi={[
                [BG.alias, b.alias],
                [BG.treaty, kodeNama(b.topId, b.treatyName)],
              ]}
            />
          </section>
        </div>
      </div>
    </Modal>
  )
}

export default function FormBisnisGrup({
  mode,
  treaty,
  onTutup,
  onTersimpan,
}: {
  mode: ModeForm
  treaty: TreatyGroup[]
  onTutup: () => void
  onTersimpan: (b: BisnisGrup) => void
}) {
  const baris = mode.jenis === 'tambah' ? null : mode.baris
  const [isi, setIsi] = useState<Isian>(isianDari(baris))
  const [galat, setGalat] = useState<unknown>(null)
  const [pesan, setPesan] = useState<string | null>(null)
  const [sibuk, setSibuk] = useState(false)

  const judul =
    mode.jenis === 'tambah'
      ? BG.judulTambah
      : mode.jenis === 'ubah'
        ? BG.judulUbah(mode.baris.id)
        : BG.judulLihat(mode.baris.id)

  if (mode.jenis === 'lihat') return <LihatBisnisGrup bisnis={mode.baris} onTutup={onTutup} />

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
      (b) => {
        setSibuk(false)
        onTersimpan(b)
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
      labelBatal={BG.batal}
      aksi={
        <button type="submit" className="btn btn--primary" disabled={sibuk}>
          {sibuk ? BG.menyimpan : BG.save}
        </button>
      }
    >
      <div className="businessgroup__form">
        {galat !== null && <Gagal galat={galat} />}
        {pesan !== null && (
          <div className="alert alert--warn" role="alert">
            {pesan}
          </div>
        )}
        <label className="field">
          <span className="field__label">{BG.id}</span>
          <input
            className="field__input field__input--readonly"
            value={baris?.id ?? ''}
            placeholder={BG.idOtomatis}
            readOnly
          />
        </label>
        <label className="field">
          <span className="field__label">
            {BG.treaty} <span aria-hidden="true">*</span>
          </span>
          <select
            className="field__input"
            value={isi.topId}
            required
            autoFocus
            onChange={(e) => ubahIsian({ topId: e.target.value })}
          >
            <option value="">{BG.pilih}</option>
            {opsiTreaty(treaty, baris?.topId ?? '', baris?.treatyName ?? '').map((o) => (
              <option key={o.value} value={o.value}>
                {o.label}
              </option>
            ))}
          </select>
        </label>
        <label className="field">
          <span className="field__label">
            {BG.nama} <span aria-hidden="true">*</span>
          </span>
          <input
            className="field__input"
            value={isi.name}
            maxLength={4000}
            required
            onChange={(e) => ubahIsian({ name: e.target.value.toUpperCase() })}
          />
        </label>
        <label className="field">
          <span className="field__label">{BG.alias}</span>
          <input
            className="field__input"
            value={isi.alias}
            maxLength={4000}
            placeholder={BG.aliasOtomatis}
            onChange={(e) => ubahIsian({ alias: e.target.value.toUpperCase() })}
          />
        </label>
        <p className="muted businessgroup__catatan">{BG.catatan}</p>
      </div>
    </Modal>
  )
}
