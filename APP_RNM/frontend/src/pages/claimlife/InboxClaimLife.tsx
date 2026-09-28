// Inbox Claim Life — F0.4, brief lanjutan 7 §1.3.
//
// Empat tab = empat assignment `Register_Flow.xml`, judulnya VERBATIM
// `pyTaskName` (358 · 343 · 268 · 313). Tab yang perannya tidak dipegang
// pelaku TIDAK dirender — bukan dirender lalu dikosongkan: tab kosong yang
// tidak pernah dapat berisi hanyalah tempat orang menunggu sia-sia.
//
// Kolomnya dari `ReportDefinition/InboxPremiumList.xml`, label VERBATIM:
// `Case ID` 721 · `Create Date/Time` 736 · `Create Operator Name` 751 ·
// `Work Status` 765; 50 baris per halaman (592), urut waktu buat DESC (733).
//
// ⚠️ Klik baris membuka layar TAHAP kasus itu. Sampai kelompok A3 masing-
// masing selesai, ia membuka halaman Detail yang ada dengan keterangan
// `BelumTersedia` DI DALAM halaman — bukan butir menu.

import { useCallback, useEffect, useState } from 'react'

import { MENU } from '../../assets/labels'
import { PERAN, TAHAP, TAHAP_ID, type KodePeran } from '../../assets/labels.claimlife'
import { Gagal, Kosong, Memuat } from '../../components/ui/dasar'
import { unduhXlsx, type KolomEksporXlsx } from '../../lib/exportXlsx'
import {
  ambilKotakMasuk,
  TAHAP_NOMOR,
  type BarisInbox,
  type HalamanInbox,
  type NomorTahap,
} from '../../services/api'

/** Satu tab: nomor tahap, judul VERBATIM, dan peran yang memegangnya. */
interface Tab {
  nomor: NomorTahap
  judul: string
  judulID: string
  peran: KodePeran
}

/**
 * Keempat tab, berurut tangga.
 *
 * ⛔ `peran` menentukan tab mana yang TAMPIL - `[terverifikasi]`
 * `Register_Flow.xml`: Input Register dan Outstanding Claim keduanya
 * `ReasLifeAdmin` (1508 `Current operator`, 1351 `ToWorklist`), Medical Check
 * `ReasLifeMedicalAdvisor` (1072), Claim Analis `ReasLifeSPV` (1198).
 */
const TAB: readonly Tab[] = [
  {
    nomor: TAHAP_NOMOR.inputRegister,
    judul: TAHAP.inputRegister,
    judulID: TAHAP_ID.inputRegister,
    peran: PERAN.admin,
  },
  {
    nomor: TAHAP_NOMOR.outstanding,
    judul: TAHAP.outstandingClaim,
    judulID: TAHAP_ID.outstandingClaim,
    peran: PERAN.admin,
  },
  {
    nomor: TAHAP_NOMOR.medicalCheck,
    judul: TAHAP.medicalCheck,
    judulID: TAHAP_ID.medicalCheck,
    peran: PERAN.medis,
  },
  {
    nomor: TAHAP_NOMOR.claimAnalis,
    judul: TAHAP.claimAnalis,
    judulID: TAHAP_ID.claimAnalis,
    peran: PERAN.spv,
  },
]

/**
 * Tab yang pelaku ini berhak lihat.
 *
 * ⚠️ Pelaku berperan GANDA melihat gabungan - bukan salah satu. Layar yang
 * memaksa memilih satu peran menyembunyikan separuh pekerjaan orang itu.
 */
export function tabUntuk(peran: readonly KodePeran[]): Tab[] {
  return TAB.filter((t) => peran.includes(t.peran))
}

/** Label kolom — VERBATIM `InboxPremiumList.xml`, barisnya di komentar. */
const KOLOM = {
  caseId: 'Case ID', // 721
  tglCreate: 'Create Date/Time', // 736
  createOpName: 'Create Operator Name', // 751
  status: 'Work Status', // 765
  /** `[InputOSClaimLife.xml:951]` */
  nomorKlaim: 'Claim No',
  /** `[tidak ada sebagai label kolom]` — tombol `Choose Policy No` 3776. */
  nomorPolis: 'Policy No',
} as const

/**
 * Kolom ekspor xlsx — butir **bg**.
 *
 * ⛔ TEPAT kolom yang TAMPIL, dan berurutan sama. Ekspor yang memuat kolom
 * yang layarnya tidak tampilkan membuat berkas dan layar menjawab pertanyaan
 * yang berbeda - dan yang memegang berkasnya tidak akan tahu mana yang benar.
 *
 * ⚠️ Tanggal keluar APA ADANYA (RFC 3339), bukan dalam bentuk layar. Itu
 * penyimpangan yang disengaja dan berlawanan arah dengan kolomnya: berkas
 * lembar-sebar diurutkan dan disaring, dan `28-09-2026 14:03` diurutkan
 * sebagai TEKS - Desember mendahului Februari. Bentuk ISO diurutkan benar
 * oleh alat mana pun yang membukanya.
 */
export const KOLOM_EKSPOR: KolomEksporXlsx<BarisInbox>[] = [
  { kunci: 'caseId', label: KOLOM.caseId },
  { kunci: 'nomorKlaim', label: KOLOM.nomorKlaim },
  { kunci: 'nomorPolis', label: KOLOM.nomorPolis },
  { kunci: 'status', label: KOLOM.status },
  { kunci: 'createOpName', label: KOLOM.createOpName },
  { kunci: 'tglCreate', label: KOLOM.tglCreate },
]

