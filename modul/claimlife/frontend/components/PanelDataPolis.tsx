// Panel data polis layar Register — A3 kelompok Register.
//
// Meniru himpunan medan `Section/InputRegisterClaimLife.xml` PERSIS: DUA PULUH
// SATU medan `.PolicyDataLife.*` (av-2, GILIRAN-11) ditambah tiga tombol, seluruhnya dengan nomor barisnya di
// `modul/claimlife/frontend/labels.ts`.
//
// ⛔ KEDUA PULUH SATU MEDAN TERIKAT KE `.PolicyDataLife.*`, bukan isian bebas. Di
// Pega ia terisi ketika polis dipilih lewat `Choose Policy No` (b3776), dan
// layar ini memperlakukannya sama: dibaca, tidak diketik. Membuatnya isian
// kosong akan menyuruh orang mengetik ulang apa yang sudah ada di sistem
// polis - dan dua salinan data polis akan berbeda pada hari pertama.
//
// ⛔ RALAT PEMBACAAN 27-09-2026. Ronde pertama panel ini menyebut ketiga
// medan `Ceding`, `Class of Business`, dan `Marketing Officer` sebagai
// isian ber-autocomplete, dan membangun rute `GET /api/rujukan/{jenis}`
// untuknya. KELIRU: keempat medan itu - ditambah `Type` - `pyReadOnly`
// `true` dan `pyEditOptions` `Read-only` di Pega, dengan `pyLabelFor`
// menunjuk `CedingCoName`, `BusinessName`, `MarketingName`, `Type`.
//
// Sebab salahnya layak dicatat: blok kontrolnya berdiri SEBELUM labelnya
// di DOM section, dan saya membaca jendela KE DEPAN dari label - yaitu
// grep dengan langkah tambahan, bukan pembacaan pohon. Kodenya dibuang
// di ketiga lapis.
//
// ⛔ BUTIR av DIJALANKAN 28-09-2026. `PolicyDataLife` kini DIBACA dari modul
// PremiumList Life lewat `GET /api/polis-life/ringkas` - tabel relasional
// `T_PREMIUM_LIST`, bukan cermin JSON-nya `[keputusan work owner]`.
//
// Sesudah av-2 (GILIRAN-11) ENAM BELAS dari dua puluh satu medan bersumber
// (dua medan retro hanya tampil bagi TP/TR). LIMA tidak, dan ketiadaannya
// tetap DINYATAKAN alih-alih diisi teks kosong:
//
//	TanggalRespon · TanggalKonfirmasi · TanggalRealisasi ·
//	TanggalKonfirmasiBalik · KetentuanUnderwriting
//
// Kelimanya properti halaman kerja Pega dan tidak punya kolom di migrasi
// 050-056 mana pun. Menghilangkan medannya dari layar membuat paritas tampak
// lengkap padahal tidak - dan tidak ada yang akan mencarinya lagi.
//
// ⚠️ VERSI POLIS IKUT DITAMPILKAN. Satu nomor polis punya banyak versi, dan
// yang dibaca adalah `PROD_KE` terbesar. Layar yang tidak dapat menyebut
// versi mana yang ditampilkannya membuat selisih angka mustahil ditelusuri.

import { REGISTER } from '../labels'
import type { PolicyDataLife } from '../api'
import { formatDate } from '../../../../inti/frontend/lib/format'

/**
 * Tanggal tampilan dd/mm/yyyy (keputusan work owner 10-10-2026). Memakai pengurai bersama inti (`formatDate`, yang
 * menjawab dd-mm-yyyy) lalu mengganti pemisahnya; nilai tak terurai / kosong = "—".
 */
export function tanggalTampil(v: string | undefined | null): string {
  const t = formatDate(v ?? '')
  return t === '' ? '—' : t.split('-').join('/')
}

/** Chip versi polis di judul panel: versi 0 = "New business", selebihnya "Endorsement /NN". */
export function labelVersiPolis(prodKe: number): string {
  return prodKe > 0 ? `Endorsement /${String(prodKe).padStart(2, '0')}` : 'New business'
}

/** Satu medan panel: label VERBATIM, nilainya, dan keadaan sumbernya. */
export interface MedanPolis {
  label: string
  nilai: string
  /** true bila backend belum punya jalan mengisinya. */
  belumBersumber?: boolean
  /** Nilai kode (PL Number) — satu baris, tidak dipecah (11-10-2026). */
  kode?: boolean
}

/**
 * Menyusun kesebelas medan dari apa yang backend BENAR-BENAR kirim.
 *
 * ⚠️ Dipisah dari komponennya supaya dapat diuji tanpa DOM - dan supaya
 * daftar medan yang belum bersumber dapat dihitung oleh uji, bukan dipercaya
 * dari komentar.
 */
