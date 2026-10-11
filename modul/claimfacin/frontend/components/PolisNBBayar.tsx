// Isi tab polis NB Fac In hanya-baca Premium Deduction / Brokerage Fee, Payment, dan Inw Fac Cedant Panels untuk FIRE,
// keluarga Aneka, dan Marine Cargo (`PolisNB.tsx`; perintah work owner 11-10-2026 "tampilan view polis yang belum
// tergenerate tolong diperbaiki. cek tampilan nya di nb facin!"). Susunan dan judul disalin dari
// `modul/nbfacin/frontend/components/{TabDeduction,TabPayment,TabCedant}.tsx` (pola, bukan impor); isian, tombol, dan
// Save tidak dibawa. Nilai Deduction / Payment = TERSIMPAN (keputusan work owner, `polisNBBayar.ts`).

import type { Baris, Halaman } from '../api'
import {
  CEDANT_NB,
  CEDANT_NB_LABEL,
  LOKASI_NB,
  OPSI_SHARE_CEDANT_NB,
  TEKS_NB,
  angkaNB,
  anakNB,
  benarNB,
  coverageNB,
  daftarNB,
  hitungSpreadingNB,
  itemNB,
  labelKode,
  nilaiNB,
  premiItemNB,
  premiRNMItemNB,
  tanggalNB,
  type HasilSpreadingNB,
} from '../polisNB'
import {
  coverageAnekaNB,
  golfNB,
  hitungSpreadingAnekaNB,
  itemDasarAnekaNB,
  objekAnekaNB,
  okupasiAnekaNB,
  type LiniNB,
} from '../polisNBAneka'
import {
  CEDANT_DETAIL_NB,
  DEDUCTION_ANEKA_NB,
  DEDUCTION_CARGO_NB,
  DEDUCTION_NB,
  MATA_UANG_NB,
  PAYMENT_NB,
  bayarPolisNB,
  bayarUangNB,
  cicilanNB,
  dasarCedantNB,
  itemDasarFireNB,
  mataUangCedantNB,
  pesanPersenCicilanNB,
  shareOfCedingNB,
} from '../polisNBBayar'
import { CARGO_NB, coverageCargoNB, hitungSpreadingCargoNB, itemDasarCargoNB } from '../polisNBCargo'
import { CentangNB, GridNB, RadioNB, TeksNB } from './PolisNBDasar'
import { RingkasanSpreadingNB } from './PolisNBTab'

const v = (b: Baris, k: string) => b[k] ?? ''
const uang = (k: string) => ({ label: k, angka: true })
const UKURAN = 10
const D = DEDUCTION_NB

/** Ringkasan spreading per lini (FIRE `hitungSpreadingNB`, Aneka `totalAneka`, Cargo `totalCargo`). */
function spreadingLini(h: Halaman, lini: LiniNB): HasilSpreadingNB {
  if (lini.aneka) return hitungSpreadingAnekaNB(h)
  if (lini.cargo) return hitungSpreadingCargoNB(h)
  return hitungSpreadingNB(h)
}

// ---------------------------------------------------------------- Premium Deduction: daftar komisi

/** Detail coverage daftar komisi: empat angka (urutan dan label per lini) + grid spreading `DEDUCTION.kolomSpreading`. */
function DetailKomisiNB({
  c,
  rows,
  detail,
}: {
  c: Baris
  rows: Baris[]
  detail: readonly (readonly [string, string])[]
}) {
  return (
    <div className="claimfacin__nb-isi">
      <dl className="claimfacin__nb-angka-lima claimfacin__nb-angka-empat">
        {detail.map(([l, k]) => (
          <div key={l}>
            <dt>{l}</dt>
            <dd>{angkaNB(v(c, k))}</dd>
          </div>
        ))}
      </dl>
      <GridNB
        kosong={TEKS_NB.spreadingKosong}
        kolom={[D.kolomSpreading[0], ...D.kolomSpreading.slice(1).map(uang)]}
        baris={rows.map((r) => [
          v(r, 'TreatyName') || v(r, 'TreatyType'),
          angkaNB(v(r, 'SharePercentage')),
          angkaNB(v(r, 'TSISpreaded')),
          angkaNB(v(r, 'ClaimEstimation')),
          angkaNB(v(r, 'PremiumSpreaded')),
        ])}
      />
    </div>
  )
}

const DETAIL_FIRE = [
  [D.detailCoverage.tsiLiability, 'TSILiability'],
  [D.detailCoverage.premium, 'Premium'],
  [D.detailCoverage.tsiNusantaraRe, 'TSINusantaraRe'],
  [D.detailCoverage.premiNusantaraRe, 'PremiNusantaraRe'],
] as const
const DA = DEDUCTION_ANEKA_NB.detailCoverage
/** Aneka / Cargo: "TSI Liability" terikat ke `.TSI` coverage (VERBATIM `CommisionItemAneka` L795). */
const DETAIL_ANEKA = [
  [DA.tsi, 'TSI'],
  [DA.tsiNusantaraRe, 'TSINusantaraRe'],
  [DA.premium, 'Premium'],
  [DA.premiNusantaraRe, 'PremiNusantaraRe'],
] as const

