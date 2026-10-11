// Isi tab detail polis NB Fac In hanya-baca (`PolisNB.tsx`): Object (+ tujuh sub-tab baris objek), Coverage (lokasi ->
// item -> coverage -> form + Deductible, ringkasan), Spreading, dan Inw Fac Cedant Panels. Susunan, judul, dan format
// angka disalin dari `modul/nbfacin/frontend/components/{TabObject,SubTab*,TabCoverage,FormCoverage,GridDeductible,
// TabSpreading,TabCedant}.tsx` (pola, bukan impor); tombol (Add / Delete / Choose / Save / Copy) tidak dibawa.

import { useState } from 'react'

import type { Baris, Halaman } from '../api'
import {
  ALAMAT_NB,
  AKUMULASI_KOSONG_NB,
  CEDANT_NB,
  CEDANT_NB_LABEL,
  FLOOD_STATUS_TAMPIL_NB,
  FORM_COV_NB,
  FORM_DEDUCTIBLE_NB,
  FORM_FEA_NB,
  FORM_ITEM_NB,
  FORM_KERUGIAN_NB,
  FORM_OKUPASI_NB,
  JUDUL_DEDUCTIBLE_NB,
  KOLOM_COV_ITEM_NB,
  KOLOM_COV_OBJEK_NB,
  KOLOM_COV_TOTAL_NB,
  KOLOM_COVERAGE_NB,
  KOLOM_DEDUCTIBLE_NB,
  KOLOM_FEA_NB,
  KOLOM_ITEM_NB,
  KOLOM_KERUGIAN_NB,
  KOLOM_KLAIM_INTERNAL_NB,
  KOLOM_OBJEK_NB,
  KOLOM_OKUPASI_NB,
  KOLOM_TOTAL_ITEM_NB,
  KONDISI_LAIN_NB,
  LAIN_SEKITAR_NB,
  LOKASI_NB,
  LOSS_RATIO_NB,
  MEDAN_SISI_NB,
  MINMAX_OR_NB,
  OBJECT_TYPE_LAINNYA_NB,
  OPSI_CONDITION_NB,
  OPSI_CONSTRUCTION_NB,
  OPSI_COVERAGE_BASIS_NB,
  OPSI_FLOOD_AREA_NB,
  OPSI_FLOOD_STATUS_NB,
  OPSI_FLOOR_NB,
  OPSI_HOUSEKEEPING_NB,
  OPSI_INDEMNITY_UNIT_NB,
  OPSI_KONDISI_DEDUCTIBLE_NB,
  OPSI_MINMAX_NB,
  OPSI_OWNERSHIP_NB,
  OPSI_ROOF_NB,
  OPSI_SHARE_CEDANT_NB,
  OPSI_TYPE_DEDUCTIBLE2_NB,
  OPSI_TYPE_DEDUCTIBLE_NB,
  OPSI_WALL_NB,
  PILIHAN_DAY_COV_NB,
  RINGKASAN_COV_ITEM_NB,
  RINGKASAN_COV_UANG_NB,
  SISI_NB,
  SPREADING_NB,
  SUBTAB_OBJEK_NB,
  TEKS_NB,
  TYPE_IN_AMOUNT_NB,
  UKURAN_OBJEK_NB,
  anakNB,
  angkaNB,
  benarNB,
  coverageNB,
  daftarNB,
  deductibleNB,
  hitungSpreadingNB,
  itemNB,
  labelKode,
  nilaiNB,
  premiItemNB,
  premiRNMItemNB,
  ringkasanCoverageNB,
  spreadingNB,
  tahunNB,
  tanggalNB,
  totalCoverageLokasiNB,
  totalTSIItemNB,
  type SubTabObjekNB,
  type TotalSpreadingNB,
} from '../polisNB'
import { CentangNB, GridNB, RadioNB, TeksNB } from './PolisNBDasar'
import { TabGeser } from './TabGeser'

const v = (b: Baris, k: string) => b[k] ?? ''
const uang = (k: string) => ({ label: k, angka: true })

/** Kolom grid berjudul angka (rata kanan) menurut label. */
const kolomAngka = (labels: readonly string[], angka: readonly string[]) =>
  labels.map((l) => (angka.includes(l) ? uang(l) : l))

// ---------------------------------------------------------------- Object

/** Tab Object (`TabObject`): grid LocationList 10 per halaman; baris terbuka = sub-tab objek. */
export function TabObjekNB({ h }: { h: Halaman }) {
  const [sub, setSub] = useState<Record<number, SubTabObjekNB>>({})
  const lokasi = daftarNB(h, LOKASI_NB)
  return (
    <div className="claimfacin__nb-isi">
      <GridNB
        kolom={KOLOM_OBJEK_NB}
        ukuran={UKURAN_OBJEK_NB}
        baris={lokasi.map((b) => [
          <input
            key="t"
            type="checkbox"
            disabled
            readOnly
            checked={benarNB(b['Property.IsTopRisk'])}
            aria-label={KOLOM_OBJEK_NB[0]}
          />,
          v(b, 'Property.ObjectNo'),
          v(b, 'Property.ObjectName'),
          v(b, 'Property.RiskLocation.ASMAddress'),
        ])}
        rinci={(i) => {
          const aktif = sub[i] ?? SUBTAB_OBJEK_NB[0]
          return (
            <div className="claimfacin__tab">
              <TabGeser tab={SUBTAB_OBJEK_NB} aktif={aktif} onPilih={(t) => setSub((s) => ({ ...s, [i]: t }))} />
              <div key={aktif} role="tabpanel" className="claimfacin__tab-isi">
                <SubObjekNB h={h} n={i + 1} b={lokasi[i] ?? {}} tab={aktif} />
              </div>
            </div>
          )
        }}
      />
    </div>
  )
}

