// Key baris popup: HARUS unik per baris. View CLAIM_MASTER_TREATY bisa memuat baris yang sama pada Treaty ID + Treaty
// Group + Class of Business tetapi berbeda kolom lain (DEV 08-10-2026: 500 baris, 386 key tiga kolom). Key kembar
// membuat React membiarkan baris lama tertinggal saat daftar berganti (temuan work owner: filter Treaty ID menampilkan
// treaty lain). Server mengirim DISTINCT atas seluruh kolom, jadi gabungan seluruh kolom unik.

export interface KolomMaster {
  treatyId: string
  classOfBusiness: string
  classOfBusinessId: string
  treatyContractName: string
  sob: string
  ceding: string
  treatyType: string
  proportionType: string
  treatyGroup: string
  treatyGroupId: string
  treatyYear: string
}

export function kunciMaster(b: KolomMaster): string {
  return [
    b.treatyId,
    b.classOfBusiness,
    b.classOfBusinessId,
    b.treatyContractName,
    b.sob,
    b.ceding,
    b.treatyType,
    b.proportionType,
    b.treatyGroup,
    b.treatyGroupId,
    b.treatyYear,
  ].join('\u0001')
}
