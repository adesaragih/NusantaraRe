// Label modul Endorsment Fac In - teks layar Pega VERBATIM (bahasa Inggris/Indonesia apa adanya), dengan asalnya.
//
// Asal (Pega, `D:\migrasi\RNM\Endorsment Fac In\`):
//   - portal: `Section/SFAPortal_EndorsementHeader.xml` (judul L1757, Create Addendum L1053-1072);
//   - form: `Section/WorkPrimaryDetails.xml` (`pyLabelFieldValue` tiap medan), `Section/EdmType1/2/3.xml`;
//   - tombol: `Section/crmNewHarnessButtons.xml` (OK L847-1102, Cancel L1210-1295);
//   - pilihan Endorsement Type: `D:\migrasi\RNM\DDL\EdmType.xml` L221-237 (prompt list, urutannya apa adanya);
//   - pilihan Detail EDM Type: `Activity/SetEdmType.xml` L373-1225 (nilai `.CARI1`, teks `.CARI2`).

/** Portal (`SFAPortalEndorsement`). */
export const PORTAL = {
  judul: 'Addendum',
  buat: 'Create Addendum',
  /** Grid portal tidak ada di korpus (sumber datanya `[pertanyaan terbuka]`) - teks penjelas sistem baru, bukan Pega. */
  tanpaDaftar: 'The endorsement case list is not available yet.',
  petunjukTanpaDaftar: 'Use Create Addendum to check a policy before an endorsement case is created.',
} as const

/** Form Create Addendum (`WorkPrimaryDetails`). */
export const FORM = {
  judul: 'Create Addendum',
  policyNo: 'Policy Number / Cover Number',
  sumber: 'Endorsement Source',
  catatan: 'Endorsement Note',
  jenis: 'Endorsement Type',
  pilihJenis: 'Choose Endorsement Type',
  /** `.Quotation.EdmTypeNew` - Pega tanpa label medan; placeholder-nya dipakai sebagai label. */
  jenisBaru: 'Choose EDM Type',
  detail: 'Detail EDM Type',
  tanggal: 'Endorsement Date',
  biaya: 'Biaya',
  survey: 'Survey',
  retro: 'Endorsement Internal Retro',
  ok: 'OK',
  batal: 'Cancel',
  // Teks penata layar sistem baru (bukan Pega) - templat NB FacIn, work owner 11-10-2026 "biar user friendly".
  bagianPolis: 'Policy',
  bagianJenis: 'Endorsement',
  bagianLain: 'Additional',
  kalender: 'Open calendar',
  formatTanggal: 'Use the date format dd/mm/yyyy',
} as const

export interface Pilihan {
  nilai: string
  label: string
}

/** Endorsement Type - `DDL\EdmType.xml` (K-029). */
export const JENIS_EDM: readonly Pilihan[] = [
  { nilai: '4', label: 'Penambahan / Pengurangan / Perubahan' },
  { nilai: '2', label: 'Batal Prorata' },
  { nilai: '1', label: 'Batal Sejak Semula' },
]

/**
 * `.Quotation.EdmTypeNew` - ⚠️ `[dugaan]` teks: prompt list properti ini TIDAK ada di korpus. Nilai 1/2/3 terbukti dari
 * visibilitas `EdmType1/2/3` (`.Quotation.EdmTypeNew==1/2/3` menampilkan daftar `TypePenambahan` / `TypePengurangan` /
 * `TypePerubahan`); nilai 4 dari bypass "edm ri slip" (`SetErrorBatalEndorsement_Act` 5.1, `CARI41=="4"`). Teksnya diambil
 * dari nama daftar itu - menunggu prompt list asli dari work owner.
 */
export const JENIS_EDM_BARU: readonly Pilihan[] = [
  { nilai: '1', label: 'Penambahan' },
  { nilai: '2', label: 'Pengurangan' },
  { nilai: '3', label: 'Perubahan' },
  { nilai: '4', label: 'RI Slip' },
]

