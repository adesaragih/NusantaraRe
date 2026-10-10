// Kontrol NAVIGASI di dalam tab — buka/tutup rincian baris dan pindah sub-tab.
//
// ⛔ MENGAPA BUKAN `<button>`. Mode lihat mengunci isi tab dengan
// `<fieldset disabled>` (`FormKontrakTreatyIn.tsx`), dan HTML mematikan SETIAP
// kontrol form di dalamnya — termasuk tombol yang hanya MELIHAT. Laporan
// pemakai 8 Oktober 2026: *"kenapa di treaty in saat view malah tidak bisa
// klik down colom contoh yg di limits"*. Di Pega rincian baris dan sub-tab
// tetap dapat dibuka di mode lihat; yang dikunci hanya isian dan tombol ubah.
//
// Kontrol di sini elemen BUKAN-form (`role="button"` / `role="tab"`): tidak
// terjangkau `fieldset disabled`, tetap dapat difokus (`tabIndex 0`), dan
// ditekan dengan klik, Enter, atau Spasi. Kelas dan atribut ARIA-nya sama
// dengan tombol yang digantikannya, jadi rupanya tidak berubah.

import type { KeyboardEvent, ReactNode } from 'react'

function tekanTombol(e: KeyboardEvent, aksi: () => void): void {
  if (e.key === 'Enter' || e.key === ' ') {
    e.preventDefault()
    aksi()
  }
}

/** Satu kontrol navigasi — buka/tutup rincian baris, kartu layer, dsb. */
export function TombolNavigasi({
  className,
  label,
  terbuka,
  kendali,
  onKlik,
  children,
}: {
  className: string
  /** `aria-label` — bila isi tombol hanya ikon. */
  label?: string
  /** `aria-expanded` — keadaan rincian yang dikendalikannya. */
  terbuka?: boolean
  /** `aria-controls` — pengenal isi yang dikendalikannya. */
  kendali?: string
  onKlik: () => void
  children: ReactNode
}) {
  return (
    <span
      role="button"
      tabIndex={0}
      className={className}
      aria-label={label}
      aria-expanded={terbuka}
      aria-controls={kendali}
      onClick={onKlik}
      onKeyDown={(e) => {
        tekanTombol(e, onKlik)
      }}
    >
      {children}
    </span>
  )
}

/**
 * Strip SUB-TAB di dalam tab — markup SAMA dengan `StripTab` inti
 * (`.tabs` / `.tabs__item` / `--aktif`, `role="tablist"`/`"tab"`,
 * `aria-selected`), tetapi butirnya bukan `<button>`: sub-tab tetap dapat
 * dipindah di mode lihat. ⛔ Strip tab UTAMA form tetap `StripTab` inti — ia
 * berdiri di LUAR `fieldset`.
 */
export function StripTabNavigasi<T extends string>({
  tab,
  aktif,
  onPilih,
}: {
  tab: readonly T[]
  aktif: T
  onPilih: (t: T) => void
}) {
  return (
    <div className="tabs" role="tablist">
      {tab.map((t) => (
        <span
          key={t}
          role="tab"
          tabIndex={0}
          aria-selected={t === aktif}
          className={'tabs__item' + (t === aktif ? ' tabs__item--aktif' : '')}
          onClick={() => {
            onPilih(t)
          }}
          onKeyDown={(e) => {
            tekanTombol(e, () => {
              onPilih(t)
            })
          }}
        >
          {t}
        </span>
      ))}
    </div>
  )
}

/**
 * Pemilih unsur INTERAKTIF di dalam baris — klik di sana TIDAK membuka /
 * menutup rincian: kotak isian, tombol, dropdown, tautan, dan panah ▸/▾
 * sendiri (ber-`role="button"`, sudah menangani kliknya).
 */
const UNSUR_INTERAKTIF =
  'a,button,input,select,textarea,label,[role="button"],[role="combobox"],[role="listbox"],[role="option"],[contenteditable="true"]'

/**
 * ⭐ Klik BARIS membuka/menutup rinciannya — permintaan pemakai 9 Oktober
 * 2026: *"ketika tekan barisnya itu anaknya langsung muncul juga tidak harus
 * menekan panahnya"*. Klik pada unsur interaktif dan teks yang sedang
 * diseleksi diabaikan.
 */
export function klikBaris(e: { target: EventTarget | null }, bolak: () => void): void {
  const t = e.target
  if (typeof Element !== 'undefined' && t instanceof Element && t.closest(UNSUR_INTERAKTIF) !== null) return
  if (typeof window !== 'undefined' && (window.getSelection()?.toString() ?? '') !== '') return
  bolak()
}
