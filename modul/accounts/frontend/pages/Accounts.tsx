// Halaman Accounts (keputusan work owner 04-10-2026): daftar `T_M_ACCOUNT` dicari dan dibagi per halaman di server
// (18 ribu baris), terbaru dulu; "Add" membuka form akun baru. Akun hanya diinput sekali - nol View, nol Edit, nol
// hapus. Tampilan sama dengan Kelola User (tema disalin ke `accounts.css`).

import { useCallback, useEffect, useState } from 'react'

import { Gagal, Halaman, Kosong, Memuat, UKURAN_HALAMAN } from '../../../../inti/frontend/components/ui/dasar'
import { ambilDaftar, ambilPilihan, type Account, type GroupBusiness, type HalamanDaftar } from '../api'
import { atauStrip } from '../aturan'
import FormAccount from '../components/FormAccount'
import { ACC } from '../labels'

/** Jeda ketik sebelum server ditanya. */
const JEDA_MS = 300

export default function Accounts() {
  const [kata, setKata] = useState('')
  const [cari, setCari] = useState('')
  const [halaman, setHalaman] = useState(1)
  const [data, setData] = useState<HalamanDaftar | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [groupBusiness, setGroupBusiness] = useState<GroupBusiness[]>([])
  const [formBuka, setFormBuka] = useState(false)
  const [pesan, setPesan] = useState<string | null>(null)

  const muat = useCallback(() => {
    setData(null)
    ambilDaftar(cari, halaman, UKURAN_HALAMAN).then(
      (d) => {
        setData(d)
        setGalat(null)
      },
      (g: unknown) => {
        setGalat(g)
      },
    )
  }, [cari, halaman])

  useEffect(() => {
    muat()
  }, [muat])

  useEffect(() => {
    ambilPilihan().then(
      (p) => {
        setGroupBusiness(p.groupBusiness)
      },
      (g: unknown) => {
        setGalat(g)
      },
    )
  }, [])

  useEffect(() => {
    const jam = setTimeout(() => {
      setCari(kata.trim())
      setHalaman(1)
    }, JEDA_MS)
    return () => {
      clearTimeout(jam)
    }
  }, [kata])

  const tersimpan = (a: Account) => {
    setFormBuka(false)
    setPesan(ACC.tersimpan(a.idView))
    muat()
  }

  return (
    <section className="inbox accounts__akar">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">{ACC.judul}</h2>
      </header>

      <div className="toolbar">
        <input
          className="field__input accounts__cari"
          type="search"
          aria-label={ACC.cari}
          placeholder={ACC.cari}
          value={kata}
          onChange={(e) => {
            setKata(e.target.value)
          }}
        />
        <span className="toolbar__spacer" />
        <button
          type="button"
          className="btn btn--primary"
          onClick={() => {
            setPesan(null)
            setFormBuka(true)
          }}
        >
          {ACC.tambah}
        </button>
      </div>

      {pesan !== null && (
        <div className="alert alert--ok" role="status">
          {pesan}
        </div>
      )}
      {data === null && galat === null && <Memuat pesan={ACC.memuat} />}
      {galat !== null && <Gagal galat={galat} />}
      {data !== null && data.total === 0 && <Kosong pesan={cari === '' ? ACC.kosong : ACC.tidakCocok} />}
      {data !== null && data.daftar.length > 0 && (
        <>
          <table className="inbox__tabel accounts__tabel">
            <thead>
              <tr>
                <th>{ACC.kolomId}</th>
                <th>{ACC.insuredName}</th>
                <th>{ACC.groupBusiness}</th>
                <th>{ACC.owner}</th>
                <th>{ACC.createDate}</th>
              </tr>
            </thead>
            <tbody>
              {data.daftar.map((a) => (
                <tr key={a.id} className="inbox__baris">
                  <td>{atauStrip(a.idView)}</td>
                  <td>
                    {atauStrip(a.insuredName)}
                    {/* Org ID di bawah nama, bukan kolom sendiri - pola Contact ID di Kelola User. */}
                    {a.orgId !== '' && <span className="muted accounts__kecil">{a.orgId}</span>}
                  </td>
                  <td>{atauStrip(a.groupBusiness)}</td>
                  <td>{a.createOp === '' ? <span className="muted">—</span> : a.createOp}</td>
                  <td>{a.createDate === '' ? <span className="muted">—</span> : a.createDate}</td>
                </tr>
              ))}
            </tbody>
          </table>
          <Halaman halaman={data.halaman} ukuran={data.ukuran} total={data.total} onPindah={setHalaman} />
        </>
      )}

      {formBuka && (
        <FormAccount
          groupBusiness={groupBusiness}
          onTutup={() => {
            setFormBuka(false)
          }}
          onTersimpan={tersimpan}
        />
      )}
    </section>
  )
}
