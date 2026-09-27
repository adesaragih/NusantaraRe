// Panel data polis layar Register — A3 kelompok Register.
//
// Meniru himpunan medan `Section/InputRegisterClaimLife.xml` PERSIS: sebelas
// medan data ditambah tiga tombol, seluruhnya dengan nomor barisnya di
// `assets/labels.ts`.
//
// ⛔ KESEBELAS MEDAN TERIKAT KE `.PolicyDataLife.*`, bukan isian bebas. Di
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
// ⛔ MEDAN YANG MENUNGGU DINYATAKAN, BUKAN DIHILANGKAN. `[DIPUTUSKAN
// work owner, butir av]` sumbernya modul **PremiumList Life**: di Pega
// `PolicyDataLife` diisi dari kasus modul itu, dan Claim Life TIDAK
// membaca cermin JSON-nya. Menghilangkan medannya dari layar membuat
// paritas tampak lengkap padahal tidak - dan tidak ada yang akan
// mencarinya lagi.

import { REGISTER } from '../assets/labels'

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
}): MedanPolis[] {
  const ada = (v: string | undefined): string => (v === undefined || v === '' ? '—' : v)
  return [
    // Punya sumber hari ini - dari `GET /api/peserta-life`.
    { label: REGISTER.namaTertanggung, nilai: ada(sumber.namaTertanggung) },
    // ⛔ `Type` IKUT MENUNGGU: ia `pyReadOnly` true di Pega (b9058,
    // `pyLabelFor` `Type` b9093) dan terikat `.PolicyDataLife.Type`.
    // Ronde pertama menjadikannya isian bebas - itu ralat yang sama.
    { label: REGISTER.type, nilai: ada(sumber.type), belumBersumber: true },
    { label: REGISTER.marketing, nilai: '—', belumBersumber: true },
    { label: REGISTER.ceding, nilai: '—', belumBersumber: true },
    { label: REGISTER.pemegangPolis, nilai: '—', belumBersumber: true },
    { label: REGISTER.kelasBisnis, nilai: '—', belumBersumber: true },
    { label: REGISTER.tanggalEmail, nilai: '—', belumBersumber: true },
    { label: REGISTER.tanggalRespon, nilai: '—', belumBersumber: true },
    { label: REGISTER.tanggalKonfirmasi, nilai: '—', belumBersumber: true },
    { label: REGISTER.status, nilai: '—', belumBersumber: true },
    { label: REGISTER.statusDiperbarui, nilai: '—', belumBersumber: true },
    { label: REGISTER.tanggalRealisasi, nilai: '—', belumBersumber: true },
  ]
}

export interface PanelProps {
  nomorPolis?: string
  namaTertanggung?: string
  mataUang?: string
  type?: string
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
      {belum > 0 && (
        <p className="polis__catatan" role="note">
          <strong>{belum} medan menunggu modul PremiumList Life.</strong>{' '}
          Di Pega medan ini <em>read-only</em> dan terisi dari kasus
          PremiumList Life lewat <code>.PolicyDataLife.*</code> saat nomor
          polis dipilih. Claim Life tidak membaca cermin JSON-nya
          (keputusan work owner <strong>av</strong>), jadi medannya
          ditampilkan apa adanya supaya himpunannya tetap sama dengan{' '}
          <code>InputRegisterClaimLife.xml</code>.
        </p>
      )}
    </section>
  )
}
