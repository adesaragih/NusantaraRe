// Isi tab polis NB Fac In hanya-baca keluarga Aneka (`PolisNB.tsx`; perintah work owner 11-10-2026 "tampilan view polis
// yang belum tergenerate tolong diperbaiki. cek tampilan nya di nb facin!"): Object (lokasi -> Risk Address, grid Golf,
// grid okupasi -> objek -> form objek berkartu), Coverage (lokasi -> okupasi / Golf -> objek -> coverage -> form + Deductible,
// Grand Total per mata uang), Spreading (lokasi -> okupasi / Golf -> objek -> coverage -> detail + baris spreading,
// ringkasan). Susunan, judul, dan urutan disalin dari `modul/nbfacin/frontend/components/{TabObjectAneka,FormObjekAneka,
// TabCoverageAneka,FormCoverageAneka,TabSpreadingAneka}.tsx` (pola, bukan impor); tombol (Add / Delete / Choose / Copy /
// Save) dan isian tidak dibawa. Angka 4 desimal (keputusan work owner 11-10-2026 untuk seluruh Claim Fac In).

import type { Baris, Halaman } from '../api'
import {
  JUDUL_DEDUCTIBLE_NB,
  LOKASI_NB,
  SPREADING_NB,
  TEKS_NB,
  angkaNB,
  anakNB,
  benarNB,
  daftarNB,
  labelKode,
  nilaiNB,
  tanggalNB,
} from '../polisNB'
import {
  FORM_COVERAGE_ANEKA_NB,
  FORM_MEMORANDUM_NB,
  GRID_COV_ANEKA_NB,
  GRID_LOKASI_ANEKA_NB,
  GRID_OBJEK_ANEKA_NB,
  GRID_OKUPASI_ANEKA_NB,
  GRUP_BESAR_ANEKA,
  GRUP_KECIL_ANEKA,
  KUNCI_KOLOM_OBJEK_ANEKA,
  MEDAN_BRIGUNA_NB,
  MEDAN_CATATAN_NB,
  MEDAN_MEMORANDUM_NB,
  OPSI_DAY_ANEKA_NB,
  OPSI_METODE_ANEKA_NB,
  OPSI_SECTION_ANEKA_NB,
  OPSI_UNIT_INDEMNITY_NB,
  SPREADING_ANEKA_NB,
  TEKS_ANEKA_NB,
  coverageAnekaNB,
  golfNB,
  hitungSpreadingAnekaNB,
  kelompokMedanAneka,
  kolomObjekAnekaNB,
  objekAnekaNB,
  okupasiAnekaNB,
  tampilOkupasiNB,
  totalCoverageAnekaNB,
  type LiniNB,
  type MedanAneka,
  type TampilanAnekaNB,
} from '../polisNBAneka'
import { GridNB, RadioNB, TeksNB } from './PolisNBDasar'
import { BarisSpreadingNB, DeductibleNB, RingkasanSpreadingNB, RiskAddressNB } from './PolisNBTab'

const v = (b: Baris, k: string) => b[k] ?? ''
const uang = (k: string) => ({ label: k, angka: true })
const UKURAN_LOKASI = 10

type KolomObjek = keyof typeof GRID_OBJEK_ANEKA_NB
const ANGKA_OBJEK: readonly KolomObjek[] = ['tsi', 'limitOfLiability']

// ---------------------------------------------------------------- Object

function nilaiMedan(b: Baris, m: MedanAneka): string {
  const x = v(b, m.kunci)
  switch (m.jenis) {
    case 'angka':
      return angkaNB(x)
    case 'tanggal':
      return tanggalNB(x)
    case 'section':
      return labelKode(OPSI_SECTION_ANEKA_NB, x)
    default:
      return x
  }
}

/** Satu kartu form objek (`kartu` FormObjekAneka nbfacin). */
function KartuAneka({ judul, medan, b }: { judul: string; medan: readonly MedanAneka[]; b: Baris }) {
  return (
    <section className="claimfacin__nb-blok">
      <h5 className="claimfacin__nb-subjudul">{judul}</h5>
      <div className="claimfacin__nb-tiga">
        {medan.map((m) => (
          <TeksNB key={m.kunci} label={m.label} nilai={nilaiMedan(b, m)} angka={m.jenis === 'angka'} />
        ))}
      </div>
    </section>
  )
}

