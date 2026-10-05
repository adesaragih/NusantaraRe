// Halaman Business Group (perintah work owner 05-10-2026) - daftar `BUSINESSGROUP` tanpa SYARIAH, urut nama, dengan
// cari, saringan Treaty Group, View, Add, dan Edit. Tanpa hapus (keputusan work owner). Menu View only: tanpa Add dan
// Edit (backend juga menolak tulisnya).

import { useCallback, useEffect, useState } from 'react'

import { Gagal, Kosong, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import { useBolehUbah } from '../../../../inti/frontend/lib/hakMenu'
import { ambilDaftar, ambilPilihan, type BisnisGrup, type TreatyGroup } from '../api'
import { kodeNama } from '../aturan'
import FormBisnisGrup, { type ModeForm } from '../components/FormBisnisGrup'
import { BG } from '../labels'
import { NAMA_BG } from '../menu'

const JEDA_CARI_MS = 300

export default function BusinessGroup() {
  const bolehUbah = useBolehUbah(NAMA_BG)
  const [kata, setKata] = useState('')
  const [cari, setCari] = useState('')
  const [saringTreaty, setSaringTreaty] = useState('')
  const [treaty, setTreaty] = useState<TreatyGroup[]>([])
  const [daftar, setDaftar] = useState<BisnisGrup[] | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [segar, setSegar] = useState(0)
  const [form, setForm] = useState<ModeForm | null>(null)
  const [pesan, setPesan] = useState<string | null>(null)

  useEffect(() => {
    const t = setTimeout(() => setCari(kata), JEDA_CARI_MS)
    return () => clearTimeout(t)
  }, [kata])

  useEffect(() => {
    ambilPilihan().then(
      (p) => setTreaty(p.treatyGroup),
      (g: unknown) => setGalat(g),
    )
  }, [])

  const muat = useCallback((q: string, t: string) => {
    ambilDaftar(q, t).then(
      (d) => {
        setDaftar(d.daftar)
        setGalat(null)
      },
      (g: unknown) => setGalat(g),
    )
  }, [])
  useEffect(() => {
    muat(cari, saringTreaty)
  }, [muat, cari, saringTreaty, segar])

  return (
    <section className="inbox businessgroup__akar">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">{BG.judul}</h2>
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
            {BG.tambah}
          </button>
        )}
      </header>

      <div className="businessgroup__alat">
        <input
          className="field__input businessgroup__cari"
          type="search"
          placeholder={BG.cari}
          aria-label={BG.cari}
          value={kata}
          onChange={(e) => setKata(e.target.value)}
        />
        <select
          className="field__input businessgroup__saring"
          aria-label={BG.treaty}
          value={saringTreaty}
          onChange={(e) => setSaringTreaty(e.target.value)}
        >
          <option value="">{BG.semuaTreaty}</option>
          {treaty.map((t) => (
            <option key={t.id} value={t.id}>
              {kodeNama(t.id, t.name)}
            </option>
          ))}
        </select>
        <span className="toolbar__spacer" />
        {daftar !== null && <span className="muted">{BG.jumlah(daftar.length)}</span>}
      </div>

      {pesan !== null && (
        <div className="alert alert--ok" role="status">
          {pesan}
        </div>
      )}
      {galat !== null && <Gagal galat={galat} />}
      {daftar === null && galat === null && <Memuat pesan={BG.memuat} />}
      {daftar !== null && daftar.length === 0 && (
        <Kosong pesan={cari.trim() === '' && saringTreaty === '' ? BG.kosong : BG.tidakCocok} />
      )}
      {daftar !== null && daftar.length > 0 && (
        <div className="businessgroup__gulir">
          <table className="inbox__tabel businessgroup__tabel">
            <thead>
              <tr>
                <th>{BG.no}</th>
                <th>{BG.id}</th>
                <th>{BG.nama}</th>
                <th>{BG.alias}</th>
                <th>{BG.treaty}</th>
                <th className="table__actions">{BG.aksi}</th>
              </tr>
            </thead>
            <tbody>
              {daftar.map((b, i) => (
                <tr key={b.id} className="inbox__baris">
                  <td>{i + 1}</td>
                  <td>{b.id}</td>
                  <td>{b.name}</td>
                  <td>{b.alias}</td>
                  <td>{kodeNama(b.topId, b.treatyName)}</td>
                  <td className="table__actions">
                    <span className="businessgroup__aksi">
                      <button
                        type="button"
                        className="btn btn--ghost btn--sm"
                        onClick={() => setForm({ jenis: 'lihat', baris: b })}
                      >
                        {BG.view}
                      </button>
                      {bolehUbah && (
                        <button
                          type="button"
                          className="btn btn--ghost btn--sm"
                          onClick={() => {
                            setPesan(null)
                            setForm({ jenis: 'ubah', baris: b })
                          }}
                        >
                          {BG.edit}
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
        <FormBisnisGrup
          mode={form}
          treaty={treaty}
          onTutup={() => setForm(null)}
          onTersimpan={(b) => {
            setForm(null)
            setPesan(BG.tersimpan(b.id))
            setSegar((n) => n + 1)
          }}
        />
      )}
    </section>
  )
}
