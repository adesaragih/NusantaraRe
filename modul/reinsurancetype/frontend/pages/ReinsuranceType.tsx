// Halaman Reinsurance Type (perintah work owner 05-10-2026) - daftar `REINSURANCETYPE` urut Type lalu nama, dengan cari,
// saringan Type dan Flag, View, Add, dan Edit. Tanpa hapus - nonaktif lewat Flag inactive (keputusan work owner). Menu
// View only: tanpa Add dan Edit (backend juga menolak tulisnya).

import { useCallback, useEffect, useState } from 'react'

import { Gagal, Kosong, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import { useBolehUbah } from '../../../../inti/frontend/lib/hakMenu'
import { ambilDaftar, type Jenis } from '../api'
import { FLAG, kelasFlag, TYPE } from '../aturan'
import FormJenis, { type ModeForm } from '../components/FormJenis'
import { RT } from '../labels'
import { NAMA_RT } from '../menu'

const JEDA_CARI_MS = 300

export default function ReinsuranceType() {
  const bolehUbah = useBolehUbah(NAMA_RT)
  const [kata, setKata] = useState('')
  const [cari, setCari] = useState('')
  const [saringType, setSaringType] = useState('')
  const [saringFlag, setSaringFlag] = useState('')
  const [daftar, setDaftar] = useState<Jenis[] | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [segar, setSegar] = useState(0)
  const [form, setForm] = useState<ModeForm | null>(null)
  const [pesan, setPesan] = useState<string | null>(null)

  useEffect(() => {
    const t = setTimeout(() => setCari(kata), JEDA_CARI_MS)
    return () => clearTimeout(t)
  }, [kata])

  const muat = useCallback((q: string, t: string, f: string) => {
    ambilDaftar(q, t, f).then(
      (d) => {
        setDaftar(d.daftar)
        setGalat(null)
      },
      (g: unknown) => setGalat(g),
    )
  }, [])
  useEffect(() => {
    muat(cari, saringType, saringFlag)
  }, [muat, cari, saringType, saringFlag, segar])

  const tersaring = cari.trim() !== '' || saringType !== '' || saringFlag !== ''

  return (
    <section className="inbox reinsurancetype__akar">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">{RT.judul}</h2>
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
            {RT.tambah}
          </button>
        )}
      </header>

      <div className="reinsurancetype__alat">
        <input
          className="field__input reinsurancetype__cari"
          type="search"
          placeholder={RT.cari}
          aria-label={RT.cari}
          value={kata}
          onChange={(e) => setKata(e.target.value)}
        />
        <select
          className="field__input reinsurancetype__saring"
          aria-label={RT.type}
          value={saringType}
          onChange={(e) => setSaringType(e.target.value)}
        >
          <option value="">{RT.semuaType}</option>
          {TYPE.map((t) => (
            <option key={t} value={t}>
              {RT.labelType(t)}
            </option>
          ))}
        </select>
        <select
          className="field__input reinsurancetype__saring"
          aria-label={RT.flag}
          value={saringFlag}
          onChange={(e) => setSaringFlag(e.target.value)}
        >
          <option value="">{RT.semuaFlag}</option>
          {FLAG.map((f) => (
            <option key={f} value={f}>
              {RT.labelFlag(f)}
            </option>
          ))}
        </select>
        <span className="toolbar__spacer" />
        {daftar !== null && <span className="muted">{RT.jumlah(daftar.length)}</span>}
      </div>

      {pesan !== null && (
        <div className="alert alert--ok" role="status">
          {pesan}
        </div>
      )}
      {galat !== null && <Gagal galat={galat} />}
      {daftar === null && galat === null && <Memuat pesan={RT.memuat} />}
      {daftar !== null && daftar.length === 0 && <Kosong pesan={tersaring ? RT.tidakCocok : RT.kosong} />}
      {daftar !== null && daftar.length > 0 && (
        <div className="reinsurancetype__gulir">
          <table className="inbox__tabel reinsurancetype__tabel">
            <thead>
              <tr>
                <th>{RT.no}</th>
                <th>{RT.id}</th>
                <th>{RT.nama}</th>
                <th>{RT.type}</th>
                <th>{RT.soa}</th>
                <th>{RT.flag}</th>
                <th>{RT.pengubah}</th>
                <th>{RT.tglUbah}</th>
                {bolehUbah && <th className="table__actions">{RT.aksi}</th>}
              </tr>
            </thead>
            <tbody>
              {daftar.map((j, i) => (
                <tr
                  key={j.id}
                  className={j.flag === 'inactive' ? 'inbox__baris reinsurancetype__baris--nonaktif' : 'inbox__baris'}
                >
                  <td>{i + 1}</td>
                  <td>{j.id}</td>
                  <td>{j.name}</td>
                  <td>{RT.labelType(j.type)}</td>
                  <td>{j.soaName}</td>
                  <td>{j.flag !== '' && <span className={kelasFlag(j.flag)}>{RT.labelFlag(j.flag)}</span>}</td>
                  <td>{j.userId}</td>
                  <td>{j.diubah}</td>
                  {bolehUbah && (
                    <td className="table__actions">
                      <span className="reinsurancetype__aksi">
                        <button
                          type="button"
                          className="btn btn--ghost btn--sm"
                          onClick={() => {
                            setPesan(null)
                            setForm({ jenis: 'ubah', baris: j })
                          }}
                        >
                          {RT.edit}
                        </button>
                      </span>
                    </td>
                  )}
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {form !== null && (
        <FormJenis
          mode={form}
          onTutup={() => setForm(null)}
          onTersimpan={(j) => {
            setForm(null)
            setPesan(RT.tersimpan(j.id))
            setSegar((n) => n + 1)
          }}
        />
      )}
    </section>
  )
}
