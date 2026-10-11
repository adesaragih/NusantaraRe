// Popup Detail (`InboxSummaryRIRisk` b8869 -> `setIDUsedBy_Act` b8886 + harness `InboxRIRisk` b8923). Section
// `InboxRIRisk` (judul "R/I RISK DETAIL" b382):
//   - form: ID (`.ID` b1347, disabled), R/I RISK NAME (`TempIDUsedBy.USEDBY` b1727, disabled - selalu ringkasan yang
//     dilihat), CONTRACT wajib b1933 (pxTextInput b1923), YEAR b2202 dan MONTH b2389 TIDAK wajib (pyMax 4), RISK (PERMIL)
//     wajib b2595 (pxNumber b2582); Save b3132 (`AddToList_Act`), Cancel b3399 hanya saat Edit (`DATASHOW = 'IsEdit'`
//     b3565, `NewRIRiskLife_Act` b3423), Clear Field b1125 (`clearInputFieldRIRisk_act` b1148);
//   - grid `BrowseRIRiskLife_RD` b7612 (idusedby b7508): ID, R/I RISK NAME, CONTRACT, YEAR, MONTH, RISK (PERMIL), EDIT
//     b9784 (`EditRIRiskLife_Act` b9808); ID menaik b7607; 200 per halaman (pyPageSizeOther b10169).
// BEDA dari R/I Comm Life: medan MONTH, YEAR tidak wajib, label "RISK (PERMIL)", tombol "EDIT", 200 baris.
// Upload di section ini tersembunyi (`1=2` b4210) - tidak ditampilkan. Menu View only: tanpa form dan EDIT.

import { useEffect, useState } from 'react'

