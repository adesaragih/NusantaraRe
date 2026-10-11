// Fungsi murni tab Object Marine Cargo (tiket 74 P1) - baris baru, perapian jawaban server, dan penerapan popup Select Ship /
// Select Agent. Bukti `NB\` = `D:\migrasi\RNM\NB FacIn\`.

import type { BarisAgentSurvei, BarisKapal, CoverageCargo, KapalCargo, KapalTambahan, ObjekCargo } from '../api'

export function kapalKosong(): KapalCargo {
  return {
    id: '',
    nmShip: '',
    shipTypeId: '',
    shipTypeName: '',
    grt: '',
    nrt: '',
    dwt: '',
    age: '',
    flagName: '',
    construction: '',
    yMake1: '',
    yMake2: '',
    remark: '',
    register: '',
    imo: '',
    formerShipName: '',
    classification: '',
    rebuilt: '',
  }
}

export function kapalTambahanBaru(): KapalTambahan {
  return { id: '', register: '', nmShip: '', voyageNo: '', shipTypeId: '', flagName: '', nrt: '', dwt: '', grt: '', yMake1: '', construction: '' }
}

/** Baris cargo baru - tombol "+" grid (`InputDtlObject_FacIn.xml` L4086 addRow AFTER). */
export function cargoBaru(): ObjekCargo {
  return {
    id: '',
    isTopRisk: false,
    tradingId: '',
    tradingNote: '',
    goodId: '',
    goodNote: '',
    conveyanceId: '',
    conveyanceName: '',
    conveyanceNote: '',
    packingId: '',
    packingNote: '',
    fromRute: '',
    toRute: '',
    sailDate: '',
    blNumber: '',
    invoiceNumber: '',
    invoiceDate: '',
    maxAge: '',
    maxAgeDetail: '',
    surveyAgent: { agentId: '', fullName: '', address: '', city: '', country: '' },
    lc: { lc: false, bank: '', number: '', condition: '', remark: '', endDate: '' },
    ship: kapalKosong(),
    additionalShip: [],
  }
}

/** Jawaban server dirapikan: halaman anak yang absen (null / tidak dikirim) diisi kosong supaya form tidak runtuh. */
export function rapikanCargo(c: Partial<ObjekCargo>): ObjekCargo {
  const b = cargoBaru()
  return {
    ...b,
    ...c,
    surveyAgent: { ...b.surveyAgent, ...(c.surveyAgent ?? {}) },
    lc: { ...b.lc, ...(c.lc ?? {}) },
    ship: { ...b.ship, ...(c.ship ?? {}) },
    additionalShip: (c.additionalShip ?? []).map((k) => ({ ...kapalTambahanBaru(), ...k })),
  }
}

/**
 * Popup Select Ship → `.PolicyData.Ship` (`NB\Activity\SelectShip_act.xml` L470-873). GRT / NRT / DWT memakai koma desimal di
 * master: `@toDecimal(@replaceAll(x, ',', '.'))` (L665 / L706 / L727). Age dihitung server (`CountAgeShip`).
 */
export function terapkanKapal(c: ObjekCargo, k: BarisKapal): ObjekCargo {
  const angka = (v: string) => v.replace(/,/g, '.')
  return {
    ...c,
    ship: {
      ...c.ship,
      id: k.id,
      nmShip: k.nama,
      shipTypeId: k.shipTypeId ?? '',
      shipTypeName: k.shipTypeName ?? '',
      formerShipName: k.formerShipName ?? '',
      classification: k.klasifikasi,
      flagName: k.flag,
      construction: k.konstruksi,
      grt: angka(k.grt),
      nrt: angka(k.nrt),
      dwt: angka(k.dwt),
      yMake1: k.tahun,
      yMake2: k.yMake2 ?? '',
      remark: k.remark ?? '',
      register: k.register ?? '',
      imo: k.imo ?? '',
      rebuilt: k.rebuilt ?? '',
      age: '',
    },
  }
}

/** Popup Select Agent → `.PolicyData.SurveyAgent` per cargo (K74-6; nilai yang di-SET autocomplete Agent ID `MarineCargoDtl.xml` L1392). */
export function terapkanAgent(c: ObjekCargo, a: BarisAgentSurvei): ObjekCargo {
  return { ...c, surveyAgent: { agentId: a.id, fullName: a.name, address: a.address, city: a.city, country: a.country } }
}

/**
 * Autocomplete Register ID baris Additional Ship (`MarineCargoDtl.xml` L16755, `BrowseShip_RD`): memasang `.REGISTER,
 * .SHIPTYPE_ID, .NM_SHIP, .FLAG_NAME, .NRT, .DWT, .GRT, .Y_MAKE1, .CONST`; Voyage No tetap.
 */
export function terapkanKapalTambahan(t: KapalTambahan, k: BarisKapal): KapalTambahan {
  return {
    ...t,
    id: k.id,
    register: k.register ?? '',
    shipTypeId: k.shipTypeId ?? '',
    nmShip: k.nama,
    flagName: k.flag,
    nrt: k.nrt,
    dwt: k.dwt,
    grt: k.grt,
    yMake1: k.tahun,
    construction: k.konstruksi,
  }
}

/** Tanggal kabel DD-MM-YYYY → dd/mm/yyyy untuk sel grid SailDate (`Date-Short-Custom-YYYY`). */
export function tampilTanggal(v: string): string {
  return v.replace(/-/g, '/')
}

/** Coverage cargo kosong (bentuk kabel; proRatePercent dari server - `AppendCoverageMarineCargo_ACT`). */
export function coverageCargoKosong(): CoverageCargo {
  return {
    coverage: '',
    coverageNote: '',
    coverageInitial: '',
    currency: '',
    rate: '',
    limitOfLiability: '',
    discountPercentage: '',
    tsi: '',
    premium: '',
    minPremium: '',
    warrantyExtension: [],
    deductibles: [],
  }
}

/** Jawaban server dirapikan: daftar anak yang absen (null) menjadi kosong. */
export function rapikanCoverageCargo(c: Partial<CoverageCargo>): CoverageCargo {
  return { ...coverageCargoKosong(), ...c, warrantyExtension: c.warrantyExtension ?? [], deductibles: c.deductibles ?? [] }
}

/**
 * Hasil hitung diterapkan ke coverage di layar: medan isian + server dari jawaban, sedangkan Warranty Extension dan
 * Deductible yang di layar dipertahankan (rantai hitung tidak menyentuhnya).
 */
export function terapkanHitungCargo(layar: CoverageCargo, hasil: CoverageCargo): CoverageCargo {
  return { ...rapikanCoverageCargo(hasil), warrantyExtension: layar.warrantyExtension, deductibles: layar.deductibles }
}