function SubObjekNB({ h, n, b, tab }: { h: Halaman; n: number; b: Baris; tab: SubTabObjekNB }) {
  switch (tab) {
    case 'Object Address':
      return <AlamatNB b={b} />
    case 'Surrounding Risk':
      return <SekitarNB b={b} />
    case 'Object Item':
      return <ItemObjekNB h={h} n={n} />
    case 'Occupation':
      return <OkupasiNB h={h} n={n} />
    case 'FEA':
      return <FeaNB h={h} n={n} />
    case 'Loss Record':
      return <KerugianNB h={h} n={n} b={b} />
    case 'Loss Record Internal':
      return <KlaimInternalNB h={h} n={n} />
  }
}

/** Sub-tab Object Address (Section ObjectDetails): tiga blok. */
function AlamatNB({ b }: { b: Baris }) {
  const p = (k: string) => v(b, 'Property.' + k)
  const bc = (k: string) => p('BuildingConstruction.' + k)
  return (
    <div className="claimfacin__nb-isi">
      <section className="claimfacin__nb-blok">
        <div className="claimfacin__nb-kolom">
          <div className="claimfacin__nb-tumpuk">
            <TeksNB label={ALAMAT_NB.objectNo} nilai={p('ObjectNo')} />
            <TeksNB label={ALAMAT_NB.objectType} nilai={p('ObjectType')} />
            {p('ObjectType') === OBJECT_TYPE_LAINNYA_NB && (
              <TeksNB label={ALAMAT_NB.objectName} nilai={p('ObjectName')} />
            )}
          </div>
          <div className="claimfacin__nb-tumpuk">
            <CentangNB label={ALAMAT_NB.materialDamage} nilai={benarNB(p('IsMaterialDamage'))} />
            <CentangNB label={ALAMAT_NB.topRisk} nilai={benarNB(p('IsTopRisk'))} />
          </div>
        </div>
      </section>
      <section className="claimfacin__nb-blok">
        <h5 className="claimfacin__nb-subjudul">{ALAMAT_NB.judulRisk}</h5>
        <div className="claimfacin__nb-kolom">
          <div className="claimfacin__nb-tumpuk">
            <TeksNB label={ALAMAT_NB.type} nilai={p('RoadType')} />
            <TeksNB label={ALAMAT_NB.address} nilai={p('RoadName')} />
            <TeksNB label={ALAMAT_NB.buildingNo} nilai={p('BuildingNo')} />
            <TeksNB label={ALAMAT_NB.zipCode} nilai={p('RiskLocation.ASMZipCode')} />
            <TeksNB label={ALAMAT_NB.country} nilai={p('Country')} />
          </div>
          <div className="claimfacin__nb-tumpuk">
            <TeksNB label={ALAMAT_NB.territory} nilai={p('RiskLocation.ASMRW')} />
            <TeksNB label={ALAMAT_NB.city} nilai={p('RiskLocation.ASMCity')} />
            <TeksNB label={ALAMAT_NB.district} nilai={p('RiskLocation.ASMDistrict')} />
            <TeksNB label={ALAMAT_NB.province} nilai={p('Province')} />
            <TeksNB label={ALAMAT_NB.riskAddressId} nilai={p('AlmRiskID')} />
          </div>
        </div>
        <TeksNB label={ALAMAT_NB.riskLocation} nilai={p('RiskLocation.ASMAddress')} lebar />
      </section>
      <section className="claimfacin__nb-blok">
        <h5 className="claimfacin__nb-subjudul">{ALAMAT_NB.judulBangunan}</h5>
        <div className="claimfacin__nb-kolom">
          <div className="claimfacin__nb-tumpuk">
            <TeksNB label={ALAMAT_NB.numberOfFloor} nilai={bc('NumberOfFloor')} />
            <TeksNB label={ALAMAT_NB.roofType} nilai={labelKode(OPSI_ROOF_NB, bc('RoofType'))} />
            <TeksNB label={ALAMAT_NB.wallType} nilai={labelKode(OPSI_WALL_NB, bc('WallType'))} />
            <TeksNB label={ALAMAT_NB.floorType} nilai={labelKode(OPSI_FLOOR_NB, bc('FloorType'))} />
          </div>
          <div className="claimfacin__nb-tumpuk">
            <TeksNB label={ALAMAT_NB.partitionType} nilai={bc('PartitionType')} />
            <TeksNB label={ALAMAT_NB.supportWallType} nilai={bc('SupportWallType')} />
            <TeksNB label={ALAMAT_NB.otherType} nilai={bc('OthersType')} />
          </div>
        </div>
      </section>
    </div>
  )
}

