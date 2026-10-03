// Tab Coverage FIRE layar Inward Facultative (tiket 43, tahap C1).
//
// Port berantai: `CoverageList.xml` (grid objek `.LocationList`: No. · Object Name · Location; baris dibuka) ->
// `PropertyItemListCoverage.xml` (grid item: Object Item Type · Currency · TSI Object Item · Total Gross Premium ·
// ‰ Total Net Rate; grid total per mata uang) -> `InputCoverageFire.xml` (grid `.CoverageList`: Coverage · ‰ Standard
// Rate · Premi; Tambah = `AddCoverageAutoFire` - lima coverage FIRE bila masih kosong; Hapus) -> `FormCoverage`.
// Save = Dtl sel `Save` (`SaveFacIn_Act`): seluruh objek dikirim `PUT …/objek` (sama dengan tab Object).
//
// Total Gross Premium item = Σ Premium coverage (`CountPremi_ACT` langkah total) - dijumlah eksak (`jumlahDesimal`).
// Total per mata uang dan ‰ Total Net Rate = dari server (pembagian di backend, tanpa float).
//
// Keputusan agent (tiket 43): P-5 perubahan di tab Coverage dan tab Object disimpan masing-masing lewat tombol Save
// tab itu; pindah tab tanpa Save membuang perubahan (sama dengan pola tab sebelumnya). P-6 tombol Copy Coverage /
// Copy Deductible / Copy Coverage From dan grid Summary = tahap berikut.

import { useEffect, useRef, useState } from 'react'

