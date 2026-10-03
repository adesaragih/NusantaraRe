// Halaman Marketing Officer (permintaan work owner 03-10-2026): halaman depan = daftar LEADER; dari leader dibuka
// anggotanya; setiap MO punya log perubahan dari MARKETINGOFFICER_LOG. Tampilan mengikuti Kelola User (tema di
// `marketingofficer.css`, akar `.marketingofficer__akar`). Tambah dan ubah lewat FormMO; nol hapus.

import { useCallback, useEffect, useState } from 'react'

import { Gagal, Kosong, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import { ambilDaftar, type BarisMO, type MarketingOfficer as BarisTersimpan } from '../api'
import { kelompokLeader, saring, saringLeader, tandaAkun, type Saringan } from '../aturan'
import FormMO from '../components/FormMO'
import LogMO from '../components/LogMO'
import { MO } from '../labels'

const SARINGAN: { nilai: Saringan; label: string }[] = [
  { nilai: 'semua', label: MO.saringSemua },
  { nilai: 'aktif', label: MO.saringAktif },
  { nilai: 'nonaktif', label: MO.saringNonaktif },
]

/** Tampilan halaman: daftar leader, anggota satu leader, atau MO tanpa leader. */
type Tampilan = { jenis: 'leader' } | { jenis: 'anggota'; leaderId: string } | { jenis: 'tanpa' }

/** Form terbuka: `baris` null = tambah; `leaderAwal` = tambah anggota leader itu. */
interface FormTerbuka {
  baris: BarisMO | null
  leaderAwal?: string
}

function Status({ b }: { b: BarisMO }) {
  const aktif = b.moStatus === '1'
  return (
    <span className="marketingofficer__status">
      <span className={`badge ${aktif ? 'marketingofficer__badge--aktif' : 'marketingofficer__badge--nonaktif'}`}>
        {aktif ? MO.aktif : MO.nonaktif}
      </span>
    </span>
  )
}

function SelNama({ b }: { b: BarisMO }) {
  return (
    <td>
      {b.clientName}
      <span className="muted marketingofficer__kecil">{b.clientId}</span>
    </td>
  )
}

function SelAkun({ b }: { b: BarisMO }) {
  const tanda = tandaAkun(b)
  return (
    <td>
      {b.aksesLogin === '' ? <span className="muted">—</span> : b.aksesLogin}
      {tanda !== '' && <span className="marketingofficer__tanda">{tanda}</span>}
      {b.emailAkun !== '' && <span className="muted marketingofficer__kecil">{b.emailAkun}</span>}
    </td>
  )
}

function SelCabang({ b }: { b: BarisMO }) {
  return (
    <td>
      {b.branchDetailName === '' ? <span className="muted">—</span> : b.branchDetailName}
      {b.teamGroup !== '' && <span className="muted marketingofficer__kecil">{`${MO.teamGroup} ${b.teamGroup}`}</span>}
    </td>
  )
}

export default function MarketingOfficer() {
  const [daftar, setDaftar] = useState<BarisMO[] | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [kueri, setKueri] = useState('')
  const [saringan, setSaringan] = useState<Saringan>('semua')
  const [tampilan, setTampilan] = useState<Tampilan>({ jenis: 'leader' })
  const [form, setForm] = useState<FormTerbuka | undefined>(undefined)
  const [log, setLog] = useState<BarisMO | undefined>(undefined)
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

  const pindah = (t: Tampilan) => {
    setTampilan(t)
    setKueri('')
    setPesan(null)
  }

  const kelompok = daftar === null ? null : kelompokLeader(daftar)
  const leaderTerbuka =
    tampilan.jenis === 'anggota' ? kelompok?.leader.find((k) => k.leader.id === tampilan.leaderId) : undefined
  // Leader yang hilang sesudah muat ulang: kembali ke daftar leader.
  const jenis = tampilan.jenis === 'anggota' && kelompok !== null && leaderTerbuka === undefined ? 'leader' : tampilan.jenis

  const barisAnggota =
    jenis === 'anggota' ? (leaderTerbuka?.anggota ?? []) : jenis === 'tanpa' ? (kelompok?.tanpaLeader ?? []) : []
  const tampilLeader = kelompok === null ? [] : saringLeader(kelompok.leader, kueri, saringan)
  const tampilAnggota = saring(barisAnggota, kueri, saringan)

  const tombolAksi = (b: BarisMO, denganAnggota: boolean) => (
    <span className="marketingofficer__aksi">
      {denganAnggota && (
        <button
          type="button"
          className="btn btn--ghost"
          onClick={() => {
            pindah({ jenis: 'anggota', leaderId: b.id })
          }}
        >
          {MO.anggota}
        </button>
      )}
      <button
        type="button"
        className="btn btn--ghost"
        onClick={() => {
          setLog(b)
        }}
      >
        {MO.log}
      </button>
      <button
        type="button"
        className="btn btn--ghost"
        onClick={() => {
          setPesan(null)
          setForm({ baris: b })
        }}
      >
        {MO.ubah}
      </button>
    </span>
  )

  return (
    <section className="inbox marketingofficer__akar">
      <header className="inbox__kepala">
        {jenis !== 'leader' && (
          <button
            type="button"
            className="btn btn--ghost marketingofficer__kembali"
            onClick={() => {
              pindah({ jenis: 'leader' })
            }}
          >
            {MO.kembali}
          </button>
        )}
        <h2 className="inbox__judul">
          {jenis === 'leader'
            ? MO.judul
            : jenis === 'tanpa'
              ? MO.tanpaLeader
              : MO.anggotaDari(leaderTerbuka?.leader.clientName ?? '')}
        </h2>
      </header>
      <p className="muted marketingofficer__sub">
        {jenis === 'leader'
          ? MO.sub
          : jenis === 'tanpa'
            ? MO.subTanpaLeader
            : `${MO.code} ${leaderTerbuka?.leader.id ?? ''} · ${MO.hitungAnggota(leaderTerbuka?.aktif ?? 0, barisAnggota.length)}`}
      </p>

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
        {jenis === 'leader' && kelompok !== null && kelompok.tanpaLeader.length > 0 && (
          <button
            type="button"
            className="btn btn--ghost"
            onClick={() => {
              pindah({ jenis: 'tanpa' })
            }}
          >
            {`${MO.tanpaLeader} (${kelompok.tanpaLeader.length})`}
          </button>
        )}
        <button
          type="button"
          className="btn btn--primary"
          onClick={() => {
            setPesan(null)
            setForm({ baris: null, leaderAwal: jenis === 'anggota' ? leaderTerbuka?.leader.id : undefined })
          }}
        >
          {jenis === 'anggota' ? MO.tambahAnggota : MO.tambah}
        </button>
      </div>

      {pesan !== null && (
        <div className="alert alert--ok" role="status">
          {pesan}
        </div>
      )}
      {daftar === null && galat === null && <Memuat pesan={MO.memuat} />}
      {galat !== null && <Gagal galat={galat} />}

      {jenis === 'leader' && kelompok !== null && (
        <>
          {kelompok.leader.length === 0 && <Kosong pesan={MO.kosongLeader} />}
          {kelompok.leader.length > 0 && tampilLeader.length === 0 && <Kosong pesan={MO.tidakCocok} />}
          {tampilLeader.length > 0 && (
            <table className="inbox__tabel marketingofficer__tabel">
              <thead>
                <tr>
                  <th>{MO.kolomCode}</th>
                  <th>{MO.kolomNama}</th>
                  <th>{MO.kolomAkun}</th>
                  <th>{MO.kolomSubBranch}</th>
                  <th>{MO.kolomAnggota}</th>
                  <th>{MO.kolomStatus}</th>
                  <th className="table__actions">{MO.kolomAksi}</th>
                </tr>
              </thead>
              <tbody>
                {tampilLeader.map((k) => (
                  <tr key={k.leader.id} className="inbox__baris">
                    <td>{k.leader.id}</td>
                    <SelNama b={k.leader} />
                    <SelAkun b={k.leader} />
                    <SelCabang b={k.leader} />
                    <td>{MO.hitungAnggota(k.aktif, k.anggota.length)}</td>
                    <td>
                      <Status b={k.leader} />
                    </td>
                    <td className="table__actions">{tombolAksi(k.leader, true)}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </>
      )}

      {jenis !== 'leader' && kelompok !== null && (
        <>
          {barisAnggota.length === 0 && <Kosong pesan={MO.kosongAnggota} />}
          {barisAnggota.length > 0 && tampilAnggota.length === 0 && <Kosong pesan={MO.tidakCocok} />}
          {tampilAnggota.length > 0 && (
            <table className="inbox__tabel marketingofficer__tabel">
              <thead>
                <tr>
                  <th>{MO.kolomCode}</th>
                  <th>{MO.kolomNama}</th>
                  <th>{MO.kolomAkun}</th>
                  <th>{MO.kolomSubBranch}</th>
                  <th>{MO.kolomStatus}</th>
                  <th className="table__actions">{MO.kolomAksi}</th>
                </tr>
              </thead>
              <tbody>
                {tampilAnggota.map((b) => (
                  <tr key={b.id} className="inbox__baris">
                    <td>{b.id}</td>
                    <SelNama b={b} />
                    <SelAkun b={b} />
                    <SelCabang b={b} />
                    <td>
                      <Status b={b} />
                    </td>
                    <td className="table__actions">{tombolAksi(b, false)}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </>
      )}

      {form !== undefined && (
        <FormMO
          key={form.baris?.id ?? `baru-${form.leaderAwal ?? ''}`}
          baris={form.baris}
          leaderAwal={form.leaderAwal}
          onTutup={() => {
            setForm(undefined)
          }}
          onTersimpan={tersimpan}
        />
      )}
      {log !== undefined && (
        <LogMO
          baris={log}
          onTutup={() => {
            setLog(undefined)
          }}
        />
      )}
    </section>
  )
}