import { Gagal, Kosong, Memuat, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { ambilDetail, tambahDetail, ubahDetail, type Halaman, type IsianRincian, type Rincian, type Ringkasan } from '../api'
import {
  DIGIT_YEAR_MONTH,
  isianDariRincian,
  ISIAN_RINCIAN_KOSONG,
  jumlahHalaman,
  periksaIsianRincian,
  tampilDesimal,
  UKURAN_HALAMAN_RINCIAN,
} from '../aturan'
import { RK } from '../labels'

export default function RincianDetail({
  ringkasan,
  bolehUbah,
  onTutup,
}: {
  ringkasan: Ringkasan
  bolehUbah: boolean
  onTutup: () => void
}) {
  const [halaman, setHalaman] = useState(1)
  const [data, setData] = useState<Halaman<Rincian> | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [segar, setSegar] = useState(0)

  // Form: `ubahID` kosong = tambah baris.
  const [ubahID, setUbahID] = useState('')
  const [isi, setIsi] = useState<IsianRincian>(ISIAN_RINCIAN_KOSONG)
  const [pesanForm, setPesanForm] = useState<string | null>(null)
  const [galatForm, setGalatForm] = useState<unknown>(null)
  const [pesan, setPesan] = useState<string | null>(null)
  const [sibuk, setSibuk] = useState(false)

  useEffect(() => {
    let batal = false
    setData(null)
    ambilDetail(ringkasan.id, halaman).then(
      (d) => {
        if (batal) return
        setData(d)
        setGalat(null)
      },
      (g: unknown) => {
        if (!batal) setGalat(g)
      },
    )
    return () => {
      batal = true
    }
  }, [ringkasan.id, halaman, segar])

  const kosongkan = () => {
    setUbahID('')
    setIsi(ISIAN_RINCIAN_KOSONG)
    setPesanForm(null)
    setGalatForm(null)
  }

  const isiMedan = (k: keyof IsianRincian) => (e: { target: { value: string } }) => setIsi((s) => ({ ...s, [k]: e.target.value }))

  const simpan = () => {
    if (sibuk) return
    const salah = periksaIsianRincian(isi)
    setPesanForm(salah)
    setGalatForm(null)
    if (salah !== null) return
    setSibuk(true)
    const janji = ubahID === '' ? tambahDetail(ringkasan.id, isi) : ubahDetail(ringkasan.id, ubahID, isi)
    janji.then(
      (k) => {
        setSibuk(false)
        kosongkan()
        setPesan(RK.detailTersimpan(k.id))
        setSegar((n) => n + 1)
      },
      (g: unknown) => {
        setSibuk(false)
        setGalatForm(g)
      },
    )
  }

  const dari = data === null ? 1 : jumlahHalaman(data.total, data.ukuran || UKURAN_HALAMAN_RINCIAN)

  return (
    <Modal judul={`${RK.detail} — ${ringkasan.id} ${ringkasan.usedby}`} onTutup={onTutup} labelBatal={RK.tutup} lebar>
      {bolehUbah && (
        <form
          className="ririsklife__kartu ririsklife__form ririsklife__form--rapat"
          onSubmit={(e) => {
            e.preventDefault()
            simpan()
          }}
        >
          <h3 className="ririsklife__subjudul">{RK.judulRincian}</h3>
          {galatForm !== null && <Gagal galat={galatForm} />}
          {pesanForm !== null && (
            <div className="alert alert--warn" role="alert">
              {pesanForm}
            </div>
          )}
          {pesan !== null && pesanForm === null && galatForm === null && (
            <div className="alert alert--ok" role="status">
              {pesan}
            </div>
          )}
          <div className="ririsklife__rincian">
            <label className="field">
              <span className="field__label">{RK.id}</span>
              <input className="field__input field__input--readonly" value={ubahID} placeholder={RK.idDetailOtomatis} readOnly disabled />
            </label>
            <label className="field">
              <span className="field__label">
                {RK.nama} <span aria-hidden="true">*</span>
              </span>
              <input className="field__input field__input--readonly" value={ringkasan.usedby} readOnly disabled />
            </label>
            <label className="field">
              <span className="field__label">
                {RK.contract} <span aria-hidden="true">*</span>
              </span>
              <input
                className="field__input"
                inputMode="numeric"
                value={isi.contract}
                placeholder={RK.contohBulat}
                required
                onChange={isiMedan('contract')}
              />
            </label>
            <label className="field">
              <span className="field__label">{RK.year}</span>
              <input
                className="field__input"
                inputMode="numeric"
                value={isi.year}
                placeholder={RK.contohBulat}
                maxLength={DIGIT_YEAR_MONTH}
                onChange={isiMedan('year')}
              />
            </label>
            <label className="field">
              <span className="field__label">{RK.month}</span>
              <input
                className="field__input"
                inputMode="numeric"
                value={isi.month}
                placeholder={RK.contohBulat}
                maxLength={DIGIT_YEAR_MONTH}
                onChange={isiMedan('month')}
              />
            </label>
            <label className="field">
              <span className="field__label">
                {RK.risk} <span aria-hidden="true">*</span>
              </span>
              <input className="field__input" inputMode="decimal" value={isi.risk} required onChange={isiMedan('risk')} />
            </label>
          </div>
          <div className="ririsklife__aksi">
            <button type="submit" className="btn btn--primary" disabled={sibuk}>
              {sibuk ? RK.menyimpan : RK.save}
            </button>
            {ubahID !== '' && (
              <button type="button" className="btn btn--ghost" onClick={kosongkan}>
                {RK.batal}
              </button>
            )}
            <button type="button" className="btn btn--ghost" onClick={kosongkan}>
              {RK.bersihkan}
            </button>
            {ubahID !== '' && <span className="muted">{RK.modeUbahDetail(ubahID)}</span>}
          </div>
        </form>
      )}

      {galat !== null && <Gagal galat={galat} />}
      {data === null && galat === null && <Memuat pesan={RK.memuat} />}
      {data !== null && data.daftar.length === 0 && <Kosong pesan={RK.kosongDetail} />}
      {data !== null && data.daftar.length > 0 && (
        <>
          <div className="ririsklife__gulir">
            <table className="inbox__tabel ririsklife__tabel">
              <thead>
                <tr>
                  <th>{RK.id}</th>
                  <th>{RK.nama}</th>
                  <th>{RK.contract}</th>
                  <th>{RK.year}</th>
                  <th>{RK.month}</th>
                  <th className="ririsklife__angka">{RK.risk}</th>
                  {bolehUbah && <th>{RK.aksi}</th>}
                </tr>
              </thead>
              <tbody>
                {data.daftar.map((k) => (
                  <tr key={k.id} className={`inbox__baris${k.id === ubahID ? ' ririsklife__baris--pilih' : ''}`}>
                    <td>{k.id}</td>
                    <td>{k.usedby}</td>
                    <td>{k.contract}</td>
                    <td>{k.year}</td>
                    <td>{k.month}</td>
                    <td className="ririsklife__angka">{tampilDesimal(k.risk)}</td>
                    {bolehUbah && (
                      <td>
                        <button
                          type="button"
                          className="btn btn--ghost btn--sm"
                          onClick={() => {
                            setUbahID(k.id)
                            setIsi(isianDariRincian(k))
                            setPesanForm(null)
                            setGalatForm(null)
                            setPesan(null)
                          }}
                        >
                          {RK.editRincian}
                        </button>
                      </td>
                    )}
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          <div className="ririsklife__halaman">
            <button type="button" className="btn btn--ghost btn--sm" disabled={halaman <= 1} onClick={() => setHalaman((h) => h - 1)}>
              {RK.sebelum}
            </button>
            <span className="muted">{RK.halaman(data.halaman, dari)}</span>
            <button type="button" className="btn btn--ghost btn--sm" disabled={halaman >= dari} onClick={() => setHalaman((h) => h + 1)}>
              {RK.sesudah}
            </button>
          </div>
        </>
      )}
    </Modal>
  )
}
