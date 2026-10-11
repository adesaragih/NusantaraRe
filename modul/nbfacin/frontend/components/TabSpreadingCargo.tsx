// Tab Spreading MARINE CARGO (tiket 74 P3) - port `NB FacIn\Section\InputInwardFacultativeDtl.xml` tab "Spreading" (tidak
// bercabang per COB: % Share RNM, Copy Spreading, ringkasan) + `CoverageSpreadingList.xml` kontainer L14765 (`IsMarineCargo`).
// Kabel `…/spreading-cargo` (sesi c3). Spesifikasi `docs/issues/74-marine-cargo.md` S4.
//
// Rantai (master-detail, grid hanya dibuka):
// - cargo `CoverageSpreadingList` (No. · Trading · Goods · Packing · Conveyance) → flow action `InputSpreadingMarineCargo_FacIn`
//   → coverage `InputCoverageSpreadingMarineCargo` (Coverage · Standard Rate (%) · Gross Premium · Premium Nusantara Re) →
//   `SpreadingItem` (TSI 100% · TSI Liability · TSI Nusantara Re · Gross Premium · Premium Nusantara Re) + grid baris spreading
//   `InputDtlSpreadingCoverage_FacIn` - SAMA dengan FIRE (baca-saja di luar UW; Add / Delete / Type Treaty / % Share di tahap UW,
//   tiket 59). Copy Spreading hanya tahap UW (K59).
//
// Tanpa spreading otomatis: `GetKapasitasTreaty` tidak punya cabang Marine; baris bawaan 10001 100% (`SetSpreadingCoverage_DT`)
// hanya lewat tombol tersembunyi (S20 `1=2`) - tidak diport. Copy Spreading per lokasi tidak ada (pemanggilnya layar FIRE UW).
// Keputusan agent: kode `TabSpreading` / `TabSpreadingAneka` tidak diubah - bagian yang sama ditulis ulang di sini dengan fungsi
// yang sudah diekspor (`RingkasanSpreading`, `lebihDari100`, `templateLebih`, `treatyKembar`, `barisLebih`, `barisSpreadingKosong`).

import { useEffect, useRef, useState } from 'react'

import { Field, Gagal, Kosong, Memuat, Pilih } from '../../../../inti/frontend/components/ui/dasar'
import { desimalSah } from '../../../../inti/frontend/lib/desimal'
import { angka4 } from '../angka4'
import {
  ambilSpreadingCargo,
  hitungBarisSpreadingCargo,
  hitungShareSpreadingCargo,
  periksaTemplateSpreading,
  salinSpreadingCargo,
  simpanSpreadingCargo,
  type BarisSpreading,
  type CoverageSpreadingCargo,
  type TampilanSpreadingCargo,
  type TemplateSpreading,
} from '../api'
import { SPREADING as S, SPREADING_CARGO as C, TEKS_ANEKA, TEKS_SPREADING as T } from '../labels'
import { barisLengkapCargo, barisSahCargo, gantiBarisCargo, kunciCoverageCargo, type PosisiCoverageCargo } from './spreadingCargo'
import PagerLokasi, { potongHalaman } from './PagerLokasi'
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

