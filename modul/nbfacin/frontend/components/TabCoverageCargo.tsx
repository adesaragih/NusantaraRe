// Tab Coverage layar Inward Facultative, Class of Business MARINE CARGO (tiket 74 P2). Kabel `GET`/`PUT …/objek-cargo` (coverage
// di dalam cargo) + `POST …/hitung-coverage-cargo` (sesi c3). Bukti `NB\` = `D:\migrasi\RNM\NB FacIn\`; spesifikasi
// `docs/issues/74-marine-cargo.md` S3.
//
// Rantai layar (master-detail expandPane):
// 1. Grid cargo `WORK!INPUTDTLCOVERAGE_FACIN` S131 (`IsMarineCargo` L388): No. / Trading / Goods / Packing / Conveyance.
// 2. Per cargo (`DATA-CARGO!INPUTDTLCOVERAGE_FACIN`, pre-activity `AppendCoverageMarineCargo_ACT`): grid coverage Coverage Note /
//    Coverage Initial - ikon tambah / hapus TERSEMBUNYI di Pega (L998 `NEVER`), jadi satu coverage per cargo (server
//    menyintesis satu coverage kosong bila belum ada) - lalu grid `TotalTSIPremiGrossList` (dihitung server).
// 3. Form coverage `InputCoverageCargo_FacIn.xml`: Coverage Note (baca-saja, diisi Select Coverage) / Coverage Initial /
//    Select Coverage; Currency / Rate (%) / Limit of Liability / Discount (%); TSI * / Premium / Min Premium / Discount (baca-saja);
//    grid Warranty Extension `.AdditionalCoverage`; grid Deductible (editor `addDeductibleCargo` = pola FIRE) + Deskripsi.
//
// Rantai on-change dikerjakan backend saat isian SELESAI: Currency 'currency' (L2787) · Rate 'rate' (L3366) · TSI 'tsi' (L5693) ·
// Discount (%) 'diskonPersen' (L5104) · Premium 'premium' (L6120 `CountTSIandRateMarineCargo_Act`).
//
// Tidak diport (dicatat di tiket 74): tombol "Copy Deductible To All" (`CopyDeductibleToAllCargo_Act` tidak ada di korpus);
// cabang nomor polis master literal (Currency Master / Get Currency Master, K74-5); MOP master Premium = 0 (K74-3).

import { useEffect, useRef, useState } from 'react'

import { Gagal, Memuat, Pilih, type Opsi } from '../../../../inti/frontend/components/ui/dasar'
import { angka4 } from '../angka4'
import {
  ambilObjekCargo,
  cariKondisiMarine,
  daftarMataUang,
  hitungCoverageCargo,
  simpanObjekCargo,
  type CoverageCargo,
  type ObjekCargo,
  type PemicuCargo,
} from '../api'
import {
  FORM_COVERAGE_CARGO as F,
  GRID_CARGO_COVERAGE,
  GRID_COVERAGE_CARGO,
  SIMPAN_OBJEK,
  TEKS_ANEKA,
  TEKS_CARGO,
  TEKS_FORM_OPPORTUNITY,
} from '../labels'
import { pesanProteksiDari } from './BlokEmailType'
import { rapikanCargo, rapikanCoverageCargo, terapkanHitungCargo } from './cargo'
import GridDeductible from './GridDeductible'
import IsianAngka4 from './IsianAngka4'
import IsianSaran from './IsianSaran'
import PagerLokasi, { potongHalaman } from './PagerLokasi'

/** Isian angka berlabel; `onSelesai` hanya bila nilainya diubah sejak fokus (event change Pega). */
function Angka({ id, label, value, onChange, onSelesai, wajib }: { id: string; label: string; value: string; onChange: (v: string) => void; onSelesai?: () => void; wajib?: boolean }) {
  const diubah = useRef(false)
  return (
    <div className="field">
      <label className="field__label" htmlFor={id}>
        {label}
        {wajib && <span className="field__req">*</span>}
      </label>
      <IsianAngka4
        id={id}
        value={value}
        onChange={(v) => {
          diubah.current = true
          onChange(v)
        }}
        onSelesai={() => {
          if (!diubah.current) return
          diubah.current = false
          onSelesai?.()
        }}
      />
    </div>
  )
}

function Tampil({ label, nilai }: { label: string; nilai: string }) {
  return (
    <div className="field">
      <span className="field__label">{label}</span>
      <div className="nbf-inward__teks">{nilai}</div>
    </div>
  )
}

