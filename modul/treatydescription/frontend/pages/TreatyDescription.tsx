// Halaman Treaty Description (perintah work owner 05-10-2026; tampilan rancangan sendiri: "buat versi kamu") - daftar
// `TREATYDESC` urut ID (hanya kolom tabel itu), dengan cari, saringan Type (Non XOL / XOL) dan Status (Active /
// Inactive), View, Add, dan Edit. Tanpa hapus. Menu View only: tanpa Add dan Edit (backend juga menolak tulisnya).

import { useCallback, useEffect, useState } from 'react'

import { Gagal, Kosong, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import { useBolehUbah } from '../../../../inti/frontend/lib/hakMenu'
import { ambilDaftar, type Desc, type Saringan } from '../api'
import { labelJenis, labelStatus } from '../aturan'
import FormDesc, { type ModeForm } from '../components/FormDesc'
import Segmen from '../components/Segmen'
import { TD } from '../labels'
import { NAMA_TD } from '../menu'

const JEDA_CARI_MS = 300

const OPSI_JENIS = [
  { nilai: '', label: TD.semua },
  { nilai: '0', label: TD.nonXol },
  { nilai: '1', label: TD.xol },
] as const
const OPSI_STATUS = [
  { nilai: '', label: TD.semua },
  { nilai: '1', label: TD.aktif },
  { nilai: '0', label: TD.nonaktif },
] as const

export default function TreatyDescription() {
  const bolehUbah = useBolehUbah(NAMA_TD)
  const [kata, setKata] = useState('')
  const [saringan, setSaringan] = useState<Saringan>({ q: '', xol: '', status: '' })
  const [daftar, setDaftar] = useState<Desc[] | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [segar, setSegar] = useState(0)
  const [form, setForm] = useState<ModeForm | null>(null)
  const [pesan, setPesan] = useState<string | null>(null)

  useEffect(() => {
    const t = setTimeout(() => setSaringan((s) => ({ ...s, q: kata })), JEDA_CARI_MS)
    return () => clearTimeout(t)
  }, [kata])

  const muat = useCallback((s: Saringan) => {
    ambilDaftar(s).then(
      (d) => {
        setDaftar(d.daftar)
        setGalat(null)
      },
      (g: unknown) => setGalat(g),
    )
  }, [])
  useEffect(() => {
    muat(saringan)
  }, [muat, saringan, segar])

  const tanpaSaringan = saringan.q.trim() === '' && saringan.xol === '' && saringan.status === ''

  return (
    <section className="inbox treatydescription__akar">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">{TD.judul}</h2>
        <span className="toolbar__spacer" />
        {bolehUbah && (
          <button
            type="button"
            className="btn btn--primary"
            onClick={() => {
              setPesan(null)
              setForm({ jenis: 'tambah' })
            }}
          >
            {TD.tambah}
          </button>
        )}
      </header>

      <div className="treatydescription__alat">
        <input
          className="field__input treatydescription__cari"
          type="search"
          placeholder={TD.cari}
          aria-label={TD.cari}
          value={kata}
          onChange={(e) => setKata(e.target.value)}
        />
        <Segmen
          label={TD.jenis}
          opsi={OPSI_JENIS}
          nilai={saringan.xol}
          onPilih={(xol) => setSaringan((s) => ({ ...s, xol }))}
        />
        <Segmen
          label={TD.status}
          opsi={OPSI_STATUS}
          nilai={saringan.status}
          onPilih={(status) => setSaringan((s) => ({ ...s, status }))}
        />
        <span className="toolbar__spacer" />
        {daftar !== null && <span className="muted">{TD.jumlah(daftar.length)}</span>}
      </div>

      {pesan !== null && (
        <div className="alert alert--ok" role="status">
          {pesan}
        </div>
      )}
      {galat !== null && <Gagal galat={galat} />}
      {daftar === null && galat === null && <Memuat pesan={TD.memuat} />}
      {daftar !== null && daftar.length === 0 && <Kosong pesan={tanpaSaringan ? TD.kosong : TD.tidakCocok} />}
      {daftar !== null && daftar.length > 0 && (
        <div className="treatydescription__gulir">
          <table className="inbox__tabel treatydescription__tabel">
            <thead>
              <tr>
                <th>{TD.no}</th>
                <th>{TD.id}</th>
                <th>{TD.nama}</th>
                <th>{TD.jenis}</th>
                <th className="treatydescription__kolom-status">{TD.status}</th>
                <th className="table__actions">{TD.aksi}</th>
              </tr>
            </thead>
            <tbody>
              {daftar.map((d, i) => (
                <tr key={d.id} className={d.aktif ? 'inbox__baris' : 'inbox__baris treatydescription__baris--nonaktif'}>
                  <td>{i + 1}</td>
                  <td>{d.id}</td>
                  <td>{d.descName}</td>
                  <td>
                    <span
                      className={
                        d.isXol === '1'
                          ? 'treatydescription__jenis treatydescription__jenis--xol'
                          : 'treatydescription__jenis'
                      }
                    >
                      {labelJenis(d.isXol)}
                    </span>
                  </td>
                  <td className="treatydescription__kolom-status">
                    <span
                      className={
                        d.aktif
                          ? 'treatydescription__status'
                          : 'treatydescription__status treatydescription__status--nonaktif'
                      }
                    >
                      {labelStatus(d.aktif)}
                    </span>
                  </td>
                  <td className="table__actions">
                    <span className="treatydescription__aksi">
                      <button
                        type="button"
                        className="btn btn--ghost btn--sm"
                        onClick={() => setForm({ jenis: 'lihat', baris: d })}
                      >
                        {TD.view}
                      </button>
                      {bolehUbah && (
                        <button
                          type="button"
                          className="btn btn--ghost btn--sm"
                          onClick={() => {
                            setPesan(null)
                            setForm({ jenis: 'ubah', baris: d })
                          }}
                        >
                          {TD.edit}
                        </button>
                      )}
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {form !== null && (
        <FormDesc
          mode={form}
          onTutup={() => setForm(null)}
          onTersimpan={(d) => {
            setForm(null)
            setPesan(TD.tersimpan(d.id))
            setSegar((n) => n + 1)
          }}
        />
      )}
    </section>
  )
}
