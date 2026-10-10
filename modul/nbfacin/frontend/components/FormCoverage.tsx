// Form satu coverage FIRE (tiket 43, tahap C1) - port `NB FacIn\Section\CoverageItem.xml` (flow action
// `CoverageItem`) dan popup Choose Coverage (local action `ChooseCoverage`, RD `BrowseCoverageFacIn_RD`).
//
// Premi TIDAK dihitung di layar: setiap perubahan medan hitung memanggil `POST …/hitung-coverage` (rumus
// `CountPremi_ACT`, desimal eksak di backend) - mode "percent" (premi dari rate) kecuali Gross Premium yang diubah
// (mode "amount", rate dihitung balik). Medan tampil per Coverage Basis mengikuti syarat sel `CoverageItem.xml`.
//
// Tata letak (tiket 46, ralat 03-10-2026 atas gambar layar Pega): blok atas berlabel kiri (sel 3-21: Coverage Basis,
// Choose Coverage, Coverage = `.OLDID`, Accumulation Code / Address, Choose Accumulation Code + Copy Accumulation,
// Conditions), lalu dua kolom berlabel kiri (sel 24: Days radio … % Limit of Liability; sel 36: % Pro Rate … Gross
// Premium). Sel ber-`ALWAYS` dengan syarat sisa tetap tampil (Indemnity 28, Indemnity Unit 40, % Loss Limit 44);
// sel `pyReadOnly` tanpa syarat = teks (Coverage, Accumulation, TSI, % Indemnity 41, % First Scale 42).
//
// Keputusan agent (tiket 43): P-1 Layering (basis 5), Zone / 4.2 Construction, View Indemnity Table (sel 50,
// `ViewIndemnity_LA`), perubahan Indemnity Unit (`SetIndemnityRate_ACT`) = tahap berikut. P-2 (ralat) Days = radio
// 365 / 366 / 360 menurut gambar layar (aturan `.Day` tidak ada di korpus). P-3 tanda wajib ‰ Gross Rate tanpa menahan
// Save (pola L-2). P-4 jeda hitung 500 ms; jawaban lama dibuang.

import { useEffect, useId, useRef, useState } from 'react'

import { Area, Field, Gagal, Kosong, Memuat, Modal, Pilih } from '../../../../inti/frontend/components/ui/dasar'
import { desimalSah } from '../../../../inti/frontend/lib/desimal'
import { angka4 } from '../angka4'
import {
  cariCoverage,
  daftarMataUang,
  hitungCoverage,
  saranKonstruksiGempa,
  saranZona,
  type BarisCoverage,
  type CoverageObjek,
  type ModeHitung,
} from '../api'
import GridDeductible from './GridDeductible'
import { dariKetikan } from './IsianAngka4'
import IsianSaran from './IsianSaran'
import IsianUang from './IsianUang'
import {
  AKUMULASI_KOSONG,
  COVERAGE_KONSTRUKSI,
  COVERAGE_ZONA,
  FORM_COV as F,
  LABEL_COVERAGE_BASIS,
  OPSI_COVERAGE_BASIS,
  OPSI_INDEMNITY_UNIT,
  PILIHAN_DAY_COVERAGE,
  POPUP_OKUPASI,
  TEKS_COVERAGE,
  ZONA_COVERAGE,
} from '../labels'
import PopupAkumulasi from './PopupAkumulasi'

const JEDA_MS = 500

