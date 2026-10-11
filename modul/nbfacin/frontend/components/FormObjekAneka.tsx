// Form satu objek Aneka (tiket 73 P1) - port `Section\InputDtlObjectAneka_FacIn.xml` (Data-Aneka, 9 kontainer tingkat atas).
//
// Kontainer yang tampil = `tampilanAneka.kontainer` dari backend (visible-when kontainer dievaluasi lewat When korpus):
// C1 L361 (Liability, Golf, Bonding, …) · C2 L9626 (Marine Hull) · C3 L17725 (All Risk / Burglary) · C4 L19831 dan C5 L23038
// (IsObjectWithQuantityYear / HE) · C6 L24492 (HE) · C7 L27728 (Aviation Hull) · C8 L32052 (Growing Trees / Yield Shortfall) ·
// C9 L34618 (Electronic Equipment). Urutan kontainer = urutan Pega.
//
// Keputusan agent (tiket 73 P1):
// - A-1 Satu properti yang muncul di lebih dari satu kontainer (mis. Boiler = C1 + C4 + C5: Object Name, Currency, TSI,
//   Section) tampil SEKALI, di kontainer pertama - di Pega keduanya kotak berbeda yang terikat ke properti yang sama.
// - A-2 Sel ber-`pyVisible=ALWAYS` dengan sisa kondisi (`IsObjectLimiliabilityAneka` / `IsObjectSectionAneka`) selalu tampil
//   (K73-5).
// - A-3 `Read-only` TANPA syarat = SELALU baca-saja `[terverifikasi]` (tangkapan layar Pega work owner 09-10-2026, form coverage:
//   TSI L4069 tampil sebagai teks). Object Description HE (L24736) baca-saja - isinya dari SetObjectDescHE_DT. Vessel Name
//   (L10329) juga Read-only tanpa syarat (di Pega data kapal masuk lewat Upload CSV `InsertUploadAneka_act`); keputusan work owner
//   09-10-2026 K73-16 "HARUSNYA BISA DUA-DUANYA": Vessel Name tetap dapat dipilih dari daftar kapal DI FORM, dan Upload CSV Aneka
//   menyusul (P6). Medan kapal lain tetap tampil-saja (`pyReadOnlyCondition=Always`; terisi dari pilihan kapal / CSV).
// - A-4 Medan C2 yang selalu baca-saja (`pyReadOnlyCondition=Always`): Production Year, Place of Manufacture, Flag, seluruh
//   `.MarineHull.*` - diisi lookup kapal; Type of Vessel / Length / Breadth / Depth / Place of Manufacture tidak di-SET lookup
//   sehingga kosong (sama dengan Pega).
// - A-5 Tanda wajib (Object Description L606, Draft Wording ID L3481, Production Year L11261) = penanda saja; penolakan dari
//   backend (proteksi 422), pola L-2 tiket 40.
// - Object Name / Merk HE dan Draft Wording ID = teks bebas (K73-10 / K73-11) sampai DDL sumbernya ada.

import { useEffect, useRef, useState } from 'react'

import { Area, Field, Pilih, type Opsi } from '../../../../inti/frontend/components/ui/dasar'
import { cariKapal, daftarTipeKendaraan, type BarisKapal, type BarisTipeKendaraan, type ObjekAneka, type TampilanAneka } from '../api'
import { FORM_MEMORANDUM as M, FORM_OBJEK_ANEKA as F, OPSI_SECTION_ANEKA, TEKS_ANEKA, TEKS_FORM_OPPORTUNITY } from '../labels'
import { terapkanKapal, ubahHe } from './aneka'
import IsianSaran from './IsianSaran'
import IsianUang from './IsianUang'
import TanggalDMY from './TanggalDMY'

/** Jenis isian satu medan. */
type Jenis = 'teks' | 'area' | 'angka' | 'uang' | 'tanggal' | 'mataUang' | 'section' | 'tampil' | 'kapal' | 'heTeks' | 'heTipe'

/** Sakelar `tampilanAneka` yang menjadi syarat tampil satu medan (visible OTHER). */
type Syarat = 'bonding' | 'engineering' | 'jumlahSatuan'

interface Medan {
  /** Jalur properti di ObjekAneka; anak = "induk.medan" (mis. "marineHull.nrtrt"). */
  jalur: string
  label: string
  jenis: Jenis
  wajib?: boolean
  syarat?: Syarat
}