/** Form satu objek (`FormObjekAneka`): kartu besar bertumpuk, kartu kecil sebaris, Loan Details, Note, Memorandum. */
export function FormObjekAnekaNB({ b, t }: { b: Baris; t: TampilanAnekaNB }) {
  if (t.kontainer.length === 0) return <div className="alert alert--info">{TEKS_ANEKA_NB.tanpaKontainer}</div>
  const g = kelompokMedanAneka(t)
  const G = TEKS_ANEKA_NB.grup
  const c1 = t.kontainer.includes('C1')
  const kecil = GRUP_KECIL_ANEKA.filter((x) => g[x].length > 0)
  return (
    <div className="claimfacin__nb-isi">
      {GRUP_BESAR_ANEKA.filter((x) => g[x].length > 0).map((x) => (
        <KartuAneka key={x} judul={G[x]} medan={g[x]} b={b} />
      ))}
      {kecil.length > 0 && (
        <div className="claimfacin__nb-aneka-baris">
          {kecil.map((x) => (
            <KartuAneka key={x} judul={G[x]} medan={g[x]} b={b} />
          ))}
        </div>
      )}
      {c1 && t.briguna && <KartuAneka judul={G.briguna} medan={MEDAN_BRIGUNA_NB} b={b} />}
      {c1 && <KartuAneka judul={G.catatan} medan={[MEDAN_CATATAN_NB]} b={b} />}
      {c1 && t.bonding && <KartuAneka judul={FORM_MEMORANDUM_NB.judul} medan={MEDAN_MEMORANDUM_NB} b={b} />}
    </div>
  )
}

function selObjek(b: Baris, k: KolomObjek): string {
  const x = v(b, KUNCI_KOLOM_OBJEK_ANEKA[k])
  if (k === 'objectName') return x || TEKS_ANEKA_NB.belumDiisi
  if (k === 'section') return labelKode(OPSI_SECTION_ANEKA_NB, x)
  return ANGKA_OBJEK.includes(k) ? angkaNB(x) : x
}

/** Grid objek (okupasi: `kolomGridObjek`; Golf: Object Name / Currency / TSI); baris terbuka = form objek. */
function GridObjekAnekaNB({
  h,
  daftar,
  kolom,
  t,
}: {
  h: Halaman
  daftar: string
  kolom: readonly KolomObjek[]
  t: TampilanAnekaNB
}) {
  const rows = daftarNB(h, daftar)
  return (
    <GridNB
      kolom={kolom.map((k) => (ANGKA_OBJEK.includes(k) ? uang(GRID_OBJEK_ANEKA_NB[k]) : GRID_OBJEK_ANEKA_NB[k]))}
      kosong={TEKS_ANEKA_NB.kosongObjek}
      baris={rows.map((b) => kolom.map((k) => selObjek(b, k)))}
      rinci={(i) => <FormObjekAnekaNB b={rows[i] ?? {}} t={t} />}
    />
  )
}

/** Form lokasi (`FormLokasi` TabObjectAneka): Risk Address (+ Top Risk), grid Golf, grid okupasi. */
function LokasiAnekaNB({ h, n, b, lini }: { h: Halaman; n: number; b: Baris; lini: LiniNB }) {
  const t = lini.tampilan
  const okup = daftarNB(h, okupasiAnekaNB(n))
  return (
    <div className="claimfacin__nb-isi">
      <RiskAddressNB b={b} topRisk={t.topRisk} />
      {t.golf && (
        <section className="claimfacin__nb-blok">
          <GridObjekAnekaNB h={h} daftar={golfNB(n)} kolom={['objectName', 'currency', 'tsi']} t={t} />
        </section>
      )}
      {tampilOkupasiNB(lini.anekaWhen, t) && (
        <section className="claimfacin__nb-blok">
          <GridNB
            kolom={GRID_OKUPASI_ANEKA_NB}
            kosong={TEKS_ANEKA_NB.kosongOkupasi}
            baris={okup.map((o) => [v(o, 'OccupationId'), v(o, 'OccupationName')])}
            rinci={(o) => <GridObjekAnekaNB h={h} daftar={objekAnekaNB(n, o + 1)} kolom={kolomObjekAnekaNB(t)} t={t} />}
          />
        </section>
      )}
    </div>
  )
}

