// Halaman Disease Life - section Pega `InboxDisease` ("DISEASE" b371), grid dan form di SATU layar:
// - pesan `InputParam.ERRMSG` b881 (pxDisplayText, tampil bila tidak kosong b966) = galat simpan;
// - "Number / ID" b1021 (`DISEASE_LIFE.Number` b1050, pxTextInput DISABLED b1069-b1070) - kosong = dibentuk server;
// - "ICD Code" b1200 (`.ICD_Code` b1229, pxTextInput b1232; perubahan -> `SetUpperCase_DT` b1265 = huruf besar; wajib =
//   keputusan WO D3, XML `pyRequired` false b1244);
// - "Disease" b1470 (`.Disease` b1499, pxTextArea b1502, wajib b1465; perubahan -> `SetUpperCase_DT` b1538);
// - Save b2102 (`AddToList_Act` b2121) dan Cancel b2365 (`NewData_DT` b2388, hanya saat `DATASHOW = 'IsEdit'` b2525);
// - grid `BrowseDiseaseLife_RD` b5094: ID b3758 (`.Number` b4321, urut, bawaan MENURUN b5012), ICD Code b3894
//   (`.ICD_Code` b4484, urut b5036), Disease b4033 (`.Disease` b4629, tidak dapat diurutkan b5058), kolom tanpa judul
//   b4154 berisi Edit b4829 (`EditList_DT` b4852 mengisi form dengan Number, Disease, ICD_Code); saring b5148; 10 baris
//   per halaman b5169, halaman bernomor b5140.
// ⛔ 97.586 baris DEV: saring, urut, dan halaman SELALU di server (kebutuhan teknis, PARITAS) - layar ini tidak pernah
// memuat seluruh tabel.
// TIDAK ada Delete, Upload CSV, ringkasan, maupun detail (XML tidak memuatnya). Menu View only: tanpa form dan Edit
// (backend juga menolak tulisnya).

import { useCallback, useEffect, useState } from 'react'