const MH = F.marineHull
const AH = F.aviationHull

/** Medan per kontainer, urut sel Pega (nomor baris sel `InputDtlObjectAneka_FacIn.xml`). */
export const KONTAINER: Record<string, Medan[]> = {
  C1: [
    { jalur: 'objectName', label: F.objectDescription, jenis: 'area', wajib: true }, // L606
    { jalur: 'section', label: F.section, jenis: 'section' }, // L920
    { jalur: 'noKontrak', label: F.noKontrak, jenis: 'teks' }, // L1187
    { jalur: 'startDate', label: F.startDate, jenis: 'tanggal' }, // L1438
    { jalur: 'endDate', label: F.endDate, jenis: 'tanggal' }, // L1732
    { jalur: 'noKtp', label: F.noKtp, jenis: 'teks' }, // L2027
    { jalur: 'noNpwp', label: F.noNpwp, jenis: 'teks' }, // L2278
    { jalur: 'currency', label: F.currency, jenis: 'mataUang' }, // L2529
    { jalur: 'tsi', label: F.tsi, jenis: 'uang' }, // L2995
    { jalur: 'draftWordingId', label: F.draftWordingId, jenis: 'teks', wajib: true, syarat: 'bonding' }, // L3481
    { jalur: 'maintenancePeriod', label: F.maintenancePeriod, jenis: 'angka', syarat: 'engineering' }, // L3778
  ],
  // C1 lanjutan (blok Briguna L4118, Others L8326, Memorandum L8750) dirender terpisah - lihat FormObjekAneka.
  C2: [
    { jalur: 'objectName', label: F.vesselName, jenis: 'kapal' }, // L10329
    { jalur: 'currency', label: F.currency, jenis: 'mataUang' }, // L10801
    { jalur: 'year', label: F.productionYear, jenis: 'tampil', wajib: true }, // L11261
    { jalur: 'tsi', label: F.vesselPrice, jenis: 'uang' }, // L12303
    { jalur: 'builtIn', label: F.builtIn, jenis: 'tampil' }, // L12833
    { jalur: 'flag', label: F.flag, jenis: 'tampil' }, // L13057
    // Blok `.MarineHull` = tiga kolom layout Pega (L14097 / L15291 / L16433: Type of Vessel · NRT · Length | Classification ·
    // GRT · Breadth | Construction · DWT · Depth), disusun per BARIS supaya kisi 3 kolom menampilkannya sama dengan layar Pega
    // (tangkapan layar work owner 09-10-2026).
    { jalur: 'marineHull.typeOfVessel', label: MH.typeOfVessel, jenis: 'tampil' }, // L14112
    { jalur: 'marineHull.classification', label: MH.classification, jenis: 'tampil' }, // L15306
    { jalur: 'marineHull.construction', label: MH.construction, jenis: 'tampil' }, // L16448
    { jalur: 'marineHull.nrtrt', label: MH.nrtrt, jenis: 'tampil' }, // L14355
    { jalur: 'marineHull.grtrt', label: MH.grtrt, jenis: 'tampil' }, // L15548
    { jalur: 'marineHull.dwtrt', label: MH.dwtrt, jenis: 'tampil' }, // L16691
    { jalur: 'marineHull.lengthm', label: MH.lengthm, jenis: 'tampil' }, // L14598
    { jalur: 'marineHull.breadthm', label: MH.breadthm, jenis: 'tampil' }, // L15741
    { jalur: 'marineHull.depthm', label: MH.depthm, jenis: 'tampil' }, // L16884
  ],
  C3: [
    { jalur: 'objectName', label: F.objectDescription, jenis: 'area' }, // L17969
    { jalur: 'currency', label: F.currency, jenis: 'mataUang' }, // L18231
    { jalur: 'tsi', label: F.tsi, jenis: 'uang' }, // L18697
    { jalur: 'quantity', label: F.quantity, jenis: 'teks' }, // L19178
    { jalur: 'invoiceNo', label: F.invoiceNo, jenis: 'teks' }, // L19368
  ],
  C4: [
    { jalur: 'objectName', label: F.objectName, jenis: 'area' }, // L20077
    { jalur: 'section', label: F.section, jenis: 'section' }, // L20313
    { jalur: 'currency', label: F.currency, jenis: 'mataUang' }, // L20579
    { jalur: 'tsi', label: F.tsi, jenis: 'uang' }, // L21046
    { jalur: 'vehicleHe.chassisNumber', label: F.chassisNumber, jenis: 'teks' }, // L21528
    { jalur: 'vehicleHe.engineNumber', label: F.engineNumber, jenis: 'teks' }, // L21767
    { jalur: 'serialNo', label: F.serialNo, jenis: 'teks' }, // L22198
    { jalur: 'maintenancePeriod', label: F.maintenancePeriod, jenis: 'angka', syarat: 'engineering' }, // L22560
  ],
  C5: [
    { jalur: 'year', label: F.year, jenis: 'teks' }, // L23527
    { jalur: 'quantity', label: F.quantity, jenis: 'teks', syarat: 'jumlahSatuan' }, // L23692 `BusinessCode!=20`
    { jalur: 'unitQuantity', label: F.unit, jenis: 'teks', syarat: 'jumlahSatuan' }, // L23884
  ],
  C6: [
    { jalur: 'objectName', label: F.objectDescription, jenis: 'tampil' }, // L24736 Read-only tanpa syarat (A-3) - diisi SetObjectDescHE_DT
    { jalur: 'vehicleHe.objectNameHe', label: F.objectNameHe, jenis: 'heTeks' }, // L24970 (K73-10)
    { jalur: 'currency', label: F.currency, jenis: 'mataUang' }, // L25354
    { jalur: 'tsi', label: F.tsi, jenis: 'uang' }, // L25820
    { jalur: 'vehicleHe.chassisNumber', label: F.chassisNumber, jenis: 'teks' }, // L26301
    { jalur: 'vehicleHe.engineNumber', label: F.engineNumber, jenis: 'teks' }, // L26540
    { jalur: 'vehicleHe.brandName', label: F.merk, jenis: 'heTeks' }, // L26779 (K73-10)
    { jalur: 'vehicleHe.typeName', label: F.type, jenis: 'heTipe' }, // L27162 `BrowseType_RD`
  ],
  C7: [
    { jalur: 'objectName', label: F.modelTipe, jenis: 'teks' }, // L28209
    { jalur: 'year', label: F.tahun, jenis: 'teks' }, // L28421
    { jalur: 'geographicalLimits', label: F.batasGeografis, jenis: 'teks' }, // L28586
    { jalur: 'aviationHull.registrationMarks', label: AH.registrationMarks, jenis: 'teks' }, // L29500
    { jalur: 'aviationHull.maxPassengers', label: AH.maxPassengers, jenis: 'teks' },
    { jalur: 'aviationHull.pilots', label: AH.pilots, jenis: 'teks' },
    { jalur: 'aviationHull.standardUses', label: AH.standardUses, jenis: 'teks' },
    { jalur: 'aviationHull.specialUses', label: AH.specialUses, jenis: 'teks' },
    { jalur: 'aviationHull.specialRentalUses', label: AH.specialRentalUses, jenis: 'teks' }, // L31172
  ],
  C8: [
    { jalur: 'objectName', label: F.objectName, jenis: 'area' }, // L32297
    { jalur: 'currency', label: F.currency, jenis: 'mataUang' }, // L32782
    { jalur: 'tsi', label: F.tsi, jenis: 'uang' }, // L33240
    { jalur: 'year', label: F.yearOfPlanting, jenis: 'teks' }, // L33712
    { jalur: 'noOfTree', label: F.noOfTree, jenis: 'angka' }, // L33889
    { jalur: 'areaHectar', label: F.areaHectar, jenis: 'angka' }, // L34066
  ],
  C9: [
    { jalur: 'objectName', label: F.objectName, jenis: 'area' }, // L34864
    { jalur: 'section', label: F.section, jenis: 'section' }, // L35100
    { jalur: 'quantity', label: F.quantity, jenis: 'teks' }, // L35338
    { jalur: 'currency', label: F.currency, jenis: 'mataUang' }, // L35528
    { jalur: 'tsi', label: F.tsi, jenis: 'uang' }, // L35997
    { jalur: 'year', label: F.year, jenis: 'teks' }, // L36554
    { jalur: 'limitIdemnityDay', label: F.limitIdemnityDay, jenis: 'angka' }, // L36988
    { jalur: 'limitIdemnityMonth', label: F.limitIdemnityMonth, jenis: 'angka' }, // L37151
  ],
}

