// Halaman Treaty Exchange Yearly (perintah work owner 05-10-2026) - daftar kurs tahunan `TREATYEXCHANGEYEARLY` urut
// tahun treaty terbaru, dengan cari, saringan tahun, View, Add, dan Edit. Tanpa hapus (keputusan work owner). Menu View
// only: tanpa Add dan Edit (backend juga menolak tulisnya).

import { useCallback, useEffect, useState } from 'react'

import { Gagal, Kosong, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import { useBolehUbah } from '../../../../inti/frontend/lib/hakMenu'
import { ambilDaftar, ambilPilihan, type Kurs, type MataUang } from '../api'
import FormKurs, { type ModeForm } from '../components/FormKurs'
import { TEY } from '../labels'
import { NAMA_TEY } from '../menu'

const JEDA_CARI_MS = 300

export default function TreatyExchangeYearly() {
  const bolehUbah = useBolehUbah(NAMA_TEY)
  const [kata, setKata] = useState('')
  const [cari, setCari] = useState('')
  const [saringTahun, setSaringTahun] = useState('')
  const [mataUang, setMataUang] = useState<MataUang[]>([])
  const [tahun, setTahun] = useState<string[]>([])
  const [daftar, setDaftar] = useState<Kurs[] | null>(null)
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
      (p) => {
        setMataUang(p.mataUang)
        setTahun(p.tahun)
      },
      (g: unknown) => setGalat(g),
    )
  }, [segar])

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
    muat(cari, saringTahun)
  }, [muat, cari, saringTahun, segar])

  return (
    <section className="inbox treatyexchangeyearly__akar">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">{TEY.judul}</h2>
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
            {TEY.tambah}
          </button>
        )}
      </header>

      <div className="treatyexchangeyearly__alat">
        <input
          className="field__input treatyexchangeyearly__cari"
          type="search"
          placeholder={TEY.cari}
          aria-label={TEY.cari}
          value={kata}
          onChange={(e) => setKata(e.target.value)}
        />
        <select
          className="field__input treatyexchangeyearly__saring"
          aria-label={TEY.tahun}
          value={saringTahun}
          onChange={(e) => setSaringTahun(e.target.value)}
        >
          <option value="">{TEY.semuaTahun}</option>
          {tahun.map((t) => (
            <option key={t} value={t}>
              {t}
            </option>
          ))}
        </select>
        <span className="toolbar__spacer" />
        {daftar !== null && <span className="muted">{TEY.jumlah(daftar.length)}</span>}
      </div>

      {pesan !== null && (
        <div className="alert alert--ok" role="status">
          {pesan}
        </div>
      )}
      {galat !== null && <Gagal galat={galat} />}
      {daftar === null && galat === null && <Memuat pesan={TEY.memuat} />}
      {daftar !== null && daftar.length === 0 && (
        <Kosong pesan={cari.trim() === '' && saringTahun === '' ? TEY.kosong : TEY.tidakCocok} />
      )}
      {daftar !== null && daftar.length > 0 && (
        <div className="treatyexchangeyearly__gulir">
          <table className="inbox__tabel treatyexchangeyearly__tabel">
            <thead>
              <tr>
                <th>{TEY.no}</th>
                <th>{TEY.id}</th>
                <th>{TEY.tahun}</th>
                <th>{TEY.mataUang}</th>
                <th>{TEY.quarter}</th>
                <th>{TEY.mulai}</th>
                <th>{TEY.akhir}</th>
                <th className="treatyexchangeyearly__angka">{TEY.toIdr}</th>
                <th className="treatyexchangeyearly__angka">{TEY.toUsd}</th>
                <th>{TEY.pengubah}</th>
                <th>{TEY.tglUbah}</th>
                {bolehUbah && <th className="table__actions">{TEY.aksi}</th>}
              </tr>
            </thead>
            <tbody>
              {daftar.map((k, i) => (
                <tr key={k.kunci} className="inbox__baris">
                  <td>{i + 1}</td>
                  <td>{k.id}</td>
                  <td>{k.treatyYear}</td>
                  <td>{k.currency}</td>
                  <td>{k.quarter}</td>
                  <td>{k.mulai}</td>
                  <td>{k.akhir}</td>
                  <td className="treatyexchangeyearly__angka">{k.toIdr}</td>
                  <td className="treatyexchangeyearly__angka">{k.toUsd}</td>
                  <td>{k.userId}</td>
                  <td>{k.diubah}</td>
                  {bolehUbah && (
                    <td className="table__actions">
                      <span className="treatyexchangeyearly__aksi">
                        <button
                          type="button"
                          className="btn btn--ghost btn--sm"
                          onClick={() => {
                            setPesan(null)
                            setForm({ jenis: 'ubah', baris: k })
                          }}
                        >
                          {TEY.edit}
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
        <FormKurs
          mode={form}
          mataUang={mataUang}
          onTutup={() => setForm(null)}
          onTersimpan={(k) => {
            setForm(null)
            setPesan(TEY.tersimpan(k.id))
            setSegar((n) => n + 1)
          }}
        />
      )}
    </section>
  )
}
