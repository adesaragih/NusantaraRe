// Tab Premium Deduction / Brokerage Fee kasus FIRE (tiket 50) - port `NB FacIn\Section\InputInwardFacultativeDtl.xml`
// tab "Premium Deduction / Brokerage Fee" (baris 74726-90306) dengan gambar layar Pega DEV (work owner 05-10-2026):
// - % Premium Deduction (`.Policy.Payment.RICommision`), % Brokerage Fee (`.PctBrokerageFee`), centang "Exclude PPN ?"
//   (`.OfferFacIn.PPnCheck`) - diubah (change / Enter) -> `POST …/deduction/hitung` aksi "persen" (`SetRIComm_Act`);
// - daftar komisi `CoverageCommisionList` (FIRE: lokasi -> `PropertyItemListCoverageCommision` -> `InputCoverageCommisionFire`
//   -> `CommisionItem` + grid spreading `InputDtlSpreadingCommisionCoverage_FacIn`), baca-saja - data dari
//   `GET …/spreading` (pohon lokasi yang sama);
// - ringkasan Spreading (`SummarySpreading_Section`) - komponen tab Spreading;
// - grid per mata uang: Total Premium Deduction (`SetTotalPremium_Act`, aksi "commision"), Total Premium Deduction 2
//   (`ProtekDeduction2_Act`, aksi "deduction2"), Total Brokerage Fee (`TotalBrokerageFee_Act`, aksi "brokerage"); kepala
//   PPh / PPN = "PPH n %" / "PPN n %" (`SetTAXNote`, tarif `POOLDATA.TAX` terbaru);
// - Save -> `PUT …/deduction`.
// Seluruh hitungan di backend (desimal eksak); layar hanya menampilkan.
//
// Keputusan work owner 05-10-2026: centang "Exclude PPN ?" dibuat seperti Pega (`SetRIComm_Act` langkah 5 selalu
// menyetel `PPnCheck = true`, jadi PPh / PPN selalu dari Brokerage Fee penuh); daftar komisi dibuat seperti Pega.
// Keputusan agent: D-1 isian dihitung saat lepas fokus / Enter (event change / enter Pega), bukan tiap ketikan. D-2
// "Deduction 2 cannot be greater than the Premium" diperiksa juga di layar sebelum memanggil server.

import { useEffect, useState } from 'react'