import { Gagal, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import { jumlahDesimal } from '../../../../inti/frontend/lib/desimal'
import { formatNumber } from '../../../../inti/frontend/lib/format'
import { ambilObjek, coverageOtomatis, simpanObjek, type CoverageObjek, type ObjekFire } from '../api'
import {
  GRID_COV_ITEM,
  GRID_COV_OBJEK,
  GRID_COV_TOTAL,
  GRID_COVERAGE,
  GRID_OBJEK,
  SIMPAN_COVERAGE,
  TEKS_FORM_OPPORTUNITY,
  TEKS_INWARD,
  TEKS_OBJEK,
} from '../labels'
import FormCoverage, { adaGalatCoverage, coverageBaru } from './FormCoverage'
import { rapikanObjek } from './TabObject'

const DESIMAL = 4

/** Total Gross Premium item = Σ Premium coverage (eksak); tanpa coverage -> nilai server. */
export function totalPremiItem(coverages: CoverageObjek[] | undefined, dariServer: string | undefined): string {
  if (!coverages || coverages.length === 0) return dariServer ?? ''
  return jumlahDesimal(coverages.map((c) => c.premium)).total
}

/** Tombol buka/tutup baris. */
function Buka({ buka, onKlik }: { buka: boolean; onKlik: () => void }) {
  return (
    <button type="button" className="btn btn--ghost btn--sm" aria-label={TEKS_OBJEK.bukaBaris} aria-expanded={buka} onClick={onKlik}>
      {buka ? '▾' : '▸'}
    </button>
  )
}

const balik = (daftar: string[], k: string) => (daftar.includes(k) ? daftar.filter((x) => x !== k) : [...daftar, k])

export default function TabCoverage({ caseId }: { caseId: string }) {
  const [objek, setObjek] = useState<ObjekFire[] | null>(null)
  const [terbuka, setTerbuka] = useState<string[]>([])
  const [galat, setGalat] = useState<unknown>(null)
  const [menyimpan, setMenyimpan] = useState(false)
  const [tersimpan, setTersimpan] = useState(false)
  const aktif = useRef(true)

  useEffect(() => {
    aktif.current = true
    ambilObjek(caseId).then(
      (h) => {
        if (aktif.current) setObjek(h.baris.map(rapikanObjek))
      },
      (err: unknown) => {
        if (aktif.current) setGalat(err)
      },
    )
    return () => {
      aktif.current = false
    }
  }, [caseId])

  function ubahCoverage(o: number, i: number, coverages: CoverageObjek[]) {
    setObjek((d) =>
      d === null
        ? d
        : d.map((x, a) => (a !== o ? x : { ...x, items: x.items.map((it, b) => (b !== i ? it : { ...it, coverages })) })),
    )
    setTersimpan(false)
  }

  async function tambah(o: number, i: number, ada: CoverageObjek[]) {
    if (ada.length > 0) {
      ubahCoverage(o, i, [...ada, coverageBaru()])
      setTerbuka((t) => [...t, `c${o}-${i}-${ada.length}`])
      return
    }
    try {
      const h = await coverageOtomatis()
      ubahCoverage(o, i, h.baris.map((b) => coverageBaru(b)))
    } catch (err) {
      setGalat(err)
    }
  }

  async function simpan() {
    if (objek === null || objek.some((x) => x.items.some((it) => adaGalatCoverage(it.coverages ?? [])))) return
    setMenyimpan(true)
    setGalat(null)
    setTersimpan(false)
    try {
      const h = await simpanObjek(caseId, objek)
      setObjek(h.baris.map(rapikanObjek))
      setTersimpan(true)
    } catch (err) {
      setGalat(err)
    } finally {
      setMenyimpan(false)
    }
  }

  if (objek === null) return galat ? <Gagal galat={galat} /> : <Memuat />

  return (
    <div className="nbf-objek">
      <Gagal galat={galat} />
      {tersimpan && <div className="alert alert--ok">{TEKS_INWARD.tersimpan}</div>}
      <div className="table-wrap">
        <table className="nbf-tabel">
          <thead>
            <tr>
              <th scope="col" />
              {GRID_COV_OBJEK.map((k) => (
                <th key={k.sel} scope="col">
                  {k.label}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {objek.length === 0 && (
              <tr>
                <td colSpan={4}>{TEKS_INWARD.kosong}</td>
              </tr>
            )}
            {objek.map((x, o) => [
              <tr key={`o${o}`}>
                <td>
                  <Buka buka={terbuka.includes(`o${o}`)} onKlik={() => setTerbuka((t) => balik(t, `o${o}`))} />
                </td>
                <td>{x.objectNo}</td>
                <td>{x.objectName}</td>
                <td>{x.riskLocation}</td>
              </tr>,
              terbuka.includes(`o${o}`) && (
                <tr key={`od${o}`} className="nbf-objek__detail">
                  <td colSpan={4}>
                    <div className="nbf-objek__isi">
                      <div className="table-wrap">
                        <table className="nbf-tabel">
                          <thead>
                            <tr>
                              <th scope="col" />
                              {GRID_COV_ITEM.map((k) => (
                                <th key={k.sel} scope="col">
                                  {k.label}
                                </th>
                              ))}
                            </tr>
                          </thead>
                          <tbody>
                            {x.items.length === 0 && (
                              <tr>
                                <td colSpan={6}>{TEKS_INWARD.kosong}</td>
                              </tr>
                            )}
                            {x.items.map((it, i) => {
                              const covs = it.coverages ?? []
                              return [
                                <tr key={`i${i}`}>
                                  <td>
                                    <Buka buka={terbuka.includes(`i${o}-${i}`)} onKlik={() => setTerbuka((t) => balik(t, `i${o}-${i}`))} />
                                  </td>
                                  <td>{it.itemType}</td>
                                  <td>{it.currency}</td>
                                  <td className="nbf-angka">{formatNumber(it.tsi, DESIMAL)}</td>
                                  <td className="nbf-angka">{formatNumber(totalPremiItem(covs, it.totalGrossPremi), DESIMAL)}</td>
                                  <td className="nbf-angka">{formatNumber(it.totalNetRate ?? '', DESIMAL)}</td>
                                </tr>,
                                terbuka.includes(`i${o}-${i}`) && (
                                  <tr key={`id${i}`} className="nbf-objek__detail">
                                    <td colSpan={6}>
                                      <div className="table-wrap">
                                        <table className="nbf-tabel">
                                          <thead>
                                            <tr>
                                              <th scope="col" />
                                              {GRID_COVERAGE.map((k) => (
                                                <th key={k.sel} scope="col">
                                                  {k.label}
                                                </th>
                                              ))}
                                              <th scope="col" className="table__actions">
                                                <button type="button" className="btn btn--ghost btn--sm" onClick={() => void tambah(o, i, covs)}>
                                                  {GRID_OBJEK.tambah}
                                                </button>
                                              </th>
                                            </tr>
                                          </thead>
                                          <tbody>
                                            {covs.length === 0 && (
                                              <tr>
                                                <td colSpan={5}>{TEKS_INWARD.kosong}</td>
                                              </tr>
                                            )}
                                            {covs.map((c, n) => [
                                              <tr key={`c${n}`}>
                                                <td>
                                                  <Buka
                                                    buka={terbuka.includes(`c${o}-${i}-${n}`)}
                                                    onKlik={() => setTerbuka((t) => balik(t, `c${o}-${i}-${n}`))}
                                                  />
                                                </td>
                                                <td>{c.oldId}</td>
                                                <td className="nbf-angka">{formatNumber(c.rateOjk, DESIMAL)}</td>
                                                <td className="nbf-angka">{formatNumber(c.premium, DESIMAL)}</td>
                                                <td className="table__actions">
                                                  <button
                                                    type="button"
                                                    className="btn btn--ghost btn--sm"
                                                    onClick={() => ubahCoverage(o, i, covs.filter((_, k) => k !== n))}
                                                  >
                                                    {GRID_OBJEK.hapus}
                                                  </button>
                                                </td>
                                              </tr>,
                                              terbuka.includes(`c${o}-${i}-${n}`) && (
                                                <tr key={`cd${n}`} className="nbf-objek__detail">
                                                  <td colSpan={5}>
                                                    <FormCoverage
                                                      caseId={caseId}
                                                      c={c}
                                                      tsiItem={it.tsi}
                                                      item={{ isAdjustable: it.isAdjustable, pctAdjustOther: it.pctAdjustOther }}
                                                      ubah={(baru) => ubahCoverage(o, i, covs.map((y, k) => (k === n ? baru : y)))}
                                                    />
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
                              ]
                            })}
                          </tbody>
                        </table>
                      </div>
                      <div className="table-wrap">
                        <table className="nbf-tabel">
                          <thead>
                            <tr>
                              {GRID_COV_TOTAL.map((k) => (
                                <th key={k.sel} scope="col">
                                  {k.label}
                                </th>
                              ))}
                            </tr>
                          </thead>
                          <tbody>
                            {(x.totalPerCurrency ?? []).length === 0 && (
                              <tr>
                                <td colSpan={4}>{TEKS_INWARD.kosong}</td>
                              </tr>
                            )}
                            {(x.totalPerCurrency ?? []).map((t) => (
                              <tr key={t.currency}>
                                <td>{t.currency}</td>
                                <td className="nbf-angka">{formatNumber(t.tsi, DESIMAL)}</td>
                                <td className="nbf-angka">{formatNumber(t.premium, DESIMAL)}</td>
                                <td className="nbf-angka">{formatNumber(t.rate, DESIMAL)}</td>
                              </tr>
                            ))}
                          </tbody>
                        </table>
                      </div>
                    </div>
                  </td>
                </tr>
              ),
            ])}
          </tbody>
        </table>
      </div>
      <div className="nbf-objek__kaki">
        <button type="button" className="btn btn--primary" onClick={() => void simpan()} disabled={menyimpan}>
          {menyimpan ? TEKS_FORM_OPPORTUNITY.menyimpan : SIMPAN_COVERAGE}
        </button>
      </div>
    </div>
  )
}