function tanggalTampil(rfc3339: string): string {
  if (rfc3339 === '') return '—'
  const d = new Date(rfc3339)
  if (Number.isNaN(d.getTime())) return rfc3339
  const p = (n: number) => String(n).padStart(2, '0')
  return `${p(d.getDate())}-${p(d.getMonth() + 1)}-${d.getFullYear()} ${p(d.getHours())}:${p(d.getMinutes())}`
}

export interface InboxProps {
  peran: readonly KodePeran[]
  /** Membuka satu kasus. */
  onBuka: (workID: string) => void
  /** Membuka halaman Register — tombol kepala, `Register_Flow.xml:155`. */
  onRegister: () => void
}

export default function InboxClaimLife({ peran, onBuka, onRegister }: InboxProps) {
  const tab = tabUntuk(peran)
  const [aktif, setAktif] = useState<NomorTahap | null>(tab[0]?.nomor ?? null)
  const [hal, setHal] = useState<HalamanInbox | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [memuat, setMemuat] = useState(false)
  const [halaman, setHalaman] = useState(1)

  const muat = useCallback(async (t: NomorTahap, h: number) => {
    setMemuat(true)
    setGalat(null)
    try {
      setHal(await ambilKotakMasuk(t, h))
    } catch (e) {
      setHal(null)
      setGalat(e)
    } finally {
      setMemuat(false)
    }
  }, [])

  useEffect(() => {
    if (aktif === null) return
    void muat(aktif, halaman)
  }, [aktif, halaman, muat])

  if (tab.length === 0) {
    return (
      <Kosong
        pesan={
          'Tidak ada antrian untuk peran Anda. Keempat antrian Claim Life ' +
          'dipegang ReasLifeAdmin, ReasLifeMedicalAdvisor, dan ReasLifeSPV.'
        }
      />
    )
  }

  return (
    <section className="inbox">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">{MENU.inbox}</h2>
        {/* Tombol kepala — VERBATIM `Register_Flow.xml:155`. */}
        <button type="button" className="inbox__register" onClick={onRegister}>
          {MENU.register}
        </button>
        {/* ⛔ EKSPOR HANYA YANG TAMPIL — butir bg. Enam kolom yang sama
            dengan tabel di bawah, berurutan sama, dan HANYA halaman yang
            sedang terbuka. Mengekspor seluruh tahap dari tombol yang
            berdiri di atas satu tab berarti berkasnya menjawab pertanyaan
            yang berbeda dari yang layar tanyakan.

            ⚠️ Dimatikan saat belum ada baris: tombol yang menghasilkan
            berkas kosong mengajari orang mengabaikannya. */}
        <button
          type="button"
          className="inbox__ekspor"
          disabled={hal === null || hal.baris.length === 0}
          onClick={() => {
            if (hal === null) return
            unduhXlsx(`inbox-${hal.namaTahap || String(aktif)}.xlsx`,
              KOLOM_EKSPOR, hal.baris)
          }}
        >
          Export xlsx
        </button>
      </header>

      <div className="inbox__tab" role="tablist">
        {tab.map((t) => (
          <button
            key={t.nomor}
            type="button"
            role="tab"
            aria-selected={aktif === t.nomor}
            className={`inbox__tab-butir${aktif === t.nomor ? ' inbox__tab-butir--aktif' : ''}`}
            onClick={() => {
              setAktif(t.nomor)
              setHalaman(1)
            }}
          >
            {/* Judul VERBATIM, terjemahan DI SAMPING. */}
            <span className="inbox__tab-judul">{t.judul}</span>
            <span className="inbox__tab-id">{t.judulID}</span>
            {aktif === t.nomor && hal !== null && (
              <span className="inbox__lencana">{hal.total}</span>
            )}
          </button>
        ))}
      </div>

      {memuat && <Memuat />}
      {galat !== null && <Gagal galat={galat} />}

      {!memuat && galat === null && hal !== null && (
        <>
          {hal.baris.length === 0 ? (
            <Kosong pesan="Tidak ada kasus di antrian ini." />
          ) : (
            <table className="inbox__tabel">
              <thead>
                <tr>
                  <th>{KOLOM.caseId}</th>
                  <th>{KOLOM.nomorKlaim}</th>
                  <th>{KOLOM.nomorPolis}</th>
                  <th>{KOLOM.status}</th>
                  <th>{KOLOM.createOpName}</th>
                  <th>{KOLOM.tglCreate}</th>
                </tr>
              </thead>
              <tbody>
                {hal.baris.map((b: BarisInbox) => (
                  <tr
                    key={b.id}
                    className="inbox__baris"
                    onClick={() => {
                      onBuka(b.id)
                    }}
                  >
                    <td>{b.caseId}</td>
                    <td>{b.nomorKlaim}</td>
                    <td>{b.nomorPolis}</td>
                    <td>{b.status}</td>
                    <td>{b.createOpName}</td>
                    <td>{tanggalTampil(b.tglCreate)}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}

          {/* 50 per halaman (592); maksimum 500 (942). */}
          <nav className="inbox__halaman">
            <button
              type="button"
              disabled={halaman <= 1}
              onClick={() => {
                setHalaman((h) => Math.max(1, h - 1))
              }}
            >
              Sebelumnya
            </button>
            <span>
              Halaman {hal.halaman} — {hal.total} kasus
            </span>
            <button
              type="button"
              disabled={hal.halaman * hal.ukuran >= hal.total}
              onClick={() => {
                setHalaman((h) => h + 1)
              }}
            >
              Berikutnya
            </button>
          </nav>
        </>
      )}
    </section>
  )
}
