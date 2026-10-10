// Tab Spreading keluarga Aneka (tiket 73 P3) - port `NB FacIn\Section\InputInwardFacultativeDtl.xml` tab "Spreading" (tidak
// bercabang per COB: % Share RNM, Copy Spreading, ringkasan) + `CoverageSpreadingList.xml` kontainer L3010 (`isGolfInsurance ||
// IsAneka || IsBondingAndCustomBonds`). Kabel `…/spreading-aneka` (sesi c3).
//
// Rantai (master-detail, tanpa tombol - grid hanya dibuka):
// - lokasi (No · Location) → okupasi `PropertyItemListCoverageSpreading.xml` L18225 (Occupation Name) → objek
//   `SpreadingObjectAneka` (Object Name · Currency · TSI Old · TSI) → coverage `SpreadingCoverageAneka` (Coverage · Gross Premium)
//   → `SpreadingItemAneka` (Gross TSI · TSI Nusantara Re · Gross Premium · Premium Nusantara Re) + grid baris spreading;
// - Golf L16393: lokasi → objek (Object Name) → coverage `InputCoverageSpreadingAneka` (Coverage Note · % Loading · % Rate · TSI ·
//   Premi) → detail yang sama.
// Grid baris spreading = form FIRE `InputDtlSpreadingCoverage_FacIn(_IsUW)` apa adanya (sama dengan `TabSpreading`: baca-saja di
// luar UW; Add / Delete / Type Treaty / % Share di tahap UW, tiket 59). Copy Spreading hanya tahap UW (K59).
//
// Keputusan work owner 09-10-2026: K73-17 spreading otomatis hanya Engineering / Golf (backend); K73-18 TSI 100% Spreaded Aneka
// kosong; K73-19 total per lokasi hanya bila `tampilTotalLokasi` (MBD tahap UW); K73-20 kuirk Golf ditiru (backend).
// Keputusan agent: kode `TabSpreading` FIRE tidak diubah (layar FIRE yang sudah berjalan) - bagian yang sama ditulis ulang di sini
// dengan memakai fungsi yang sudah diekspor (`RingkasanSpreading`, `lebihDari100`, `templateLebih`, `treatyKembar`,
// `barisLebih`, `barisSpreadingKosong`).

import { useEffect, useRef, useState } from 'react'

import { Field, Gagal, Kosong, Memuat, Pilih } from '../../../../inti/frontend/components/ui/dasar'
import { desimalSah } from '../../../../inti/frontend/lib/desimal'
import { angka4 } from '../angka4'
import {
  ambilSpreadingAneka,
  hitungBarisSpreadingAneka,
  hitungShareSpreadingAneka,
  periksaTemplateSpreading,
  salinLokasiSpreadingAneka,
  salinSpreadingAneka,
  simpanSpreadingAneka,
  type BarisSpreading,
  type CoverageSpreadingAneka,
  type ObjekSpreadingAneka,
  type TampilanSpreadingAneka,
  type TemplateSpreading,
  type TotalSpreading,
} from '../api'
import { SPREADING as S, SPREADING_ANEKA as A, TEKS_ANEKA, TEKS_SPREADING as T } from '../labels'
import CopySpreadingLokasi from './CopySpreadingLokasi'
import { barisLengkapAneka, barisSahAneka, gantiBarisAneka, kunciCoverageSpreading, type PosisiCoverageSpreading } from './spreadingAneka'
import { RingkasanSpreading, barisLebih, barisSpreadingKosong, lebihDari100, templateLebih, treatyKembar } from './TabSpreading'

const JEDA_MS = 500
const angkaSah = (v: string) => v.trim() === '' || desimalSah(v.trim())
const angka = (v: string | undefined) => angka4(v ?? '')

function Buka({ buka, onKlik }: { buka: boolean; onKlik: () => void }) {
  return (
    <button
      type="button"
      className={'btn btn--ghost btn--sm nbf-buka' + (buka ? ' nbf-buka--terbuka' : '')}
      aria-label={TEKS_ANEKA.bukaBaris}
      aria-expanded={buka}
      onClick={onKlik}
    >
      ▸
    </button>
  )
}