/**
 * Nilai awal coverage (work owner 08-10-2026, tangkapan layar Pega: "ini nilai default saat isi coverage, dan jika tidak
 * di isi, dikasih nilainya, jangan null"): Indemnity 0, ‰ Gross Rate / % Discount / ‰ Net Rate / Limit of Liability /
 * % Limit of Liability / Discount / Gross Premium 0, % Loss Limit 100, Indemnity Unit "0" (Day (s)). Asal di Pega tidak ada
 * di korpus (`CountPremi_ACT` hanya memakai LossLimit 100 sebagai nilai lokal). Backend mengisi nilai yang sama untuk medan
 * kosong saat simpan.
 * Medan khusus basis (work owner 08-10-2026 "perlu juga, ikutin seperti xml dari pega"): Sub Limit 100 = `CountPremi_ACT`
 * L3638 `Local.SubLimit = 100` bila tidak > 0 (pola sama dengan LossLimit; 0 akan memicu `seterrormessagesublimit_ACT`
 * `.Sublimit < 5`); First Loss / EML-PML 0 = XML tanpa default, rumus `TSILiability = TSI * FirstLoss|EmlPml / 100` memperlakukan
 * kosong sebagai 0. First Scale TIDAK diberi nilai awal: dihitung dari tabel First Loss Scale (`SetFirstLossScale_ACT` L1073).
 */
export const NILAI_AWAL_COVERAGE = {
  indemnity: '0',
  rate: '0',
  discountPercentage: '0',
  netRate: '0',
  limitOfLiability: '0',
  pctLol: '0',
  lostLimit: '100',
  discount: '0',
  premium: '0',
  unit: '0',
  firstLoss: '0',
  emlPml: '0',
  sublimit: '100',
} as const

/**
 * Coverage dari server (data lama / hasil hitung) dengan medan kosong -> `NILAI_AWAL_COVERAGE` saat ditampilkan (work owner
 * 08-10-2026 "ini yang tidak ada isinya di ksh nilai 0 jgn ksong gitu"). Pasangan Rate / Premium dan % Discount / Discount
 * hanya diisi bila KEDUANYA kosong: Rate 0 di samping premi terisi membuat Save memilih mode percent dan menimpa premi
 * (backend `hitungSimpan` A158); % Discount 0 di samping diskon nominal sama halnya untuk hitung ulang diskon.
 */
export function isiNilaiAwal(c: CoverageObjek): CoverageObjek {
  const hasil: CoverageObjek = { ...c }
  const kosong = (v: string | undefined) => (v ?? '') === ''
  const pasangan: Partial<Record<keyof CoverageObjek, boolean>> = {
    rate: kosong(c.rate) && kosong(c.premium),
    premium: kosong(c.rate) && kosong(c.premium),
    discountPercentage: kosong(c.discountPercentage) && kosong(c.discount),
    discount: kosong(c.discountPercentage) && kosong(c.discount),
  }
  for (const [k, v] of Object.entries(NILAI_AWAL_COVERAGE) as [keyof typeof NILAI_AWAL_COVERAGE, string][]) {
    if (pasangan[k] === false) continue
    if ((c[k] ?? '') === '') hasil[k] = v
  }
  return hasil
}

/**
 * Coverage kosong: basis Sum Insured, % Indemnity 100 (`AddCoverageAutoFire`); Indemnity Unit "0" = pilihan pertama yang
 * dikirim dropdown tanpa pilihan kosong (K-6); medan angka lain = `NILAI_AWAL_COVERAGE`.
 */
export function coverageBaru(b?: BarisCoverage): CoverageObjek {
  return {
    coverage: b?.id ?? '', oldId: b?.oldId ?? '', coverageNote: b?.nama ?? '', coverageBasis: '1', day: '', tsi: '',
    rateOjk: '', tsiLiability: '', proRatePercent: '', indemnityPercentage: '100', firstScale: '', conditions: '',
    ...NILAI_AWAL_COVERAGE,
  }
}

/** Medan angka yang diisi pengguna. */
const MEDAN_ANGKA = [
  'indemnity', 'rate', 'firstLoss', 'discountPercentage', 'netRate', 'limitOfLiability', 'pctLol', 'indemnityPercentage',
  'firstScale', 'sublimit', 'lostLimit', 'emlPml', 'discount', 'premium',
] as const

/** Angka sah: kosong atau desimal bertitik. */
export const angkaSah = (v: string) => v.trim() === '' || desimalSah(v.trim())

/** Ada coverage dengan isian angka tidak sah. */
export function adaGalatCoverage(c: CoverageObjek[]): boolean {
  return c.some((x) => MEDAN_ANGKA.some((k) => !angkaSah(x[k])))
}

