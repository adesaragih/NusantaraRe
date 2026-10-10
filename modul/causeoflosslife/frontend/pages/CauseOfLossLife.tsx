// Halaman Cause Of Loss Life - section Pega `InboxCauseofLossLife` ("CAUSE OF LOSS" b343), grid dan form di SATU layar:
// - "ID" b823 (`.ID` b852, pxTextInput DISABLED b871-b872) - kosong = dibentuk server;
// - "Cause of Loss" b1001 (`.CauseofLoss` b1027, pxTextInput b1029, wajib b996 / b1044; huruf TIDAK diubah - XML tanpa
//   pengubah huruf);
// - pesan `InputParam.ERRMSG` b2018 (pxDisplayText, tampil bila tidak kosong b2105) = galat simpan;
// - Save b5366 (`AddToList_Act` b5385) dan Cancel b5624 (`NewData_DT` b5647, hanya saat `DATASHOW = 'IsEdit'` b5787);
// - grid `BrowseCauseofLossLife_RD` b3981: ID b2989 (`.ID` b3415), Cause of Loss b3127 (`.CauseofLoss` b3562), kolom
//   tanpa judul b3245 berisi Edit b3747 (`EditList_DT` b3770 mengisi form dengan ID dan CauseofLoss); ID menaik tetap
//   (b3923), kolom tidak dapat diurutkan (b3925 / b3947 / b3969), tanpa saring (b4036); 10 baris per halaman b4057,
//   halaman bernomor b4028.
// TIDAK ada Delete, Upload CSV, saring, ringkasan, maupun detail (XML tidak memuatnya). Menu View only: tanpa form dan
// Edit (backend juga menolak tulisnya).

import { useCallback, useEffect, useState } from 'react'

import { Gagal, Kosong, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import { useBolehUbah } from '../../../../inti/frontend/lib/hakMenu'
import { ambilDaftar, tambah, ubah, type Halaman } from '../api'
import { BATAS_CAUSE_OF_LOSS, jumlahHalaman, periksaCauseOfLoss } from '../aturan'
import { COL } from '../labels'
import { NAMA_COL } from '../menu'

export default function CauseOfLossLife() {
  const bolehUbah = useBolehUbah(NAMA_COL)
  const [halaman, setHalaman] = useState(1)
  const [data, setData] = useState<Halaman | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [segar, setSegar] = useState(0)
  const [pesan, setPesan] = useState<string | null>(null)

  // Form: `ubahID` kosong = Add; terisi = Edit (`DATASHOW = 'IsEdit'`).
  const [ubahID, setUbahID] = useState('')
  const [nama, setNama] = useState('')
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
    setNama('')
    setGalatForm(null)
    setPesanForm(null)
  }

  const simpan = () => {
    if (sibuk) return
    const salah = periksaCauseOfLoss(nama)
    setPesanForm(salah)
    setGalatForm(null)
    if (salah !== null) return
    setSibuk(true)
    const isi = nama.trim()
    const janji = ubahID === '' ? tambah(isi) : ubah(ubahID, isi)
    janji.then(
      (c) => {
        setSibuk(false)
        kosongkanForm()
        setPesan(COL.tersimpan(c.id))
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
    <section className="inbox causeoflosslife__akar">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">{COL.judul}</h2>
      </header>

      {bolehUbah && (
        <form
          className="causeoflosslife__kartu causeoflosslife__form"
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
          <div className="causeoflosslife__baris">
            <label className="field causeoflosslife__nomor">
              <span className="field__label">{COL.id}</span>
              <input className="field__input field__input--readonly" value={ubahID} placeholder={COL.idOtomatis} readOnly disabled />
            </label>
            <label className="field causeoflosslife__isian">
              <span className="field__label">
                {COL.causeOfLoss} <span aria-hidden="true">*</span>
              </span>
              <input
                className="field__input"
                type="text"
                value={nama}
                maxLength={BATAS_CAUSE_OF_LOSS}
                required
                onChange={(e) => setNama(e.target.value)}
              />
            </label>
          </div>
          <div className="causeoflosslife__aksi">
            <button type="submit" className="btn btn--primary" disabled={sibuk}>
              {sibuk ? COL.menyimpan : COL.save}
            </button>
            {ubahID !== '' && (
              <button type="button" className="btn btn--ghost" onClick={kosongkanForm}>
                {COL.batal}
              </button>
            )}
            {ubahID !== '' && <span className="muted">{COL.modeUbah(ubahID)}</span>}
          </div>
        </form>
      )}

      <div className="causeoflosslife__alat">
        <span className="toolbar__spacer" />
        {data !== null && <span className="muted">{COL.jumlah(data.total)}</span>}
      </div>

      {pesan !== null && (
        <div className="alert alert--ok" role="status">
          {pesan}
        </div>
      )}
      {galat !== null && <Gagal galat={galat} />}
      {data === null && galat === null && <Memuat pesan={COL.memuat} />}
      {data !== null && data.daftar.length === 0 && <Kosong pesan={COL.kosong} />}
      {data !== null && data.daftar.length > 0 && (
        <>
          <div className="causeoflosslife__gulir">
            <table className="inbox__tabel causeoflosslife__tabel">
              <thead>
                <tr>
                  <th>{COL.id}</th>
                  <th>{COL.causeOfLoss}</th>
                  {bolehUbah && <th className="table__actions" aria-label={COL.edit} />}
                </tr>
              </thead>
              <tbody>
                {data.daftar.map((c) => (
                  <tr key={c.id} className="inbox__baris">
                    <td>{c.id}</td>
                    <td className="causeoflosslife__nama">{c.causeOfLoss}</td>
                    {bolehUbah && (
                      <td className="table__actions">
                        <span className="causeoflosslife__aksi">
                          <button
                            type="button"
                            className="btn btn--ghost btn--sm"
                            onClick={() => {
                              setPesan(null)
                              setGalatForm(null)
                              setPesanForm(null)
                              setUbahID(c.id)
                              setNama(c.causeOfLoss)
                            }}
                          >
                            {COL.edit}
                          </button>
                        </span>
                      </td>
                    )}
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          <div className="causeoflosslife__halaman">
            <button type="button" className="btn btn--ghost btn--sm" disabled={halaman <= 1} onClick={() => setHalaman((h) => h - 1)}>
              {COL.sebelum}
            </button>
            <span className="muted">{COL.halaman(data.halaman, dari)}</span>
            <button type="button" className="btn btn--ghost btn--sm" disabled={halaman >= dari} onClick={() => setHalaman((h) => h + 1)}>
              {COL.sesudah}
            </button>
          </div>
        </>
      )}
    </section>
  )
}