import { Gagal, Kosong, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import { useBolehUbah } from '../../../../inti/frontend/lib/hakMenu'
import { ambilDaftar, tambah, ubah, type Halaman, type Saringan } from '../api'
import {
  BATAS_DISEASE,
  BATAS_ICD,
  hurufBesar,
  jumlahHalaman,
  periksaIsian,
  pilihUrut,
  SARINGAN_AWAL,
  tandaArah,
} from '../aturan'
import { DSL } from '../labels'
import { NAMA_DSL } from '../menu'

const JEDA_CARI_MS = 300

export default function DiseaseLife() {
  const bolehUbah = useBolehUbah(NAMA_DSL)
  const [ketikICD, setKetikICD] = useState('')
  const [ketikDisease, setKetikDisease] = useState('')
  const [saring, setSaring] = useState<Saringan>(SARINGAN_AWAL)
  const [data, setData] = useState<Halaman | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [segar, setSegar] = useState(0)
  const [pesan, setPesan] = useState<string | null>(null)

  // Form: `ubahID` kosong = Add; terisi = Edit (`DATASHOW = 'IsEdit'`).
  const [ubahID, setUbahID] = useState('')
  const [icdCode, setIcdCode] = useState('')
  const [disease, setDisease] = useState('')
  const [galatForm, setGalatForm] = useState<unknown>(null)
  const [pesanForm, setPesanForm] = useState<string | null>(null)
  const [sibuk, setSibuk] = useState(false)

  useEffect(() => {
    const t = setTimeout(() => setSaring((s) => ({ ...s, icdCode: ketikICD, disease: ketikDisease, halaman: 1 })), JEDA_CARI_MS)
    return () => clearTimeout(t)
  }, [ketikICD, ketikDisease])

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
    setIcdCode('')
    setDisease('')
    setGalatForm(null)
    setPesanForm(null)
  }

  const simpan = () => {
    if (sibuk) return
    const salah = periksaIsian(icdCode, disease)
    setPesanForm(salah)
    setGalatForm(null)
    if (salah !== null) return
    setSibuk(true)
    const icd = hurufBesar(icdCode.trim())
    const nama = hurufBesar(disease.trim())
    const janji = ubahID === '' ? tambah(icd, nama) : ubah(ubahID, icd, nama)
    janji.then(
      (p) => {
        setSibuk(false)
        kosongkanForm()
        setPesan(DSL.tersimpan(p.id))
        setSegar((n) => n + 1)
      },
      (g: unknown) => {
        setSibuk(false)
        setGalatForm(g)
      },
    )
  }

  const tersaring = saring.icdCode.trim() !== '' || saring.disease.trim() !== ''
  const dari = data === null ? 1 : jumlahHalaman(data.total, data.ukuran)

  return (
    <section className="inbox diseaselife__akar">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">{DSL.judul}</h2>
      </header>

      {bolehUbah && (
        <form
          className="diseaselife__kartu diseaselife__form"
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
          <div className="diseaselife__baris">
            <label className="field diseaselife__nomor">
              <span className="field__label">{DSL.nomorId}</span>
              <input className="field__input field__input--readonly" value={ubahID} placeholder={DSL.idOtomatis} readOnly disabled />
            </label>
            <label className="field diseaselife__isian">
              <span className="field__label">
                {DSL.icdCode} <span aria-hidden="true">*</span>
              </span>
              <input
                className="field__input"
                type="text"
                value={icdCode}
                maxLength={BATAS_ICD}
                required
                onChange={(e) => setIcdCode(e.target.value)}
                onBlur={() => setIcdCode((v) => hurufBesar(v))}
              />
            </label>
            <label className="field diseaselife__isian">
              <span className="field__label">
                {DSL.disease} <span aria-hidden="true">*</span>
              </span>
              <textarea
                className="field__input diseaselife__teks"
                value={disease}
                rows={3}
                maxLength={BATAS_DISEASE}
                required
                onChange={(e) => setDisease(e.target.value)}
                onBlur={() => setDisease((v) => hurufBesar(v))}
              />
            </label>
          </div>
          <div className="diseaselife__aksi">
            <button type="submit" className="btn btn--primary" disabled={sibuk}>
              {sibuk ? DSL.menyimpan : DSL.save}
            </button>
            {ubahID !== '' && (
              <button type="button" className="btn btn--ghost" onClick={kosongkanForm}>
                {DSL.batal}
              </button>
            )}
            {ubahID !== '' && <span className="muted">{DSL.modeUbah(ubahID)}</span>}
          </div>
        </form>
      )}

      <div className="diseaselife__alat">
        <input
          className="field__input diseaselife__saring"
          type="search"
          placeholder={DSL.saringICD}
          aria-label={DSL.saringICD}
          value={ketikICD}
          onChange={(e) => setKetikICD(e.target.value)}
        />
        <input
          className="field__input diseaselife__cari"
          type="search"
          placeholder={DSL.saringDisease}
          aria-label={DSL.saringDisease}
          value={ketikDisease}
          onChange={(e) => setKetikDisease(e.target.value)}
        />
        <span className="toolbar__spacer" />
        {data !== null && <span className="muted">{DSL.jumlah(data.total)}</span>}
      </div>

      {pesan !== null && (
        <div className="alert alert--ok" role="status">
          {pesan}
        </div>
      )}
      {galat !== null && <Gagal galat={galat} />}
      {data === null && galat === null && <Memuat pesan={DSL.memuat} />}
      {data !== null && data.daftar.length === 0 && <Kosong pesan={tersaring ? DSL.tidakCocok : DSL.kosong} />}
      {data !== null && data.daftar.length > 0 && (
        <>
          <div className="diseaselife__gulir">
            <table className="inbox__tabel diseaselife__tabel">
              <thead>
                <tr>
                  <th>
                    <button type="button" className="diseaselife__urut" onClick={() => setSaring((s) => pilihUrut(s, 'id'))}>
                      {DSL.id}
                      {tandaArah(saring, 'id')}
                    </button>
                  </th>
                  <th>
                    <button type="button" className="diseaselife__urut" onClick={() => setSaring((s) => pilihUrut(s, 'icd'))}>
                      {DSL.icdCode}
                      {tandaArah(saring, 'icd')}
                    </button>
                  </th>
                  <th>{DSL.disease}</th>
                  {bolehUbah && <th className="table__actions" aria-label={DSL.edit} />}
                </tr>
              </thead>
              <tbody>
                {data.daftar.map((p) => (
                  <tr key={p.id} className="inbox__baris">
                    <td>{p.id}</td>
                    <td>{p.icdCode}</td>
                    <td className="diseaselife__nama">{p.disease}</td>
                    {bolehUbah && (
                      <td className="table__actions">
                        <span className="diseaselife__aksi">
                          <button
                            type="button"
                            className="btn btn--ghost btn--sm"
                            onClick={() => {
                              setPesan(null)
                              setGalatForm(null)
                              setPesanForm(null)
                              setUbahID(p.id)
                              setIcdCode(p.icdCode)
                              setDisease(p.disease)
                            }}
                          >
                            {DSL.edit}
                          </button>
                        </span>
                      </td>
                    )}
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          <div className="diseaselife__halaman">
            <button
              type="button"
              className="btn btn--ghost btn--sm"
              disabled={saring.halaman <= 1}
              onClick={() => setSaring((s) => ({ ...s, halaman: s.halaman - 1 }))}
            >
              {DSL.sebelum}
            </button>
            <span className="muted">{DSL.halaman(data.halaman, dari)}</span>
            <button
              type="button"
              className="btn btn--ghost btn--sm"
              disabled={saring.halaman >= dari}
              onClick={() => setSaring((s) => ({ ...s, halaman: s.halaman + 1 }))}
            >
              {DSL.sesudah}
            </button>
          </div>
        </>
      )}
    </section>
  )
}