/** Tab Object Aneka (`TabObjectAneka`): grid lokasi No / Top Risk / Location, 10 per halaman. */
export function TabObjekAnekaNB({ h, lini }: { h: Halaman; lini: LiniNB }) {
  const lokasi = daftarNB(h, LOKASI_NB)
  return (
    <div className="claimfacin__nb-isi">
      <GridNB
        kolom={GRID_LOKASI_ANEKA_NB}
        ukuran={UKURAN_LOKASI}
        kosong={TEKS_ANEKA_NB.kosongLokasi}
        baris={lokasi.map((b) => [
          v(b, 'Property.ObjectNo'),
          <input
            key="t"
            type="checkbox"
            disabled
            readOnly
            checked={benarNB(b['Property.IsTopRisk'])}
            aria-label={GRID_LOKASI_ANEKA_NB[1]}
          />,
          v(b, 'Property.RiskLocation.ASMAddress'),
        ])}
        rinci={(i) => <LokasiAnekaNB h={h} n={i + 1} b={lokasi[i] ?? {}} lini={lini} />}
      />
    </div>
  )
}

// ---------------------------------------------------------------- Coverage

const FC = FORM_COVERAGE_ANEKA_NB

/** Form satu coverage (`FormCoverageAneka`): satu kisi mengalir, Sub Limit selebar penuh, lalu Deductible. */
function FormCoverageAnekaNB({
  h,
  daftar,
  c,
  b,
  t,
}: {
  h: Halaman
  daftar: string
  c: number
  b: Baris
  t: TampilanAnekaNB
}) {
  const metode = v(b, 'CalculateMethod_FacIn')
  const a = (k: string) => angkaNB(v(b, k))
  return (
    <div className="claimfacin__nb-isi claimfacin__nb-cov">
      <div className="claimfacin__nb-tiga">
        {!t.tanamanHasil && <TeksNB label={FC.coverageNote} nilai={v(b, 'CoverageNote')} />}
        <TeksNB label={FC.calculateMethod} nilai={labelKode(OPSI_METODE_ANEKA_NB, metode)} />
        {metode === '1' && <TeksNB label={FC.proRate} nilai={angkaNB(nilaiNB(h, 'ProRatePercent'))} angka />}
        {metode === '2' && <TeksNB label={FC.shortPeriod} nilai={a('PctShortPeriod')} angka />}
        <RadioNB label={FC.day} pilihan={OPSI_DAY_ANEKA_NB} nilai={v(b, 'Day')} />
        {t.currencyCoverage && <TeksNB label={FC.currency} nilai={v(b, 'Currency.Name')} />}
        <TeksNB label={FC.tsi} nilai={a('TSI')} angka />
        <TeksNB label={FC.loading} nilai={a('Loading')} angka />
        <TeksNB label={FC.scale} nilai={a('LostLimit')} angka />
        <TeksNB label={FC.rate} nilai={a('Rate')} angka />
        {t.mbd && <TeksNB label={FC.indemnity} nilai={a('Indemnity')} angka />}
        {t.mbd && <TeksNB label={FC.unit} nilai={labelKode(OPSI_UNIT_INDEMNITY_NB, v(b, 'Unit'))} />}
        {t.mbd && <TeksNB label={FC.pctIndemnity} nilai={a('IndemnityPercentage')} angka />}
        <TeksNB label={FC.premium} nilai={a('Premium')} angka />
        <TeksNB label={FC.pctDiscount} nilai={a('DiscountPercentage')} angka />
        <TeksNB label={FC.discount} nilai={a('Discount')} angka />
        <TeksNB label={FC.lol} nilai={a('LimitofLiability')} angka />
        <TeksNB label={FC.pctLol} nilai={a('PctLoL')} angka />
        <TeksNB label={FC.subLimit} nilai={v(b, 'SubLimitNote')} lebar />
      </div>
      <h5 className="claimfacin__nb-subjudul">{JUDUL_DEDUCTIBLE_NB}</h5>
      <DeductibleNB rows={daftarNB(h, anakNB(daftar, c, 'DeductibleList'))} />
    </div>
  )
}

