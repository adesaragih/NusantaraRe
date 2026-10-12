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

import { MENU, PERAN, type KodePeran } from '../../../../inti/frontend/labels'
import TabelKomite from '../components/TabelKomite'
import { REGISTER_BANYAK, TABEL_KOMITE, TAHAP } from '../labels'
import { Halaman, Kosong, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import { PitaGalat } from '../components/PitaGalat'
import { unduhXlsx, type KolomEksporXlsx } from '../../../../inti/frontend/lib/exportXlsx'
import {
  ambilKotakMasuk,
  kataCariInbox,
  MAKS_CARI_INBOX,
  TAHAP_NOMOR,
  type BarisInbox,
  type HalamanInbox,
  type NomorTahap,
} from '../api'

/** Satu tab: nomor tahap, judul VERBATIM, dan peran yang memegangnya. */
interface Tab {
  nomor: NomorTahap
  /** Nama tahap VERBATIM korpus (kunci uji, Export, kabar) — TIDAK ditampilkan sebagai label tab sejak 12-10-2026. */
  judul: string
  /** Label tampilan tab (keputusan work owner 12-10-2026, bukan nama korpus): Claim Admin / Medical Advisor / Claim Analyst. */
  label: string
  peran: KodePeran
  /**
   * Tab disembunyikan selama antriannya KOSONG (keputusan work owner 10-10-2026 "tab input register hilangkan saja,
   * karena sudah ada tombol register"). ⚠️ Bukan dihapus mutlak: `Send Back to Register` (PanelPindahTahap, b21404)
   * masih menaruh kasus di tahap Input Register - tanpa tab, kasus itu hilang dari pandangan. Tab muncul lagi
   * otomatis begitu ada kasus di sana.
   */
  sembunyiBilaKosong?: boolean
}

/**
 * Label tampilan tab inbox (keputusan work owner 12-10-2026: "ubah nama tabnya menjadi Claim Admin, Medical Advisor, Claim
 * Analyst"). BUKAN nama korpus: nilai TAHAP di basis data, kode tahap, dan logika tidak berubah.
 */
export const LABEL_TAB_INBOX = {
  admin: 'Claim Admin',
  medis: 'Medical Advisor',
  analis: 'Claim Analyst',
} as const

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
    label: TAHAP.inputRegister,
    peran: PERAN.admin,
    sembunyiBilaKosong: true,
  },
  {
    nomor: TAHAP_NOMOR.outstanding,
    judul: TAHAP.outstandingClaim,
    label: LABEL_TAB_INBOX.admin,
    peran: PERAN.admin,
  },
  {
    nomor: TAHAP_NOMOR.medicalCheck,
    judul: TAHAP.medicalCheck,
    label: LABEL_TAB_INBOX.medis,
    peran: PERAN.medis,
  },
  {
    nomor: TAHAP_NOMOR.claimAnalis,
    judul: TAHAP.claimAnalis,
    label: LABEL_TAB_INBOX.analis,
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

/**
 * Tab yang TAMPIL: tab ber-`sembunyiBilaKosong` (Input Register) hanya bila jumlah kasusnya sudah diketahui > 0
 * (10-10-2026). Jumlah yang belum terbaca = disembunyikan.
 */
/** Placeholder kotak cari inbox (keputusan work owner 12-10-2026). */
export const CARI_INBOX = 'Search case ID, claim no, PL number, insured name'

/** Jeda ketik sebelum kata kunci dikirim (pola Claim Fac In). */
export const JEDA_CARI_MS = 300

export function tabTampil(tab: readonly Tab[], jumlah: Partial<Record<NomorTahap, number>>): Tab[] {
  return tab.filter((t) => t.sembunyiBilaKosong !== true || (jumlah[t.nomor] ?? 0) > 0)
}

/** Label kolom — VERBATIM `InboxPremiumList.xml`, barisnya di komentar. */
const KOLOM = {
  caseId: 'Case ID', // 721
  tglCreate: 'Create Date/Time', // 736
  createOpName: 'Create Operator Name', // 751
  status: 'Work Status', // 765
  /** `[InputOSClaimLife.xml:951]` */
  nomorKlaim: 'Claim No',
  /** 11-10-2026: "kolom policy no ganti dengan pl number ambil dari polis". */
  plNumber: 'PL Number',
  /** Pega `pxUpdateOpName`; diturunkan dari jejak (lihat repository.BarisInbox.UpdateOpName). */
  updateOpName: 'Update Operator Name',
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
  // Urutan keputusan work owner 11-10-2026: caseid, claim no, pl_number, create op name, update op name, create date
  // time, work status.
  { kunci: 'caseId', label: KOLOM.caseId },
  { kunci: 'nomorKlaim', label: KOLOM.nomorKlaim },
  { kunci: 'plNumber', label: KOLOM.plNumber },
  { kunci: 'createOpName', label: KOLOM.createOpName },
  { kunci: 'updateOpName', label: KOLOM.updateOpName },
  { kunci: 'tglCreate', label: KOLOM.tglCreate },
  { kunci: 'status', label: KOLOM.status },
]

/** Waktu buat dd/mm/yyyy HH:mm (tampilan 10-10-2026, seperti inbox lain); ekspor xlsx tetap RFC 3339. */
export function tanggalTampil(rfc3339: string): string {
  if (rfc3339 === '') return '—'
  const d = new Date(rfc3339)
  if (Number.isNaN(d.getTime())) return rfc3339
  const p = (n: number) => String(n).padStart(2, '0')
  return `${p(d.getDate())}/${p(d.getMonth() + 1)}/${d.getFullYear()} ${p(d.getHours())}:${p(d.getMinutes())}`
}

export interface InboxProps {
  peran: readonly KodePeran[]
  /** Membuka satu kasus. */
  onBuka: (workID: string) => void
  /** Membuka halaman Register — tombol kepala, `Register_Flow.xml:155`. */
  onRegister: () => void
  /** `PropsRute.onBukaModul` - tabel komite di bawah inbox: kasus komite dibuka di tempat (menu komite dihapus). */
  onBukaModul?: (modul: string, id: string) => boolean
  /** Case ID yang baru saja didaftarkan (state rute, 10-10-2026) - pemberitahuan sukses; kosong = tidak tampil. */
  terdaftar?: readonly string[]
  /** Menutup pemberitahuan sukses. */
  onTutupTerdaftar?: () => void
  /** Pita sukses lain (mis. "Claim … sent to Medical Check.", 11-10-2026); kosong = tidak tampil. */
  kabar?: string
  /** Menutup pita kabar. */
  onTutupKabar?: () => void
}

/**
 * Pemberitahuan sukses Register (keputusan work owner 10-10-2026): alert hijau inti (`alert alert--ok`) dengan tombol
 * tutup - pola yang dipakai aplikasi (inti tidak punya toast).
 */
export function PemberitahuanTerdaftar({ ids, onTutup }: { ids: readonly string[]; onTutup?: () => void }) {
  if (ids.length === 0) return null
  return (
    <div className="alert alert--ok inbox__terdaftar" role="status">
      <span>
        {REGISTER_BANYAK.terdaftar.split('{n}').join(String(ids.length)).split('{ids}').join(ids.join(', '))}
      </span>
      {onTutup !== undefined && (
        <button type="button" className="btn btn--ghost btn--sm" onClick={onTutup}>
          {REGISTER_BANYAK.tutup}
        </button>
      )}
    </div>
  )
}

/** Pita sukses umum Inbox (pola PemberitahuanTerdaftar). */
export function PemberitahuanKabar({ kabar, onTutup }: { kabar: string; onTutup?: () => void }) {
  if (kabar === '') return null
  return (
    <div className="alert alert--ok inbox__terdaftar" role="status">
      <span>{kabar}</span>
      {onTutup !== undefined && (
        <button type="button" className="btn btn--ghost btn--sm" onClick={onTutup}>
          {REGISTER_BANYAK.tutup}
        </button>
      )}
    </div>
  )
}

export default function InboxClaimLife({
  peran,
  onBuka,
  onRegister,
  onBukaModul,
  terdaftar = [],
  onTutupTerdaftar,
  kabar = '',
  onTutupKabar,
}: InboxProps) {
  const semuaTab = tabUntuk(peran)
  // Jumlah kasus per tab (badge, 10-10-2026): halaman 1 tiap antrian, dibaca sekali saat inbox dibuka.
  const [jumlah, setJumlah] = useState<Partial<Record<NomorTahap, number>>>({})
  const tab = tabTampil(semuaTab, jumlah)
  // Tab bawaan = tab tampil pertama (Outstanding Claim untuk Admin).
  const [aktif, setAktif] = useState<NomorTahap | null>(tabTampil(semuaTab, {})[0]?.nomor ?? null)
  const [hal, setHal] = useState<HalamanInbox | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [memuat, setMemuat] = useState(false)
  const [halaman, setHalaman] = useState(1)
  // Kotak cari (work owner 12-10-2026 "tambahkan juga filter untuk mencari seperti claim facin"): `kata` = isian,
  // `cari` = kata kunci terkirim sesudah jeda 300 ms (pola Claim Fac In JEDA_CARI_MS); berlaku untuk antrian aktif.
  const [kata, setKata] = useState('')
  const [cari, setCari] = useState('')
  useEffect(() => {
    const t = setTimeout(() => {
      const k = kataCariInbox(kata)
      setCari((lama) => {
        if (lama !== k) setHalaman(1)
        return k
      })
    }, JEDA_CARI_MS)
    return () => clearTimeout(t)
  }, [kata])

  const muat = useCallback(async (t: NomorTahap, h: number, c: string) => {
    setMemuat(true)
    setGalat(null)
    try {
      setHal(await ambilKotakMasuk(t, h, undefined, c))
    } catch (e) {
      setHal(null)
      setGalat(e)
    } finally {
      setMemuat(false)
    }
  }, [])

  useEffect(() => {
    if (aktif === null) return
    void muat(aktif, halaman, cari)
  }, [aktif, halaman, cari, muat])

  const kunciTab = semuaTab.map((t) => t.nomor).join(',')
  useEffect(() => {
    let hidup = true
    for (const t of tabUntuk(peran)) {
      void ambilKotakMasuk(t.nomor, 1)
        .then((h) => {
          if (hidup) setJumlah((j) => ({ ...j, [t.nomor]: h.total }))
        })
        .catch(() => {
          // Badge tidak tampil bila hitungannya gagal; antrian aktif tetap menampilkan galatnya sendiri.
        })
    }
    return () => {
      hidup = false
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [kunciTab])

  // Jumlah antrian aktif mengikuti halaman yang baru dimuat — HANYA tanpa kata kunci: badge = jumlah seluruh antrian,
  // hasil cari tampil di paginasi (12-10-2026).
  useEffect(() => {
    if (aktif !== null && hal !== null && cari === '')
      setJumlah((j) => (j[aktif] === hal.total ? j : { ...j, [aktif]: hal.total }))
  }, [aktif, hal, cari])

  if (semuaTab.length === 0) {
    return (
      // 12-10-2026 (work owner "rapikan tampilan inbox committee ini"): akun tanpa antrian Claim Life langsung mendapat kotak
      // inbox ber-kulit berjudul Committee Inbox bila punya kasus komite (TabelKomite membaca daftarKomite); kalimat netral
      // hanya bila daftar komite pun kosong.
      <section className="inbox claimlife__akar">
        <TabelKomite
          peran={peran}
          onBukaModul={onBukaModul}
          utama
          kosong={<Kosong pesan={TABEL_KOMITE.kosongAkun} />}
        />
      </section>
    )
  }

  return (
    // claimlife__akar / claimlife__tabel (12-10-2026): kulit tampilan sama dengan inbox Claim Fac In - kelas saja, struktur tetap.
    <section className="inbox claimlife__akar">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">{MENU.inbox}</h2>
        {/* Tombol kepala — VERBATIM `Register_Flow.xml:155`. */}
        <div className="inbox__aksi-kepala">
        {/* Register HANYA di tab Claim Admin (keputusan work owner 12-10-2026); Export xlsx di semua tab. */}
        {aktif === TAHAP_NOMOR.outstanding && (
        <button type="button" className="btn btn--primary" onClick={onRegister}>
          {MENU.register}
        </button>
        )}
        {/* ⛔ EKSPOR HANYA YANG TAMPIL — butir bg. Enam kolom yang sama
            dengan tabel di bawah, berurutan sama, dan HANYA halaman yang
            sedang terbuka. Mengekspor seluruh tahap dari tombol yang
            berdiri di atas satu tab berarti berkasnya menjawab pertanyaan
            yang berbeda dari yang layar tanyakan.

            ⚠️ Dimatikan saat belum ada baris: tombol yang menghasilkan
            berkas kosong mengajari orang mengabaikannya. */}
        <button
          type="button"
          className="btn btn--ghost"
          disabled={hal === null || hal.baris.length === 0}
          onClick={() => {
            if (hal === null) return
            unduhXlsx(`inbox-${hal.namaTahap || String(aktif)}.xlsx`,
              KOLOM_EKSPOR, hal.baris)
          }}
        >
          Export xlsx
        </button>
        </div>
      </header>

      <PemberitahuanTerdaftar ids={terdaftar} onTutup={onTutupTerdaftar} />
      <PemberitahuanKabar kabar={kabar} onTutup={onTutupKabar} />

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
            {/* Judul VERBATIM (Inggris) saja; subjudul Indonesia dibuang 10-10-2026. Badge jumlah di setiap tab. */}
            <span className="inbox__tab-judul">{t.label}</span>
            {jumlah[t.nomor] !== undefined && <span className="inbox__lencana">{jumlah[t.nomor]}</span>}
          </button>
        ))}
      </div>

      {/* Kotak cari (12-10-2026) — letak dan gaya seperti toolbar Claim Fac In; saklar Teknik Claim Fac In tidak ditiru. */}
      <div className="toolbar claimlife__toolbar">
        <input
          className="field__input claimlife__cari"
          type="search"
          aria-label={CARI_INBOX}
          placeholder={CARI_INBOX}
          maxLength={MAKS_CARI_INBOX}
          value={kata}
          onChange={(e) => setKata(e.target.value)}
        />
        {kata !== '' && (
          <button type="button" className="btn btn--ghost btn--sm" onClick={() => setKata('')}>
            Clear
          </button>
        )}
        <span className="toolbar__spacer" />
      </div>

      {memuat && <Memuat />}
      {/* PitaGalat (11-10-2026): tempat galat yang sama dengan layar Claim Life lain, di atas isi antrian. */}
      <PitaGalat galat={galat} onTutup={() => setGalat(null)} />

      {!memuat && galat === null && hal !== null && (
        <>
          {hal.baris.length === 0 ? (
            cari === '' ? (
              <Kosong pesan="No cases in this queue." />
            ) : (
              <Kosong pesan="No cases match the search." />
            )
          ) : (
            <div className="inbox__gulir">
            <table className="inbox__tabel claimlife__tabel">
              <thead>
                <tr>
                  <th>{KOLOM.caseId}</th>
                  <th>{KOLOM.nomorKlaim}</th>
                  <th>{KOLOM.plNumber}</th>
                  <th>{KOLOM.createOpName}</th>
                  <th>{KOLOM.updateOpName}</th>
                  <th>{KOLOM.tglCreate}</th>
                  <th>{KOLOM.status}</th>
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
                    {/* CLAIM_NO kosong sampai Save to RNM (10-10-2026). */}
                    <td>{b.nomorKlaim === '' ? '—' : b.nomorKlaim}</td>
                    <td>{b.plNumber === '' ? '—' : b.plNumber}</td>
                    <td>{b.createOpName}</td>
                    <td>{b.updateOpName === '' ? '—' : b.updateOpName}</td>
                    <td>{tanggalTampil(b.tglCreate)}</td>
                    <td>{b.status}</td>
                  </tr>
                ))}
              </tbody>
            </table>
            </div>
          )}

          {/* 50 per halaman (592); maksimum 500 (942). Paginasi inti (Previous / Page x of N / Next), seperti inbox lain. */}
          <Halaman halaman={hal.halaman} ukuran={hal.ukuran} total={hal.total} onPindah={setHalaman} />
        </>
      )}
      <TabelKomite peran={peran} onBukaModul={onBukaModul} />
    </section>
  )
}
