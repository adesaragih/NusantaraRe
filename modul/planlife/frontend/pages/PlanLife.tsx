// Halaman Plan - section Pega `InboxProductType` ("Plan" b331), form dan grid di SATU layar:
// - pesan `OutputParam.ERRMSG` b5415 (read-only, tampil bila tidak kosong b5502) = galat simpan;
// - "Plan Name" b812 (`.CoverName` b841, pxTextInput b844);
// - "Business" b1079 (`.Business` b1105, pxAutoComplete b1107, `pyAllowFreeFormInput` b1125) dari
//   `BrowseBusinessLife_RD` b1183 - kolom OLDID b1200, Note b1231, ID b1265 (-> .BusinessID b1268);
// - "Benefit" b1432 (`.Benefit` b1461, pxAutoComplete b1464, `pyAllowFreeFormInput` b1478) dari `BrowseBenefitLife_RD`
//   b1537 - kolom Benefit b1552 (Number -> .BenefitID b1563-b1564, tersembunyi b1562);
// - Save b6403 (`SaveProductTypeLife_Act` b6422) dan New b6667 (`NewProductTypeLife_act` b6686, hanya saat
//   `DATASHOW = 'IsEdit'` b6827);
// - grid `BrowseProductTypeLife_RD` b4325: Plan Name b2978, Business b3117, Benefit b3255, kolom tanpa judul b3371
//   berisi Edit b4027 (`EditProductTypeLife_Act` b4045); 10 baris b4401, halaman bernomor b4373; tanpa saring b4380;
//   urut Plan Name / Benefit, Business tidak (b4247 / b4269 / b4291).
// Teks bebas boleh diketik (free-form); backend mencocokkan ulang ke master dan menolak (422) yang tidak cocok (K4).
// TIDAK ada Delete, Upload, ringkasan, detail, maupun medan ID (XML tidak memuatnya). View only: tanpa form dan Edit.

import { useCallback, useEffect, useState } from 'react'