/** Blok `IsCreditBriguna` (C1 L4118 / L4417, sel L4932-7687). */
const MEDAN_BRIGUNA: Medan[] = [
  { jalur: 'briguna.loanAccountNumber', label: F.briguna.loanAccountNumber, jenis: 'teks' },
  { jalur: 'briguna.loanName', label: F.briguna.loanName, jenis: 'teks' },
  { jalur: 'briguna.job', label: F.briguna.job, jenis: 'teks' },
  { jalur: 'briguna.depositAccountNo', label: F.briguna.depositAccountNo, jenis: 'teks' },
  { jalur: 'briguna.idCardNo', label: F.briguna.idCardNo, jenis: 'teks' },
  { jalur: 'briguna.institution', label: F.briguna.institution, jenis: 'teks' },
  { jalur: 'briguna.loanType', label: F.briguna.loanType, jenis: 'teks' },
  { jalur: 'briguna.dateOfBirth', label: F.briguna.dateOfBirth, jenis: 'tanggal' },
  { jalur: 'briguna.salary', label: F.briguna.salary, jenis: 'uang' },
]

/** Sub-form `InputDtlMemorandumAccept` (C1 L8750, ALWAYS; isinya pxTextInput L1068-2352, pxTextArea L2536-3633). */
const MEDAN_MEMORANDUM: Medan[] = [
  { jalur: 'memorandum.applicationNumber', label: M.applicationNumber, jenis: 'teks' },
  { jalur: 'memorandum.accumulation', label: M.accumulation, jenis: 'teks' },
  { jalur: 'memorandum.collateral', label: M.collateral, jenis: 'teks' },
  { jalur: 'memorandum.insuredSiup', label: M.insuredSiup, jenis: 'teks' },
  { jalur: 'memorandum.keyPerson', label: M.keyPerson, jenis: 'teks' },
  { jalur: 'memorandum.originalRate', label: M.originalRate, jenis: 'teks' },
  { jalur: 'memorandum.principalsOutGo', label: M.principalsOutGo, jenis: 'teks' },
  { jalur: 'memorandum.typeOfWording', label: M.typeOfWording, jenis: 'teks' },
  { jalur: 'memorandum.buildingStatus', label: M.buildingStatus, jenis: 'area' },
  { jalur: 'memorandum.supportingDocument', label: M.supportingDocument, jenis: 'area' },
  { jalur: 'memorandum.principalsData', label: M.principalsData, jenis: 'area' },
  { jalur: 'memorandum.conclusionsProposals', label: M.conclusionsProposals, jenis: 'area' },
  { jalur: 'memorandum.principalsBackground', label: M.principalsBackground, jenis: 'area' },
  { jalur: 'memorandum.summaryFinancial', label: M.summaryFinancial, jenis: 'area' },
  { jalur: 'memorandum.principalsGuarantee', label: M.principalsGuarantee, jenis: 'area' },
]