/** Sub-tab Surrounding Risk: empat sisi + Other Description. */
function SekitarNB({ b }: { b: Baris }) {
  const s = (k: string) => v(b, 'Property.SurroundingRisk.' + k)
  return (
    <div className="claimfacin__nb-isi">
      <div className="claimfacin__nb-sisi">
        {SISI_NB.map((x) => (
          <section key={x.awalan} className="claimfacin__nb-blok">
            <h5 className="claimfacin__nb-subjudul">{x.judul}</h5>
            <TeksNB label={MEDAN_SISI_NB.occupation} nilai={s(x.awalan + 'Occupation')} />
            <div className="claimfacin__nb-sebaris">
              <TeksNB
                label={MEDAN_SISI_NB.construction}
                nilai={labelKode(OPSI_CONSTRUCTION_NB, s(x.awalan + 'Construction'))}
              />
              <TeksNB label={MEDAN_SISI_NB.distance} nilai={s(x.awalan + 'Distance')} />
            </div>
            <TeksNB label={MEDAN_SISI_NB.note} nilai={s(x.awalan + 'Note')} />
          </section>
        ))}
      </div>
      <div className="claimfacin__nb-sisi">
        <section className="claimfacin__nb-blok">
          <h5 className="claimfacin__nb-subjudul">{LAIN_SEKITAR_NB.judul}</h5>
          <div className="claimfacin__nb-sebaris">
            <TeksNB
              label={LAIN_SEKITAR_NB.ownership}
              nilai={labelKode(OPSI_OWNERSHIP_NB, v(b, 'Property.Ownership'))}
            />
            <TeksNB
              label={LAIN_SEKITAR_NB.housekeepingStatus}
              nilai={labelKode(OPSI_HOUSEKEEPING_NB, s('HousekeepingStatus'))}
            />
          </div>
          <div className="claimfacin__nb-sebaris">
            <TeksNB
              label={LAIN_SEKITAR_NB.floodAreaStatus}
              nilai={labelKode(OPSI_FLOOD_STATUS_NB, s('FloodAreaStatus'))}
            />
            {s('FloodAreaStatus') === FLOOD_STATUS_TAMPIL_NB && (
              <TeksNB label={LAIN_SEKITAR_NB.floodArea} nilai={labelKode(OPSI_FLOOD_AREA_NB, s('FloodArea'))} />
            )}
          </div>
          <TeksNB label={LAIN_SEKITAR_NB.housekeepingRemark} nilai={s('HousekeepingRemark')} />
        </section>
        <section className="claimfacin__nb-blok">
          <h5 className="claimfacin__nb-subjudul">{LAIN_SEKITAR_NB.judul}</h5>
          <p className="claimfacin__nb-faktor">{LAIN_SEKITAR_NB.faktor}</p>
          <CentangNB
            label={LAIN_SEKITAR_NB.productionProcess}
            nilai={benarNB(v(b, 'Property.IsProductionProcessFlag'))}
          />
          <CentangNB label={LAIN_SEKITAR_NB.hotWork} nilai={benarNB(v(b, 'Property.IsHotWorkProcessFlag'))} />
          <CentangNB label={LAIN_SEKITAR_NB.flammable} nilai={benarNB(v(b, 'Property.IsFlammableItemFlag'))} />
        </section>
      </div>
    </div>
  )
}

/** Sub-tab Object Item: grid PropertyItemList + form baris + grid Total (Σ TSI per mata uang). */
function ItemObjekNB({ h, n }: { h: Halaman; n: number }) {
  const items = daftarNB(h, itemNB(n))
  return (
    <div className="claimfacin__nb-isi">
      <GridNB
        kolom={kolomAngka(KOLOM_ITEM_NB, ['TSI Object Item'])}
        baris={items.map((b) => [
          v(b, 'ItemType'),
          labelKode(OPSI_CONDITION_NB, v(b, 'Condition')),
          v(b, 'PropertyYear'),
          v(b, 'Unit'),
          v(b, 'Currency'),
          angkaNB(v(b, 'TSIObjectItem')),
        ])}
        rinci={(i) => <FormItemNB b={items[i] ?? {}} />}
      />
      <GridNB
        kolom={kolomAngka(KOLOM_TOTAL_ITEM_NB, ['Total TSI'])}
        baris={totalTSIItemNB(items).map((t) => [t.currency, angkaNB(t.total)])}
      />
    </div>
  )
}

function FormItemNB({ b }: { b: Baris }) {
  const adjustable = benarNB(v(b, 'IsAdjustableFlag'))
  return (
    <div className="claimfacin__nb-isi">
      <div className="claimfacin__nb-sebaris">
        <TeksNB label={FORM_ITEM_NB.itemType} nilai={v(b, 'ItemType')} />
      </div>
      <TeksNB label={FORM_ITEM_NB.note} nilai={v(b, 'PropertiItemNote')} lebar />
      <div className="claimfacin__nb-sebaris">
        <TeksNB label={FORM_ITEM_NB.year} nilai={v(b, 'PropertyYear')} />
        <TeksNB label={FORM_ITEM_NB.unit} nilai={v(b, 'Unit')} />
        <TeksNB label={FORM_ITEM_NB.condition} nilai={labelKode(OPSI_CONDITION_NB, v(b, 'Condition'))} />
      </div>
      <div className="claimfacin__nb-sebaris">
        <TeksNB label={FORM_ITEM_NB.currency} nilai={v(b, 'Currency')} />
        <TeksNB label={FORM_ITEM_NB.tsi} nilai={angkaNB(v(b, 'TSIObjectItem'))} angka />
      </div>
      <div className="claimfacin__nb-sebaris">
        <TeksNB label={FORM_ITEM_NB.yearOfPlanting} nilai={v(b, 'Year')} />
        <TeksNB label={FORM_ITEM_NB.noOfTree} nilai={v(b, 'NoOfTree')} />
        <TeksNB label={FORM_ITEM_NB.areaHectar} nilai={v(b, 'AreaHectar')} />
      </div>
      <TeksNB label={`${FORM_ITEM_NB.remark} ${v(b, 'ItemType')}`} nilai={v(b, 'Remark')} lebar />
      <div className="claimfacin__nb-sebaris">
        <CentangNB label={FORM_ITEM_NB.adjustable} nilai={adjustable} />
        <TeksNB label={FORM_ITEM_NB.pctAdjust} nilai={adjustable ? v(b, 'PctAdjustOther') : v(b, 'PctAdjust2')} />
      </div>
    </div>
  )
}