import { Gagal, Kosong, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import { useBolehUbah } from '../../../../inti/frontend/lib/hakMenu'
import {
  ambilDaftar,
  ambilPilihanBenefit,
  ambilPilihanBusiness,
  tambah,
  ubah,
  type Halaman,
  type Isian,
  type PilihanBenefit,
  type PilihanBusiness,
  type Saringan,
} from '../api'
import {
  BATAS_NAMA,
  gantiUrut,
  ISIAN_KOSONG,
  isianDariPlan,
  jumlahHalaman,
  periksaIsian,
  SARINGAN_AWAL,
  saringBenefit,
  saringBusiness,
  tandaUrut,
} from '../aturan'
import { PL } from '../labels'
import { NAMA_PL } from '../menu'

type Terbuka = 'business' | 'benefit' | null

export default function PlanLife() {
  const bolehUbah = useBolehUbah(NAMA_PL)
  const [saring, setSaring] = useState<Saringan>(SARINGAN_AWAL)
  const [data, setData] = useState<Halaman | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [segar, setSegar] = useState(0)
  const [pesan, setPesan] = useState<string | null>(null)

  // Form: `ubahID` kosong = Add; terisi = Edit (`DATASHOW = 'IsEdit'`).
  const [ubahID, setUbahID] = useState('')
  const [isi, setIsi] = useState<Isian>(ISIAN_KOSONG)
  const [galatForm, setGalatForm] = useState<unknown>(null)
  const [pesanForm, setPesanForm] = useState<string | null>(null)
  const [sibuk, setSibuk] = useState(false)

  const [bizz, setBizz] = useState<PilihanBusiness[]>([])
  const [bens, setBens] = useState<PilihanBenefit[]>([])
  const [terbuka, setTerbuka] = useState<Terbuka>(null)

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

  useEffect(() => {
    if (!bolehUbah) return
    ambilPilihanBusiness().then(setBizz, (g: unknown) => setGalatForm(g))
    ambilPilihanBenefit().then(setBens, (g: unknown) => setGalatForm(g))
  }, [bolehUbah, segar])

  /** New (`NewProductTypeLife_act`): form kembali kosong (Add). */
  const kosongkanForm = () => {
    setUbahID('')
    setIsi(ISIAN_KOSONG)
    setGalatForm(null)
    setPesanForm(null)
    setTerbuka(null)
  }

  const simpan = () => {
    if (sibuk) return
    const salah = periksaIsian(isi)
    setPesanForm(salah)
    setGalatForm(null)
    if (salah !== null) return
    setSibuk(true)
    const kirim = { coverName: isi.coverName.trim(), business: isi.business.trim(), benefit: isi.benefit.trim() }
    const janji = ubahID === '' ? tambah(kirim) : ubah(ubahID, kirim)
    janji.then(
      (p) => {
        setSibuk(false)
        kosongkanForm()
        setPesan(PL.tersimpan(p.coverName))
        setSegar((n) => n + 1)
      },
      (g: unknown) => {
        setSibuk(false)
        setGalatForm(g)
      },
    )
  }

  const dari = data === null ? 1 : jumlahHalaman(data.total, data.ukuran)
  const pilihanBiz = saringBusiness(bizz, isi.business)
  const pilihanBen = saringBenefit(bens, isi.benefit)

  return (
    <section className="inbox planlife__akar">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">{PL.judul}</h2>
      </header>

      {bolehUbah && (
        <form
          className="planlife__kartu planlife__form"
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
          <div className="planlife__baris">
            <label className="field">
              <span className="field__label">{PL.planName}</span>
              <input
                className="field__input"
                value={isi.coverName}
                maxLength={BATAS_NAMA}
                onChange={(e) => setIsi((s) => ({ ...s, coverName: e.target.value }))}
              />
            </label>
            <div className="field planlife__isian">
              <label className="field__label" htmlFor="planlife-business">
                {PL.business}
              </label>
              <input
                id="planlife-business"
                className="field__input"
                value={isi.business}
                maxLength={BATAS_NAMA}
                autoComplete="off"
                placeholder={PL.pilihDariDaftar}
                onFocus={() => setTerbuka('business')}
                onBlur={() => setTimeout(() => setTerbuka((t) => (t === 'business' ? null : t)), 150)}
                onChange={(e) => {
                  setIsi((s) => ({ ...s, business: e.target.value }))
                  setTerbuka('business')
                }}
              />
              {terbuka === 'business' && (
                <div className="planlife__pilihan" role="listbox" aria-label={PL.business}>
                  {pilihanBiz.length === 0 && <div className="muted planlife__pilihan-kosong">{PL.tanpaPilihan}</div>}
                  {pilihanBiz.length > 0 && (
                    <div className="planlife__pilihan-kepala">
                      <span>{PL.oldId}</span>
                      <span>{PL.note}</span>
                      <span>{PL.id}</span>
                    </div>
                  )}
                  {pilihanBiz.map((b) => (
                    <button
                      key={b.id}
                      type="button"
                      role="option"
                      aria-selected={false}
                      className="planlife__pilihan-butir"
                      onMouseDown={(e) => e.preventDefault()}
                      onClick={() => {
                        setIsi((s) => ({ ...s, business: b.note }))
                        setTerbuka(null)
                      }}
                    >
                      <span>{b.oldId}</span>
                      <span>{b.note}</span>
                      <span>{b.id}</span>
                    </button>
                  ))}
                </div>
              )}
            </div>
            <div className="field planlife__isian">
              <label className="field__label" htmlFor="planlife-benefit">
                {PL.benefit}
              </label>
              <input
                id="planlife-benefit"
                className="field__input"
                value={isi.benefit}
                maxLength={BATAS_NAMA}
                autoComplete="off"
                placeholder={PL.pilihDariDaftar}
                onFocus={() => setTerbuka('benefit')}
                onBlur={() => setTimeout(() => setTerbuka((t) => (t === 'benefit' ? null : t)), 150)}
                onChange={(e) => {
                  setIsi((s) => ({ ...s, benefit: e.target.value }))
                  setTerbuka('benefit')
                }}
              />
              {terbuka === 'benefit' && (
                <div className="planlife__pilihan" role="listbox" aria-label={PL.benefit}>
                  {pilihanBen.length === 0 && <div className="muted planlife__pilihan-kosong">{PL.tanpaPilihan}</div>}
                  {pilihanBen.map((b) => (
                    <button
                      key={b.id}
                      type="button"
                      role="option"
                      aria-selected={false}
                      className="planlife__pilihan-butir planlife__pilihan-butir--satu"
                      onMouseDown={(e) => e.preventDefault()}
                      onClick={() => {
                        setIsi((s) => ({ ...s, benefit: b.benefit }))
                        setTerbuka(null)
                      }}
                    >
                      <span>{b.benefit}</span>
                    </button>
                  ))}
                </div>
              )}
            </div>
          </div>
          <div className="planlife__aksi">
            <button type="submit" className="btn btn--primary" disabled={sibuk}>
              {sibuk ? PL.menyimpan : PL.save}
            </button>
            {ubahID !== '' && (
              <button type="button" className="btn btn--ghost" onClick={kosongkanForm}>
                {PL.baru}
              </button>
            )}
            {ubahID !== '' && <span className="muted">{PL.modeUbah(isi.coverName)}</span>}
          </div>
        </form>
      )}

      <div className="planlife__alat">
        <span className="toolbar__spacer" />
        {data !== null && <span className="muted">{PL.jumlah(data.total)}</span>}
      </div>

      {pesan !== null && (
        <div className="alert alert--ok" role="status">
          {pesan}
        </div>
      )}
      {galat !== null && <Gagal galat={galat} />}
      {data === null && galat === null && <Memuat pesan={PL.memuat} />}
      {data !== null && data.daftar.length === 0 && <Kosong pesan={PL.kosong} />}
      {data !== null && data.daftar.length > 0 && (
        <>
          <div className="planlife__gulir">
            <table className="inbox__tabel planlife__tabel">
              <thead>
                <tr>
                  <th>
                    <button type="button" className="planlife__urut" onClick={() => setSaring((s) => gantiUrut(s, 'covername'))}>
                      {PL.planName}
                      {tandaUrut(saring, 'covername')}
                    </button>
                  </th>
                  <th>{PL.business}</th>
                  <th>
                    <button type="button" className="planlife__urut" onClick={() => setSaring((s) => gantiUrut(s, 'benefit'))}>
                      {PL.benefit}
                      {tandaUrut(saring, 'benefit')}
                    </button>
                  </th>
                  {bolehUbah && <th className="table__actions" aria-label={PL.edit} />}
                </tr>
              </thead>
              <tbody>
                {data.daftar.map((p) => (
                  <tr key={p.id} className="inbox__baris">
                    <td>{p.coverName}</td>
                    <td>{p.business}</td>
                    <td className="planlife__benefit">{p.benefit}</td>
                    {bolehUbah && (
                      <td className="table__actions">
                        <span className="planlife__aksi">
                          <button
                            type="button"
                            className="btn btn--ghost btn--sm"
                            onClick={() => {
                              setPesan(null)
                              setGalatForm(null)
                              setPesanForm(null)
                              setUbahID(p.id)
                              setIsi(isianDariPlan(p))
                            }}
                          >
                            {PL.edit}
                          </button>
                        </span>
                      </td>
                    )}
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          <div className="planlife__halaman">
            <button
              type="button"
              className="btn btn--ghost btn--sm"
              disabled={saring.halaman <= 1}
              onClick={() => setSaring((s) => ({ ...s, halaman: s.halaman - 1 }))}
            >
              {PL.sebelum}
            </button>
            <span className="muted">{PL.halaman(data.halaman, dari)}</span>
            <button
              type="button"
              className="btn btn--ghost btn--sm"
              disabled={saring.halaman >= dari}
              onClick={() => setSaring((s) => ({ ...s, halaman: s.halaman + 1 }))}
            >
              {PL.sesudah}
            </button>
          </div>
        </>
      )}
    </section>
  )
}