/** Daftar komisi FIRE (`DaftarKomisi`): lokasi -> item -> coverage -> detail. */
function KomisiFireNB({ h }: { h: Halaman }) {
  const lokasi = daftarNB(h, LOKASI_NB)
  return (
    <GridNB
      kolom={D.kolomLokasi}
      ukuran={UKURAN}
      baris={lokasi.map((b) => [
        v(b, 'Property.ObjectNo'),
        v(b, 'Property.ObjectName'),
        v(b, 'Property.RiskLocation.ASMAddress'),
      ])}
      rinci={(l) => {
        const items = daftarNB(h, itemNB(l + 1))
        return (
          <GridNB
            kolom={[D.kolomItem[0], D.kolomItem[1], ...D.kolomItem.slice(2).map(uang)]}
            baris={items.map((it, j) => {
              const cov = daftarNB(h, coverageNB(l + 1, j + 1))
              return [
                v(it, 'ItemType'),
                v(it, 'Currency'),
                angkaNB(v(it, 'TSIObjectItem')),
                angkaNB(premiItemNB(cov, it)),
                angkaNB(premiRNMItemNB(cov)),
              ]
            })}
            rinci={(j) => {
              const daftar = coverageNB(l + 1, j + 1)
              const cov = daftarNB(h, daftar)
              return (
                <GridNB
                  kolom={[D.kolomCoverage[0], ...D.kolomCoverage.slice(1).map(uang)]}
                  baris={cov.map((c) => [
                    v(c, 'OLDID'),
                    angkaNB(v(c, 'Rate')),
                    angkaNB(v(c, 'Premium')),
                    angkaNB(v(c, 'PremiNusantaraRe')),
                  ])}
                  rinci={(m) => (
                    <DetailKomisiNB
                      c={cov[m] ?? {}}
                      rows={daftarNB(h, anakNB(daftar, m + 1, 'SpreadingList'))}
                      detail={DETAIL_FIRE}
                    />
                  )}
                />
              )
            }}
          />
        )
      }}
    />
  )
}

/** Objek -> coverage daftar komisi Aneka (Golf: kolom Object Name saja). */
function KomisiObjekAnekaNB({ h, daftar, golf }: { h: Halaman; daftar: string; golf: boolean }) {
  const A = DEDUCTION_ANEKA_NB
  const rows = daftarNB(h, daftar)
  return (
    <GridNB
      kolom={golf ? A.kolomObjekGolf : [A.kolomObjek[0], A.kolomObjek[1], uang(A.kolomObjek[2])]}
      baris={rows.map((b) =>
        golf ? [v(b, 'ObjectName')] : [v(b, 'ObjectName'), v(b, 'Currency.Name'), angkaNB(v(b, 'TSI'))],
      )}
      rinci={(i) => {
        const cov = coverageAnekaNB(daftar, i + 1)
        const cs = daftarNB(h, cov)
        return (
          <GridNB
            kolom={[A.kolomCoverage[0], uang(A.kolomCoverage[1])]}
            baris={cs.map((c) => [v(c, 'CoverageNote'), angkaNB(v(c, 'Premium'))])}
            rinci={(m) => (
              <DetailKomisiNB
                c={cs[m] ?? {}}
                rows={daftarNB(h, anakNB(cov, m + 1, 'SpreadingList'))}
                detail={DETAIL_ANEKA}
              />
            )}
          />
        )
      }}
    />
  )
}

/** Daftar komisi Aneka (`DaftarKomisiAneka`): lokasi -> okupasi -> objek / objek Golf -> coverage -> detail. */
function KomisiAnekaNB({ h }: { h: Halaman }) {
  const A = DEDUCTION_ANEKA_NB
  const lokasi = daftarNB(h, LOKASI_NB)
  return (
    <GridNB
      kolom={A.kolomLokasi}
      ukuran={UKURAN}
      baris={lokasi.map((b) => [v(b, 'Property.ObjectNo'), v(b, 'Property.RiskLocation.ASMAddress')])}
      rinci={(l) => {
        const okup = daftarNB(h, okupasiAnekaNB(l + 1))
        const golf = daftarNB(h, golfNB(l + 1))
        return (
          <div className="claimfacin__nb-isi">
            {okup.length > 0 && (
              <GridNB
                kolom={A.kolomOkupasi}
                baris={okup.map((o) => [v(o, 'OccupationName')])}
                rinci={(o) => <KomisiObjekAnekaNB h={h} daftar={objekAnekaNB(l + 1, o + 1)} golf={false} />}
              />
            )}
            {golf.length > 0 && <KomisiObjekAnekaNB h={h} daftar={golfNB(l + 1)} golf />}
          </div>
        )
      }}
    />
  )
}

