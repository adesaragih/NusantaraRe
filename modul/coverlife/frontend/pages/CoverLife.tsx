// Halaman Cover Life - section Pega `InboxCoverLife` ("COVER" b368), grid dan form di SATU layar:
// - "Cover" b848 (`.Cover` b877, pxTextInput b880, wajib b843 / b895; huruf TIDAK diubah - XML tanpa pengubah huruf);
// - "Note" b1026 (`.Note` b1055, pxTextArea b1058, tidak wajib b1020 / b1070);
// - form TANPA medan ID (hanya Cover dan Note) - ID dibentuk server; saat Edit ID ditampilkan sebagai teks status;
// - pesan `InputParam.ERRMSG` b2052 (pxDisplayText, tampil bila tidak kosong b2138) = galat simpan;
// - Save b5407 (`AddToList_Act` b5426) dan Cancel b5671 (`NewData_DT` b5694, hanya saat `DATASHOW = 'IsEdit'` b5831);
// - grid `BrowseCoverLife_RD` b4026: ID b3023 (`.ID` b3449), Cover b3161 (`.Cover` b3596), kolom tanpa judul b3277 berisi
//   Edit b3784 (`EditList_DT` b3807 mengisi form dengan ID, Cover, Note); Note TIDAK tampil di grid; ID menaik tetap
//   (b3967), kolom tidak dapat diurutkan (b3969 / b3991 / b4013), tanpa saring (b4080); 50 baris per halaman b4101,
//   halaman bernomor b4072.
// TIDAK ada Delete, Upload CSV, saring, ringkasan, maupun detail (XML tidak memuatnya). Menu View only: tanpa form dan
// Edit (backend juga menolak tulisnya).

import { useCallback, useEffect, useState } from 'react'

import { Gagal, Kosong, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import { useBolehUbah } from '../../../../inti/frontend/lib/hakMenu'
import { ambilDaftar, tambah, ubah, type Halaman } from '../api'
import { BATAS_COVER, BATAS_NOTE, jumlahHalaman, periksaIsian } from '../aturan'
import { CVL } from '../labels'
import { NAMA_CVL } from '../menu'

export default function CoverLife() {
  const bolehUbah = useBolehUbah(NAMA_CVL)
  const [halaman, setHalaman] = useState(1)
  const [data, setData] = useState<Halaman | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [segar, setSegar] = useState(0)
  const [pesan, setPesan] = useState<string | null>(null)

  // Form: `ubahID` kosong = Add; terisi = Edit (`DATASHOW = 'IsEdit'`).
  const [ubahID, setUbahID] = useState('')
  const [cover, setCover] = useState('')
  const [note, setNote] = useState('')
  const [galatForm, setGalatForm] = useState<unknown>(null)
  const [pesanForm, setPesanForm] = useState<string | null>(null)
  const [sibuk, setSibuk] = useState(false)

  const muat = useCallback((h: number) => {
    ambilDaftar(h).then(
      (d) => {
        setData(d)
        setGalat(null)
      },
      (g: unknown) => setGalat(g),
    )
  }, [])
  useEffect(() => {
    muat(halaman)
  }, [muat, halaman, segar])

  /** Cancel (`NewData_DT`): form kembali kosong (Add). */
  const kosongkanForm = () => {
    setUbahID('')
    setCover('')
    setNote('')
    setGalatForm(null)
    setPesanForm(null)
  }

  const simpan = () => {
    if (sibuk) return
    const salah = periksaIsian(cover, note)
    setPesanForm(salah)
    setGalatForm(null)
    if (salah !== null) return
    setSibuk(true)
    const isi = cover.trim()
    const janji = ubahID === '' ? tambah(isi, note) : ubah(ubahID, isi, note)
    janji.then(
      (c) => {
        setSibuk(false)
        kosongkanForm()
        setPesan(CVL.tersimpan(c.id))
        setSegar((n) => n + 1)
      },
      (g: unknown) => {
        setSibuk(false)
        setGalatForm(g)
      },
    )
  }

  const dari = data === null ? 1 : jumlahHalaman(data.total, data.ukuran)

  return (
    <section className="inbox coverlife__akar">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">{CVL.judul}</h2>
      </header>

      {bolehUbah && (
        <form
          className="coverlife__kartu coverlife__form"
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
          <div className="coverlife__baris">
            <label className="field coverlife__isian">
              <span className="field__label">
                {CVL.cover} <span aria-hidden="true">*</span>
              </span>
              <input
                className="field__input"
                type="text"
                value={cover}
                maxLength={BATAS_COVER}
                required
                onChange={(e) => setCover(e.target.value)}
              />
            </label>
            <label className="field coverlife__isian">
              <span className="field__label">{CVL.note}</span>
              <textarea
                className="field__input coverlife__teks"
                value={note}
                rows={3}
                maxLength={BATAS_NOTE}
                onChange={(e) => setNote(e.target.value)}
              />
            </label>
          </div>
          <div className="coverlife__aksi">
            <button type="submit" className="btn btn--primary" disabled={sibuk}>
              {sibuk ? CVL.menyimpan : CVL.save}
            </button>
            {ubahID !== '' && (
              <button type="button" className="btn btn--ghost" onClick={kosongkanForm}>
                {CVL.batal}
              </button>
            )}
            {ubahID !== '' && <span className="muted">{CVL.modeUbah(ubahID)}</span>}
          </div>
        </form>
      )}

      <div className="coverlife__alat">
        <span className="toolbar__spacer" />
        {data !== null && <span className="muted">{CVL.jumlah(data.total)}</span>}
      </div>

      {pesan !== null && (
        <div className="alert alert--ok" role="status">
          {pesan}
        </div>
      )}
      {galat !== null && <Gagal galat={galat} />}
      {data === null && galat === null && <Memuat pesan={CVL.memuat} />}
      {data !== null && data.daftar.length === 0 && <Kosong pesan={CVL.kosong} />}
      {data !== null && data.daftar.length > 0 && (
        <>
          <div className="coverlife__gulir">
            <table className="inbox__tabel coverlife__tabel">
              <thead>
                <tr>
                  <th>{CVL.id}</th>
                  <th>{CVL.cover}</th>
                  {bolehUbah && <th className="table__actions" aria-label={CVL.edit} />}
                </tr>
              </thead>
              <tbody>
                {data.daftar.map((c) => (
                  <tr key={c.id} className="inbox__baris">
                    <td>{c.id}</td>
                    <td className="coverlife__nama">{c.cover}</td>
                    {bolehUbah && (
                      <td className="table__actions">
                        <span className="coverlife__aksi">
                          <button
                            type="button"
                            className="btn btn--ghost btn--sm"
                            onClick={() => {
                              setPesan(null)
                              setGalatForm(null)
                              setPesanForm(null)
                              setUbahID(c.id)
                              setCover(c.cover)
                              setNote(c.note)
                            }}
                          >
                            {CVL.edit}
                          </button>
                        </span>
                      </td>
                    )}
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          <div className="coverlife__halaman">
            <button type="button" className="btn btn--ghost btn--sm" disabled={halaman <= 1} onClick={() => setHalaman((h) => h - 1)}>
              {CVL.sebelum}
            </button>
            <span className="muted">{CVL.halaman(data.halaman, dari)}</span>
            <button type="button" className="btn btn--ghost btn--sm" disabled={halaman >= dari} onClick={() => setHalaman((h) => h + 1)}>
              {CVL.sesudah}
            </button>
          </div>
        </>
      )}
    </section>
  )
}
