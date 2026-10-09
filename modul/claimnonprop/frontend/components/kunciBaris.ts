// Disalin dari `modul/claimprop/frontend/components/kunciBaris.ts` (pola, bukan impor): keputusan work owner yang disebut di bawah
// adalah keputusan layar Claim Prop yang ditiru Claim Non Prop (prompt Claim Non Prop tahap 1).
// Key baris popup: HARUS unik per baris. Popup master bisa memuat baris yang sama pada Treaty ID + Treaty Group + Class
// of Business tetapi berbeda kolom lain (DEV Claim Prop 08-10-2026: 500 baris, 386 key tiga kolom). Key kembar
// membuat React membiarkan baris lama tertinggal saat daftar berganti (temuan work owner: filter Treaty ID menampilkan
// treaty lain). Server mengirim DISTINCT atas seluruh kolom, jadi gabungan seluruh kolom unik.

/** Kolom baris popup ChooseMasterTNonProp - `models.BarisMaster` (server: DISTINCT atas kesembilan kolom). */
export interface KolomMaster {
  treatyId: string
  treatyContractName: string
  proportionType: string
  classOfBusinessId: string
  classOfBusiness: string
  sob: string
  ceding: string
  treatyYear: string
  treatyGroup: string
}

export function kunciMaster(b: KolomMaster): string {
  return [
    b.treatyId,
    b.treatyContractName,
    b.proportionType,
    b.classOfBusinessId,
    b.classOfBusiness,
    b.sob,
    b.ceding,
    b.treatyYear,
    b.treatyGroup,
  ].join('\u0001')
}
