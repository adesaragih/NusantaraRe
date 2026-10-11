// Tab Coverage layar Inward Facultative, keluarga Aneka (tiket 73 P2). Data = `GET`/`PUT …/objek-aneka` (coverage di bawah
// tiap objek), hitung = `POST …/hitung-coverage-aneka` (sesi c3).
//
// Port berantai (master-detail, pola tab Coverage FIRE):
// 1. Grid lokasi - `DDL\InputDtlCoverage_FacIn.xml` kontainer 7 (L25072 `IsAneka || IsGolfInsurance || IsBonding || IsHE ||
//    IsCustomBonds`): No. · Location.
// 2. `DDL\ObjectOccupation_FacIn.xml` (kontainer !IsEDM): grid "List Okupasi" (Occupation Name), grid Golf (Object Name ·
//    Currency · TSI), total per mata uang (`.Property.TotalTSIPremiGrossList`, server).
// 3. `NB FacIn\Section\ObjectDtlAneka_FacIn.xml`: grid objek (Object Name · Currency · TSI).
// 4. `DDL\InputCoverageAneka_FacIn.xml` (Data-Aneka): grid `.CoverageList` (Coverage · Gross Premium; Add L1710, Trash L2566)
//    + tombol salin kontainer `IsAneka`: Copy Coverage From This Section (L7852, `.IdxLocation = 1`), Copy Deductible To All
//    (L9247, `.IdxLocation = 1 && .pxListSubscript = 1`), Copy Deductible (This Section) (L9497, `.IdxLocation = 1`; aktivitas
//    sama).
// 5. `FormCoverageAneka`.
// Save: seluruh lokasi + % Pro Rate case dikirim `PUT …/objek-aneka` (medan RNM coverage dihitung server saat PUT).
//
// Keputusan agent (tiket 73 P2):
// - D-1 Tidak dibuat (rule tidak ada di ekspor NB, K73-8): Copy Coverage To All / For TPL (`CopyCoverageToAllLocation_Act`
//   Data-Aneka), Copy Coverage To All Golf (`CopyCoverageToALL_HIO_Act`, hanya Endorsment), kontainer "Copy From" (Location
//   Address / Object Item Name / Copy to Another Object Item - `CopyCoverage_Act` dan `D_FilterAneka` hanya Endorsment).
// - D-2 Salin langsung `Obj-Save` di Pega; di sini tersimpan lewat Save tab (pola K53-2).
// - D-3 Total per mata uang = nilai server; diperbarui sesudah Save.

import { useEffect, useRef, useState } from 'react'

