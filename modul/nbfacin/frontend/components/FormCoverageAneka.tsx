// Form satu coverage Aneka (tiket 73 P2) - port `NB FacIn\Section\InputCoverageAneka_FacIn.xml` (Data-Coverage, 20.118 baris).
//
// Tata letak = layout Pega (work owner 09-10-2026, tangkapan layar Pega "SAMAIN KAYA GAMBAR2"): SATU kisi 3 kolom, sel mengalir
// berurutan (sel tersembunyi tidak menyisakan tempat), Sub Limit selebar penuh, lalu garis + Deductible.
//
// Baca-saja `[terverifikasi]` (`pyEditOptions=Read-only` + `pyReadOnlyCondition`): syarat `IsClaim` = dapat diisi di NB;
// TANPA syarat = SELALU baca-saja (tangkapan layar Pega: TSI tampil sebagai teks): Currency L3557, TSI L4069, % Indemnity L6149.
// Coverage Note L1709 juga tanpa syarat (di Pega diisi tombol Choose Coverage), tetapi keputusan work owner 09-10-2026 "ini
// balikin aja, sama seperti sblm nya": Coverage Note = isian dengan saran sambil mengetik (C-3), tanpa tombol / popup.
//
// Setiap sel yang di Pega punya aksi on-change memanggil `hitung(coverage, pemicu)` saat isian SELESAI (lepas fokus / Enter /
// pilih) - rantai aksinya dikerjakan backend `POST …/hitung-coverage-aneka` (nama pemicu = kabel sesi c3):
// Calculate Method L2364 'metode' · % Pro Rate L2649 'proRate' · % Short Period L2941 'shortPeriod' · Day L3258 'day' ·
// % Loading L4501 'loading' · % Scale L4793 'scale' · % Rate L5168 'rate' · Indemnity / Unit L5548 / L5843 'indemnity' ·
// Premium L6339 'premium' · % Discount L6609 'diskonPersen' · Discount L6981 'diskonJumlah' · Limit of Liability L7649 'lol' ·
// % LoL L7922 'pctLol'. (Currency / TSI baca-saja: pemicunya tetap dipakai Copy Coverage From This Section.)
//
// Keputusan agent (tiket 73 P2):
// - C-1 % Pro Rate terikat `pyWorkPage.OfferFacIn.ProRatePercent` (nilai CASE, L2649/L2713): isiannya mengubah nilai case;
//   `fillPremiBond` membaca `.ProRatePercent` COVERAGE (persis Pega) - coverage tidak disalin dari ketikan.
// - C-3 Tombol "Choose Coverage" (L1221, local action `ChooseCov_Aneka` - section tidak ada di ekspor, K73-8) TIDAK dibuat:
//   Coverage Note = saran sambil mengetik dari RD `BrowseCoverageFacIn_RD` (Type ANEKA, BizCode case), pilih `.ID` → coverage,
//   `.NamaCoverage` → Coverage Note (work owner 09-10-2026).

import { useRef } from 'react'

import { Area, Pilih, type Opsi } from '../../../../inti/frontend/components/ui/dasar'
import { angka4 } from '../angka4'
import { cariCoverageAneka, type BarisCoverage, type CoverageAneka, type Deductible, type PemicuAneka, type TampilanAneka } from '../api'
import { FORM_COVERAGE_ANEKA as F, OPSI_DAY_ANEKA, OPSI_METODE_ANEKA, OPSI_UNIT_INDEMNITY } from '../labels'
import GridDeductible from './GridDeductible'
import IsianAngka4 from './IsianAngka4'
import IsianSaran from './IsianSaran'

const OPSI_METODE: Opsi[] = OPSI_METODE_ANEKA.map((o) => ({ value: o.value, label: o.label }))
const OPSI_UNIT: Opsi[] = OPSI_UNIT_INDEMNITY.map((o) => ({ value: o.value, label: o.label }))

