// Halaman Benefit - section Pega `InboxBenefit` ("INSURANCE BENEFIT" b313), grid dan form di SATU layar:
// - pesan `InputParam.ERRMSG` b823 (read-only, tampil bila tidak kosong b908) = galat simpan;
// - "Number / ID" b964 (`.Number` b993, pxTextInput DISABLED b1012-b1013) - kosong = dibentuk server;
// - "Benefit" b1143 (`.Benefit` b1172, pxTextArea b1175, wajib b1137; perubahan -> `SetUpperCase_DT` b1211 = huruf besar);
// - Save b1772 (`AddToList_Act` b1791) dan Cancel b2035 (`NewData_DT` b2058, hanya saat `DATASHOW = 'IsEdit'` b2198);
// - grid `BrowseBenefitLife_RD` b4451: ID b3431 (`.Number` b3858, urut, bawaan MENURUN b4391), Benefit b3570
//   (`.Benefit` b4020, tidak dapat diurutkan b4415), kolom tanpa judul b3708 berisi Edit b4220 (`EditList_DT` b4243
//   mengisi form dengan Number dan Benefit); 10 baris per halaman b4526, halaman bernomor b4498.
// TIDAK ada Delete, Upload CSV, ringkasan, maupun detail (XML tidak memuatnya). Menu View only: tanpa form dan Edit
// (backend juga menolak tulisnya).

import { useCallback, useEffect, useState } from 'react'

