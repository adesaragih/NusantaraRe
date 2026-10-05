// Halaman Treaty Group (perintah work owner 05-10-2026) - daftar `TREATYGROUP` (urut ORDERNO di backend; Order No tidak
// tampil), dengan cari, saringan OJK Business, View (beserta grup bisnis anaknya), Add, dan Edit. Tanpa hapus (keputusan work owner). Menu View only:
// tanpa Add dan Edit (backend juga menolak tulisnya).

import { useCallback, useEffect, useState } from 'react'

import { Gagal, Kosong, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import { useBolehUbah } from '../../../../inti/frontend/lib/hakMenu'
import { ambilDaftar, ambilPilihan, type Grup, type Ojk } from '../api'
import { kodeNama } from '../aturan'
import FormGrup, { type ModeForm } from '../components/FormGrup'
import { TG } from '../labels'
import { NAMA_TG } from '../menu'

const JEDA_CARI_MS = 300

export default function TreatyGroup() {
  const bolehUbah = useBolehUbah(NAMA_TG)
  const [kata, setKata] = useState('')
  const [cari, setCari] = useState('')
  const [saringOjk, setSaringOjk] = useState('')
  const [ojk, setOjk] = useState<Ojk[]>([])
  const [daftar, setDaftar] = useState<Grup[] | null>(null)
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
      (p) => setOjk(p.ojk),
      (g: unknown) => setGalat(g),
    )
  }, [])

  const muat = useCallback((q: string, o: string) => {
    ambilDaftar(q, o).then(
      (d) => {
        setDaftar(d.daftar)
        setGalat(null)
      },
      (g: unknown) => setGalat(g),
    )
  }, [])
  useEffect(() => {
    muat(cari, saringOjk)
  }, [muat, cari, saringOjk, segar])

  return (
    <section className="inbox treatygroup__akar">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">{TG.judul}</h2>
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
            {TG.tambah}
          </button>
        )}
      </header>

      <div className="treatygroup__alat">
        <input
          className="field__input treatygroup__cari"
          type="search"
          placeholder={TG.cari}
          aria-label={TG.cari}
          value={kata}
          onChange={(e) => setKata(e.target.value)}
        />
        <select
          className="field__input treatygroup__saring"
          aria-label={TG.ojk}
          value={saringOjk}
          onChange={(e) => setSaringOjk(e.target.value)}
        >
          <option value="">{TG.semuaOjk}</option>
          {ojk.map((o) => (
            <option key={o.id} value={o.id}>
              {kodeNama(o.id, o.name)}
            </option>
          ))}
        </select>
        <span className="toolbar__spacer" />
        {daftar !== null && <span className="muted">{TG.jumlah(daftar.length)}</span>}
      </div>

      {pesan !== null && (
        <div className="alert alert--ok" role="status">
          {pesan}
        </div>
      )}
      {galat !== null && <Gagal galat={galat} />}
      {daftar === null && galat === null && <Memuat pesan={TG.memuat} />}
      {daftar !== null && daftar.length === 0 && (
        <Kosong pesan={cari.trim() === '' && saringOjk === '' ? TG.kosong : TG.tidakCocok} />
      )}
      {daftar !== null && daftar.length > 0 && (
        <div className="treatygroup__gulir">
          <table className="inbox__tabel treatygroup__tabel">
            <thead>
              <tr>
                <th>{TG.no}</th>
                <th>{TG.id}</th>
                <th>{TG.nama}</th>
                <th>{TG.soa}</th>
                <th>{TG.ojk}</th>
                <th>{TG.coa}</th>
                <th>{TG.pengubah}</th>
                <th>{TG.tglUbah}</th>
                {bolehUbah && <th className="table__actions">{TG.aksi}</th>}
              </tr>
            </thead>
            <tbody>
              {daftar.map((g, i) => (
                <tr key={g.id} className="inbox__baris">
                  <td>{i + 1}</td>
                  <td>{g.id}</td>
                  <td>{g.name}</td>
                  <td>{g.soaName}</td>
                  <td>{kodeNama(g.ojkId, g.ojkName)}</td>
                  <td>{kodeNama(g.coaId, g.coaName)}</td>
                  <td>{g.userId}</td>
                  <td>{g.diubah}</td>
                  {bolehUbah && (
                    <td className="table__actions">
                      <span className="treatygroup__aksi">
                        <button
                          type="button"
                          className="btn btn--ghost btn--sm"
                          onClick={() => {
                            setPesan(null)
                            setForm({ jenis: 'ubah', baris: g })
                          }}
                        >
                          {TG.edit}
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
        <FormGrup
          mode={form}
          ojk={ojk}
          onTutup={() => setForm(null)}
          onTersimpan={(d) => {
            setForm(null)
            setPesan(TG.tersimpan(d.id))
            setSegar((n) => n + 1)
          }}
        />
      )}
    </section>
  )
}
