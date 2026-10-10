// Halaman R/I Comm Life - section Pega `InboxSummaryRIComm` ("R/I COMM SUMMARY" b382): form (ERRMSG b915, ID disabled
// b1109, USEDBY wajib b1289, Save b1799 / Cancel b2092 hanya saat Edit b2261), tombol Upload CSV / View Upload / Simpan
// Upload, dan grid `BrowseRICommSummary` (filter, urut, ID menaik, 50 per halaman; Edit / Detail / Delete per baris).
// Menu View only: tanpa form, Upload, Edit, dan Delete (backend juga menolak tulisnya); Detail tetap.

import { useCallback, useEffect, useState } from 'react'

import { Gagal, Kosong, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import { useBolehUbah } from '../../../../inti/frontend/lib/hakMenu'
import { ambilDaftar, simpanUnggah, tambah, ubah, type Halaman, type Ringkasan, type Saringan } from '../api'
import { gantiUrut, jumlahHalaman, nomorAwal, periksaNama, SARINGAN_AWAL, tandaUrut, type KolomUrut } from '../aturan'
import HasilUnggah from '../components/HasilUnggah'
import KonfirmasiHapus from '../components/KonfirmasiHapus'
import KomisiDetail from '../components/KomisiDetail'
import UnggahCSV, { type BerkasCSV } from '../components/UnggahCSV'
import { RC } from '../labels'
import { NAMA_RC } from '../menu'

const JEDA_CARI_MS = 300

type Popup =
  | { jenis: 'detail'; baris: Ringkasan }
  | { jenis: 'hapus'; baris: Ringkasan }
  | { jenis: 'unggah' }
  | { jenis: 'pratinjau'; isi: string }

export default function RICommLife() {
  const bolehUbah = useBolehUbah(NAMA_RC)
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
        setPesan(RC.tersimpan(r.id))
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
        setPesan(RC.unggahTersimpan(h.disimpan, h.ringkasanBaru))
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
      <button type="button" className="ricommlife__urut" onClick={() => setSaring((s) => gantiUrut(s, kolom))}>
        {label}
        {tandaUrut(saring, kolom)}
      </button>
    </th>
  )

  return (
    <section className="inbox ricommlife__akar">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">{RC.judul}</h2>
      </header>

      {bolehUbah && (
        <form
          className="ricommlife__kartu ricommlife__form"
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
          <div className="ricommlife__baris">
            <label className="field">
              <span className="field__label">{RC.id}</span>
              <input
                className="field__input field__input--readonly"
                value={ubahID}
                placeholder={RC.idOtomatis}
                readOnly
                disabled
              />
            </label>
            <label className="field">
              <span className="field__label">
                {RC.usedby} <span aria-hidden="true">*</span>
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
          <div className="ricommlife__aksi">
            <button type="submit" className="btn btn--primary" disabled={sibuk}>
              {sibuk ? RC.menyimpan : RC.save}
            </button>
            {ubahID !== '' && (
              <button type="button" className="btn btn--ghost" onClick={kosongkanForm}>
                {RC.batal}
              </button>
            )}
            {ubahID !== '' && <span className="muted">{RC.modeUbah(ubahID)}</span>}
          </div>
        </form>
      )}

      {bolehUbah && (
        <div className="ricommlife__kartu">
          <div className="ricommlife__aksi">
            <button type="button" className="btn btn--ghost" onClick={() => setPopup({ jenis: 'unggah' })}>
              {RC.unggah}
            </button>
            <button
              type="button"
              className="btn btn--ghost"
              disabled={berkas === null}
              onClick={() => berkas !== null && setPopup({ jenis: 'pratinjau', isi: berkas.isi })}
            >
              {RC.lihatUnggah}
            </button>
            <button
              type="button"
              className="btn btn--primary"
              disabled={berkas === null || menyimpanUnggah}
              onClick={simpanBerkas}
            >
              {menyimpanUnggah ? RC.menyimpan : RC.simpanUnggah}
            </button>
          </div>
          <p className="muted ricommlife__catatan">
            {berkas === null ? RC.belumAdaBerkas : RC.berkasDipilih(berkas.nama, Math.ceil(berkas.ukuran / 1024))}
          </p>
          <p className="muted ricommlife__catatan">{RC.format}</p>
          {galatUnggah !== null && <Gagal galat={galatUnggah} />}
        </div>
      )}

      <div className="ricommlife__alat">
        <input
          className="field__input ricommlife__saring"
          type="search"
          placeholder={RC.saringID}
          aria-label={RC.saringID}
          value={ketikID}
          onChange={(e) => setKetikID(e.target.value)}
        />
        <input
          className="field__input ricommlife__cari"
          type="search"
          placeholder={RC.saringNama}
          aria-label={RC.saringNama}
          value={ketikNama}
          onChange={(e) => setKetikNama(e.target.value)}
        />
        <span className="toolbar__spacer" />
        {data !== null && <span className="muted">{RC.jumlah(data.total)}</span>}
      </div>

      {pesan !== null && (
        <div className="alert alert--ok" role="status">
          {pesan}
        </div>
      )}
      {galat !== null && <Gagal galat={galat} />}
      {data === null && galat === null && <Memuat pesan={RC.memuat} />}
      {data !== null && data.daftar.length === 0 && <Kosong pesan={tersaring ? RC.tidakCocok : RC.kosong} />}
      {data !== null && data.daftar.length > 0 && (
        <>
          <div className="ricommlife__gulir">
            <table className="inbox__tabel ricommlife__tabel">
              <thead>
                <tr>
                  <th>No</th>
                  {kepala('id', RC.id)}
                  {kepala('usedby', RC.nama)}
                  {kepala('operatorid', RC.operator)}
                  {kepala('modifieddate', RC.tanggal)}
                  <th className="table__actions">{RC.aksi}</th>
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
                      <span className="ricommlife__aksi">
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
                            {RC.edit}
                          </button>
                        )}
                        <button
                          type="button"
                          className="btn btn--ghost btn--sm"
                          onClick={() => setPopup({ jenis: 'detail', baris: r })}
                        >
                          {RC.detail}
                        </button>
                        {bolehUbah && (
                          <button
                            type="button"
                            className="btn btn--ghost btn--sm ricommlife__hapus"
                            onClick={() => {
                              setPesan(null)
                              setPopup({ jenis: 'hapus', baris: r })
                            }}
                          >
                            {RC.hapus}
                          </button>
                        )}
                      </span>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          <div className="ricommlife__halaman">
            <button
              type="button"
              className="btn btn--ghost btn--sm"
              disabled={saring.halaman <= 1}
              onClick={() => setSaring((s) => ({ ...s, halaman: s.halaman - 1 }))}
            >
              {RC.sebelum}
            </button>
            <span className="muted">{RC.halaman(data.halaman, dari)}</span>
            <button
              type="button"
              className="btn btn--ghost btn--sm"
              disabled={saring.halaman >= dari}
              onClick={() => setSaring((s) => ({ ...s, halaman: s.halaman + 1 }))}
            >
              {RC.sesudah}
            </button>
          </div>
        </>
      )}

      {popup?.jenis === 'detail' && (
        <KomisiDetail ringkasan={popup.baris} bolehUbah={bolehUbah} onTutup={() => setPopup(null)} />
      )}
      {popup?.jenis === 'hapus' && (
        <KonfirmasiHapus
          ringkasan={popup.baris}
          onTutup={() => setPopup(null)}
          onTerhapus={(id, n) => {
            setPopup(null)
            if (ubahID === id) kosongkanForm()
            setPesan(RC.terhapus(id, n))
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
