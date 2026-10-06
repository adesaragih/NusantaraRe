// ⛔ DIPINDAHKAN dari `pages/FormKontrakTreatyIn.tsx` 5 Oktober 2026 —
// pemindahan MURNI, nol perubahan perilaku.

import { FieldTanggal } from '../../../../inti/frontend/components/ui/dasar'

/**
 * Medan yang KUNCINYA tidak ada di dokumen warisan kontrak ini.
 *
 * ⛔ MATI dengan keterangan, bukan kotak kosong. Kotak kosong terbaca
 * "belum diisi"; medan mati terbaca "tidak ada di sistem lama". Bedanya
 * menentukan apa yang orang tagih — dan sapuan 3 Oktober 2026 menemukan ia
 * BUKAN kasus langka: `ContractRefNo` tidak ada di 1.112 dari 1.854 kontrak,
 * `TreatyLeader` di 1.195.
 *
 * Pola yang sama sudah dipakai tombol `Choose Ceding`, dan sebabnya tertulis
 * di sana.
 */
/**
 * Pembungkus `FieldTanggal` yang MENYEMBUNYIKAN `dd/mm/yyyy` saat kosong.
 *
 * ⛔ Masalahnya nyata: `<input type="date">` yang kosong menuliskan
 * `dd/mm/yyyy` sendiri, dan di layar ia terbaca seperti nilai — kontrak
 * tanpa tanggal terlihat seolah bertanggal.
 *
 * ⚠️ `FieldTanggal` tinggal di `inti/frontend/components/ui/dasar.tsx`, dan
 * berkas itu DIPAKAI BERSAMA — ia tidak disunting. Jadi kaitnya dipasang di
 * sini: satu `<span>` yang menambahkan kelas saat nilainya kosong, lalu
 * aturan CSS modul yang mewarnai `::-webkit-datetime-edit` transparan
 * selama kosong DAN tidak difokus.
 *
 * ⛔ CSS SAJA TIDAK CUKUP, dan itu terukur bukan dikira: tidak ada pemilih
 * yang membedakan `<input type="date">` kosong dari yang terisi.
 * `:placeholder-shown` tidak cocok untuk medan tanggal, medan tanggal
 * kosong tanpa `required` itu `:valid`, dan React menyetel `value` sebagai
 * PROPERTI sehingga `[value='']` tidak cocok. Kaitnya harus datang dari
 * markup, dan markup yang boleh disunting adalah yang ini.
 *
 * ⚠️ BATASNYA DINYATAKAN: `::-webkit-datetime-edit` hanya ada di
 * Chromium/WebKit. Di Firefox `dd/mm/yyyy` tetap tampak. Tidak ada padanan
 * standarnya, dan menggantinya dengan `type="text"` akan membuang pemilih
 * tanggal bawaan — harga yang lebih mahal daripada masalahnya.
 */
export default function TanggalRedup({
  label,
  value,
  onChange,
}: {
  label: string
  value: string
  onChange: (v: string) => void
}) {
  return (
    <span className={'trin__tanggal' + (value === '' ? ' trin__tanggal--kosong' : '')}>
      <FieldTanggal label={label} value={value} onChange={onChange} />
    </span>
  )
}