/** Detail EDM Type per `EdmTypeNew` - `SetEdmType` apa adanya (K-046: Pengurangan tanpa 12, langkahnya menulis ke Penambahan). */
export const DETAIL_EDM: Readonly<Record<string, readonly Pilihan[]>> = {
  '1': [
    { nilai: '1', label: 'Extend Period' },
    { nilai: '2', label: 'Adjustment TSI / Add Object / Rate / Premium' },
    { nilai: '4', label: 'Adjustment Spreading' },
    { nilai: '6', label: 'Adjustment Period' },
    { nilai: '7', label: 'Adjustment Currency' },
    { nilai: '8', label: 'Adjustment Deduction' },
    { nilai: '11', label: 'Adjustment Share Cedant' },
    { nilai: '12', label: 'Adjustment PPN/PPH' },
  ],
  '2': [
    { nilai: '2', label: 'Adjustment TSI / Add Object / Rate / Premium' },
    { nilai: '4', label: 'Adjustment Spreading' },
    { nilai: '6', label: 'Adjustment Period' },
    { nilai: '7', label: 'Adjustment Currency' },
    { nilai: '8', label: 'Adjustment Deduction' },
    { nilai: '11', label: 'Adjustment Share Cedant' },
  ],
  '3': [
    { nilai: '9', label: 'Adjustment Insured Name' },
    { nilai: '0', label: 'Adjustment Reff. Number' },
    { nilai: '11', label: 'Adjustment Share Cedant' },
  ],
}

/**
 * `.EndorsementInternalRetro` - pxRadioButtons, bawaan `0` (WorkPrimaryDetails L10212); gerbang 6 membaca `==1`.
 * ⚠️ Prompt list properti ini TIDAK ada di korpus: kodenya ditampilkan apa adanya sampai teksnya diketahui.
 */
export const PILIHAN_RETRO: readonly Pilihan[] = [
  { nilai: '1', label: '1' },
  { nilai: '0', label: '0' },
]

/** Teks pemeriksaan - teks sistem baru (bukan Pega): Pega hanya mematikan tombol OK dan memasang pesan di medan. */
export const PERIKSA = {
  memeriksa: 'Checking the policy…',
  lolos: 'The policy passed every check.',
  belumLengkap: 'Some checks could not be run yet - the endorsement case cannot be created until they are.',
  judulGerbang: 'Policy checks',
  gerbang: {
    1: 'Policy already cancelled by endorsement',
    2: 'Unfinished endorsement on this policy',
    3: 'Claim on this policy',
    4: 'Payment on this policy (cancellation only)',
    5: 'Renewal of this policy',
    6: 'Fac out spreading (internal retro only)',
  } as Readonly<Record<number, string>>,
  hasil: {
    lolos: 'passed',
    tolak: 'rejected',
    diklep: 'exception applies',
    dilewati: 'skipped',
    'belum-diperiksa': 'not checked yet',
  } as Readonly<Record<string, string>>,
  bisnis: 'Line of business',
  periode: (mulai: string, akhir: string) => (mulai === '' && akhir === '' ? 'Policy period' : `Policy period ${mulai} – ${akhir}`),
  /** OK mati - Pega mematikannya selama `CARIGROUP` belum `true` (crmNewHarnessButtons L885). */
  okMati: 'Fill in the required fields and pass every policy check first.',
  /** Alasan OK mati, satu per keadaan - ditampilkan di bawah kepala form. */
  alasan: {
    wajib: (medan: string[]) => `Required: ${medan.join(', ')}.`,
    memeriksa: 'Checking the policy…',
    ditolak: 'The policy did not pass the checks below.',
    belumLengkap: 'Some checks could not be run yet.',
  },
  belumDiperiksa: 'Enter the policy number to run the policy checks.',
  menyimpan: 'Creating the endorsement case…',
  /** Pesan portal sesudah kasus lahir - teks sistem baru. */
  dibuat: (id: string) => `Endorsement case ${id} created.`,
  tutup: 'Close',
} as const
