// Halaman Treaty Group OJK (perintah work owner 05-10-2026) - daftar lini bisnis OJK (`TREATYGROUPOJK`), dengan cari,
// View, Add, dan Edit. Urutannya tetap ORDERNO (backend), tetapi Order No tidak tampil ("ORDERNO hide aja"). Tanpa hapus (keputusan work owner). Menu View only: tanpa Add dan Edit (backend juga
// menolak tulisnya).

import { useCallback, useEffect, useState } from 'react'

import { Gagal, Kosong, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import { useBolehUbah } from '../../../../inti/frontend/lib/hakMenu'
import { ambilDaftar, type Ojk } from '../api'
import FormOjk, { type ModeForm } from '../components/FormOjk'
import { TGO } from '../labels'
import { NAMA_TGO } from '../menu'

const JEDA_CARI_MS = 300

export default function TreatyGroupOjk() {
  const bolehUbah = useBolehUbah(NAMA_TGO)
  const [kata, setKata] = useState('')
  const [cari, setCari] = useState('')
  const [daftar, setDaftar] = useState<Ojk[] | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [segar, setSegar] = useState(0)
  const [form, setForm] = useState<ModeForm | null>(null)
  const [pesan, setPesan] = useState<string | null>(null)

  useEffect(() => {
    const t = setTimeout(() => setCari(kata), JEDA_CARI_MS)
    return () => clearTimeout(t)
  }, [kata])

  const muat = useCallback((q: string) => {
    ambilDaftar(q).then(
      (d) => {
        setDaftar(d.daftar)
        setGalat(null)
      },
      (g: unknown) => setGalat(g),
    )
  }, [])
  useEffect(() => {
    muat(cari)
  }, [muat, cari, segar])

  return (
    <section className="inbox treatygroupojk__akar">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">{TGO.judul}</h2>
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
            {TGO.tambah}
          </button>
        )}
      </header>

      <div className="treatygroupojk__alat">
        <input
          className="field__input treatygroupojk__cari"
          type="search"
          placeholder={TGO.cari}
          aria-label={TGO.cari}
          value={kata}
          onChange={(e) => setKata(e.target.value)}
        />
        <span className="toolbar__spacer" />
        {daftar !== null && <span className="muted">{TGO.jumlah(daftar.length)}</span>}
      </div>

      {pesan !== null && (
        <div className="alert alert--ok" role="status">
          {pesan}
        </div>
      )}
      {galat !== null && <Gagal galat={galat} />}
      {daftar === null && galat === null && <Memuat pesan={TGO.memuat} />}
      {daftar !== null && daftar.length === 0 && <Kosong pesan={cari.trim() === '' ? TGO.kosong : TGO.tidakCocok} />}
      {daftar !== null && daftar.length > 0 && (
        <div className="treatygroupojk__gulir">
          <table className="inbox__tabel treatygroupojk__tabel">
            <thead>
              <tr>
                <th>{TGO.no}</th>
                <th>{TGO.id}</th>
                <th>{TGO.nama}</th>
                <th>{TGO.namaIdn}</th>
                {bolehUbah && <th className="table__actions">{TGO.aksi}</th>}
              </tr>
            </thead>
            <tbody>
              {daftar.map((o, i) => (
                <tr key={o.id} className="inbox__baris">
                  <td>{i + 1}</td>
                  <td>{o.id}</td>
                  <td>{o.name}</td>
                  <td>{o.nameIdn}</td>
                  {bolehUbah && (
                    <td className="table__actions">
                      <span className="treatygroupojk__aksi">
                        <button
                          type="button"
                          className="btn btn--ghost btn--sm"
                          onClick={() => {
                            setPesan(null)
                            setForm({ jenis: 'ubah', baris: o })
                          }}
                        >
                          {TGO.edit}
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
        <FormOjk
          mode={form}
          onTutup={() => setForm(null)}
          onTersimpan={(o) => {
            setForm(null)
            setPesan(TGO.tersimpan(o.id))
            setSegar((n) => n + 1)
          }}
        />
      )}
    </section>
  )
}