/** Sub-tab Occupation (`Property.OccupationList`). */
function OkupasiNB({ h, n }: { h: Halaman; n: number }) {
  const rows = daftarNB(h, anakNB(LOKASI_NB, n, 'Property.OccupationList'))
  return (
    <GridNB
      kolom={KOLOM_OKUPASI_NB}
      baris={rows.map((b) => [v(b, 'OccupationId'), v(b, 'OccupationName'), v(b, 'TableOfLimit.Description')])}
      rinci={(i) => {
        const b = rows[i] ?? {}
        return (
          <div className="claimfacin__nb-tiga">
            <TeksNB label={FORM_OKUPASI_NB.occupationId} nilai={v(b, 'OccupationId')} />
            <TeksNB label={FORM_OKUPASI_NB.occupationName} nilai={v(b, 'OccupationName')} lebar />
            <TeksNB label={FORM_OKUPASI_NB.konstruksi} nilai={v(b, 'TableOfLimit.Description')} lebar />
          </div>
        )
      }}
    />
  )
}

/** FEAList di baris lokasi; Pega menaruh empat medan di bawah `DataFEA` (tabel datar tanpa awalan). */
const fea = (b: Baris, k: string) => b['DataFEA.' + k] ?? b[k] ?? ''

/** Sub-tab FEA (`LocationList(n).FEAList`). */
function FeaNB({ h, n }: { h: Halaman; n: number }) {
  const rows = daftarNB(h, anakNB(LOKASI_NB, n, 'FEAList'))
  return (
    <GridNB
      kolom={kolomAngka(KOLOM_FEA_NB, [
        'APAR',
        'Sprinkler',
        'Smoke Detector & Alarm',
        'Hydrant',
        'Private Truck Brigade',
      ])}
      baris={rows.map((b) => [
        v(b, 'APAR'),
        v(b, 'Sprinkler'),
        v(b, 'SmokeDetector'),
        v(b, 'Hydrant'),
        fea(b, 'PrivateTruckBrigade'),
        v(b, 'InfoFEA'),
      ])}
      rinci={(i) => {
        const b = rows[i] ?? {}
        return (
          <div className="claimfacin__nb-tiga">
            <TeksNB label={FORM_FEA_NB.apar} nilai={v(b, 'APAR')} />
            <TeksNB label={FORM_FEA_NB.sprinkler} nilai={v(b, 'Sprinkler')} />
            <TeksNB label={FORM_FEA_NB.smokeDetector} nilai={v(b, 'SmokeDetector')} />
            <TeksNB label={FORM_FEA_NB.hydrant} nilai={v(b, 'Hydrant')} />
            <TeksNB label={FORM_FEA_NB.privateTruckBrigade} nilai={fea(b, 'PrivateTruckBrigade')} />
            <TeksNB label={FORM_FEA_NB.privateFireBrigade} nilai={fea(b, 'PrivateFireBrigade')} />
            <TeksNB label={FORM_FEA_NB.teamSopSafety} nilai={fea(b, 'TeamSOPSafety')} />
            <TeksNB label={FORM_FEA_NB.teamSopRiskManagement} nilai={fea(b, 'TeamSOPRiskManagement')} />
            <TeksNB label={FORM_FEA_NB.info} nilai={v(b, 'InfoFEA')} lebar />
          </div>
        )
      }}
    />
  )
}

