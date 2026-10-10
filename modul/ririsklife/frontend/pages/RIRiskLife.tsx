// Halaman R/I Risk - section Pega `InboxSummaryRIRisk` ("R/I RISK SUMMARY" b375): form (ERRMSG b903, ID disabled b1073,
// R/I RISK NAME b1225/b1256 wajib b1274), tombol Upload CSV b3156 / View Upload b3676 / Simpan Upload b4690, dan grid
// `BrowseRIRiskSummary` b7127 (filter, urut, ID menaik b9806, 50 per halaman b10030; Edit b8589 / Detail b8869 /
// Delete b9624 per baris).
// Save b1796 (`AddToListSummary_Act`) / Cancel b2063 (hanya saat Edit, `DATASHOW = 'IsEdit'` b2231): di XML wadahnya
// TERSEMBUNYI (`pyContainerVisibleWhen` `1=2` b1511) tetapi DITAMPILKAN atas keputusan work owner 08-10-2026, sama
// dengan ricommlife; Edit (`EditListSummary_DT` b8616) mengisi form lalu Save mengganti nama (ikut ke rinciannya).
// Label `Format excel` b5495 tetap tidak dirender (`1=2` b5263).
// Menu View only: tanpa form, Upload, Edit, dan Delete (backend juga menolak tulisnya); Detail tetap.

import { useCallback, useEffect, useState } from 'react'

