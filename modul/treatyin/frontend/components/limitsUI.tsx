// Perangkat tampilan tab Limits — dipakai cabang Prop dan Non-Prop.
//
// ⭐ Permintaan pemakai 6 Oktober 2026: "terapkan UI/UX hebat … jangan hanya
// panah". Yang dibangun di sini:
//   - kepala bagian: judul + lencana jumlah + tombol Add berikon;
//   - kartu lipat: nomor urut, judul (atau "belum dipilih"), ringkasan
//     nilai sebagai chip, tombol Detail berchevron yang berputar, Delete
//     berikon di sisi kanan;
//   - bagian dalam rincian berjudul kecil, supaya medan yang berkaitan
//     terbaca sebagai satu kelompok.
//
// ⛔ Isi kartu yang tertutup TETAP di DOM (`hidden`), seperti `<details>`
// sebelumnya: keadaan di dalamnya (sub-tab aktif, isian) tidak hilang saat
// dilipat, dan pembaca layar/uji statis melihat strukturnya utuh.
//
// ⛔ Warna hanya dari token `inti` (`--primary`, `--border`, `--surface-2`…),
// jadi mode gelap ikut tanpa satu baris tambahan.

import { useId, useState, type ReactNode } from 'react'

import { IkonChevron } from '../../../../inti/frontend/components/ui/dasar'
import { TombolNavigasi } from './navigasi'

/** Sifat bersama ikon garis 24×24 — sama dengan `dasar.tsx`. */
const sifat = (ukuran: number) => ({
  width: ukuran,
  height: ukuran,
  viewBox: '0 0 24 24',
  fill: 'none',
  stroke: 'currentColor',
  strokeWidth: 2,
  strokeLinecap: 'round' as const,
  strokeLinejoin: 'round' as const,
  'aria-hidden': true,
})

/** Chevron ke kanan — kartu tertutup. */
export function IkonChevronKanan({ ukuran = 16 }: { ukuran?: number }) {
  return (
    <svg {...sifat(ukuran)} strokeWidth={1.8}>
      <path d="m9 6 6 6-6 6" />
    </svg>
  )
}

export function IkonTambah({ ukuran = 14 }: { ukuran?: number }) {
  return (
    <svg {...sifat(ukuran)}>
      <path d="M12 5v14M5 12h14" />
    </svg>
  )
}

export function IkonHapus({ ukuran = 14 }: { ukuran?: number }) {
  return (
    <svg {...sifat(ukuran)}>
      <path d="M4 7h16M10 11v6M14 11v6M6 7l1 12a2 2 0 0 0 2 2h6a2 2 0 0 0 2-2l1-12M9 7V4h6v3" />
    </svg>
  )
}

/**
 * Tombol Add — label teks ekspor (`Add`, `add Layer`, …).
 *
 * ⛔ 8 Oktober 2026 — TANPA ikon dan tanpa rupa kapsul merah: tombol polos
 * seperti Pega ("jangan ada design tambahan"). Kelas `tl-tambah` tetap
 * sebagai penanda (uji dan CSS lama), rupanya disamakan `btn btn--sm`.
 */
export function TombolTambah({ label, onClick, disabled }: { label: string; onClick: () => void; disabled?: boolean }) {
  return (
    <button type="button" className="btn btn--sm tl-tambah" onClick={onClick} disabled={disabled}>
      {label}
    </button>
  )
}

/**
 * Tombol Delete/Remove — `labelAkses` menyebut APA yang dihapus. ⛔ Tanpa
 * ikon (8 Oktober 2026), seperti `TombolTambah`.
 */
export function TombolHapus({
  label,
  labelAkses,
  onClick,
  disabled,
}: {
  label: string
  labelAkses?: string
  onClick: () => void
  disabled?: boolean
}) {
  return (
    <button type="button" className="btn btn--sm tl-hapus" aria-label={labelAkses} onClick={onClick} disabled={disabled}>
      {label}
    </button>
  )
}

/** Kepala bagian: judul + lencana jumlah + aksi di kanan. */
export function KepalaBagian({ judul, jumlah, aksi }: { judul: string; jumlah?: number; aksi?: ReactNode }) {
  return (
    <div className="tl-kepala">
      <span className="tl-kepala__judul">
        {judul}
        {jumlah !== undefined && <span className="tl-lencana">{jumlah}</span>}
      </span>
      {aksi}
    </div>
  )
}

/** Satu nilai ringkas di kepala kartu. */
export function Chip({ label, nilai }: { label: string; nilai: string }) {
  if (nilai === '') return null
  return (
    <span className="tl-chip">
      <span className="tl-chip__label">{label}</span>
      <span className="tl-chip__nilai">{nilai}</span>
    </span>
  )
}

/**
 * Kartu lipat — satu baris grid `expandPane` Pega.
 *
 * ⛔ Tombol buka-tutup dan tombol Delete adalah DUA tombol bersaudara, bukan
 * bersarang: tombol di dalam tombol tidak sah dan pembaca layar tersesat.
 */
export function KartuLipat({
  nomor,
  judul,
  judulKosong,
  meta,
  bukaAwal = false,
  aksi,
  labelDetail = 'Detail',
  children,
}: {
  nomor?: number
  judul: string
  /** Teks bila judul kosong — baris baru yang belum dipilih. */
  judulKosong: string
  meta?: ReactNode
  bukaAwal?: boolean
  /** Tombol di kanan kepala (Delete). */
  aksi?: ReactNode
  labelDetail?: string
  children: ReactNode
}) {
  const [buka, setBuka] = useState(bukaAwal)
  const idIsi = useId()
  return (
    <section className={'tl-kartu' + (buka ? ' tl-kartu--buka' : '')}>
      <div className="tl-kartu__kepala">
        {/* ⭐ Kontrol NAVIGASI, bukan `<button>` — tetap dapat dibuka di mode
            lihat (`fieldset disabled`). Lihat `navigasi.tsx`. */}
        <TombolNavigasi
          className="tl-kartu__lipat"
          terbuka={buka}
          kendali={idIsi}
          onKlik={() => {
            setBuka(!buka)
          }}
        >
          <span className="tl-kartu__chevron">{buka ? <IkonChevron /> : <IkonChevronKanan />}</span>
          {nomor !== undefined && <span className="tl-nomor">{nomor}</span>}
          <span className={'tl-kartu__judul' + (judul === '' ? ' tl-kartu__judul--kosong' : '')}>{judul === '' ? judulKosong : judul}</span>
          <span className="tl-kartu__petunjuk">{labelDetail}</span>
        </TombolNavigasi>
        {meta !== undefined && <div className="tl-kartu__meta">{meta}</div>}
        {aksi !== undefined && <div className="tl-kartu__aksi">{aksi}</div>}
      </div>
      <div className="tl-kartu__isi" id={idIsi} hidden={!buka}>
        {children}
      </div>
    </section>
  )
}

/** Bagian di dalam rincian — judul kecil + isi. */
export function Bagian({ judul, aksi, children }: { judul: string; aksi?: ReactNode; children: ReactNode }) {
  return (
    <div className="tl-bagian">
      <div className="tl-bagian__kepala">
        <h5 className="tl-bagian__judul">{judul}</h5>
        {aksi}
      </div>
      {children}
    </div>
  )
}