/** Sub-tab Loss Record (`Property.ListCauseOfLoss`) + blok LR baris lokasi. */
function KerugianNB({ h, n, b }: { h: Halaman; n: number; b: Baris }) {
  const rows = daftarNB(h, anakNB(LOKASI_NB, n, 'Property.ListCauseOfLoss'))
  const insured = h.nilai['OfferFacIn.QuotationData.InsuredName'] ?? ''
  return (
    <div className="claimfacin__nb-isi">
      <GridNB
        kolom={kolomAngka(KOLOM_KERUGIAN_NB, ['Total of Loss', 'Total Claim'])}
        baris={rows.map((r) => [
          tanggalNB(v(r, 'DateOfLoss')),
          v(r, 'CoinsData.CoinsName'),
          v(r, 'LossObject'),
          v(r, 'Currency'),
          angkaNB(v(r, 'Amount')),
          angkaNB(v(r, 'Claim')),
        ])}
        rinci={(i) => {
          const r = rows[i] ?? {}
          return (
            <div className="claimfacin__nb-dua">
              <TeksNB label={FORM_KERUGIAN_NB.insuredName} nilai={insured} />
              <TeksNB label={FORM_KERUGIAN_NB.dateOfLoss} nilai={tanggalNB(v(r, 'DateOfLoss'))} />
              <TeksNB label={FORM_KERUGIAN_NB.lossObject} nilai={v(r, 'LossObject')} />
              <TeksNB label={FORM_KERUGIAN_NB.currency} nilai={v(r, 'Currency')} />
              <TeksNB label={FORM_KERUGIAN_NB.claim} nilai={angkaNB(v(r, 'Claim'))} angka />
              <TeksNB label={FORM_KERUGIAN_NB.preventionOfLoss} nilai={v(r, 'PreventionOfLoss')} />
              <TeksNB label={FORM_KERUGIAN_NB.causeOfLoss} nilai={v(r, 'CauseOfLoss')} />
              <TeksNB label={FORM_KERUGIAN_NB.remarks} nilai={v(r, 'Remarks')} />
              <TeksNB label={FORM_KERUGIAN_NB.detail} nilai={v(r, 'Detail')} lebar />
            </div>
          )
        }}
      />
      <div className="claimfacin__nb-empat">
        {LOSS_RATIO_NB.map((x) => (
          <TeksNB key={x.kunci} label={x.label} nilai={angkaNB(v(b, x.kunci))} angka />
        ))}
      </div>
    </div>
  )
}

/** Sub-tab Loss Record Internal (`Property.ListCauseOfLossClaim`) - grid tanpa buka / aksi. */
function KlaimInternalNB({ h, n }: { h: Halaman; n: number }) {
  const rows = daftarNB(h, anakNB(LOKASI_NB, n, 'Property.ListCauseOfLossClaim'))
  return (
    <GridNB
      kolom={kolomAngka(KOLOM_KLAIM_INTERNAL_NB, [
        'Premium',
        'O/S Claim',
        'Acccepted Claim',
        'Incurred Claim',
        'Loss Ratio',
      ])}
      baris={rows.map((r) => [
        tahunNB(v(r, 'DateOfLoss')),
        tanggalNB(v(r, 'DateOfLoss')),
        v(r, 'LocationNo'),
        v(r, 'Location'),
        v(r, 'Currency'),
        angkaNB(v(r, 'Premium')),
        angkaNB(v(r, 'EstiamsiKlaim')),
        angkaNB(v(r, 'AkseptasiKlaim')),
        angkaNB(v(r, 'InccuredKlaim')),
        angkaNB(v(r, 'LossRatio')),
        v(r, 'Remark'),
      ])}
    />
  )
}

// ---------------------------------------------------------------- Coverage

/** Tab Coverage (`TabCoverage`): empat tingkat buka-tutup + ringkasan hijau. */
export function TabCoverageNB({ h }: { h: Halaman }) {
  const lokasi = daftarNB(h, LOKASI_NB)
  const r = ringkasanCoverageNB(h)
  return (
    <div className="claimfacin__nb-isi">
      <GridNB
        kolom={KOLOM_COV_OBJEK_NB}
        baris={lokasi.map((b) => [
          v(b, 'Property.ObjectNo'),
          v(b, 'Property.ObjectName'),
          v(b, 'Property.RiskLocation.ASMAddress'),
        ])}
        rinci={(i) => <CoverageLokasiNB h={h} n={i + 1} />}
      />
      <section className="claimfacin__nb-ringkasan">
        <GridNB
          kolom={kolomAngka(RINGKASAN_COV_ITEM_NB, ['Total TSI', 'Total Premium'])}
          baris={r.item.map((x) => [x.itemType, x.currency, angkaNB(x.tsi), angkaNB(x.premium)])}
        />
        <GridNB
          kolom={kolomAngka(RINGKASAN_COV_UANG_NB, ['Total TSI', 'Total Premium', 'Rate'])}
          baris={r.mataUang.map((x) => [x.currency, angkaNB(x.tsi), angkaNB(x.premium), angkaNB(x.rate)])}
        />
      </section>
    </div>
  )
}

function CoverageLokasiNB({ h, n }: { h: Halaman; n: number }) {
  const items = daftarNB(h, itemNB(n))
  return (
    <div className="claimfacin__nb-isi">
      <GridNB
        kolom={kolomAngka(KOLOM_COV_ITEM_NB, ['TSI Object Item', 'Total Gross Premium', '‰ Total Net Rate'])}
        baris={items.map((b, j) => [
          v(b, 'ItemType'),
          v(b, 'Currency'),
          angkaNB(v(b, 'TSIObjectItem')),
          angkaNB(premiItemNB(daftarNB(h, coverageNB(n, j + 1)), b)),
          angkaNB(v(b, 'TotalNetRate')),
        ])}
        rinci={(j) => <CoverageItemNB h={h} n={n} j={j + 1} item={items[j] ?? {}} />}
      />
      <GridNB
        kolom={kolomAngka(KOLOM_COV_TOTAL_NB, ['Total TSI', 'Total Premi', 'Total Gross Rate'])}
        baris={totalCoverageLokasiNB(h, n).map((t) => [
          t.currency,
          angkaNB(t.tsi),
          angkaNB(t.premium),
          angkaNB(t.rate),
        ])}
      />
    </div>
  )
}