/** Daftar komisi Marine Cargo (`DaftarKomisiCargo`): cargo -> coverage (Coverage / TSI / Gross Premium) -> detail. */
function KomisiCargoNB({ h }: { h: Halaman }) {
  const C = DEDUCTION_CARGO_NB
  const cargo = daftarNB(h, CARGO_NB)
  return (
    <GridNB
      kolom={C.kolomCargo}
      ukuran={UKURAN}
      baris={cargo.map((b, i) => [
        String(i + 1),
        v(b, 'TradingNote'),
        v(b, 'GoodNote'),
        v(b, 'PackingNote'),
        v(b, 'ConveyanceNote'),
      ])}
      rinci={(i) => {
        const daftar = coverageCargoNB(i + 1)
        const cs = daftarNB(h, daftar)
        return (
          <GridNB
            kolom={[C.kolomCoverage[0], uang(C.kolomCoverage[1]), uang(C.kolomCoverage[2])]}
            baris={cs.map((c) => [v(c, 'CoverageNote'), angkaNB(v(c, 'TSI')), angkaNB(v(c, 'Premium'))])}
            rinci={(m) => (
              <DetailKomisiNB
                c={cs[m] ?? {}}
                rows={daftarNB(h, anakNB(daftar, m + 1, 'SpreadingList'))}
                detail={DETAIL_ANEKA}
              />
            )}
          />
        )
      }}
    />
  )
}

/** Tab Premium Deduction / Brokerage Fee (`TabDeduction`): persen, daftar komisi, ringkasan, dua kotak per mata uang. */
export function TabDeductionNB({ h, lini }: { h: Halaman; lini: LiniNB }) {
  const uangList = daftarNB(h, MATA_UANG_NB)
  return (
    <div className="claimfacin__nb-isi">
      <div className="claimfacin__nb-blok claimfacin__nb-kolom">
        <TeksNB kiri label={D.pctPremiumDeduction} nilai={angkaNB(bayarPolisNB(h, 'RICommision'))} angka />
        <div className="claimfacin__nb-tumpuk">
          <TeksNB kiri label={D.pctBrokerageFee} nilai={angkaNB(bayarPolisNB(h, 'PctBrokerageFee'))} angka />
          <CentangNB label={D.excludePpn} nilai={benarNB(nilaiNB(h, 'PPnCheck'))} />
        </div>
      </div>
      {lini.aneka ? <KomisiAnekaNB h={h} /> : lini.cargo ? <KomisiCargoNB h={h} /> : <KomisiFireNB h={h} />}
      <RingkasanSpreadingNB s={spreadingLini(h, lini)} />
      <div className="claimfacin__nb-kolom">
        <GridNB
          kolom={[D.kolomMataUang, uang(D.totalPremiumDeduction), uang(D.totalPremiumDeduction2)]}
          baris={uangList.map((b) => [
            v(b, 'Name'),
            angkaNB(bayarUangNB(b, 'Commision')),
            angkaNB(bayarUangNB(b, 'Deduction2')),
          ])}
        />
        <GridNB
          kolom={[D.kolomMataUang, uang(D.totalBrokerageFee), uang(D.pph), uang(D.ppn)]}
          baris={uangList.map((b) => [
            v(b, 'Name'),
            angkaNB(bayarUangNB(b, 'BrokerageFee')),
            angkaNB(bayarUangNB(b, 'PPh')),
            angkaNB(bayarUangNB(b, 'PPN')),
          ])}
        />
      </div>
    </div>
  )
}

// ---------------------------------------------------------------- Payment

