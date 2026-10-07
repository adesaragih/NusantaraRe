// Tab utama Loss Record kasus FIRE (tiket 55) - port `NB FacIn\Section\InputInwardFacultativeDtl.xml` kontainer "Loss
// Record" -> `ObjectList_LossRecord` + `SummaryLossRecord_Section`, ubah baris = `InputLossRecord_Sec`, tata letak = gambar
// layar Pega dari work owner 06-10-2026:
// - bar "Please upload file with .csv format": Download Template CSV (isi = berkas work owner `templatecsvlossrecord`),
//   Upload CSV (baris diganti utuh seperti `SaveLossRecord`, dibaca di layar - `csvLossRecord.ts`), Save Data;
//   tombol Save tab di kaki (Dtl S89) - keduanya aksi yang sama (keputusan 2: satu simpan yang selalu menghitung ulang);
// - Premi PAR / MB / PL / Other COB (4 desimal); Premi PAR bawaan dari backend;
// - grid Uw Year · Class Of Business · Date of Loss · Cause Of Loss · Currency · Loss Amount · Loss Amount ( IDR ),
//   Tambah / Hapus, ▸ membuka ubah baris (Insured Name baca-saja, Date of Loss, Currency, Loss Amount, Cause of Loss,
//   Class Of Business saran `BUSINESS.NOTE`; Uw Year dan Loss Amount ( IDR ) dihitung backend);
// - kotak hijau Total Loss per Class Of Business + Currency; kotak biru Loss Ratio 1 / 3 / 5 Year (Gross / Net).
// Seluruh hitungan (Uw Year, kurs, total, loss ratio) di backend; layar hanya menampilkan.
//
// Keputusan work owner 06-10-2026: tab hanya untuk akun ber-POSITION "IT" (`aksesLossRecord`, ditegakkan backend);
// cacat Pega dibetulkan (ember 3 / 5 tahun, satu Save yang selalu menghitung ulang, Premi PAR bawaan, kurs).
// Keputusan agent L55-1: Uw Year dan Loss Amount ( IDR ) baca-saja (Pega: dapat diketik tetapi dihitung ulang saat Date of
// Loss / Loss Amount berubah). L55-2: Currency dan Cause of Loss dapat diubah di baris baru (Pega: baca-saja, sehingga baris
// tambahan manual tidak bisa diberi mata uang). L55-3: CSV yang galat tidak mengganti grid sama sekali.
// L55-7: Class Of Business wajib dari daftar `BUSINESS.NOTE` - ketikan (IsianSaran `wajibDariDaftar`) dan CSV
// (`cocokkanKelasBisnis`).

import { useEffect, useRef, useState, type ChangeEvent } from 'react'

import { Gagal, Memuat, Pilih } from '../../../../inti/frontend/components/ui/dasar'
import { angka4 } from '../angka4'
import {
  ambilLossRecord,
  daftarMataUang,
  hitungLossRecord,
  saranClassOfBusiness,
  simpanLossRecord,
  type BarisLossRecord,
  type TampilanLossRecord,
} from '../api'
import {
  LOSS_RECORD as L,
  TEKS_FORM_OPPORTUNITY,
  TEKS_LOSS_RECORD as T,
  TEMPLATE_CSV_LOSS_RECORD,
} from '../labels'
import { bacaCsvLossRecord, cocokkanKelasBisnis } from './csvLossRecord'
import IsianAngka4 from './IsianAngka4'
import IsianSaran from './IsianSaran'
import TanggalDMY from './TanggalDMY'

const JEDA_MS = 300

/** Baris baru dari tombol Tambah (`addRow`). */
export const barisLossRecordBaru = (): BarisLossRecord => ({
  uwYear: '',
  classOfBusiness: '',
  dateOfLoss: '',
  causeOfLoss: '',
  currency: '',
  lossAmount: '',
  lossAmountIdr: '',
})

/** Kabel DD-MM-YYYY -> tampil DD/MM/YYYY. */
const tanggal = (v: string) => v.replace(/-/g, '/')

/** Persen loss ratio: "6239.6" -> "6.239,6000%"; kosong -> "". */
export const persenRasio = (v: string) => (v.trim() === '' ? '' : `${angka4(v)}%`)

/** Unduh template CSV (padanan `Rule-File-Binary.downloadFile`). */
function unduhTemplate() {
  const url = URL.createObjectURL(new Blob([TEMPLATE_CSV_LOSS_RECORD], { type: 'text/csv' }))
  const a = document.createElement('a')
  a.href = url
  a.download = T.namaBerkas
  a.click()
  URL.revokeObjectURL(url)
}

