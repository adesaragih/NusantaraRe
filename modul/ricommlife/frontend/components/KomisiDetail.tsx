// Popup Detail (b9042 -> `setIDUsedBy_Act` b9059 + harness `InboxRIComm` b9096). Section `InboxRIComm` (judul
// "R/I COMM DETAIL" b349):
//   - form: ID (`.ID` b1325, disabled), R/I COMM NAME (`TempIDUsedBy.USEDBY` b1718, disabled - selalu ringkasan yang
//     dilihat), CONTRACT wajib b1923, YEAR wajib b2201 (pyMax 4 b2206), COMM wajib b2412; Save b2931 (`AddToList_Act`),
//     Cancel b3230 hanya saat Edit (`DATASHOW = 'IsEdit'` b3396), Clear Field b1100;
//   - grid `BrowseRICommLife_RD` (idusedby b7396): ID, R/I COMM NAME, CONTRACT, YEAR, COMM, Edit b9295; ID menaik
//     b9515, 50 per halaman b9716.
// Upload di section ini tersembunyi (`1=2` b4041) - tidak ditampilkan. Menu View only: tanpa form dan Edit.

import { useEffect, useState } from 'react'

import { Gagal, Kosong, Memuat, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { ambilDetail, tambahDetail, ubahDetail, type Halaman, type IsianKomisi, type Komisi, type Ringkasan } from '../api'
import { DIGIT_YEAR, isianDariKomisi, ISIAN_KOMISI_KOSONG, jumlahHalaman, periksaIsianKomisi, tampilDesimal } from '../aturan'
import { RC } from '../labels'

export default function KomisiDetail({
  ringkasan,
  bolehUbah,
  onTutup,
}: {
  ringkasan: Ringkasan
  bolehUbah: boolean
  onTutup: () => void
}) {
  const [halaman, setHalaman] = useState(1)
  const [data, setData] = useState<Halaman<Komisi> | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [segar, setSegar] = useState(0)

  // Form: `ubahID` kosong = tambah baris.
  const [ubahID, setUbahID] = useState('')
  const [isi, setIsi] = useState<IsianKomisi>(ISIAN_KOMISI_KOSONG)
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
    setIsi(ISIAN_KOMISI_KOSONG)
    setPesanForm(null)
    setGalatForm(null)
  }

  const isiMedan = (k: keyof IsianKomisi) => (e: { target: { value: string } }) => setIsi((s) => ({ ...s, [k]: e.target.value }))

  const simpan = () => {
    if (sibuk) return
    const salah = periksaIsianKomisi(isi)
    setPesanForm(salah)
    setGalatForm(null)
    if (salah !== null) return
    setSibuk(true)
    const janji = ubahID === '' ? tambahDetail(ringkasan.id, isi) : ubahDetail(ringkasan.id, ubahID, isi)
    janji.then(
      (k) => {
        setSibuk(false)
        kosongkan()
        setPesan(RC.detailTersimpan(k.id))
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
    <Modal judul={`${RC.detail} — ${ringkasan.id} ${ringkasan.usedby}`} onTutup={onTutup} labelBatal={RC.tutup} lebar>
      {bolehUbah && (
        <form
          className="ricommlife__kartu ricommlife__form ricommlife__form--rapat"
          onSubmit={(e) => {
            e.preventDefault()
            simpan()
          }}
        >
          <h3 className="ricommlife__subjudul">{RC.judulRincian}</h3>
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
          <div className="ricommlife__rincian">
            <label className="field">
              <span className="field__label">{RC.id}</span>
              <input className="field__input field__input--readonly" value={ubahID} placeholder={RC.idDetailOtomatis} readOnly disabled />
            </label>
            <label className="field">
              <span className="field__label">
                {RC.nama} <span aria-hidden="true">*</span>
              </span>
              <input className="field__input field__input--readonly" value={ringkasan.usedby} readOnly disabled />
            </label>
            <label className="field">
              <span className="field__label">
                {RC.contract} <span aria-hidden="true">*</span>
              </span>
              <input
                className="field__input"
                inputMode="numeric"
                value={isi.contract}
                placeholder={RC.contohBulat}
                required
                onChange={isiMedan('contract')}
              />
            </label>
            <label className="field">
              <span className="field__label">
                {RC.year} <span aria-hidden="true">*</span>
              </span>
              <input
                className="field__input"
                inputMode="numeric"
                value={isi.year}
                placeholder={RC.contohBulat}
                maxLength={DIGIT_YEAR}
                required
                onChange={isiMedan('year')}
              />
            </label>
            <label className="field">
              <span className="field__label">
                {RC.comm} <span aria-hidden="true">*</span>
              </span>
              <input className="field__input" inputMode="decimal" value={isi.comm} required onChange={isiMedan('comm')} />
            </label>
          </div>
          <div className="ricommlife__aksi">
            <button type="submit" className="btn btn--primary" disabled={sibuk}>
              {sibuk ? RC.menyimpan : RC.save}
            </button>
            {ubahID !== '' && (
              <button type="button" className="btn btn--ghost" onClick={kosongkan}>
                {RC.batal}
              </button>
            )}
            <button type="button" className="btn btn--ghost" onClick={kosongkan}>
              {RC.bersihkan}
            </button>
            {ubahID !== '' && <span className="muted">{RC.modeUbahDetail(ubahID)}</span>}
          </div>
        </form>
      )}

      {galat !== null && <Gagal galat={galat} />}
      {data === null && galat === null && <Memuat pesan={RC.memuat} />}
      {data !== null && data.daftar.length === 0 && <Kosong pesan={RC.kosongDetail} />}
      {data !== null && data.daftar.length > 0 && (
        <>
          <div className="ricommlife__gulir">
            <table className="inbox__tabel ricommlife__tabel">
              <thead>
                <tr>
                  <th>{RC.id}</th>
                  <th>{RC.nama}</th>
                  <th>{RC.contract}</th>
                  <th>{RC.year}</th>
                  <th className="ricommlife__angka">{RC.comm}</th>
                  {bolehUbah && <th>{RC.aksi}</th>}
                </tr>
              </thead>
              <tbody>
                {data.daftar.map((k) => (
                  <tr key={k.id} className={`inbox__baris${k.id === ubahID ? ' ricommlife__baris--pilih' : ''}`}>
                    <td>{k.id}</td>
                    <td>{k.usedby}</td>
                    <td>{k.contract}</td>
                    <td>{k.year}</td>
                    <td className="ricommlife__angka">{tampilDesimal(k.comm)}</td>
                    {bolehUbah && (
                      <td>
                        <button
                          type="button"
                          className="btn btn--ghost btn--sm"
                          onClick={() => {
                            setUbahID(k.id)
                            setIsi(isianDariKomisi(k))
                            setPesanForm(null)
                            setGalatForm(null)
                            setPesan(null)
                          }}
                        >
                          {RC.edit}
                        </button>
                      </td>
                    )}
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          <div className="ricommlife__halaman">
            <button type="button" className="btn btn--ghost btn--sm" disabled={halaman <= 1} onClick={() => setHalaman((h) => h - 1)}>
              {RC.sebelum}
            </button>
            <span className="muted">{RC.halaman(data.halaman, dari)}</span>
            <button type="button" className="btn btn--ghost btn--sm" disabled={halaman >= dari} onClick={() => setHalaman((h) => h + 1)}>
              {RC.sesudah}
            </button>
          </div>
        </>
      )}
    </Modal>
  )
}