export default function TabSpreadingCargo({ caseId, tahapUw = false }: { caseId: string; tahapUw?: boolean }) {
  const [v, setV] = useState<TampilanSpreadingCargo | null>(null)
  const [share, setShare] = useState('')
  const [galat, setGalat] = useState<unknown>(null)
  const [menghitung, setMenghitung] = useState(false)
  const [menyimpan, setMenyimpan] = useState(false)
  const [tersimpan, setTersimpan] = useState(false)
  const [terbuka, setTerbuka] = useState<string[]>([])
  const [halCargo, setHalCargo] = useState(1)
  const [pesanTemplate, setPesanTemplate] = useState<string[]>([])
  const nomor = useRef(0)
  const nomorTemplate = useRef(0)
  const jadwal = useRef<number | undefined>(undefined)

  useEffect(() => {
    let batal = false
    ambilSpreadingCargo(caseId).then(
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
  function hitung(janji: () => Promise<TampilanSpreadingCargo>) {
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
    jadwal.current = window.setTimeout(() => hitung(() => hitungShareSpreadingCargo(caseId, nilai.trim())), JEDA_MS)
  }

  /** Baris spreading satu coverage diubah (tahap UW): hitung-baris berjeda; jawaban yang masih berjalan dibuang. */
  function ubahBaris(p: PosisiCoverageCargo, baris: BarisSpreading[]) {
    if (v === null) return
    const berikut = { ...v, cargo: gantiBarisCargo(v.cargo, p, baris), percentShare: share.trim() }
    setV(berikut)
    setTersimpan(false)
    nomor.current++
    setMenghitung(false)
    window.clearTimeout(jadwal.current)
    if (!barisSahCargo(berikut.cargo) || !angkaSah(share) || share.trim() === '' || lebihDari100(share)) return
    jadwal.current = window.setTimeout(() => hitung(() => hitungBarisSpreadingCargo(caseId, berikut)), JEDA_MS)
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
      const h = await simpanSpreadingCargo(caseId, { ...v, percentShare: share.trim() })
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
  const barisTidakSah = tahapUw && !barisSahCargo(v.cargo)

  // Fungsi render biasa (bukan komponen di dalam komponen - isian kehilangan fokus tiap render).
  /** Detail coverage (`SpreadingItem`, sama dengan FIRE) + grid baris spreading. */
  function detailCoverage(p: PosisiCoverageCargo, c: CoverageSpreadingCargo) {
    return (
      <div className="nbf-spread__coverage">
        <dl className="nbf-spread__angka nbf-aneka-angka">
          {(
            [
              [S.detailCoverage.tsi, c.tsi],
              [S.detailCoverage.tsiLiability, c.tsiLiability],
              [S.detailCoverage.tsiNusantaraRe, c.tsiNusantaraRe],
              [S.detailCoverage.premium, c.premium],
              [S.detailCoverage.premiNusantaraRe, c.premiNusantaraRe],
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

  /** Grid coverage satu cargo (`InputCoverageSpreadingMarineCargo`). */
  function gridCoverage(coverages: CoverageSpreadingCargo[], ci: number) {
    return (
      <div className="nbf-cov-wrap">
        <table className="nbf-tabel nbf-aneka-grid">
          <Kepala kolom={C.kolomCoverage} />
          <tbody>
            {coverages.length === 0 && (
              <tr>
                <td colSpan={C.kolomCoverage.length + 1}>{T.kosong}</td>
              </tr>
            )}
            {coverages.map((c, n) => {
              const p = { c: ci, n }
              const kc = kunciCoverageCargo(p)
              return [
                <tr key={kc} className={terbuka.includes(kc) ? 'nbf-baris--terbuka' : undefined}>
                  <td>
                    <Buka buka={terbuka.includes(kc)} onKlik={() => balik(kc)} />
                  </td>
                  <td>{c.coverageNote}</td>
                  <td className="nbf-angka">{angka(c.rate)}</td>
                  <td className="nbf-angka">{angka(c.premium)}</td>
                  <td className="nbf-angka">{angka(c.premiNusantaraRe)}</td>
                </tr>,
                terbuka.includes(kc) && (
                  <tr key={`${kc}-d`} className="nbf-objek__detail">
                    <td colSpan={C.kolomCoverage.length + 1}>{detailCoverage(p, c)}</td>
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
                onClick={() => hitung(() => salinSpreadingCargo(caseId, share.trim(), v.template))}
              >
                {S.salinSemua}
              </button>
            </div>
          </div>
        )}
      </section>

      {/* Cargo → coverage → spreading; 10 cargo per halaman (work owner 11-10-2026) */}
      <PagerLokasi total={v.cargo.length} halaman={halCargo} onPindah={setHalCargo} />
      <div className="nbf-cov-wrap">
        <table className="nbf-tabel nbf-aneka-grid">
          <Kepala kolom={C.kolomCargo} />
          <tbody>
            {v.cargo.length === 0 && (
              <tr>
                <td colSpan={C.kolomCargo.length + 1}>
                  <Kosong pesan={T.kosong} />
                </td>
              </tr>
            )}
            {potongHalaman(v.cargo, halCargo).map(([x, ci]) => {
              const k = `g${ci}`
              return [
                <tr key={k} className={terbuka.includes(k) ? 'nbf-baris--terbuka' : undefined}>
                  <td>
                    <Buka buka={terbuka.includes(k)} onKlik={() => balik(k)} />
                  </td>
                  <td>{x.no}</td>
                  <td>{x.tradingNote}</td>
                  <td>{x.goodNote}</td>
                  <td>{x.packingNote}</td>
                  <td>{x.conveyanceNote}</td>
                </tr>,
                terbuka.includes(k) && (
                  <tr key={`${k}-d`} className="nbf-objek__detail">
                    <td colSpan={C.kolomCargo.length + 1}>
                      <div className="nbf-lapis nbf-lapis--objek">{gridCoverage(x.coverages, ci)}</div>
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

      {barisTidakSah && !barisLengkapCargo(v.cargo) && <div className="field__error">{T.barisBelumLengkap}</div>}
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
