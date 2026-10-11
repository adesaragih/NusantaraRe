// Isi tab polis NB Fac In hanya-baca Marine Cargo (`PolisNB.tsx`; perintah work owner 11-10-2026 "tampilan view polis
// yang belum tergenerate tolong diperbaiki. cek tampilan nya di nb facin!"): Object (grid cargo -> form cargo: kepala
// Trading / Conveyance / Goods / Packing, Survey Agent, LC, Ship Data, data pelayaran, Additional Ship), Coverage (cargo ->
// coverage + form selalu terbuka, Warranty Extension, Deductible ber-Deskripsi, Grand Total), Spreading (cargo -> coverage
// -> detail + baris spreading, ringkasan `totalCargo`). Susunan dan judul disalin dari `modul/nbfacin/frontend/components/
// {TabObjectCargo,TabCoverageCargo,TabSpreadingCargo}.tsx` (pola, bukan impor); tombol (Add / Delete / Select Agent /
// Select Ship / Select Coverage / Save / Copy) tidak dibawa.

import type { ReactNode } from 'react'

import type { Baris, Halaman } from '../api'
import {
  JUDUL_DEDUCTIBLE_NB,
  SPREADING_NB,
  TEKS_NB,
  angkaNB,
  anakNB,
  benarNB,
  daftarNB,
  nilaiNB,
  tanggalNB,
} from '../polisNB'
import {
  CARGO_NB,
  FORM_CARGO_NB,
  FORM_COVERAGE_CARGO_NB,
  GRID_CARGO_COVERAGE_NB,
  GRID_CARGO_NB,
  GRID_COVERAGE_CARGO_NB,
  SPREADING_CARGO_NB,
  TEKS_CARGO_NB,
  coverageCargoNB,
  hitungSpreadingCargoNB,
  kapalTambahanNB,
  medanCargo,
  totalGrossCargoNB,
} from '../polisNBCargo'
import { CentangNB, GridNB, TeksNB } from './PolisNBDasar'
import { BarisSpreadingNB, DeductibleNB, DetailCoverageSpreadingNB, RingkasanSpreadingNB } from './PolisNBTab'

const v = (b: Baris, k: string) => b[k] ?? ''
const uang = (k: string) => ({ label: k, angka: true })
const UKURAN_CARGO = 10
const F = FORM_CARGO_NB

/** Sel cargo: Goods kosong = "(belum diisi)" seperti nbfacin. */
const goods = (b: Baris) => v(b, 'GoodNote') || TEKS_CARGO_NB.belumDiisi

// ---------------------------------------------------------------- Object

function Kartu({ judul, children }: { judul?: string; children: ReactNode }) {
  return (
    <section className="claimfacin__nb-blok">
      {judul && <h5 className="claimfacin__nb-subjudul">{judul}</h5>}
      {children}
    </section>
  )
}