function CoverageItemNB({ h, n, j, item }: { h: Halaman; n: number; j: number; item: Baris }) {
  const cov = daftarNB(h, coverageNB(n, j))
  return (
    <GridNB
      kolom={kolomAngka(KOLOM_COVERAGE_NB, ['‰ Standard Rate', 'Premi'])}
      baris={cov.map((c) => [
        <span key="k" className="claimfacin__nb-kode">
          {v(c, 'OLDID')}
        </span>,
        angkaNB(v(c, 'RateOJK')),
        angkaNB(v(c, 'Premium')),
      ])}
      rinci={(c) => <FormCoverageNB h={h} n={n} j={j} c={c + 1} b={cov[c] ?? {}} item={item} />}
    />
  )
}

/** Form coverage hanya-baca (`FormCoverage`): blok atas label kiri, dua kolom menurut Coverage Basis, grid Deductible. */
function FormCoverageNB({
  h,
  n,
  j,
  c,
  b,
  item,
}: {
  h: Halaman
  n: number
  j: number
  c: number
  b: Baris
  item: Baris
}) {
  const basis = v(b, 'CoverageBasis')
  const tsi = v(b, 'TSI') !== '' ? v(b, 'TSI') : v(item, 'TSIObjectItem')
  return (
    <div className="claimfacin__nb-isi claimfacin__nb-cov">
      <div className="claimfacin__nb-tumpuk">
        <TeksNB kiri label={FORM_COV_NB.coverageBasis} nilai={labelKode(OPSI_COVERAGE_BASIS_NB, basis)} />
        <TeksNB kiri label={FORM_COV_NB.coverage} nilai={v(b, 'OLDID')} />
        <TeksNB kiri label={FORM_COV_NB.accumulationCode} nilai={v(b, 'AccumulationCode') || AKUMULASI_KOSONG_NB} />
        <TeksNB
          kiri
          label={FORM_COV_NB.accumulationAddress}
          nilai={v(b, 'AccumulationDescription') || AKUMULASI_KOSONG_NB}
        />
        <TeksNB kiri label={FORM_COV_NB.conditions} nilai={v(b, 'Conditions')} />
      </div>
      <div className="claimfacin__nb-kolom claimfacin__nb-cov-kolom">
        <div className="claimfacin__nb-tumpuk">
          <RadioNB kiri label={FORM_COV_NB.day} pilihan={PILIHAN_DAY_COV_NB} nilai={v(b, 'Day')} />
          <TeksNB kiri label={FORM_COV_NB.tsi} nilai={angkaNB(tsi)} angka />
          <TeksNB kiri label={FORM_COV_NB.indemnity} nilai={v(b, 'Indemnity')} />
          <TeksNB kiri label={FORM_COV_NB.rate} nilai={angkaNB(v(b, 'Rate'))} angka />
          {basis === '2' && <TeksNB kiri label={FORM_COV_NB.firstLoss} nilai={angkaNB(v(b, 'FirstLoss'))} angka />}
          <TeksNB kiri label={FORM_COV_NB.discountPercentage} nilai={angkaNB(v(b, 'DiscountPercentage'))} angka />
          {basis === '2' && (
            <TeksNB kiri label={FORM_COV_NB.tsiLiability} nilai={angkaNB(v(b, 'TSILiability'))} angka />
          )}
          <TeksNB kiri label={FORM_COV_NB.netRate} nilai={angkaNB(v(b, 'NetRate'))} angka />
          <TeksNB kiri label={FORM_COV_NB.limitOfLiability} nilai={angkaNB(v(b, 'LimitofLiability'))} angka />
          <TeksNB kiri label={FORM_COV_NB.pctLol} nilai={angkaNB(v(b, 'PctLoL'))} angka />
        </div>
        <div className="claimfacin__nb-tumpuk">
          <TeksNB kiri label={FORM_COV_NB.proRate} nilai={angkaNB(v(b, 'ProRatePercent'))} angka />
          <TeksNB kiri label={FORM_COV_NB.unit} nilai={labelKode(OPSI_INDEMNITY_UNIT_NB, v(b, 'Unit') || '0')} />
          <TeksNB kiri label={FORM_COV_NB.indemnityPercentage} nilai={angkaNB(v(b, 'IndemnityPercentage'))} angka />
          {basis === '2' && <TeksNB kiri label={FORM_COV_NB.firstScale} nilai={angkaNB(v(b, 'FirstScale'))} angka />}
          {(basis === '4' || basis === '5') && (
            <TeksNB kiri label={FORM_COV_NB.sublimit} nilai={angkaNB(v(b, 'Sublimit'))} angka />
          )}
          <TeksNB kiri label={FORM_COV_NB.lostLimit} nilai={angkaNB(v(b, 'LostLimit'))} angka />
          {basis === '3' && <TeksNB kiri label={FORM_COV_NB.emlPml} nilai={angkaNB(v(b, 'EmlPml'))} angka />}
          <TeksNB kiri label={FORM_COV_NB.discount} nilai={angkaNB(v(b, 'Discount'))} angka />
          <TeksNB kiri label={FORM_COV_NB.premium} nilai={angkaNB(v(b, 'Premium'))} angka />
        </div>
      </div>
      <h5 className="claimfacin__nb-subjudul">{JUDUL_DEDUCTIBLE_NB}</h5>
      <DeductibleNB rows={daftarNB(h, deductibleNB(n, j, c))} />
    </div>
  )
}

const kondisiDeductible = (d: Baris) =>
  v(d, 'Condition') === KONDISI_LAIN_NB
    ? v(d, 'InputCondition')
    : labelKode(OPSI_KONDISI_DEDUCTIBLE_NB, v(d, 'Condition'))