import { Gagal, Kosong, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import { useBolehUbah } from '../../../../inti/frontend/lib/hakMenu'
import { ambilDaftar, tambah, ubah, type Halaman, type Saringan } from '../api'
import { balikArah, BATAS_BENEFIT, hurufBesar, jumlahHalaman, periksaBenefit, SARINGAN_AWAL, tandaArah } from '../aturan'
import { BN } from '../labels'
import { NAMA_BN } from '../menu'

const JEDA_CARI_MS = 300

export default function BenefitLife() {
  const bolehUbah = useBolehUbah(NAMA_BN)
  const [ketikID, setKetikID] = useState('')
  const [ketikBenefit, setKetikBenefit] = useState('')
  const [saring, setSaring] = useState<Saringan>(SARINGAN_AWAL)
  const [data, setData] = useState<Halaman | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [segar, setSegar] = useState(0)
  const [pesan, setPesan] = useState<string | null>(null)

  // Form: `ubahID` kosong = Add; terisi = Edit (`DATASHOW = 'IsEdit'`).
  const [ubahID, setUbahID] = useState('')
  const [benefit, setBenefit] = useState('')
  const [galatForm, setGalatForm] = useState<unknown>(null)
  const [pesanForm, setPesanForm] = useState<string | null>(null)
  const [sibuk, setSibuk] = useState(false)

  useEffect(() => {
    const t = setTimeout(() => setSaring((s) => ({ ...s, id: ketikID, benefit: ketikBenefit, halaman: 1 })), JEDA_CARI_MS)
    return () => clearTimeout(t)
  }, [ketikID, ketikBenefit])

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

  /** Cancel (`NewData_DT`): form kembali kosong (Add). */
  const kosongkanForm = () => {
    setUbahID('')
    setBenefit('')
    setGalatForm(null)
    setPesanForm(null)
  }

  const simpan = () => {
    if (sibuk) return
    const salah = periksaBenefit(benefit)
    setPesanForm(salah)
    setGalatForm(null)
    if (salah !== null) return
    setSibuk(true)
    const isi = hurufBesar(benefit.trim())
    const janji = ubahID === '' ? tambah(isi) : ubah(ubahID, isi)
    janji.then(
      (b) => {
        setSibuk(false)
        kosongkanForm()
        setPesan(BN.tersimpan(b.id))
        setSegar((n) => n + 1)
      },
      (g: unknown) => {
        setSibuk(false)
        setGalatForm(g)
      },
    )
  }

  const tersaring = saring.id.trim() !== '' || saring.benefit.trim() !== ''
  const dari = data === null ? 1 : jumlahHalaman(data.total, data.ukuran)

  return (
    <section className="inbox benefitlife__akar">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">{BN.judul}</h2>
      </header>

      {bolehUbah && (
        <form
          className="benefitlife__kartu benefitlife__form"
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
          <div className="benefitlife__baris">
            <label className="field benefitlife__nomor">
              <span className="field__label">{BN.nomorId}</span>
              <input className="field__input field__input--readonly" value={ubahID} placeholder={BN.idOtomatis} readOnly disabled />
            </label>
            <label className="field benefitlife__isian">
              <span className="field__label">
                {BN.benefit} <span aria-hidden="true">*</span>
              </span>
              <textarea
                className="field__input benefitlife__teks"
                value={benefit}
                rows={3}
                maxLength={BATAS_BENEFIT}
                required
                onChange={(e) => setBenefit(e.target.value)}
                onBlur={() => setBenefit((v) => hurufBesar(v))}
              />
            </label>
          </div>
          <div className="benefitlife__aksi">
            <button type="submit" className="btn btn--primary" disabled={sibuk}>
              {sibuk ? BN.menyimpan : BN.save}
            </button>
            {ubahID !== '' && (
              <button type="button" className="btn btn--ghost" onClick={kosongkanForm}>
                {BN.batal}
              </button>
            )}
            {ubahID !== '' && <span className="muted">{BN.modeUbah(ubahID)}</span>}
          </div>
        </form>
      )}

      <div className="benefitlife__alat">
        <input
          className="field__input benefitlife__saring"
          type="search"
          placeholder={BN.saringID}
          aria-label={BN.saringID}
          value={ketikID}
          onChange={(e) => setKetikID(e.target.value)}
        />
        <input
          className="field__input benefitlife__cari"
          type="search"
          placeholder={BN.saringBenefit}
          aria-label={BN.saringBenefit}
          value={ketikBenefit}
          onChange={(e) => setKetikBenefit(e.target.value)}
        />
        <span className="toolbar__spacer" />
        {data !== null && <span className="muted">{BN.jumlah(data.total)}</span>}
      </div>

      {pesan !== null && (
        <div className="alert alert--ok" role="status">
          {pesan}
        </div>
      )}
      {galat !== null && <Gagal galat={galat} />}
      {data === null && galat === null && <Memuat pesan={BN.memuat} />}
      {data !== null && data.daftar.length === 0 && <Kosong pesan={tersaring ? BN.tidakCocok : BN.kosong} />}
      {data !== null && data.daftar.length > 0 && (
        <>
          <div className="benefitlife__gulir">
            <table className="inbox__tabel benefitlife__tabel">
              <thead>
                <tr>
                  <th>
                    <button type="button" className="benefitlife__urut" onClick={() => setSaring((s) => balikArah(s))}>
                      {BN.id}
                      {tandaArah(saring)}
                    </button>
                  </th>
                  <th>{BN.benefit}</th>
                  {bolehUbah && <th className="table__actions" aria-label={BN.edit} />}
                </tr>
              </thead>
              <tbody>
                {data.daftar.map((b) => (
                  <tr key={b.id} className="inbox__baris">
                    <td>{b.id}</td>
                    <td className="benefitlife__benefit">{b.benefit}</td>
                    {bolehUbah && (
                      <td className="table__actions">
                        <span className="benefitlife__aksi">
                          <button
                            type="button"
                            className="btn btn--ghost btn--sm"
                            onClick={() => {
                              setPesan(null)
                              setGalatForm(null)
                              setPesanForm(null)
                              setUbahID(b.id)
                              setBenefit(b.benefit)
                            }}
                          >
                            {BN.edit}
                          </button>
                        </span>
                      </td>
                    )}
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          <div className="benefitlife__halaman">
            <button
              type="button"
              className="btn btn--ghost btn--sm"
              disabled={saring.halaman <= 1}
              onClick={() => setSaring((s) => ({ ...s, halaman: s.halaman - 1 }))}
            >
              {BN.sebelum}
            </button>
            <span className="muted">{BN.halaman(data.halaman, dari)}</span>
            <button
              type="button"
              className="btn btn--ghost btn--sm"
              disabled={saring.halaman >= dari}
              onClick={() => setSaring((s) => ({ ...s, halaman: s.halaman + 1 }))}
            >
              {BN.sesudah}
            </button>
          </div>
        </>
      )}
    </section>
  )
}
