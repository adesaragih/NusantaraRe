// Halaman R/I Rate Life - section Pega `InboxSummaryRIRate` ("R/I RATE SUMMARY" b382): form (ID disabled b1104,
// R/I RATE NAME wajib b1285, Save b1809 / Cancel b2109 - container `1=2` di Pega, DITAMPILKAN atas perintah CRUD work
// owner 05-10-2026), tombol Upload CSV / View Upload / Simpan Upload, dan grid `BrowseRateLifeSummary` (filter, urut,
// 50 per halaman; Edit / Detail / Delete per baris). Menu View only: tanpa form, Upload, Edit, dan Delete (backend
// juga menolak tulisnya); Detail tetap.

import { useCallback, useEffect, useState } from 'react'

import { Gagal, Kosong, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import { useBolehUbah } from '../../../../inti/frontend/lib/hakMenu'
import { ambilDaftar, simpanUnggah, tambah, ubah, type Halaman, type Ringkasan, type Saringan } from '../api'
import { gantiUrut, jumlahHalaman, nomorAwal, periksaNama, SARINGAN_AWAL, tandaUrut, type KolomUrut } from '../aturan'
import HasilUnggah from '../components/HasilUnggah'
import KonfirmasiHapus from '../components/KonfirmasiHapus'
import RateDetail from '../components/RateDetail'
import UnggahCSV, { type BerkasCSV } from '../components/UnggahCSV'
import { RR } from '../labels'
import { NAMA_RR } from '../menu'

const JEDA_CARI_MS = 300

type Popup =
  | { jenis: 'detail'; baris: Ringkasan }
  | { jenis: 'hapus'; baris: Ringkasan }
  | { jenis: 'unggah' }
  | { jenis: 'pratinjau'; isi: string }

export default function RIRateLife() {
  const bolehUbah = useBolehUbah(NAMA_RR)
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
        setPesan(RR.tersimpan(r.id))
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
        setPesan(RR.unggahTersimpan(h.disimpan, h.ringkasanBaru))
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
      <button type="button" className="riratelife__urut" onClick={() => setSaring((s) => gantiUrut(s, kolom))}>
        {label}
        {tandaUrut(saring, kolom)}
      </button>
    </th>
  )

  return (
    <section className="inbox riratelife__akar">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">{RR.judul}</h2>
      </header>

      {bolehUbah && (
        <form
          className="riratelife__kartu riratelife__form"
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
          <div className="riratelife__baris">
            <label className="field">
              <span className="field__label">{RR.id}</span>
              <input
                className="field__input field__input--readonly"
                value={ubahID}
                placeholder={RR.idOtomatis}
                readOnly
                disabled
              />
            </label>
            <label className="field">
              <span className="field__label">
                {RR.nama} <span aria-hidden="true">*</span>
              </span>
              <input
                className="field__input"
                value={nama}
                maxLength={200}
                required
                onChange={(e) => setNama(e.target.value)}
              />
            </label>
          </div>
          <div className="riratelife__aksi">
            <button type="submit" className="btn btn--primary" disabled={sibuk}>
              {sibuk ? RR.menyimpan : RR.save}
            </button>
            {ubahID !== '' && (
              <button type="button" className="btn btn--ghost" onClick={kosongkanForm}>
                {RR.batal}
              </button>
            )}
            {ubahID !== '' && <span className="muted">{RR.modeUbah(ubahID)}</span>}
          </div>
        </form>
      )}

      {bolehUbah && (
        <div className="riratelife__kartu">
          <div className="riratelife__aksi">
            <button type="button" className="btn btn--ghost" onClick={() => setPopup({ jenis: 'unggah' })}>
              {RR.unggah}
            </button>
            <button
              type="button"
              className="btn btn--ghost"
              disabled={berkas === null}
              onClick={() => berkas !== null && setPopup({ jenis: 'pratinjau', isi: berkas.isi })}
            >
              {RR.lihatUnggah}
            </button>
            <button
              type="button"
              className="btn btn--primary"
              disabled={berkas === null || menyimpanUnggah}
              onClick={simpanBerkas}
            >
              {menyimpanUnggah ? RR.menyimpan : RR.simpanUnggah}
            </button>
          </div>
          <p className="muted riratelife__catatan">
            {berkas === null ? RR.belumAdaBerkas : RR.berkasDipilih(berkas.nama, Math.ceil(berkas.ukuran / 1024))}
          </p>
          <p className="muted riratelife__catatan">{RR.format}</p>
          {galatUnggah !== null && <Gagal galat={galatUnggah} />}
        </div>
      )}

      <div className="riratelife__alat">
        <input
          className="field__input riratelife__saring"
          type="search"
          placeholder={RR.saringID}
          aria-label={RR.saringID}
          value={ketikID}
          onChange={(e) => setKetikID(e.target.value)}
        />
        <input
          className="field__input riratelife__cari"
          type="search"
          placeholder={RR.saringNama}
          aria-label={RR.saringNama}
          value={ketikNama}
          onChange={(e) => setKetikNama(e.target.value)}
        />
        <span className="toolbar__spacer" />
        {data !== null && <span className="muted">{RR.jumlah(data.total)}</span>}
      </div>

      {pesan !== null && (
        <div className="alert alert--ok" role="status">
          {pesan}
        </div>
      )}
      {galat !== null && <Gagal galat={galat} />}
      {data === null && galat === null && <Memuat pesan={RR.memuat} />}
      {data !== null && data.daftar.length === 0 && <Kosong pesan={tersaring ? RR.tidakCocok : RR.kosong} />}
      {data !== null && data.daftar.length > 0 && (
        <>
          <div className="riratelife__gulir">
            <table className="inbox__tabel riratelife__tabel">
              <thead>
                <tr>
                  <th>No</th>
                  {kepala('id', RR.id)}
                  {kepala('usedby', RR.nama)}
                  {kepala('operatorid', RR.operator)}
                  {kepala('modifieddate', RR.tanggal)}
                  <th className="table__actions">{RR.aksi}</th>
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
                      <span className="riratelife__aksi">
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
                            {RR.edit}
                          </button>
                        )}
                        <button
                          type="button"
                          className="btn btn--ghost btn--sm"
                          onClick={() => setPopup({ jenis: 'detail', baris: r })}
                        >
                          {RR.detail}
                        </button>
                        {bolehUbah && (
                          <button
                            type="button"
                            className="btn btn--ghost btn--sm riratelife__hapus"
                            onClick={() => {
                              setPesan(null)
                              setPopup({ jenis: 'hapus', baris: r })
                            }}
                          >
                            {RR.hapus}
                          </button>
                        )}
                      </span>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          <div className="riratelife__halaman">
            <button
              type="button"
              className="btn btn--ghost btn--sm"
              disabled={saring.halaman <= 1}
              onClick={() => setSaring((s) => ({ ...s, halaman: s.halaman - 1 }))}
            >
              {RR.sebelum}
            </button>
            <span className="muted">{RR.halaman(data.halaman, dari)}</span>
            <button
              type="button"
              className="btn btn--ghost btn--sm"
              disabled={saring.halaman >= dari}
              onClick={() => setSaring((s) => ({ ...s, halaman: s.halaman + 1 }))}
            >
              {RR.sesudah}
            </button>
          </div>
        </>
      )}

      {popup?.jenis === 'detail' && <RateDetail ringkasan={popup.baris} bolehUbah={bolehUbah} onTutup={() => setPopup(null)} />}
      {popup?.jenis === 'hapus' && (
        <KonfirmasiHapus
          ringkasan={popup.baris}
          onTutup={() => setPopup(null)}
          onTerhapus={(id, n) => {
            setPopup(null)
            if (ubahID === id) kosongkanForm()
            setPesan(RR.terhapus(id, n))
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
