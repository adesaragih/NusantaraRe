// Popup Add / Edit / View Treaty Description - hanya kolom TREATYDESC. View = kepala bernama + lencana ID dan kartu
// klasifikasi (seperti View Treaty Exchange Yearly; "buat rapih", 05-10-2026). Add / Edit: ID dibuat backend;
// Description Name huruf besar; Type dan Status bersegmen.

import { useState } from 'react'

import { Gagal, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { tambah, ubah, type Desc, type Isian, type Jenis, type Status } from '../api'
import { BATAS_NAMA, isianDari, labelJenis, labelStatus, periksaIsian, rapikanIsian } from '../aturan'
import { TD } from '../labels'
import Segmen from './Segmen'

export type ModeForm = { jenis: 'tambah' } | { jenis: 'ubah'; baris: Desc } | { jenis: 'lihat'; baris: Desc }

const OPSI_JENIS: readonly { nilai: Jenis; label: string }[] = [
  { nilai: '0', label: TD.nonXol },
  { nilai: '1', label: TD.xol },
]
const OPSI_STATUS: readonly { nilai: Status; label: string }[] = [
  { nilai: '1', label: TD.aktif },
  { nilai: '0', label: TD.nonaktif },
]

function Pasangan({ isi }: { isi: [string, string][] }) {
  return (
    <dl className="treatydescription__pasangan">
      {isi.map(([label, nilai]) => (
        <div key={label}>
          <dt>{label}</dt>
          <dd>{nilai === '' ? '—' : nilai}</dd>
        </div>
      ))}
    </dl>
  )
}

function Lihat({ baris, onTutup }: { baris: Desc; onTutup: () => void }) {
  return (
    <Modal judul={TD.judulLihat(baris.id)} onTutup={onTutup} labelBatal={TD.tutup} lebar>
      <div className="treatydescription__lihat">
        <div className="treatydescription__lihat-kepala">
          <div className="treatydescription__lihat-mata">
            <p className="treatydescription__lihat-kode">{baris.descName === '' ? '—' : baris.descName}</p>
          </div>
          <span className="treatydescription__lencana">
            {TD.id} {baris.id}
          </span>
        </div>
        <div className="treatydescription__lihat-kartu">
          <section className="treatydescription__kartu treatydescription__kartu--penuh">
            <h3 className="treatydescription__judul-kartu">{TD.klasifikasi}</h3>
            <Pasangan
              isi={[
                [TD.jenis, labelJenis(baris.isXol)],
                [TD.status, labelStatus(baris.aktif)],
              ]}
            />
          </section>
        </div>
      </div>
    </Modal>
  )
}

export default function FormDesc({
  mode,
  onTutup,
  onTersimpan,
}: {
  mode: ModeForm
  onTutup: () => void
  onTersimpan: (d: Desc) => void
}) {
  const baris = mode.jenis === 'tambah' ? null : mode.baris
  const [isi, setIsi] = useState<Isian>(isianDari(baris))
  const [galat, setGalat] = useState<unknown>(null)
  const [pesan, setPesan] = useState<string | null>(null)
  const [sibuk, setSibuk] = useState(false)

  if (mode.jenis === 'lihat') return <Lihat baris={mode.baris} onTutup={onTutup} />

  const judul = mode.jenis === 'tambah' ? TD.judulTambah : TD.judulUbah(mode.baris.id)

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
      labelBatal={TD.batal}
      aksi={
        <button type="submit" className="btn btn--primary" disabled={sibuk}>
          {sibuk ? TD.menyimpan : TD.save}
        </button>
      }
    >
      <div className="treatydescription__form">
        {galat !== null && <Gagal galat={galat} />}
        {pesan !== null && (
          <div className="alert alert--warn" role="alert">
            {pesan}
          </div>
        )}
        <label className="field">
          <span className="field__label">{TD.id}</span>
          <input
            className="field__input field__input--readonly"
            value={baris?.id ?? ''}
            placeholder={TD.idOtomatis}
            readOnly
          />
        </label>
        <label className="field">
          <span className="field__label">
            {TD.nama} <span aria-hidden="true">*</span>
          </span>
          <input
            className="field__input"
            value={isi.descName}
            maxLength={BATAS_NAMA}
            required
            autoFocus
            onChange={(e) => ubahIsian({ descName: e.target.value.toUpperCase() })}
          />
        </label>
        <div className="treatydescription__pilihan">
          <span className="field__label">{TD.jenis}</span>
          <Segmen label={TD.jenis} opsi={OPSI_JENIS} nilai={isi.isXol} onPilih={(isXol) => ubahIsian({ isXol })} />
        </div>
        <div className="treatydescription__pilihan">
          <span className="field__label">{TD.status}</span>
          <Segmen
            label={TD.status}
            opsi={OPSI_STATUS}
            nilai={isi.statusAktif}
            onPilih={(statusAktif) => ubahIsian({ statusAktif })}
          />
        </div>
        <p className="muted treatydescription__catatan">{TD.catatan}</p>
      </div>
    </Modal>
  )
}
