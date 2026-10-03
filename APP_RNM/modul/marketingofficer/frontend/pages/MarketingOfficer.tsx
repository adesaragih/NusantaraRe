// Halaman Marketing Officer - daftar seluruh baris `POOLDATA.MARKETINGOFFICER` dengan tombol tambah dan ubah
// (perintah work owner 03-10-2026: modul baru untuk insert/update tabel itu). Nol hapus: nonaktif = Active tidak
// dicentang.

import { useCallback, useEffect, useState } from 'react'

import { Gagal, Kosong, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import { ambilDaftar, type BarisMO, type MarketingOfficer as BarisTersimpan } from '../api'
import { NILAI_LEADER, saring, tandaAkun, type Saringan } from '../aturan'
import FormMO from '../components/FormMO'
import { MO } from '../labels'

const SARINGAN: { nilai: Saringan; label: string }[] = [
  { nilai: 'semua', label: MO.saringSemua },
  { nilai: 'aktif', label: MO.saringAktif },
  { nilai: 'nonaktif', label: MO.saringNonaktif },
]

export default function MarketingOfficer() {
  const [daftar, setDaftar] = useState<BarisMO[] | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [kueri, setKueri] = useState('')
  const [saringan, setSaringan] = useState<Saringan>('semua')
  // `undefined` = form tertutup, `null` = tambah.
  const [form, setForm] = useState<BarisMO | null | undefined>(undefined)
  const [pesan, setPesan] = useState<string | null>(null)

  const muat = useCallback(() => {
    ambilDaftar().then(
      (d) => {
        setDaftar(d.daftar)
        setGalat(null)
      },
      (g: unknown) => {
        setGalat(g)
      },
    )
  }, [])

  useEffect(() => {
    muat()
  }, [muat])

  const tersimpan = (m: BarisTersimpan) => {
    setForm(undefined)
    setPesan(MO.tersimpan(m.id))
    muat()
  }

  const tampil = daftar === null ? [] : saring(daftar, kueri, saringan)

  return (
    <section className="inbox">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">{MO.judul}</h2>
      </header>
      <p className="muted marketingofficer__sub">{MO.sub}</p>

      <div className="toolbar">
        <input
          className="field__input marketingofficer__cari"
          type="search"
          aria-label={MO.cari}
          placeholder={MO.cari}
          value={kueri}
          onChange={(e) => {
            setKueri(e.target.value)
          }}
        />
        <span className="marketingofficer__saring" role="group">
          {SARINGAN.map((s) => (
            <button
              key={s.nilai}
              type="button"
              className={saringan === s.nilai ? 'btn btn--primary' : 'btn btn--ghost'}
              aria-pressed={saringan === s.nilai}
              onClick={() => {
                setSaringan(s.nilai)
              }}
            >
              {s.label}
            </button>
          ))}
        </span>
        <span className="toolbar__spacer" />
        <button
          type="button"
          className="btn btn--primary"
          onClick={() => {
            setPesan(null)
            setForm(null)
          }}
        >
          {MO.tambah}
        </button>
      </div>

      {pesan !== null && (
        <div className="alert alert--ok" role="status">
          {pesan}
        </div>
      )}
      {daftar === null && galat === null && <Memuat pesan={MO.memuat} />}
      {galat !== null && <Gagal galat={galat} />}
      {daftar !== null && daftar.length === 0 && <Kosong pesan={MO.kosong} />}
      {daftar !== null && daftar.length > 0 && tampil.length === 0 && <Kosong pesan={MO.tidakCocok} />}

      {tampil.length > 0 && (
        <div className="marketingofficer__tabel">
          <table className="inbox__tabel">
            <thead>
              <tr>
                <th>{MO.kolomCode}</th>
                <th>{MO.kolomNama}</th>
                <th>{MO.kolomAkun}</th>
                <th>{MO.kolomLeader}</th>
                <th>{MO.kolomSubBranch}</th>
                <th>{MO.kolomStatus}</th>
                <th className="table__actions">{MO.kolomAksi}</th>
              </tr>
            </thead>
            <tbody>
              {tampil.map((b) => {
                const tanda = tandaAkun(b)
                return (
                  <tr key={b.id} className="inbox__baris">
                    <td>{b.id}</td>
                    <td>
                      {b.clientName}
                      <span className="muted marketingofficer__kecil">{b.clientId}</span>
                    </td>
                    <td>
                      {b.aksesLogin === '' ? <span className="muted">—</span> : b.aksesLogin}
                      {tanda !== '' && <span className="marketingofficer__tanda">{tanda}</span>}
                      {b.emailAkun !== '' && <span className="muted marketingofficer__kecil">{b.emailAkun}</span>}
                    </td>
                    <td>{b.clientId2 === NILAI_LEADER ? <strong>{MO.adalahLeader}</strong> : b.moLeader}</td>
                    <td>
                      {b.branchDetailName === '' ? <span className="muted">—</span> : b.branchDetailName}
                      {b.teamGroup !== '' && (
                        <span className="muted marketingofficer__kecil">{`${MO.teamGroup} ${b.teamGroup}`}</span>
                      )}
                    </td>
                    <td>{b.moStatus === '1' ? MO.aktif : MO.nonaktif}</td>
                    <td className="table__actions">
                      <button
                        type="button"
                        className="btn btn--ghost"
                        onClick={() => {
                          setPesan(null)
                          setForm(b)
                        }}
                      >
                        {MO.ubah}
                      </button>
                    </td>
                  </tr>
                )
              })}
            </tbody>
          </table>
        </div>
      )}

      {form !== undefined && (
        <FormMO
          key={form?.id ?? ''}
          baris={form}
          onTutup={() => {
            setForm(undefined)
          }}
          onTersimpan={tersimpan}
        />
      )}
    </section>
  )
}