import { Gagal, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import { desimalSah, jumlahDesimal, pangkasNolEkor } from '../../../../inti/frontend/lib/desimal'
import { angka4 } from '../angka4'
import {
  ambilDeduction,
  ambilSpreading,
  hitungDeduction,
  simpanDeduction,
  type AksiDeduction,
  type BarisDeduction,
  type LokasiSpreading,
  type TampilanDeduction,
  type TampilanSpreading,
} from '../api'
import { DEDUCTION as D, TEKS_DEDUCTION as T } from '../labels'
import IsianAngka4 from './IsianAngka4'
import { RingkasanSpreading } from './TabSpreading'

const angka = (v: string | undefined) => angka4(v ?? '')

/** Kepala kolom pajak = `"PPH " + Tax.pxResults(1).PPH + " %"` (`SetTAXNote`); tarif tak ada -> label saja. */
export function kepalaPajak(nama: string, tarif: string): string {
  const t = tarif.trim()
  return t === '' ? nama : `${nama} ${desimalSah(t) ? pangkasNolEkor(t, 0) : t} %`
}

/** a > b eksak (teks desimal); bukan desimal -> false. */
export function lebihBesar(a: string, b: string): boolean {
  const x = a.trim()
  const y = b.trim()
  if (!desimalSah(x) || !desimalSah(y)) return false
  const selisih = jumlahDesimal([x, y.startsWith('-') ? y.slice(1) : `-${y.replace(/^\+/, '')}`]).total
  return !selisih.startsWith('-') && /[1-9]/.test(selisih)
}

/** Deduction 2 > Premium baris itu (`ProtekDeduction2_Act`). */
export const deduction2Lebih = (b: BarisDeduction) => lebihBesar(b.deduction2, b.premium)

/** Tombol buka/tutup baris drill-down. */
function Buka({ buka, onKlik }: { buka: boolean; onKlik: () => void }) {
  return (
    <button
      type="button"
      className={'btn btn--ghost btn--sm nbf-buka' + (buka ? ' nbf-buka--terbuka' : '')}
      aria-label={T.bukaBaris}
      aria-expanded={buka}
      onClick={onKlik}
    >
      ▸
    </button>
  )
}

/** Daftar komisi (`CoverageCommisionList` FIRE) - lokasi -> item -> coverage -> detail, baca-saja. */
function DaftarKomisi({ lokasi }: { lokasi: LokasiSpreading[] }) {
  const [terbuka, setTerbuka] = useState<string[]>([])
  const ganti = (k: string) => setTerbuka((t) => (t.includes(k) ? t.filter((x) => x !== k) : [...t, k]))
  return (
    <div className="nbf-cov-wrap">
      <table className="nbf-tabel">
        <thead>
          <tr>
            <th scope="col" />
            {D.kolomLokasi.map((k) => (
              <th key={k} scope="col">
                {k}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {lokasi.length === 0 && (
            <tr>
              <td colSpan={4}>{T.kosong}</td>
            </tr>
          )}
          {lokasi.map((l, a) => [
            <tr key={`l${a}`} className={terbuka.includes(`l${a}`) ? 'nbf-baris--terbuka' : undefined}>
              <td>
                <Buka buka={terbuka.includes(`l${a}`)} onKlik={() => ganti(`l${a}`)} />
              </td>
              <td>{l.objectNo}</td>
              <td>{l.objectName}</td>
              <td>{l.location}</td>
            </tr>,
            terbuka.includes(`l${a}`) && (
              <tr key={`ld${a}`} className="nbf-objek__detail">
                <td colSpan={4}>
                  <div className="nbf-cov-wrap">
                    <table className="nbf-tabel">
                      <thead>
                        <tr>
                          <th scope="col" />
                          {D.kolomItem.map((k) => (
                            <th key={k} scope="col">
                              {k}
                            </th>
                          ))}
                        </tr>
                      </thead>
                      <tbody>
                        {l.items.map((it, b) => [
                          <tr key={`i${a}-${b}`}>
                            <td>
                              <Buka buka={terbuka.includes(`i${a}-${b}`)} onKlik={() => ganti(`i${a}-${b}`)} />
                            </td>
                            <td>{it.itemType}</td>
                            <td>{it.currency}</td>
                            <td className="nbf-angka">{angka(it.tsi)}</td>
                            <td className="nbf-angka">{angka(it.totalGrossPremi)}</td>
                            <td className="nbf-angka">{angka(it.totalPremiumRnm)}</td>
                          </tr>,
                          terbuka.includes(`i${a}-${b}`) && (
                            <tr key={`id${a}-${b}`} className="nbf-objek__detail">
                              <td colSpan={6}>
                                <div className="nbf-cov-wrap">
                                  <table className="nbf-tabel">
                                    <thead>
                                      <tr>
                                        <th scope="col" />
                                        {D.kolomCoverage.map((k) => (
                                          <th key={k} scope="col">
                                            {k}
                                          </th>
                                        ))}
                                      </tr>
                                    </thead>
                                    <tbody>
                                      {it.coverages.map((c, k) => [
                                        <tr key={`c${a}-${b}-${k}`}>
                                          <td>
                                            <Buka buka={terbuka.includes(`c${a}-${b}-${k}`)} onKlik={() => ganti(`c${a}-${b}-${k}`)} />
                                          </td>
                                          <td>{c.oldId}</td>
                                          <td className="nbf-angka">{angka(c.rate)}</td>
                                          <td className="nbf-angka">{angka(c.premium)}</td>
                                          <td className="nbf-angka">{angka(c.premiNusantaraRe)}</td>
                                        </tr>,
                                        terbuka.includes(`c${a}-${b}-${k}`) && (
                                          <tr key={`cd${a}-${b}-${k}`} className="nbf-objek__detail">
                                            <td colSpan={5}>
                                              <dl className="nbf-deduct__rinci">
                                                <div>
                                                  <dt>{D.detailCoverage.tsiLiability}</dt>
                                                  <dd>{angka(c.tsiLiability)}</dd>
                                                </div>
                                                <div>
                                                  <dt>{D.detailCoverage.premium}</dt>
                                                  <dd>{angka(c.premium)}</dd>
                                                </div>
                                                <div>
                                                  <dt>{D.detailCoverage.tsiNusantaraRe}</dt>
                                                  <dd>{angka(c.tsiNusantaraRe)}</dd>
                                                </div>
                                                <div>
                                                  <dt>{D.detailCoverage.premiNusantaraRe}</dt>
                                                  <dd>{angka(c.premiNusantaraRe)}</dd>
                                                </div>
                                              </dl>
                                              <div className="nbf-cov-wrap">
                                                <table className="nbf-tabel">
                                                  <thead>
                                                    <tr>
                                                      {D.kolomSpreading.map((h) => (
                                                        <th key={h} scope="col">
                                                          {h}
                                                        </th>
                                                      ))}
                                                    </tr>
                                                  </thead>
                                                  <tbody>
                                                    {c.spreading.length === 0 && (
                                                      <tr>
                                                        <td colSpan={5}>{T.kosong}</td>
                                                      </tr>
                                                    )}
                                                    {c.spreading.map((s, m) => (
                                                      <tr key={`${s.treatyType}-${m}`}>
                                                        <td>{s.treatyName}</td>
                                                        <td className="nbf-angka">{angka(s.sharePercentage)}</td>
                                                        <td className="nbf-angka">{angka(s.tsiSpreaded)}</td>
                                                        <td className="nbf-angka">{angka(s.claimEstimation)}</td>
                                                        <td className="nbf-angka">{angka(s.premiumSpreaded)}</td>
                                                      </tr>
                                                    ))}
                                                  </tbody>
                                                </table>
                                              </div>
                                            </td>
                                          </tr>
                                        ),
                                      ])}
                                    </tbody>
                                  </table>
                                </div>
                              </td>
                            </tr>
                          ),
                        ])}
                      </tbody>
                    </table>
                  </div>
                </td>
              </tr>
            ),
          ])}
        </tbody>
      </table>
    </div>
  )
}

export default function TabDeduction({ caseId }: { caseId: string }) {
  const [v, setV] = useState<TampilanDeduction | null>(null)
  const [sp, setSp] = useState<TampilanSpreading | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [menghitung, setMenghitung] = useState(false)
  const [menyimpan, setMenyimpan] = useState(false)
  const [tersimpan, setTersimpan] = useState(false)

  useEffect(() => {
    let batal = false
    ambilDeduction(caseId).then(
      (h) => {
        if (!batal) setV(h)
      },
      (err: unknown) => {
        if (!batal) setGalat(err)
      },
    )
    ambilSpreading(caseId).then(
      (h) => {
        if (!batal) setSp(h)
      },
      (err: unknown) => {
        if (!batal) setGalat(err)
      },
    )
    return () => {
      batal = true
    }
  }, [caseId])

  if (v === null) return galat ? <Gagal galat={galat} /> : <Memuat />

  const ubah = (baru: Partial<TampilanDeduction>) => {
    setV({ ...v, ...baru })
    setTersimpan(false)
  }
  const ubahBaris = (n: number, b: Partial<BarisDeduction>) => ubah({ mataUang: v.mataUang.map((x, k) => (k === n ? { ...x, ...b } : x)) })

  async function hitung(aksi: AksiDeduction, idx: number, isi: TampilanDeduction) {
    setMenghitung(true)
    setGalat(null)
    try {
      setV(await hitungDeduction(caseId, aksi, idx, isi))
    } catch (err) {
      setGalat(err)
    } finally {
      setMenghitung(false)
    }
  }

  async function simpan() {
    if (v === null || v.mataUang.some(deduction2Lebih)) return
    setMenyimpan(true)
    setGalat(null)
    try {
      setV(await simpanDeduction(caseId, v))
      setTersimpan(true)
    } catch (err) {
      setGalat(err)
    } finally {
      setMenyimpan(false)
    }
  }

  const isianPersen = (label: string, kunci: 'pctPremiumDeduction' | 'pctBrokerageFee') => (
    <div className="nbf-deduct__medan">
      <label htmlFor={`deduct-${kunci}`}>{label}</label>
      <IsianAngka4
        id={`deduct-${kunci}`}
        value={v[kunci]}
        onChange={(x) => ubah({ [kunci]: x })}
        onSelesai={() => void hitung('persen', 0, v)}
      />
    </div>
  )

  return (
    <div className="nbf-objek nbf-deduct">
      <Gagal galat={galat} />
      {tersimpan && <div className="alert alert--ok">{T.tersimpan}</div>}
      {v.pesan.map((p) => (
        <div key={p} className="alert alert--warn">
          {p}
        </div>
      ))}

      <section className="nbf-lapis nbf-deduct__atas">
        {isianPersen(D.pctPremiumDeduction, 'pctPremiumDeduction')}
        <div className="nbf-deduct__kanan">
          {isianPersen(D.pctBrokerageFee, 'pctBrokerageFee')}
          <label className="nbf-inward__pilihan">
            <input
              type="checkbox"
              checked={v.ppnCheck}
              onChange={(e) => {
                const isi = { ...v, ppnCheck: e.target.checked }
                ubah({ ppnCheck: e.target.checked })
                void hitung('persen', 0, isi)
              }}
            />{' '}
            {D.excludePpn}
          </label>
        </div>
        {menghitung && <span className="muted">{T.menghitung}</span>}
      </section>

      <section className="nbf-lapis">
        {sp === null ? <Memuat /> : <DaftarKomisi lokasi={sp.lokasi} />}
      </section>

      {sp !== null && <RingkasanSpreading v={sp} />}

      <div className="nbf-deduct__grid">
        <section className="nbf-deduct__kotak">
          <table className="nbf-tabel">
            <thead>
              <tr>
                <th scope="col">{D.kolomMataUang}</th>
                <th scope="col">{D.totalPremiumDeduction}</th>
                <th scope="col">{D.totalPremiumDeduction2}</th>
              </tr>
            </thead>
            <tbody>
              {v.mataUang.length === 0 && (
                <tr>
                  <td colSpan={3}>{T.kosong}</td>
                </tr>
              )}
              {v.mataUang.map((b, n) => (
                <tr key={b.currency}>
                  <td>{b.currency}</td>
                  <td>
                    <IsianAngka4
                      ariaLabel={`${D.totalPremiumDeduction} ${b.currency}`}
                      value={b.commision}
                      onChange={(x) => ubahBaris(n, { commision: x })}
                      onSelesai={() => void hitung('commision', n, v)}
                    />
                  </td>
                  <td>
                    <IsianAngka4
                      ariaLabel={`${D.totalPremiumDeduction2} ${b.currency}`}
                      value={b.deduction2}
                      onChange={(x) => ubahBaris(n, { deduction2: x })}
                      onSelesai={() => {
                        if (!deduction2Lebih(b)) void hitung('deduction2', n, v)
                      }}
                    />
                    {deduction2Lebih(b) && <div className="field__error">{T.deduction2Lebih}</div>}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </section>
        <section className="nbf-deduct__kotak">
          <table className="nbf-tabel">
            <thead>
              <tr>
                <th scope="col">{D.kolomMataUang}</th>
                <th scope="col">{D.totalBrokerageFee}</th>
                <th scope="col">{kepalaPajak(D.pph, v.tarifPph)}</th>
                <th scope="col">{kepalaPajak(D.ppn, v.tarifPpn)}</th>
              </tr>
            </thead>
            <tbody>
              {v.mataUang.length === 0 && (
                <tr>
                  <td colSpan={4}>{T.kosong}</td>
                </tr>
              )}
              {v.mataUang.map((b, n) => (
                <tr key={b.currency}>
                  <td>{b.currency}</td>
                  <td>
                    <IsianAngka4
                      ariaLabel={`${D.totalBrokerageFee} ${b.currency}`}
                      value={b.brokerageFee}
                      onChange={(x) => ubahBaris(n, { brokerageFee: x })}
                      onSelesai={() => void hitung('brokerage', n, v)}
                    />
                  </td>
                  <td className="nbf-angka">{angka(b.pph)}</td>
                  <td className="nbf-angka">{angka(b.ppn)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </section>
      </div>

      <div className="nbf-objek__kaki">
        <button type="button" className="btn btn--primary" onClick={() => void simpan()} disabled={menyimpan || menghitung}>
          {menyimpan ? T.menyimpan : T.simpan}
        </button>
      </div>
    </div>
  )
}