/** Form satu cargo (`FormCargo` TabObjectCargo). */
function FormCargoNB({ h, n, b }: { h: Halaman; n: number; b: Baris }) {
  const m = (k: string) => medanCargo(b, k)
  const tambahan = kapalTambahanNB(h, n)
  return (
    <div className="claimfacin__nb-isi">
      <Kartu>
        <CentangNB label={F.topRisk} nilai={benarNB(v(b, 'IsTopRisk'))} />
        <div className="claimfacin__nb-kolom">
          <TeksNB label={F.trading} nilai={v(b, 'TradingNote')} />
          <TeksNB label={F.conveyance} nilai={v(b, 'ConveyanceNote')} />
          <TeksNB label={F.goods} nilai={v(b, 'GoodNote')} />
          <TeksNB label={F.packing} nilai={v(b, 'PackingNote')} />
        </div>
      </Kartu>
      <div className="claimfacin__nb-kolom">
        <div className="claimfacin__nb-tumpuk">
          <Kartu judul={F.grupAgent}>
            <TeksNB label={F.fullName} nilai={m('SurveyAgent.pyFullName')} />
            <TeksNB label={F.address} nilai={m('SurveyAgent.pyAddress')} />
            <TeksNB label={F.city} nilai={m('SurveyAgent.pyCity')} />
            <TeksNB label={F.country} nilai={m('SurveyAgent.pyCountry')} />
          </Kartu>
          <Kartu judul={F.grupLc}>
            <CentangNB label={F.lcStatus} nilai={benarNB(m('LC.LC'))} />
            <TeksNB label={F.lcBank} nilai={m('LC.LCBANK')} />
            <TeksNB label={F.lcNumber} nilai={m('LC.LCNumber')} />
            <TeksNB label={F.lcCondition} nilai={m('LC.LCCondition')} />
            <TeksNB label={F.lcRemark} nilai={m('LC.LCRemark')} />
            <TeksNB label={F.lcEndDate} nilai={tanggalNB(m('LC.LCEndDate'))} />
          </Kartu>
          <Kartu judul={F.grupKapal}>
            <TeksNB label={F.shipName} nilai={m('Ship.NM_SHIP')} />
            <TeksNB label={F.shipType} nilai={m('Ship.SHIPTYPE_NAME')} />
            <TeksNB label={F.grt} nilai={m('Ship.GRT')} />
            <TeksNB label={F.nrt} nilai={m('Ship.NRT')} />
            <TeksNB label={F.dwt} nilai={m('Ship.DWT')} />
            <TeksNB label={F.age} nilai={m('Ship.AGE')} />
            <TeksNB label={F.flagName} nilai={m('Ship.FLAG_NAME')} />
            <TeksNB label={F.construction} nilai={m('Ship.CONST')} />
            <TeksNB label={F.yMake1} nilai={m('Ship.Y_MAKE1')} />
            <TeksNB label={F.remark} nilai={m('Ship.REMARK')} />
          </Kartu>
        </div>
        <div className="claimfacin__nb-tumpuk">
          <Kartu>
            <TeksNB label={F.sailDate} nilai={tanggalNB(m('SailDate'))} />
            <TeksNB label={F.fromRute} nilai={v(b, 'FromRute')} />
            <TeksNB label={F.toRute} nilai={v(b, 'ToRute')} />
            <TeksNB label={F.blNumber} nilai={m('BLNumber')} />
            <TeksNB label={F.invoiceNumber} nilai={m('InvoiceNumber')} />
            <TeksNB label={F.invoiceDate} nilai={tanggalNB(m('InvoiceDate'))} />
            <TeksNB label={F.maxAge} nilai={m('MaxAge')} />
            <TeksNB label={F.maxAgeDetail} nilai={m('MaxAgeDetail')} />
          </Kartu>
        </div>
      </div>
      <Kartu judul={F.grupKapalTambahan}>
        <GridNB
          kolom={F.kolomKapalTambahan}
          baris={tambahan.map((s) => [v(s, 'ID'), v(s, 'NM_SHIP'), v(s, 'VOYAGE_NO')])}
        />
      </Kartu>
    </div>
  )
}

/** Tab Object Marine Cargo (`TabObjectCargo`): No. / Top Risk / SailDate / Goods / Ship Name / Conveyance. */
export function TabObjekCargoNB({ h }: { h: Halaman }) {
  const cargo = daftarNB(h, CARGO_NB)
  return (
    <div className="claimfacin__nb-isi">
      <GridNB
        kolom={GRID_CARGO_NB}
        ukuran={UKURAN_CARGO}
        kosong={TEKS_CARGO_NB.kosong}
        baris={cargo.map((b, i) => [
          String(i + 1),
          <input
            key="t"
            type="checkbox"
            disabled
            readOnly
            checked={benarNB(v(b, 'IsTopRisk'))}
            aria-label={GRID_CARGO_NB[1]}
          />,
          tanggalNB(medanCargo(b, 'SailDate')),
          goods(b),
          medanCargo(b, 'Ship.NM_SHIP'),
          v(b, 'ConveyanceNote'),
        ])}
        rinci={(i) => <FormCargoNB h={h} n={i + 1} b={cargo[i] ?? {}} />}
      />
    </div>
  )
}

// ---------------------------------------------------------------- Coverage

const FC = FORM_COVERAGE_CARGO_NB

/** Form satu coverage cargo (`FormCoverageCargo`) + Warranty Extension + Deductible ber-Deskripsi. */
function FormCoverageCargoNB({ h, daftar, k, c }: { h: Halaman; daftar: string; k: number; c: Baris }) {
  const a = (x: string) => angkaNB(v(c, x))
  const warranty = daftarNB(h, anakNB(daftar, k, 'AdditionalCoverage'))
  return (
    <div className="claimfacin__nb-isi claimfacin__nb-cov">
      <TeksNB label={FC.coverageNote} nilai={v(c, 'CoverageNote')} />
      <TeksNB label={FC.coverageInitial} nilai={v(c, 'CoverageInitial')} />
      <div className="claimfacin__nb-kolom">
        <div className="claimfacin__nb-tumpuk">
          <TeksNB label={FC.currency} nilai={v(c, 'Currency.Name') || '—'} />
          <TeksNB label={FC.rate} nilai={a('Rate')} angka />
          <TeksNB label={FC.limitOfLiability} nilai={a('LimitofLiability')} angka />
          <TeksNB label={FC.discountPercentage} nilai={a('DiscountPercentage')} angka />
        </div>
        <div className="claimfacin__nb-tumpuk">
          <TeksNB label={FC.tsi} nilai={a('TSI')} angka />
          <TeksNB label={FC.premium} nilai={a('Premium')} angka />
          <TeksNB label={FC.minPremium} nilai={a('MinPremium')} angka />
          <TeksNB label={FC.discount} nilai={a('Discount')} angka />
        </div>
      </div>
      <h5 className="claimfacin__nb-subjudul">{FC.grupWarranty}</h5>
      <GridNB kolom={FC.kolomWarranty} baris={warranty.map((w) => [v(w, 'Coverage'), v(w, 'CoverageNote')])} />
      <h5 className="claimfacin__nb-subjudul">{JUDUL_DEDUCTIBLE_NB}</h5>
      <DeductibleNB rows={daftarNB(h, anakNB(daftar, k, 'DeductibleList'))} deskripsi />
    </div>
  )
}

