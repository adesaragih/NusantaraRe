// Panel data polis layar Register — A3 kelompok Register.
//
// Meniru himpunan medan `Section/InputRegisterClaimLife.xml` PERSIS: DUA PULUH
// SATU medan `.PolicyDataLife.*` (av-2, GILIRAN-11) ditambah tiga tombol, seluruhnya dengan nomor barisnya di
// `assets/labels.ts`.
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

import { REGISTER } from '../../assets/labels.claimlife'
import type { PolicyDataLife } from '../../services/api'

/** Satu medan panel: label VERBATIM, nilainya, dan keadaan sumbernya. */
export interface MedanPolis {
  label: string
  nilai: string
  /** true bila backend belum punya jalan mengisinya. */
  belumBersumber?: boolean
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
    // `TYPE_CEDING_NAME` yang ditampilkan; kodenya bila katanya kosong.
    medan(REGISTER.sistemReasuransi, p === null ? undefined : p.typeCedingName || p.typeCeding),
    medan(REGISTER.metodePremi, p?.proRateType),
    medan(REGISTER.marketing, p?.marketingName),
    medan(REGISTER.wpc, p?.wpc),
    ...(treaty
      ? [medan(REGISTER.retroName, p?.retroName), medan(REGISTER.securityReinsurer, p?.securityReinsurer)]
      : []),
    medan(REGISTER.ceding, p?.cedingCoName),
    medan(REGISTER.pemegangPolis, p?.policyHolderName),
    medan(REGISTER.sob, p?.sobName),
    medan(REGISTER.kelasBisnis, p?.businessName),
    medan(REGISTER.idProduk, p?.productNameId),
    medan(REGISTER.namaProduk, p?.productName),
    medan(REGISTER.tanggalEmail, p?.dateReceived),
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

export interface PanelProps {
  nomorPolis?: string
  namaTertanggung?: string
  mataUang?: string
  type?: string
  /** Data polis dari PremiumList Life; null berarti belum terbaca. */
  polis?: PolicyDataLife | null
}

export function PanelDataPolis(p: PanelProps) {
  const medan = medanPolis(p)
  const belum = medan.filter((m) => m.belumBersumber === true).length

  return (
    <section className="polis">
      <h3 className="polis__judul">Data Polis</h3>
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
      {p.polis !== null && p.polis !== undefined && (
        <p className="polis__versi" role="status">
          Polis {p.polis.nomorPolis} — versi {p.polis.prodKe}
          {p.polis.productName !== '' && <> · produk {p.polis.productName}</>}
        </p>
      )}
      {belum > 0 && (
        <p className="polis__catatan" role="note">
          {p.polis === null || p.polis === undefined ? (
            <>
              <strong>Data polis belum terbaca.</strong> Medan ini{' '}
              <em>read-only</em> dan terisi dari modul PremiumList Life lewat{' '}
              <code>.PolicyDataLife.*</code> saat nomor polis dipilih. Pilih
              nomor polis, atau polis itu memang belum ada di modul tersebut.
            </>
          ) : (
            <>
              {/* ⛔ KALIMAT BERBEDA untuk sebab yang berbeda. "Menunggu modul
                  PremiumList Life" sesudah modulnya ADA adalah kalimat yang
                  salah - dan kalimat yang salah membuat orang berhenti
                  membacanya. */}
              <strong>{belum} medan belum punya kolom di mana pun.</strong>{' '}
              <code>{p.polis.medanTanpaSumber.join(', ')}</code> adalah properti
              halaman kerja Pega dan tidak punya kolom di migrasi 050–056 mana
              pun. Medannya tetap ditampilkan supaya himpunannya sama dengan{' '}
              <code>InputRegisterClaimLife.xml</code>, dengan ketiadaannya
              dinyatakan alih-alih diisi teks kosong.
            </>
          )}
        </p>
      )}
    </section>
  )
}