import { Gagal, Kosong, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import { useBolehUbah } from '../../../../inti/frontend/lib/hakMenu'
import { ambilDaftar, simpanUnggah, tambah, ubah, type Halaman, type Ringkasan, type Saringan } from '../api'
import { gantiUrut, jumlahHalaman, nomorAwal, periksaNama, SARINGAN_AWAL, tandaUrut, type KolomUrut } from '../aturan'
import HasilUnggah from '../components/HasilUnggah'
import KonfirmasiHapus from '../components/KonfirmasiHapus'
import RincianDetail from '../components/RincianDetail'
import UnggahCSV, { type BerkasCSV } from '../components/UnggahCSV'
import { RK } from '../labels'
import { NAMA_RK } from '../menu'

const JEDA_CARI_MS = 300

type Popup =
  | { jenis: 'detail'; baris: Ringkasan }
  | { jenis: 'hapus'; baris: Ringkasan }
  | { jenis: 'unggah' }
  | { jenis: 'pratinjau'; isi: string }

export default function RIRiskLife() {
  const bolehUbah = useBolehUbah(NAMA_RK)
  const [ketikID, setKetikID] = useState('')
  const [ketikNama, setKetikNama] = useState('')
  const [saring, setSaring] = useState<Saringan>(SARINGAN_AWAL)
  const [data, setData] = useState<Halaman<Ringkasan> | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [segar, setSegar] = useState(0)
  const [popup, setPopup] = useState<Popup | null>(null)
  const [pesan, setPesan] = useState<string | null>(null)

  // Form ringkasan: `ubahID` kosong = Add.
  const [ubahID, setUbahID] = useState('')
  const [nama, setNama] = useState('')
  const [galatForm, setGalatForm] = useState<unknown>(null)
  const [pesanForm, setPesanForm] = useState<string | null>(null)
  const [sibuk, setSibuk] = useState(false)

  const [berkas, setBerkas] = useState<BerkasCSV | null>(null)
  const [galatUnggah, setGalatUnggah] = useState<unknown>(null)
  const [menyimpanUnggah, setMenyimpanUnggah] = useState(false)

  useEffect(() => {
    const t = setTimeout(() => setSaring((s) => ({ ...s, id: ketikID, usedby: ketikNama, halaman: 1 })), JEDA_CARI_MS)
    return () => clearTimeout(t)
  }, [ketikID, ketikNama])

  const muat = useCallback((s: Saringan) => {
    ambilDaftar(s).then(
      (d) => {
        setData(d)
        setGalat(null)
      },
      (g: unknown) => setGalat(g),
    )
  }, [])
  useEffect(() => {
    muat(saring)
  }, [muat, saring, segar])

  const kosongkanForm = () => {
    setUbahID('')
    setNama('')
    setGalatForm(null)
    setPesanForm(null)
  }

  const simpan = () => {
    if (sibuk) return
    const salah = periksaNama(nama)
    setPesanForm(salah)
    setGalatForm(null)
    if (salah !== null) return
    setSibuk(true)
    const janji = ubahID === '' ? tambah(nama.trim()) : ubah(ubahID, nama.trim())
    janji.then(
      (r) => {
        setSibuk(false)
        kosongkanForm()
        setPesan(RK.tersimpan(r.id))
        setSegar((n) => n + 1)
      },
      (g: unknown) => {
        setSibuk(false)
        setGalatForm(g)
      },
    )
  }

  const simpanBerkas = () => {
    if (berkas === null || menyimpanUnggah) return
    setMenyimpanUnggah(true)
    setGalatUnggah(null)
    simpanUnggah(berkas.isi).then(
      (h) => {
        setMenyimpanUnggah(false)
        setBerkas(null)
        setPesan(RK.unggahTersimpan(h.disimpan, h.ringkasanBaru))
        setSegar((n) => n + 1)
      },
      (g: unknown) => {
        setMenyimpanUnggah(false)
        setGalatUnggah(g)
      },
    )
  }

  const tersaring = saring.id.trim() !== '' || saring.usedby.trim() !== ''
  const dari = data === null ? 1 : jumlahHalaman(data.total, data.ukuran)
  const kepala = (kolom: KolomUrut, label: string) => (
    <th>
      <button type="button" className="ririsklife__urut" onClick={() => setSaring((s) => gantiUrut(s, kolom))}>
        {label}
        {tandaUrut(saring, kolom)}
      </button>
    </th>
  )

  return (
    <section className="inbox ririsklife__akar">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">{RK.judul}</h2>
      </header>

      {bolehUbah && (
        <form
          className="ririsklife__kartu ririsklife__form"
          onSubmit={(e) => {
            e.preventDefault()
            simpan()
          }}
        >
          {galatForm !== null && <Gagal galat={galatForm} />}
          {pesanForm !== null && (
            <div className="alert alert--warn" role="alert">
              {pesanForm}
            </div>
          )}
          <div className="ririsklife__baris">
            <label className="field">
              <span className="field__label">{RK.id}</span>
              <input className="field__input field__input--readonly" value={ubahID} placeholder={RK.idOtomatis} readOnly disabled />
            </label>
            <label className="field">
              <span className="field__label">
                {RK.nama} <span aria-hidden="true">*</span>
              </span>
              <input className="field__input" value={nama} maxLength={200} required onChange={(e) => setNama(e.target.value)} />
            </label>
          </div>
          <div className="ririsklife__aksi">
            <button type="submit" className="btn btn--primary" disabled={sibuk}>
              {sibuk ? RK.menyimpan : RK.save}
            </button>
            {ubahID !== '' && (
              <button type="button" className="btn btn--ghost" onClick={kosongkanForm}>
                {RK.batal}
              </button>
            )}
            {ubahID !== '' && <span className="muted">{RK.modeUbah(ubahID)}</span>}
          </div>
        </form>
      )}

      {bolehUbah && (
        <div className="ririsklife__kartu">
          <div className="ririsklife__aksi">
            <button type="button" className="btn btn--ghost" onClick={() => setPopup({ jenis: 'unggah' })}>
              {RK.unggah}
            </button>
            <button
              type="button"
              className="btn btn--ghost"
              disabled={berkas === null}
              onClick={() => berkas !== null && setPopup({ jenis: 'pratinjau', isi: berkas.isi })}
            >
              {RK.lihatUnggah}
            </button>
            <button
              type="button"
              className="btn btn--primary"
              disabled={berkas === null || menyimpanUnggah}
              onClick={simpanBerkas}
            >
              {menyimpanUnggah ? RK.menyimpan : RK.simpanUnggah}
            </button>
          </div>
          <p className="muted ririsklife__catatan">
            {berkas === null ? RK.belumAdaBerkas : RK.berkasDipilih(berkas.nama, Math.ceil(berkas.ukuran / 1024))}
          </p>
          {galatUnggah !== null && <Gagal galat={galatUnggah} />}
        </div>
      )}

      <div className="ririsklife__alat">
        <input
          className="field__input ririsklife__saring"
          type="search"
          placeholder={RK.saringID}
          aria-label={RK.saringID}
          value={ketikID}
          onChange={(e) => setKetikID(e.target.value)}
        />
        <input
          className="field__input ririsklife__cari"
          type="search"
          placeholder={RK.saringNama}
          aria-label={RK.saringNama}
          value={ketikNama}
          onChange={(e) => setKetikNama(e.target.value)}
        />
        <span className="toolbar__spacer" />
        {data !== null && <span className="muted">{RK.jumlah(data.total)}</span>}
      </div>

      {pesan !== null && (
        <div className="alert alert--ok" role="status">
          {pesan}
        </div>
      )}
      {galat !== null && <Gagal galat={galat} />}
      {data === null && galat === null && <Memuat pesan={RK.memuat} />}
      {data !== null && data.daftar.length === 0 && <Kosong pesan={tersaring ? RK.tidakCocok : RK.kosong} />}
      {data !== null && data.daftar.length > 0 && (
        <>
          <div className="ririsklife__gulir">
            <table className="inbox__tabel ririsklife__tabel">
              <thead>
                <tr>
                  <th>No</th>
                  {kepala('id', RK.id)}
                  {kepala('usedby', RK.nama)}
                  {kepala('operatorid', RK.operator)}
                  {kepala('modifieddate', RK.tanggal)}
                  <th className="table__actions">{RK.aksi}</th>
                </tr>
              </thead>
              <tbody>
                {data.daftar.map((r, i) => (
                  <tr key={r.id} className="inbox__baris">
                    <td>{nomorAwal(data.halaman, data.ukuran) + i}</td>
                    <td>{r.id}</td>
                    <td>{r.usedby}</td>
                    <td>{r.operatorId}</td>
                    <td>{r.diubah}</td>
                    <td className="table__actions">
                      <span className="ririsklife__aksi">
                        {bolehUbah && (
                          <button
                            type="button"
                            className="btn btn--ghost btn--sm"
                            onClick={() => {
                              setPesan(null)
                              setGalatForm(null)
                              setPesanForm(null)
                              setUbahID(r.id)
                              setNama(r.usedby)
                            }}
                          >
                            {RK.edit}
                          </button>
                        )}
                        <button
                          type="button"
                          className="btn btn--ghost btn--sm"
                          onClick={() => setPopup({ jenis: 'detail', baris: r })}
                        >
                          {RK.detail}
                        </button>
                        {bolehUbah && (
                          <button
                            type="button"
                            className="btn btn--ghost btn--sm ririsklife__hapus"
                            onClick={() => {
                              setPesan(null)
                              setPopup({ jenis: 'hapus', baris: r })
                            }}
                          >
                            {RK.hapus}
                          </button>
                        )}
                      </span>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          <div className="ririsklife__halaman">
            <button
              type="button"
              className="btn btn--ghost btn--sm"
              disabled={saring.halaman <= 1}
              onClick={() => setSaring((s) => ({ ...s, halaman: s.halaman - 1 }))}
            >
              {RK.sebelum}
            </button>
            <span className="muted">{RK.halaman(data.halaman, dari)}</span>
            <button
              type="button"
              className="btn btn--ghost btn--sm"
              disabled={saring.halaman >= dari}
              onClick={() => setSaring((s) => ({ ...s, halaman: s.halaman + 1 }))}
            >
              {RK.sesudah}
            </button>
          </div>
        </>
      )}

      {popup?.jenis === 'detail' && (
        <RincianDetail ringkasan={popup.baris} bolehUbah={bolehUbah} onTutup={() => setPopup(null)} />
      )}
      {popup?.jenis === 'hapus' && (
        <KonfirmasiHapus
          ringkasan={popup.baris}
          onTutup={() => setPopup(null)}
          onTerhapus={(id, n) => {
            setPopup(null)
            if (ubahID === id) kosongkanForm()
            setPesan(RK.terhapus(id, n))
            setSegar((x) => x + 1)
          }}
        />
      )}
      {popup?.jenis === 'unggah' && (
        <UnggahCSV
          onTutup={() => setPopup(null)}
          onDipilih={(b) => {
            setPopup(null)
            setGalatUnggah(null)
            setPesan(null)
            setBerkas(b)
          }}
        />
      )}
      {popup?.jenis === 'pratinjau' && <HasilUnggah isi={popup.isi} onTutup={() => setPopup(null)} />}
    </section>
  )
}