import { Gagal, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import { angka4 } from '../angka4'
import {
  ambilObjekAneka,
  daftarMataUang,
  hitungCoverageAneka,
  simpanObjekAneka,
  type CoverageAneka,
  type Deductible,
  type LokasiAneka,
  type ObjekAneka,
  type PemicuAneka,
  type TampilanAneka,
} from '../api'
import {
  GRID_COV_LOKASI_ANEKA,
  GRID_COV_OBJEK_ANEKA,
  GRID_COV_OKUPASI_ANEKA,
  GRID_COV_TOTAL_ANEKA,
  GRID_COVERAGE_ANEKA,
  GRID_OBJEK,
  SALIN_ANEKA,
  SIMPAN_COVERAGE,
  TEKS_ANEKA,
  TEKS_COV_ANEKA,
  TEKS_FORM_OPPORTUNITY,
  TEKS_INWARD,
} from '../labels'
import { rapikanLokasiAneka, tampilOkupasi } from './aneka'
import {
  gantiObjek,
  rapikanCoverageAneka,
  salinCoverageSeksi,
  salinDeductibleAneka,
  setelCurrencyDeductibleTerakhir,
  tambahCoverage,
  ubahCoverageDi,
  type PosisiObjekAneka,
} from './coverageAneka'
import FormCoverageAneka from './FormCoverageAneka'
import PagerLokasi, { potongHalaman } from './PagerLokasi'

/** Kunci baris terbuka (string) per lapisan. */
const kunciObjek = (p: PosisiObjekAneka) => `o${p.l}-${p.k ?? 'g'}-${p.j}`

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

function Kepala({ kolom, aksi }: { kolom: readonly string[]; aksi?: React.ReactNode }) {
  return (
    <thead>
      <tr>
        <th scope="col" />
        {kolom.map((k) => (
          <th key={k} scope="col">
            {k}
          </th>
        ))}
        {aksi !== undefined && (
          <th scope="col" className="table__actions">
            {aksi}
          </th>
        )}
      </tr>
    </thead>
  )
}

/** Objek okupasi / Golf dengan coverage-nya dilengkapi (data P1 / lama). */
function rapikanCoverageObjek(o: ObjekAneka): ObjekAneka {
  return o.coverages ? { ...o, coverages: o.coverages.map(rapikanCoverageAneka) } : { ...o, coverages: [] }
}

function rapikanSemua(l: LokasiAneka[]): LokasiAneka[] {
  return l.map(rapikanLokasiAneka).map((x) => ({
    ...x,
    okupasi: x.okupasi.map((k) => ({ ...k, objek: k.objek.map(rapikanCoverageObjek) })),
    objekGolf: x.objekGolf.map(rapikanCoverageObjek),
  }))
}

export default function TabCoverageAneka({ caseId, tampilan, aneka }: { caseId: string; tampilan: TampilanAneka; aneka: boolean }) {
  const [lokasi, setLokasi] = useState<LokasiAneka[] | null>(null)
  const [proRate, setProRateState] = useState('')
  const proRateRef = useRef('')
  const [mataUang, setMataUang] = useState<string[]>([])
  const [terbuka, setTerbuka] = useState<string[]>([])
  const [halLokasi, setHalLokasi] = useState(1)
  const [galat, setGalat] = useState<unknown>(null)
  const [menghitung, setMenghitung] = useState(0)
  const [menyimpan, setMenyimpan] = useState(false)
  const [tersimpan, setTersimpan] = useState(false)
  const [pesanSalin, setPesanSalin] = useState<{ kunci: string; teks: string } | null>(null)
  const nomorHitung = useRef(new Map<string, number>())

  const setProRate = (v: string) => {
    proRateRef.current = v
    setProRateState(v)
  }

  useEffect(() => {
    let batal = false
    ambilObjekAneka(caseId).then(
      (h) => {
        if (batal) return
        setLokasi(rapikanSemua(h.lokasi))
        setProRate(h.proRatePercent ?? '')
      },
      (err: unknown) => {
        if (!batal) setGalat(err)
      },
    )
    daftarMataUang().then(
      (h) => {
        if (!batal) setMataUang(h.baris)
      },
      () => undefined,
    )
    return () => {
      batal = true
    }
  }, [caseId])

  const balik = (k: string) => setTerbuka((t) => (t.includes(k) ? t.filter((x) => x !== k) : [...t, k]))
  const ganti = (f: (l: LokasiAneka[]) => LokasiAneka[]) => {
    setLokasi((l) => (l === null ? l : f(l)))
    setTersimpan(false)
  }

  /** Rantai on-change satu coverage di backend; jawaban usang (ada hitung lebih baru untuk coverage yang sama) dibuang. */
  function hitung(p: PosisiObjekAneka, n: number, c: CoverageAneka, pemicu: PemicuAneka) {
    const kunci = `${kunciObjek(p)}-${n}`
    const nomor = (nomorHitung.current.get(kunci) ?? 0) + 1
    nomorHitung.current.set(kunci, nomor)
    setMenghitung((x) => x + 1)
    setGalat(null)
    hitungCoverageAneka(caseId, { proRatePercent: proRateRef.current, baris: [{ coverage: c, pemicu }] })
      .then(
        (h) => {
          if (nomorHitung.current.get(kunci) !== nomor) return
          setProRate(h.proRatePercent ?? proRateRef.current)
          const hasil = h.baris[0]
          if (hasil === undefined) return
          // Deductible tidak disentuh rantai hitung: yang di layar dipertahankan.
          ganti((l) => {
            const kini = l[p.l] === undefined ? undefined : (p.k === null ? l[p.l]!.objekGolf[p.j] : l[p.l]!.okupasi[p.k]?.objek[p.j])?.coverages?.[n]
            return ubahCoverageDi(l, p, n, { ...rapikanCoverageAneka(hasil), deductibles: kini?.deductibles ?? c.deductibles })
          })
        },
        (err: unknown) => setGalat(err),
      )
      .finally(() => setMenghitung((x) => x - 1))
  }

  function ubahDeductible(p: PosisiObjekAneka, n: number, lama: Deductible[], baru: Deductible[]) {
    ganti((l) => {
      const hasil = gantiObjek(l, p, (o) => ({ ...o, coverages: (o.coverages ?? []).map((c, i) => (i === n ? { ...c, deductibles: baru } : c)) }))
      // addDeductibleAneka_ACT: Add deductible menyetel Currency deductible terakhir SETIAP coverage.
      return baru.length > lama.length ? setelCurrencyDeductibleTerakhir(hasil) : hasil
    })
  }

  /** Copy Coverage From This Section: pasang salinan, lalu hitung ulang semua coverage sasaran (pemicu 'tsi', satu permintaan). */
  function salinSeksi(sumber: PosisiObjekAneka) {
    if (lokasi === null) return
    const h = salinCoverageSeksi(lokasi, sumber)
    setLokasi(h.lokasi)
    setTersimpan(false)
    setPesanSalin({ kunci: kunciObjek(sumber), teks: TEKS_COV_ANEKA.seksiTersalin(h.sasaran.length) })
    const baris: { p: PosisiObjekAneka; n: number; c: CoverageAneka }[] = []
    for (const p of h.sasaran) {
      const o = p.k === null ? undefined : h.lokasi[p.l]?.okupasi[p.k]?.objek[p.j]
      ;(o?.coverages ?? []).forEach((c, n) => baris.push({ p, n, c }))
    }
    if (baris.length === 0) return
    setMenghitung((x) => x + 1)
    hitungCoverageAneka(caseId, { proRatePercent: proRateRef.current, baris: baris.map((b) => ({ coverage: b.c, pemicu: 'tsi' as const })) })
      .then(
        (r) => {
          setProRate(r.proRatePercent ?? proRateRef.current)
          ganti((l) =>
            baris.reduce((acc, b, i) => {
              const hasil = r.baris[i]
              return hasil === undefined ? acc : ubahCoverageDi(acc, b.p, b.n, { ...rapikanCoverageAneka(hasil), deductibles: b.c.deductibles })
            }, l),
          )
        },
        (err: unknown) => setGalat(err),
      )
      .finally(() => setMenghitung((x) => x - 1))
  }

  function salinDeductible(sumber: PosisiObjekAneka) {
    if (lokasi === null) return
    const h = salinDeductibleAneka(lokasi, sumber)
    setLokasi(h.lokasi)
    setTersimpan(false)
    setPesanSalin({
      kunci: kunciObjek(sumber),
      teks: h.jumlah === 0 ? TEKS_COV_ANEKA.deductibleKosong : TEKS_COV_ANEKA.deductibleTersalin(h.jumlah),
    })
  }

  async function simpan() {
    if (lokasi === null) return
    setMenyimpan(true)
    setGalat(null)
    setTersimpan(false)
    try {
      const h = await simpanObjekAneka(caseId, lokasi, proRateRef.current)
      setLokasi(rapikanSemua(h.lokasi))
      setProRate(h.proRatePercent ?? '')
      setTersimpan(true)
    } catch (err) {
      setGalat(err)
    } finally {
      setMenyimpan(false)
    }
  }

  if (lokasi === null) return galat !== null ? <Gagal galat={galat} /> : <Memuat />

  // Fungsi render biasa (BUKAN komponen di dalam komponen): komponen yang didefinisikan ulang tiap render di-mount ulang
  // dan isian kehilangan fokus setiap ketikan.
  /** Daftar coverage satu objek + tombol salin + form coverage yang dibuka. */
  function daftarCoverage(p: PosisiObjekAneka, o: ObjekAneka) {
    const covs = o.coverages ?? []
    const ko = kunciObjek(p)
    // Tombol salin hanya di kontainer IsAneka (bukan Golf murni) dan objek okupasi.
    const bolehSalin = p.k !== null && aneka
    return (
      <div className="nbf-lapis">
        <div className="nbf-cov-wrap">
          <table className="nbf-tabel nbf-aneka-grid nbf-aneka-grid--aksi">
            <Kepala
              kolom={GRID_COVERAGE_ANEKA}
              aksi={
                <button type="button" className="btn nbf-tema__tombol btn--ghost btn--sm" onClick={() => ganti((l) => tambahCoverage(l, p))}>
                  {GRID_OBJEK.tambah}
                </button>
              }
            />
            <tbody>
              {covs.length === 0 && (
                <tr>
                  <td colSpan={4}>{TEKS_COV_ANEKA.kosongCoverage}</td>
                </tr>
              )}
              {covs.map((c, n) => {
                const kc = `${ko}-c${n}`
                return [
                  <tr key={kc} className={terbuka.includes(kc) ? 'nbf-baris--terbuka' : undefined}>
                    <td>
                      <Buka buka={terbuka.includes(kc)} onKlik={() => balik(kc)} />
                    </td>
                    <td>{c.coverageNote === '' ? <span className="nbf-aneka-grid__kosong">{TEKS_ANEKA.belumDiisi}</span> : c.coverageNote}</td>
                    <td className="nbf-angka">{angka4(c.premium)}</td>
                    <td className="table__actions">
                      <button
                        type="button"
                        className="btn nbf-tema__tombol btn--ghost btn--sm"
                        onClick={() => ganti((l) => gantiObjek(l, p, (x) => ({ ...x, coverages: (x.coverages ?? []).filter((_, i) => i !== n) })))}
                      >
                        {GRID_OBJEK.hapus}
                      </button>
                    </td>
                  </tr>,
                  terbuka.includes(kc) && (
                    <tr key={`${kc}-d`} className="nbf-objek__detail">
                      <td colSpan={4}>
                        <FormCoverageAneka
                          caseId={caseId}
                          idDasar={kc}
                          c={c}
                          ubah={(baru) => ganti((l) => ubahCoverageDi(l, p, n, baru))}
                          hitung={(baru, pemicu) => hitung(p, n, baru, pemicu)}
                          tampilan={tampilan}
                          mataUang={mataUang}
                          proRate={proRate}
                          ubahProRate={(v) => {
                            setProRate(v)
                            setTersimpan(false)
                          }}
                          ubahDeductible={(d) => ubahDeductible(p, n, c.deductibles, d)}
                        />
                      </td>
                    </tr>
                  ),
                ]
              })}
            </tbody>
          </table>
        </div>
        {bolehSalin && p.l === 0 && (
          <div className="nbf-cov-salin">
            <button type="button" className="btn nbf-tema__tombol btn--sm" onClick={() => salinSeksi(p)}>
              {SALIN_ANEKA.dariSeksi}
            </button>
            {p.j === 0 && (
              <button type="button" className="btn nbf-tema__tombol btn--sm" onClick={() => salinDeductible(p)}>
                {SALIN_ANEKA.deductibleSemua}
              </button>
            )}
            <button type="button" className="btn nbf-tema__tombol btn--sm" onClick={() => salinDeductible(p)}>
              {SALIN_ANEKA.deductibleSeksi}
            </button>
          </div>
        )}
        {pesanSalin?.kunci === ko && (
          <div className="alert alert--ok" role="status">
            {pesanSalin.teks}
          </div>
        )}
      </div>
    )
  }

  /** Grid objek (okupasi atau Golf) → daftar coverage. */
  function gridObjekCoverage(daftar: ObjekAneka[], posisi: (j: number) => PosisiObjekAneka) {
    return (
      <div className="nbf-cov-wrap">
        <table className="nbf-tabel nbf-aneka-grid">
          <Kepala kolom={GRID_COV_OBJEK_ANEKA} />
          <tbody>
            {daftar.length === 0 && (
              <tr>
                <td colSpan={4}>{TEKS_ANEKA.kosongObjek}</td>
              </tr>
            )}
            {daftar.map((o, j) => {
              const p = posisi(j)
              const k = kunciObjek(p)
              return [
                <tr key={k} className={terbuka.includes(k) ? 'nbf-baris--terbuka' : undefined}>
                  <td>
                    <Buka buka={terbuka.includes(k)} onKlik={() => balik(k)} />
                  </td>
                  <td>{o.objectName === '' ? <span className="nbf-aneka-grid__kosong">{TEKS_ANEKA.belumDiisi}</span> : o.objectName}</td>
                  <td>{o.currency}</td>
                  <td className="nbf-angka">{angka4(o.tsi)}</td>
                </tr>,
                terbuka.includes(k) && (
                  <tr key={`${k}-d`} className="nbf-objek__detail">
                    <td colSpan={4}>
                      {daftarCoverage(p, o)}
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
    <div className="nbf-objek nbf-cov-tab">
      <Gagal galat={galat} />
      {menghitung > 0 && <div className="muted">{TEKS_COV_ANEKA.menghitung}</div>}
      {tersimpan && <div className="alert alert--ok">{TEKS_COV_ANEKA.tersimpan}</div>}
      {/* Lokasi - 10 per halaman (work owner 11-10-2026) */}
      <PagerLokasi total={lokasi.length} halaman={halLokasi} onPindah={setHalLokasi} />
      <div className="nbf-cov-wrap">
        <table className="nbf-tabel nbf-aneka-grid nbf-aneka-grid--lokasi">
          <Kepala kolom={GRID_COV_LOKASI_ANEKA} />
          <tbody>
            {lokasi.length === 0 && (
              <tr>
                <td colSpan={3}>{TEKS_INWARD.kosong}</td>
              </tr>
            )}
            {potongHalaman(lokasi, halLokasi).map(([l, a]) => {
              const kl = `l${a}`
              return [
                <tr key={kl} className={terbuka.includes(kl) ? 'nbf-baris--terbuka' : undefined}>
                  <td>
                    <Buka buka={terbuka.includes(kl)} onKlik={() => balik(kl)} />
                  </td>
                  <td>{l.objectNo}</td>
                  <td>{l.riskLocation}</td>
                </tr>,
                terbuka.includes(kl) && (
                  <tr key={`${kl}-d`} className="nbf-objek__detail">
                    <td colSpan={3}>
                      <div className="nbf-objek__isi nbf-lapis nbf-lapis--objek">
                        {tampilOkupasi(aneka, tampilan) && (
                          <div className="nbf-cov-wrap">
                            <table className="nbf-tabel nbf-aneka-grid">
                              <Kepala kolom={GRID_COV_OKUPASI_ANEKA} />
                              <tbody>
                                {l.okupasi.length === 0 && (
                                  <tr>
                                    <td colSpan={2}>{TEKS_ANEKA.kosongOkupasi}</td>
                                  </tr>
                                )}
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
                                        <td colSpan={2}>
                                          {gridObjekCoverage(k.objek, (j) => ({ l: a, k: b, j }))}
                                        </td>
                                      </tr>
                                    ),
                                  ]
                                })}
                              </tbody>
                            </table>
                          </div>
                        )}
                        {tampilan.golf && gridObjekCoverage(l.objekGolf, (j) => ({ l: a, k: null, j }))}
                        <div className="nbf-cov-wrap">
                          <table className="nbf-tabel nbf-aneka-total">
                            <thead>
                              <tr>
                                {GRID_COV_TOTAL_ANEKA.map((k) => (
                                  <th key={k} scope="col">
                                    {k}
                                  </th>
                                ))}
                              </tr>
                            </thead>
                            <tbody>
                              {(l.totalPerCurrency ?? []).length === 0 && (
                                <tr>
                                  <td colSpan={4}>{TEKS_INWARD.kosong}</td>
                                </tr>
                              )}
                              {(l.totalPerCurrency ?? []).map((t) => (
                                <tr key={t.currency}>
                                  <td>{t.currency}</td>
                                  <td className="nbf-angka">{angka4(t.tsi)}</td>
                                  <td className="nbf-angka">{angka4(t.premium)}</td>
                                  <td className="nbf-angka">{angka4(t.rate)}</td>
                                </tr>
                              ))}
                            </tbody>
                          </table>
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
      <div className="nbf-objek__kaki">
        <button
          type="button"
          className="btn nbf-tema__tombol nbf-tema__tombol--utama btn--primary"
          onClick={() => void simpan()}
          disabled={menyimpan || menghitung > 0}
        >
          {menyimpan ? TEKS_FORM_OPPORTUNITY.menyimpan : SIMPAN_COVERAGE}
        </button>
      </div>
    </div>
  )
}