/** `awal` hanya untuk uji render statis (tanpa DOM, efek tidak berjalan); layar memuat lewat `ambilLossRecord`. */
export default function TabLossRecord({
  caseId,
  insuredName,
  awal = null,
}: {
  caseId: string
  insuredName: string
  awal?: TampilanLossRecord | null
}) {
  const [v, setV] = useState<TampilanLossRecord | null>(awal)
  const [galat, setGalat] = useState<unknown>(null)
  const [galatCsv, setGalatCsv] = useState<string[]>([])
  const [infoCsv, setInfoCsv] = useState('')
  const [terbuka, setTerbuka] = useState<number[]>([])
  const [mataUang, setMataUang] = useState<string[]>([])
  const [menghitung, setMenghitung] = useState(false)
  const [menyimpan, setMenyimpan] = useState(false)
  const [tersimpan, setTersimpan] = useState(false)
  const nomor = useRef(0)
  const jadwal = useRef<number | undefined>(undefined)
  const berkas = useRef<HTMLInputElement>(null)

  useEffect(() => {
    let batal = false
    ambilLossRecord(caseId).then(
      (h) => {
        if (!batal) setV(h)
      },
      (err: unknown) => {
        if (!batal) setGalat(err)
      },
    )
    daftarMataUang().then(
      (h) => {
        if (!batal) setMataUang(h.baris)
      },
      (err: unknown) => {
        if (!batal) setGalat(err)
      },
    )
    return () => {
      batal = true
      window.clearTimeout(jadwal.current)
    }
  }, [caseId])

  if (v === null) return galat ? <Gagal galat={galat} /> : <Memuat />

  /** Ubah lalu jadwalkan hitung ulang di backend (urutan baris dipertahankan); jawaban lama dibuang. */
  function ubahHitung(baru: TampilanLossRecord) {
    setV(baru)
    setTersimpan(false)
    window.clearTimeout(jadwal.current)
    jadwal.current = window.setTimeout(() => {
      const n = ++nomor.current
      setMenghitung(true)
      setGalat(null)
      hitungLossRecord(caseId, baru).then(
        (h) => {
          if (n === nomor.current) {
            setV(h)
            setMenghitung(false)
          }
        },
        (err: unknown) => {
          if (n === nomor.current) {
            setGalat(err)
            setMenghitung(false)
          }
        },
      )
    }, JEDA_MS)
  }

  const ubahBaris = (i: number, b: Partial<BarisLossRecord>) => ubahHitung({ ...v, baris: v.baris.map((x, k) => (k === i ? { ...x, ...b } : x)) })

  async function simpan() {
    if (v === null) return
    window.clearTimeout(jadwal.current)
    nomor.current++
    setMenyimpan(true)
    setGalat(null)
    try {
      setV(await simpanLossRecord(caseId, v))
      setTerbuka([])
      setTersimpan(true)
      setInfoCsv('')
    } catch (err) {
      setGalat(err)
    } finally {
      setMenyimpan(false)
      setMenghitung(false)
    }
  }

  async function unggah(e: ChangeEvent<HTMLInputElement>) {
    const f = e.target.files?.[0]
    e.target.value = ''
    if (!f || v === null) return
    setInfoCsv(T.membaca)
    setGalatCsv([])
    let h = bacaCsvLossRecord(await f.text())
    if (h.galat.length === 0) {
      setInfoCsv(T.memeriksaKelas)
      try {
        h = await cocokkanKelasBisnis(h, saranClassOfBusiness)
      } catch (err) {
        setGalat(err)
        setInfoCsv('')
        return
      }
    }
    if (h.galat.length > 0) {
      setGalatCsv(h.galat)
      setInfoCsv('')
      return
    }
    setTerbuka([])
    setInfoCsv(`${h.baris.length}${T.csvTerbaca}`)
    ubahHitung({ ...v, baris: h.baris })
  }

  const premi = (label: string, kunci: 'premiPar' | 'premiMb' | 'premiPl' | 'premiOther') => (
    <div className="nbf-lr__premi-baris">
      <label htmlFor={`lr-${kunci}`}>{label}</label>
      <IsianAngka4 id={`lr-${kunci}`} value={v[kunci]} onChange={(x) => ubahHitung({ ...v, [kunci]: x })} />
    </div>
  )

  const rasio: [string, string][] = [
    [v.rasio.gross1, v.rasio.net1],
    [v.rasio.gross3, v.rasio.net3],
    [v.rasio.gross5, v.rasio.net5],
  ]

  return (
    <div className="nbf-objek nbf-lr">
      <Gagal galat={galat} />
      {tersimpan && <div className="alert alert--ok">{T.tersimpan}</div>}
      {v.pesan.map((p) => (
        <div key={p} className="alert alert--warn">
          {p}
        </div>
      ))}

      <section className="nbf-lr__unggah">
        <h5>{L.judulUnggah}</h5>
        <div className="nbf-lr__tombol">
          <button type="button" className="btn btn--ghost btn--sm" onClick={unduhTemplate}>
            {L.unduhTemplate}
          </button>
          <button type="button" className="btn btn--sm" onClick={() => berkas.current?.click()} disabled={menyimpan}>
            {L.unggah}
          </button>
          <input ref={berkas} type="file" accept=".csv,text/csv" hidden aria-label={L.unggah} onChange={(e) => void unggah(e)} />
          <button type="button" className="btn btn--primary btn--sm" onClick={() => void simpan()} disabled={menyimpan || menghitung}>
            {menyimpan ? T.menyimpan : L.simpan}
          </button>
          {menghitung && <span className="muted">{T.menghitung}</span>}
        </div>
        {infoCsv !== '' && <p className="nbf-lr__info">{infoCsv}</p>}
        {galatCsv.length > 0 && (
          <ul className="nbf-lr__galat">
            {galatCsv.map((g) => (
              <li key={g}>{g}</li>
            ))}
          </ul>
        )}
      </section>

      <section className="nbf-lr__premi">
        {premi(L.premiPar, 'premiPar')}
        {premi(L.premiMb, 'premiMb')}
        {premi(L.premiPl, 'premiPl')}
        {premi(L.premiOther, 'premiOther')}
      </section>

      <div className="nbf-cov-wrap">
        <table className="nbf-tabel nbf-lr__grid">
          <thead>
            <tr>
              <th scope="col" />
              <th scope="col">{L.uwYear}</th>
              <th scope="col">{L.classOfBusiness}</th>
              <th scope="col">{L.dateOfLoss}</th>
              <th scope="col">{L.causeOfLoss}</th>
              <th scope="col">{L.currency}</th>
              <th scope="col">{L.lossAmount}</th>
              <th scope="col">{L.lossAmountIdr}</th>
              <th scope="col">
                <button
                  type="button"
                  className="btn btn--ghost btn--sm"
                  onClick={() => {
                    setTerbuka((t) => [...t, v.baris.length])
                    ubahHitung({ ...v, baris: [...v.baris, barisLossRecordBaru()] })
                  }}
                >
                  ＋ {T.tambah}
                </button>
              </th>
            </tr>
          </thead>
          <tbody>
            {v.baris.length === 0 && (
              <tr>
                <td colSpan={9}>{T.kosong}</td>
              </tr>
            )}
            {v.baris.map((b, i) => {
              const buka = terbuka.includes(i)
              return [
                <tr key={`b${i}`} className={buka ? 'nbf-baris--terbuka' : undefined}>
                  <td>
                    <button
                      type="button"
                      className={'btn btn--ghost btn--sm nbf-buka' + (buka ? ' nbf-buka--terbuka' : '')}
                      aria-label={T.bukaBaris}
                      aria-expanded={buka}
                      onClick={() => setTerbuka((t) => (buka ? t.filter((x) => x !== i) : [...t, i]))}
                    >
                      ▸
                    </button>
                  </td>
                  <td>{b.uwYear}</td>
                  <td>{b.classOfBusiness}</td>
                  <td>{tanggal(b.dateOfLoss)}</td>
                  <td>{b.causeOfLoss}</td>
                  <td>{b.currency}</td>
                  <td className="nbf-angka">{angka4(b.lossAmount)}</td>
                  <td className="nbf-angka">{angka4(b.lossAmountIdr)}</td>
                  <td>
                    <button
                      type="button"
                      className="btn btn--ghost btn--sm"
                      onClick={() => {
                        setTerbuka([])
                        ubahHitung({ ...v, baris: v.baris.filter((_, k) => k !== i) })
                      }}
                    >
                      ✕ {T.hapus}
                    </button>
                  </td>
                </tr>,
                buka && (
                  <tr key={`d${i}`} className="nbf-objek__detail">
                    <td colSpan={9}>
                      <div className="nbf-lr__ubah">
                        <dl className="nbf-lr__ubah-info">
                          <div>
                            <dt>{L.insuredName}</dt>
                            <dd>{insuredName || '—'}</dd>
                          </div>
                          <div>
                            <dt>{L.uwYear}</dt>
                            <dd>{b.uwYear || '—'}</dd>
                          </div>
                          <div>
                            <dt>{L.lossAmountIdr}</dt>
                            <dd className="nbf-angka">{b.lossAmountIdr === '' ? '—' : angka4(b.lossAmountIdr)}</dd>
                          </div>
                        </dl>
                        <div className="nbf-lr__ubah-isian">
                          <TanggalDMY
                            label={L.dateOfLoss}
                            value={b.dateOfLoss}
                            onChange={(x) => ubahBaris(i, { dateOfLoss: x })}
                            labelKalender={TEKS_FORM_OPPORTUNITY.kalender}
                            pesanFormat={TEKS_FORM_OPPORTUNITY.formatTanggal}
                          />
                          <Pilih
                            label={L.currency}
                            value={b.currency}
                            onChange={(x) => ubahBaris(i, { currency: x })}
                            opsi={mataUang.map((m) => ({ value: m, label: m }))}
                          />
                          <div className="field">
                            <label className="field__label" htmlFor={`lr-jumlah-${i}`}>
                              {L.lossAmount}
                            </label>
                            <IsianAngka4 id={`lr-jumlah-${i}`} value={b.lossAmount} onChange={(x) => ubahBaris(i, { lossAmount: x })} />
                          </div>
                          <div className="field">
                            <label className="field__label" htmlFor={`lr-sebab-${i}`}>
                              {L.causeOfLossUbah}
                            </label>
                            <input
                              id={`lr-sebab-${i}`}
                              className="field__input"
                              value={b.causeOfLoss}
                              onChange={(e) => ubahBaris(i, { causeOfLoss: e.target.value.toUpperCase() })}
                            />
                          </div>
                          <div className="nbf-lr__ubah-lebar">
                            <IsianSaran
                              label={L.classOfBusiness}
                              value={b.classOfBusiness}
                              onKetik={(x) => ubahBaris(i, { classOfBusiness: x })}
                              cari={saranClassOfBusiness}
                              onPilih={(s) => ubahBaris(i, { classOfBusiness: s.label })}
                              muatSaatFokus
                              wajibDariDaftar
                            />
                          </div>
                        </div>
                      </div>
                    </td>
                  </tr>
                ),
              ]
            })}
          </tbody>
        </table>
      </div>

      <section className="nbf-lr__total">
        <div className="nbf-cov-wrap">
          <table className="nbf-tabel">
            <thead>
              <tr>
                <th scope="col">{L.classOfBusiness}</th>
                <th scope="col">{L.currency}</th>
                <th scope="col">{L.totalLoss}</th>
                <th scope="col">{L.totalLossIdr}</th>
              </tr>
            </thead>
            <tbody>
              {v.totalCob.length === 0 && (
                <tr>
                  <td colSpan={4}>{T.kosong}</td>
                </tr>
              )}
              {v.totalCob.map((t) => (
                <tr key={`${t.classOfBusiness}|${t.currency}`}>
                  <td>{t.classOfBusiness}</td>
                  <td>{t.currency}</td>
                  <td className="nbf-angka">{angka4(t.lossAmount)}</td>
                  <td className="nbf-angka">{angka4(t.lossAmountIdr)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </section>

      <section className="nbf-lr__rasio">
        {L.rasio.map((judul, k) => (
          <div key={judul} className="nbf-lr__rasio-baris">
            <span className="nbf-lr__rasio-judul">{judul}</span>
            <dl>
              <dt>{L.gross}</dt>
              <dd>{persenRasio(rasio[k]?.[0] ?? '')}</dd>
            </dl>
            <dl>
              <dt>{L.net}</dt>
              <dd>{persenRasio(rasio[k]?.[1] ?? '')}</dd>
            </dl>
          </div>
        ))}
      </section>

      <div className="nbf-objek__kaki">
        {menghitung && <span className="muted">{T.menghitung}</span>}
        <button type="button" className="btn btn--primary" onClick={() => void simpan()} disabled={menyimpan || menghitung}>
          {menyimpan ? T.menyimpan : L.simpanTab}
        </button>
      </div>
    </div>
  )
}