/** Kartu pengelompokan tampilan (work owner 09-10-2026 "RAPIKAN BIAR USER FRIENDLY"; judul `TEKS_ANEKA.grup`, bukan dari Pega). */
export type Grup = keyof typeof TEKS_ANEKA.grup

/** Kartu lebar (bertumpuk) dan kartu kecil (berjajar satu baris). */
const GRUP_BESAR: Grup[] = ['objek', 'kapal', 'kendaraan', 'pesawat']
const GRUP_KECIL: Grup[] = ['nilai', 'identitas', 'periode'] // urutan work owner 09-10-2026: Sum Insured, Identity, Period

/** Kartu satu medan: medan kapal tampil-saja (C2), anak `.VehicleHE` / `.AviationHull`, periode, nilai, identitas; sisanya Object. */
export function grupDari(m: Medan): Grup {
  const [induk] = m.jalur.split('.')
  if (induk === 'marineHull' || (m.jenis === 'tampil' && m.jalur !== 'objectName')) return 'kapal'
  if (induk === 'vehicleHe') return 'kendaraan'
  if (induk === 'aviationHull') return 'pesawat'
  if (m.jalur === 'startDate' || m.jalur === 'endDate') return 'periode'
  if (m.jalur === 'currency' || m.jalur === 'tsi') return 'nilai'
  if (m.jalur === 'noKtp' || m.jalur === 'noNpwp') return 'identitas'
  return 'objek'
}

/**
 * Medan per kartu, urutan dalam kartu = urutan sel Pega, kecuali Object Name / Description yang selalu pertama di kartu Object
 * (HE: kontainer C5 Year / Quantity / Unit berada sebelum C6 di Pega).
 */