function Kepala({ kolom }: { kolom: readonly string[] }) {
  return (
    <thead>
      <tr>
        <th scope="col" />
        {kolom.map((k) => (
          <th key={k} scope="col">
            {k}
          </th>
        ))}
      </tr>
    </thead>
  )
}

/** Total per lokasi (`.Property.TotalTSIPremiSpreadRNM`) - kolom `SPREADING.kolomTotal` (FIRE). */
function TabelTotalLokasi({ baris }: { baris: TotalSpreading[] }) {
  return (
    <div className="nbf-cov-wrap">
      <table className="nbf-tabel">
        <thead>
          <tr>
            {S.kolomTotal.map((k) => (
              <th key={k} scope="col">
                {k}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {baris.length === 0 && (
            <tr>
              <td colSpan={S.kolomTotal.length}>{T.kosong}</td>
            </tr>
          )}
          {baris.map((t, i) => (
            <tr key={`${t.currency}-${t.treatyType}-${i}`}>
              <td>{t.currency}</td>
              <td>{t.treatyName}</td>
              <td className="nbf-angka">{angka(t.sharePercentage)}</td>
              <td className="nbf-angka">{angka(t.claimSpreaded)}</td>
              <td className="nbf-angka">{angka(t.tsiSpreaded)}</td>
              <td className="nbf-angka">{angka(t.claimEstimation)}</td>
              <td className="nbf-angka">{angka(t.premiumSpreaded)}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

export default function TabSpreadingAneka({ caseId, tahapUw = false }: { caseId: string; tahapUw?: boolean }) {
  const [v, setV] = useState<TampilanSpreadingAneka | null>(null)
  const [share, setShare] = useState('')
  const [galat, setGalat] = useState<unknown>(null)
  const [menghitung, setMenghitung] = useState(false)
  const [menyimpan, setMenyimpan] = useState(false)
  const [tersimpan, setTersimpan] = useState(false)
  const [terbuka, setTerbuka] = useState<string[]>([])
  const [pesanTemplate, setPesanTemplate] = useState<string[]>([])
  const nomor = useRef(0)
  const nomorTemplate = useRef(0)
  const jadwal = useRef<number | undefined>(undefined)

  useEffect(() => {
    let batal = false
    ambilSpreadingAneka(caseId).then(
      (h) => {
        if (batal) return
        setV(h)
        setShare(h.percentShare)
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

  const balik = (k: string) => setTerbuka((t) => (t.includes(k) ? t.filter((x) => x !== k) : [...t, k]))

  /** Satu panggilan hitung; jawaban lama dibuang. */
  function hitung(janji: () => Promise<TampilanSpreadingAneka>) {
    const n = ++nomor.current
    setMenghitung(true)
    setGalat(null)
    janji().then(
      (h) => {
        if (n !== nomor.current) return
        setV(h)
        setMenghitung(false)
      },
      (err: unknown) => {
        if (n !== nomor.current) return
        setGalat(err)
        setMenghitung(false)
      },
    )
  }

  function ubahShare(nilai: string) {
    setShare(nilai)
    setTersimpan(false)
    window.clearTimeout(jadwal.current)
    if (!angkaSah(nilai) || nilai.trim() === '' || lebihDari100(nilai)) return
    jadwal.current = window.setTimeout(() => hitung(() => hitungShareSpreadingAneka(caseId, nilai.trim())), JEDA_MS)
  }

  /** Baris spreading satu coverage diubah (tahap UW): hitung-baris berjeda; jawaban yang masih berjalan dibuang. */
  function ubahBaris(p: PosisiCoverageSpreading, baris: BarisSpreading[]) {
    if (v === null) return
    const berikut = { ...v, lokasi: gantiBarisAneka(v.lokasi, p, baris), percentShare: share.trim() }
    setV(berikut)
    setTersimpan(false)
    nomor.current++
    setMenghitung(false)
    window.clearTimeout(jadwal.current)
    if (!barisSahAneka(berikut.lokasi) || !angkaSah(share) || share.trim() === '' || lebihDari100(share)) return
    jadwal.current = window.setTimeout(() => hitung(() => hitungBarisSpreadingAneka(caseId, berikut)), JEDA_MS)
  }

  /** Template Copy Spreading; `periksa` = Type Treaty berubah / baris dihapus (proteksi template di server). */
  function ubahTemplate(template: TemplateSpreading[], periksa = false) {
    setV((x) => (x ? { ...x, template } : x))
    setTersimpan(false)
    if (!periksa) return
    const n = ++nomorTemplate.current
    if (template.every((x) => x.treatyType.trim() === '')) {
      setPesanTemplate([])
      return
    }
    periksaTemplateSpreading(caseId, template).then(
      (h) => {
        if (n === nomorTemplate.current) setPesanTemplate(h.pesan)
      },
      (err: unknown) => {
        if (n === nomorTemplate.current) setGalat(err)
      },
    )
  }

  async function simpan() {
    if (v === null) return
    setMenyimpan(true)
    setGalat(null)
    try {
      const h = await simpanSpreadingAneka(caseId, { ...v, percentShare: share.trim() })
      setV(h)
      setShare(h.percentShare)
      setTersimpan(true)
    } catch (err) {
      setGalat(err)
    } finally {
      setMenyimpan(false)
    }
  }

  if (v === null) return galat ? <Gagal galat={galat} /> : <Memuat />

  const opsiTreaty = v.treaty.map((t) => ({ value: t.id, label: t.name }))
  const lebih = templateLebih(v.template)
  const kembar = treatyKembar(v.template)
  const barisTidakSah = tahapUw && !barisSahAneka(v.lokasi)

  // Fungsi render biasa (bukan komponen di dalam komponen - isian kehilangan fokus tiap render).
  /** Detail coverage (`SpreadingItemAneka`) + grid baris spreading. */
  function detailCoverage(p: PosisiCoverageSpreading, c: CoverageSpreadingAneka) {
    return (
      <div className="nbf-spread__coverage">
        <dl className="nbf-spread__angka nbf-aneka-angka">
          {(
            [
              [A.detailCoverage.tsi, c.tsi],
              [A.detailCoverage.tsiNusantaraRe, c.tsiNusantaraRe],
              [A.detailCoverage.premium, c.premium],
              [A.detailCoverage.premiNusantaraRe, c.premiNusantaraRe],
            ] as const
          ).map(([lbl, nilai]) => (
            <div key={lbl}>
              <dt>{lbl}</dt>
              <dd className="nbf-angka">{angka(nilai)}</dd>
            </div>
          ))}
        </dl>
        <div className="nbf-cov-wrap">
          <table className="nbf-tabel">
            <thead>
              <tr>
                {S.kolomSpreading.map((k) => (
                  <th key={k} scope="col">
                    {k}
                  </th>
                ))}
                {tahapUw && (
                  <th scope="col" className="table__actions">
                    <button type="button" className="btn nbf-tema__tombol btn--sm" onClick={() => ubahBaris(p, [...c.spreading, barisSpreadingKosong()])}>
                      {S.tambah}
                    </button>
                  </th>
                )}
              </tr>
            </thead>
            <tbody>
              {c.spreading.length === 0 && (
                <tr>
                  <td colSpan={tahapUw ? 7 : 6}>{T.kosong}</td>
                </tr>
              )}
              {c.spreading.map((b, k) => (
                <tr key={tahapUw ? k : `${b.treatyType}-${k}`}>
                  {tahapUw ? (
                    <>
                      <td>
                        <Pilih
                          label={S.kolomSpreading[0]}
                          value={b.treatyType}
                          kosong={T.pilihTreaty}
                          opsi={opsiTreaty}
                          onChange={(id) =>
                            ubahBaris(
                              p,
                              c.spreading.map((x, m) => (m === k ? { ...x, treatyType: id, treatyName: v!.treaty.find((y) => y.id === id)?.name ?? '' } : x)),
                            )
                          }
                        />
                      </td>
                      <td>
                        <Field
                          label={S.kolomSpreading[1]}
                          value={b.sharePercentage}
                          onChange={(s) => ubahBaris(p, c.spreading.map((x, m) => (m === k ? { ...x, sharePercentage: s } : x)))}
                          error={angkaSah(b.sharePercentage) ? undefined : T.angka}
                        />
                      </td>
                    </>
                  ) : (
                    <>
                      <td>{b.treatyName || b.treatyType}</td>
                      <td className="nbf-angka">{angka(b.sharePercentage)}</td>
                    </>
                  )}
                  <td className="nbf-angka">{angka(b.tsiGrossSpreaded)}</td>
                  <td className="nbf-angka">{angka(b.tsiSpreaded)}</td>
                  <td className="nbf-angka">{angka(b.claimEstimation)}</td>
                  <td className="nbf-angka">{angka(b.premiumSpreaded)}</td>
                  {tahapUw && (
                    <td className="table__actions">
                      <button type="button" className="btn nbf-tema__tombol btn--sm" onClick={() => ubahBaris(p, c.spreading.filter((_, m) => m !== k))}>
                        {S.hapus}
                      </button>
                    </td>
                  )}
                </tr>
              ))}
            </tbody>
          </table>
        </div>
        {tahapUw && barisLebih(c.spreading) && <div className="field__error">{T.shareBarisLebih}</div>}
      </div>
    )
  }

  /** Grid coverage satu objek: Aneka (Coverage · Gross Premium) atau Golf (Coverage Note · % Loading · % Rate · TSI · Premi). */
  function gridCoverage(o: ObjekSpreadingAneka, posisi: (n: number) => PosisiCoverageSpreading, golf: boolean) {
    const kolom = golf ? A.kolomCoverageGolf : A.kolomCoverage
    return (
      <div className="nbf-cov-wrap">
        <table className="nbf-tabel nbf-aneka-grid">
          <Kepala kolom={kolom} />
          <tbody>
            {o.coverages.length === 0 && (
              <tr>
                <td colSpan={kolom.length + 1}>{T.kosong}</td>
              </tr>
            )}
            {o.coverages.map((c, n) => {
              const p = posisi(n)
              const kc = kunciCoverageSpreading(p)
              return [
                <tr key={kc} className={terbuka.includes(kc) ? 'nbf-baris--terbuka' : undefined}>
                  <td>
                    <Buka buka={terbuka.includes(kc)} onKlik={() => balik(kc)} />
                  </td>
                  <td>{c.coverageNote}</td>
                  {golf && <td className="nbf-angka">{angka(c.loading)}</td>}
                  {golf && <td className="nbf-angka">{angka(c.rate)}</td>}
                  {golf && <td className="nbf-angka">{angka(c.tsi)}</td>}
                  <td className="nbf-angka">{angka(c.premium)}</td>
                </tr>,
                terbuka.includes(kc) && (
                  <tr key={`${kc}-d`} className="nbf-objek__detail">
                    <td colSpan={kolom.length + 1}>{detailCoverage(p, c)}</td>
                  </tr>
                ),
              ]
            })}
          </tbody>
        </table>
      </div>
    )
  }

  /** Grid objek: Aneka (Object Name · Currency · TSI Old · TSI) atau Golf (Object Name). */
  function gridObjek(daftar: ObjekSpreadingAneka[], l: number, k: number | null) {
    const golf = k === null
    const kolom = golf ? A.kolomObjekGolf : A.kolomObjek
    return (
      <div className="nbf-cov-wrap">
        <table className="nbf-tabel nbf-aneka-grid">
          <Kepala kolom={kolom} />
          <tbody>
            {daftar.length === 0 && (
              <tr>
                <td colSpan={kolom.length + 1}>{TEKS_ANEKA.kosongObjek}</td>
              </tr>
            )}
            {daftar.map((o, j) => {
              const ko = `o${l}-${k ?? 'g'}-${j}`
              return [
                <tr key={ko} className={terbuka.includes(ko) ? 'nbf-baris--terbuka' : undefined}>
                  <td>
                    <Buka buka={terbuka.includes(ko)} onKlik={() => balik(ko)} />
                  </td>
                  <td>{o.objectName}</td>
                  {!golf && <td>{o.currency}</td>}
                  {!golf && <td className="nbf-angka">{angka(o.tsi)}</td>}
                </tr>,
                terbuka.includes(ko) && (
                  <tr key={`${ko}-d`} className="nbf-objek__detail">
                    <td colSpan={kolom.length + 1}>
                      <div className="nbf-lapis">{gridCoverage(o, (n) => ({ l, k, j, n }), golf)}</div>
                    </td>
                  </tr>
                ),
              ]
            })}
          </tbody>
        </table>
      </div>
    )
  }

  return (
    <div className="nbf-objek nbf-cov-tab nbf-spread">
      <Gagal galat={galat} />
      {tersimpan && <div className="alert alert--ok">{T.tersimpan}</div>}
      {v.pesan.map((p) => (
        <div key={p} className="alert alert--warn">
          {p}
        </div>
      ))}

      {/* % Share RNM + Copy Spreading (tab, tidak bercabang per COB) */}
      <section className="nbf-lapis nbf-spread__atas">
        <div className="nbf-spread__share">
          <Field
            label={S.percentShare}
            value={share}
            onChange={ubahShare}
            required
            error={!angkaSah(share) ? T.angka : lebihDari100(share) ? T.percentShareLebih : undefined}
          />
          {menghitung && <span className="muted">{T.menghitung}</span>}
        </div>
        {tahapUw && (
          <div className="nbf-spread__copy">
            <h5 className="nbf-akum__kartu-judul">{S.judulCopy}</h5>
            <div className="nbf-cov-wrap">
              <table className="nbf-tabel">
                <thead>
                  <tr>
                    {S.kolomTemplate.map((k) => (
                      <th key={k} scope="col">
                        {k}
                      </th>
                    ))}
                    <th scope="col" className="table__actions">
                      <button
                        type="button"
                        className="btn nbf-tema__tombol btn--sm"
                        onClick={() => ubahTemplate([...v.template, { treatyType: '', treatyName: '', sharePercentage: '' }])}
                      >
                        {S.tambah}
                      </button>
                    </th>
                  </tr>
                </thead>
                <tbody>
                  {v.template.length === 0 && (
                    <tr>
                      <td colSpan={3}>{T.kosong}</td>
                    </tr>
                  )}
                  {v.template.map((t, i) => (
                    <tr key={i}>
                      <td>
                        <Pilih
                          label={S.kolomTemplate[0]}
                          value={t.treatyType}
                          kosong={T.pilihTreaty}
                          opsi={opsiTreaty}
                          onChange={(id) =>
                            ubahTemplate(
                              v.template.map((x, k) => (k === i ? { ...x, treatyType: id, treatyName: v.treaty.find((y) => y.id === id)?.name ?? '' } : x)),
                              true,
                            )
                          }
                        />
                      </td>
                      <td>
                        <Field
                          label={S.kolomTemplate[1]}
                          value={t.sharePercentage}
                          onChange={(s) => ubahTemplate(v.template.map((x, k) => (k === i ? { ...x, sharePercentage: s } : x)))}
                          error={angkaSah(t.sharePercentage) ? undefined : T.angka}
                        />
                      </td>
                      <td className="table__actions">
                        <button type="button" className="btn nbf-tema__tombol btn--sm" onClick={() => ubahTemplate(v.template.filter((_, k) => k !== i), true)}>
                          {S.hapus}
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
            {lebih && <div className="field__error">{T.shareLebih}</div>}
            {kembar && <div className="field__error">{T.treatySama}</div>}
            {pesanTemplate
              .filter((p) => !(kembar && p === T.treatySama))
              .map((p) => (
                <div key={p} className="field__error">
                  {p}
                </div>
              ))}
            <div>
              <button
                type="button"
                className="btn nbf-tema__tombol nbf-tema__tombol--utama btn--primary btn--sm"
                disabled={lebih || kembar || pesanTemplate.length > 0 || v.template.length === 0 || menghitung}
                onClick={() => hitung(() => salinSpreadingAneka(caseId, share.trim(), v.template))}
              >
                {S.salinSemua}
              </button>
            </div>
          </div>
        )}
      </section>

      {/* Lokasi → okupasi | Golf → objek → coverage → spreading */}
      <div className="nbf-cov-wrap">
        <table className="nbf-tabel nbf-aneka-grid nbf-aneka-grid--lokasi">
          <Kepala kolom={A.kolomLokasi} />
          <tbody>
            {v.lokasi.length === 0 && (
              <tr>
                <td colSpan={3}>
                  <Kosong pesan={T.kosong} />
                </td>
              </tr>
            )}
            {v.lokasi.map((l, a) => {
              const kl = `l${a}`
              return [
                <tr key={kl} className={terbuka.includes(kl) ? 'nbf-baris--terbuka' : undefined}>
                  <td>
                    <Buka buka={terbuka.includes(kl)} onKlik={() => balik(kl)} />
                  </td>
                  <td>{l.objectNo}</td>
                  <td>{l.location}</td>
                </tr>,
                terbuka.includes(kl) && (
                  <tr key={`${kl}-d`} className="nbf-objek__detail">
                    <td colSpan={3}>
                      <div className="nbf-lapis nbf-lapis--objek">
                        {/* K73-40: Copy Spreading per lokasi, tahap UW - template case di layar dipertahankan. */}
                        {tahapUw && (
                          <CopySpreadingLokasi
                            treaty={v.treaty}
                            menghitung={menghitung}
                            onSalin={(t) =>
                              hitung(() => salinLokasiSpreadingAneka(caseId, share.trim(), a + 1, t).then((h) => ({ ...h, template: v.template })))
                            }
                          />
                        )}
                        {l.okupasi.length > 0 && (
                          <div className="nbf-cov-wrap">
                            <table className="nbf-tabel nbf-aneka-grid">
                              <Kepala kolom={A.kolomOkupasi} />
                              <tbody>
                                {l.okupasi.map((k, b) => {
                                  const kk = `k${a}-${b}`
                                  return [
                                    <tr key={kk} className={terbuka.includes(kk) ? 'nbf-baris--terbuka' : undefined}>
                                      <td>
                                        <Buka buka={terbuka.includes(kk)} onKlik={() => balik(kk)} />
                                      </td>
                                      <td>{k.occupationName}</td>
                                    </tr>,
                                    terbuka.includes(kk) && (
                                      <tr key={`${kk}-d`} className="nbf-objek__detail">
                                        <td colSpan={2}>{gridObjek(k.objek, a, b)}</td>
                                      </tr>
                                    ),
                                  ]
                                })}
                              </tbody>
                            </table>
                          </div>
                        )}
                        {l.objekGolf.length > 0 && gridObjek(l.objekGolf, a, null)}
                        {v.tampilTotalLokasi && (
                          <div className="nbf-spread__total">
                            <TabelTotalLokasi baris={l.total} />
                          </div>
                        )}
                      </div>
                    </td>
                  </tr>
                ),
              ]
            })}
          </tbody>
        </table>
      </div>

      {/* Ringkasan (SummarySpreading_Section, tidak bercabang per COB) */}
      <RingkasanSpreading v={v} />

      {barisTidakSah && !barisLengkapAneka(v.lokasi) && <div className="field__error">{T.barisBelumLengkap}</div>}
      <div className="nbf-objek__kaki">
        <button
          type="button"
          className="btn nbf-tema__tombol nbf-tema__tombol--utama btn--primary"
          onClick={() => void simpan()}
          disabled={menyimpan || !angkaSah(share) || lebihDari100(share) || barisTidakSah}
        >
          {menyimpan ? T.menghitung : S.simpan}
        </button>
      </div>
    </div>
  )
}