/** Isian angka berlabel; `onSelesai` hanya bila nilainya diubah sejak fokus (event change Pega). */
function Angka({ id, label, value, onChange, onSelesai }: { id: string; label: string; value: string; onChange: (v: string) => void; onSelesai?: () => void }) {
  const diubah = useRef(false)
  return (
    <div className="field">
      <label className="field__label" htmlFor={id}>
        {label}
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

/** Sel baca-saja (label + teks); angka rata kanan seperti Pega. */
function Tampil({ label, nilai, angka }: { label: string; nilai: string; angka?: boolean }) {
  return (
    <div className="field">
      <span className="field__label">{label}</span>
      <div className={'nbf-inward__teks' + (angka ? ' nbf-angka' : '')}>{nilai === '' ? '—' : nilai}</div>
    </div>
  )
}

export default function FormCoverageAneka({
  caseId,
  idDasar,
  c,
  ubah,
  hitung,
  tampilan,
  mataUang,
  proRate,
  ubahProRate,
  ubahDeductible,
}: {
  caseId: string
  /** Awalan id unik isian (satu form per coverage di layar). */
  idDasar: string
  c: CoverageAneka
  ubah: (c: CoverageAneka) => void
  /** Jalankan rantai on-change di backend untuk coverage ini (hasilnya dipasang pemanggil). */
  hitung: (c: CoverageAneka, pemicu: PemicuAneka) => void
  tampilan: TampilanAneka
  mataUang: string[]
  /** % Pro Rate tingkat case. */
  proRate: string
  ubahProRate: (v: string) => void
  /** DeductibleList berubah (pemanggil menerapkan addDeductibleAneka_ACT bila bertambah). */
  ubahDeductible: (d: Deductible[]) => void
}) {
  const coverageCari = useRef<BarisCoverage[]>([])
  const set = (k: keyof CoverageAneka) => (v: string) => ubah({ ...c, [k]: v })
  const angka = (k: keyof CoverageAneka, label: string, pemicu?: PemicuAneka) => (
    <Angka
      id={`${idDasar}-${k}`}
      label={label}
      value={String(c[k] ?? '')}
      onChange={set(k)}
      onSelesai={pemicu === undefined ? undefined : () => hitung(c, pemicu)}
    />
  )
  const pilihLalu = (k: keyof CoverageAneka, pemicu: PemicuAneka) => (v: string) => {
    const baru = { ...c, [k]: v }
    ubah(baru)
    hitung(baru, pemicu)
  }
  const metode = c.calculateMethodFacIn

  return (
    <div className="nbf-ringkas nbf-aneka-cov">
      <div className="nbf-aneka-cov__grid">
        {!tampilan.tanamanHasil && (
          // L1709 Coverage Note (visible `!IsYieldShortfall`) - saran sambil mengetik (C-3).
          <IsianSaran
            label={F.coverageNote}
            value={c.coverageNote}
            muatSaatFokus
            // Work owner 09-10-2026: "jika yang diketik tidak ada dari daftar coverage, dihapus saja ketikannya" - teks yang tidak
            // sama persis dengan satu saran dikosongkan saat kotak ditinggal; mengetik melepas kode coverage pilihan sebelumnya.
            wajibDariDaftar
            onKetik={(v) => ubah({ ...c, coverageNote: v, coverage: '', oldId: '' })}
            cari={(q) =>
              cariCoverageAneka(caseId, q).then((h) => {
                coverageCari.current = h.baris
                return { baris: h.baris.map((b) => ({ id: b.id, label: b.nama })) }
              })
            }
            onPilih={(s) => {
              const b = coverageCari.current.find((x) => x.id === s.id)
              ubah({ ...c, coverage: s.id, coverageNote: s.label, oldId: b?.oldId ?? c.oldId })
            }}
          />
        )}
        <Pilih label={F.calculateMethod} value={metode} onChange={pilihLalu('calculateMethodFacIn', 'metode')} opsi={OPSI_METODE} />
        {metode === '1' && (
          // L2649 `pyWorkPage.OfferFacIn.ProRatePercent` (case, C-1).
          <Angka id={`${idDasar}-proRate`} label={F.proRate} value={proRate} onChange={ubahProRate} onSelesai={() => hitung(c, 'proRate')} />
        )}
        {metode === '2' && angka('pctShortPeriod', F.shortPeriod, 'shortPeriod')}
        <fieldset className="nbf-aneka-radio">
          <legend className="field__label">{F.day}</legend>
          {OPSI_DAY_ANEKA.map((d) => (
            <label key={d} className="nbf-inward__pilihan">
              <input type="radio" name={`${idDasar}-day`} value={d} checked={c.day === d} onChange={() => pilihLalu('day', 'day')(d)} /> {d}
            </label>
          ))}
        </fieldset>
        {tampilan.currencyCoverage && <Tampil label={F.currency} nilai={c.currency} />}
        <Tampil label={F.tsi} nilai={angka4(c.tsi)} angka />
        {angka('loading', F.loading, 'loading')}
        {angka('lostLimit', F.scale, 'scale')}
        {angka('rate', F.rate, 'rate')}
        {tampilan.mbd && angka('indemnity', F.indemnity, 'indemnity')}
        {tampilan.mbd && <Pilih label={F.unit} value={c.unit} onChange={pilihLalu('unit', 'indemnity')} opsi={OPSI_UNIT} />}
        {tampilan.mbd && <Tampil label={F.pctIndemnity} nilai={angka4(c.indemnityPercentage)} angka />}
        {angka('premium', F.premium, 'premium')}
        {angka('discountPercentage', F.pctDiscount, 'diskonPersen')}
        {angka('discount', F.discount, 'diskonJumlah')}
        {angka('limitOfLiability', F.lol, 'lol')}
        {angka('pctLol', F.pctLol, 'pctLol')}
        <div className="nbf-aneka-cov__lebar">
          <Area label={F.subLimit} value={c.subLimitNote} onChange={set('subLimitNote')} baris={4} />
        </div>
      </div>

      <hr className="nbf-aneka-cov__garis" />
      <GridDeductible daftar={c.deductibles} currency={c.currency} mataUang={mataUang} ubah={ubahDeductible} />

    </div>
  )
}