export function kelompokkan(medan: Medan[]): Record<Grup, Medan[]> {
  const hasil = Object.fromEntries(Object.keys(TEKS_ANEKA.grup).map((g) => [g, [] as Medan[]])) as Record<Grup, Medan[]>
  for (const m of medan) hasil[grupDari(m)].push(m)
  hasil.objek.sort((a, b) => Number(b.jalur === 'objectName') - Number(a.jalur === 'objectName'))
  return hasil
}

/** Medan unik yang tampil, urut kontainer (A-1); syarat OTHER dievaluasi dari `tampilanAneka`. */
export function medanTampil(t: TampilanAneka): Medan[] {
  const sudah = new Set<string>()
  const hasil: Medan[] = []
  for (const k of ['C1', 'C2', 'C3', 'C4', 'C5', 'C6', 'C7', 'C8', 'C9']) {
    if (!t.kontainer.includes(k)) continue
    for (const m of KONTAINER[k] ?? []) {
      if (m.syarat !== undefined && !t[m.syarat]) continue
      if (sudah.has(m.jalur)) continue
      sudah.add(m.jalur)
      hasil.push(m)
    }
  }
  return hasil
}

/** Baca satu jalur ("induk.medan" atau "medan"). */
export function baca(o: ObjekAneka, jalur: string): string {
  const [a = '', b] = jalur.split('.')
  const v = b === undefined ? (o as unknown as Record<string, unknown>)[a] : ((o as unknown as Record<string, Record<string, unknown>>)[a] ?? {})[b]
  return typeof v === 'string' ? v : ''
}

/** Tulis satu jalur - salinan baru. */
export function tulis(o: ObjekAneka, jalur: string, nilai: string): ObjekAneka {
  const [a = '', b] = jalur.split('.')
  if (b === undefined) return { ...o, [a]: nilai }
  const anak = (o as unknown as Record<string, Record<string, string>>)[a]
  return { ...o, [a]: { ...anak, [b]: nilai } }
}

function Tampil({ label, nilai, wajib }: { label: string; nilai: string; wajib?: boolean }) {
  return (
    <div className="field">
      <span className="field__label">
        {label}
        {wajib && <span className="field__req">*</span>}
      </span>
      <div className="nbf-inward__teks">{nilai}</div>
    </div>
  )
}

const OPSI_SECTION: Opsi[] = OPSI_SECTION_ANEKA.map((o) => ({ value: o.value, label: o.label }))

/** Daftar + nilai tersimpan yang tidak dikenal daftar (data lama), supaya tidak hilang diam-diam saat Save. */
function denganTersimpan(daftar: Opsi[], nilai: string): Opsi[] {
  return nilai === '' || daftar.some((o) => o.value === nilai) ? daftar : [...daftar, { value: nilai, label: nilai }]
}