export function medanPolis(sumber: {
  nomorPolis?: string
  namaTertanggung?: string
  mataUang?: string
  type?: string
  /** Data polis dari PremiumList Life; null berarti belum terbaca. */
  polis?: PolicyDataLife | null
}): MedanPolis[] {
  const ada = (v: string | undefined): string => (v === undefined || v === '' ? '—' : v)
  const p = sumber.polis ?? null
  // ⛔ Nilai dari polis dipakai HANYA bila polisnya terbaca. Jatuh kembali ke
  // `sumber.type` saat polis belum ada menjaga layar tetap berguna untuk
  // klaim atas polis yang belum ada di modul PremiumList Life.
  const dariPolis = (v: string | undefined): MedanPolis =>
    p === null ? { nilai: '—', label: '', belumBersumber: true } : { nilai: ada(v), label: '' }
  const medan = (label: string, v: string | undefined): MedanPolis => ({
    ...dariPolis(v),
    label,
  })
  // ⛔ `pyCondition` b10541/b10824 `.PolicyDataLife.Type=='TP'||'TR'`: dua
  // medan retro TAMPIL hanya bagi treaty. Type diambil dari polis, atau dari
  // pilihan pemakai bila polisnya belum terbaca.
  const tipe = p?.type ?? sumber.type ?? ''
  const treaty = tipe === 'TP' || tipe === 'TR'
  const tanpaKolom = (label: string): MedanPolis => ({ label, nilai: '—', belumBersumber: true })
  return [
    // Punya sumber hari ini - dari `GET /api/peserta-life`.
    { label: REGISTER.namaTertanggung, nilai: ada(sumber.namaTertanggung) },
    // ⛔ `Type` `pyReadOnly` true di Pega (b9058, `pyLabelFor` `Type` b9093)
    // dan terikat `.PolicyDataLife.Type`. Ronde pertama menjadikannya isian
    // bebas - itu ralat yang sama. Ia dan medan sesudahnya bersumber
    // `T_PREMIUM_LIST` lewat butir av dan av-2.
    //
    // ⚠️ `Type` jatuh kembali ke `sumber.type` bila polisnya belum terbaca:
    // klaim dapat didaftarkan atas polis yang belum ada di PremiumList Life.
    p === null
      ? { label: REGISTER.type, nilai: ada(sumber.type), belumBersumber: true }
      : { label: REGISTER.type, nilai: ada(p.type) },
    // Urut section, b9301 sampai b14806.
    // ⛔ `System Reinsurance` (b9301) dan `Premium Method` (b9694) TIDAK ditampilkan sejak 10-10-2026 (keputusan work
    // owner "system reinsurance, premium method tidak perlu di tampilkan") - tampilan saja; datanya tetap di
    // PolicyDataLife.
    medan(REGISTER.marketing, p?.marketingName),
    medan(REGISTER.wpc, p === null ? undefined : tanggalTampil(p.wpc)),
    ...(treaty
      ? [medan(REGISTER.retroName, p?.retroName), medan(REGISTER.securityReinsurer, p?.securityReinsurer)]
      : []),
    medan(REGISTER.ceding, p?.cedingCoName),
    medan(REGISTER.pemegangPolis, p?.policyHolderName),
    medan(REGISTER.sob, p?.sobName),
    medan(REGISTER.kelasBisnis, p?.businessName),
    medan(REGISTER.idProduk, p?.productNameId),
    medan(REGISTER.namaProduk, p?.productName),
    medan(REGISTER.tanggalEmail, p === null ? undefined : tanggalTampil(p.dateReceived)),
    // ⛔ KELIMA INI TETAP TANPA SUMBER: nol kolom di migrasi 050-056 mana pun
    // (`MedanTanpaSumberPolicyData` di backend menyebut kelimanya).
    tanpaKolom(REGISTER.tanggalRespon),
    tanpaKolom(REGISTER.tanggalKonfirmasi),
    medan(REGISTER.status, p?.status),
    tanpaKolom(REGISTER.konfirmasiBalik),
    medan(REGISTER.statusDiperbarui, p?.statusUpdate),
    tanpaKolom(REGISTER.tanggalRealisasi),
    tanpaKolom(REGISTER.catatanUnderwriter),
  ]
}

/**
 * Tab Policy Details layar Outstanding — TIGA KOLOM seperti screenshot Pega CLMLF-9293 (keputusan work owner 11-10-2026),
 * urut `InputOSClaimLife.xml` b6734–b12679.
 *
 * ⛔ Berbeda dari panel Register: System Reinsurance (`.PolicyDataLife.TypeCeding` b7102) dan Premium Method
 * (`.PolicyDataLife.ProRateType` b7495) DITAMPILKAN di sini (screenshot 11-10-2026); keputusan 10-10-2026 yang
 * menyembunyikannya tetap berlaku untuk panel Register.
 *
 * ⚠️ Type tampil sebagai KODE (QP/QR/TP/TR): dropdown b6905 memakai daftar properti (`pyListSource associated`) yang
 * rule-nya tidak ada di korpus, jadi pemetaan kode → nama tidak dapat diverifikasi.
 */