function Tampil({ label, nilai, angka }: { label: string; nilai: string; angka?: boolean }) {
  return (
    <div className="field">
      <span className="field__label">{label}</span>
      <div className={'nbf-inward__teks' + (angka ? ' nbf-angka' : '')}>{nilai}</div>
    </div>
  )
}

function PopupCoverage({ onTutup, onPilih }: { onTutup: () => void; onPilih: (b: BarisCoverage) => void }) {
  const [kotak, setKotak] = useState('')
  const [hasil, setHasil] = useState<BarisCoverage[] | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const nomor = useRef(0)

  useEffect(() => {
    const n = ++nomor.current
    const jadwal = window.setTimeout(
      () => {
        cariCoverage(kotak).then(
          (h) => {
            if (n === nomor.current) setHasil(h.baris)
          },
          (err: unknown) => {
            if (n === nomor.current) setGalat(err)
          },
        )
      },
      kotak === '' ? 0 : 400,
    )
    return () => window.clearTimeout(jadwal)
  }, [kotak])

  return (
    <Modal judul={F.pilihCoverage.label} onTutup={onTutup} lebar>
      <div className="nbf-popup__kotak">
        <Field label={TEKS_COVERAGE.cariCoverage} value={kotak} onChange={setKotak} />
      </div>
      <div className="table-wrap">
        <table className="nbf-tabel">
          <thead>
            <tr>
              <th scope="col">{F.coverage.label}</th>
              <th scope="col">{POPUP_OKUPASI.kolom[1].label}</th>
              <th scope="col" />
            </tr>
          </thead>
          {hasil && hasil.length > 0 && (
            <tbody>
              {hasil.map((b) => (
                <tr key={b.id}>
                  <td>{b.oldId}</td>
                  <td>{b.nama}</td>
                  <td className="table__actions">
                    <button type="button" className="btn nbf-tema__tombol btn--ghost btn--sm" onClick={() => onPilih(b)}>
                      {POPUP_OKUPASI.pilih.label}
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          )}
        </table>
      </div>
      {hasil === null && galat === null && <Memuat />}
      {hasil !== null && hasil.length === 0 && <Kosong pesan={TEKS_COVERAGE.tanpaCoverage} />}
      <Gagal galat={galat} />
    </Modal>
  )
}

export default function FormCoverage({
  caseId,
  c: cAsli,
  tsiItem,
  item,
  pertama,
  onSalinAkumulasi,
  zipRisiko = '',
  lokasi,
  ubah,
}: {
  caseId: string
  c: CoverageObjek
  /** TSI Object Item (teks desimal). */
  tsiItem: string
  /** Item pemilik (PctAdjustment `CountPremi_ACT`; mata uang = awal deductible baru). */
  item?: { isAdjustable: boolean; pctAdjustOther: string; currency?: string }
  /** Coverage pertama item - tombol Copy Accumulation (sel 19, `.pxListSubscript==1`). */
  pertama?: boolean
  /** Copy Accumulation: salin ke semua coverage di objek yang sama (`CopyAccumulationCode_Act`). */
  onSalinAkumulasi?: () => void
  /** Zip lokasi risiko objek (`RiskLocation.ASMZipCode`) - saringan awal & pemeriksaan Choose Accumulation. */
  zipRisiko?: string
  /** Data objek untuk `CountPremiRate_ACT` (callAct "rate", tiket 52): `.Property.Province` dan OccupationId okupasinya. */
  lokasi?: { province: string; occupationIds: string[] }
  ubah: (c: CoverageObjek) => void
}) {
  // Medan kosong tampil (dan ikut terkirim saat diubah) dengan nilai awal; backend mengisi nilai yang sama saat Save.
  const c = isiNilaiAwal(cAsli)
  const [pilih, setPilih] = useState(false)
  const [pilihAkumulasi, setPilihAkumulasi] = useState(false)
  // Nama grup radio Days unik per form (beberapa coverage dapat terbuka bersamaan).
  const namaHari = `nbfacin-cov-day-${useId()}`
  const [menghitung, setMenghitung] = useState(false)
  const [galat, setGalat] = useState<unknown>(null)
  const nomor = useRef(0)
  const jadwal = useRef<number | undefined>(undefined)
  // `Param.DiscountStatus`: medan diskon terakhir yang diubah pengguna.
  const modeDiskon = useRef<ModeHitung>('percent')

  useEffect(() => () => window.clearTimeout(jadwal.current), [])

  // Daftar mata uang untuk Currency deductible (RD `BrowseCurrency_RD`).
  const [mataUang, setMataUang] = useState<string[]>([])
  useEffect(() => {
    let batal = false
    daftarMataUang().then(
      (h) => {
        if (!batal) setMataUang(h.baris)
      },
      () => {},
    )
    return () => {
      batal = true
    }
  }, [])

  /** Ubah lalu jadwalkan hitung ulang di backend; `callAct="rate"` = Coverage / Zone diubah (tiket 52). */
  function ubahHitung(baru: CoverageObjek, mode: ModeHitung, callAct?: 'rate') {
    ubah(baru)
    window.clearTimeout(jadwal.current)
    if (MEDAN_ANGKA.some((k) => !angkaSah(baru[k])) || baru.coverageBasis === '5') return
    jadwal.current = window.setTimeout(() => {
      const n = ++nomor.current
      setMenghitung(true)
      setGalat(null)
      hitungCoverage(caseId, {
        coverage: baru,
        tsi: tsiItem,
        mode,
        modeDiskon: modeDiskon.current,
        isAdjustable: item?.isAdjustable,
        pctAdjustOther: item?.pctAdjustOther,
        ...(callAct ? { callAct, province: lokasi?.province ?? '', occupationIds: lokasi?.occupationIds ?? [] } : {}),
      }).then(
        (h) => {
          if (n === nomor.current) {
            ubah(h)
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

  // Mode amount (Gross Premium diketik): Rate dikosongkan sampai hitung mengisinya - Save sebelum hitung selesai tetap
  // amount di backend (hitungSimpan A158: rate kosong + premi terisi), tidak tertimpa Rate lama / nilai awal 0 (tiket 68).
  const set = (k: keyof CoverageObjek, mode: ModeHitung = 'percent') => (v: string) =>
    ubahHitung(
      { ...c, [k]: (MEDAN_ANGKA as readonly string[]).includes(k) ? dariKetikan(v) : v, ...(mode === 'amount' ? { rate: '' } : {}) },
      mode,
    )
  const angka = (k: (typeof MEDAN_ANGKA)[number], label: string, mode: ModeHitung = 'percent', required = false) => (
    <Field label={label} value={c[k]} onChange={set(k, mode)} required={required} error={angkaSah(c[k]) ? undefined : TEKS_COVERAGE.angka} />
  )
  /** Medan UANG: berformat ribuan saat diketik (IsianUang, permintaan work owner 03-10-2026). */
  const uang = (k: 'limitOfLiability' | 'premium', label: string, mode: ModeHitung = 'percent') => (
    <IsianUang label={label} value={c[k]} onChange={set(k, mode)} error={angkaSah(c[k]) ? undefined : TEKS_COVERAGE.angka} />
  )
  /** Medan diskon: menandai mode diskon (`DiscountStatus`) sebelum hitung. */
  const angkaDiskon = (k: 'discountPercentage' | 'discount', label: string, md: ModeHitung) => {
    const props = {
      label,
      value: c[k],
      onChange: (v: string) => {
        modeDiskon.current = md
        ubahHitung({ ...c, [k]: dariKetikan(v) }, 'percent')
      },
      error: angkaSah(c[k]) ? undefined : TEKS_COVERAGE.angka,
    }
    // Discount = uang (berformat ribuan); % Discount = persen (isian biasa).
    return k === 'discount' ? <IsianUang {...props} /> : <Field {...props} />
  }
  const b = c.coverageBasis

  return (
    <div className="nbf-objek__isi nbf-cov">
      {/* Blok atas (sel 1, "Stacked with labels left"): sel 3 -> 21. */}
      <div className="nbf-cov__atas nbf-labelkiri">
        <Pilih label={LABEL_COVERAGE_BASIS} value={b} onChange={set('coverageBasis')} opsi={OPSI_COVERAGE_BASIS} />
        <div className="nbf-opp__tombol">
          <button type="button" className="btn nbf-tema__tombol btn--sm" onClick={() => setPilih(true)}>
            {F.pilihCoverage.label}
          </button>
        </div>
        {/* Sel 5 = `.OLDID` baca-saja. */}
        <Tampil label={F.coverage.label} nilai={c.oldId} />
        {/* Sel 6 / 11 / 12: Zone - autocomplete selagi kosong; terisi -> teks + ikon "Switch to edit mode" (SetShowZone_DT). */}
        {COVERAGE_ZONA.includes(c.coverage) &&
          ((c.zone ?? '') === '' ? (
            <IsianSaran
              label={ZONA_COVERAGE.zone}
              value=""
              onKetik={() => {}}
              cari={(q) => saranZona(c.coverage, c.oldId, q)}
              muatSaatFokus
              onPilih={(s) => ubahHitung({ ...c, zone: s.label }, 'percent', 'rate')}
            />
          ) : (
            <div className="field">
              <span className="field__label">{ZONA_COVERAGE.zone}</span>
              <div className="nbf-cov__zona">
                <div className="nbf-inward__teks">{c.zone}</div>
                <button
                  type="button"
                  className="btn nbf-tema__tombol btn--ghost btn--sm"
                  title={ZONA_COVERAGE.ubahZona}
                  aria-label={ZONA_COVERAGE.ubahZona}
                  onClick={() => ubah({ ...c, zone: '', constructionEarthquake: '' })}
                >
                  ✎
                </button>
              </div>
            </div>
          ))}
        {/* Sel 13: 4.2 Construction - mengubahnya menghitung ulang premi (tarif EQS db-link tidak dibaca; tiket 52 butir 2). */}
        {COVERAGE_KONSTRUKSI.includes(c.coverage) && (
          <IsianSaran
            label={ZONA_COVERAGE.konstruksi}
            value={c.constructionEarthquake ?? ''}
            onKetik={(v) => ubah({ ...c, constructionEarthquake: v })}
            cari={(q) => saranKonstruksiGempa(q)}
            muatSaatFokus
            onPilih={(s) => ubahHitung({ ...c, constructionEarthquake: s.label }, 'percent')}
          />
        )}
        <Tampil label={F.accumulationCode.label} nilai={c.accumulationCode || AKUMULASI_KOSONG} />
        <Tampil label={F.accumulationAddress.label} nilai={c.accumulationDescription || AKUMULASI_KOSONG} />
        {/* data-ubah-akumulasi: tetap aktif di tahap Underwriting (tiket 60, `CoverageItem_IsUW` kontainer tombol L7247). */}
        <div className="nbf-opp__tombol">
          <button type="button" className="btn nbf-tema__tombol btn--sm" data-ubah-akumulasi="" onClick={() => setPilihAkumulasi(true)}>
            {F.pilihAkumulasi.label}
          </button>
          {/* Sel 19: hanya coverage pertama (`.pxListSubscript==1`). */}
          {pertama && onSalinAkumulasi && (
            <button type="button" className="btn nbf-tema__tombol btn--sm" data-ubah-akumulasi="" onClick={onSalinAkumulasi}>
              {F.salinAkumulasi.label}
            </button>
          )}
        </div>
        <Area label={F.conditions.label} value={c.conditions} onChange={(v) => ubah({ ...c, conditions: v })} baris={3} />
      </div>
      {b === '5' && <div className="alert alert--warn">{TEKS_COVERAGE.layeringBelum}</div>}
      {/* Dua kolom (sel 22, "Inline grid double"): kiri sel 24 (26-35), kanan sel 36 (38-47). */}
      <div className="nbf-cov__kolom">
        <div className="nbf-cov__tumpuk nbf-labelkiri">
          <div className="field">
            <span className="field__label">{F.day.label}</span>
            <div className="nbf-inward__radio" role="radiogroup" aria-label={F.day.label}>
              {PILIHAN_DAY_COVERAGE.map((d) => (
                <label key={d} className="nbf-inward__pilihan">
                  <input type="radio" name={namaHari} value={d} checked={c.day === d} onChange={() => set('day')(d)} /> {d}
                </label>
              ))}
            </div>
          </div>
          <Tampil label={F.tsi.label} nilai={angka4(c.tsi || tsiItem)} angka />
          {angka('indemnity', F.indemnity.label)}
          {angka('rate', F.rate.label, 'percent', true)}
          {b === '2' && angka('firstLoss', F.firstLoss.label)}
          {angkaDiskon('discountPercentage', F.discountPercentage.label, 'percent')}
          {b === '2' && <Tampil label={F.tsiLiability.label} nilai={angka4(c.tsiLiability)} angka />}
          {angka('netRate', F.netRate.label)}
          {uang('limitOfLiability', F.limitOfLiability.label)}
          {angka('pctLol', F.pctLol.label)}
        </div>
        <div className="nbf-cov__tumpuk nbf-labelkiri">
          <Tampil label={F.proRate.label} nilai={angka4(c.proRatePercent)} angka />
          {/* Sel 40: tanpa pilihan kosong (`pyHasNoSelection` false) - kosong tampil sebagai pilihan pertama "0". */}
          <div className="field">
            <label className="field__label" htmlFor={`${namaHari}-unit`}>
              {F.unit.label}
            </label>
            <select
              id={`${namaHari}-unit`}
              className="field__input"
              value={c.unit || OPSI_INDEMNITY_UNIT[0]!.value}
              onChange={(e) => ubah({ ...c, unit: e.target.value })}
            >
              {OPSI_INDEMNITY_UNIT.map((o) => (
                <option key={o.value} value={o.value}>
                  {o.label}
                </option>
              ))}
            </select>
          </div>
          {/* Sel 41 / 42: baca-saja tanpa syarat. */}
          <Tampil label={F.indemnityPercentage.label} nilai={angka4(c.indemnityPercentage)} angka />
          {b === '2' && <Tampil label={F.firstScale.label} nilai={angka4(c.firstScale)} angka />}
          {(b === '4' || b === '5') && angka('sublimit', F.sublimit.label)}
          {angka('lostLimit', F.lostLimit.label)}
          {b === '3' && angka('emlPml', F.emlPml.label)}
          {angkaDiskon('discount', F.discount.label, 'amount')}
          {uang('premium', F.premium.label, 'amount')}
          {menghitung && <span className="muted">{TEKS_COVERAGE.menghitung}</span>}
          <Gagal galat={galat} />
        </div>
      </div>
      {/* Deductible (tahap C3, tiket 45): tidak memicu hitung premi. */}
      <GridDeductible
        daftar={c.deductibles ?? []}
        currency={item?.currency ?? ''}
        mataUang={mataUang}
        ubah={(deductibles) => ubah({ ...c, deductibles })}
      />
      {pilihAkumulasi && (
        <div data-ubah-akumulasi="">
          <PopupAkumulasi
            zipRisiko={zipRisiko}
            onTutup={() => setPilihAkumulasi(false)}
            onPilih={(kode, alamat) => {
              setPilihAkumulasi(false)
              ubah({ ...c, accumulationCode: kode, accumulationDescription: alamat })
            }}
          />
        </div>
      )}
      {pilih && (
        <PopupCoverage
          onTutup={() => setPilih(false)}
          onPilih={(row) => {
            setPilih(false)
            // Coverage diubah: SetShowZone_DT mengosongkan Zone / 4.2 Construction, lalu CountPremi_ACT(CallAct="rate").
            ubahHitung(
              { ...c, coverage: row.id, oldId: row.oldId, coverageNote: row.nama, zone: '', constructionEarthquake: '' },
              'percent',
              'rate',
            )
          }}
        />
      )}
    </div>
  )
}