/** Grid Warranty Extension `.AdditionalCoverage` (L7548): Coverage (autocomplete `BrowseMarineCondition_RD`) / Deductible (= ConditionName, baca-saja). */
function GridWarranty({ daftar, ubah }: { daftar: CoverageCargo['warrantyExtension']; ubah: (d: CoverageCargo['warrantyExtension']) => void }) {
  const ganti = (i: number, w: CoverageCargo['warrantyExtension'][number]) => ubah(daftar.map((x, j) => (j === i ? w : x)))
  return (
    <section className="nbf-alamat__blok nbf-aneka-form__kartu">
      <h5 className="nbf-objek__judul">{F.grupWarranty}</h5>
      <table className="nbf-tabel nbf-aneka-grid nbf-aneka-grid--aksi nbf-cargo-kapal">
        <thead>
          <tr>
            {F.kolomWarranty.map((k) => (
              <th key={k} scope="col">
                {k}
              </th>
            ))}
            <th scope="col">
              <button type="button" className="btn nbf-tema__tombol btn--ghost btn--sm" onClick={() => ubah([...daftar, { coverage: '', coverageNote: '' }])}>
                {F.tambah}
              </button>
            </th>
          </tr>
        </thead>
        <tbody>
          {daftar.map((w, i) => (
            <tr key={i}>
              <td>
                <IsianSaran
                  label={F.kolomWarranty[0]}
                  value={w.coverage}
                  onKetik={(v) => ganti(i, { coverage: v, coverageNote: '' })}
                  cari={(q) => cariKondisiMarine(q).then((h) => ({ baris: h.baris.map((b) => ({ id: b.id, label: b.id, ekstra: b.conditionName })) }))}
                  onPilih={(s) => ganti(i, { coverage: s.id, coverageNote: s.ekstra ?? '' })}
                  tampil={(s) => `${s.id} - ${s.ekstra ?? ''}`}
                  muatSaatFokus
                  wajibDariDaftar
                />
              </td>
              <td>{w.coverageNote}</td>
              <td>
                <button type="button" className="btn nbf-tema__tombol btn--ghost btn--sm" onClick={() => ubah(daftar.filter((_, j) => j !== i))}>
                  {F.hapus}
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </section>
  )
}

/** Form satu coverage cargo - `InputCoverageCargo_FacIn.xml`. */
export function FormCoverageCargo({
  c,
  idx,
  ubah,
  hitung,
  mataUang,
}: {
  c: CoverageCargo
  idx: string
  ubah: (c: CoverageCargo) => void
  hitung: (c: CoverageCargo, pemicu: PemicuCargo) => void
  mataUang: string[]
}) {
  const id = (m: string) => `nbfacin-covcargo-${idx}-${m}`
  const opsiUang: Opsi[] = mataUang.map((m) => ({ value: m, label: m }))
  const set = (k: 'rate' | 'limitOfLiability' | 'discountPercentage' | 'tsi' | 'premium' | 'minPremium') => (v: string) => ubah({ ...c, [k]: v })
  // Tata letak = tangkapan layar Pega (work owner 11-10-2026): Coverage Note, Coverage Initial (selebar penuh), lalu dua kolom -
  // kiri Currency / Rate (%) / Limit of Liability / Discount (%), kanan TSI / Premium / Min Premium / Discount.
  return (
    <div className="nbf-ringkas nbf-aneka-form">
      <section className="nbf-alamat__blok nbf-aneka-form__kartu">
        {/* Tombol Select Coverage (popup) DIGANTI isian saran yang bisa diketik (work owner 11-10-2026): pilih dari daftar
            `BrowseMarineCondition_RD` mengisi `.Coverage` + `.CoverageNote` (= SelectCoverageMarine_Act); ketikan yang tidak sama
            persis dengan satu saran dihapus saat kotak ditinggal (pola Coverage Note Aneka 09-10-2026). */}
        <IsianSaran
          label={F.coverageNote}
          value={c.coverageNote}
          muatSaatFokus
          wajibDariDaftar
          onKetik={(v) => ubah({ ...c, coverageNote: v, coverage: '' })}
          cari={(q) => cariKondisiMarine(q).then((h) => ({ baris: h.baris.map((b) => ({ id: b.id, label: b.conditionName })) }))}
          onPilih={(s) => ubah({ ...c, coverage: s.id, coverageNote: s.label })}
        />
        <div className="field">
          <label className="field__label" htmlFor={id('initial')}>
            {F.coverageInitial}
          </label>
          <input id={id('initial')} className="field__input" value={c.coverageInitial} onChange={(e) => ubah({ ...c, coverageInitial: e.target.value })} />
        </div>
        <div className="nbf-cargo-cov">
          <div className="nbf-cargo-cov__kolom">
            <Pilih
              label={F.currency}
              value={c.currency}
              onChange={(v) => {
                const baru = { ...c, currency: v }
                ubah(baru)
                hitung(baru, 'currency')
              }}
              opsi={opsiUang}
              kosong="—"
            />
            <Angka id={id('rate')} label={F.rate} value={c.rate} onChange={set('rate')} onSelesai={() => hitung(c, 'rate')} />
            <Angka id={id('lol')} label={F.limitOfLiability} value={c.limitOfLiability} onChange={set('limitOfLiability')} />
            <Angka id={id('diskon')} label={F.discountPercentage} value={c.discountPercentage} onChange={set('discountPercentage')} onSelesai={() => hitung(c, 'diskonPersen')} />
          </div>
          <div className="nbf-cargo-cov__kolom">
            <Angka id={id('tsi')} label={F.tsi} value={c.tsi} onChange={set('tsi')} onSelesai={() => hitung(c, 'tsi')} wajib />
            <Angka id={id('premium')} label={F.premium} value={c.premium} onChange={set('premium')} onSelesai={() => hitung(c, 'premium')} />
            <Angka id={id('minpremi')} label={F.minPremium} value={c.minPremium} onChange={set('minPremium')} />
            <Tampil label={F.discount} nilai={angka4(c.discount ?? '')} />
          </div>
        </div>
      </section>
      <GridWarranty daftar={c.warrantyExtension} ubah={(warrantyExtension) => ubah({ ...c, warrantyExtension })} />
      <GridDeductible daftar={c.deductibles} currency={c.currency} mataUang={mataUang} ubah={(deductibles) => ubah({ ...c, deductibles })} deskripsi={F.deskripsi} />
    </div>
  )
}

export default function TabCoverageCargo({ caseId }: { caseId: string }) {
  const [cargo, setCargo] = useState<ObjekCargo[] | null>(null)
  const [mataUang, setMataUang] = useState<string[]>([])
  const [halCargo, setHalCargo] = useState(1)
  const [buka, setBuka] = useState<number | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [proteksi, setProteksi] = useState<string[]>([])
  const [menyimpan, setMenyimpan] = useState(false)
  const [menghitung, setMenghitung] = useState(0)
  const [tersimpan, setTersimpan] = useState(false)
  const nomorHitung = useRef(new Map<string, number>())

  const rapikan = (d: ObjekCargo[]) => d.map((c) => ({ ...rapikanCargo(c), coverages: (c.coverages ?? []).map(rapikanCoverageCargo) }))

  useEffect(() => {
    let batal = false
    ambilObjekCargo(caseId).then(
      (h) => {
        if (!batal) setCargo(rapikan(h.cargo))
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

  if (cargo === null) return galat !== null ? <Gagal galat={galat} /> : <Memuat />

  const ganti = (f: (d: ObjekCargo[]) => ObjekCargo[]) => {
    setCargo((d) => (d === null ? d : f(d)))
    setTersimpan(false)
  }
  const ubahCoverage = (i: number, n: number, f: (c: CoverageCargo) => CoverageCargo) =>
    ganti((d) => d.map((x, j) => (j === i ? { ...x, coverages: (x.coverages ?? []).map((c, k) => (k === n ? f(c) : c)) } : x)))

  /** Rantai on-change satu coverage di backend; jawaban usang (ada hitung lebih baru untuk coverage yang sama) dibuang. */
  function hitung(i: number, n: number, c: CoverageCargo, pemicu: PemicuCargo) {
    const kunci = `${i}-${n}`
    const nomor = (nomorHitung.current.get(kunci) ?? 0) + 1
    nomorHitung.current.set(kunci, nomor)
    setMenghitung((x) => x + 1)
    setGalat(null)
    hitungCoverageCargo(caseId, [{ coverage: c, pemicu }])
      .then(
        (h) => {
          const hasil = h.baris[0]
          if (nomorHitung.current.get(kunci) !== nomor || hasil === undefined) return
          ubahCoverage(i, n, (kini) => terapkanHitungCargo(kini, hasil))
        },
        (err: unknown) => setGalat(err),
      )
      .finally(() => setMenghitung((x) => x - 1))
  }

  async function simpan(daftar: ObjekCargo[]) {
    setMenyimpan(true)
    setGalat(null)
    setProteksi([])
    setTersimpan(false)
    try {
      const h = await simpanObjekCargo(caseId, daftar)
      setCargo(rapikan(h.cargo))
      setTersimpan(true)
    } catch (err) {
      const p = pesanProteksiDari(err)
      if (p !== null) setProteksi(p)
      else setGalat(err)
    } finally {
      setMenyimpan(false)
    }
  }

  const kolom = GRID_CARGO_COVERAGE.length + 1

  return (
    <div className="nbf-objek nbf-obj-tab">
      <Gagal galat={galat} />
      {proteksi.length > 0 && (
        <div className="alert alert--error" role="alert">
          {TEKS_ANEKA.ditolak}
          <ul>
            {proteksi.map((p, i) => (
              <li key={i}>{p}</li>
            ))}
          </ul>
        </div>
      )}
      {tersimpan && <div className="alert alert--ok">{TEKS_ANEKA.tersimpan}</div>}
      {/* Cargo - 10 per halaman (work owner 11-10-2026) */}
      <PagerLokasi total={cargo.length} halaman={halCargo} onPindah={setHalCargo} />
      <div className="nbf-cov-wrap">
        <table className="nbf-tabel nbf-aneka-grid">
          <thead>
            <tr>
              <th scope="col" />
              {GRID_CARGO_COVERAGE.map((k) => (
                <th key={k} scope="col">
                  {k}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {cargo.length === 0 && (
              <tr>
                <td colSpan={kolom}>{TEKS_CARGO.kosong}</td>
              </tr>
            )}
            {potongHalaman(cargo, halCargo).map(([x, i]) => [
              <tr key={i} className={buka === i ? 'nbf-baris--terbuka' : undefined}>
                <td>
                  <button
                    type="button"
                    className={'btn btn--ghost btn--sm nbf-buka' + (buka === i ? ' nbf-buka--terbuka' : '')}
                    aria-label={TEKS_ANEKA.bukaBaris}
                    aria-expanded={buka === i}
                    onClick={() => setBuka(buka === i ? null : i)}
                  >
                    ▸
                  </button>
                </td>
                <td>{i + 1}</td>
                <td>{x.tradingNote}</td>
                <td>{x.goodNote === '' ? <span className="nbf-aneka-grid__kosong">{TEKS_ANEKA.belumDiisi}</span> : x.goodNote}</td>
                <td>{x.packingNote}</td>
                <td>{x.conveyanceNote}</td>
              </tr>,
              buka === i && (
                <tr key={`${i}-detail`} className="nbf-objek__detail">
                  <td colSpan={kolom}>
                    <div className="nbf-lapis nbf-lapis--objek">
                      <table className="nbf-tabel nbf-aneka-grid nbf-cargo-covgrid">
                        <thead>
                          <tr>
                            {GRID_COVERAGE_CARGO.kolom.map((k) => (
                              <th key={k} scope="col">
                                {k}
                              </th>
                            ))}
                          </tr>
                        </thead>
                        <tbody>
                          {(x.coverages ?? []).map((c, n) => [
                            <tr key={n}>
                              <td>{c.coverageNote === '' ? <span className="nbf-aneka-grid__kosong">{TEKS_ANEKA.belumDiisi}</span> : c.coverageNote}</td>
                              <td>{c.coverageInitial}</td>
                            </tr>,
                            <tr key={`${n}-form`} className="nbf-objek__detail">
                              <td colSpan={2}>
                                <FormCoverageCargo
                                  c={c}
                                  idx={`${i}-${n}`}
                                  mataUang={mataUang}
                                  ubah={(baru) => ubahCoverage(i, n, () => baru)}
                                  hitung={(baru, pemicu) => hitung(i, n, baru, pemicu)}
                                />
                              </td>
                            </tr>,
                          ])}
                        </tbody>
                      </table>
                      {(x.totalGross ?? []).length > 0 && (
                        <table className="nbf-tabel nbf-aneka-grid nbf-aneka-total">
                          <thead>
                            <tr>
                              {GRID_COVERAGE_CARGO.total.map((k) => (
                                <th key={k} scope="col">
                                  {k}
                                </th>
                              ))}
                            </tr>
                          </thead>
                          <tbody>
                            {(x.totalGross ?? []).map((t, k) => (
                              <tr key={k}>
                                <td>{t.currency}</td>
                                <td className="nbf-angka">{angka4(t.tsi)}</td>
                                <td className="nbf-angka">{angka4(t.premium)}</td>
                                <td className="nbf-angka">{angka4(t.rate)}</td>
                              </tr>
                            ))}
                          </tbody>
                        </table>
                      )}
                    </div>
                  </td>
                </tr>
              ),
            ])}
          </tbody>
        </table>
      </div>
      <div className="nbf-objek__kaki">
        <button
          type="button"
          className="btn nbf-tema__tombol nbf-tema__tombol--utama btn--primary"
          onClick={() => void simpan(cargo)}
          disabled={menyimpan || menghitung > 0}
        >
          {menyimpan ? TEKS_FORM_OPPORTUNITY.menyimpan : SIMPAN_OBJEK.label}
        </button>
      </div>
    </div>
  )
}
