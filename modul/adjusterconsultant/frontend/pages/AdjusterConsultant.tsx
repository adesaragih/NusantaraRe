// Halaman Adjuster Consultant (perintah work owner 05-10-2026) - padanan layar master Pega `MstAdjusterConsultant`
// (Claim Fac In / Claim Prop): daftar ID, Name, Address, Telp No, Edit Date dengan Add dan Edit. Tambahan keputusan work
// owner: cari, saringan status (bawaan Active), View, dan Activate / Deactivate pengganti Delete (Delete Pega membuka
// konfirmasi hapus Treaty Group - salin-tempel - jadi adjuster tidak pernah terhapus). Menu View only: tanpa Add, Edit,
// Activate / Deactivate (backend juga menolak tulisnya).

import { useCallback, useEffect, useState } from 'react'

import { Gagal, Kosong, Memuat, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { useBolehUbah } from '../../../../inti/frontend/lib/hakMenu'
import { ambilDaftar, setelAktif, type Adjuster, type Status } from '../api'
import { namaTampil, STATUS_AWAL, URUTAN_STATUS } from '../aturan'
import FormAdjuster, { type ModeForm } from '../components/FormAdjuster'
import { ADJ } from '../labels'
import { NAMA_ADJ } from '../menu'

const JEDA_CARI_MS = 300

const LABEL_STATUS: Record<Status, string> = { active: ADJ.aktif, inactive: ADJ.nonaktif, '': ADJ.semua }

export default function AdjusterConsultant() {
  const bolehUbah = useBolehUbah(NAMA_ADJ)
  const [kata, setKata] = useState('')
  const [cari, setCari] = useState('')
  const [status, setStatus] = useState<Status>(STATUS_AWAL)
  const [daftar, setDaftar] = useState<Adjuster[] | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [segar, setSegar] = useState(0)
  const [form, setForm] = useState<ModeForm | null>(null)
  const [pesan, setPesan] = useState<string | null>(null)
  const [akanUbahStatus, setAkanUbahStatus] = useState<Adjuster | null>(null)
  const [memproses, setMemproses] = useState(false)
  const [galatStatus, setGalatStatus] = useState<unknown>(null)

  useEffect(() => {
    const t = setTimeout(() => setCari(kata), JEDA_CARI_MS)
    return () => clearTimeout(t)
  }, [kata])

  const muat = useCallback((q: string, s: Status) => {
    ambilDaftar(q, s).then(
      (d) => {
        setDaftar(d.daftar)
        setGalat(null)
      },
      (g: unknown) => setGalat(g),
    )
  }, [])
  useEffect(() => {
    muat(cari, status)
  }, [muat, cari, status, segar])

  const jalankanStatus = () => {
    if (akanUbahStatus === null || memproses) return
    const target = akanUbahStatus
    setMemproses(true)
    setGalatStatus(null)
    setelAktif(target.id, !target.active).then(
      (a) => {
        setMemproses(false)
        setAkanUbahStatus(null)
        setPesan(a.active ? ADJ.diaktifkan(a.id) : ADJ.dinonaktifkan(a.id))
        setSegar((n) => n + 1)
      },
      (g: unknown) => {
        setMemproses(false)
        setGalatStatus(g)
      },
    )
  }

  return (
    <section className="inbox adjusterconsultant__akar">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">{ADJ.judul}</h2>
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
            {ADJ.tambah}
          </button>
        )}
      </header>

      <div className="adjusterconsultant__alat">
        <input
          className="field__input adjusterconsultant__cari"
          type="search"
          placeholder={ADJ.cari}
          aria-label={ADJ.cari}
          value={kata}
          onChange={(e) => setKata(e.target.value)}
        />
        <div className="adjusterconsultant__segmen" role="radiogroup" aria-label={ADJ.status}>
          {URUTAN_STATUS.map((s) => (
            <button
              key={s === '' ? 'semua' : s}
              type="button"
              role="radio"
              aria-checked={status === s}
              className={
                status === s
                  ? 'adjusterconsultant__segmen-butir adjusterconsultant__segmen-butir--aktif'
                  : 'adjusterconsultant__segmen-butir'
              }
              onClick={() => setStatus(s)}
            >
              {LABEL_STATUS[s]}
            </button>
          ))}
        </div>
        <span className="toolbar__spacer" />
        {daftar !== null && <span className="muted">{ADJ.jumlah(daftar.length)}</span>}
      </div>

      {pesan !== null && (
        <div className="alert alert--ok" role="status">
          {pesan}
        </div>
      )}
      {galat !== null && <Gagal galat={galat} />}
      {daftar === null && galat === null && <Memuat pesan={ADJ.memuat} />}
      {daftar !== null && daftar.length === 0 && <Kosong pesan={cari.trim() === '' && status === '' ? ADJ.kosong : ADJ.tidakCocok} />}
      {daftar !== null && daftar.length > 0 && (
        <div className="adjusterconsultant__gulir">
          <table className="inbox__tabel adjusterconsultant__tabel">
            <thead>
              <tr>
                <th>{ADJ.no}</th>
                <th>{ADJ.id}</th>
                <th>{ADJ.nama}</th>
                <th>{ADJ.alamat}</th>
                <th>{ADJ.telp}</th>
                <th>{ADJ.pengubah}</th>
                <th>{ADJ.tglUbah}</th>
                <th className="adjusterconsultant__kolom-status">{ADJ.status}</th>
                <th className="table__actions">{ADJ.aksi}</th>
              </tr>
            </thead>
            <tbody>
              {daftar.map((a, i) => (
                <tr key={a.id} className={a.active ? 'inbox__baris' : 'inbox__baris adjusterconsultant__baris--nonaktif'}>
                  <td>{i + 1}</td>
                  <td>{a.id}</td>
                  <td>{a.name}</td>
                  <td className="adjusterconsultant__alamat">{a.address}</td>
                  <td>{a.telpNo}</td>
                  <td>{a.username}</td>
                  <td>{a.editDate}</td>
                  <td className="adjusterconsultant__kolom-status">
                    <span
                      className={
                        a.active ? 'adjusterconsultant__status' : 'adjusterconsultant__status adjusterconsultant__status--nonaktif'
                      }
                    >
                      {a.active ? ADJ.aktif : ADJ.nonaktif}
                    </span>
                  </td>
                  <td className="table__actions">
                    <span className="adjusterconsultant__aksi">
                      <button type="button" className="btn btn--ghost btn--sm" onClick={() => setForm({ jenis: 'lihat', baris: a })}>
                        {ADJ.view}
                      </button>
                      {bolehUbah && (
                        <>
                          <button
                            type="button"
                            className="btn btn--ghost btn--sm"
                            onClick={() => {
                              setPesan(null)
                              setForm({ jenis: 'ubah', baris: a })
                            }}
                          >
                            {ADJ.edit}
                          </button>
                          <button
                            type="button"
                            className={a.active ? 'btn btn--ghost btn--sm adjusterconsultant__hapus' : 'btn btn--ghost btn--sm'}
                            onClick={() => {
                              setPesan(null)
                              setGalatStatus(null)
                              setAkanUbahStatus(a)
                            }}
                          >
                            {a.active ? ADJ.deactivate : ADJ.activate}
                          </button>
                        </>
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
        <FormAdjuster
          mode={form}
          onTutup={() => setForm(null)}
          onTersimpan={(a) => {
            setForm(null)
            setPesan(ADJ.tersimpan(a.id))
            setSegar((n) => n + 1)
          }}
        />
      )}

      {akanUbahStatus !== null && (
        <Modal
          judul={akanUbahStatus.active ? ADJ.judulNonaktif : ADJ.judulAktif}
          onTutup={() => setAkanUbahStatus(null)}
          onKirim={jalankanStatus}
          labelBatal={ADJ.batal}
          aksi={
            <button type="submit" className={akanUbahStatus.active ? 'btn btn--danger' : 'btn btn--primary'} disabled={memproses}>
              {memproses ? ADJ.memproses : akanUbahStatus.active ? ADJ.deactivate : ADJ.activate}
            </button>
          }
        >
          {galatStatus !== null && <Gagal galat={galatStatus} />}
          <p>
            {akanUbahStatus.active ? ADJ.kalimatNonaktif(namaTampil(akanUbahStatus)) : ADJ.kalimatAktif(namaTampil(akanUbahStatus))}
          </p>
        </Modal>
      )}
    </section>
  )
}