function CoverageCargoNB({ h, n }: { h: Halaman; n: number }) {
  const daftar = coverageCargoNB(n)
  const cov = daftarNB(h, daftar)
  const total = totalGrossCargoNB(cov)
  const [u, tsi, prem, rate] = GRID_COVERAGE_CARGO_NB.total
  return (
    <div className="claimfacin__nb-isi">
      <GridNB
        kolom={GRID_COVERAGE_CARGO_NB.kolom}
        baris={cov.map((c) => [v(c, 'CoverageNote'), v(c, 'CoverageInitial')])}
        rinci={(k) => <FormCoverageCargoNB h={h} daftar={daftar} k={k + 1} c={cov[k] ?? {}} />}
        selaluBuka
      />
      {total.length > 0 && (
        <section className="claimfacin__nb-total">
          <GridNB
            kolom={[u, uang(tsi), uang(prem), uang(rate)]}
            baris={total.map((r) => [r.currency, angkaNB(r.tsi), angkaNB(r.premium), angkaNB(r.rate)])}
          />
        </section>
      )}
    </div>
  )
}

const barisCargo = (b: Baris, i: number) => [
  String(i + 1),
  v(b, 'TradingNote'),
  goods(b),
  v(b, 'PackingNote'),
  v(b, 'ConveyanceNote'),
]

/** Tab Coverage Marine Cargo (`TabCoverageCargo`): No. / Trading / Goods / Packing / Conveyance. */
export function TabCoverageCargoNB({ h }: { h: Halaman }) {
  const cargo = daftarNB(h, CARGO_NB)
  return (
    <div className="claimfacin__nb-isi">
      <GridNB
        kolom={GRID_CARGO_COVERAGE_NB}
        ukuran={UKURAN_CARGO}
        kosong={TEKS_CARGO_NB.kosong}
        baris={cargo.map(barisCargo)}
        rinci={(i) => <CoverageCargoNB h={h} n={i + 1} />}
      />
    </div>
  )
}

// ---------------------------------------------------------------- Spreading

/** Tab Spreading Marine Cargo (`TabSpreadingCargo`): % Share RNM, cargo -> coverage -> detail + baris spreading, ringkasan. */
export function TabSpreadingCargoNB({ h }: { h: Halaman }) {
  const cargo = daftarNB(h, CARGO_NB)
  const s = hitungSpreadingCargoNB(h)
  const [c0, ...cn] = SPREADING_CARGO_NB.kolomCoverage
  return (
    <div className="claimfacin__nb-isi">
      <div className="claimfacin__nb-blok claimfacin__nb-share">
        <TeksNB kiri label={SPREADING_NB.percentShare} nilai={angkaNB(nilaiNB(h, 'PercentShare'))} angka />
      </div>
      <GridNB
        kolom={SPREADING_CARGO_NB.kolomCargo}
        ukuran={UKURAN_CARGO}
        kosong={TEKS_NB.spreadingKosong}
        baris={cargo.map(barisCargo)}
        rinci={(i) => {
          const daftar = coverageCargoNB(i + 1)
          const cov = daftarNB(h, daftar)
          return (
            <GridNB
              kolom={[c0 ?? '', ...cn.map(uang)]}
              baris={cov.map((c) => [
                v(c, 'CoverageNote'),
                angkaNB(v(c, 'Rate')),
                angkaNB(v(c, 'Premium')),
                angkaNB(v(c, 'PremiNusantaraRe')),
              ])}
              rinci={(k) => (
                <div className="claimfacin__nb-isi">
                  <DetailCoverageSpreadingNB c={cov[k] ?? {}} />
                  <BarisSpreadingNB rows={daftarNB(h, anakNB(daftar, k + 1, 'SpreadingList'))} />
                </div>
              )}
            />
          )
        }}
      />
      <RingkasanSpreadingNB s={s} />
    </div>
  )
}