/** Tab Payment (`TabPayment`, sama untuk semua lini): Installment, grid per mata uang -> rincian cicilan. */
export function TabPaymentNB({ h }: { h: Halaman }) {
  const P = PAYMENT_NB
  const uangList = daftarNB(h, MATA_UANG_NB)
  return (
    <div className="claimfacin__nb-isi">
      <h4 className="claimfacin__nb-judul">{P.judul}</h4>
      <div className="claimfacin__nb-blok claimfacin__nb-share">
        <TeksNB kiri label={P.installment} nilai={bayarPolisNB(h, 'Installment')} />
      </div>
      <GridNB
        kolom={[P.kolomMataUang[0], uang(P.kolomMataUang[1]), uang(P.kolomMataUang[2])]}
        baris={uangList.map((b) => [
          v(b, 'Name'),
          angkaNB(bayarUangNB(b, 'Premium')),
          angkaNB(v(b, 'SumTotalPayment')),
        ])}
        rinci={(n) => {
          const cicilan = cicilanNB(h, n + 1)
          const pesan = pesanPersenCicilanNB(cicilan)
          return (
            <div className="claimfacin__nb-isi">
              <h5 className="claimfacin__nb-subjudul">{P.judulCicilan}</h5>
              {pesan !== '' && <div className="field__error">{pesan}</div>}
              <GridNB
                kolom={[
                  P.no,
                  P.dueDate,
                  ...[P.persen, P.premium, P.deduction, P.deduction2, P.brokerage, P.pph, P.ppn, P.discount].map(uang),
                  uang(P.totalPayment),
                ]}
                baris={cicilan.map((c) => [
                  v(c, 'InstallmentNo'),
                  tanggalNB(v(c, 'DueDate')),
                  angkaNB(v(c, 'InstallmentPercentage')),
                  angkaNB(v(c, 'Premium')),
                  angkaNB(v(c, 'RICommision')),
                  angkaNB(v(c, 'Deduction2')),
                  angkaNB(v(c, 'BrokerageFee')),
                  angkaNB(v(c, 'PPh')),
                  angkaNB(v(c, 'PPN')),
                  angkaNB(v(c, 'Discount')),
                  angkaNB(v(c, 'PaymentTotal')),
                ])}
              />
            </div>
          )
        }}
      />
    </div>
  )
}

// ---------------------------------------------------------------- Inw Fac Cedant Panels

/**
 * Tab Inw Fac Cedant Panels (`TabCedant`): blok atas (Total TSI / Premi RNM ber-tag IDR, nilai tersimpan), grid cedant
 * -> rincian (Ceding, % Share "of", CurrencyList cedant dihitung), ringkasan per mata uang (Aneka / Cargo tanpa Top Risk).
 */
export function TabCedantLengkapNB({ h, lini }: { h: Halaman; lini: LiniNB }) {
  const C = CEDANT_NB_LABEL
  const X = CEDANT_DETAIL_NB
  const rows = daftarNB(h, CEDANT_NB)
  const tipe = nilaiNB(h, 'ShareCedantType')
  const items = lini.aneka ? itemDasarAnekaNB(h) : lini.cargo ? itemDasarCargoNB(h) : itemDasarFireNB(h)
  const dasar = dasarCedantNB(h, items)
  const dari = shareOfCedingNB(
    tipe,
    nilaiNB(h, 'PercentShare'),
    h.nilai['OfferFacIn.QuotationData.ShareOfCeding'] ?? '',
  )
  return (
    <div className="claimfacin__nb-isi">
      <div className="claimfacin__nb-blok claimfacin__nb-cedant">
        <TeksNB kiri label={C.sob} nilai={h.nilai['OfferFacIn.QuotationData.SobName'] ?? ''} />
        <RadioNB
          kiri
          label={C.shareCedantType}
          pilihan={OPSI_SHARE_CEDANT_NB.map((o) => o.value)}
          opsiLabel={(p) => labelKode(OPSI_SHARE_CEDANT_NB, p)}
          nilai={tipe}
        />
        <TeksNB kiri label={C.percentShare} nilai={angkaNB(nilaiNB(h, 'PercentShare'))} angka />
        <TeksNB
          kiri
          label={`${C.totalTsi} ${X.mataUangTotal}`}
          nilai={angkaNB(nilaiNB(h, 'TotalTSINusaReSpreading'))}
          angka
        />
        <TeksNB
          kiri
          label={`${C.totalPremi} ${X.mataUangTotal}`}
          nilai={angkaNB(nilaiNB(h, 'TotalPremiNusaRe'))}
          angka
        />
      </div>
      <div className="claimfacin__nb-cedant-grid">
        <GridNB
          kolom={[C.kolomCeding, uang(C.kolomShare)]}
          baris={rows.map((r) => [v(r, 'CedingCoName'), angkaNB(v(r, 'ShareCeding'))])}
          rinci={(i) => {
            const r = rows[i] ?? {}
            const cu = mataUangCedantNB(v(r, 'ShareCeding'), tipe, dasar)
            return (
              <div className="claimfacin__nb-isi">
                <TeksNB kiri label={X.kolomCeding} nilai={v(r, 'CedingCoName')} />
                <TeksNB kiri label={X.kolomShare} nilai={`${angkaNB(v(r, 'ShareCeding'))} ${X.dari} ${dari}`} />
                <GridNB
                  kolom={[X.kolomMataUang[0], uang(X.kolomMataUang[1]), uang(X.kolomMataUang[2])]}
                  baris={cu.map((u) => [u.currency, angkaNB(u.tsi), angkaNB(u.premium)])}
                />
              </div>
            )
          }}
        />
      </div>
      <RingkasanSpreadingNB s={spreadingLini(h, lini)} hanyaMataUang tanpaTopRisk={!lini.fire} />
    </div>
  )
}
