// Bentuk unggah berkas modal `ASM Attach Content` — DISALIN dari menu Master
// Product Name Life (`modul/masterproductnamelife/frontend/bentuk.ts`
// `gabungBerkas` / `unggahBerurutan`), permintaan pemakai 8 Oktober 2026:
// "contoh yang ada di menu product name life itu aja ditiru samakan yah krn
// tim saya menggunakan hal tersebut".
//
// ⛔ Disalin, tidak diimpor: `modul/X` hanya boleh mengimpor `inti/**` dan
// `modul/X/**` (`inti/frontend/lapisan.guard.test.ts`).

/**
 * Pilihan berkas digabung: berkas yang namanya (tanpa beda huruf besar/kecil)
 * sudah ada dilewati, urutan pilihan dipertahankan.
 */
export function gabungBerkas(lama: readonly File[], baru: readonly File[]): File[] {
  const hasil = [...lama]
  const ada = new Set(lama.map((f) => f.name.toLowerCase()))
  for (const f of baru) {
    const k = f.name.toLowerCase()
    if (ada.has(k)) continue
    ada.add(k)
    hasil.push(f)
  }
  return hasil
}

/**
 * Unggah SATU PER SATU, berurutan, lewat rute unggah yang ada (Attach per
 * berkas); kegagalan satu berkas dikumpulkan dan tidak menghentikan sisanya.
 * `mulai` dipanggil sebelum tiap berkas (indeks dari 0).
 */
export async function unggahBerurutan(
  berkas: readonly File[],
  kirim: (f: File) => Promise<unknown>,
  mulai?: (f: File, i: number) => void,
): Promise<{ berkas: File; galat: unknown }[]> {
  const gagal: { berkas: File; galat: unknown }[] = []
  for (const [i, f] of berkas.entries()) {
    mulai?.(f, i)
    try {
      await kirim(f)
    } catch (galat) {
      gagal.push({ berkas: f, galat })
    }
  }
  return gagal
}