/** Grid `.CoverageList` objek (Coverage / Gross Premium); baris terbuka = form coverage. */
function GridCoverageAnekaNB({ h, objek, n, t }: { h: Halaman; objek: string; n: number; t: TampilanAnekaNB }) {
  const daftar = coverageAnekaNB(objek, n)
  const rows = daftarNB(h, daftar)
  return (
    <GridNB
      kolom={[GRID_COV_ANEKA_NB.coverage[0], uang(GRID_COV_ANEKA_NB.coverage[1])]}
      kosong={GRID_COV_ANEKA_NB.kosongCoverage}
      baris={rows.map((b) => [v(b, 'CoverageNote'), angkaNB(v(b, 'Premium'))])}
      rinci={(c) => <FormCoverageAnekaNB h={h} daftar={daftar} c={c + 1} b={rows[c] ?? {}} t={t} />}
    />
  )
}

/** Grid objek tab Coverage (Object Name / Currency / TSI); baris terbuka = daftar coverage. */
function GridCovObjekAnekaNB({ h, daftar, t }: { h: Halaman; daftar: string; t: TampilanAnekaNB }) {
  const rows = daftarNB(h, daftar)
  const [k1, k2, k3] = GRID_COV_ANEKA_NB.objek
  return (
    <GridNB
      kolom={[k1, k2, uang(k3)]}
      baris={rows.map((b) => [v(b, 'ObjectName'), v(b, 'Currency.Name'), angkaNB(v(b, 'TSI'))])}
      rinci={(i) => <GridCoverageAnekaNB h={h} objek={daftar} n={i + 1} t={t} />}
    />
  )
}

function CoverageLokasiAnekaNB({ h, n, lini }: { h: Halaman; n: number; lini: LiniNB }) {
  const t = lini.tampilan
  const okup = daftarNB(h, okupasiAnekaNB(n))
  const total = totalCoverageAnekaNB(h, n)
  const [u, tsi, prem, rate] = GRID_COV_ANEKA_NB.total
  return (
    <div className="claimfacin__nb-isi">
      {tampilOkupasiNB(lini.anekaWhen, t) && (
        <GridNB
          kolom={GRID_COV_ANEKA_NB.okupasi}
          baris={okup.map((o) => [v(o, 'OccupationName')])}
          rinci={(o) => <GridCovObjekAnekaNB h={h} daftar={objekAnekaNB(n, o + 1)} t={t} />}
        />
      )}
      {t.golf && <GridCovObjekAnekaNB h={h} daftar={golfNB(n)} t={t} />}
      <section className="claimfacin__nb-total">
        <GridNB
          kolom={[u, uang(tsi), uang(prem), uang(rate)]}
          baris={total.map((r) => [r.currency, angkaNB(r.tsi), angkaNB(r.premium), angkaNB(r.rate)])}
        />
      </section>
    </div>
  )
}

/** Tab Coverage Aneka (`TabCoverageAneka`): grid lokasi No. / Location, 10 per halaman. */
export function TabCoverageAnekaNB({ h, lini }: { h: Halaman; lini: LiniNB }) {
  const lokasi = daftarNB(h, LOKASI_NB)
  return (
    <div className="claimfacin__nb-isi">
      <GridNB
        kolom={GRID_COV_ANEKA_NB.lokasi}
        ukuran={UKURAN_LOKASI}
        baris={lokasi.map((b) => [v(b, 'Property.ObjectNo'), v(b, 'Property.RiskLocation.ASMAddress')])}
        rinci={(i) => <CoverageLokasiAnekaNB h={h} n={i + 1} lini={lini} />}
      />
    </div>
  )
}

// ---------------------------------------------------------------- Spreading

const DS = SPREADING_ANEKA_NB.detailCoverage

