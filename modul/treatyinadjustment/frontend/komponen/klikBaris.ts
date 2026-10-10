// Klik baris grid membuka rincian — SALINAN `klikBaris` modul Treaty In
// (`modul/treatyin/frontend/components/navigasi.tsx`); modul tidak boleh
// saling impor.
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
