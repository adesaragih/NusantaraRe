// Label modul Treaty Contract Out — tiket 02 →.
//
// ⛔ VERBATIM dari korpus `D:\XML\RNM_BRD\Treaty Contract Out\`, dengan path +
// baris di sebelahnya (berkas korpus satu tag per baris; nomor baris mentah =
// nomor `sed -e 's/></>\n</g'`). Label yang diketik dari ingatan adalah label
// yang bergeser, dan yang bergeser tidak berbunyi.
//
// ⛔ Nol kata "Old" dan "testing" di nama apa pun (penyimpangan sadar 8);
// dijaga `components/treaty-contract-out/namaJujur.test.ts`.

/**
 * Menu kelompok **Treaty Contract Out** — tiga harness portal (kelas
 * `Data-Portal`), butirnya VERBATIM `<pyLabel>` tiap harness.
 *
 * ⚠️ Nama kelompok = nama FOLDER korpus, pola `MODUL` di `labels.ts`.
 * Kelompok ini belum ada di `MODUL` sampai tiket 03 menambahkannya ADITIF.
 */
export const MENU_TCO = {
  /** Nama folder korpus `D:\XML\RNM_BRD\Treaty Contract Out`. */
  kelompok: 'Treaty Contract Out',
  /** `Harness/InboxTreatyContract.xml` b151 `<pyLabel>`. */
  inboxTreatyContract: 'InboxTreatyContract',
  /** `Harness/InboxTreatyContractReinsType.xml` b151 `<pyLabel>`. */
  inboxTreatyContractReinsType: 'InboxTreatyContractReinsType',
  /** `Harness/InboxTreatyContractDescription.xml` b359 `<pyLabel>`. */
  inboxTreatyContractDescription: 'InboxTreatyContractDescription',
} as const

/**
 * Pemilih jenis reasuransi — tiket 02.
 *
 * Dua label korpus untuk medan yang SAMA (`ReinsTypeID` dari master
 * `REINSURANCETYPE`): form kontrak memakai `ReinsType`, form tahun memakai
 * `Reinsurance Type`. Keduanya dibawa apa adanya; pemanggil memilih.
 */
export const JENIS_REASURANSI_TCO = {
  /** `Section/InputTreatyContractReinsType.xml` b2652 `<pyLabelFieldValue>`. */
  reinsType: 'ReinsType',
  /** `Section/InputTreatyContract.xml` b7925 `<pyLabelFieldValue>` (dan b7948 `<pyLabelPreview>`). */
  reinsuranceType: 'Reinsurance Type',
  /** `[tidak ada di korpus]` — kosakata kami; keadaan yang ADR-0015 tuntut terlihat. */
  masterKosong: 'Master jenis reasuransi kosong atau tidak terbaca.',
  /** `[tidak ada di korpus]` — teks pilihan kosong `Pilih`. */
  belumDipilih: '-- pilih jenis reasuransi --',
} as const