export default function FormObjekAneka({
  o,
  ubah,
  tampilan,
  mataUang,
}: {
  o: ObjekAneka
  ubah: (o: ObjekAneka) => void
  tampilan: TampilanAneka
  /** Pilihan Currency (`GET /api/nbfacin/mata-uang`, RD `BrowseCurrency_RD`). */
  mataUang: string[]
}) {
  // Hasil pencarian kapal terakhir: pilih saran → baris kapal lengkap (IsianSaran hanya membawa id / label).
  const kapal = useRef<BarisKapal[]>([])
  const [tipe, setTipe] = useState<BarisTipeKendaraan[] | null>(null)
  const pakaiTipe = tampilan.kontainer.includes('C6')

  useEffect(() => {
    if (!pakaiTipe) return
    let batal = false
    daftarTipeKendaraan().then(
      (h) => {
        if (!batal) setTipe(h.baris)
      },
      () => {
        if (!batal) setTipe([])
      },
    )
    return () => {
      batal = true
    }
  }, [pakaiTipe])

  const opsiUang: Opsi[] = mataUang.map((m) => ({ value: m, label: m }))

  function isian(m: Medan) {
    const v = baca(o, m.jalur)
    const set = (nilai: string) => ubah(tulis(o, m.jalur, nilai))
    switch (m.jenis) {
      case 'area':
        return <Area key={m.jalur} label={m.label} value={v} onChange={set} baris={3} required={m.wajib} />
      case 'angka':
        return <Field key={m.jalur} label={m.label} type="number" value={v} onChange={set} />
      case 'uang':
        return <IsianUang key={m.jalur} label={m.label} value={v} onChange={set} />
      case 'tanggal':
        return (
          <TanggalDMY
            key={m.jalur}
            label={m.label}
            value={v}
            onChange={set}
            labelKalender={TEKS_FORM_OPPORTUNITY.kalender}
            pesanFormat={TEKS_FORM_OPPORTUNITY.formatTanggal}
          />
        )
      case 'mataUang':
        return <Pilih key={m.jalur} label={m.label} value={v} onChange={set} opsi={denganTersimpan(opsiUang, v)} />
      case 'section':
        return <Pilih key={m.jalur} label={m.label} value={v} onChange={set} opsi={denganTersimpan(OPSI_SECTION, v)} />
      case 'tampil':
        return <Tampil key={m.jalur} label={m.label} nilai={v} wajib={m.wajib} />
      case 'kapal':
        return (
          <IsianSaran
            key={m.jalur}
            label={m.label}
            value={v}
            onKetik={set}
            muatSaatFokus
            cari={(q) =>
              cariKapal(q).then((h) => {
                kapal.current = h.baris
                return { baris: h.baris.map((k) => ({ id: k.id, label: k.nama })) }
              })
            }
            onPilih={(s) => {
              const k = kapal.current.find((x) => x.id === s.id)
              if (k) ubah(terapkanKapal(o, k))
            }}
          />
        )
      case 'heTeks': {
        // Object Name HE / Merk: teks bebas (K73-10); ID lookup dikosongkan karena teks tidak lagi dari daftar.
        const anak = m.jalur.split('.')[1] as 'objectNameHe' | 'brandName'
        const idnya = anak === 'objectNameHe' ? 'objectIdHe' : 'brand'
        return <Field key={m.jalur} label={m.label} value={v} onChange={(nilai) => ubah(ubahHe(o, { [anak]: nilai, [idnya]: '' }))} />
      }
      case 'heTipe':
        return (
          <IsianSaran
            key={m.jalur}
            label={m.label}
            value={v}
            muatSaatFokus
            onKetik={(nilai) => ubah(ubahHe(o, { typeName: nilai, type: '' }))}
            cari={(q) => {
              const kata = q.trim().toUpperCase()
              const baris = (tipe ?? []).filter((t) => kata === '' || t.typeName.toUpperCase().includes(kata))
              return Promise.resolve({ baris: baris.map((t) => ({ id: t.id, label: t.typeName })) })
            }}
            onPilih={(s) => ubah(ubahHe(o, { typeName: s.label, type: s.id }))}
          />
        )
      default:
        return <Field key={m.jalur} label={m.label} value={v} onChange={set} required={m.wajib} />
    }
  }

  if (tampilan.kontainer.length === 0) return <div className="alert alert--info">{TEKS_ANEKA.tanpaKontainer}</div>

  const medan = medanTampil(tampilan)
  const c1 = tampilan.kontainer.includes('C1')
  const kelompok = kelompokkan(medan)
  const kartu = (g: Grup, isi: Medan[]) =>
    isi.length > 0 && (
      <section key={g} className={`nbf-alamat__blok nbf-aneka-form__kartu nbf-aneka-form__kartu--${g}`}>
        <h5 className="nbf-objek__judul">{TEKS_ANEKA.grup[g]}</h5>
        <div className="nbf-aneka-form__grid">{isi.map(isian)}</div>
      </section>
    )
  const kecil = GRUP_KECIL.filter((g) => kelompok[g].length > 0)
  return (
    <div className="nbf-ringkas nbf-aneka-form">
      {GRUP_BESAR.map((g) => kartu(g, kelompok[g]))}
      {kecil.length > 0 && <div className="nbf-aneka-form__baris">{kecil.map((g) => kartu(g, kelompok[g]))}</div>}
      {c1 && tampilan.briguna && kartu('briguna', MEDAN_BRIGUNA)}
      {c1 && kartu('catatan', [{ jalur: 'others', label: F.others, jenis: 'area' } /* L8326 */])}
      {c1 && tampilan.bonding && (
        <section className="nbf-alamat__blok nbf-aneka-form__kartu">
          <h5 className="nbf-objek__judul">{M.judul}</h5>
          <div className="nbf-aneka-form__grid">{MEDAN_MEMORANDUM.map(isian)}</div>
        </section>
      )}
    </div>
  )
}