/** Detail coverage tab Spreading (`SpreadingItemAneka`): empat angka lalu baris spreading. */
function DetailSpreadingAnekaNB({ h, daftar, c, b }: { h: Halaman; daftar: string; c: number; b: Baris }) {
  return (
    <div className="claimfacin__nb-isi">
      <dl className="claimfacin__nb-angka-lima claimfacin__nb-angka-empat">
        {(
          [
            [DS.tsi, 'TSI'],
            [DS.tsiNusantaraRe, 'TSINusantaraRe'],
            [DS.premium, 'Premium'],
            [DS.premiNusantaraRe, 'PremiNusantaraRe'],
          ] as const
        ).map(([l, k]) => (
          <div key={k}>
            <dt>{l}</dt>
            <dd>{angkaNB(v(b, k))}</dd>
          </div>
        ))}
      </dl>
      <BarisSpreadingNB rows={daftarNB(h, anakNB(daftar, c, 'SpreadingList'))} />
    </div>
  )
}

/** Objek -> coverage tab Spreading; Golf memakai kolom coverage Golf. */
function SpreadingObjekAnekaNB({ h, daftar, golf }: { h: Halaman; daftar: string; golf: boolean }) {
  const rows = daftarNB(h, daftar)
  const S = SPREADING_ANEKA_NB
  return (
    <GridNB
      kolom={golf ? S.kolomObjekGolf : [S.kolomObjek[0], S.kolomObjek[1], uang(S.kolomObjek[2])]}
      baris={rows.map((b) =>
        golf ? [v(b, 'ObjectName')] : [v(b, 'ObjectName'), v(b, 'Currency.Name'), angkaNB(v(b, 'TSI'))],
      )}
      rinci={(i) => {
        const cov = coverageAnekaNB(daftar, i + 1)
        const cs = daftarNB(h, cov)
        return (
          <GridNB
            kolom={
              golf
                ? [S.kolomCoverageGolf[0], ...S.kolomCoverageGolf.slice(1).map(uang)]
                : [S.kolomCoverage[0], uang(S.kolomCoverage[1])]
            }
            baris={cs.map((c) =>
              golf
                ? [
                    v(c, 'CoverageNote'),
                    angkaNB(v(c, 'Loading')),
                    angkaNB(v(c, 'Rate')),
                    angkaNB(v(c, 'TSI')),
                    angkaNB(v(c, 'Premium')),
                  ]
                : [v(c, 'CoverageNote'), angkaNB(v(c, 'Premium'))],
            )}
            rinci={(m) => <DetailSpreadingAnekaNB h={h} daftar={cov} c={m + 1} b={cs[m] ?? {}} />}
          />
        )
      }}
    />
  )
}

function SpreadingLokasiAnekaNB({ h, n }: { h: Halaman; n: number }) {
  const okup = daftarNB(h, okupasiAnekaNB(n))
  const golf = daftarNB(h, golfNB(n))
  return (
    <div className="claimfacin__nb-isi">
      {okup.length > 0 && (
        <GridNB
          kolom={SPREADING_ANEKA_NB.kolomOkupasi}
          baris={okup.map((o) => [v(o, 'OccupationName')])}
          rinci={(o) => <SpreadingObjekAnekaNB h={h} daftar={objekAnekaNB(n, o + 1)} golf={false} />}
        />
      )}
      {golf.length > 0 && <SpreadingObjekAnekaNB h={h} daftar={golfNB(n)} golf />}
    </div>
  )
}

/** Tab Spreading Aneka (`TabSpreadingAneka`): % Share RNM, lokasi No / Location, ringkasan (totalAneka + ringkasan FIRE). */
export function TabSpreadingAnekaNB({ h }: { h: Halaman }) {
  const lokasi = daftarNB(h, LOKASI_NB)
  const s = hitungSpreadingAnekaNB(h)
  return (
    <div className="claimfacin__nb-isi">
      <div className="claimfacin__nb-blok claimfacin__nb-share">
        <TeksNB kiri label={SPREADING_NB.percentShare} nilai={angkaNB(nilaiNB(h, 'PercentShare'))} angka />
      </div>
      <GridNB
        kolom={SPREADING_ANEKA_NB.kolomLokasi}
        ukuran={UKURAN_LOKASI}
        kosong={TEKS_NB.spreadingKosong}
        baris={lokasi.map((b) => [v(b, 'Property.ObjectNo'), v(b, 'Property.RiskLocation.ASMAddress')])}
        rinci={(i) => <SpreadingLokasiAnekaNB h={h} n={i + 1} />}
      />
      <RingkasanSpreadingNB s={s} />
    </div>
  )
}
