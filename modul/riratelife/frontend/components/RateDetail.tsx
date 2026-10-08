// Popup "Rate Detail" (b11446; tombol Detail b11398 -> `setIDUsedBy_Act` + harness `InboxRIRate`). Section
// `InboxRIRate` (View Detail.xml, judul "R/I RATE DETAIL" b367):
//   - form: ID (`.ID` disabled), R/I RATE NAME (`TempIDUsedBy.USEDBY` b1747, disabled - selalu ringkasan yang dilihat),
//     GENDER radio b1938, CONTRACT wajib b2141, AGE b2400, RATE wajib b2579; Save b3118 (`AddToList_Act`), Cancel b3418
//     hanya saat Edit (`DATASHOW = 'IsEdit'` b3638), Clear Field b1064;
//   - grid `BrowseRateLife_RD` (idusedby b7568): ID, R/I RATE NAME, GENDER, CONTRACT, AGE, RATE, Edit b9773; 20 per
//     halaman b10206, ID menurun.
// Upload di section ini tersembunyi (`1=2` b4029) - tidak ditampilkan. Menu View only: tanpa form dan Edit.

import { useEffect, useState } from 'react'

import { Gagal, Kosong, Memuat, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { ambilRate, tambahRate, ubahRate, type Halaman, type IsianRate, type Rate, type Ringkasan } from '../api'
import { GENDER, isianDariRate, ISIAN_RATE_KOSONG, jumlahHalaman, periksaIsianRate } from '../aturan'
import { RR } from '../labels'

export default function RateDetail({
  ringkasan,
  bolehUbah,
  onTutup,
}: {
  ringkasan: Ringkasan
  bolehUbah: boolean
  onTutup: () => void
}) {
  const [halaman, setHalaman] = useState(1)
  const [data, setData] = useState<Halaman<Rate> | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [segar, setSegar] = useState(0)

  // Form: `ubahID` kosong = tambah baris.
  const [ubahID, setUbahID] = useState('')
  const [isi, setIsi] = useState<IsianRate>(ISIAN_RATE_KOSONG)
  const [pesanForm, setPesanForm] = useState<string | null>(null)
  const [galatForm, setGalatForm] = useState<unknown>(null)
  const [pesan, setPesan] = useState<string | null>(null)
  const [sibuk, setSibuk] = useState(false)

  useEffect(() => {
    let batal = false
    setData(null)
    ambilRate(ringkasan.id, halaman).then(
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
    setIsi(ISIAN_RATE_KOSONG)
    setPesanForm(null)
    setGalatForm(null)
  }

  const isiMedan = (k: keyof IsianRate) => (e: { target: { value: string } }) => setIsi((s) => ({ ...s, [k]: e.target.value }))

  const simpan = () => {
    if (sibuk) return
    const salah = periksaIsianRate(isi)
    setPesanForm(salah)
    setGalatForm(null)
    if (salah !== null) return
    setSibuk(true)
    const janji = ubahID === '' ? tambahRate(ringkasan.id, isi) : ubahRate(ringkasan.id, ubahID, isi)
    janji.then(
      (r) => {
        setSibuk(false)
        kosongkan()
        setPesan(RR.rateTersimpan(r.id))
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
    <Modal judul={`${RR.judulDetail} — ${ringkasan.id} ${ringkasan.usedby}`} onTutup={onTutup} labelBatal={RR.tutup} lebar>
      {bolehUbah && (
        <form
          className="riratelife__kartu riratelife__form riratelife__form--rapat"
          onSubmit={(e) => {
            e.preventDefault()
            simpan()
          }}
        >
          <h3 className="riratelife__subjudul">{RR.judulRincian}</h3>
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
          <div className="riratelife__rincian">
            <label className="field">
              <span className="field__label">{RR.id}</span>
              <input className="field__input field__input--readonly" value={ubahID} placeholder={RR.idRateOtomatis} readOnly disabled />
            </label>
            <label className="field">
              <span className="field__label">
                {RR.nama} <span aria-hidden="true">*</span>
              </span>
              <input className="field__input field__input--readonly" value={ringkasan.usedby} readOnly disabled />
            </label>
            <fieldset className="field riratelife__pilihan">
              <legend className="field__label">{RR.gender}</legend>
              <div className="riratelife__radio-baris">
                {GENDER.map((g) => (
                  <label key={g} className="riratelife__radio">
                    <input type="radio" name="riratelife-gender" value={g} checked={isi.gender === g} onChange={isiMedan('gender')} />
                    {g}
                  </label>
                ))}
              </div>
            </fieldset>
            <label className="field">
              <span className="field__label">
                {RR.contract} <span aria-hidden="true">*</span>
              </span>
              <input
                className="field__input"
                inputMode="numeric"
                value={isi.contract}
                placeholder={RR.contohBulat}
                required
                onChange={isiMedan('contract')}
              />
            </label>
            <label className="field">
              <span className="field__label">{RR.age}</span>
              <input className="field__input" inputMode="numeric" value={isi.age} placeholder={RR.contohBulat} onChange={isiMedan('age')} />
            </label>
            <label className="field">
              <span className="field__label">
                {RR.rate} <span aria-hidden="true">*</span>
              </span>
              <input
                className="field__input"
                inputMode="decimal"
                value={isi.rate}
                placeholder={RR.contohRate}
                required
                onChange={isiMedan('rate')}
              />
            </label>
          </div>
          <div className="riratelife__aksi">
            <button type="submit" className="btn btn--primary" disabled={sibuk}>
              {sibuk ? RR.menyimpan : RR.save}
            </button>
            {ubahID !== '' && (
              <button type="button" className="btn btn--ghost" onClick={kosongkan}>
                {RR.batal}
              </button>
            )}
            <button type="button" className="btn btn--ghost" onClick={kosongkan}>
              {RR.bersihkan}
            </button>
            {ubahID !== '' && <span className="muted">{RR.modeUbahRate(ubahID)}</span>}
          </div>
        </form>
      )}

      {galat !== null && <Gagal galat={galat} />}
      {data === null && galat === null && <Memuat pesan={RR.memuat} />}
      {data !== null && data.daftar.length === 0 && <Kosong pesan={RR.kosongDetail} />}
      {data !== null && data.daftar.length > 0 && (
        <>
          <div className="riratelife__gulir">
            <table className="inbox__tabel riratelife__tabel">
              <thead>
                <tr>
                  <th>{RR.id}</th>
                  <th>{RR.nama}</th>
                  <th>{RR.gender}</th>
                  <th>{RR.contract}</th>
                  <th>{RR.age}</th>
                  <th className="riratelife__angka">{RR.rate}</th>
                  {bolehUbah && <th>{RR.aksi}</th>}
                </tr>
              </thead>
              <tbody>
                {data.daftar.map((r) => (
                  <tr key={r.id} className={`inbox__baris${r.id === ubahID ? ' riratelife__baris--pilih' : ''}`}>
                    <td>{r.id}</td>
                    <td>{r.usedby}</td>
                    <td>{r.gender}</td>
                    <td>{r.contract}</td>
                    <td>{r.age}</td>
                    <td className="riratelife__angka">{r.rate}</td>
                    {bolehUbah && (
                      <td>
                        <button
                          type="button"
                          className="btn btn--ghost btn--sm"
                          onClick={() => {
                            setUbahID(r.id)
                            setIsi(isianDariRate(r))
                            setPesanForm(null)
                            setGalatForm(null)
                            setPesan(null)
                          }}
                        >
                          {RR.edit}
                        </button>
                      </td>
                    )}
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          <div className="riratelife__halaman">
            <button type="button" className="btn btn--ghost btn--sm" disabled={halaman <= 1} onClick={() => setHalaman((h) => h - 1)}>
              {RR.sebelum}
            </button>
            <span className="muted">{RR.halaman(data.halaman, dari)}</span>
            <button type="button" className="btn btn--ghost btn--sm" disabled={halaman >= dari} onClick={() => setHalaman((h) => h + 1)}>
              {RR.sesudah}
            </button>
          </div>
        </>
      )}
    </Modal>
  )
}