/** Grid Deductible (`GridDeductible`) + form baris dua kolom. */
function DeductibleNB({ rows }: { rows: Baris[] }) {
  return (
    <GridNB
      kolom={kolomAngka(KOLOM_DEDUCTIBLE_NB, ['Pct Deductible', 'Amount'])}
      baris={rows.map((d) => [
        angkaNB(v(d, 'PctDeductible')),
        labelKode(OPSI_TYPE_DEDUCTIBLE_NB, v(d, 'TypeDeductible')),
        labelKode(OPSI_MINMAX_NB, v(d, 'MinMax')),
        v(d, 'Currency'),
        angkaNB(v(d, 'Amount')),
        kondisiDeductible(d),
        v(d, 'TimeExcess'),
      ])}
      rinci={(i) => {
        const d = rows[i] ?? {}
        const inAmount = v(d, 'TypeDeductible') === TYPE_IN_AMOUNT_NB
        return (
          <div className="claimfacin__nb-kolom">
            <div className="claimfacin__nb-tumpuk">
              <TeksNB
                label={FORM_DEDUCTIBLE_NB.typeDeductible}
                nilai={labelKode(OPSI_TYPE_DEDUCTIBLE_NB, v(d, 'TypeDeductible'))}
              />
              {!inAmount && (
                <TeksNB label={FORM_DEDUCTIBLE_NB.pctDeductible} nilai={angkaNB(v(d, 'PctDeductible'))} angka />
              )}
              {!inAmount && (
                <TeksNB label={FORM_DEDUCTIBLE_NB.minMax} nilai={labelKode(OPSI_MINMAX_NB, v(d, 'MinMax'))} />
              )}
              <TeksNB label={FORM_DEDUCTIBLE_NB.currency} nilai={v(d, 'Currency')} />
              {v(d, 'MinMax') === MINMAX_OR_NB && (
                <>
                  <TeksNB
                    label={FORM_DEDUCTIBLE_NB.typeDeductible}
                    nilai={labelKode(OPSI_TYPE_DEDUCTIBLE2_NB, v(d, 'TypeDeductible2'))}
                  />
                  <TeksNB label={FORM_DEDUCTIBLE_NB.pctDeductible} nilai={v(d, 'PctDeductible2')} />
                </>
              )}
            </div>
            <div className="claimfacin__nb-tumpuk">
              <TeksNB
                label={FORM_DEDUCTIBLE_NB.condition}
                nilai={labelKode(OPSI_KONDISI_DEDUCTIBLE_NB, v(d, 'Condition'))}
              />
              <TeksNB label={FORM_DEDUCTIBLE_NB.amount} nilai={angkaNB(v(d, 'Amount'))} angka />
              {v(d, 'Condition') === KONDISI_LAIN_NB && (
                <TeksNB label={FORM_DEDUCTIBLE_NB.condition} nilai={v(d, 'InputCondition')} />
              )}
              <TeksNB label={FORM_DEDUCTIBLE_NB.timeExcess} nilai={v(d, 'TimeExcess')} />
            </div>
          </div>
        )
      }}
    />
  )
}

// ---------------------------------------------------------------- Spreading

const barisTotal = (t: TotalSpreadingNB) => [
  t.currency,
  t.treatyName,
  angkaNB(t.sharePercentage),
  angkaNB(t.claimSpreaded),
  angkaNB(t.tsiSpreaded),
  angkaNB(t.claimEstimation),
  angkaNB(t.premiumSpreaded),
]

const ANGKA_TOTAL = [
  '(%) Share',
  'Total TSI Top Risk Spreaded',
  'Total TSI Spreaded',
  'Total Limit of Liability',
  'Total First Loss',
  'Total Premium Spreaded',
]

/** Tab Spreading (`TabSpreading`): % Share RNM, lokasi -> item (+ total kuning) -> coverage -> spreading, ringkasan hijau. */
export function TabSpreadingNB({ h }: { h: Halaman }) {
  const lokasi = daftarNB(h, LOKASI_NB)
  const s = hitungSpreadingNB(h)
  return (
    <div className="claimfacin__nb-isi">
      <div className="claimfacin__nb-blok claimfacin__nb-share">
        <TeksNB kiri label={SPREADING_NB.percentShare} nilai={angkaNB(nilaiNB(h, 'PercentShare'))} angka />
      </div>
      <GridNB
        kolom={SPREADING_NB.kolomLokasi}
        baris={lokasi.map((b) => [
          v(b, 'Property.ObjectNo'),
          v(b, 'Property.ObjectName'),
          v(b, 'Property.RiskLocation.ASMAddress'),
        ])}
        rinci={(i) => <SpreadingLokasiNB h={h} n={i + 1} total={s.lokasi[i] ?? []} />}
      />
      <section className="claimfacin__nb-ringkasan">
        <GridNB
          kolom={kolomAngka(SPREADING_NB.kolomRingkasanTreaty, ANGKA_TOTAL)}
          baris={s.treaty.map((t) => {
            const r = barisTotal(t)
            return [...r.slice(0, 6), angkaNB(t.claimAmountIDR), r[6]]
          })}
        />
        <GridNB
          kolom={kolomAngka(SPREADING_NB.kolomRingkasanMataUang, ANGKA_TOTAL)}
          baris={s.mataUang.map((t) => [
            t.currency,
            angkaNB(t.tsiTopRisk),
            angkaNB(t.tsi),
            angkaNB(t.lol),
            angkaNB(t.premium),
          ])}
        />
      </section>
    </div>
  )
}