export function medanPolisOutstanding(polis: PolicyDataLife | null, plNumber: string): MedanPolis[][] {
  const ada = (v: string | undefined): string => (v === undefined || v === '' ? '—' : v)
  const p = polis
  const medan = (label: string, v: string | undefined): MedanPolis =>
    p === null ? { label, nilai: '—', belumBersumber: true } : { label, nilai: ada(v) }
  const tanpaKolom = (label: string): MedanPolis => ({ label, nilai: '—', belumBersumber: true })
  const treaty = p?.type === 'TP' || p?.type === 'TR'
  return [
    [
      // 11-10-2026 "rapikan kolom pl number": label "PL Number" (bukan nama kolom), nilai kode satu baris.
      { label: 'PL Number', nilai: ada(plNumber), kode: true },
      medan(REGISTER.type, p?.type),
      medan(REGISTER.sistemReasuransi, p?.typeCedingName || p?.typeCeding),
      medan(REGISTER.metodePremi, p?.proRateType),
      medan(REGISTER.marketing, p?.marketingName),
      medan(REGISTER.wpc, p === null ? undefined : tanggalTampil(p.wpc)),
      ...(treaty
        ? [medan(REGISTER.retroName, p?.retroName), medan(REGISTER.securityReinsurer, p?.securityReinsurer)]
        : []),
    ],
    [
      medan(REGISTER.ceding, p?.cedingCoName),
      medan(REGISTER.pemegangPolis, p?.policyHolderName),
      medan(REGISTER.sob, p?.sobName),
      medan(REGISTER.kelasBisnis, p?.businessName),
      medan(REGISTER.idProduk, p?.productNameId),
      medan(REGISTER.namaProduk, p?.productName),
    ],
    [
      medan(REGISTER.tanggalEmail, p === null ? undefined : tanggalTampil(p.dateReceived)),
      tanpaKolom(REGISTER.tanggalRespon),
      tanpaKolom(REGISTER.tanggalKonfirmasi),
      medan(REGISTER.status, p?.status),
      tanpaKolom(REGISTER.konfirmasiBalik),
      medan(REGISTER.statusDiperbarui, p?.statusUpdate),
      tanpaKolom(REGISTER.tanggalRealisasi),
      tanpaKolom(REGISTER.catatanUnderwriter),
    ],
  ]
}

/** Panel tiga kolom label-nilai tab Policy Details layar Outstanding. */
export function PanelPolisOutstanding({
  polis,
  plNumber,
  polisTidakDitemukan,
}: {
  polis: PolicyDataLife | null
  plNumber: string
  polisTidakDitemukan?: { idPega: string }
}) {
  return (
    <section className="os__form os__form--tiga os__form--polis" aria-label="Policy Details">
      {medanPolisOutstanding(polis, plNumber).map((kolom, i) => (
        <dl key={i} className="os__kolom">
          {kolom.map((m) => (
            <div key={m.label} className={`os__baris-form${m.belumBersumber === true ? ' polis__medan--belum' : ''}`}>
              <dt>{m.label}</dt>
              <dd className={m.kode === true ? 'os__nilai-pl' : undefined}>{m.nilai}</dd>
            </div>
          ))}
        </dl>
      ))}
      {polis === null && polisTidakDitemukan !== undefined && (
        <p className="alert alert--warn os__lebar" role="note">
          Policy data not found for this insured ({polisTidakDitemukan.idPega}).
        </p>
      )}
    </section>
  )
}

export interface PanelProps {
  nomorPolis?: string
  namaTertanggung?: string
  mataUang?: string
  type?: string
  /** Data polis dari PremiumList Life; null berarti belum terbaca. */
  polis?: PolicyDataLife | null
  /** Header polis peserta tidak ditemukan (10-10-2026) - satu kalimat sederhana menyebut IDPEGA-nya. */
  polisTidakDitemukan?: { idPega: string }
}

export function PanelDataPolis(p: PanelProps) {
  const medan = medanPolis(p)

  return (
    <section className="polis">
      {/* Judul + chip versi (dulu baris teknis "Polis … — versi … · produk …", dihapus 10-10-2026). */}
      <h3 className="polis__judul">
        Data Polis
        {p.polis != null && <span className="polis__chip">{labelVersiPolis(p.polis.prodKe)}</span>}
      </h3>
      {/* Grid label-nilai: label kecil abu-abu di atas, nilai di bawah; nilai panjang membungkus. */}
      <dl className="polis__daftar">
        {medan.map((m) => (
          <div
            key={m.label}
            className={`polis__medan${m.belumBersumber === true ? ' polis__medan--belum' : ''}`}
          >
            <dt>{m.label}</dt>
            <dd>{m.nilai}</dd>
          </div>
        ))}
      </dl>
      {/* ⛔ Keputusan work owner 10-10-2026: klaim hanya untuk polis Resolved-Completed - Find Insured tetap menampilkan
          pesertanya, panel memperingatkan (Register menolak 422). */}
      {p.polis != null && p.polis.status !== 'Resolved-Completed' && (
        <p className="alert alert--warn" role="note">
          Policy not Resolved-Completed (status: {p.polis.status === '' ? '—' : p.polis.status}). A claim cannot be registered yet.
        </p>
      )}
      {p.polis == null && p.polisTidakDitemukan !== undefined && (
        <p className="alert alert--warn" role="note">
          Policy data not found for this insured ({p.polisTidakDitemukan.idPega}).
        </p>
      )}
    </section>
  )
}