function SpreadingLokasiNB({ h, n, total }: { h: Halaman; n: number; total: TotalSpreadingNB[] }) {
  const items = daftarNB(h, itemNB(n))
  return (
    <div className="claimfacin__nb-isi">
      <GridNB
        kolom={kolomAngka(SPREADING_NB.kolomItem, ['TSI Object Item', 'Total Gross Premium', 'Total Premium RNM'])}
        baris={items.map((b, j) => {
          const cov = daftarNB(h, coverageNB(n, j + 1))
          return [
            v(b, 'ItemType'),
            v(b, 'Currency'),
            angkaNB(v(b, 'TSIObjectItem')),
            angkaNB(premiItemNB(cov, b)),
            angkaNB(premiRNMItemNB(cov)),
          ]
        })}
        rinci={(j) => <SpreadingItemNB h={h} n={n} j={j + 1} />}
      />
      <section className="claimfacin__nb-total">
        <GridNB kolom={kolomAngka(SPREADING_NB.kolomTotal, ANGKA_TOTAL)} baris={total.map(barisTotal)} />
      </section>
    </div>
  )
}

function SpreadingItemNB({ h, n, j }: { h: Halaman; n: number; j: number }) {
  const cov = daftarNB(h, coverageNB(n, j))
  const D = SPREADING_NB.detailCoverage
  return (
    <GridNB
      kolom={kolomAngka(SPREADING_NB.kolomCoverage, ['‰ Standard Rate', 'Gross Premium', 'Premium RNM'])}
      baris={cov.map((c) => [
        <span key="k" className="claimfacin__nb-kode">
          {v(c, 'OLDID')}
        </span>,
        angkaNB(v(c, 'Rate')),
        angkaNB(v(c, 'Premium')),
        angkaNB(v(c, 'PremiNusantaraRe')),
      ])}
      rinci={(m) => {
        const c = cov[m] ?? {}
        const rows = daftarNB(h, spreadingNB(n, j, m + 1))
        return (
          <div className="claimfacin__nb-isi">
            <dl className="claimfacin__nb-angka-lima">
              {(
                [
                  [D.tsi, 'TSI'],
                  [D.tsiLiability, 'TSILiability'],
                  [D.tsiNusantaraRe, 'TSINusantaraRe'],
                  [D.premium, 'Premium'],
                  [D.premiNusantaraRe, 'PremiNusantaraRe'],
                ] as const
              ).map(([l, k]) => (
                <div key={k}>
                  <dt>{l}</dt>
                  <dd>{angkaNB(v(c, k))}</dd>
                </div>
              ))}
            </dl>
            <GridNB
              kosong={TEKS_NB.spreadingKosong}
              kolom={kolomAngka(SPREADING_NB.kolomSpreading, [
                '% Share',
                'TSI 100% Spreaded',
                'TSI Spreaded (RNM)',
                'Limit Of Liability',
                'Premium Spreaded',
              ])}
              baris={rows.map((r) => [
                v(r, 'TreatyName') || v(r, 'TreatyType'),
                angkaNB(v(r, 'SharePercentage')),
                angkaNB(v(r, 'TSIGrossSpreaded')),
                angkaNB(v(r, 'TSISpreaded')),
                angkaNB(v(r, 'ClaimEstimation')),
                angkaNB(v(r, 'PremiumSpreaded')),
              ])}
            />
          </div>
        )
      }}
    />
  )
}

// ---------------------------------------------------------------- Inw Fac Cedant Panels

/** Tab Inw Fac Cedant Panels (`TabCedant`): blok atas label kiri + grid CedingCedantList. */
export function TabCedantNB({ h }: { h: Halaman }) {
  const C = CEDANT_NB_LABEL
  const rows = daftarNB(h, CEDANT_NB)
  return (
    <div className="claimfacin__nb-isi">
      <div className="claimfacin__nb-blok claimfacin__nb-cedant">
        <TeksNB kiri label={C.sob} nilai={h.nilai['OfferFacIn.QuotationData.SobName'] ?? ''} />
        <RadioNB
          kiri
          label={C.shareCedantType}
          pilihan={OPSI_SHARE_CEDANT_NB.map((o) => o.value)}
          opsiLabel={(p) => labelKode(OPSI_SHARE_CEDANT_NB, p)}
          nilai={nilaiNB(h, 'ShareCedantType')}
        />
        <TeksNB kiri label={C.percentShare} nilai={angkaNB(nilaiNB(h, 'PercentShare'))} angka />
        <TeksNB kiri label={C.totalTsi} nilai={angkaNB(nilaiNB(h, 'TotalTSINusaReSpreading'))} angka />
        <TeksNB kiri label={C.totalPremi} nilai={angkaNB(nilaiNB(h, 'TotalPremiNusaRe'))} angka />
      </div>
      <div className="claimfacin__nb-cedant-grid">
        <GridNB
          kolom={[C.kolomCeding, uang(C.kolomShare)]}
          baris={rows.map((r) => [v(r, 'CedingCoName'), angkaNB(v(r, 'ShareCeding'))])}
        />
      </div>
    </div>
  )
}
